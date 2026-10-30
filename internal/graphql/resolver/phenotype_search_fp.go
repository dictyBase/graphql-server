package resolver

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	ORD "github.com/IBM/fp-go/v2/ord"
	PR "github.com/IBM/fp-go/v2/pair"
	P "github.com/IBM/fp-go/v2/predicate"
	R "github.com/IBM/fp-go/v2/record"
	T "github.com/IBM/fp-go/v2/tuple"
	ioeutils "github.com/dictyBase/fp-go-loom/ioeitherutils"
	MO "github.com/dictyBase/fp-go-loom/matchopt"
	predarrays "github.com/dictyBase/fp-go-loom/predicate/array"
	predord "github.com/dictyBase/fp-go-loom/predicate/ord"
	anno "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// phenotypeSearchContext carries one phenotype search request through
// the search pipeline. Only the caller-supplied fields are set at the
// seed; all derived fields are attached by lenses so every pipeline step
// receives a single state parameter.
type phenotypeSearchContext struct {
	client         anno.TaggedAnnotationServiceClient
	stockClient    pb.StockServiceClient
	gctx           context.Context
	query          string
	limit          *int
	params         *anno.PhenotypeSearchParameters
	collection     *anno.PhenotypeSearchResultCollection
	strainMetaByID map[string]strainSearchMetadata
	results        []*models.StockSearchResult
	effectiveLimit int
}

// phenotypeFieldSpec holds the GraphQL field enum of one matched term
// field and the extractor of the text the user interface shows for it.
type phenotypeFieldSpec struct {
	field models.StockSearchFieldEnum
	text  func(*anno.PhenotypeSearchResult) string
}

type phenotypeSearchFieldPair = PR.Pair[*anno.PhenotypeSearchResult, phenotypeFieldSpec]

// ── state lenses ────────────────────────────────────────────────────────────

var (
	phenotypeSearchParamsLens = L.MakeLens(
		func(s phenotypeSearchContext) *anno.PhenotypeSearchParameters {
			return s.params
		},
		func(
			s phenotypeSearchContext,
			params *anno.PhenotypeSearchParameters,
		) phenotypeSearchContext {
			s.params = params
			return s
		},
	)

	phenotypeSearchCollectionLens = L.MakeLens(
		func(s phenotypeSearchContext) *anno.PhenotypeSearchResultCollection {
			return s.collection
		},
		func(
			s phenotypeSearchContext,
			collection *anno.PhenotypeSearchResultCollection,
		) phenotypeSearchContext {
			s.collection = collection
			return s
		},
	)

	phenotypeStrainMetaLens = L.MakeLens(
		func(s phenotypeSearchContext) map[string]strainSearchMetadata {
			return s.strainMetaByID
		},
		func(
			s phenotypeSearchContext,
			metaByID map[string]strainSearchMetadata,
		) phenotypeSearchContext {
			s.strainMetaByID = metaByID
			return s
		},
	)

	phenotypeSearchResultsLens = L.MakeLens(
		func(s phenotypeSearchContext) []*models.StockSearchResult {
			return s.results
		},
		func(
			s phenotypeSearchContext,
			results []*models.StockSearchResult,
		) phenotypeSearchContext {
			s.results = results
			return s
		},
	)

	phenotypeSearchLimitLens = L.MakeLens(
		func(s phenotypeSearchContext) int {
			return s.effectiveLimit
		},
		func(
			s phenotypeSearchContext,
			limit int,
		) phenotypeSearchContext {
			s.effectiveLimit = limit
			return s
		},
	)
)

// ── field enum mapping ──────────────────────────────────────────────────────

// phenotypeFieldSpecs maps the four phenotype term fields onto their
// GraphQL StockSearchFieldEnum counterparts and display text. The
// unspecified value is absent on purpose, so the lookup below fails. A
// TERM_COMMENT match shows the definition, because the protobuf response
// carries no comment text.
var phenotypeFieldSpecs = map[anno.PhenotypeSearchField]phenotypeFieldSpec{
	anno.PhenotypeSearchField_PHENOTYPE_SEARCH_FIELD_LABEL: {
		field: models.StockSearchFieldEnumTermLabel,
		text:  (*anno.PhenotypeSearchResult).GetName,
	},
	anno.PhenotypeSearchField_PHENOTYPE_SEARCH_FIELD_SYNONYM: {
		field: models.StockSearchFieldEnumTermSynonym,
		text:  (*anno.PhenotypeSearchResult).GetName,
	},
	anno.PhenotypeSearchField_PHENOTYPE_SEARCH_FIELD_DEFINITION: {
		field: models.StockSearchFieldEnumTermDefinition,
		text:  (*anno.PhenotypeSearchResult).GetDefinition,
	},
	anno.PhenotypeSearchField_PHENOTYPE_SEARCH_FIELD_COMMENT: {
		field: models.StockSearchFieldEnumTermComment,
		text:  (*anno.PhenotypeSearchResult).GetDefinition,
	},
}

