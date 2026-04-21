package resolver

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/graphql"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks/clients"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/graphql/resolverutils"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// stockClientRegistry returns a fixed stock client for resolver tests
// that need to verify exact parameters passed to ListPlasmids.
type stockClientRegistry struct {
	*mocks.MockRegistry
	stockClient pb.StockServiceClient
}

func (r *stockClientRegistry) GetStockClient(key string) pb.StockServiceClient {
	return r.stockClient
}

func TestPlasmid(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	q := &QueryResolver{
		Registry: &mocks.MockRegistry{},
		Logger:   mocks.TestLogger(),
	}
	plasmidID := "DBP123456"
	p, err := q.Plasmid(context.Background(), plasmidID)
	assert.NoError(err, "expect no error from getting plasmid by ID")
	assert.Equal(p.ID, plasmidID, "should match plasmid ID")
	assert.Equal(
		p.CreatedBy,
		mocks.MockPlasmidAttributes.CreatedBy,
		"should match created_by",
	)
	assert.Equal(
		p.UpdatedBy,
		mocks.MockPlasmidAttributes.UpdatedBy,
		"should match updated_by",
	)
	assert.Equal(
		p.Summary,
		&mocks.MockPlasmidAttributes.Summary,
		"should match summary",
	)
	assert.Equal(
		p.EditableSummary,
		&mocks.MockPlasmidAttributes.EditableSummary,
		"should match editable summary",
	)
	assert.Equal(
		p.Depositor,
		mocks.MockPlasmidAttributes.Depositor,
		"should match depositor (he's gold)",
	)
	assert.ElementsMatch(
		p.Genes,
		mocks.MockPlasmidAttributes.Genes,
		"should match genes list",
	)
	assert.ElementsMatch(
		p.Dbxrefs,
		sliceConverter(mocks.MockPlasmidAttributes.Dbxrefs),
		"should match dbxrefs",
	)
	assert.ElementsMatch(
		p.Publications,
		mocks.MockPlasmidAttributes.Publications,
		"should match publications",
	)
	assert.Equal(
		p.ImageMap,
		&mocks.MockPlasmidAttributes.ImageMap,
		"should match image map",
	)
	assert.Equal(
		p.Sequence,
		&mocks.MockPlasmidAttributes.Sequence,
		"should match sequence",
	)
	assert.Equal(p.Name, mocks.MockPlasmidAttributes.Name, "should match name")
}

func TestStrain(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	q := &QueryResolver{
		Registry: &mocks.MockRegistry{},
		Logger:   mocks.TestLogger(),
	}
	strainID := "DBS123456"
	p, err := q.Strain(context.Background(), strainID)
	assert.NoError(err, "expect no error from getting strain by ID")
	assert.Equal(p.ID, strainID, "should match strain ID")
	assert.Equal(
		p.CreatedBy,
		mocks.MockStrainAttributes.CreatedBy,
		"should match created_by",
	)
	assert.Equal(
		p.UpdatedBy,
		mocks.MockStrainAttributes.UpdatedBy,
		"should match updated_by",
	)
	assert.Equal(
		p.Summary,
		&mocks.MockStrainAttributes.Summary,
		"should match summary",
	)
	assert.Equal(
		p.EditableSummary,
		&mocks.MockStrainAttributes.EditableSummary,
		"should match editable summary",
	)
	assert.Equal(
		p.Depositor,
		mocks.MockStrainAttributes.Depositor,
		"should match depositor (he's gold)",
	)
	assert.ElementsMatch(
		p.Genes,
		mocks.MockStrainAttributes.Genes,
		"should match genes list",
	)
	assert.ElementsMatch(
		p.Dbxrefs,
		mocks.MockStrainAttributes.Dbxrefs,
		"should match dbxrefs",
	)
	assert.ElementsMatch(
		p.Publications,
		mocks.MockStrainAttributes.Publications,
		"should match publications",
	)
	assert.Equal(
		p.Label,
		mocks.MockStrainAttributes.Label,
		"should match label",
	)
	assert.Equal(
		p.Species,
		mocks.MockStrainAttributes.Species,
		"should match species",
	)
	assert.Equal(
		p.Plasmid,
		&mocks.MockStrainAttributes.Plasmid,
		"should match plasmid",
	)
	assert.ElementsMatch(
		p.Names,
		mocks.MockStrainAttributes.Names,
		"should match names",
	)
}

