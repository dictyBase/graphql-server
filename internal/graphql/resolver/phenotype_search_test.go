package resolver

import (
	"context"
	"testing"

	anno "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks/clients"
	"github.com/dictyBase/graphql-server/internal/graphql/models"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// phenotypeClientRegistry returns fixed annotation and stock clients for
// phenotype search resolver tests that verify exact request parameters.
type phenotypeClientRegistry struct {
	*mocks.MockRegistry
	annoClient  anno.TaggedAnnotationServiceClient
	stockClient pb.StockServiceClient
}

func (r *phenotypeClientRegistry) GetAnnotationClient(
	key string,
) anno.TaggedAnnotationServiceClient {
	return r.annoClient
}

func (r *phenotypeClientRegistry) GetStockClient(
	key string,
) pb.StockServiceClient {
	return r.stockClient
}

func TestSearchPhenotypes(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	resolver := &QueryResolver{
		Registry: &mocks.MockRegistry{},
		Logger:   mocks.TestLogger(),
	}

	result, err := resolver.SearchPhenotypes(context.Background(), "delayed", nil)
	require.NoError(err)
	require.Len(result.Results, 2, "should return two phenotype matches")
	require.NotNil(result.Limit, "should resolve limit from meta")
	require.Equal(50, *result.Limit, "should match effective limit from meta")
	require.Equal(2, result.TotalCount, "should count the returned rows")

	first := result.Results[0]
	require.Equal(mocks.MockStockSearchID, first.ID, "should map entry_id to id")
	require.Equal(
		models.StockEntityTypeStrain,
		first.Entity,
		"should map strain entity",
	)
	require.Equal(
		models.StockSearchFieldEnumTermLabel,
		first.Field,
		"should map a term label match to TERM_LABEL",
	)
	require.Equal(
		"delayed culmination",
		first.DisplayText,
		"should show the term name for TERM_LABEL",
	)
	require.Equal(1250.5, first.Score)
	require.Equal(
		"GWDI strain carrying a csA deletion.",
		first.Summary,
		"should enrich summary from ListStrainsByIds",
	)
	require.Equal("AX4", first.StrainLabel, "should enrich strain label from ListStrainsByIds")
	require.Equal("gwdi", first.StockType, "should map REMI-seq strain to gwdi")

	second := result.Results[1]
	require.Equal(mocks.MockAnnoSearchStrainID, second.ID)
	require.Equal(
		models.StockSearchFieldEnumTermDefinition,
		second.Field,
		"should map a term definition match to TERM_DEFINITION",
	)
	require.Equal(
		"Cell shape differs from the wild-type shape.",
		second.DisplayText,
		"should show the definition for TERM_DEFINITION",
	)
	require.Equal(300.25, second.Score)
	require.Equal("Catenin null mutant strain.", second.Summary)
	require.Equal("aarA-", second.StrainLabel)
	require.Equal("regular", second.StockType, "should map general strain to regular")
}

func TestSearchPhenotypesFieldMapping(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchPhenotypes(
			mock.Anything,
			mock.AnythingOfType("*annotation.PhenotypeSearchParameters"),
		).
		Return(&anno.PhenotypeSearchResultCollection{
			Data: []*anno.PhenotypeSearchResult{
				{
					EntryId:    mocks.MockStockSearchID,
					Field:      anno.PhenotypeSearchField_PHENOTYPE_SEARCH_FIELD_SYNONYM,
					Name:       "late culmination",
					Definition: "Strain completes development later than wild type.",
					Score:      800.5,
				},
				{
					EntryId:    mocks.MockAnnoSearchStrainID,
					Field:      anno.PhenotypeSearchField_PHENOTYPE_SEARCH_FIELD_COMMENT,
					Name:       "aberrant morphology",
					Definition: "Cell shape differs from the wild-type shape.",
					Score:      200.25,
				},
			},
			Meta: &anno.Meta{
				Limit: 50,
				Total: 2,
			},
		}, nil)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		ListStrainsByIds(
			mock.Anything,
			mock.AnythingOfType("*stock.StockIdList"),
		).
		Return(mocks.MockSearchStrainList(), nil)

	reg := &phenotypeClientRegistry{
		MockRegistry: &mocks.MockRegistry{},
		annoClient:   mockedAnnoClient,
		stockClient:  mockedStockClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	result, err := resolver.SearchPhenotypes(
		context.Background(),
		"culmination",
		nil,
	)
	require.NoError(err)
	require.Len(result.Results, 2)

	synonym := result.Results[0]
	require.Equal(
		models.StockSearchFieldEnumTermSynonym,
		synonym.Field,
		"should map a term synonym match to TERM_SYNONYM",
	)
	require.Equal(
		"late culmination",
		synonym.DisplayText,
		"should show the term name for TERM_SYNONYM",
	)

	comment := result.Results[1]
	require.Equal(
		models.StockSearchFieldEnumTermComment,
		comment.Field,
		"should map a term comment match to TERM_COMMENT",
	)
	require.Equal(
		"Cell shape differs from the wild-type shape.",
		comment.DisplayText,
		"should show the definition for TERM_COMMENT",
	)
}

func TestSearchPhenotypesPassesExpectedParams(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchPhenotypes(
			mock.Anything,
			mock.MatchedBy(func(params *anno.PhenotypeSearchParameters) bool {
				attr := params.Data.Attributes
				return params.Data.Type == "phenotype" &&
					attr.Query == "delayed" &&
					attr.Limit == 25
			}),
		).
		Return(mocks.MockPhenotypeSearchResultCollection(), nil)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		ListStrainsByIds(
			mock.Anything,
			&pb.StockIdList{
				Id: []string{mocks.MockAnnoSearchStrainID, mocks.MockStockSearchID},
			},
		).
		Return(mocks.MockSearchStrainList(), nil)

	reg := &phenotypeClientRegistry{
		MockRegistry: &mocks.MockRegistry{},
		annoClient:   mockedAnnoClient,
		stockClient:  mockedStockClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}
	limit := 25

	result, err := resolver.SearchPhenotypes(
		context.Background(),
		"delayed",
		&limit,
	)
	require.NoError(err)
	require.Len(result.Results, 2)
}

func TestSearchPhenotypesDefaultLimitPassesZero(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchPhenotypes(
			mock.Anything,
			mock.MatchedBy(func(params *anno.PhenotypeSearchParameters) bool {
				return params.Data.Attributes.Limit == 0
			}),
		).
		Return(mocks.MockPhenotypeSearchResultCollection(), nil)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		ListStrainsByIds(
			mock.Anything,
			mock.AnythingOfType("*stock.StockIdList"),
		).
		Return(mocks.MockSearchStrainList(), nil)

	reg := &phenotypeClientRegistry{
		MockRegistry: &mocks.MockRegistry{},
		annoClient:   mockedAnnoClient,
		stockClient:  mockedStockClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	result, err := resolver.SearchPhenotypes(context.Background(), "delayed", nil)
	require.NoError(err)
	require.NotNil(result.Limit)
	require.Equal(50, *result.Limit, "should resolve limit from meta")
}

func TestSearchPhenotypesClampsLimitAboveFifty(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchPhenotypes(
			mock.Anything,
			mock.MatchedBy(func(params *anno.PhenotypeSearchParameters) bool {
				return params.Data.Attributes.Limit == 50
			}),
		).
		Return(mocks.MockPhenotypeSearchResultCollection(), nil)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		ListStrainsByIds(
			mock.Anything,
			mock.AnythingOfType("*stock.StockIdList"),
		).
		Return(mocks.MockSearchStrainList(), nil)

	reg := &phenotypeClientRegistry{
		MockRegistry: &mocks.MockRegistry{},
		annoClient:   mockedAnnoClient,
		stockClient:  mockedStockClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}
	limit := 100

	result, err := resolver.SearchPhenotypes(
		context.Background(),
		"delayed",
		&limit,
	)
	require.NoError(err)
	require.NotNil(result.Limit)
	require.Equal(50, *result.Limit, "should echo the service limit of 50")
}

func TestSearchPhenotypesShortQuery(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedStockClient := clients.NewStockServiceClient(t)
	reg := &phenotypeClientRegistry{
		MockRegistry: &mocks.MockRegistry{},
		annoClient:   mockedAnnoClient,
		stockClient:  mockedStockClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	result, err := resolver.SearchPhenotypes(context.Background(), "a", nil)
	require.Error(err)
	require.Equal(codes.InvalidArgument, status.Code(err))
	require.ErrorContains(err, "expect a search query of at least 2 characters")
	require.NotNil(result, "should return an empty list with the error")
}

func TestSearchPhenotypesEmptyResults(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchPhenotypes(
			mock.Anything,
			mock.AnythingOfType("*annotation.PhenotypeSearchParameters"),
		).
		Return(mocks.MockEmptyPhenotypeSearchResultCollection(), nil)

	mockedStockClient := clients.NewStockServiceClient(t)
	reg := &phenotypeClientRegistry{
		MockRegistry: &mocks.MockRegistry{},
		annoClient:   mockedAnnoClient,
		stockClient:  mockedStockClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	result, err := resolver.SearchPhenotypes(context.Background(), "delayed", nil)
	require.NoError(err)
	require.Empty(result.Results)
	require.Equal(0, result.TotalCount)
}

func TestSearchPhenotypesDropsHitWithoutStrain(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchPhenotypes(
			mock.Anything,
			mock.AnythingOfType("*annotation.PhenotypeSearchParameters"),
		).
		Return(mocks.MockPhenotypeSearchResultCollection(), nil)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		ListStrainsByIds(
			mock.Anything,
			mock.AnythingOfType("*stock.StockIdList"),
		).
		Return(&pb.StrainList{}, nil)

	reg := &phenotypeClientRegistry{
		MockRegistry: &mocks.MockRegistry{},
		annoClient:   mockedAnnoClient,
		stockClient:  mockedStockClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	result, err := resolver.SearchPhenotypes(context.Background(), "delayed", nil)
	require.NoError(err)
	require.Empty(result.Results, "should drop hits whose strain is gone")
	require.Equal(0, result.TotalCount)
}

func TestSearchPhenotypesStrainListNotFound(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchPhenotypes(
			mock.Anything,
			mock.AnythingOfType("*annotation.PhenotypeSearchParameters"),
		).
		Return(mocks.MockPhenotypeSearchResultCollection(), nil)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		ListStrainsByIds(
			mock.Anything,
			mock.AnythingOfType("*stock.StockIdList"),
		).
		Return(nil, status.Error(codes.NotFound, "no such strains"))

	reg := &phenotypeClientRegistry{
		MockRegistry: &mocks.MockRegistry{},
		annoClient:   mockedAnnoClient,
		stockClient:  mockedStockClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	result, err := resolver.SearchPhenotypes(context.Background(), "delayed", nil)
	require.NoError(err, "should recover a missing strain list as empty metadata")
	require.Empty(result.Results)
	require.Equal(0, result.TotalCount)
}

func TestSearchPhenotypesStrainListError(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchPhenotypes(
			mock.Anything,
			mock.AnythingOfType("*annotation.PhenotypeSearchParameters"),
		).
		Return(mocks.MockPhenotypeSearchResultCollection(), nil)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		ListStrainsByIds(
			mock.Anything,
			mock.AnythingOfType("*stock.StockIdList"),
		).
		Return(nil, status.Error(codes.Internal, "stock service down"))

	reg := &phenotypeClientRegistry{
		MockRegistry: &mocks.MockRegistry{},
		annoClient:   mockedAnnoClient,
		stockClient:  mockedStockClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	result, err := resolver.SearchPhenotypes(context.Background(), "delayed", nil)
	require.Error(err)
	require.ErrorContains(err, "fetch strain metadata for query")
	require.ErrorContains(err, "stock service down")
	require.NotNil(result, "should return an empty list with the error")
}

func TestSearchPhenotypesInvalidFieldEnum(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchPhenotypes(
			mock.Anything,
			mock.AnythingOfType("*annotation.PhenotypeSearchParameters"),
		).
		Return(&anno.PhenotypeSearchResultCollection{
			Data: []*anno.PhenotypeSearchResult{
				{
					EntryId:    mocks.MockStockSearchID,
					Field:      anno.PhenotypeSearchField_PHENOTYPE_SEARCH_FIELD_UNSPECIFIED,
					Name:       "delayed culmination",
					Definition: "Strain completes development later than wild type.",
					Score:      1250.5,
				},
			},
			Meta: &anno.Meta{
				Limit: 50,
				Total: 1,
			},
		}, nil)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		ListStrainsByIds(
			mock.Anything,
			mock.AnythingOfType("*stock.StockIdList"),
		).
		Return(mocks.MockSearchStrainList(), nil)

	reg := &phenotypeClientRegistry{
		MockRegistry: &mocks.MockRegistry{},
		annoClient:   mockedAnnoClient,
		stockClient:  mockedStockClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	result, err := resolver.SearchPhenotypes(context.Background(), "delayed", nil)
	require.Error(err)
	require.ErrorContains(err, "unexpected phenotype search field")
	require.NotNil(result, "should return an empty list with the error")
}

func TestSearchPhenotypesError(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchPhenotypes(
			mock.Anything,
			mock.AnythingOfType("*annotation.PhenotypeSearchParameters"),
		).
		Return(nil, status.Error(codes.Internal, "backend down"))

	mockedStockClient := clients.NewStockServiceClient(t)
	reg := &phenotypeClientRegistry{
		MockRegistry: &mocks.MockRegistry{},
		annoClient:   mockedAnnoClient,
		stockClient:  mockedStockClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	result, err := resolver.SearchPhenotypes(context.Background(), "delayed", nil)
	require.Error(err)
	require.ErrorContains(err, "search phenotypes for query")
	require.ErrorContains(err, "backend down")
	require.NotNil(result, "should return an empty list with the error")
}
