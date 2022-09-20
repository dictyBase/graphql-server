package publication

import (
	"context"
	"time"

	"github.com/dictyBase/graphql-server/internal/graphql/models"

	"github.com/dictyBase/go-genproto/dictybaseapis/publication"
	"github.com/sirupsen/logrus"
)

type PublicationResolver struct {
	Logger *logrus.Entry
}

func (r *PublicationResolver) Authors(
	ctx context.Context,
	obj *models.Publication,
) ([]*publication.Author, error) {
	return obj.Authors, nil
}
