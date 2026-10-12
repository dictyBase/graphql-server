package stock

import (
	"context"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	ioeutils "github.com/dictyBase/fp-go-loom/ioeitherutils"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/sirupsen/logrus"
)

// StockSearchResultResolver resolves the computed fields of one stock
// search result.
type StockSearchResultResolver struct {
	AnnotationClient annotation.TaggedAnnotationServiceClient
	Logger           *logrus.Entry
}

// InStock reports whether the matched stock is in the stock center
// inventory. The entity of the result selects the inventory annotation
// set. A missing annotation means that the stock is not in stock.
func (srr *StockSearchResultResolver) InStock(
	ctx context.Context,
	obj *models.StockSearchResult,
) (bool, error) {
	result := F.Pipe2(
		inStockForEntity(ctx, srr.AnnotationClient, obj.ID, obj.Entity),
		ioeutils.ToEither[error, bool],
		E.Fold(onInventoryError, onInventorySuccess),
	)

	if result.F1 != nil {
		errorutils.AddGQLError(ctx, result.F1)
		srr.Logger.Error(result.F1)
		return false, result.F1
	}
	return result.F2, nil
}
