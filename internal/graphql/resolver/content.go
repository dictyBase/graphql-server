package resolver

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/content"
	"github.com/dictyBase/graphql-server/internal/authentication"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
)

var (
	contentCreatorScope = "write:content"
	contentEditorScope  = "edit:content"
	contentDeleteScope  = "delete:content"
	noncharReg          = regexp.MustCompile("[^a-z0-9]+")
)

func Slugify(name string) string {
	return strings.Trim(
		noncharReg.ReplaceAllString(strings.ToLower(name), "-"),
		"-",
	)
}

func (mrs *MutationResolver) CreateContent(
	ctx context.Context,
	input *models.CreateContentInput,
) (*pb.Content, error) {
	if err := authentication.ValidateContent(ctx, "scope", contentCreatorScope); err != nil {
		errorutils.AddGQLError(ctx, err)
		mrs.Logger.Error(err)
		return nil, err
	}
	cnt, err := mrs.GetContentClient(registry.CONTENT).
		StoreContent(ctx, &pb.StoreContentRequest{
			Data: &pb.StoreContentRequest_Data{
				Type: "contents",
				Attributes: &pb.NewContentAttributes{
					Name:      input.Name,
					CreatedBy: input.CreatedBy,
					Content:   input.Content,
					Namespace: input.Namespace,
					Slug: Slugify(
						fmt.Sprintf("%s %s", input.Namespace, input.Name),
					),
				},
			},
		})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		mrs.Logger.Error(err)
		return nil, err
	}
	mrs.Logger.Debugf(
		"successfully created new content with ID %d",
		cnt.Data.Id,
	)
	return cnt, nil
}

func (mrs *MutationResolver) UpdateContent(
	ctx context.Context,
	input *models.UpdateContentInput,
) (*pb.Content, error) {
	if err := authentication.ValidateContent(ctx, "scope", contentEditorScope); err != nil {
		errorutils.AddGQLError(ctx, err)
		mrs.Logger.Error(err)
		return nil, err
	}
	cid, err := strconv.ParseInt(input.ID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf(
			"error in parsing string %s to int %s",
			input.ID,
			err,
		)
	}
	cnt, err := mrs.GetContentClient(registry.CONTENT).
		UpdateContent(ctx, &pb.UpdateContentRequest{
			Id: cid,
			Data: &pb.UpdateContentRequest_Data{
				Type: "contents",
				Id:   cid,
				Attributes: &pb.ExistingContentAttributes{
					UpdatedBy: input.UpdatedBy,
					Content:   input.Content,
				},
			},
		})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		mrs.Logger.Error(err)
		return nil, err
	}
	ucnt, err := mrs.GetContentClient(registry.CONTENT).
		GetContent(ctx, &pb.ContentIdRequest{Id: cnt.Data.Id})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		mrs.Logger.Error(err)
		return nil, err
	}
	mrs.Logger.Debugf("successfully updated content with ID %d", ucnt.Data.Id)
	return ucnt, nil
}

func (mrs *MutationResolver) DeleteContent(
	ctx context.Context,
	id string,
) (*models.DeleteContent, error) {
	if err := authentication.ValidateContent(ctx, "scope", contentDeleteScope); err != nil {
		errorutils.AddGQLError(ctx, err)
		mrs.Logger.Error(err)
		return nil, err
	}
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error in parsing string %s to int %s", id, err)
	}
	if _, err := mrs.GetContentClient(registry.CONTENT).DeleteContent(ctx, &pb.ContentIdRequest{Id: cid}); err != nil {
		return &models.DeleteContent{
			Success: false,
		}, err
	}
	mrs.Logger.Debugf("successfully deleted content with ID %s", id)
	return &models.DeleteContent{
		Success: true,
	}, nil
}

func (qrs *QueryResolver) Content(
	ctx context.Context,
	id string,
) (*pb.Content, error) {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error in parsing string %s to int %s", id, err)
	}
	content, err := qrs.GetContentClient(registry.CONTENT).
		GetContent(ctx, &pb.ContentIdRequest{Id: cid})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		qrs.Logger.Error(err)
		return nil, err
	}
	qrs.Logger.Debugf("successfully found content with ID %s", id)
	return content, nil
}

func (qrs *QueryResolver) ContentBySlug(
	ctx context.Context,
	slug string,
) (*pb.Content, error) {
	content, err := qrs.GetContentClient(registry.CONTENT).
		GetContentBySlug(ctx, &pb.ContentRequest{Slug: slug})
	if err != nil {
		errorutils.AddGQLError(ctx, err)
		qrs.Logger.Error(err)
		return nil, err
	}
	qrs.Logger.Debugf("successfully found content with slug %s", slug)
	return content, nil
}
