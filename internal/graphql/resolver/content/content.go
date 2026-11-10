package content

import (
	"context"
	"fmt"
	"strconv"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/dictyBase/aphgrpc"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/content"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/authentication"
	"github.com/dictyBase/graphql-server/internal/graphql/resolverutils"
	"github.com/sirupsen/logrus"
)

type ContentResolver struct {
	Client     pb.ContentServiceClient
	UserClient authentication.LogtoClient
	Logger     *logrus.Entry
	StorageURL string
	Bucket     string
	BucketPath string
}

func (rcs *ContentResolver) ID(
	ctx context.Context,
	obj *pb.Content,
) (string, error) {
	return strconv.FormatInt(obj.Data.Id, 10), nil
}

func (rcs *ContentResolver) Name(
	ctx context.Context,
	obj *pb.Content,
) (string, error) {
	return obj.Data.Attributes.Name, nil
}

func (rcs *ContentResolver) Slug(
	ctx context.Context,
	obj *pb.Content,
) (string, error) {
	return obj.Data.Attributes.Slug, nil
}

func (rcs *ContentResolver) userByEmail(
	ctx context.Context,
	email string,
) (*user.User, error) {
	res := F.Pipe5(
		resolverutils.UserLookupState{},
		resolverutils.LookupCtxLens.Set(ctx),
		resolverutils.LookupQueryLens.Set(email),
		resolverutils.LookupClientLens.Set(rcs.UserClient),
		resolverutils.LookupLoggerLens.Set(rcs.Logger),
		resolverutils.UserByEmail,
	)
	return res.F2, res.F1
}

func (rcs *ContentResolver) CreatedBy(
	ctx context.Context,
	obj *pb.Content,
) (*user.User, error) {
	user, err := rcs.userByEmail(ctx, obj.Data.Attributes.CreatedBy)
	if err != nil {
		return user, fmt.Errorf("error in getting created by user %s", err)
	}
	return user, nil
}

func (rcs *ContentResolver) UpdatedBy(
	ctx context.Context,
	obj *pb.Content,
) (*user.User, error) {
	user, err := rcs.userByEmail(ctx, obj.Data.Attributes.UpdatedBy)
	if err != nil {
		return user, fmt.Errorf("error in getting updated by user %s", err)
	}
	return user, nil
}

func (rcs *ContentResolver) CreatedAt(
	ctx context.Context,
	obj *pb.Content,
) (*time.Time, error) {
	time := aphgrpc.ProtoTimeStamp(obj.Data.Attributes.CreatedAt)
	return &time, nil
}

func (rcs *ContentResolver) UpdatedAt(
	ctx context.Context,
	obj *pb.Content,
) (*time.Time, error) {
	time := aphgrpc.ProtoTimeStamp(obj.Data.Attributes.UpdatedAt)
	return &time, nil
}

func (rcs *ContentResolver) Content(
	_ context.Context,
	obj *pb.Content,
) (string, error) {
	rawContent := obj.GetData().GetAttributes().GetContent()
	return F.Pipe10(
		contentStorageState{
			storageURL: rcs.StorageURL,
			bucket:     rcs.Bucket,
			bucketPath: rcs.BucketPath,
			rawContent: rawContent,
		},
		storageURLLens.Modify(trimTrailingSlash),
		bucketLens.Modify(trimSlashes),
		bucketPathLens.Modify(trimSlashes),
		O.FromPredicate(isStorageConfigComplete),
		O.Let(segmentPrefixLens.Set, buildSegmentPrefix),
		O.Let(targetPrefixLens.Set, buildTargetPrefix),
		O.Bind(patternLens.Set, compileStorageURLPattern),
		O.Let(replacerLens.Set, buildPrefixReplacer),
		O.Map(rewriteStorageURLsInState),
		O.GetOrElse(F.Constant(rawContent)),
	), nil
}

func (rcs *ContentResolver) Namespace(
	ctx context.Context,
	obj *pb.Content,
) (string, error) {
	return obj.Data.Attributes.Namespace, nil
}
