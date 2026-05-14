package providers

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/FerdaneOgut/jobradar/scraper-service/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WeWorkRemotelyProvider struct{}

type wwrRSS struct {
	Channel wwrChannel `xml:"channel"`
}

type wwrChannel struct {
	Items []wwrItem `xml:"item"`
}

type wwrItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Region      string `xml:"region"`
	Type        string `xml:"type"`
	Company     string `xml:"company"`
	Description string `xml:"description"`
}

func (p *WeWorkRemotelyProvider) Name() string {
	return "weworkremotely"
}

func extractCompanyAndTitle(raw string) (string, string) {
	parts := strings.SplitN(raw, ":", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
	return "", raw
}

func (p *WeWorkRemotelyProvider) FetchJobs() ([]models.Job, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get("https://weworkremotely.com/remote-jobs.rss")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rss wwrRSS
	if err := xml.NewDecoder(resp.Body).Decode(&rss); err != nil {
		return nil, err
	}

	var jobs []models.Job
	for i, item := range rss.Channel.Items {
		company, title := extractCompanyAndTitle(item.Title)

		jobs = append(jobs, models.Job{
			ID:          primitive.NewObjectID(),
			ExternalID:  fmt.Sprintf("weworkremotely-%d", i),
			Title:       title,
			Company:     company,
			Location:    item.Region,
			Description: item.Description,
			URL:         item.Link,
			Source:      "weworkremotely",
			IsRemote:    true,
			CreatedAt:   time.Now(),
		})
	}

	return jobs, nil
}
