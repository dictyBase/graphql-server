package resolver

import (
	"context"
	"fmt"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	ORD "github.com/IBM/fp-go/v2/ord"
	P "github.com/IBM/fp-go/v2/predicate"
	R "github.com/IBM/fp-go/v2/record"
	T "github.com/IBM/fp-go/v2/tuple"
	MO "github.com/dictyBase/fp-go-loom/matchopt"
	predord "github.com/dictyBase/fp-go-loom/predicate/ord"
	anno "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
)

const defaultSuggestionLimit = 5

// listStockSuggestionsContext carries the autocomplete request through
// the suggestion pipeline. Only the caller-supplied fields are set at
// the seed; all intermediate and derived fields are attached by lenses
// so every pipeline step receives a single state parameter.
type listStockSuggestionsContext struct {
	client            pb.StockServiceClient
	annoClient        anno.TaggedAnnotationServiceClient
	gctx              context.Context
	query             string
	limit             *int
	entity            *models.StockEntityType
	params            *pb.StockAutocompleteParameters
	annoParams        *anno.AnnotationAutocompleteParameters
	collection        *pb.StockSuggestionCollection
	annoCollection    *anno.AnnotationSuggestionCollection
	stockSuggestions  []*models.StockSuggestion
	annoSuggestions   []*models.StockSuggestion
	effectiveLimit    int
	mergedSuggestions []*models.StockSuggestion
}

type (
	stockSuggestionEnums = T.Tuple2[models.StockEntityType, models.StockSearchFieldEnum]
)

// ── state lenses ────────────────────────────────────────────────────────────

var (
	suggestionParamsLens = L.MakeLens(
		func(s listStockSuggestionsContext) *pb.StockAutocompleteParameters {
			return s.params
		},
		func(
			s listStockSuggestionsContext,
			params *pb.StockAutocompleteParameters,
		) listStockSuggestionsContext {
			s.params = params
			return s
		},
	)
	suggestionCollectionLens = L.MakeLens(
		func(s listStockSuggestionsContext) *pb.StockSuggestionCollection {
			return s.collection
		},
		func(
			s listStockSuggestionsContext,
			coll *pb.StockSuggestionCollection,
		) listStockSuggestionsContext {
			s.collection = coll
			return s
		},
	)
	annoSuggestionParamsLens = L.MakeLens(
		func(s listStockSuggestionsContext) *anno.AnnotationAutocompleteParameters {
			return s.annoParams
		},
		func(
			s listStockSuggestionsContext,
			params *anno.AnnotationAutocompleteParameters,
		) listStockSuggestionsContext {
			s.annoParams = params
			return s
		},
	)
	annoSuggestionCollectionLens = L.MakeLens(
		func(s listStockSuggestionsContext) *anno.AnnotationSuggestionCollection {
			return s.annoCollection
		},
		func(
			s listStockSuggestionsContext,
			coll *anno.AnnotationSuggestionCollection,
		) listStockSuggestionsContext {
			s.annoCollection = coll
			return s
		},
	)
	stockSuggestionsLens = L.MakeLens(
		func(s listStockSuggestionsContext) []*models.StockSuggestion {
			return s.stockSuggestions
		},
		func(
			s listStockSuggestionsContext,
			items []*models.StockSuggestion,
		) listStockSuggestionsContext {
			s.stockSuggestions = items
			return s
		},
	)
	annoSuggestionsLens = L.MakeLens(
		func(s listStockSuggestionsContext) []*models.StockSuggestion {
			return s.annoSuggestions
		},
		func(
			s listStockSuggestionsContext,
			items []*models.StockSuggestion,
		) listStockSuggestionsContext {
			s.annoSuggestions = items
			return s
		},
	)
	effectiveLimitLens = L.MakeLens(
		func(s listStockSuggestionsContext) int {
			return s.effectiveLimit
		},
		func(
			s listStockSuggestionsContext,
			limit int,
		) listStockSuggestionsContext {
			s.effectiveLimit = limit
			return s
		},
	)
	mergedSuggestionsLens = L.MakeLens(
		func(s listStockSuggestionsContext) []*models.StockSuggestion {
			return s.mergedSuggestions
		},
		func(
			s listStockSuggestionsContext,
			items []*models.StockSuggestion,
		) listStockSuggestionsContext {
			s.mergedSuggestions = items
			return s
		},
	)
)

// ── enum conversions ────────────────────────────────────────────────────────

