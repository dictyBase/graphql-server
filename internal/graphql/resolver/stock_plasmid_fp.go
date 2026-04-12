package resolver

import (
	"context"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	T "github.com/IBM/fp-go/v2/tuple"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/graphql/resolver/stock"
)

type plasmidCollectionIOE = IOE.IOEither[error, *pb.PlasmidCollection]

func buildPlasmidStockParameters(
	cursor, limit int64,
	filter string,
) *pb.StockParameters {
	return &pb.StockParameters{
		Cursor: cursor,
		Limit:  limit,
		Filter: filter,
	}
}

func callListPlasmids(
	client pb.StockServiceClient,
	ctx context.Context,
	params *pb.StockParameters,
) plasmidCollectionIOE {
	return IOE.TryCatchError(func() (*pb.PlasmidCollection, error) {
		return client.ListPlasmids(ctx, params)
	})
}

func convertPlasmidCollection(
	list *pb.PlasmidCollection,
) []*models.Plasmid {
	return F.Pipe1(
		list.Data,
		A.Map(func(item *pb.PlasmidCollection_Data) *models.Plasmid {
			return stock.ConvertToPlasmidModel(item.Id, item.Attributes)
		}),
	)
}

func buildPlasmidListResult(
	list *pb.PlasmidCollection,
	cursor int64,
) *models.PlasmidListWithCursor {
	plasmids := convertPlasmidCollection(list)
	lmt := int(list.Meta.Limit)
	return &models.PlasmidListWithCursor{
		Plasmids:       plasmids,
		NextCursor:     int(list.Meta.NextCursor),
		PreviousCursor: int(cursor),
		Limit:          &lmt,
		TotalCount:     int(list.Meta.Total),
	}
}

func foldPlasmidListResult(
	ioe IOE.IOEither[error, *models.PlasmidListWithCursor],
) T.Tuple2[error, *models.PlasmidListWithCursor] {
	return F.Pipe1(
		ioe(),
		E.Fold(
			func(err error) T.Tuple2[error, *models.PlasmidListWithCursor] {
				return T.MakeTuple2(err, &models.PlasmidListWithCursor{})
			},
			func(data *models.PlasmidListWithCursor) T.Tuple2[error, *models.PlasmidListWithCursor] {
				return T.MakeTuple2[error](nil, data)
			},
		),
	)
}
