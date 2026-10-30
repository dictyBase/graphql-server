package resolver

import (
	"context"

	"github.com/dictyBase/graphql-server/internal/graphql/models"
)

// SearchPhenotypes resolves the searchPhenotypes GraphQL query. Task 3
// replaces this stub with the full implementation.
func (qrs *QueryResolver) SearchPhenotypes(
	ctx context.Context,
	query string,
	limit *int,
) (*models.StockSearchResultList, error) {
	panic("not implemented: SearchPhenotypes")
}
