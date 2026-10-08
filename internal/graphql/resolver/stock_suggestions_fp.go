package resolver

import (
	"context"
	"fmt"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	R "github.com/IBM/fp-go/v2/record"
	T "github.com/IBM/fp-go/v2/tuple"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
)

// listStockSuggestionsContext carries the autocomplete request through
// the suggestion pipeline. Only the caller-supplied fields are set at
// the seed; the params and collection are attached by the lens setters.
type listStockSuggestionsContext struct {
	client     pb.StockServiceClient
	gctx       context.Context
	query      string
	limit      *int
	entity     *models.StockEntityType
	params     *pb.StockAutocompleteParameters
	collection *pb.StockSuggestionCollection
}

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

// convertSuggestionItem projects one protobuf suggestion onto the GraphQL
// type. Both enum fields must resolve, otherwise the whole result fails.
func convertSuggestionItem(
	item *pb.StockSuggestion,
) E.Either[error, *models.StockSuggestion] {
	return F.Pipe1(
		toGQLEntity(item.Entity),
		E.Chain(func(entity models.StockEntityType) E.Either[error, *models.StockSuggestion] {
			return F.Pipe1(
				toGQLSearchField(item.Field),
				E.Map[error](func(field models.StockSearchFieldEnum) *models.StockSuggestion {
					return &models.StockSuggestion{
						ID:          item.Id,
						Entity:      entity,
						Field:       field,
						DisplayText: item.DisplayText,
						Score:       item.Score,
					}
				}),
			)
		}),
	)
}

// ── pipeline steps ──────────────────────────────────────────────────────────

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

// suggestionParamsFromState derives the request params from the caller
// supplied fields. It fails when the entity filter is unknown. The
// Either lookup lifts straight into IOEither, so the resolver pipeline
// stays in IOEither until the forced edge.
func suggestionParamsFromState(
	s listStockSuggestionsContext,
) IOE.IOEither[error, *pb.StockAutocompleteParameters] {
	return F.Pipe5(
		s.entity,
		O.FromNillable2[models.StockEntityType],
		O.GetOrElse(F.Constant(models.StockEntityTypeAll)),
		toKnownProtoEntity,
		E.Map[error](buildSuggestionParameters(s)),
		IOE.FromEither[error, *pb.StockAutocompleteParameters],
	)
}

func fetchSuggestionCollection(
	s listStockSuggestionsContext,
) IOE.IOEither[error, *pb.StockSuggestionCollection] {
	return F.Pipe1(
		IOE.TryCatchError(func() (*pb.StockSuggestionCollection, error) {
			return s.client.AutocompleteStock(s.gctx, s.params)
		}),
		IOE.MapLeft[*pb.StockSuggestionCollection](func(err error) error {
			return fmt.Errorf("fetch stock suggestions for query %q: %w", s.query, err)
		}),
	)
}

var buildSuggestionList = F.Curry2(
	func(
		s listStockSuggestionsContext,
		suggestions []*models.StockSuggestion,
	) *models.StockSuggestionList {
		meta := s.collection.Meta
		lmt := int(meta.Limit)
		return &models.StockSuggestionList{
			Suggestions: suggestions,
			Limit:       &lmt,
			TotalCount:  int(meta.Total),
		}
	},
)

// toStockSuggestionList projects the fetched collection onto the GraphQL
// list type. One bad suggestion fails the whole result, because the
// GraphQL enums carry no invalid value to fall back to.
func toStockSuggestionList(
	s listStockSuggestionsContext,
) E.Either[error, *models.StockSuggestionList] {
	return F.Pipe2(
		s.collection.Data,
		E.TraverseArray(convertSuggestionItem),
		E.Map[error](buildSuggestionList(s)),
	)
}
