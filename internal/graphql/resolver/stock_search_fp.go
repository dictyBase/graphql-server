package resolver

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	ORD "github.com/IBM/fp-go/v2/ord"
	PR "github.com/IBM/fp-go/v2/pair"
	P "github.com/IBM/fp-go/v2/predicate"
	R "github.com/IBM/fp-go/v2/record"
	T "github.com/IBM/fp-go/v2/tuple"
	MO "github.com/dictyBase/fp-go-loom/matchopt"
	predarrays "github.com/dictyBase/fp-go-loom/predicate/array"
	predord "github.com/dictyBase/fp-go-loom/predicate/ord"
	anno "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/graphql/resolverutils"
	"github.com/dictyBase/graphql-server/internal/registry"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	defaultSearchLimit       = 50
	maxAnnotationSearchLimit = int64(50)
	// bacterialSearchAnnoLimit caps the single ListAnnotations query that
	// resolves bacterial food source strains for full search results.
	bacterialSearchAnnoLimit = int64(100)
	// minAnnotationSearchQueryRunes is the minimum trimmed query length, in
	// runes, that modware-annotation accepts for a full search.
	minAnnotationSearchQueryRunes = 2
	// searchStockTypeRegular is the GraphQL stock_type for non-GWDI,
	// non-bacterial, and non-GoldenBraid stocks.
	searchStockTypeRegular = "regular"
	// searchStockTypeGwdi is the GraphQL stock_type for REMI-seq strains.
	searchStockTypeGwdi = "gwdi"
	// searchStockTypeBacterial is the GraphQL stock_type for bacterial
	// strains.
	searchStockTypeBacterial = "bacterial"
	// searchStockTypeGoldenBraid is the GraphQL stock_type for GB vector
	// plasmids.
	searchStockTypeGoldenBraid = "goldenbraid"
)

// strainSearchMetadata holds the stored summary, strain label, and mapped
// stock type fetched in batch from the stock service for annotation search
// hits.
type strainSearchMetadata struct {
	summary     string
	strainLabel string
	stockType   string
}

// listStockSearchContext carries the full search request through the
// search pipeline. Only the caller-supplied fields are set at the seed;
// all intermediate and derived fields are attached by lenses so every
// pipeline step receives a single state parameter.
type listStockSearchContext struct {
	client         pb.StockServiceClient
	annoClient     anno.TaggedAnnotationServiceClient
	gctx           context.Context
	query          string
	limit          *int
	entity         *models.StockEntityType
	params         *pb.StockSearchParameters
	annoParams     *anno.AnnotationSearchParameters
	collection     *pb.StockSearchResultCollection
	annoCollection *anno.AnnotationSearchResultCollection
	strainMetaByID map[string]strainSearchMetadata
	stockResults   []*models.StockSearchResult
	annoResults    []*models.StockSearchResult
	effectiveLimit int
	mergedResults  []*models.StockSearchResult
	bacterialIDs   map[string]bool
}

type (
	stockSearchEnums    = T.Tuple2[models.StockEntityType, models.StockSearchFieldEnum]
	stockSearchPair     = PR.Pair[*pb.StockSearchResult, stockSearchEnums]
	annoSearchFieldPair = PR.Pair[*anno.AnnotationSearchResult, models.StockSearchFieldEnum]
	annoSearchMetaPair  = PR.Pair[annoSearchFieldPair, strainSearchMetadata]
	searchParamPair     = PR.Pair[listStockSearchContext, pb.StockEntity]
	strainMetaEntry     = R.Entry[string, strainSearchMetadata]
)

// ── state lenses ────────────────────────────────────────────────────────────

