package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/FerdaneOgut/jobradar/scraper-service/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RemoteOKProvider struct{}

type remoteOKJob struct {
	ID          string   `json:"id"`
	Position    string   `json:"position"`
	Company     string   `json:"company"`
	Location    string   `json:"location"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Tags        []string `json:"tags"`
	Salary      string   `json:"salary"`
}

func (p *RemoteOKProvider) Name() string {
	return "remoteok"
}

func (p *RemoteOKProvider) FetchJobs() ([]models.Job, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", "https://remoteok.com/api", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "jobradar/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rawJobs []remoteOKJob
	if err := json.NewDecoder(resp.Body).Decode(&rawJobs); err != nil {
		return nil, err
	}

	var jobs []models.Job
	for _, rj := range rawJobs {
		if rj.ID == "" {
			continue
		}
		jobs = append(jobs, models.Job{
			ID:          primitive.NewObjectID(),
			ExternalID:  fmt.Sprintf("remoteok-%s", rj.ID),
			Title:       rj.Position,
			Company:     rj.Company,
			Location:    rj.Location,
			Description: rj.Description,
			URL:         rj.URL,
			Source:      "remoteok",
			Tags:        rj.Tags,
			IsRemote:    true,
			Salary:      rj.Salary,
			CreatedAt:   time.Now(),
		})
	}

	return jobs, nil
}
