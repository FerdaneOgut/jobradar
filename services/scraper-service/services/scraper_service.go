package services

import (
	"context"
	"log"

	"github.com/FerdaneOgut/jobradar/scraper-service/providers"
	"github.com/FerdaneOgut/jobradar/scraper-service/repository"
)

type ScraperService struct {
	repo      *repository.JobRepository
	providers []providers.JobProvider
}

func NewScraperService(repo *repository.JobRepository) *ScraperService {
	return &ScraperService{
		repo: repo,
		providers: []providers.JobProvider{
			&providers.RemoteOKProvider{},
			&providers.RemotiveProvider{},
			&providers.ArbeitnowProvider{},
			&providers.WeWorkRemotelyProvider{},
		},
	}
}

func (s *ScraperService) ScrapeAll() {
	ctx := context.Background()

	for _, provider := range s.providers {
		log.Printf("Scraping from %s...", provider.Name())

		jobs, err := provider.FetchJobs()
		if err != nil {
			log.Printf("Error scraping %s: %v", provider.Name(), err)
			continue
		}

		saved := 0
		skipped := 0

		for _, job := range jobs {
			exists, err := s.repo.IsExist(ctx, job.ExternalID)
			if err != nil {
				log.Printf("Error checking job existence: %v", err)
				continue
			}

			if exists {
				skipped++
				continue
			}

			if err := s.repo.Save(ctx, job); err != nil {
				log.Printf("Error saving job: %v", err)
				continue
			}

			saved++
		}

		log.Printf("%s: saved %d, skipped %d", provider.Name(), saved, skipped)
	}
}