var (
	searchParamsLens = L.MakeLens(
		func(s listStockSearchContext) *pb.StockSearchParameters {
			return s.params
		},
		func(
			s listStockSearchContext,
			params *pb.StockSearchParameters,
		) listStockSearchContext {
			s.params = params
			return s
		},
	)
	searchCollectionLens = L.MakeLens(
		func(s listStockSearchContext) *pb.StockSearchResultCollection {
			return s.collection
		},
		func(
			s listStockSearchContext,
			coll *pb.StockSearchResultCollection,
		) listStockSearchContext {
			s.collection = coll
			return s
		},
	)
	annoSearchParamsLens = L.MakeLens(
		func(s listStockSearchContext) *anno.AnnotationSearchParameters {
			return s.annoParams
		},
		func(
			s listStockSearchContext,
			params *anno.AnnotationSearchParameters,
		) listStockSearchContext {
			s.annoParams = params
			return s
		},
	)
	annoSearchCollectionLens = L.MakeLens(
		func(s listStockSearchContext) *anno.AnnotationSearchResultCollection {
			return s.annoCollection
		},
		func(
			s listStockSearchContext,
			coll *anno.AnnotationSearchResultCollection,
		) listStockSearchContext {
			s.annoCollection = coll
			return s
		},
	)
	strainMetaByIDLens = L.MakeLens(
		func(s listStockSearchContext) map[string]strainSearchMetadata {
			return s.strainMetaByID
		},
		func(
			s listStockSearchContext,
			meta map[string]strainSearchMetadata,
		) listStockSearchContext {
			s.strainMetaByID = meta
			return s
		},
	)
	stockResultsLens = L.MakeLens(
		func(s listStockSearchContext) []*models.StockSearchResult {
			return s.stockResults
		},
		func(
			s listStockSearchContext,
			items []*models.StockSearchResult,
		) listStockSearchContext {
			s.stockResults = items
			return s
		},
	)
	annoResultsLens = L.MakeLens(
		func(s listStockSearchContext) []*models.StockSearchResult {
			return s.annoResults
		},
		func(
			s listStockSearchContext,
			items []*models.StockSearchResult,
		) listStockSearchContext {
			s.annoResults = items
			return s
		},
	)
	effectiveSearchLimitLens = L.MakeLens(
		func(s listStockSearchContext) int {
			return s.effectiveLimit
		},
		func(
			s listStockSearchContext,
			limit int,
		) listStockSearchContext {
			s.effectiveLimit = limit
			return s
		},
	)
	mergedResultsLens = L.MakeLens(
		func(s listStockSearchContext) []*models.StockSearchResult {
			return s.mergedResults
		},
		func(
			s listStockSearchContext,
			items []*models.StockSearchResult,
		) listStockSearchContext {
			s.mergedResults = items
			return s
		},
	)
	bacterialIDsLens = L.MakeLens(
		func(s listStockSearchContext) map[string]bool {
			return s.bacterialIDs
		},
		func(
			s listStockSearchContext,
			ids map[string]bool,
		) listStockSearchContext {
			s.bacterialIDs = ids
			return s
		},
	)
)

// ── enum conversions ────────────────────────────────────────────────────────

// Proto annotation search fields emitted by SearchAnnotation.
const (
	annoFieldSynonym     = anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_SYNONYM
	annoFieldGenotype    = anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_GENOTYPE
	annoFieldSysName     = anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_SYSTEMATIC_NAME
	annoFieldGenMod      = anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_GENETIC_MODIFICATION
	annoFieldMutMethod   = anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_MUTAGENESIS_METHOD
	annoFieldPhenotype   = anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_PHENOTYPE
	annoFieldEnvironment = anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_ENVIRONMENT
	annoFieldAssay       = anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_ASSAY
	annoFieldChar        = anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_CHARACTERISTIC
	annoFieldPublicNote  = anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_PUBLIC_NOTE
)

// protoToGQLFullAnnotationSearchField maps all 10 full-search annotation
// fields onto their GraphQL StockSearchFieldEnum counterparts.
var protoToGQLFullAnnotationSearchField = map[anno.AnnotationSearchField]models.StockSearchFieldEnum{
	annoFieldSynonym:     models.StockSearchFieldEnumSynonym,
	annoFieldGenotype:    models.StockSearchFieldEnumGenotype,
	annoFieldSysName:     models.StockSearchFieldEnumSystematicName,
	annoFieldGenMod:      models.StockSearchFieldEnumGeneticModification,
	annoFieldMutMethod:   models.StockSearchFieldEnumMutagenesisMethod,
	annoFieldPhenotype:   models.StockSearchFieldEnumPhenotype,
	annoFieldEnvironment: models.StockSearchFieldEnumEnvironment,
	annoFieldAssay:       models.StockSearchFieldEnumAssay,
	annoFieldChar:        models.StockSearchFieldEnumCharacteristic,
	annoFieldPublicNote:  models.StockSearchFieldEnumPublicNote,
}

