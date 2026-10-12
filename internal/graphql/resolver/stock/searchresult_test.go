package stock

import (
	"context"
	"testing"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks/clients"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func searchResultResolver(
	annoClient *clients.TaggedAnnotationServiceClient,
) *StockSearchResultResolver {
	return &StockSearchResultResolver{
		AnnotationClient: annoClient,
		Logger:           mocks.TestLogger(),
	}
}

func mockSearchResult(entity models.StockEntityType) *models.StockSearchResult {
	return &models.StockSearchResult{
		ID:          "DBS123456",
		Entity:      entity,
		Field:       models.StockSearchFieldEnumStockID,
		DisplayText: "DBS123456",
		Score:       1,
		Summary:     "GWDI strain carrying a csA deletion.",
		StrainLabel: "AX4",
	}
}

func TestSearchResultStrainInStock(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	r := searchResultResolver(mocks.MockedInStockClient())

	inStock, err := r.InStock(
		context.Background(),
		mockSearchResult(models.StockEntityTypeStrain),
	)

	assert.NoError(err, "expect no error from the inventory lookup")
	assert.True(inStock, "should return true after finding the inventory entry")
}

func TestSearchResultStrainNotInStock(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	r := searchResultResolver(mocks.MockedNotInStockClient())

	inStock, err := r.InStock(
		context.Background(),
		mockSearchResult(models.StockEntityTypeStrain),
	)

	assert.NoError(err, "expect no error from a missing inventory entry")
	assert.False(inStock, "should return false after a not-found inventory entry")
}

func TestSearchResultPlasmidInStock(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	r := searchResultResolver(mocks.MockedInStockClient())

	inStock, err := r.InStock(
		context.Background(),
		mockSearchResult(models.StockEntityTypePlasmid),
	)

	assert.NoError(err, "expect no error from the inventory lookup")
	assert.True(inStock, "should return true after finding the inventory entry")
}

func TestSearchResultPlasmidNotInStock(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	r := searchResultResolver(mocks.MockedNotInStockClient())

	inStock, err := r.InStock(
		context.Background(),
		mockSearchResult(models.StockEntityTypePlasmid),
	)

	assert.NoError(err, "expect no error from a missing inventory entry")
	assert.False(inStock, "should return false after a not-found inventory entry")
}

func TestSearchResultInStockError(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	r := searchResultResolver(mocks.MockedInStockErrorClient())

	inStock, err := r.InStock(
		context.Background(),
		mockSearchResult(models.StockEntityTypeStrain),
	)

	require.Error(err, "expect an error from the inventory service failure")
	require.ErrorContains(err, "error getting DBS123456 from inventory")
	require.False(inStock, "should return false on an inventory service failure")
}

func TestSearchResultInStockUnknownEntity(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	r := searchResultResolver(mocks.MockedInStockClient())

	inStock, err := r.InStock(
		context.Background(),
		mockSearchResult(models.StockEntityTypeAll),
	)

	require.Error(err, "expect an error for an unknown stock entity")
	require.ErrorContains(err, "unexpected stock entity")
	require.False(inStock, "should return false for an unknown stock entity")
}

func TestSearchResultInStockPassesInventoryRef(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	strainAnno := new(clients.TaggedAnnotationServiceClient)
	strainAnno.On(
		"GetEntryAnnotation",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.MatchedBy(func(req *annotation.EntryAnnotationRequest) bool {
			return req.Tag == registry.StrainInvTag &&
				req.Ontology == registry.StrainInvOnto &&
				req.EntryId == "DBS123456"
		}),
	).Return(mocks.MockInStockAnno, nil)

	plasmidAnno := new(clients.TaggedAnnotationServiceClient)
	plasmidAnno.On(
		"GetEntryAnnotation",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.MatchedBy(func(req *annotation.EntryAnnotationRequest) bool {
			return req.Tag == registry.PlasmidInvTag &&
				req.Ontology == registry.PlasmidInvOnto &&
				req.EntryId == "DBP0000120"
		}),
	).Return(mocks.MockInStockAnno, nil)

	strainResult, strainErr := searchResultResolver(strainAnno).InStock(
		context.Background(),
		mockSearchResult(models.StockEntityTypeStrain),
	)
	require.NoError(strainErr)
	require.True(strainResult)
	strainAnno.AssertExpectations(t)

	plasmidResult, plasmidErr := searchResultResolver(plasmidAnno).InStock(
		context.Background(),
		&models.StockSearchResult{
			ID:     "DBP0000120",
			Entity: models.StockEntityTypePlasmid,
			Field:  models.StockSearchFieldEnumName,
		},
	)
	require.NoError(plasmidErr)
	require.True(plasmidResult)
	plasmidAnno.AssertExpectations(t)
}
