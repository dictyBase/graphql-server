package fetch

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/repository"
	"github.com/sirupsen/logrus"
)

const (
	baseGoaURLPrefix = "https://www.ebi.ac.uk/QuickGO/services/annotation/search"
	baseGoaURLParams = "?includeFields=goName&limit=100&geneProductId="
)

var baseGoaURL = fmt.Sprintf("%s%s", baseGoaURLPrefix, baseGoaURLParams)

type quickGo struct {
	NumberOfHits int      `json:"numberOfHits"`
	Results      []result `json:"results"`
	PageInfo     pageInfo `json:"pageInfo"`
}

type result struct {
	ID            string      `json:"id"`
	GeneProductID string      `json:"geneProductId"`
	Qualifier     string      `json:"qualifier"`
	GoID          string      `json:"goId"`
	GoName        string      `json:"goName"`
	GoEvidence    string      `json:"goEvidence"`
	GoAspect      string      `json:"goAspect"`
	EvidenceCode  string      `json:"evidenceCode"`
	Reference     string      `json:"reference"`
	WithFrom      []with      `json:"withFrom"`
	TaxonID       int         `json:"taxonId"`
	TaxonName     string      `json:"taxonName"`
	AssignedBy    string      `json:"assignedBy"`
	Extensions    []extension `json:"extensions"`
	Symbol        string      `json:"symbol"`
	Date          string      `json:"date"`
}

type with struct {
	ConnectedXRefs []withXRef `json:"connectedXrefs"`
}

type extension struct {
	ConnectedXRefs []extensionXRef `json:"connectedXrefs"`
}

type withXRef struct {
	DB string `json:"db"`
	ID string `json:"id"`
}

type extensionXRef struct {
	DB       string `json:"db"`
	ID       string `json:"id"`
	Relation string `json:"relation"`
}

type pageInfo struct {
	ResultsPerPage int `json:"resultsPerPage"`
	Current        int `json:"current"`
	Total          int `json:"total"`
}

type FetchAndBuildAnnotationsParams struct {
	Ctx       context.Context
	Gene      string
	UniprotID string
	Redis     repository.Repository
	Logger    *logrus.Entry
}

func fetchAndUnmarshalJSON(url string, target interface{}) error {
	res, err := GetResp(url)
	if err != nil {
		return fmt.Errorf("error fetching data: %w", err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return fmt.Errorf("error decoding JSON: %w", err)
	}
	return nil
}

func fetchGOAs(url string) (*quickGo, error) {
	goa := new(quickGo)
	err := fetchAndUnmarshalJSON(url, goa)
	return goa, err
}

func FetchAndBuildAnnotations(
	params *FetchAndBuildAnnotationsParams,
) ([]*models.GOAnnotation, error) {
	url := fmt.Sprintf("%s%s", baseGoaURL, params.UniprotID)
	geneOntology, err := fetchGOAs(url)
	if err != nil {
		params.Logger.WithFields(logrus.Fields{
			"url":   url,
			"error": err,
		}).Error("failed to fetch gene ontology annotations")
		errorutils.AddGQLError(params.Ctx, err)
		return nil, fmt.Errorf(
			"error fetching gene ontology annotations: %w",
			err,
		)
	}

	annotations := make([]*models.GOAnnotation, 0, len(geneOntology.Results))
	for _, result := range geneOntology.Results {
		annotations = append(annotations, buildGOAnnotation(result))
	}
	return annotations, nil
}

func buildGOAnnotation(result result) *models.GOAnnotation {
	annotation := &models.GOAnnotation{
		ID:           result.ID,
		Type:         result.GoAspect,
		Date:         result.Date,
		EvidenceCode: result.GoEvidence,
		GoTerm:       result.GoName,
		Qualifier:    result.Qualifier,
		Publication:  result.Reference,
		AssignedBy:   result.AssignedBy,
	}

	if len(result.WithFrom) > 0 {
		annotation.With = getWith(result.WithFrom)
	}
	if len(result.Extensions) > 0 {
		annotation.Extensions = getExtensions(result.Extensions)
	}

	return annotation
}

func getWith(
	with []with,
) []*models.With {
	wm := []*models.With{}
	for _, v := range with {
		for _, w := range v.ConnectedXRefs {
			wm = append(wm, &models.With{
				ID: w.ID,
				Db: w.DB,
			})
		}
	}
	return wm
}

func getExtensions(
	extensions []extension,
) []*models.Extension {
	ext := []*models.Extension{}
	for _, v := range extensions {
		for _, e := range v.ConnectedXRefs {
			ext = append(ext, &models.Extension{
				ID:       e.ID,
				Db:       e.DB,
				Relation: e.Relation,
			})
		}
	}
	return ext
}