// ── query validation ────────────────────────────────────────────────────────

// trimmedPhenotypeQueryRuneCount counts the runes of the query once
// surrounding whitespace is trimmed.
func trimmedPhenotypeQueryRuneCount(s phenotypeSearchContext) int {
	return F.Pipe2(
		s.query,
		strings.TrimSpace,
		utf8.RuneCountInString,
	)
}

var (
	// hasMinPhenotypeQueryRunes holds when the trimmed query reaches
	// the length that modware-annotation accepts for a full search.
	hasMinPhenotypeQueryRunes = F.Pipe2(
		minAnnotationSearchQueryRunes,
		ORD.Geq(predord.IntOrd),
		P.ContraMap(trimmedPhenotypeQueryRuneCount),
	)

	// errShortPhenotypeSearchQuery is the client-facing error for a
	// query below the documented minimum, with the message the
	// annotation service uses.
	errShortPhenotypeSearchQuery = status.Error(
		codes.InvalidArgument,
		"expect a search query of at least 2 characters",
	)

	// checkPhenotypeSearchQuery rejects a query below the documented
	// minimum with codes.InvalidArgument.
	checkPhenotypeSearchQuery = E.FromPredicate(
		hasMinPhenotypeQueryRunes,
		F.Constant1[phenotypeSearchContext](errShortPhenotypeSearchQuery),
	)
)

// ── request building ────────────────────────────────────────────────────────

// buildPhenotypeSearchParameters builds the protobuf request from the
// caller-supplied state. The limit passes 0 when the query omits it, so
// the service applies its default of 50, and it never exceeds 50,
// because PhenotypeSearchAttributes enforces lte: 50 in protovalidate.
func buildPhenotypeSearchParameters(
	s phenotypeSearchContext,
) *anno.PhenotypeSearchParameters {
	return &anno.PhenotypeSearchParameters{
		Data: &anno.PhenotypeSearchParameters_Data{
			Type: "phenotype",
			Attributes: &anno.PhenotypeSearchAttributes{
				Query: s.query,
				Limit: resolveAnnotationSearchLimit(s.limit),
			},
		},
	}
}

// ── backend calls ───────────────────────────────────────────────────────────

// fetchPhenotypeSearchCollection calls the SearchPhenotypes rpc for the
// request on the state.
func fetchPhenotypeSearchCollection(
	s phenotypeSearchContext,
) IOE.IOEither[error, *anno.PhenotypeSearchResultCollection] {
	return F.Pipe1(
		IOE.TryCatchError(func() (*anno.PhenotypeSearchResultCollection, error) {
			return s.client.SearchPhenotypes(s.gctx, s.params)
		}),
		IOE.MapLeft[*anno.PhenotypeSearchResultCollection](func(err error) error {
			return fmt.Errorf("search phenotypes for query %q: %w", s.query, err)
		}),
	)
}

// phenotypeSearchEntryID extracts the strain entry id of one match.
func phenotypeSearchEntryID(item *anno.PhenotypeSearchResult) string {
	return item.GetEntryId()
}

// phenotypeSearchCollectionData returns the matched rows on the state.
func phenotypeSearchCollectionData(
	s phenotypeSearchContext,
) []*anno.PhenotypeSearchResult {
	return s.collection.GetData()
}