var (
	// gqlToProtoEntity maps the GraphQL entity filter to the protobuf one.
	// ALL and an absent filter both cover both kinds of stock, which the
	// proto encodes as STOCK_ENTITY_UNSPECIFIED.
	gqlToProtoEntity = map[models.StockEntityType]pb.StockEntity{
		models.StockEntityTypeAll:     pb.StockEntity_STOCK_ENTITY_UNSPECIFIED,
		models.StockEntityTypeStrain:  pb.StockEntity_STOCK_ENTITY_STRAIN,
		models.StockEntityTypePlasmid: pb.StockEntity_STOCK_ENTITY_PLASMID,
	}

	// protoToGQLEntity maps a matched stock kind back to the GraphQL enum.
	protoToGQLEntity = map[pb.StockEntity]models.StockEntityType{
		pb.StockEntity_STOCK_ENTITY_STRAIN:  models.StockEntityTypeStrain,
		pb.StockEntity_STOCK_ENTITY_PLASMID: models.StockEntityTypePlasmid,
	}

	// protoToGQLSearchField maps a matched field back to the GraphQL enum.
	protoToGQLSearchField = map[pb.StockSearchField]models.StockSearchFieldEnum{
		pb.StockSearchField_STOCK_SEARCH_FIELD_STOCK_ID:  models.StockSearchFieldEnumStockID,
		pb.StockSearchField_STOCK_SEARCH_FIELD_GENES:     models.StockSearchFieldEnumGenes,
		pb.StockSearchField_STOCK_SEARCH_FIELD_DBXREFS:   models.StockSearchFieldEnumDbxrefs,
		pb.StockSearchField_STOCK_SEARCH_FIELD_LABEL:     models.StockSearchFieldEnumLabel,
		pb.StockSearchField_STOCK_SEARCH_FIELD_NAMES:     models.StockSearchFieldEnumNames,
		pb.StockSearchField_STOCK_SEARCH_FIELD_SPECIES:   models.StockSearchFieldEnumSpecies,
		pb.StockSearchField_STOCK_SEARCH_FIELD_PLASMID:   models.StockSearchFieldEnumPlasmid,
		pb.StockSearchField_STOCK_SEARCH_FIELD_NAME:      models.StockSearchFieldEnumName,
		pb.StockSearchField_STOCK_SEARCH_FIELD_SUMMARY:   models.StockSearchFieldEnumSummary,
		pb.StockSearchField_STOCK_SEARCH_FIELD_DEPOSITOR: models.StockSearchFieldEnumDepositor,
	}

	// protoToGQLAnnotationSearchField maps a matched annotation field
	// back to the GraphQL enum.
	protoToGQLAnnotationSearchField = map[anno.AnnotationSearchField]models.StockSearchFieldEnum{
		anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_SYNONYM:  models.StockSearchFieldEnumSynonym,
		anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_GENOTYPE: models.StockSearchFieldEnumGenotype,
	}
)

// toKnownProtoEntity resolves a GraphQL entity filter against the lookup
// table. A filter missing from the table fails instead of silently falling
// back, so an enum drift between the GraphQL schema and the table surfaces
// as an error.
func toKnownProtoEntity(
	entity models.StockEntityType,
) E.Either[error, pb.StockEntity] {
	return F.Pipe2(
		gqlToProtoEntity,
		R.Lookup[pb.StockEntity](entity),
		E.FromOption[pb.StockEntity](func() error {
			return fmt.Errorf("unknown stock entity filter %s", entity)
		}),
	)
}

// toGQLEntity maps a matched stock kind back to the GraphQL enum. The
// stock service contract says it never returns an unspecified kind, so an
// unknown value is an error rather than a silent default.
func toGQLEntity(
	entity pb.StockEntity,
) E.Either[error, models.StockEntityType] {
	return F.Pipe2(
		protoToGQLEntity,
		R.Lookup[models.StockEntityType](entity),
		E.FromOption[models.StockEntityType](func() error {
			return fmt.Errorf(
				"unexpected stock entity %s in stock suggestion",
				entity,
			)
		}),
	)
}

// toGQLSearchField maps a matched field back to the GraphQL enum. The
// stock service contract says it never returns an unspecified field, so an
// unknown value is an error rather than a silent default.
func toGQLSearchField(
	field pb.StockSearchField,
) E.Either[error, models.StockSearchFieldEnum] {
	return F.Pipe2(
		protoToGQLSearchField,
		R.Lookup[models.StockSearchFieldEnum](field),
		E.FromOption[models.StockSearchFieldEnum](func() error {
			return fmt.Errorf(
				"unexpected search field %s in stock suggestion",
				field,
			)
		}),
	)
}

