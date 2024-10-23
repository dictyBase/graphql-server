package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/repository"
	"github.com/sirupsen/logrus"
)

const (
	geneHash               = "GENE2NAME/geneids"
	goHash                 = "GO2NAME/goids"
	uniprotHash            = "UNIPROT2NAME/uniprot"
	geneUniprotHash        = "GENE2UNIPROT/gene"
	geneGOAnnotationPrefix = "gene:goa:"
)

type CachedAnnotationsParams struct {
	Ctx    context.Context
	Gene   string
	Redis  repository.Repository
	Logger *logrus.Entry
}

type CacheAnnotationsParams struct {
	Ctx         context.Context
	Gene        string
	Annotations []*models.GOAnnotation
	Redis       repository.Repository
	Logger      *logrus.Entry
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
