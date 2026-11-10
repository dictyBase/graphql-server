package stock

import (
	"context"
	"testing"

	"github.com/dictyBase/graphql-server/internal/authentication"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks/clients"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func plasmidResolver(
	annoClient *clients.TaggedAnnotationServiceClient,
) *PlasmidResolver {
	return &PlasmidResolver{
		Client:           mocks.MockedStockClient(),
		UserClient:       authentication.NoopClient{},
		AnnotationClient: annoClient,
		Registry:         &mocks.MockRegistry{},
		Logger:           mocks.TestLogger(),
	}
}

func TestPlasmidInStock(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	r := plasmidResolver(mocks.MockedInStockClient())
	p, err := r.InStock(context.Background(),
		ConvertToPlasmidModel(
			"DBP0000120",
			mocks.MockPlasmidInputWithParams("kenny@bania.com"),
		))

	assert.NoError(err, "expect no error from getting plasmid inventory")
	assert.True(p, "should return true after finding inventory")
}

func TestPlasmidGenes(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	r := plasmidResolver(mocks.MockedGenModClient())
	g, err := r.Genes(
		context.Background(),
		ConvertToPlasmidModel(
			"DBP0000120",
			mocks.MockPlasmidInputWithParams("kenny@bania.com"),
		),
	)
	assert.NoError(err, "expect no error from getting associated genes")
	genes := []*models.Gene{
		{ID: "DDB_G0285425", Name: "DDB_G0285425"},
	}
	assert.ElementsMatch(g, genes, "should match associated genes")
}

func TestPlasmidDepositorAuthDisabledNonEmpty(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	r := plasmidResolver(mocks.MockedAnnotationClient())
	d, err := r.Depositor(
		context.Background(),
		ConvertToPlasmidModel(
			"DBP0391806",
			mocks.MockPlasmidInputWithParams("gadi@bcm.edu"),
		),
	)
	req.NoError(err)
	req.NotNil(d)
	attrs := d.GetData().GetAttributes()
	req.Equal("Gadi", attrs.GetFirstName())
	req.Empty(attrs.GetLastName())
	req.Empty(attrs.GetEmail())
	req.NotEqual("Dicty", attrs.GetFirstName())
}

func TestPlasmidDepositorAuthDisabledEmpty(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	r := plasmidResolver(mocks.MockedAnnotationClient())
	d, err := r.Depositor(
		context.Background(),
		ConvertToPlasmidModel(
			"DBP0000120",
			mocks.MockPlasmidInputWithParams(""),
		),
	)
	req.NoError(err)
	req.NotNil(d)
	attrs := d.GetData().GetAttributes()
	req.Equal("Dicty", attrs.GetFirstName())
	req.Equal("Stock Center", attrs.GetLastName())
}

func TestPlasmidCreatedByAndUpdatedByAuthDisabled(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	r := plasmidResolver(mocks.MockedAnnotationClient())
	obj := ConvertToPlasmidModel(
		"DBP0391806",
		mocks.MockPlasmidInputWithParams("gadi@bcm.edu"),
	)
	cb, err := r.CreatedBy(context.Background(), obj)
	req.NoError(err)
	req.Equal("Dicty", cb.GetData().GetAttributes().GetFirstName())
	req.Equal("Stock Center", cb.GetData().GetAttributes().GetLastName())

	ub, err := r.UpdatedBy(context.Background(), obj)
	req.NoError(err)
	req.Equal("Dicty", ub.GetData().GetAttributes().GetFirstName())
	req.Equal("Stock Center", ub.GetData().GetAttributes().GetLastName())
}

func TestPlasmidNotInStock(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	r := plasmidResolver(mocks.MockedNotInStockClient())

	inStock, err := r.InStock(context.Background(),
		ConvertToPlasmidModel(
			"DBP0000120",
			mocks.MockPlasmidInputWithParams("kenny@bania.com"),
		))

	assert.NoError(err, "expect no error from a missing inventory entry")
	assert.False(inStock, "should return false after a not-found inventory entry")
}

func TestPlasmidInStockError(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	r := plasmidResolver(mocks.MockedInStockErrorClient())

	inStock, err := r.InStock(context.Background(),
		ConvertToPlasmidModel(
			"DBP0000120",
			mocks.MockPlasmidInputWithParams("kenny@bania.com"),
		))

	require.Error(err, "expect an error from the inventory service failure")
	require.ErrorContains(err, "from inventory")
	require.False(inStock, "should return false on an inventory service failure")
}