func TestListPlasmids(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	q := &QueryResolver{
		Registry: &mocks.MockRegistry{},
		Logger:   mocks.TestLogger(),
	}
	cursor := 0
	limit := 10
	filter := &models.PlasmidListFilter{
		PlasmidType: models.PlasmidTypeAll,
	}
	p, err := q.ListPlasmids(context.Background(), &cursor, &limit, filter)
	assert.NoError(err, "expect no error from getting list of plasmids")
	assert.Equal(p.Limit, &limit, "should match limit")
	assert.Equal(p.PreviousCursor, cursor, "should match previous cursor")
	assert.Equal(
		p.NextCursor,
		0,
		"should not have value for next cursor since less results than limit",
	)
	assert.Equal(p.TotalCount, 3, "should match total count (length) of items")
	assert.Len(p.Plasmids, 3, "should have three plasmids")
}

func TestListPlasmidsWithAnnotation(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	q := &QueryResolver{
		Registry: &mocks.MockRegistry{},
		Logger:   mocks.TestLogger(),
	}
	cursor := 0
	limit := 10
	s, err := q.ListPlasmidsWithAnnotation(
		context.Background(),
		&cursor,
		&limit,
		"plasmid_inventory",
		"plasmid inventory",
	)
	assert.NoError(err, "expect no error from getting list of plasmids")
	assert.Equal(s.Limit, &limit, "should match limit")
	assert.Equal(s.PreviousCursor, cursor, "should match previous cursor")
	assert.Equal(
		s.NextCursor,
		0,
		"should not have value for next cursor since less results than limit",
	)
	assert.Equal(s.TotalCount, 4, "should match total count (length) of items")
	assert.Len(s.Plasmids, 4, "should have four plasmids")
}

func TestListStrains(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	q := &QueryResolver{
		Registry: &mocks.MockRegistry{},
		Logger:   mocks.TestLogger(),
	}
	cursor := 0
	limit := 10
	// filter := "type===strain"
	s, err := q.ListStrains(
		context.Background(),
		&cursor,
		&limit,
		&models.StrainListFilter{StrainType: models.StrainTypeGwdi},
	)
	assert.NoError(err, "expect no error from getting list of strains")
	assert.Equal(s.Limit, &limit, "should match limit")
	assert.Equal(s.PreviousCursor, cursor, "should match previous cursor")
	assert.Equal(
		s.NextCursor,
		0,
		"should not have value for next cursor since less results than limit",
	)
	assert.Equal(s.TotalCount, 3, "should match total count (length) of items")
	assert.Len(s.Strains, 3, "should have three strains")
}

/* func TestListStrainsWithAnnotation(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	q := &QueryResolver{
		Registry: &mocks.MockRegistry{},
		Logger:   mocks.TestLogger(),
	}
	cursor := 0
	limit := 10
	s, err := q.ListStrainsWithAnnotation(
		context.Background(),
		&cursor,
		&limit,
		"phenotype",
		"delayed culmination",
	)
	assert.NoError(err, "expect no error from getting list of strains")
	assert.Equal(s.Limit, &limit, "should match limit")
	assert.Equal(s.PreviousCursor, cursor, "should match previous cursor")
	assert.Equal(
		s.NextCursor,
		0,
		"should not have value for next cursor since less results than limit",
	)
	assert.Equal(s.TotalCount, 4, "should match total count (length) of items")
	assert.Len(s.Strains, 4, "should have four strains")
} */