// fetchPhenotypeStrainMetadataMap fetches the summary, the strain label,
// and the stock type of the matched strains in one ListStrainsByIds
// call. One flat pipe derives the sorted, deduplicated strain ids, runs
// the rpc, recovers a missing strain list, and enriches errors.
func fetchPhenotypeStrainMetadataMap(
	s phenotypeSearchContext,
) IOE.IOEither[error, map[string]strainSearchMetadata] {
	return F.Pipe9(
		s.collection.GetData(),
		A.Map(phenotypeSearchEntryID),
		A.StrictUniq[string],
		A.Sort(stringOrd),
		IOE.Of[error, []string],
		IOE.Map[error](func(ids []string) *pb.StockIdList {
			return &pb.StockIdList{Id: ids}
		}),
		IOE.Chain(func(req *pb.StockIdList) IOE.IOEither[error, *pb.StrainList] {
			return IOE.TryCatchError(func() (*pb.StrainList, error) {
				return s.stockClient.ListStrainsByIds(s.gctx, req)
			})
		}),
		IOE.OrElse(recoverStrainListNotFound),
		IOE.MapLeft[*pb.StrainList](func(err error) error {
			return fmt.Errorf(
				"fetch strain metadata for query %q: %w",
				s.query,
				err,
			)
		}),
		IOE.Map[error](buildStrainMetaMap),
	)
}

var (
	// isNonEmptyPhenotypeSearchResults holds when the search returned
	// at least one row.
	isNonEmptyPhenotypeSearchResults = predarrays.IsNonEmpty[*anno.PhenotypeSearchResult]()

	hasPhenotypeSearchHits = F.Pipe1(
		isNonEmptyPhenotypeSearchResults,
		P.ContraMap(phenotypeSearchCollectionData),
	)

	// emptyPhenotypeStrainMetadata is the Kleisli arm used when the
	// search returned no rows. It builds the empty metadata map and
	// ignores the state parameter.
	emptyPhenotypeStrainMetadata = F.Pipe2(
		R.Empty[string, strainSearchMetadata](),
		IOE.Of[error, map[string]strainSearchMetadata],
		F.Constant1[phenotypeSearchContext, IOE.IOEither[error, map[string]strainSearchMetadata]],
	)

	// fetchPhenotypeStrainMetadata skips ListStrainsByIds when
	// SearchPhenotypes returned no rows.
	fetchPhenotypeStrainMetadata = F.Pipe1(
		hasPhenotypeSearchHits,
		P.Fold(
			emptyPhenotypeStrainMetadata,
			fetchPhenotypeStrainMetadataMap,
		),
	)
)

// ── result projection ───────────────────────────────────────────────────────

// validatePhenotypeSearchItem pairs one protobuf phenotype search result
// with its field spec. Any unspecified or unknown field value fails.
func validatePhenotypeSearchItem(
	item *anno.PhenotypeSearchResult,
) E.Either[error, phenotypeSearchFieldPair] {
	return F.Pipe3(
		phenotypeFieldSpecs,
		R.Lookup[phenotypeFieldSpec](item.Field),
		E.FromOption[phenotypeFieldSpec](func() error {
			return fmt.Errorf(
				"unexpected phenotype search field %s in phenotype search result",
				item.Field,
			)
		}),
		E.Map[error](func(sp phenotypeFieldSpec) phenotypeSearchFieldPair {
			return PR.MakePair(item, sp)
		}),
	)
}

// buildPhenotypeSearchResult joins one validated hit with its strain
// metadata and projects both onto the shared full search result type.
// A hit whose strain id is absent from metaByID yields None and is
// dropped by A.FilterMap, because a result row cannot show data of a
// strain that no longer exists. The in_stock field stays zero here,
// because the StockSearchResult field resolver computes it.
func buildPhenotypeSearchResult(
	metaByID map[string]strainSearchMetadata,
	p phenotypeSearchFieldPair,
) O.Option[*models.StockSearchResult] {
	item := PR.Head(p)
	spec := PR.Tail(p)

	return F.Pipe2(
		metaByID,
		R.Lookup[strainSearchMetadata](item.EntryId),
		O.Map(func(meta strainSearchMetadata) *models.StockSearchResult {
			return &models.StockSearchResult{
				ID:          item.EntryId,
				Entity:      models.StockEntityTypeStrain,
				Field:       spec.field,
				DisplayText: spec.text(item),
				Score:       item.Score,
				Summary:     meta.summary,
				StrainLabel: meta.strainLabel,
				StockType:   meta.stockType,
			}
		}),
	)
}

