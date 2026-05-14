package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/FerdaneOgut/jobradar/scraper-service/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RemotiveProvider struct{}

type remotiveResponse struct {
	Jobs []remotiveJob `json:"jobs"`
}

type remotiveJob struct {
	ID            int      `json:"id"`
	URL           string   `json:"url"`
	Title         string   `json:"title"`
	CompanyName   string   `json:"company_name"`
	CandidateCity string   `json:"candidate_required_location"`
	Description   string   `json:"description"`
	Tags          []string `json:"tags"`
	Salary        string   `json:"salary"`
}

func (p *RemotiveProvider) Name() string {
	return "remotive"
}

func (p *RemotiveProvider) FetchJobs() ([]models.Job, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get("https://remotive.com/api/remote-jobs?limit=50")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result remotiveResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var jobs []models.Job
	for _, rj := range result.Jobs {
		jobs = append(jobs, models.Job{
			ID:          primitive.NewObjectID(),
			ExternalID:  fmt.Sprintf("remotive-%d", rj.ID),
			Title:       rj.Title,
			Company:     rj.CompanyName,
			Location:    rj.CandidateCity,
			Description: rj.Description,
			URL:         rj.URL,
			Source:      "remotive",
			Tags:        rj.Tags,
			IsRemote:    true,
			Salary:      rj.Salary,
			CreatedAt:   time.Now(),
		})
	}

	return jobs, nil
}