func TestCreatePlasmid(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	m := &MutationResolver{
		Registry: &mocks.MockRegistry{},
		Logger:   mocks.TestLogger(),
	}
	summary := "test summary"
	esummary := "editable test summary"
	depositor := "kenny@bania.com"
	input := &models.CreatePlasmidInput{
		CreatedBy:       "art@vandelay.com",
		UpdatedBy:       "art@vandelay.com",
		Summary:         &summary,
		EditableSummary: &esummary,
		Depositor:       &depositor,
		InStock:         true,
	}
	p, err := m.CreatePlasmid(context.Background(), input)
	assert.NoError(err, "expect no error from creating new plasmid")
	assert.Equal(p.ID, "DBP123456", "should match plasmid ID")
	assert.Equal(p.CreatedBy, input.CreatedBy, "should match created_by")
	assert.Equal(p.UpdatedBy, input.UpdatedBy, "should match updated_by")
	assert.Equal(*p.Summary, *input.Summary, "should match summary")
	assert.Equal(
		*p.EditableSummary,
		*input.EditableSummary,
		"should match editable summary",
	)
	assert.Equal(
		p.Depositor,
		*input.Depositor,
		"should match depositor (he's gold)",
	)
}

func TestCreateStrain(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	m := &MutationResolver{
		Registry: &mocks.MockRegistry{},
		Logger:   mocks.TestLogger(),
	}
	summary := "test summary"
	esummary := "editable test summary"
	depositor := "kenny@bania.com"
	input := &models.CreateStrainInput{
		CreatedBy:       "art@vandelay.com",
		UpdatedBy:       "art@vandelay.com",
		Summary:         &summary,
		EditableSummary: &esummary,
		Depositor:       &depositor,
		SystematicName:  "test1",
		Label:           "test99",
		Species:         "human",
		InStock:         true,
	}
	p, err := m.CreateStrain(context.Background(), input)
	assert.NoError(err, "expect no error from creating new strain")
	assert.Equal(p.ID, "DBS123456", "should match strain ID")
	assert.Equal(p.CreatedBy, input.CreatedBy, "should match created_by")
	assert.Equal(p.UpdatedBy, input.UpdatedBy, "should match updated_by")
	assert.Equal(*p.Summary, *input.Summary, "should match summary")
	assert.Equal(
		*p.EditableSummary,
		*input.EditableSummary,
		"should match editable summary",
	)
	assert.Equal(
		p.Depositor,
		*input.Depositor,
		"should match depositor (he's gold)",
	)
	assert.Equal(p.Label, input.Label, "should match label")
	assert.Equal(p.Species, input.Species, "should match species")
}

func TestUpdatePlasmid(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	m := &MutationResolver{
		Registry: &mocks.MockRegistry{},
		Logger:   mocks.TestLogger(),
	}
	summary := "updated summary"
	esummary := "editable updated summary"
	depositor := "puddy@nyrangers.com"
	input := &models.UpdatePlasmidInput{
		UpdatedBy:       "h.e.@pennypacker.com",
		Summary:         &summary,
		EditableSummary: &esummary,
		Depositor:       &depositor,
	}
	p, err := m.UpdatePlasmid(context.Background(), "DBP123456", input)
	assert.NoError(err, "expect no error from creating new plasmid")
	assert.Equal(
		p.UpdatedBy,
		input.UpdatedBy,
		"should match updated updated_by",
	)
	assert.Equal(p.Summary, input.Summary, "should match updated summary")
	assert.Equal(
		p.EditableSummary,
		input.EditableSummary,
		"should match updated editable summary",
	)
	assert.Equal(
		p.Depositor,
		*input.Depositor,
		"should match updated depositor (he's gold)",
	)
	assert.ElementsMatch(
		p.Genes,
		mocks.MockUpdatePlasmidAttributes.Genes,
		"should match existing genes list",
	)
	assert.ElementsMatch(
		p.Dbxrefs,
		sliceConverter(mocks.MockUpdatePlasmidAttributes.Dbxrefs),
		"should match existing dbxrefs",
	)
	assert.ElementsMatch(
		p.Publications,
		mocks.MockUpdatePlasmidAttributes.Publications,
		"should match existing publications",
	)
	assert.Equal(
		p.ImageMap,
		&mocks.MockUpdatePlasmidAttributes.ImageMap,
		"should match existing image map",
	)
	assert.Equal(
		p.Sequence,
		&mocks.MockUpdatePlasmidAttributes.Sequence,
		"should match existing sequence",
	)
	assert.Equal(
		p.Name,
		mocks.MockUpdatePlasmidAttributes.Name,
		"should match name",
	)
}

