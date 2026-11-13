package resolver

import (
	"strings"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	IO "github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/optics/lens"
	P "github.com/IBM/fp-go/v2/predicate"
	Str "github.com/IBM/fp-go/v2/string"
	predarrays "github.com/dictyBase/fp-go-loom/predicate/array"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
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
)

var (
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

	// saveOrWarnSearchHistory saves the text when a repository exists and
	// logs a warning when it does not.
	saveOrWarnSearchHistory = F.Pipe1(
		hasSearchHistoryRepo,
		P.Fold(warnMissingSearchHistoryRepo, saveStockSearchHistory),
	)

	// recordStockSearchHistory saves the search text of a successful
	// search that has results. It never fails the search.
	recordStockSearchHistory = F.Pipe1(
		shouldTrySearchHistory,
		P.Fold(
			F.Constant1[listStockSearchContext](IO.Of(F.VOID)),
			saveOrWarnSearchHistory,
		),
	)
)

// saveStockSearchHistory writes the normalized text with the time of the
// write as the score. A Redis error is logged and dropped.
func saveStockSearchHistory(s listStockSearchContext) IO.IO[IO.Void] {
	return F.Pipe1(
		IOE.TryCatchError(func() (IO.Void, error) {
			return F.VOID, searchHistoryLens.Get(s).AddRecentMember(
				StockSearchHistoryKey,
				normalizedSearchText(s),
				time.Now().UnixMilli(),
				stockSearchHistoryCapacity,
			)
		}),
		IOE.Fold(
			F.Bind1st(logSearchHistoryFailure, s),
			F.Constant1[IO.Void](IO.Of(F.VOID)),
		),
	)
}

func logSearchHistoryFailure(
	s listStockSearchContext,
	err error,
) IO.IO[IO.Void] {
	return IO.FromImpure(func() {
		searchLoggerLens.Get(s).Warnf(
			"could not save stock search history for %q: %s",
			normalizedSearchText(s),
			err,
		)
	})
}

func warnMissingSearchHistoryRepo(s listStockSearchContext) IO.IO[IO.Void] {
	return IO.FromImpure(func() {
		searchLoggerLens.Get(s).Warnf(
			"no redis repository, stock search history not saved for %q",
			normalizedSearchText(s),
		)
	})
}