// ontologyTermToSearchStockType maps special ontology term labels onto their
// GraphQL stock_type strings. Any other term defaults to searchStockTypeRegular.
var ontologyTermToSearchStockType = map[string]string{
	registry.GwdiStrainTag:         searchStockTypeGwdi,
	registry.BacterialStrainTag:    searchStockTypeBacterial,
	registry.GoldenBraidPlasmidTag: searchStockTypeGoldenBraid,
}

// resolveSearchStockType maps an ontology term label onto its GraphQL
// stock_type string, defaulting to "regular".
func resolveSearchStockType(term string) string {
	return F.Pipe2(
		ontologyTermToSearchStockType,
		R.Lookup[string](term),
		O.GetOrElse(F.Constant(searchStockTypeRegular)),
	)
}

// toGQLFullAnnotationSearchField maps a matched annotation full-search field
// back to the GraphQL enum. Any unspecified or unknown value fails.
func toGQLFullAnnotationSearchField(
	field anno.AnnotationSearchField,
) E.Either[error, models.StockSearchFieldEnum] {
	return F.Pipe2(
		protoToGQLFullAnnotationSearchField,
		R.Lookup[models.StockSearchFieldEnum](field),
		E.FromOption[models.StockSearchFieldEnum](func() error {
			return fmt.Errorf(
				"unexpected annotation search field %s in stock search result",
				field,
			)
		}),
	)
}

// ── limit helpers ───────────────────────────────────────────────────────────

var (
	// zeroSearchLimit lets the stock and annotation services apply their own
	// default of 50 when the GraphQL query omits the limit.
	zeroSearchLimit = int64(0)
	// int64Ord orders int64 limits.
	int64Ord = ORD.FromStrictCompare[int64]()
	// isAboveAnnotationSearchMaxLimit holds for limits that exceed the
	// annotation service cap.
	isAboveAnnotationSearchMaxLimit = F.Pipe1(
		maxAnnotationSearchLimit,
		ORD.Gt(int64Ord),
	)
	// clampAnnotationSearchLimit keeps an absent limit at zero, caps limits
	// above maxAnnotationSearchLimit and passes every other limit through.
	clampAnnotationSearchLimit = MO.Match(
		zeroSearchLimit,
		[]O.Kleisli[int64, int64]{
			MO.Const(isAboveAnnotationSearchMaxLimit, maxAnnotationSearchLimit),
			MO.Default(F.Identity[int64]),
		},
	)
)

// resolveSearchLimit converts the optional GraphQL limit. Zero lets the
// stock service apply its own default of 50, so an absent limit passes
// zero.
func resolveSearchLimit(limit *int) int64 {
	return F.Pipe3(
		limit,
		O.FromNillable2[int],
		O.Map(toInt64),
		O.GetOrElse(F.Constant(zeroSearchLimit)),
	)
}

// resolveAnnotationSearchLimit converts the optional GraphQL limit for
// SearchAnnotation, capping values above 50 at maxAnnotationSearchLimit
// because AnnotationSearchAttributes enforces lte: 50 in protovalidate.
func resolveAnnotationSearchLimit(limit *int) int64 {
	return F.Pipe2(
		limit,
		resolveSearchLimit,
		clampAnnotationSearchLimit,
	)
}

// ── result helpers ──────────────────────────────────────────────────────────

func onStockSearchListError(
	err error,
) T.Tuple2[error, *models.StockSearchResultList] {
	return T.MakeTuple2(err, &models.StockSearchResultList{})
}

func onStockSearchListSuccess(
	data *models.StockSearchResultList,
) T.Tuple2[error, *models.StockSearchResultList] {
	return T.MakeTuple2[error](nil, data)
}

// ── PR.Pair builders ────────────────────────────────────────────────────────

func buildStockSearchResult(
	p stockSearchPair,
) *models.StockSearchResult {
	item := PR.Head(p)
	enums := PR.Tail(p)

	return &models.StockSearchResult{
		ID:          item.Id,
		Entity:      enums.F1,
		Field:       enums.F2,
		DisplayText: item.DisplayText,
		Score:       item.Score,
		Summary:     item.Summary,
		StrainLabel: item.StrainLabel,
		StockType:   resolveSearchStockType(item.StockType),
	}
}

