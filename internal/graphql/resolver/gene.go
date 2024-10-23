package resolver

import (
	"context"
	"fmt"

	"github.com/dictyBase/graphql-server/internal/graphql/cache"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/sirupsen/logrus"
)

func (qrs *QueryResolver) GeneOntologyAnnotation(
	ctx context.Context,
	gene string,
) ([]*models.GOAnnotation, error) {
	redis := qrs.GetRedisRepository(registry.REDISREPO)

	// Try to get from cache first
	annotations, found, err := cache.GetCachedAnnotations(
		&cache.CachedAnnotationsParams{
			Ctx:    ctx,
			Gene:   gene,
			Redis:  redis,
			Logger: qrs.Logger,
		},
	)
	if err != nil {
		return nil, err
	}
	if found {
		return annotations, nil
	}

	// Get UniProt ID for the gene
	uniprotID, err := cache.GetUniprotIDForGene(gene, redis)
	if err != nil {
		qrs.Logger.WithFields(logrus.Fields{
			"gene":  gene,
			"error": err,
		}).Error("failed to get UniProt ID for gene")
		errorutils.AddGQLError(ctx, err)
		return nil, fmt.Errorf("error getting UniProt ID: %w", err)
	}

	// Fetch and build annotations
	annotations, err = cache.FetchAndBuildAnnotations(
		&cache.FetchAndBuildAnnotationsParams{
			Ctx:       ctx,
			Gene:      gene,
			UniprotID: uniprotID,
			Redis:     redis,
			Logger:    qrs.Logger,
		},
	)
	if err != nil {
		return nil, err
	}

	// Cache the results
	if err := cache.CacheAnnotations(&cache.CacheAnnotationsParams{
		Ctx:         ctx,
		Gene:        gene,
		Annotations: annotations,
		Redis:       redis,
		Logger:      qrs.Logger,
	}); err != nil {
		return nil, err
	}

	return annotations, nil
}

func (qrs *QueryResolver) Gene(
	ctx context.Context,
	geneID string,
) (*models.Gene, error) {
	redis := qrs.GetRedisRepository(registry.REDISREPO)
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
