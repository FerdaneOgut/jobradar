package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/FerdaneOgut/jobradar/scraper-service/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ArbeitnowProvider struct{}

type arbeitnowResponse struct {
	Data []arbeitnowJob `json:"data"`
}

type arbeitnowJob struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	CompanyName string   `json:"company_name"`
	Location    string   `json:"location"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Tags        []string `json:"tags"`
	Remote      bool     `json:"remote"`
}

func (p *ArbeitnowProvider) Name() string {
	return "arbeitnow"
}

func (p *ArbeitnowProvider) FetchJobs() ([]models.Job, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get("https://www.arbeitnow.com/api/job-board-api?page=1")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result arbeitnowResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var jobs []models.Job
	for _, aj := range result.Data {
		jobs = append(jobs, models.Job{
			ID:          primitive.NewObjectID(),
			ExternalID:  fmt.Sprintf("arbeitnow-%s", aj.Slug),
			Title:       aj.Title,
			Company:     aj.CompanyName,
			Location:    aj.Location,
			Description: aj.Description,
			URL:         aj.URL,
			Source:      "arbeitnow",
			Tags:        aj.Tags,
			IsRemote:    aj.Remote,
			CreatedAt:   time.Now(),
		})
	}

	return jobs, nil
}