// convertSearchResultItem projects one protobuf stock search result onto the
// GraphQL type by pairing the item with its resolved (entity, field) tuple,
// sequencing the Pair into Either, and mapping buildStockSearchResult.
func convertSearchResultItem(
	item *pb.StockSearchResult,
) E.Either[error, *models.StockSearchResult] {
	entity := toGQLEntity(item.Entity)
	field := toGQLSearchField(item.Field)
	enums := E.SequenceT2(entity, field)

	return F.Pipe2(
		PR.MakePair(item, enums),
		PR.Sequence(E.Map[error, stockSearchEnums, stockSearchPair]),
		E.Map[error](buildStockSearchResult),
	)
}

// validateAnnotationSearchItem pairs one protobuf annotation search result
// with its validated GraphQL field enum.
func validateAnnotationSearchItem(
	item *anno.AnnotationSearchResult,
) E.Either[error, annoSearchFieldPair] {
	field := toGQLFullAnnotationSearchField(item.Field)

	return F.Pipe1(
		PR.MakePair(item, field),
		PR.Sequence(E.Map[error, models.StockSearchFieldEnum, annoSearchFieldPair]),
	)
}

func buildAnnotationStockSearchResult(
	p annoSearchMetaPair,
) *models.StockSearchResult {
	fieldPair := PR.Head(p)
	meta := PR.Tail(p)
	item := PR.Head(fieldPair)
	field := PR.Tail(fieldPair)

	return &models.StockSearchResult{
		ID:          item.EntryId,
		Entity:      models.StockEntityTypeStrain,
		Field:       field,
		DisplayText: item.DisplayText,
		Score:       item.Score,
		Summary:     meta.summary,
		StrainLabel: meta.strainLabel,
		StockType:   meta.stockType,
	}
}

// enrichAnnotationSearchPair pairs a validated annotation search hit with its
// strain metadata from metaByID. Hits whose strain ID is absent from metaByID
// yield None and are dropped by A.FilterMap.
func enrichAnnotationSearchPair(
	metaByID map[string]strainSearchMetadata,
	fieldPair annoSearchFieldPair,
) O.Option[*models.StockSearchResult] {
	item := PR.Head(fieldPair)
	metaOpt := F.Pipe1(
		metaByID,
		R.Lookup[strainSearchMetadata](item.EntryId),
	)

	return F.Pipe2(
		PR.MakePair(fieldPair, metaOpt),
		PR.Sequence(O.Map[strainSearchMetadata, annoSearchMetaPair]),
		O.Map(buildAnnotationStockSearchResult),
	)
}

// ── pipeline steps ──────────────────────────────────────────────────────────

func buildSearchParameters(
	p searchParamPair,
) *pb.StockSearchParameters {
	s := PR.Head(p)
	entity := PR.Tail(p)

	return &pb.StockSearchParameters{
		Data: &pb.StockSearchParameters_Data{
			Type: "stock",
			Attributes: &pb.StockSearchAttributes{
				Query:  s.query,
				Limit:  resolveSearchLimit(s.limit),
				Entity: entity,
			},
		},
	}
}

// searchParamsFromState derives the stock search request params from the
// caller-supplied fields by pairing state with the resolved proto entity,
// sequencing the Pair into Either, and mapping buildSearchParameters.
func searchParamsFromState(
	s listStockSearchContext,
) E.Either[error, *pb.StockSearchParameters] {
	entity := resolveProtoEntity(s.entity)

	return F.Pipe2(
		PR.MakePair(s, entity),
		PR.Sequence(E.Map[error, pb.StockEntity, searchParamPair]),
		E.Map[error](buildSearchParameters),
	)
}

func buildAnnotationSearchParameters(
	s listStockSearchContext,
) *anno.AnnotationSearchParameters {
	return &anno.AnnotationSearchParameters{
		Data: &anno.AnnotationSearchParameters_Data{
			Type: "annotation",
			Attributes: &anno.AnnotationSearchAttributes{
				Query: s.query,
				Limit: resolveAnnotationSearchLimit(s.limit),
			},
		},
	}
}

