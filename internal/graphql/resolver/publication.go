package resolver

import (
	"context"
	"fmt"

	"github.com/dictyBase/graphql-server/internal/graphql/cache"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/fetch"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
)

// Publication is the resolver for getting an individual publication by ID.
func (qrs *QueryResolver) Publication(
	ctx context.Context,
	id string,
) (*models.Publication, error) {
	pub, err := fetch.FetchPublication(
		ctx, qrs.GetRedisRepository(cache.RedisKey),
		qrs.Registry.GetAPIEndpoint(registry.PUBLICATION), id,
	)
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		qrs.Logger.Error(err)
		return nil, fmt.Errorf(
			"error in fetching publication %s",
			err,
		)
	}
	return pub, nil
}

// AllPublications is the resolver for the allPublications field.
func (qrs *QueryResolver) AllPublications(
	ctx context.Context,
	gene string,
	limit *int,
	sortBy *string,
) (*models.NumberOfPublicationsWithGene, error) {
	return &models.NumberOfPublicationsWithGene{NumPubs: 0}, nil
}

// ListRecentPublications is the resolver for the listRecentPublications field.
func (qrs *QueryResolver) ListRecentPublications(
	ctx context.Context,
	limit int,
) ([]*models.Publication, error) {
	return []*models.Publication{}, nil
}

// ListPublicationsWithGene is the resolver for the listPublicationsWithGene field.
func (qrs *QueryResolver) ListPublicationsWithGene(
	ctx context.Context,
	gene string,
) ([]*models.PublicationWithGene, error) {
	return []*models.PublicationWithGene{}, nil
}
