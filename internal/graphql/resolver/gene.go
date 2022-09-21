package resolver

import (
	"context"

	"github.com/dictyBase/graphql-server/internal/graphql/cache"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
)

func (q *QueryResolver) Gene(
	ctx context.Context,
	gene string,
) (*models.Gene, error) {
	g := &models.Gene{}
	redis := q.GetRedisRepository(cache.RedisKey)
	gn, err := cache.GetGeneFromCache(ctx, redis, gene)
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		q.Logger.Error(err)
		return g, err
	}
	return gn, nil
}

func (q *QueryResolver) AllOrthologs(
	ctx context.Context,
	gene string,
) (*models.Gene, error) {
	return &models.Gene{}, nil
}

func (q *QueryResolver) GeneralInformation(
	ctx context.Context,
	gene string,
) (*models.Gene, error) {
	return &models.Gene{}, nil
}

// GetAssociatedSequnces is the resolver for the getAssociatedSequnces field.
func (q *QueryResolver) GetAssociatedSequnces(
	ctx context.Context,
	gene string,
) (*models.Gene, error) {
	return &models.Gene{}, nil
}

// GetLinks is the resolver for the getLinks field.
func (q *QueryResolver) GetLinks(
	ctx context.Context,
	gene string,
) (*models.Gene, error) {
	return &models.Gene{}, nil
}

// GetProteinInformation is the resolver for the getProteinInformation field.
func (q *QueryResolver) GetProteinInformation(
	ctx context.Context,
	gene string,
) (*models.Gene, error) {
	return &models.Gene{}, nil
}

// ListGeneProductInfo is the resolver for the listGeneProductInfo field.
func (q *QueryResolver) ListGeneProductInfo(
	ctx context.Context,
	gene string,
) (*models.Gene, error) {
	return &models.Gene{}, nil
}

// ListRecentGenes is the resolver for the listRecentGenes field.
func (q *QueryResolver) ListRecentGenes(
	ctx context.Context,
	limit int,
) ([]*models.Gene, error) {
	return []*models.Gene{&models.Gene{}}, nil
}
