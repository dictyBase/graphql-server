package resolver

import (
	"context"
	"fmt"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	T "github.com/IBM/fp-go/v2/tuple"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
)

// listStockSearchContext carries the full search request through the
// search pipeline. Only the caller-supplied fields are set at the seed;
// the params and collection are attached by the lens setters. The enum
// conversions are shared with the autocomplete pipeline in
// stock_suggestions_fp.go.
type listStockSearchContext struct {
	client     pb.StockServiceClient
	gctx       context.Context
	query      string
	limit      *int
	entity     *models.StockEntityType
	params     *pb.StockSearchParameters
	collection *pb.StockSearchResultCollection
}

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
)

// zeroSearchLimit lets the stock service apply its own default of 50
// when the GraphQL query omits the limit.
var zeroSearchLimit = int64(0)

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

// convertSearchResultItem projects one protobuf search result onto the
// GraphQL type. Both enum fields must resolve, otherwise the whole result
// fails.
func convertSearchResultItem(
	item *pb.StockSearchResult,
) E.Either[error, *models.StockSearchResult] {
	return F.Pipe1(
		toGQLEntity(item.Entity),
		E.Chain(func(entity models.StockEntityType) E.Either[error, *models.StockSearchResult] {
			return F.Pipe1(
				toGQLSearchField(item.Field),
				E.Map[error](func(field models.StockSearchFieldEnum) *models.StockSearchResult {
					return &models.StockSearchResult{
						ID:          item.Id,
						Entity:      entity,
						Field:       field,
						DisplayText: item.DisplayText,
						Score:       item.Score,
						Summary:     item.Summary,
						StrainLabel: item.StrainLabel,
					}
				}),
			)
		}),
	)
}

// ── pipeline steps ──────────────────────────────────────────────────────────

var buildSearchParameters = F.Curry2(
	func(
		s listStockSearchContext,
		entity pb.StockEntity,
	) *pb.StockSearchParameters {
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
	},
)

// searchParamsFromState derives the request params from the caller
// supplied fields. It fails when the entity filter is unknown.
func searchParamsFromState(
	s listStockSearchContext,
) E.Either[error, *pb.StockSearchParameters] {
	return F.Pipe4(
		s.entity,
		O.FromNillable2[models.StockEntityType],
		O.GetOrElse(F.Constant(models.StockEntityTypeAll)),
		toKnownProtoEntity,
		E.Map[error](buildSearchParameters(s)),
	)
}

func fetchSearchCollection(
	s listStockSearchContext,
) IOE.IOEither[error, *pb.StockSearchResultCollection] {
	return F.Pipe1(
		IOE.TryCatchError(func() (*pb.StockSearchResultCollection, error) {
			return s.client.SearchStock(s.gctx, s.params)
		}),
		IOE.MapLeft[*pb.StockSearchResultCollection](func(err error) error {
			return fmt.Errorf("search stocks for query %q: %w", s.query, err)
		}),
	)
}

var buildSearchResultList = F.Curry2(
	func(
		s listStockSearchContext,
		results []*models.StockSearchResult,
	) *models.StockSearchResultList {
		meta := s.collection.Meta
		lmt := int(meta.Limit)
		return &models.StockSearchResultList{
			Results:    results,
			Limit:      &lmt,
			TotalCount: int(meta.Total),
		}
	},
)

// toStockSearchList projects the fetched collection onto the GraphQL
// list type. One bad result fails the whole response, because the
// GraphQL enums carry no invalid value to fall back to.
func toStockSearchList(
	s listStockSearchContext,
) E.Either[error, *models.StockSearchResultList] {
	return F.Pipe2(
		s.collection.Data,
		E.TraverseArray(convertSearchResultItem),
		E.Map[error](buildSearchResultList(s)),
	)
}
