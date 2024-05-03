package organism

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dictyBase/graphql-server/internal/graphql/fetch"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/sirupsen/logrus"
)

type OrganismResolver struct {
	Logger       *logrus.Entry
	DownloadsURL string
}

type downloads struct {
	Data []downloadsData `json:"data"`
}

type downloadsData struct {
	Type       string              `json:"type"`
	ID         string              `json:"id"`
	Attributes downloadsAttributes `json:"attributes"`
}

type downloadsAttributes struct {
	Title string          `json:"title"`
	Items []downloadsItem `json:"items"`
}

type downloadsItem struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func fetchDownloads(url string) (*downloads, error) {
	d := new(downloads)
	res, err := fetch.GetResp(url)
	if err != nil {
		return d, err
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(d); err != nil {
		return d, fmt.Errorf("error in decoding json %s", err)
	}
	return d, nil
}

func (r *OrganismResolver) Downloads(
	ctx context.Context,
	obj *models.Organism,
) ([]*models.Download, error) {
	ds := []*models.Download{}
	if r.DownloadsURL == "" {
		r.DownloadsURL = fmt.Sprintf(
			"https://raw.githubusercontent.com/dictyBase/migration-data/master/downloads/%s.staging.json",
			obj.TaxonID,
		)
	}
	res, err := fetchDownloads(r.DownloadsURL)
	if err != nil {
		return ds, err
	}
	for _, val := range res.Data {
		items := []*models.DownloadItem{}
		for _, item := range val.Attributes.Items {
			items = append(items, &models.DownloadItem{
				Title: item.Title,
				URL:   item.URL,
			})
		}
		ds = append(ds, &models.Download{
			Title: val.Attributes.Title,
			Items: items,
		})
	}
	return ds, nil
}
