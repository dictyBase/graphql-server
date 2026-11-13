package resolver

import (
	"context"
	"errors"
	"strings"
	"time"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IO "github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	N "github.com/IBM/fp-go/v2/number"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
	Str "github.com/IBM/fp-go/v2/string"
	MO "github.com/dictyBase/fp-go-loom/matchopt"
	predarrays "github.com/dictyBase/fp-go-loom/predicate/array"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/dictyBase/graphql-server/internal/repository"
	"github.com/sirupsen/logrus"
)

const (
	// StockSearchHistoryKey is the Redis sorted set that holds the saved
	// stock search texts. The member is the normalized text and the score
	// is the save time in Unix milliseconds.
	StockSearchHistoryKey = "stock-center:search-history"
	// stockSearchHistoryCapacity is the number of distinct texts that
	// Redis keeps, newest first.
	stockSearchHistoryCapacity = int64(250)
	// defaultRecentSearchLimit is the number of entries that
	// recentStockSearches returns when the client sends no limit.
	defaultRecentSearchLimit = 5
	// maxRecentSearchLimit is the stored capacity of the history.
	maxRecentSearchLimit = int(stockSearchHistoryCapacity)
)

// recentSearchState carries the request of recentStockSearches through
// its pipeline.
type recentSearchState struct {
	history repository.Repository
	limit   *int
}

var (
	recentHistoryLens = L.MakeLens(
		func(s recentSearchState) repository.Repository { return s.history },
		func(s recentSearchState, v repository.Repository) recentSearchState {
			s.history = v
			return s
		},
	)
	recentLimitLens = L.MakeLens(
		func(s recentSearchState) *int { return s.limit },
		func(s recentSearchState, v *int) recentSearchState {
			s.limit = v
			return s
		},
	)
	recentQueryLens = L.MakeLens(
		func(m models.RecentStockSearch) string { return m.Query },
		func(m models.RecentStockSearch, v string) models.RecentStockSearch {
			m.Query = v
			return m
		},
	)
	recentSearchedAtLens = L.MakeLens(
		func(m models.RecentStockSearch) time.Time { return m.SearchedAt },
		func(m models.RecentStockSearch, v time.Time) models.RecentStockSearch {
			m.SearchedAt = v
			return m
		},
	)

	isBelowRecentSearchLimit = N.LessThan(1)
	isAboveRecentSearchLimit = N.MoreThan(maxRecentSearchLimit)
	// clampRecentSearchLimit gives the default for a limit below 1, the
	// maximum for a limit above it, and passes every other limit.
	clampRecentSearchLimit = MO.Match(
		defaultRecentSearchLimit,
		[]O.Kleisli[int, int]{
			MO.Const(isBelowRecentSearchLimit, defaultRecentSearchLimit),
			MO.Const(isAboveRecentSearchLimit, maxRecentSearchLimit),
			MO.Default(F.Identity[int]),
		},
	)
	// recentCount converts the optional GraphQL limit. An absent limit or
	// a limit below 1 gives defaultRecentSearchLimit, and a limit above
	// maxRecentSearchLimit gives maxRecentSearchLimit.
	recentCount = F.Flow5(
		recentLimitLens.Get,
		O.FromNillable2[int],
		O.Map(clampRecentSearchLimit),
		O.GetOrElse(F.Constant(defaultRecentSearchLimit)),
		toInt64,
	)
	hasRecentHistoryRepo = F.Pipe1(
		P.IsNonZero[repository.Repository](),
		P.ContraMap(recentHistoryLens.Get),
	)
	// toRecentStockSearches maps the stored members and replaces the nil
	// slice that A.Map gives for an empty input with an empty slice.
	toRecentStockSearches = F.Flow3(
		A.Map(toRecentStockSearch),
		O.FromPredicate(predarrays.IsNonEmpty[*models.RecentStockSearch]()),
		O.GetOrElse(F.Constant([]*models.RecentStockSearch{})),
	)

	searchQueryLens = L.MakeLens(
		func(s listStockSearchContext) string { return s.query },
		func(s listStockSearchContext, v string) listStockSearchContext {
			s.query = v
			return s
		},
	)
	searchHistoryLens = L.MakeLens(
		func(s listStockSearchContext) repository.Repository {
			return s.history
		},
		func(
			s listStockSearchContext,
			v repository.Repository,
		) listStockSearchContext {
			s.history = v
			return s
		},
	)
	searchLoggerLens = L.MakeLens(
		func(s listStockSearchContext) *logrus.Entry { return s.logger },
		func(s listStockSearchContext, v *logrus.Entry) listStockSearchContext {
			s.logger = v
			return s
		},
	)

	// normalizeSearchHistoryText trims the text and folds it to lower
	// case, so that "  csA " and "CSA" are one history entry.
	normalizeSearchHistoryText = F.Flow2(strings.TrimSpace, strings.ToLower)
	normalizedSearchText       = F.Flow2(
		searchQueryLens.Get,
		normalizeSearchHistoryText,
	)

	hasSearchHistoryRepo = F.Pipe1(
		P.IsNonZero[repository.Repository](),
		P.ContraMap(searchHistoryLens.Get),
	)
	hasMergedSearchResults = F.Pipe1(
		predarrays.IsNonEmpty[*models.StockSearchResult](),
		P.ContraMap(mergedResultsLens.Get),
	)
	hasSearchHistoryText = F.Pipe1(
		Str.IsNonEmpty,
		P.ContraMap(normalizedSearchText),
	)
	shouldTrySearchHistory = F.Pipe1(
		hasMergedSearchResults,
		P.And(hasSearchHistoryText),
	)

	// errSkipSearchHistory marks a search that has nothing to save. It is
	// recovered silently. errNoSearchHistoryRepo is logged as a warning.
	errSkipSearchHistory   = errors.New("skip stock search history")
	errNoSearchHistoryRepo = errors.New("no redis repository")

	// recoverSkippedSearchHistory turns the skip marker into success and
	// keeps every other error.
	recoverSkippedSearchHistory = F.Pipe1(
		F.Bind2nd(errors.Is, errSkipSearchHistory),
		P.Fold(
			IOE.Left[IO.Void, error],
			F.Constant1[error](IOE.Of[error](F.VOID)),
		),
	)
)

