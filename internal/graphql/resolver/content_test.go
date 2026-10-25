package resolver

import (
	"context"
	"testing"

	"github.com/dictyBase/go-genproto/dictybaseapis/content"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks/clients"
	contentres "github.com/dictyBase/graphql-server/internal/graphql/resolver/content"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type contentClientRegistry struct {
	*mocks.MockRegistry
	contentClient content.ContentServiceClient
}

func (r *contentClientRegistry) GetContentClient(
	key string,
) content.ContentServiceClient {
	return r.contentClient
}

func TestListContentByNamespace(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	assert := assert.New(t)

	mockedContent := new(clients.ContentServiceClient)
	mockedContent.On(
		"ListContents",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.MatchedBy(func(params *content.ListParameters) bool {
			return params.Limit == 10 &&
				params.Filter == "namespace===news"
		}),
	).Return(mocks.MockContentCollection(), nil)

	reg := &contentClientRegistry{
		MockRegistry:  &mocks.MockRegistry{},
		contentClient: mockedContent,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}
	limit := 10

	result, err := resolver.ListContentByNamespace(
		context.Background(), "news", &limit,
	)
	require.NoError(err)
	assert.Len(result, 3)
	assert.Equal(int64(100), result[0].Data.Id)
	assert.Equal("test-page", result[0].Data.Attributes.Name)
	assert.Equal("news", result[0].Data.Attributes.Namespace)
	mockedContent.AssertExpectations(t)
}

func TestListContentByNamespaceDefaultLimit(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedContent := new(clients.ContentServiceClient)
	mockedContent.On(
		"ListContents",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.MatchedBy(func(params *content.ListParameters) bool {
			return params.Limit == 10 &&
				params.Filter == "namespace===news"
		}),
	).Return(mocks.MockContentCollection(), nil)

	reg := &contentClientRegistry{
		MockRegistry:  &mocks.MockRegistry{},
		contentClient: mockedContent,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	result, err := resolver.ListContentByNamespace(
		context.Background(), "news", nil,
	)
	require.NoError(err)
	require.Len(result, 3)
	mockedContent.AssertExpectations(t)
}

func TestListContentByNamespaceEmpty(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	assert := assert.New(t)

	mockedContent := new(clients.ContentServiceClient)
	mockedContent.On(
		"ListContents",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.MatchedBy(func(params *content.ListParameters) bool {
			return params.Limit == 5 &&
				params.Filter == "namespace===docs"
		}),
	).Return(&content.ContentCollection{}, nil)

	reg := &contentClientRegistry{
		MockRegistry:  &mocks.MockRegistry{},
		contentClient: mockedContent,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}
	limit := 5

	result, err := resolver.ListContentByNamespace(
		context.Background(), "docs", &limit,
	)
	require.NoError(err)
	assert.Empty(result)
	mockedContent.AssertExpectations(t)
}

func TestResolverContentWiresStorageConfig(t *testing.T) {
	t.Parallel()
	reg := registry.NewRegistry()
	reg.AddAPIEndpoint(registry.S3STORAGE, "https://slimestorage.dictycr.org")
	reg.AddRecord(registry.S3Bucket, "editor")
	reg.AddRecord(registry.S3BucketPath, "assets")
	reg.AddAPIConnection(registry.CONTENT, nil)
	reg.AddAuthClient(registry.AUTH, nil)

	r := NewResolver(reg, nil, mocks.TestLogger())
	cr, ok := r.Content().(*contentres.ContentResolver)
	require.True(t, ok)
	require.Equal(t, "https://slimestorage.dictycr.org", cr.StorageURL)
	require.Equal(t, "editor", cr.Bucket)
	require.Equal(t, "assets", cr.BucketPath)
}

func TestResolverContentEmptyRegistryDoesNotPanic(t *testing.T) {
	t.Parallel()
	reg := &mocks.MockRegistry{}
	r := NewResolver(reg, nil, mocks.TestLogger())
	require.NotPanics(t, func() {
		_ = r.Content()
	})
}
