package resolver

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/99designs/gqlgen/graphql"
	anno "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks/clients"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/graphql/resolverutils"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/dictyBase/graphql-server/internal/repository"
	"github.com/sirupsen/logrus"
	logrustest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	resolver, _ := newHookedHistoryResolver(fake)
	return resolver
}

// newHookedHistoryResolver also returns the hook that captures every
// log entry of the resolver.
func newHookedHistoryResolver(
	fake *mocks.HistoryRepository,
) (*QueryResolver, *logrustest.Hook) {
	reg := &searchHistoryRegistry{Registry: &mocks.MockRegistry{}}
	if fake != nil {
		reg.history = fake
	}
	logger, hook := logrustest.NewNullLogger()
	return &QueryResolver{Registry: reg, Logger: logrus.NewEntry(logger)}, hook
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
	ctx := graphql.WithResponseContext(
		context.Background(),
		graphql.DefaultErrorPresenter,
		graphql.DefaultRecover,
	)
	baseline, err := newHistoryResolver(nil).SearchStocks(ctx, "csa", nil, nil)
	require.NoError(err)
	fake := &mocks.HistoryRepository{Err: errors.New("redis down")}
	resolver, hook := newHookedHistoryResolver(fake)

	result, err := resolver.SearchStocks(ctx, "csa", nil, nil)

	require.NoError(err)
	require.Equal(baseline, result, "redis error must not change the result")
	require.Empty(graphql.GetErrors(ctx), "redis error must not add a GraphQL error")
	require.Len(fake.Calls, 1)
	warn := hook.LastEntry()
	require.NotNil(warn)
	require.Equal(logrus.WarnLevel, warn.Level)
	require.Contains(warn.Message, "could not save stock search history")
	require.Contains(warn.Message, "redis down")
}

func TestSearchStocksHistoryNilRepoSkips(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	ctx := graphql.WithResponseContext(
		context.Background(),
		graphql.DefaultErrorPresenter,
		graphql.DefaultRecover,
	)
	resolver, hook := newHookedHistoryResolver(nil)

	result, err := resolver.SearchStocks(ctx, "csa", nil, nil)

	require.NoError(err)
	require.Equal(4, result.TotalCount)
	require.Empty(graphql.GetErrors(ctx), "missing repository must not add a GraphQL error")
	warn := hook.LastEntry()
	require.NotNil(warn)
	require.Equal(logrus.WarnLevel, warn.Level)
	require.Contains(warn.Message, "no redis repository")
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

func TestSearchStocksHistorySkipsBacterialLookupFailure(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		SearchStock(
			mock.Anything,
			mock.AnythingOfType("*stock.StockSearchParameters"),
		).
		Return(mocks.MockStockSearchResultCollection(), nil)
	mockedStockClient.EXPECT().
		ListStrainsByIds(
			mock.Anything,
			mock.AnythingOfType("*stock.StockIdList"),
		).
		Return(mocks.MockSearchStrainList(), nil)
	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchAnnotation(
			mock.Anything,
			mock.AnythingOfType("*annotation.AnnotationSearchParameters"),
		).
		Return(mocks.MockAnnotationSearchResultCollection(), nil)
	mockedAnnoClient.EXPECT().
		ListAnnotations(
			mock.Anything,
			&anno.ListParameters{
				Limit:  100,
				Filter: resolverutils.BacterialAnnotationFilter(),
			},
		).
		Return(nil, status.Error(codes.Internal, "annotation service down"))
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
	require.Empty(fake.Calls, "search that fails at the bacterial lookup must not be saved")
}

// newRecentResolver builds a resolver whose history repository is fake.
// A nil fake gives a nil repository.
func newRecentResolver(
	fake *mocks.HistoryRepository,
) (*QueryResolver, *logrustest.Hook) {
	return newHookedHistoryResolver(fake)
}

func TestResolveRecentSearchLimit(t *testing.T) {
	t.Parallel()
	one, seven, max, over, huge := 1, 7, 250, 251, 1000
	zero, neg := 0, -3
	for name, tc := range map[string]struct {
		in   *int
		want int
	}{
		"nil":      {nil, 5},
		"zero":     {&zero, 5},
		"negative": {&neg, 5},
		"one":      {&one, 1},
		"seven":    {&seven, 7},
		"max":      {&max, 250},
		"over":     {&over, 250},
		"huge":     {&huge, 250},
	} {
		require.Equal(t, tc.want, resolveRecentSearchLimit(tc.in), name)
	}
}

func TestRecentStockSearchesMapsEntries(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	fake := &mocks.HistoryRepository{Members: []repository.ScoredMember{
		{Member: "csa", Score: 1763049600123},
		{Member: "ax4", Score: 1763049500000},
	}}
	resolver, _ := newRecentResolver(fake)

	result, err := resolver.RecentStockSearches(context.Background(), nil)

	require.NoError(err)
	require.Len(result, 2)
	require.Equal("csa", result[0].Query)
	require.True(result[0].SearchedAt.Equal(time.UnixMilli(1763049600123)))
	require.Equal(time.UTC, result[0].SearchedAt.Location())
	require.Equal("ax4", result[1].Query)
	require.Equal([]int64{5}, fake.ReadCounts)
}

func TestRecentStockSearchesPassesLimit(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	fake := &mocks.HistoryRepository{}
	resolver, _ := newRecentResolver(fake)
	three, huge := 3, 1000

	_, err := resolver.RecentStockSearches(context.Background(), &three)
	require.NoError(err)
	_, err = resolver.RecentStockSearches(context.Background(), &huge)
	require.NoError(err)

	require.Equal([]int64{3, 250}, fake.ReadCounts)
}

func TestRecentStockSearchesEmptyHistory(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	resolver, _ := newRecentResolver(&mocks.HistoryRepository{})

	result, err := resolver.RecentStockSearches(context.Background(), nil)

	require.NoError(err)
	require.NotNil(result)
	require.Len(result, 0)
}

func TestRecentStockSearchesRepoErrorReturnsEmpty(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	ctx := graphql.WithResponseContext(
		context.Background(),
		graphql.DefaultErrorPresenter,
		graphql.DefaultRecover,
	)
	fake := &mocks.HistoryRepository{ReadErr: errors.New("redis down")}
	resolver, hook := newRecentResolver(fake)

	result, err := resolver.RecentStockSearches(ctx, nil)

	require.NoError(err)
	require.NotNil(result)
	require.Len(result, 0)
	require.Empty(graphql.GetErrors(ctx), "redis error must not add a GraphQL error")
	warn := hook.LastEntry()
	require.NotNil(warn)
	require.Equal(logrus.WarnLevel, warn.Level)
	require.Contains(warn.Message, "could not read stock search history")
	require.Contains(warn.Message, "redis down")
}

func TestRecentStockSearchesNilRepoReturnsEmpty(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	resolver, hook := newRecentResolver(nil)

	result, err := resolver.RecentStockSearches(context.Background(), nil)

	require.NoError(err)
	require.NotNil(result)
	require.Len(result, 0)
	warn := hook.LastEntry()
	require.NotNil(warn)
	require.Equal(logrus.WarnLevel, warn.Level)
}
