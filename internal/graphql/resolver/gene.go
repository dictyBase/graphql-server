package resolver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dictyBase/graphql-server/internal/graphql/cache"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/fetch"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/repository"
)

const (
	geneHash         = "GENE2NAME/geneids"
	goHash           = "GO2NAME/goids"
	uniprotHash      = "UNIPROT2NAME/uniprot"
	geneUniprotHash  = "GENE2UNIPROT/gene"
	baseGoaURLPrefix = "https://www.ebi.ac.uk/QuickGO/services/annotation/search"
	baseGoaURLParams = "?includeFields=goName&limit=100&geneProductId="
)

var baseGoaURL = fmt.Sprintf("%s%s", baseGoaURLPrefix, baseGoaURLParams)

var dbHashMap = map[string]string{
	"dictyBase": geneHash,
	"GO":        goHash,
	"UniProtKB": uniprotHash,
}

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
	// TargetSets    []string    `json:"targetSets"`
	// Synonyms      []string    `json:"synonyms"`
	// Name          string      `json:"name"`
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

func fetchAndUnmarshalJSON(url string, target interface{}) error {
	res, err := fetch.GetResp(url)
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

func getValFromHash(
	hash, key string,
	cache repository.Repository,
) (string, bool, error) {
	exists, err := cache.HExists(hash, key)
	if err != nil {
		return "", false, fmt.Errorf(
			"error checking hash existence for %s/%s: %w",
			hash,
			key,
			err,
		)
	}
	if !exists {
		return "", false, nil
	}
	name, err := cache.HGet(hash, key)
	if err != nil {
		return "", true, fmt.Errorf(
			"error getting value from hash %s/%s: %w",
			hash,
			key,
			err,
		)
	}
	return name, true, nil
}

func getNameFromDB(db, id string, cache repository.Repository) string {
	hash, ok := dbHashMap[db]
	if !ok {
		return ""
	}

	key := id
	if db == "GO" {
		key = fmt.Sprintf("%s:%s", db, id)
	}

	val, _, err := getValFromHash(hash, key, cache)
	if err != nil {
		// Since the original function returned empty string on error,
		// we'll maintain that behavior but could log the error if needed
		return ""
	}
	return val
}

func getWith(with []with, repo repository.Repository) []*models.With {
	wm := []*models.With{}
	for _, v := range with {
		for _, w := range v.ConnectedXRefs {
			wm = append(wm, &models.With{
				ID:   w.ID,
				Db:   w.DB,
				Name: getNameFromDB(w.DB, w.ID, repo),
			})
		}
	}
	return wm
}

func getExtensions(
	extensions []extension,
	repo repository.Repository,
) []*models.Extension {
	ext := []*models.Extension{}
	for _, v := range extensions {
		for _, e := range v.ConnectedXRefs {
			ext = append(ext, &models.Extension{
				ID:       e.ID,
				Db:       e.DB,
				Relation: e.Relation,
				Name:     getNameFromDB(e.DB, e.ID, repo),
			})
		}
	}
	return ext
}

func (qrs *QueryResolver) GeneOntologyAnnotation(
	ctx context.Context,
	gene string,
) ([]*models.GOAnnotation, error) {
	redis := qrs.GetRedisRepository(cache.RedisKey)
	uniprotID, exists, err := getValFromHash(geneUniprotHash, gene, redis)
	if err != nil {
		return nil, fmt.Errorf(
			"error getting UniProt ID for gene %s: %w",
			gene,
			err,
		)
	}
	if !exists {
		return nil, fmt.Errorf("no UniProt ID found for gene %s", gene)
	}

	url := fmt.Sprintf("%s%s", baseGoaURL, uniprotID)
	geneOntology, err := fetchGOAs(url)
	if err != nil {
		return nil, fmt.Errorf(
			"error fetching gene ontology annotations: %w",
			err,
		)
	}

	annotations := make([]*models.GOAnnotation, 0, len(geneOntology.Results))
	for _, result := range geneOntology.Results {
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

		if result.WithFrom != nil {
			annotation.With = getWith(result.WithFrom, redis)
		}
		if result.Extensions != nil {
			annotation.Extensions = getExtensions(result.Extensions, redis)
		}

		annotations = append(annotations, annotation)
	}

	return annotations, nil
}

func (qrs *QueryResolver) Gene(
	ctx context.Context,
	geneID string,
) (*models.Gene, error) {
	redis := qrs.GetRedisRepository(cache.RedisKey)
	gene, err := cache.GetGeneFromCache(ctx, redis, geneID)
	if err != nil {
		qrs.Logger.WithError(err).Error("Failed to get gene from cache")
		errorutils.AddGQLError(ctx, err)
		return nil, fmt.Errorf("error retrieving gene information: %w", err)
	}
	return gene, nil
}

func (qrs *QueryResolver) AllOrthologs(
	ctx context.Context,
	gene string,
) (*models.Gene, error) {
	return &models.Gene{}, nil
}

func (qrs *QueryResolver) GeneralInformation(
	ctx context.Context,
	gene string,
) (*models.Gene, error) {
	return &models.Gene{}, nil
}

// GetAssociatedSequnces is the resolver for the getAssociatedSequnces field.
func (qrs *QueryResolver) GetAssociatedSequnces(
	ctx context.Context,
	gene string,
) (*models.Gene, error) {
	return &models.Gene{}, nil
}

// GetLinks is the resolver for the getLinks field.
func (qrs *QueryResolver) GetLinks(
	ctx context.Context,
	gene string,
) (*models.Gene, error) {
	return &models.Gene{}, nil
}

// GetProteinInformation is the resolver for the getProteinInformation field.
func (qrs *QueryResolver) GetProteinInformation(
	ctx context.Context,
	gene string,
) (*models.Gene, error) {
	return &models.Gene{}, nil
}

// ListGeneProductInfo is the resolver for the listGeneProductInfo field.
func (qrs *QueryResolver) ListGeneProductInfo(
	ctx context.Context,
	gene string,
) (*models.Gene, error) {
	return &models.Gene{}, nil
}

// ListRecentGenes is the resolver for the listRecentGenes field.
func (qrs *QueryResolver) ListRecentGenes(
	ctx context.Context,
	limit int,
) ([]*models.Gene, error) {
	return []*models.Gene{}, nil
}

// GeneGeneralInformation is the resolver for the geneGeneralInformation field.
func (qrs *QueryResolver) GeneGeneralInformation(
	ctx context.Context,
	gene string,
) (*models.GeneGeneralInfo, error) {
	return &models.GeneGeneralInfo{}, nil
}
