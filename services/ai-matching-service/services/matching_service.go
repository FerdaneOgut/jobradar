package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/FerdaneOgut/jobradar/ai-matching-service/models"
	"github.com/FerdaneOgut/jobradar/ai-matching-service/providers"
	"github.com/FerdaneOgut/jobradar/ai-matching-service/repository"
	"github.com/ledongthuc/pdf"
	"github.com/nguyenthenguyen/docx"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MatchingService struct {
	repo        *repository.MatchRepository
	llmProvider providers.LLMProvider
}

func NewMatchingService(repo *repository.MatchRepository, llmProvider providers.LLMProvider) *MatchingService {
	return &MatchingService{
		repo:        repo,
		llmProvider: llmProvider,
	}
}

func (s *MatchingService) MatchAllUsers() {
	ctx := context.Background()

	users, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		return
	}

	jobs, err := s.repo.GetRecentJobs(ctx)
	if err != nil {
		log.Printf("Error fetching jobs: %v", err)
		return
	}

	if len(jobs) == 0 {
		log.Println("No recent jobs to match")
		return
	}

	// Limit to 20 jobs per run to stay within free tier
	if len(jobs) > 20 {
		jobs = jobs[:20]
	}

	log.Printf("Matching %d users against %d jobs", len(users), len(jobs))

	for _, user := range users {
		cvText, err := readCV(user.CvPath)
		if err != nil {
			log.Printf("Error reading CV for user %s: %v", user.ID, err)
			continue
		}

		for _, job := range jobs {
			match, err := s.llmProvider.Match(cvText, job.Title, job.Description)
			if err != nil {
				log.Printf("Error matching job %s: %v", job.ID, err)
				// Wait before retrying next job
				time.Sleep(5 * time.Second)
				continue
			}

			if match.Score >= 50 {
				result := models.MatchResult{
					ID:        bson.NewObjectID(),
					UserID:    user.ID.Hex(),
					JobID:     job.ID.Hex(),
					Job:       job,
					Score:     match.Score,
					Reason:    match.Reason,
					CreatedAt: time.Now(),
				}

				if err := s.repo.SaveMatch(ctx, result); err != nil {
					log.Printf("Error saving match: %v", err)
				} else {
					log.Printf("Match saved: user=%s job=%s score=%d", user.ID.Hex(), job.Title, match.Score)
				}
			}

			// Rate limit: 4 requests per minute to stay safe
			time.Sleep(15 * time.Second)
		}
	}
}

func readCV(path string) (string, error) {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".pdf":
		return readPDF(path)
	case ".docx":
		return readDOCX(path)
	case ".txt":
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(data), nil
	default:
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}
}

func readPDF(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var text string
	for i := 1; i <= r.NumPage(); i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		content, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}
		text += content
	}

	if text == "" {
		return "", fmt.Errorf("could not extract text from PDF")
	}

	return text, nil
}

func readDOCX(path string) (string, error) {
	r, err := docx.ReadDocxFile(path)
	if err != nil {
		return "", err
	}
	defer r.Close()

	return r.Editable().GetContent(), nil
}