func TestUpdateStrain(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	m := &MutationResolver{
		Registry: &mocks.MockRegistry{},
		Logger:   mocks.TestLogger(),
	}
	summary := "updated summary"
	esummary := "editable updated summary"
	depositor := "puddy@nyrangers.com"
	input := &models.UpdateStrainInput{
		UpdatedBy:       "h.e.@pennypacker.com",
		Summary:         &summary,
		EditableSummary: &esummary,
		Depositor:       &depositor,
	}
	p, err := m.UpdateStrain(context.Background(), "DBS123456", input)
	assert.NoError(err, "expect no error from creating new strain")
	assert.Equal(
		p.UpdatedBy,
		input.UpdatedBy,
		"should match updated updated_by",
	)
	assert.Equal(p.Summary, input.Summary, "should match updated summary")
	assert.Equal(
		p.EditableSummary,
		input.EditableSummary,
		"should match updated editable summary",
	)
	assert.Equal(
		p.Depositor,
		*input.Depositor,
		"should match updated depositor (he's gold)",
	)
	assert.ElementsMatch(
		p.Genes,
		mocks.MockUpdateStrainAttributes.Genes,
		"should match existing genes list",
	)
	assert.ElementsMatch(
		p.Dbxrefs,
		mocks.MockUpdateStrainAttributes.Dbxrefs,
		"should match existing dbxrefs",
	)
	assert.ElementsMatch(
		p.Publications,
		mocks.MockUpdateStrainAttributes.Publications,
		"should match existing publications",
	)
	assert.Equal(
		p.Label,
		mocks.MockUpdateStrainAttributes.Label,
		"should match existing label",
	)
	assert.Equal(
		p.Species,
		mocks.MockUpdateStrainAttributes.Species,
		"should match existing species",
	)
	assert.Equal(
		p.Plasmid,
		&mocks.MockUpdateStrainAttributes.Plasmid,
		"should match plasmid",
	)
}

func TestGetOntology(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	assert.Equal(
		resolverutils.GetOntology("phenotype"),
		registry.PhenoOntology,
		"should return phenotype ontology",
	)
	assert.Equal(
		resolverutils.GetOntology("characteristic"),
		registry.StrainCharOnto,
		"should return strain characteristics ontology",
	)
	assert.Equal(
		resolverutils.GetOntology("strain_inventory"),
		registry.StrainInvOnto,
		"should return strain inventory ontology",
	)
	assert.Equal(
		resolverutils.GetOntology("plasmid_inventory"),
		registry.PlasmidInvOnto,
		"should return plasmid inventory ontology",
	)
	assert.Equal(
		resolverutils.GetOntology("banana"),
		"invalid ontology",
		"should return invalid ontology",
	)
}

func sliceConverter(s []string) []*string {
	c := []*string{}
	// need to use for loop here, not range
	// https://github.com/golang/go/issues/22791#issuecomment-345391395
	for i := 0; i < len(s); i++ {
		c = append(c, &s[i])
	}
	return c
}

func TestListPlasmidsPassesExpectedFilter(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := new(clients.StockServiceClient)
	mockedStockClient.On(
		"ListPlasmids",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.MatchedBy(func(params *pb.StockParameters) bool {
			return params.Cursor == 0 &&
				params.Limit == 10 &&
				params.Filter == "summary=~GoldenBraid"
		}),
	).Return(mocks.MockPlasmidCollection(), nil)

	reg := &stockClientRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}
	cursor := 0
	limit := 10
	summary := "GoldenBraid"

	result, err := resolver.ListPlasmids(context.Background(), &cursor, &limit, &models.PlasmidListFilter{
		Summary:     &summary,
		PlasmidType: models.PlasmidTypeAll,
	})
	require.NoError(err)
	require.Len(result.Plasmids, 3)
	mockedStockClient.AssertExpectations(t)
}