// recordStockSearchHistory saves the search text of a successful search
// that has results. It never fails the search: a Redis error or a missing
// repository is logged as a warning.
func recordStockSearchHistory(s listStockSearchContext) IO.IO[IO.Void] {
	return F.Pipe7(
		IOE.Of[error](s),
		IOE.ChainEitherK(E.FromPredicate(
			shouldTrySearchHistory,
			F.Constant1[listStockSearchContext](errSkipSearchHistory),
		)),
		IOE.ChainEitherK(E.FromPredicate(
			hasSearchHistoryRepo,
			F.Constant1[listStockSearchContext](errNoSearchHistoryRepo),
		)),
		IOE.Chain(saveSearchHistory),
		IOE.OrElse(recoverSkippedSearchHistory),
		IOE.Swap[error, IO.Void],
		IOE.ChainConsumer[IO.Void](s.warnSearchHistoryFailure),
		IOE.Fold(
			IO.Of[IO.Void],
			F.Constant1[struct{}](IO.Of(F.VOID)),
		),
	)
}

// saveSearchHistory writes the normalized text with the time of the write
// as the score.
func saveSearchHistory(
	s listStockSearchContext,
) IOE.IOEither[error, IO.Void] {
	return IOE.TryCatchError(func() (IO.Void, error) {
		return F.VOID, searchHistoryLens.Get(s).AddRecentMember(
			StockSearchHistoryKey,
			normalizedSearchText(s),
			time.Now().UnixMilli(),
			stockSearchHistoryCapacity,
		)
	})
}

// warnSearchHistoryFailure is the consumer of a save failure. It logs a
// warning on the logger of the state.
func (s listStockSearchContext) warnSearchHistoryFailure(err error) {
	searchLoggerLens.Get(s).Warnf(
		"could not save stock search history for %q: %s",
		normalizedSearchText(s),
		err,
	)
}

// toRecentStockSearch maps one stored member to its GraphQL entry. The
// time is UTC, because models.MarshalTimestamp writes a literal Z.
func toRecentStockSearch(m repository.ScoredMember) *models.RecentStockSearch {
	searchedAt := time.UnixMilli(m.Score).UTC()

	return F.Pipe3(
		models.RecentStockSearch{},
		recentQueryLens.Set(m.Member),
		recentSearchedAtLens.Set(searchedAt),
		F.Ref[models.RecentStockSearch],
	)
}

func readRecentSearches(
	s recentSearchState,
) IOE.IOEither[error, []repository.ScoredMember] {
	return IOE.TryCatchError(func() ([]repository.ScoredMember, error) {
		return recentHistoryLens.Get(s).RecentMembers(
			StockSearchHistoryKey,
			recentCount(s),
		)
	})
}

// RecentStockSearches returns the newest distinct texts of the
// searchStocks query. A read error or a missing repository gives an empty
// list and a logged warning, never a GraphQL error.
func (qrs *QueryResolver) RecentStockSearches(
	_ context.Context,
	limit *int,
) ([]*models.RecentStockSearch, error) {
	history := qrs.GetRedisRepository(registry.REDISREPO)
	list := F.Pipe8(
		recentSearchState{},
		recentHistoryLens.Set(history),
		recentLimitLens.Set(limit),
		IOE.Of[error, recentSearchState],
		IOE.ChainEitherK(E.FromPredicate(
			hasRecentHistoryRepo,
			F.Constant1[recentSearchState](errNoSearchHistoryRepo),
		)),
		IOE.Chain(readRecentSearches),
		IOE.Map[error](toRecentStockSearches),
		IOE.TapLeft[[]*models.RecentStockSearch](F.Flow2(
			IO.FromConsumer(qrs.warnRecentSearchesFailure),
			IOE.FromIO[error, IO.Void],
		)),
		IOE.GetOrElse(F.Constant1[error](IO.Of([]*models.RecentStockSearch{}))),
	)()

	return list, nil
}

// warnRecentSearchesFailure logs a warning for a failed history read.
func (qrs *QueryResolver) warnRecentSearchesFailure(err error) {
	qrs.Logger.Warnf("could not read stock search history: %s", err)
}
