package resolver

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks/clients"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/dictyBase/graphql-server/internal/repository"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// searchHistoryRegistry overrides only GetRedisRepository. A nil history
// yields a nil repository, which the resolver must skip.
type searchHistoryRegistry struct {
	registry.Registry
	history repository.Repository
}

func (r *searchHistoryRegistry) GetRedisRepository(
	_ string,
) repository.Repository {
	return r.history
}

// newHistoryResolver builds a resolver over the default mock clients
// with the given history repository. A nil fake gives a nil repository.
func newHistoryResolver(fake *mocks.HistoryRepository) *QueryResolver {
	reg := &searchHistoryRegistry{Registry: &mocks.MockRegistry{}}
	if fake != nil {
		reg.history = fake
	}
	return &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}
}

func TestNormalizeSearchHistoryText(t *testing.T) {
	t.Parallel()
	for in, want := range map[string]string{
		"  CsA  ":    "csa",
		"GFP-Tag":    "gfp-tag",
		"ΣΑ":         "σα",
		"\t\n":       "",
		"already-ok": "already-ok",
	} {
		require.Equal(t, want, normalizeSearchHistoryText(in), "input %q", in)
	}
}

func TestSearchStocksHistorySavesText(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	fake := &mocks.HistoryRepository{}
	resolver := newHistoryResolver(fake)

	before := time.Now().UnixMilli()
	result, err := resolver.SearchStocks(context.Background(), "csa", nil, nil)
	after := time.Now().UnixMilli()

	require.NoError(err)
	require.Equal(4, result.TotalCount)
	require.Len(fake.Calls, 1)
	call := fake.Calls[0]
	require.Equal(StockSearchHistoryKey, call.Key)
	require.Equal("csa", call.Member)
	require.Equal(int64(250), call.Capacity)
	require.GreaterOrEqual(call.Score, before)
	require.LessOrEqual(call.Score, after)
}

func TestSearchStocksHistoryNormalizesText(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	fake := &mocks.HistoryRepository{}
	resolver := newHistoryResolver(fake)

	_, err := resolver.SearchStocks(context.Background(), "  CsA  ", nil, nil)

	require.NoError(err)
	require.Len(fake.Calls, 1)
	require.Equal("csa", fake.Calls[0].Member)
}

func TestSearchStocksHistoryErrorKeepsResult(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	fake := &mocks.HistoryRepository{Err: errors.New("redis down")}
	resolver := newHistoryResolver(fake)

	result, err := resolver.SearchStocks(context.Background(), "csa", nil, nil)

	require.NoError(err)
	require.Equal(4, result.TotalCount)
	require.Len(fake.Calls, 1)
}

func TestSearchStocksHistoryNilRepoSkips(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	resolver := newHistoryResolver(nil)

	result, err := resolver.SearchStocks(context.Background(), "csa", nil, nil)

	require.NoError(err)
	require.Equal(4, result.TotalCount)
}

func TestSearchStocksHistorySkipsBlankQuery(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	fake := &mocks.HistoryRepository{}
	resolver := newHistoryResolver(fake)

	_, err := resolver.SearchStocks(context.Background(), "   ", nil, nil)

	require.NoError(err)
	require.Empty(fake.Calls, "blank text must not be saved")
}

func TestSearchStocksHistorySkipsZeroResults(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		SearchStock(
			mock.Anything,
			mock.AnythingOfType("*stock.StockSearchParameters"),
		).
		Return(&stock.StockSearchResultCollection{
			Data: []*stock.StockSearchResult{},
			Meta: &stock.Meta{Limit: 50, Total: 0},
		}, nil)
	fake := &mocks.HistoryRepository{}
	resolver := &QueryResolver{
		Registry: &searchHistoryRegistry{
			Registry: &bacterialStrainRegistry{
				MockRegistry: &mocks.MockRegistry{ConnMap: nil},
				stockClient:  mockedStockClient,
				annoClient:   clients.NewTaggedAnnotationServiceClient(t),
			},
			history: fake,
		},
		Logger: mocks.TestLogger(),
	}
	entity := models.StockEntityTypePlasmid

	result, err := resolver.SearchStocks(
		context.Background(), "nomatch", nil, &entity,
	)

	require.NoError(err)
	require.Equal(0, result.TotalCount)
	require.Empty(fake.Calls, "zero-result search must not be saved")
}

func TestSearchStocksHistorySkipsFailedSearch(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		SearchStock(
			mock.Anything,
			mock.AnythingOfType("*stock.StockSearchParameters"),
		).
		Return(nil, errors.New("boom"))
	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchAnnotation(
			mock.Anything,
			mock.AnythingOfType("*annotation.AnnotationSearchParameters"),
		).
		Return(mocks.MockAnnotationSearchResultCollection(), nil)
	fake := &mocks.HistoryRepository{}
	resolver := &QueryResolver{
		Registry: &searchHistoryRegistry{
			Registry: &bacterialStrainRegistry{
				MockRegistry: &mocks.MockRegistry{ConnMap: nil},
				stockClient:  mockedStockClient,
				annoClient:   mockedAnnoClient,
			},
			history: fake,
		},
		Logger: mocks.TestLogger(),
	}

	_, err := resolver.SearchStocks(context.Background(), "csa", nil, nil)

	require.Error(err)
	require.Empty(fake.Calls, "failed search must not be saved")
}
