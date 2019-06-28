package resolver

import (
	"context"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/publication"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/utils"
)

// Publication is the resolver for getting an individual publication by ID.
func (q *QueryResolver) Publication(ctx context.Context, id string) (*pb.Publication, error) {
	p := &pb.Publication{}
	pub, err := utils.FetchPublication(ctx, q.Registry, id)
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		q.Logger.Error(err)
		return p, err
	}
	return pub, nil
}
