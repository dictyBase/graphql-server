package resolver

import (
	"context"
	"fmt"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	O "github.com/IBM/fp-go/v2/option"
	T "github.com/IBM/fp-go/v2/tuple"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
)

// listStockSearchContext carries the full search request through the
// search pipeline. Only the caller-supplied fields are set at the seed;
// the params and collection are filled by the pipeline steps. The enum
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

// buildStockSearchParams fills the request params on the context from
// the caller-supplied fields.
func buildStockSearchParams(
	ctx listStockSearchContext,
) E.Either[error, listStockSearchContext] {
	return F.Pipe1(
		toProtoStockEntity(ctx.entity),
		E.Map[error](func(entity pb.StockEntity) listStockSearchContext {
			ctx.params = &pb.StockSearchParameters{
				Data: &pb.StockSearchParameters_Data{
					Type: "stock",
					Attributes: &pb.StockSearchAttributes{
						Query:  ctx.query,
						Limit:  resolveSearchLimit(ctx.limit),
						Entity: entity,
					},
				},
			}
			return ctx
		}),
	)
}

func fetchStockSearchCollection(
	ctx listStockSearchContext,
) IOE.IOEither[error, listStockSearchContext] {
	return F.Pipe2(
		IOE.TryCatchError(func() (*pb.StockSearchResultCollection, error) {
			return ctx.client.SearchStock(ctx.gctx, ctx.params)
		}),
		IOE.MapLeft[*pb.StockSearchResultCollection](func(err error) error {
			return fmt.Errorf("search stocks for query %q: %w", ctx.query, err)
		}),
		IOE.Map[error](func(coll *pb.StockSearchResultCollection) listStockSearchContext {
			ctx.collection = coll
			return ctx
		}),
	)
}

// toStockSearchList projects the fetched collection onto the GraphQL
// list type. One bad result fails the whole response, because the
// GraphQL enums carry no invalid value to fall back to.
func toStockSearchList(
	ctx listStockSearchContext,
) E.Either[error, *models.StockSearchResultList] {
	return F.Pipe2(
		ctx.collection.Data,
		E.TraverseArray(convertSearchResultItem),
		E.Map[error](func(results []*models.StockSearchResult) *models.StockSearchResultList {
			meta := ctx.collection.Meta
			lmt := int(meta.Limit)
			return &models.StockSearchResultList{
				Results:    results,
				Limit:      &lmt,
				TotalCount: int(meta.Total),
			}
		}),
	)
}