// withSearchParams attaches both the stock and annotation search request
// parameters to the state in one flat Either pipe, lifted into the IOEither
// pipeline with IOE.ChainEitherK.
func withSearchParams(
	s listStockSearchContext,
) E.Either[error, listStockSearchContext] {
	return F.Pipe3(
		s,
		E.Of[error, listStockSearchContext],
		E.Bind(searchParamsLens.Set, searchParamsFromState),
		E.Let[error](
			annoSearchParamsLens.Set,
			buildAnnotationSearchParameters,
		),
	)
}

func fetchSearchCollection(
	s listStockSearchContext,
) IOE.IOEither[error, *pb.StockSearchResultCollection] {
	rawFetch := func() (*pb.StockSearchResultCollection, error) {
		return s.client.SearchStock(s.gctx, s.params)
	}

	return F.Pipe1(
		IOE.TryCatchError(rawFetch),
		IOE.MapLeft[*pb.StockSearchResultCollection](func(err error) error {
			return fmt.Errorf("search stocks for query %q: %w", s.query, err)
		}),
	)
}

func fetchAnnotationSearchCollection(
	s listStockSearchContext,
) IOE.IOEither[error, *anno.AnnotationSearchResultCollection] {
	rawFetch := func() (*anno.AnnotationSearchResultCollection, error) {
		return s.annoClient.SearchAnnotation(s.gctx, s.annoParams)
	}

	return F.Pipe1(
		IOE.TryCatchError(rawFetch),
		IOE.MapLeft[*anno.AnnotationSearchResultCollection](func(err error) error {
			return fmt.Errorf(
				"search annotations for query %q: %w",
				s.query,
				err,
			)
		}),
	)
}

func emptyAnnotationSearchCollection(
	s listStockSearchContext,
) IOE.IOEither[error, *anno.AnnotationSearchResultCollection] {
	return IOE.Of[error](&anno.AnnotationSearchResultCollection{
		Data: []*anno.AnnotationSearchResult{},
		Meta: &anno.Meta{
			Limit: s.annoParams.Data.Attributes.Limit,
			Total: 0,
		},
	})
}

func searchProtoEntity(s listStockSearchContext) pb.StockEntity {
	return s.params.Data.Attributes.Entity
}

// trimmedSearchQueryRuneCount counts the runes of the query once surrounding
// whitespace is trimmed.
func trimmedSearchQueryRuneCount(s listStockSearchContext) int {
	return F.Pipe2(
		s.query,
		strings.TrimSpace,
		utf8.RuneCountInString,
	)
}

var (
	// isPlasmidSearchRequest holds when the request filters to PLASMID.
	isPlasmidSearchRequest = F.Pipe1(
		isPlasmidProtoEntity,
		P.ContraMap(searchProtoEntity),
	)

	// isShortSearchQuery holds when the trimmed query is shorter than
	// modware-annotation accepts, even though modware-stock does.
	isShortSearchQuery = F.Pipe2(
		minAnnotationSearchQueryRunes,
		ORD.Lt(predord.IntOrd),
		P.ContraMap(trimmedSearchQueryRuneCount),
	)

	// skipsAnnotationSearch holds when SearchAnnotation must not be called.
	skipsAnnotationSearch = F.Pipe1(
		isPlasmidSearchRequest,
		P.Or(isShortSearchQuery),
	)

	// fetchAnnotationSearchForEntity skips SearchAnnotation when the request
	// filters to PLASMID, because annotation full search only indexes strain
	// annotations, or when the query is too short for annotation full search,
	// so stock-only hits still come back instead of an InvalidArgument error.
	fetchAnnotationSearchForEntity = F.Pipe1(
		skipsAnnotationSearch,
		P.Fold(
			fetchAnnotationSearchCollection,
			emptyAnnotationSearchCollection,
		),
	)
)

// fetchSearchCollections runs the stock and annotation full-search RPCs
// concurrently and attaches both collections to the state struct.
func fetchSearchCollections(
	s listStockSearchContext,
) IOE.IOEither[error, listStockSearchContext] {
	stockSearch := fetchSearchCollection(s)
	annoSearch := fetchAnnotationSearchForEntity(s)

	return F.Pipe2(
		IOE.Of[error](s),
		IOE.ApSL(searchCollectionLens, stockSearch),
		IOE.ApSL(annoSearchCollectionLens, annoSearch),
	)
}