// toGQLAnnotationSearchField maps a matched annotation field back to the
// GraphQL enum. AutocompleteAnnotation only emits SYNONYM and GENOTYPE, so
// any other value is an error rather than a silent default.
func toGQLAnnotationSearchField(
	field anno.AnnotationSearchField,
) E.Either[error, models.StockSearchFieldEnum] {
	return F.Pipe2(
		protoToGQLAnnotationSearchField,
		R.Lookup[models.StockSearchFieldEnum](field),
		E.FromOption[models.StockSearchFieldEnum](func() error {
			return fmt.Errorf(
				"unexpected annotation search field %s in stock suggestion",
				field,
			)
		}),
	)
}

func toInt64(value int) int64 { return int64(value) }

// zeroSuggestionLimit lets the stock service apply its own default of 5
// when the GraphQL query omits the limit.
var zeroSuggestionLimit = int64(0)

// resolveSuggestionLimit converts the optional GraphQL limit. Zero lets the
// stock service apply its own default of 5, so an absent limit passes zero.
func resolveSuggestionLimit(limit *int) int64 {
	return F.Pipe3(
		limit,
		O.FromNillable2[int],
		O.Map(toInt64),
		O.GetOrElse(F.Constant(zeroSuggestionLimit)),
	)
}

// ── result helpers ──────────────────────────────────────────────────────────

func onStockSuggestionListError(
	err error,
) T.Tuple2[error, *models.StockSuggestionList] {
	return T.MakeTuple2(err, &models.StockSuggestionList{})
}

func onStockSuggestionListSuccess(
	data *models.StockSuggestionList,
) T.Tuple2[error, *models.StockSuggestionList] {
	return T.MakeTuple2[error](nil, data)
}

// ── curried builders ───────────────────────────────────────────────────────

// buildStockSuggestion is a curried builder: partial application with the
// stock item yields a function of the resolved (entity, field) enums, so the
// pipeline stays data-last without a pair carrier.
var buildStockSuggestion = F.Curry2(
	func(
		item *pb.StockSuggestion,
		enums stockSuggestionEnums,
	) *models.StockSuggestion {
		return &models.StockSuggestion{
			ID:          item.Id,
			Entity:      enums.F1,
			Field:       enums.F2,
			DisplayText: item.DisplayText,
			Score:       item.Score,
		}
	},
)

// convertSuggestionItem projects one protobuf suggestion onto the GraphQL
// type by sequencing the resolved (entity, field) enum lookups into a tuple
// Either and mapping the curried builder over it.
func convertSuggestionItem(
	item *pb.StockSuggestion,
) E.Either[error, *models.StockSuggestion] {
	entity := toGQLEntity(item.Entity)
	field := toGQLSearchField(item.Field)

	return F.Pipe1(
		E.SequenceT2(entity, field),
		E.Map[error](buildStockSuggestion(item)),
	)
}

// buildAnnotationStockSuggestion is a curried builder: partial application
// with the annotation item yields a function of the resolved field enum, so
// the pipeline stays data-last without a pair carrier.
var buildAnnotationStockSuggestion = F.Curry2(
	func(
		item *anno.AnnotationSuggestion,
		field models.StockSearchFieldEnum,
	) *models.StockSuggestion {
		return &models.StockSuggestion{
			ID:          item.EntryId,
			Entity:      models.StockEntityTypeStrain,
			Field:       field,
			DisplayText: item.DisplayText,
			Score:       item.Score,
		}
	},
)

// convertAnnotationSuggestionItem projects one protobuf annotation suggestion
// onto the GraphQL StockSuggestion type by resolving its field enum and
// mapping the curried builder over it.
func convertAnnotationSuggestionItem(
	item *anno.AnnotationSuggestion,
) E.Either[error, *models.StockSuggestion] {
	return F.Pipe2(
		item.Field,
		toGQLAnnotationSearchField,
		E.Map[error](buildAnnotationStockSuggestion(item)),
	)
}

// ── pipeline steps ──────────────────────────────────────────────────────────

func resolveProtoEntity(
	entity *models.StockEntityType,
) E.Either[error, pb.StockEntity] {
	return F.Pipe3(
		entity,
		O.FromNillable2[models.StockEntityType],
		O.GetOrElse(F.Constant(models.StockEntityTypeAll)),
		toKnownProtoEntity,
	)
}

