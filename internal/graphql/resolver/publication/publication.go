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

func (r *PublicationResolver) ID(
	ctx context.Context,
	obj *models.Publication,
) (string, error) {
	return obj.ID, nil
}

func (r *PublicationResolver) Doi(
	ctx context.Context,
	obj *models.Publication,
) (*string, error) {
	return obj.Doi, nil
}

func (r *PublicationResolver) Title(
	ctx context.Context,
	obj *models.Publication,
) (string, error) {
	return obj.Title, nil
}

func (r *PublicationResolver) Abstract(
	ctx context.Context,
	obj *models.Publication,
) (string, error) {
	return obj.Abstract, nil
}

func (r *PublicationResolver) Journal(
	ctx context.Context,
	obj *models.Publication,
) (string, error) {
	return obj.Journal, nil
}

func (r *PublicationResolver) PubDate(
	ctx context.Context,
	obj *models.Publication,
) (*time.Time, error) {
	return obj.PubDate, nil
}

func (r *PublicationResolver) Volume(
	ctx context.Context,
	obj *models.Publication,
) (*string, error) {
	return obj.Volume, nil
}

func (r *PublicationResolver) Pages(
	ctx context.Context,
	obj *models.Publication,
) (*string, error) {
	return obj.Pages, nil
}

func (r *PublicationResolver) Issn(
	ctx context.Context,
	obj *models.Publication,
) (*string, error) {
	return obj.Issn, nil
}

func (r *PublicationResolver) PubType(
	ctx context.Context,
	obj *models.Publication,
) (string, error) {
	return obj.PubType, nil
}

func (r *PublicationResolver) Source(
	ctx context.Context,
	obj *models.Publication,
) (string, error) {
	return obj.Source, nil
}

func (r *PublicationResolver) Issue(
	ctx context.Context,
	obj *models.Publication,
) (*string, error) {
	return obj.Issue, nil
}

func (r *PublicationResolver) Status(
	ctx context.Context,
	obj *models.Publication,
) (*string, error) {
	return obj.Status, nil
}

func (r *PublicationResolver) Authors(
	ctx context.Context,
	obj *models.Publication,
) ([]*publication.Author, error) {
	return obj.Authors, nil
}