// convertPhenotypeSearchResultsFromState validates every matched row,
// then joins the valid rows with their strain metadata. It fails when
// one row carries an unknown field enum.
func convertPhenotypeSearchResultsFromState(
	s phenotypeSearchContext,
) E.Either[error, []*models.StockSearchResult] {
	return F.Pipe2(
		s.collection.GetData(),
		E.TraverseArray(validatePhenotypeSearchItem),
		E.Map[error](A.FilterMap(
			F.Bind1st(buildPhenotypeSearchResult, s.strainMetaByID),
		)),
	)
}

// phenotypeMetaSearchLimit reads the effective limit the service
// reports in meta when it is positive.
func phenotypeMetaSearchLimit(s phenotypeSearchContext) O.Option[int] {
	return F.Pipe2(
		s.collection.GetMeta().GetLimit(),
		intFromInt64,
		O.FromPredicate(isPositiveInt),
	)
}

// effectivePhenotypeSearchLimit picks the limit the service reports in
// meta, with the default of 50 as the fallback. The meta limit is
// authoritative, because the service applies its own default and cap.
var effectivePhenotypeSearchLimit = MO.Match(
	defaultSearchLimit,
	[]O.Kleisli[phenotypeSearchContext, int]{
		phenotypeMetaSearchLimit,
	},
)

// buildPhenotypeSearchResultList projects a completed handler run onto
// the shared full search result list.
func buildPhenotypeSearchResultList(
	s phenotypeSearchContext,
) *models.StockSearchResultList {
	limit := s.effectiveLimit

	return &models.StockSearchResultList{
		Results:    s.results,
		Limit:      &limit,
		TotalCount: len(s.results),
	}
}

// ── boundary fold ───────────────────────────────────────────────────────────

func onPhenotypeSearchListError(
	err error,
) T.Tuple2[error, *models.StockSearchResultList] {
	return T.MakeTuple2(err, &models.StockSearchResultList{})
}

func onPhenotypeSearchListSuccess(
	data *models.StockSearchResultList,
) T.Tuple2[error, *models.StockSearchResultList] {
	return T.MakeTuple2[error](nil, data)
}

// SearchPhenotypes resolves the searchPhenotypes GraphQL query. It
// returns the results of a full phenotype search over strain phenotype
// annotations. The match runs against the cvterm label, the synonyms,
// the definition, and the comments of active terms in the Dicty
// Phenotypes, Dicty Environment, and Dictyostelium Assay ontologies.
// Results are ranked, deduplicated by entry_id, and carry the same data
// as full strain search results.
//
// The query must be at least 2 characters. The limit defaults to 50 on
// the service when omitted and never exceeds 50.
//
//nolint:dupl // shares request shape with SearchStocks
func (qrs *QueryResolver) SearchPhenotypes(
	ctx context.Context,
	query string,
	limit *int,
) (*models.StockSearchResultList, error) {
	result := F.Pipe9(
		IOE.Of[error](phenotypeSearchContext{
			client:      qrs.GetAnnotationClient(registry.ANNOTATION),
			stockClient: qrs.GetStockClient(registry.STOCK),
			gctx:        ctx,
			query:       query,
			limit:       limit,
		}),
		IOE.ChainEitherK(checkPhenotypeSearchQuery),
		IOE.Let[error](
			phenotypeSearchParamsLens.Set,
			buildPhenotypeSearchParameters,
		),
		IOE.Bind(
			phenotypeSearchCollectionLens.Set,
			fetchPhenotypeSearchCollection,
		),
		IOE.Bind(phenotypeStrainMetaLens.Set, fetchPhenotypeStrainMetadata),
		IOE.ChainEitherK(F.Flow2(
			E.Of[error, phenotypeSearchContext],
			E.Bind(
				phenotypeSearchResultsLens.Set,
				convertPhenotypeSearchResultsFromState,
			),
		)),
		IOE.Let[error](
			phenotypeSearchLimitLens.Set,
			effectivePhenotypeSearchLimit,
		),
		IOE.Map[error](buildPhenotypeSearchResultList),
		ioeutils.ToEither[error, *models.StockSearchResultList],
		E.Fold(onPhenotypeSearchListError, onPhenotypeSearchListSuccess),
	)

	if result.F1 != nil {
		errorutils.AddGQLError(ctx, result.F1)
		qrs.Logger.Error(result.F1)
		return result.F2, result.F1
	}

	qrs.Logger.Debugf(
		"successfully retrieved %d phenotype search results",
		result.F2.TotalCount,
	)
	return result.F2, nil
}