// buildSuggestionParameters is a curried builder: partial application with
// the request state yields a function of the resolved proto entity, so the
// pipeline stays data-last without a pair carrier.
var buildSuggestionParameters = F.Curry2(
	func(
		s listStockSuggestionsContext,
		entity pb.StockEntity,
	) *pb.StockAutocompleteParameters {
		return &pb.StockAutocompleteParameters{
			Data: &pb.StockAutocompleteParameters_Data{
				Type: "stock",
				Attributes: &pb.StockAutocompleteAttributes{
					Query:  s.query,
					Limit:  resolveSuggestionLimit(s.limit),
					Entity: entity,
				},
			},
		}
	},
)

// suggestionParamsFromState derives the stock request params from the caller
// supplied fields by resolving the proto entity and mapping the curried
// builder over it.
func suggestionParamsFromState(
	s listStockSuggestionsContext,
) E.Either[error, *pb.StockAutocompleteParameters] {
	return F.Pipe2(
		s.entity,
		resolveProtoEntity,
		E.Map[error](buildSuggestionParameters(s)),
	)
}

func buildAnnotationSuggestionParameters(
	s listStockSuggestionsContext,
) *anno.AnnotationAutocompleteParameters {
	return &anno.AnnotationAutocompleteParameters{
		Data: &anno.AnnotationAutocompleteParameters_Data{
			Type: "annotation",
			Attributes: &anno.AnnotationAutocompleteAttributes{
				Query: s.query,
				Limit: resolveSuggestionLimit(s.limit),
			},
		},
	}
}

// withSuggestionParams attaches both the stock and annotation request
// parameters to the state in one flat Either pipe, lifted into the IOEither
// pipeline with IOE.ChainEitherK.
func withSuggestionParams(
	s listStockSuggestionsContext,
) E.Either[error, listStockSuggestionsContext] {
	return F.Pipe3(
		s,
		E.Of[error, listStockSuggestionsContext],
		E.Bind(suggestionParamsLens.Set, suggestionParamsFromState),
		E.Let[error](
			annoSuggestionParamsLens.Set,
			buildAnnotationSuggestionParameters,
		),
	)
}

func fetchSuggestionCollection(
	s listStockSuggestionsContext,
) IOE.IOEither[error, *pb.StockSuggestionCollection] {
	rawFetch := func() (*pb.StockSuggestionCollection, error) {
		return s.client.AutocompleteStock(s.gctx, s.params)
	}

	return F.Pipe1(
		IOE.TryCatchError(rawFetch),
		IOE.MapLeft[*pb.StockSuggestionCollection](func(err error) error {
			return fmt.Errorf("fetch stock suggestions for query %q: %w", s.query, err)
		}),
	)
}

func fetchAnnotationSuggestionCollection(
	s listStockSuggestionsContext,
) IOE.IOEither[error, *anno.AnnotationSuggestionCollection] {
	rawFetch := func() (*anno.AnnotationSuggestionCollection, error) {
		return s.annoClient.AutocompleteAnnotation(s.gctx, s.annoParams)
	}

	return F.Pipe1(
		IOE.TryCatchError(rawFetch),
		IOE.MapLeft[*anno.AnnotationSuggestionCollection](func(err error) error {
			return fmt.Errorf(
				"fetch annotation suggestions for query %q: %w",
				s.query,
				err,
			)
		}),
	)
}

func emptyAnnotationSuggestionCollection(
	s listStockSuggestionsContext,
) IOE.IOEither[error, *anno.AnnotationSuggestionCollection] {
	return IOE.Of[error](&anno.AnnotationSuggestionCollection{
		Data: []*anno.AnnotationSuggestion{},
		Meta: &anno.Meta{
			Limit: s.annoParams.Data.Attributes.Limit,
			Total: 0,
		},
	})
}

func suggestionProtoEntity(s listStockSuggestionsContext) pb.StockEntity {
	return s.params.Data.Attributes.Entity
}

var (
	protoEntityEq = EQ.FromStrictEquals[pb.StockEntity]()

	isPlasmidProtoEntity = F.Pipe1(
		pb.StockEntity_STOCK_ENTITY_PLASMID,
		EQ.Equals(protoEntityEq),
	)

	isPlasmidSuggestionRequest = F.Pipe1(
		isPlasmidProtoEntity,
		P.ContraMap(suggestionProtoEntity),
	)

	// fetchAnnotationSuggestionsForEntity skips AutocompleteAnnotation when
	// the request filters to PLASMID, because annotation autocomplete only
	// indexes strain synonyms and genotypes.
	fetchAnnotationSuggestionsForEntity = F.Pipe1(
		isPlasmidSuggestionRequest,
		P.Fold(
			fetchAnnotationSuggestionCollection,
			emptyAnnotationSuggestionCollection,
		),
	)
)

