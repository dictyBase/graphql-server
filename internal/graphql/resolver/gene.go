package resolver

import (
	"context"

	"github.com/dictyBase/graphql-server/internal/graphql/cache"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
)

func (qrs *QueryResolver) Gene(
	ctx context.Context,
	gene string,
) (*models.Gene, error) {
	g := &models.Gene{}
	redis := qrs.GetRedisRepository(cache.RedisKey)
	gn, err := cache.GetGeneFromCache(ctx, redis, gene)
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		qrs.Logger.Error(err)
		return g, err
	}
	return gn, nil
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
