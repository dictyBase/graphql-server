package resolver

import (
	"context"
	"fmt"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	T "github.com/IBM/fp-go/v2/tuple"
	"github.com/dictyBase/aphgrpc"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/models"

	"github.com/dictyBase/graphql-server/internal/graphql/resolverutils"
)

type listPlasmidsContext struct {
	client pb.StockServiceClient
	gctx   context.Context
	cursor *int
	limit  *int
	filter *models.PlasmidListFilter
}

type withListPlasmidParams struct {
	listPlasmidsContext
	cus int64
	lmt int64
}

type withListPlasmidFilter struct {
	withListPlasmidParams
	filterQuery string
}

type withListPlasmidCollection struct {
	withListPlasmidFilter
	collection *pb.PlasmidCollection
}

var setListPlasmidParams = F.Curry2(
	func(params T.Tuple2[int64, int64], ctx listPlasmidsContext) withListPlasmidParams {
		return withListPlasmidParams{
			listPlasmidsContext: ctx,
			cus:                 params.F1,
			lmt:                 params.F2,
		}
	},
)

var setListPlasmidFilter = F.Curry2(
	func(query string, ctx withListPlasmidParams) withListPlasmidFilter {
		return withListPlasmidFilter{
			withListPlasmidParams: ctx,
			filterQuery:           query,
		}
	},
)

var setListPlasmidCollection = F.Curry2(
	func(coll *pb.PlasmidCollection, ctx withListPlasmidFilter) withListPlasmidCollection {
		return withListPlasmidCollection{
			withListPlasmidFilter: ctx,
			collection:            coll,
		}
	},
)

func toEither[ERR, A any](ioe IOE.IOEither[ERR, A]) E.Either[ERR, A] {
	return ioe()
}

func onPlasmidListError(
	err error,
) T.Tuple2[error, *models.PlasmidListWithCursor] {
	return T.MakeTuple2(err, &models.PlasmidListWithCursor{})
}

func onPlasmidListSuccess(
	data *models.PlasmidListWithCursor,
) T.Tuple2[error, *models.PlasmidListWithCursor] {
	return T.MakeTuple2[error](nil, data)
}

func computeListPlasmidParams(
	ctx listPlasmidsContext,
) T.Tuple2[int64, int64] {
	return T.MakeTuple2(
		resolverutils.GetCursorFP(ctx.cursor),
		resolverutils.GetLimitFP(ctx.limit),
	)
}

func buildListPlasmidFilterQuery(ctx withListPlasmidParams) IOE.IOEither[error, string] {
	return F.Pipe7(
		ctx.filter,
		E.FromNillable[models.PlasmidListFilter](fmt.Errorf("nil filter")),
		E.Chain(resolverutils.CheckIDField),
		E.Chain(resolverutils.CheckInStockField),
		E.Chain(resolverutils.CheckUnverifiedPlasmidType),
		E.Chain(resolverutils.CheckValidPlasmidType),
		E.Map[error](resolverutils.BuildPlasmidFieldQuery),
		IOE.FromEither[error, string],
	)
}

func fetchListPlasmidCollection(
	ctx withListPlasmidFilter,
) IOE.IOEither[error, *pb.PlasmidCollection] {
	return IOE.TryCatchError(func() (*pb.PlasmidCollection, error) {
		return ctx.client.ListPlasmids(
			ctx.gctx,
			&pb.StockParameters{
				Cursor: ctx.cus,
				Limit:  ctx.lmt,
				Filter: ctx.filterQuery,
			})
	})
}

func extractListPlasmidResult(
	ctx withListPlasmidCollection,
) *models.PlasmidListWithCursor {
	plasmids := F.Pipe1(
		ctx.collection.Data,
		A.Map(func(item *pb.PlasmidCollection_Data) *models.Plasmid {
			dbxrefs := F.Pipe1(
				item.Attributes.Dbxrefs,
				A.Map(func(s string) *string { return &s }),
			)
			return &models.Plasmid{
				ID:              item.Id,
				CreatedAt:       aphgrpc.ProtoTimeStamp(item.Attributes.CreatedAt),
				UpdatedAt:       aphgrpc.ProtoTimeStamp(item.Attributes.UpdatedAt),
				Summary:         &item.Attributes.Summary,
				EditableSummary: &item.Attributes.EditableSummary,
				Dbxrefs:         dbxrefs,
				ImageMap:        &item.Attributes.ImageMap,
				Sequence:        &item.Attributes.Sequence,
				Name:            item.Attributes.Name,
				LazyStock: models.LazyStock{
					CreatedBy:    item.Attributes.CreatedBy,
					UpdatedBy:    item.Attributes.UpdatedBy,
					Depositor:    item.Attributes.Depositor,
					Genes:        item.Attributes.Genes,
					Publications: item.Attributes.Publications,
				},
			}
		}),
	)
	lmt := int(ctx.collection.Meta.Limit)
	return &models.PlasmidListWithCursor{
		Plasmids:       plasmids,
		NextCursor:     int(ctx.collection.Meta.NextCursor),
		PreviousCursor: int(ctx.cus),
		Limit:          &lmt,
		TotalCount:     int(ctx.collection.Meta.Total),
	}
}