func convertStockSuggestionsFromState(
	s listStockSuggestionsContext,
) E.Either[error, []*models.StockSuggestion] {
	return F.Pipe1(
		s.collection.GetData(),
		E.TraverseArray(convertSuggestionItem),
	)
}

func convertAnnoSuggestionsFromState(
	s listStockSuggestionsContext,
) E.Either[error, []*models.StockSuggestion] {
	return F.Pipe1(
		s.annoCollection.GetData(),
		E.TraverseArray(convertAnnotationSuggestionItem),
	)
}

// ── ordering and limit helpers ───────────────────────────────────────────────

func suggestionScore(item *models.StockSuggestion) float64 {
	return item.Score
}

func suggestionID(item *models.StockSuggestion) string {
	return item.ID
}

func suggestionField(item *models.StockSuggestion) string {
	return item.Field.String()
}

func suggestionDisplayText(item *models.StockSuggestion) string {
	return item.DisplayText
}

func intFromInt64(value int64) int {
	return int(value)
}

func requestedSuggestionLimit(s listStockSuggestionsContext) O.Option[int] {
	return F.Pipe2(
		s.limit,
		O.FromNillable2[int],
		O.Chain(O.FromPredicate(isPositiveInt)),
	)
}

func metaSuggestionLimit(s listStockSuggestionsContext) O.Option[int] {
	return F.Pipe2(
		s.collection.GetMeta().GetLimit(),
		intFromInt64,
		O.FromPredicate(isPositiveInt),
	)
}

var (
	isPositiveInt = F.Pipe1(0, ORD.Gt(predord.IntOrd))

	effectiveSuggestionLimit = MO.Match(
		defaultSuggestionLimit,
		[]O.Kleisli[listStockSuggestionsContext, int]{
			requestedSuggestionLimit,
			metaSuggestionLimit,
		},
	)

	stringOrd = ORD.FromStrictCompare[string]()

	byScoreDesc = F.Pipe2(
		predord.Float64Ord,
		ORD.ContraMap(suggestionScore),
		ORD.Reverse[*models.StockSuggestion],
	)
	byIDAsc = F.Pipe1(
		stringOrd,
		ORD.ContraMap(suggestionID),
	)
	byFieldAsc = F.Pipe1(
		stringOrd,
		ORD.ContraMap(suggestionField),
	)
	byDisplayTextAsc = F.Pipe1(
		stringOrd,
		ORD.ContraMap(suggestionDisplayText),
	)
	suggestionOrderings = []ORD.Ord[*models.StockSuggestion]{
		byScoreDesc,
		byIDAsc,
		byFieldAsc,
		byDisplayTextAsc,
	}
)

// mergeSuggestionsFromState concatenates stock and annotation suggestions on
// the state struct, sorts them by (score DESC, id ASC, field ASC, display_text
// ASC), deduplicates by stock ID keeping the highest-ranked suggestion, and
// truncates to s.effectiveLimit.
func mergeSuggestionsFromState(
	s listStockSuggestionsContext,
) []*models.StockSuggestion {
	return F.Pipe5(
		s.stockSuggestions,
		A.Concat(s.annoSuggestions),
		A.SortBy(suggestionOrderings),
		A.Uniq(suggestionID),
		A.Slice[*models.StockSuggestion](0, s.effectiveLimit),
		A.Match(
			A.Empty[*models.StockSuggestion],
			F.Identity[[]*models.StockSuggestion],
		),
	)
}

func buildMergedSuggestionList(
	s listStockSuggestionsContext,
) *models.StockSuggestionList {
	limit := s.effectiveLimit

	return &models.StockSuggestionList{
		Suggestions: s.mergedSuggestions,
		Limit:       &limit,
		TotalCount:  len(s.mergedSuggestions),
	}
}

// toStockSuggestionList converts both collections, resolves the effective
// limit, merges and deduplicates the suggestions, and projects the final
// GraphQL list in a single flat point-free Either pipe, chained into the
// IOEither pipeline via IOE.ChainEitherK.
func toStockSuggestionList(
	s listStockSuggestionsContext,
) E.Either[error, *models.StockSuggestionList] {
	return F.Pipe6(
		s,
		E.Of[error, listStockSuggestionsContext],
		E.Bind(stockSuggestionsLens.Set, convertStockSuggestionsFromState),
		E.Bind(annoSuggestionsLens.Set, convertAnnoSuggestionsFromState),
		E.Let[error](effectiveLimitLens.Set, effectiveSuggestionLimit),
		E.Let[error](mergedSuggestionsLens.Set, mergeSuggestionsFromState),
		E.Map[error](buildMergedSuggestionList),
	)
}
