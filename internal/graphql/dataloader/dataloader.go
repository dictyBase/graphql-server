package dataloader

//go:generate go run github.com/vektah/dataloaden StrainLoader string *github.com/dictyBase/graphql-server/internal/graphql/models.Strain

import (
	"context"
	"time"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/graphql/resolver/stock"
	"github.com/dictyBase/graphql-server/internal/registry"
)

type contextKey string

const key = contextKey("dataloaders")

type Loaders struct {
	StrainByID *StrainLoader
}

func newLoaders(ctx context.Context, nr registry.Registry) *Loaders {
	return &Loaders{
		StrainByID: newStrainByID(ctx, nr),
	}
}

// Retriever retrieves dataloaders from the request context.
type Retriever interface {
	Retrieve(context.Context) *Loaders
}

type retriever struct {
	key contextKey
}

func (r *retriever) Retrieve(ctx context.Context) *Loaders {
	return ctx.Value(r.key).(*Loaders)
}

// NewRetriever instantiates a new implementation of Retriever.
func NewRetriever() Retriever {
	return &retriever{key: key}
}

func newStrainByID(_ context.Context, nr registry.Registry) *StrainLoader {
	return NewStrainLoader(StrainLoaderConfig{
		MaxBatch: 100,
		Wait:     100 * time.Millisecond,
		Fetch: func(ids []string) ([]*models.Strain, []error) {
			strains := make([]*models.Strain, 0)
			sl, err := nr.GetStockClient(registry.STOCK).ListStrainsByIds(
				context.Background(),
				&pb.StockIdList{Id: ids},
			)
			if err != nil {
				return strains, []error{err}
			}
			for _, sd := range sl.Data {
				strains = append(
					strains,
					stock.ConvertToStrainModel(sd.Id, sd.Attributes))
			}
			return strains, nil
		},
	})
}