func TestListPlasmidsNilFilter(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := new(clients.StockServiceClient)
	mockedStockClient.On(
		"ListPlasmids",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.MatchedBy(func(params *pb.StockParameters) bool {
			return params.Cursor == 0 &&
				params.Limit == 10 &&
				params.Filter == ""
		}),
	).Return(mocks.MockPlasmidCollection(), nil)

	reg := &stockClientRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}
	cursor := 0
	limit := 10

	result, err := resolver.ListPlasmids(context.Background(), &cursor, &limit, nil)
	require.NoError(err)
	require.Len(result.Plasmids, 3)
	mockedStockClient.AssertExpectations(t)
}

func TestListPlasmidsUnsupportedInStockFilter(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	resolver := &QueryResolver{
		Registry: &mocks.MockRegistry{},
		Logger:   mocks.TestLogger(),
	}
	cursor := 0
	limit := 10
	inStock := true

	ctx := graphql.WithResponseContext(context.Background(), graphql.DefaultErrorPresenter, graphql.DefaultRecover)
	result, err := resolver.ListPlasmids(ctx, &cursor, &limit, &models.PlasmidListFilter{
		InStock:     &inStock,
		PlasmidType: models.PlasmidTypeAll,
	})
	require.Error(err)
	require.Contains(err.Error(), "in_stock filter is not yet supported")
	require.Empty(result.Plasmids)
}

func TestListPlasmidsUnsupportedIDFilter(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	resolver := &QueryResolver{
		Registry: &mocks.MockRegistry{},
		Logger:   mocks.TestLogger(),
	}
	cursor := 0
	limit := 10
	id := "DBP123456"

	ctx := graphql.WithResponseContext(context.Background(), graphql.DefaultErrorPresenter, graphql.DefaultRecover)
	result, err := resolver.ListPlasmids(ctx, &cursor, &limit, &models.PlasmidListFilter{
		ID:          &id,
		PlasmidType: models.PlasmidTypeAll,
	})
	require.Error(err)
	require.Contains(err.Error(), "id filter is not yet supported")
	require.Empty(result.Plasmids)
}

func TestListPlasmidsRegularTypeFilter(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := new(clients.StockServiceClient)
	mockedStockClient.On(
		"ListPlasmids",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.MatchedBy(func(params *pb.StockParameters) bool {
			return params.Cursor == 0 &&
				params.Limit == 10 &&
				params.Filter == "tag===vector"
		}),
	).Return(mocks.MockPlasmidCollection(), nil)

	reg := &stockClientRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}
	cursor := 0
	limit := 10

	result, err := resolver.ListPlasmids(context.Background(), &cursor, &limit, &models.PlasmidListFilter{
		PlasmidType: models.PlasmidTypeRegular,
	})
	require.NoError(err)
	require.Len(result.Plasmids, 3)
	mockedStockClient.AssertExpectations(t)
}

func TestListPlasmidsGoldenBraidTypeFilter(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := new(clients.StockServiceClient)
	mockedStockClient.On(
		"ListPlasmids",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.MatchedBy(func(params *pb.StockParameters) bool {
			return params.Cursor == 0 &&
				params.Limit == 10 &&
				params.Filter == "tag===GB vector"
		}),
	).Return(mocks.MockPlasmidCollection(), nil)

	reg := &stockClientRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}
	cursor := 0
	limit := 10

	result, err := resolver.ListPlasmids(context.Background(), &cursor, &limit, &models.PlasmidListFilter{
		PlasmidType: models.PlasmidTypeGoldenBraid,
	})
	require.NoError(err)
	require.Len(result.Plasmids, 3)
	mockedStockClient.AssertExpectations(t)
}