func annoSearchEntryID(item *anno.AnnotationSearchResult) string {
	return item.GetEntryId()
}

func uniqueAnnotationEntryIDs(
	s listStockSearchContext,
) []string {
	return F.Pipe3(
		s.annoCollection.GetData(),
		A.Map(annoSearchEntryID),
		A.StrictUniq[string],
		A.Sort(stringOrd),
	)
}

func toStrainMetaEntry(item *pb.StrainList_Data) strainMetaEntry {
	id := item.GetId()
	attrs := item.GetAttributes()
	prop := attrs.GetDictyStrainProperty()
	meta := strainSearchMetadata{
		summary:     attrs.GetSummary(),
		strainLabel: attrs.GetLabel(),
		stockType:   resolveSearchStockType(prop),
	}

	return PR.MakePair(id, meta)
}

func buildStrainMetaMap(list *pb.StrainList) map[string]strainSearchMetadata {
	return F.Pipe2(
		list.GetData(),
		A.Map(toStrainMetaEntry),
		R.FromEntries[string, strainSearchMetadata],
	)
}

func fetchStrainMetadataMap(
	s listStockSearchContext,
) IOE.IOEither[error, map[string]strainSearchMetadata] {
	ids := uniqueAnnotationEntryIDs(s)
	req := &pb.StockIdList{Id: ids}
	rawFetch := func() (*pb.StrainList, error) {
		return s.client.ListStrainsByIds(s.gctx, req)
	}

	return F.Pipe3(
		IOE.TryCatchError(rawFetch),
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

func emptyStrainMetadataMap(
	_ listStockSearchContext,
) IOE.IOEither[error, map[string]strainSearchMetadata] {
	empty := R.Empty[string, strainSearchMetadata]()

	return IOE.Of[error](empty)
}

func annoSearchCollectionData(
	s listStockSearchContext,
) []*anno.AnnotationSearchResult {
	return s.annoCollection.GetData()
}

var (
	searchCodeEq = EQ.FromStrictEquals[codes.Code]()

	isNotFoundCode = F.Pipe1(
		codes.NotFound,
		EQ.Equals(searchCodeEq),
	)

	isSearchNotFoundStatus = F.Flow2(
		status.Code,
		isNotFoundCode,
	)

	// recoverStrainListNotFound turns a codes.NotFound response from
	// ListStrainsByIds (returned by modware-stock when all requested strain
	// IDs are missing) into an empty StrainList so stale annotation hits are
	// dropped without failing the search.
	recoverStrainListNotFound = F.Pipe1(
		isSearchNotFoundStatus,
		P.Fold(
			IOE.Left[*pb.StrainList, error],
			func(error) IOE.IOEither[error, *pb.StrainList] {
				return IOE.Of[error](&pb.StrainList{})
			},
		),
	)

	isNonEmptyAnnoSearchResults = predarrays.IsNonEmpty[*anno.AnnotationSearchResult]()

	hasAnnotationSearchHits = F.Pipe1(
		isNonEmptyAnnoSearchResults,
		P.ContraMap(annoSearchCollectionData),
	)

	// fetchStrainMetadataForAnnotationHits skips ListStrainsByIds when
	// SearchAnnotation returned no rows.
	fetchStrainMetadataForAnnotationHits = F.Pipe1(
		hasAnnotationSearchHits,
		P.Fold(
			emptyStrainMetadataMap,
			fetchStrainMetadataMap,
		),
	)
)

func convertStockSearchResultsFromState(
	s listStockSearchContext,
) E.Either[error, []*models.StockSearchResult] {
	return F.Pipe1(
		s.collection.GetData(),
		E.TraverseArray(convertSearchResultItem),
	)
}

func convertAnnoSearchResultsFromState(
	s listStockSearchContext,
) E.Either[error, []*models.StockSearchResult] {
	enrich := F.Bind1st(enrichAnnotationSearchPair, s.strainMetaByID)
	filterEnriched := A.FilterMap(enrich)

	return F.Pipe2(
		s.annoCollection.GetData(),
		E.TraverseArray(validateAnnotationSearchItem),
		E.Map[error](filterEnriched),
	)
}

// ── ordering and limit helpers ──────────────────────────────────────────────

func searchResultScore(item *models.StockSearchResult) float64 {
	return item.Score
}

func searchResultID(item *models.StockSearchResult) string {
	return item.ID
}

func searchResultField(item *models.StockSearchResult) string {
	return item.Field.String()
}

func searchResultDisplayText(item *models.StockSearchResult) string {
	return item.DisplayText
}

func searchResultEntity(item *models.StockSearchResult) models.StockEntityType {
	return item.Entity
}

func searchResultStockType(item *models.StockSearchResult) string {
	return item.StockType
}

func mergedSearchResults(s listStockSearchContext) []*models.StockSearchResult {
	return s.mergedResults
}

func requestedSearchLimit(s listStockSearchContext) O.Option[int] {
	return F.Pipe2(
		s.limit,
		O.FromNillable2[int],
		O.Chain(O.FromPredicate(isPositiveInt)),
	)
}

func metaSearchLimit(s listStockSearchContext) O.Option[int] {
	return F.Pipe2(
		s.collection.GetMeta().GetLimit(),
		intFromInt64,
		O.FromPredicate(isPositiveInt),
	)
}

var (
	effectiveSearchLimit = MO.Match(
		defaultSearchLimit,
		[]O.Kleisli[listStockSearchContext, int]{
			requestedSearchLimit,
			metaSearchLimit,
		},
	)

	bySearchScoreDesc = F.Pipe2(
		predord.Float64Ord,
		ORD.ContraMap(searchResultScore),
		ORD.Reverse[*models.StockSearchResult],
	)
	bySearchIDAsc = F.Pipe1(
		stringOrd,
		ORD.ContraMap(searchResultID),
	)
	bySearchFieldAsc = F.Pipe1(
		stringOrd,
		ORD.ContraMap(searchResultField),
	)
	bySearchDisplayTextAsc = F.Pipe1(
		stringOrd,
		ORD.ContraMap(searchResultDisplayText),
	)
	searchResultOrderings = []ORD.Ord[*models.StockSearchResult]{
		bySearchScoreDesc,
		bySearchIDAsc,
		bySearchFieldAsc,
		bySearchDisplayTextAsc,
	}
)

// mergeSearchResultsFromState concatenates stock and enriched annotation
// search results on the state struct, sorts them by (score DESC, id ASC,
// field ASC, display_text ASC), deduplicates by stock ID keeping the
// highest-ranked result, and truncates to s.effectiveLimit.
func mergeSearchResultsFromState(
	s listStockSearchContext,
) []*models.StockSearchResult {
	return F.Pipe5(
		s.stockResults,
		A.Concat(s.annoResults),
		A.SortBy(searchResultOrderings),
		A.Uniq(searchResultID),
		A.Slice[*models.StockSearchResult](0, s.effectiveLimit),
		A.Match(
			A.Empty[*models.StockSearchResult],
			F.Identity[[]*models.StockSearchResult],
		),
	)
}

var (
	searchEntityEq = EQ.FromStrictEquals[models.StockEntityType]()
	stockTypeStrEq = EQ.FromStrictEquals[string]()

	isStrainSearchResult = F.Pipe1(
		models.StockEntityTypeStrain,
		EQ.Equals(searchEntityEq),
	)
	isRegularStockType = F.Pipe1(
		searchStockTypeRegular,
		EQ.Equals(stockTypeStrEq),
	)
	isStrainResultItem = F.Pipe1(
		isStrainSearchResult,
		P.ContraMap(searchResultEntity),
	)
	isRegularResultItem = F.Pipe1(
		isRegularStockType,
		P.ContraMap(searchResultStockType),
	)
	isRegularStrainResultItem = F.Pipe1(
		isStrainResultItem,
		P.And(isRegularResultItem),
	)
	hasRegularStrainResult   = A.Any(isRegularStrainResultItem)
	hasRegularStrainInMerged = F.Pipe1(
		hasRegularStrainResult,
		P.ContraMap(mergedSearchResults),
	)
	isNotPlasmidSearchRequest = F.Pipe1(
		isPlasmidSearchRequest,
		P.Not,
	)
	needsBacterialStrainLookup = F.Pipe1(
		isNotPlasmidSearchRequest,
		P.And(hasRegularStrainInMerged),
	)
	recoverBacterialListNotFound = F.Pipe1(
		isSearchNotFoundStatus,
		P.Fold(
			IOE.Left[*anno.TaggedAnnotationCollection, error],
			func(error) IOE.IOEither[error, *anno.TaggedAnnotationCollection] {
				return IOE.Of[error](&anno.TaggedAnnotationCollection{})
			},
		),
	)
	fetchBacterialStrainIDsForSearch = F.Pipe1(
		needsBacterialStrainLookup,
		P.Fold(
			emptyBacterialIDMap,
			fetchBacterialStrainIDMap,
		),
	)
)

func toBacterialIDEntry(
	item *anno.TaggedAnnotationCollection_Data,
) R.Entry[string, bool] {
	id := item.GetAttributes().GetEntryId()

	return PR.MakePair(id, true)
}

func buildBacterialIDMap(
	coll *anno.TaggedAnnotationCollection,
) map[string]bool {
	return F.Pipe2(
		coll.GetData(),
		A.Map(toBacterialIDEntry),
		R.FromEntries[string, bool],
	)
}

func emptyBacterialIDMap(
	_ listStockSearchContext,
) IOE.IOEither[error, map[string]bool] {
	empty := R.Empty[string, bool]()

	return IOE.Of[error](empty)
}

func fetchBacterialStrainIDMap(
	s listStockSearchContext,
) IOE.IOEither[error, map[string]bool] {
	req := &anno.ListParameters{
		Limit:  bacterialSearchAnnoLimit,
		Filter: resolverutils.BacterialAnnotationFilter(),
	}
	rawFetch := func() (*anno.TaggedAnnotationCollection, error) {
		return s.annoClient.ListAnnotations(s.gctx, req)
	}

	return F.Pipe3(
		IOE.TryCatchError(rawFetch),
		IOE.OrElse(recoverBacterialListNotFound),
		IOE.MapLeft[*anno.TaggedAnnotationCollection](func(err error) error {
			return fmt.Errorf(
				"fetch bacterial strain annotations for query %q: %w",
				s.query,
				err,
			)
		}),
		IOE.Map[error](buildBacterialIDMap),
	)
}

func withBacterialStockType(
	item *models.StockSearchResult,
) *models.StockSearchResult {
	copied := *item
	copied.StockType = searchStockTypeBacterial

	return &copied
}

func hasBacterialEntryID(
	bacterialIDs map[string]bool,
) P.Predicate[*models.StockSearchResult] {
	hasID := func(id string) bool {
		return F.Pipe2(
			bacterialIDs,
			R.Lookup[bool](id),
			O.GetOrElse(F.Constant(false)),
		)
	}

	return F.Pipe1(
		hasID,
		P.ContraMap(searchResultID),
	)
}

func promoteBacterialStrainItem(
	bacterialIDs map[string]bool,
) func(*models.StockSearchResult) *models.StockSearchResult {
	isMember := hasBacterialEntryID(bacterialIDs)
	isTarget := F.Pipe1(
		isRegularStrainResultItem,
		P.And(isMember),
	)

	return F.Pipe1(
		isTarget,
		P.Fold(
			F.Identity[*models.StockSearchResult],
			withBacterialStockType,
		),
	)
}

func buildMergedSearchResultList(
	s listStockSearchContext,
) *models.StockSearchResultList {
	limit := s.effectiveLimit
	promote := promoteBacterialStrainItem(s.bacterialIDs)
	results := F.Pipe1(
		s.mergedResults,
		A.Map(promote),
	)

	return &models.StockSearchResultList{
		Results:    results,
		Limit:      &limit,
		TotalCount: len(results),
	}
}

// withMergedSearchResults converts both collections, resolves the effective
// limit, and merges and deduplicates the search results on the state struct in
// a single flat point-free Either pipe, chained into the IOEither pipeline via
// IOE.ChainEitherK.
func withMergedSearchResults(
	s listStockSearchContext,
) E.Either[error, listStockSearchContext] {
	return F.Pipe5(
		s,
		E.Of[error, listStockSearchContext],
		E.Bind(stockResultsLens.Set, convertStockSearchResultsFromState),
		E.Bind(annoResultsLens.Set, convertAnnoSearchResultsFromState),
		E.Let[error](effectiveSearchLimitLens.Set, effectiveSearchLimit),
		E.Let[error](mergedResultsLens.Set, mergeSearchResultsFromState),
	)
}
