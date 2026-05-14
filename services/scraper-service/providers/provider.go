package providers

import "github.com/FerdaneOgut/jobradar/scraper-service/models"

type JobProvider interface {
	FetchJobs() ([]models.Job, error)
	Name() string
}
