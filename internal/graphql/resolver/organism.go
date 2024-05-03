package resolver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dictyBase/graphql-server/internal/graphql/fetch"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
)

type organisms struct {
	Data []organismData `json:"data"`
}

type organismData struct {
	Type       string             `json:"type"`
	ID         string             `json:"id"`
	Attributes organismAttributes `json:"attributes"`
}

type organismAttributes struct {
	TaxonID        string     `json:"taxon_id"`
	ScientificName string     `json:"scientific_name"`
	Citations      []citation `json:"citations"`
}

type citation struct {
	Authors string `json:"authors"`
	Title   string `json:"title"`
	Journal string `json:"journal"`
	Link    string `json:"link"`
}

func fetchOrganisms(ctx context.Context, url string) (*organisms, error) {
	o := new(organisms)
	res, err := fetch.GetResp(url)
	if err != nil {
		return o, err
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(o); err != nil {
		return o, fmt.Errorf("error in decoding json %s", err)
	}
	return o, nil
}

func (qrs *QueryResolver) Organism(
	ctx context.Context,
	taxonID string,
) (*models.Organism, error) {
	o := &models.Organism{}
	c := []*models.Citation{}
	url := qrs.GetAPIEndpoint("organism")
	d, err := fetchOrganisms(ctx, url)
	if err != nil {
		return o, err
	}
	for _, val := range d.Data {
		if val.ID == taxonID {
			o.TaxonID = val.Attributes.TaxonID
			o.ScientificName = val.Attributes.ScientificName
			for _, ci := range val.Attributes.Citations {
				c = append(c, &models.Citation{
					Authors:  ci.Authors,
					Journal:  ci.Journal,
					PubmedID: ci.Link[len(ci.Link)-8:], // just get ID from URL
					Title:    ci.Title,
				})
			}
			o.Citations = c
		}
	}
	return o, nil
}

func (qrs *QueryResolver) ListOrganisms(
	ctx context.Context,
) ([]*models.Organism, error) {
	orgs := []*models.Organism{}
	url := qrs.GetAPIEndpoint("organism")
	d, err := fetchOrganisms(ctx, url)
	if err != nil {
		return orgs, err
	}
	for _, val := range d.Data {
		c := []*models.Citation{}
		for _, ci := range val.Attributes.Citations {
			c = append(c, &models.Citation{
				Authors:  ci.Authors,
				Journal:  ci.Journal,
				PubmedID: ci.Link[len(ci.Link)-8:], // just get ID from URL
				Title:    ci.Title,
			})
		}
		orgs = append(orgs, &models.Organism{
			TaxonID:        val.Attributes.TaxonID,
			ScientificName: val.Attributes.ScientificName,
			Citations:      c,
		})
	}
	return orgs, nil
}
