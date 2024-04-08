package resolver

import (
	"context"
	"fmt"
	"time"

	"github.com/dictyBase/graphql-server/internal/authentication"

	"github.com/minio/minio-go/v7"

	"github.com/99designs/gqlgen/graphql"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/google/uuid"
)

func (mrs *MutationResolver) UploadFile(
	ctx context.Context,
	file graphql.Upload,
) (*models.ImageFile, error) {
	if err := authentication.ValidateContent(ctx, "scope", authentication.ContentCreatorScope); err != nil {
		errorutils.AddGQLError(ctx, err)
		mrs.Logger.Error(err)
		return nil, err
	}
	rndId, err := uuid.NewRandom()
	if err != nil {
		iderr := fmt.Errorf("error in generating random uuid %s", err)
		errorutils.AddGQLError(ctx, iderr)
		return nil, err
	}
	bucket := mrs.Registry.GetRecord(registry.S3Bucket)
	fileInBucket := fmt.Sprintf(
		"%s/%s/%s",
		mrs.Registry.GetRecord(registry.S3BucketPath),
		time.Now().Format(time.DateOnly),
		rndId.String(),
	)
	uploadInfo, err := mrs.Registry.GetS3Client(registry.S3CLIENT).PutObject(
		ctx, bucket, fileInBucket, file.File, file.Size, minio.PutObjectOptions{},
	)
	if err != nil {
		uerr := fmt.Errorf("error in uploading file %s", file.Filename)
		errorutils.AddGQLError(ctx, uerr)
		return nil, uerr
	}
	mrs.Logger.Debug(uploadInfo)
	return &models.ImageFile{
		URL: fmt.Sprintf(
			"%s/%s/%s",
			mrs.Registry.GetAPIEndpoint(registry.S3STORAGE),
			bucket,
			fileInBucket,
		)}, nil
}
