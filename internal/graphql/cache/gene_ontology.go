package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/fetch"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/repository"
	"github.com/sirupsen/logrus"
)

const (
	geneHash               = "GENE2NAME/geneids"
	goHash                 = "GO2NAME/goids"
	uniprotHash            = "UNIPROT2NAME/uniprot"
	geneUniprotHash        = "GENE2UNIPROT/gene"
	baseGoaURLPrefix       = "https://www.ebi.ac.uk/QuickGO/services/annotation/search"
	baseGoaURLParams       = "?includeFields=goName&limit=100&geneProductId="
	geneGOAnnotationPrefix = "gene:goa:"
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

type CachedAnnotationsParams struct {
	Ctx    context.Context
	Gene   string
	Redis  repository.Repository
	Logger *logrus.Entry
}

type FetchAndBuildAnnotationsParams struct {
	Ctx       context.Context
	Gene      string
	UniprotID string
	Redis     repository.Repository
	Logger    *logrus.Entry
}

type CacheAnnotationsParams struct {
	Ctx         context.Context
	Gene        string
	Annotations []*models.GOAnnotation
	Redis       repository.Repository
	Logger      *logrus.Entry
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
	redisRepo repository.Repository,
) (string, bool, error) {
	exists, err := redisRepo.HExists(hash, key)
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
	name, err := redisRepo.HGet(hash, key)
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

func getNameFromDB(
	db, id string,
	redisRepo repository.Repository,
) (string, bool, error) {
	hash, ok := dbHashMap[db]
	if !ok {
		return "", false, nil
	}

	key := id
	if db == "GO" {
		key = fmt.Sprintf("%s:%s", db, id)
	}

	return getValFromHash(hash, key, redisRepo)
}

func getWith(
	with []with,
	repo repository.Repository,
) ([]*models.With, bool, error) {
	wm := []*models.With{}
	hasValidEntries := false
	for _, v := range with {
		for _, w := range v.ConnectedXRefs {
			name, exists, err := getNameFromDB(w.DB, w.ID, repo)
			if err != nil {
				return nil, false, fmt.Errorf(
					"error getting name for %s/%s: %w",
					w.DB,
					w.ID,
					err,
				)
			}
			if !exists {
				continue
			}
			hasValidEntries = true
			wm = append(wm, &models.With{
				ID:   w.ID,
				Db:   w.DB,
				Name: name,
			})
		}
	}
	return wm, hasValidEntries, nil
}

func getExtensions(
	extensions []extension,
	repo repository.Repository,
) ([]*models.Extension, bool, error) {
	ext := []*models.Extension{}
	hasValidEntries := false
	for _, v := range extensions {
		for _, e := range v.ConnectedXRefs {
			name, exists, err := getNameFromDB(e.DB, e.ID, repo)
			if err != nil {
				return nil, false, fmt.Errorf(
					"error getting name for %s/%s: %w",
					e.DB,
					e.ID,
					err,
				)
			}
			if !exists {
				continue
			}
			hasValidEntries = true
			ext = append(ext, &models.Extension{
				ID:       e.ID,
				Db:       e.DB,
				Relation: e.Relation,
				Name:     name,
			})
		}
	}
	return ext, hasValidEntries, nil
}

func GetUniprotIDForGene(
	gene string,
	redis repository.Repository,
) (string, error) {
	uniprotID, exists, err := getValFromHash(geneUniprotHash, gene, redis)
	if err != nil {
		return "", fmt.Errorf(
			"error getting UniProt ID for gene %s: %w",
			gene,
			err,
		)
	}
	if !exists {
		return "", fmt.Errorf("no UniProt ID found for gene %s", gene)
	}
	return uniprotID, nil
}

func getCacheKey(geneID string) string {
	return fmt.Sprintf("%s%s", geneGOAnnotationPrefix, geneID)
}

func GetCachedAnnotations(
	params *CachedAnnotationsParams,
) ([]*models.GOAnnotation, bool, error) {
	cacheKey := getCacheKey(params.Gene)
	exists, err := params.Redis.Exists(cacheKey)
	if err != nil {
		params.Logger.WithFields(logrus.Fields{
			"gene":  params.Gene,
			"error": err,
		}).Error("failed to check cache for gene ontology annotations")
		errorutils.AddGQLError(params.Ctx, err)
		return nil, false, fmt.Errorf("error checking cache: %w", err)
	}
	if !exists {
		return nil, false, nil
	}

	cached, err := params.Redis.Get(cacheKey)
	if err != nil {
		params.Logger.WithFields(logrus.Fields{
			"gene":  params.Gene,
			"error": err,
		}).Error("failed to get cached gene ontology annotations")
		errorutils.AddGQLError(params.Ctx, err)
		return nil, true, fmt.Errorf("error retrieving from cache: %w", err)
	}

	var annotations []*models.GOAnnotation
	if err := json.Unmarshal([]byte(cached), &annotations); err != nil {
		params.Logger.WithFields(logrus.Fields{
			"gene":  params.Gene,
			"error": err,
		}).Error("failed to unmarshal cached gene ontology annotations")
		errorutils.AddGQLError(params.Ctx, err)
		return nil, true, fmt.Errorf("error parsing cached data: %w", err)
	}
	return annotations, true, nil
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
		annotation, err := buildGOAnnotation(result, params.Redis)
		if err != nil {
			params.Logger.WithFields(logrus.Fields{
				"gene":          params.Gene,
				"annotation_id": result.ID,
				"error":         err,
			}).Error("failed to build GO annotation")
			errorutils.AddGQLError(params.Ctx, err)
			return nil, fmt.Errorf("error building annotation: %w", err)
		}
		annotations = append(annotations, annotation)
	}
	return annotations, nil
}

func CacheAnnotations(params *CacheAnnotationsParams) error {
	cached, err := json.Marshal(params.Annotations)
	if err != nil {
		params.Logger.WithFields(logrus.Fields{
			"gene":  params.Gene,
			"error": err,
		}).Error("failed to marshal annotations for caching")
		errorutils.AddGQLError(params.Ctx, err)
		return fmt.Errorf("error preparing data for cache: %w", err)
	}

	if err := params.Redis.SetWithTTL(getCacheKey(params.Gene), string(cached), 14*24*time.Hour); err != nil {
		params.Logger.WithFields(logrus.Fields{
			"gene":  params.Gene,
			"error": err,
		}).Error("failed to cache gene ontology annotations")
		errorutils.AddGQLError(params.Ctx, err)
		return fmt.Errorf("error caching annotations: %w", err)
	}
	return nil
}

func buildGOAnnotation(
	result result,
	redis repository.Repository,
) (*models.GOAnnotation, error) {
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
		with, exists, err := getWith(result.WithFrom, redis)
		if err != nil {
			return nil, fmt.Errorf("error getting with data: %w", err)
		}
		if exists {
			annotation.With = with
		}
	}

	if result.Extensions != nil {
		extensions, exists, err := getExtensions(result.Extensions, redis)
		if err != nil {
			return nil, fmt.Errorf("error getting extensions data: %w", err)
		}
		if exists {
			annotation.Extensions = extensions
		}
	}

	return annotation, nil
}
