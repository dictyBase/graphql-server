package resolver

import (
	"errors"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/dictyBase/graphql-server/internal/graphql/generated"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks"
	"github.com/dictyBase/graphql-server/internal/repository"
	"github.com/stretchr/testify/require"
)

const recentStockSearchesQuery = `query($limit: Int) {
	recentStockSearches(limit: $limit) { query searchedAt }
}`

type recentSearchesResponse struct {
	RecentStockSearches []struct {
		Query      string
		SearchedAt string
	}
}

// newHistoryGraphQLClient serves the generated executable schema over the
// real resolver tree with the fake history repository.
func newHistoryGraphQLClient(fake *mocks.HistoryRepository) *client.Client {
	reg := &searchHistoryRegistry{Registry: &mocks.MockRegistry{}}
	if fake != nil {
		reg.history = fake
	}
	resolver := NewResolver(reg, nil, mocks.TestLogger())
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{Resolvers: resolver}),
	)
	return client.New(srv)
}

func TestRecentStockSearchesGraphQL(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	fake := &mocks.HistoryRepository{Members: []repository.ScoredMember{
		{Member: "csa", Score: 1763049600123},
		{Member: "ax4", Score: 1763049500000},
	}}
	var resp recentSearchesResponse

	err := newHistoryGraphQLClient(fake).Post(
		recentStockSearchesQuery, &resp, client.Var("limit", 2),
	)

	require.NoError(err)
	require.Len(resp.RecentStockSearches, 2)
	require.Equal("csa", resp.RecentStockSearches[0].Query)
	require.Equal(
		"2025-11-13T16:00:00.123Z",
		resp.RecentStockSearches[0].SearchedAt,
		"searchedAt must be UTC with a literal Z",
	)
	require.Equal("ax4", resp.RecentStockSearches[1].Query)
	require.Equal([]int64{2}, fake.ReadCounts)
}

func TestRecentStockSearchesGraphQLDefaultLimit(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	fake := &mocks.HistoryRepository{}
	var resp recentSearchesResponse

	err := newHistoryGraphQLClient(fake).Post(
		`{ recentStockSearches { query } }`, &resp,
	)

	require.NoError(err)
	require.NotNil(resp.RecentStockSearches)
	require.Len(resp.RecentStockSearches, 0)
	require.Equal([]int64{5}, fake.ReadCounts)
}

func TestRecentStockSearchesGraphQLRedisDownGivesEmptyList(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	fake := &mocks.HistoryRepository{ReadErr: errors.New("redis down")}
	var resp recentSearchesResponse

	err := newHistoryGraphQLClient(fake).Post(
		recentStockSearchesQuery, &resp,
	)

	require.NoError(err, "redis failure must not add a GraphQL error")
	require.NotNil(resp.RecentStockSearches, "list must be [] not null")
	require.Len(resp.RecentStockSearches, 0)
}

func TestRecentStockSearchesGraphQLNoRepoGivesEmptyList(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	var resp recentSearchesResponse

	err := newHistoryGraphQLClient(nil).Post(recentStockSearchesQuery, &resp)

	require.NoError(err)
	require.NotNil(resp.RecentStockSearches)
	require.Len(resp.RecentStockSearches, 0)
}
