package resolver

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/graphql"

	anno "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks/clients"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/graphql/resolverutils"
	"github.com/dictyBase/graphql-server/internal/registry"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// The stock client mock serves a plasmid whose depositor is
	// kenny@bania.com (see mockPlasmidWithParams); assert against the
	// literal instead of the shared fixture so the test does not depend
	// on mock-side mutation ordering.
	assert.Equal(
		p.Depositor,
		"kenny@bania.com",
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
	// Range over the index, never over the values: ranging values would
	// take the address of the loop copy, see
	// https://github.com/golang/go/issues/22791#issuecomment-345391395
	for i := range len(s) {
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

// bacterialStrainRegistry injects both annotation and stock clients for
// tests that exercise the bacterial annotation pipeline.
type bacterialStrainRegistry struct {
	*mocks.MockRegistry
	annoClient  anno.TaggedAnnotationServiceClient
	stockClient pb.StockServiceClient
}

func (r *bacterialStrainRegistry) GetAnnotationClient(
	_ string,
) anno.TaggedAnnotationServiceClient {
	return r.annoClient
}

func (r *bacterialStrainRegistry) GetStockClient(
	_ string,
) pb.StockServiceClient {
	return r.stockClient
}

// TestListStrainsBacterialRoutesToAnnotationService verifies that
// StrainTypeBACTERIAL calls ListAnnotations with the bacterial food source
// filter, deduplicates IDs, then calls ListStrainsByIds.
func TestListStrainsBacterialRoutesToAnnotationService(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedAnno := new(clients.TaggedAnnotationServiceClient)
	mockedAnno.On(
		"ListAnnotations",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.MatchedBy(func(p *anno.ListParameters) bool {
			return p.Filter == resolverutils.BacterialAnnotationFilter()
		}),
	).Return(mocks.MockBacterialStrainListAnno(), nil)

	mockedStock := new(clients.StockServiceClient)
	mockedStock.On(
		"ListStrainsByIds",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		// 3 annotations deduplicate to 2 unique IDs
		mock.MatchedBy(func(ids *pb.StockIdList) bool { return len(ids.Id) == 2 }),
	).Return(mocks.MockStrainList(), nil)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		annoClient:   mockedAnno,
		stockClient:  mockedStock,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}
	cursor, limit := 0, 10

	result, err := resolver.ListStrains(
		context.Background(), &cursor, &limit,
		&models.StrainListFilter{StrainType: models.StrainTypeBacterial},
	)

	require.NoError(err)
	require.Len(result.Strains, 2)
	require.Equal(2, result.TotalCount)
	mockedAnno.AssertExpectations(t)
	mockedStock.AssertExpectations(t)
	mockedStock.AssertNotCalled(t, "ListStrains")
}

// TestListStrainsAllExcludesBacterial verifies that StrainTypeAll filter DSL
// does not contain any bacterial tag and the annotation service is never called.
func TestListStrainsAllExcludesBacterial(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStock := new(clients.StockServiceClient)
	mockedStock.On(
		"ListStrains",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.MatchedBy(func(p *pb.StockParameters) bool {
			return !strings.Contains(p.Filter, registry.BacterialStrainTag) &&
				!strings.Contains(p.Filter, registry.BacterialFoodSourceTag) &&
				strings.Contains(p.Filter, registry.GeneralStrainTag) &&
				strings.Contains(p.Filter, registry.GwdiStrainTag)
		}),
	).Return(mocks.MockStrainCollection(), nil)

	mockedAnno := new(clients.TaggedAnnotationServiceClient)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStock,
		annoClient:   mockedAnno,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}
	cursor, limit := 0, 10

	result, err := resolver.ListStrains(
		context.Background(), &cursor, &limit,
		&models.StrainListFilter{StrainType: models.StrainTypeAll},
	)
	require.NoError(err)
	require.Len(result.Strains, 3)
	mockedStock.AssertExpectations(t)
	mockedAnno.AssertNotCalled(t, "ListAnnotations")
}

func TestListStockSuggestions(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	resolver := &QueryResolver{
		Registry: &mocks.MockRegistry{},
		Logger:   mocks.TestLogger(),
	}

	result, err := resolver.ListStockSuggestions(
		context.Background(),
		"ax4",
		nil,
		nil,
	)
	require.NoError(err)
	require.Equal(4, result.TotalCount, "should match merged unique count")
	require.NotNil(result.Limit, "should resolve limit from meta")
	require.Equal(5, *result.Limit, "should match effective limit from meta")
	require.Len(result.Suggestions, 4, "should have four merged suggestions")

	first := result.Suggestions[0]
	require.Equal("DBS0236126", first.ID, "should keep highest score for duplicate DBS0236126")
	require.Equal(
		models.StockEntityTypeStrain,
		first.Entity,
		"should map strain entity",
	)
	require.Equal(
		models.StockSearchFieldEnumStockID,
		first.Field,
		"should map stock_id field",
	)
	require.Equal("DBS0236126", first.DisplayText, "should match display text")
	require.Equal(1.0, first.Score, "should match score")

	second := result.Suggestions[1]
	require.Equal("DBS0235559", second.ID, "should include annotation suggestion")
	require.Equal(models.StockEntityTypeStrain, second.Entity)
	require.Equal(models.StockSearchFieldEnumSynonym, second.Field)
	require.Equal(
		"catenin null",
		second.DisplayText,
		"should preserve annotation display text",
	)
	require.Equal(0.75, second.Score, "should preserve annotation score")

	third := result.Suggestions[2]
	require.Equal("DBS0236127", third.ID, "should include distinct strain row")
	require.Equal(models.StockEntityTypeStrain, third.Entity)
	require.Equal(models.StockSearchFieldEnumLabel, third.Field)

	fourth := result.Suggestions[3]
	require.Equal("DBP0236126", fourth.ID, "should include plasmid suggestion")
	require.Equal(
		models.StockEntityTypePlasmid,
		fourth.Entity,
		"should map plasmid entity",
	)
	require.Equal(
		models.StockSearchFieldEnumName,
		fourth.Field,
		"should map plasmid name field",
	)
}

func TestListStockSuggestionsPassesExpectedParams(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := new(clients.StockServiceClient)
	mockedStockClient.On(
		"AutocompleteStock",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.MatchedBy(func(params *pb.StockAutocompleteParameters) bool {
			attr := params.Data.Attributes
			return params.Data.Type == "stock" &&
				attr.Query == "dbp1" &&
				attr.Limit == 25 &&
				attr.Entity == pb.StockEntity_STOCK_ENTITY_PLASMID
		}),
	).Return(mocks.MockStockSuggestionCollection(), nil)

	mockedAnnoClient := new(clients.TaggedAnnotationServiceClient)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}
	limit := 25
	entity := models.StockEntityTypePlasmid

	result, err := resolver.ListStockSuggestions(
		context.Background(),
		"dbp1",
		&limit,
		&entity,
	)
	require.NoError(err)
	require.Len(result.Suggestions, 3)
	mockedStockClient.AssertExpectations(t)
	mockedAnnoClient.AssertNotCalled(t, "AutocompleteAnnotation")
}

func TestListStockSuggestionsDefaultsToBothEntities(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := new(clients.StockServiceClient)
	mockedStockClient.On(
		"AutocompleteStock",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.MatchedBy(func(params *pb.StockAutocompleteParameters) bool {
			attr := params.Data.Attributes
			return attr.Query == "ax4" &&
				attr.Limit == 0 &&
				attr.Entity == pb.StockEntity_STOCK_ENTITY_UNSPECIFIED
		}),
	).Return(mocks.MockStockSuggestionCollection(), nil)

	mockedAnnoClient := new(clients.TaggedAnnotationServiceClient)
	mockedAnnoClient.On(
		"AutocompleteAnnotation",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.MatchedBy(func(params *anno.AnnotationAutocompleteParameters) bool {
			attr := params.Data.Attributes
			return params.Data.Type == "annotation" &&
				attr.Query == "ax4" &&
				attr.Limit == 0
		}),
	).Return(mocks.MockAnnotationSuggestionCollection(), nil)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}
	entity := models.StockEntityTypeAll

	result, err := resolver.ListStockSuggestions(
		context.Background(),
		"ax4",
		nil,
		&entity,
	)
	require.NoError(err)
	require.Len(result.Suggestions, 4)
	mockedStockClient.AssertExpectations(t)
	mockedAnnoClient.AssertExpectations(t)
}

func TestListStockSuggestionsError(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := new(clients.StockServiceClient)
	mockedStockClient.On(
		"AutocompleteStock",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.AnythingOfType("*stock.StockAutocompleteParameters"),
	).Return(nil, fmt.Errorf("stock service unavailable"))

	mockedAnnoClient := new(clients.TaggedAnnotationServiceClient)
	mockedAnnoClient.On(
		"AutocompleteAnnotation",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.AnythingOfType("*annotation.AnnotationAutocompleteParameters"),
	).Return(mocks.MockAnnotationSuggestionCollection(), nil)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	result, err := resolver.ListStockSuggestions(
		context.Background(),
		"ax4",
		nil,
		nil,
	)
	require.Error(err)
	require.ErrorContains(err, "stock service unavailable")
	require.NotNil(result, "should return empty suggestion list on error")
	mockedStockClient.AssertExpectations(t)
	mockedAnnoClient.AssertExpectations(t)
}

func TestListStockSuggestionsInvalidEntityEnum(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := new(clients.StockServiceClient)
	mockedStockClient.On(
		"AutocompleteStock",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.AnythingOfType("*stock.StockAutocompleteParameters"),
	).Return(&pb.StockSuggestionCollection{
		Data: []*pb.StockSuggestion{
			{
				Id:          "DBS0236126",
				Entity:      pb.StockEntity(99),
				Field:       pb.StockSearchField_STOCK_SEARCH_FIELD_STOCK_ID,
				DisplayText: "DBS0236126",
				Score:       1,
			},
		},
		Meta: &pb.Meta{Limit: 5, Total: 1},
	}, nil)

	mockedAnnoClient := new(clients.TaggedAnnotationServiceClient)
	mockedAnnoClient.On(
		"AutocompleteAnnotation",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.AnythingOfType("*annotation.AnnotationAutocompleteParameters"),
	).Return(mocks.MockAnnotationSuggestionCollection(), nil)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	_, err := resolver.ListStockSuggestions(
		context.Background(),
		"ax4",
		nil,
		nil,
	)
	require.Error(err)
	require.ErrorContains(err, "unexpected stock entity")
	mockedStockClient.AssertExpectations(t)
}

func TestListStockSuggestionsInvalidFieldEnum(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := new(clients.StockServiceClient)
	mockedStockClient.On(
		"AutocompleteStock",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.AnythingOfType("*stock.StockAutocompleteParameters"),
	).Return(&pb.StockSuggestionCollection{
		Data: []*pb.StockSuggestion{
			{
				Id:          "DBS0236126",
				Entity:      pb.StockEntity_STOCK_ENTITY_STRAIN,
				Field:       pb.StockSearchField(99),
				DisplayText: "DBS0236126",
				Score:       1,
			},
		},
		Meta: &pb.Meta{Limit: 5, Total: 1},
	}, nil)

	mockedAnnoClient := new(clients.TaggedAnnotationServiceClient)
	mockedAnnoClient.On(
		"AutocompleteAnnotation",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.AnythingOfType("*annotation.AnnotationAutocompleteParameters"),
	).Return(mocks.MockAnnotationSuggestionCollection(), nil)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	_, err := resolver.ListStockSuggestions(
		context.Background(),
		"ax4",
		nil,
		nil,
	)
	require.Error(err)
	require.ErrorContains(err, "unexpected search field")
	mockedStockClient.AssertExpectations(t)
}

func TestListStockSuggestionsStrainCallsBothServices(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := new(clients.StockServiceClient)
	mockedStockClient.On(
		"AutocompleteStock",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.MatchedBy(func(params *pb.StockAutocompleteParameters) bool {
			attr := params.Data.Attributes
			return params.Data.Type == "stock" &&
				attr.Query == "ax4" &&
				attr.Limit == 10 &&
				attr.Entity == pb.StockEntity_STOCK_ENTITY_STRAIN
		}),
	).Return(mocks.MockStockSuggestionCollection(), nil)

	mockedAnnoClient := new(clients.TaggedAnnotationServiceClient)
	mockedAnnoClient.On(
		"AutocompleteAnnotation",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.MatchedBy(func(params *anno.AnnotationAutocompleteParameters) bool {
			attr := params.Data.Attributes
			return params.Data.Type == "annotation" &&
				attr.Query == "ax4" &&
				attr.Limit == 10
		}),
	).Return(mocks.MockAnnotationSuggestionCollection(), nil)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}
	limit := 10
	entity := models.StockEntityTypeStrain

	result, err := resolver.ListStockSuggestions(
		context.Background(),
		"ax4",
		&limit,
		&entity,
	)
	require.NoError(err)
	require.Len(result.Suggestions, 4)
	require.Equal(10, *result.Limit)
	mockedStockClient.AssertExpectations(t)
	mockedAnnoClient.AssertExpectations(t)
}

func TestListStockSuggestionsBothEmpty(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := new(clients.StockServiceClient)
	mockedStockClient.On(
		"AutocompleteStock",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.AnythingOfType("*stock.StockAutocompleteParameters"),
	).Return(&pb.StockSuggestionCollection{
		Data: nil,
		Meta: &pb.Meta{Limit: 5, Total: 0},
	}, nil)

	mockedAnnoClient := new(clients.TaggedAnnotationServiceClient)
	mockedAnnoClient.On(
		"AutocompleteAnnotation",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.AnythingOfType("*annotation.AnnotationAutocompleteParameters"),
	).Return(&anno.AnnotationSuggestionCollection{
		Data: nil,
		Meta: &anno.Meta{Limit: 5, Total: 0},
	}, nil)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	result, err := resolver.ListStockSuggestions(
		context.Background(),
		"zzz",
		nil,
		nil,
	)
	require.NoError(err)
	require.Empty(result.Suggestions)
	require.Equal(5, *result.Limit)
	require.Equal(0, result.TotalCount)
	mockedStockClient.AssertExpectations(t)
	mockedAnnoClient.AssertExpectations(t)
}

func TestListStockSuggestionsAnnotationError(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := new(clients.StockServiceClient)
	mockedStockClient.On(
		"AutocompleteStock",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.AnythingOfType("*stock.StockAutocompleteParameters"),
	).Return(mocks.MockStockSuggestionCollection(), nil)

	mockedAnnoClient := new(clients.TaggedAnnotationServiceClient)
	mockedAnnoClient.On(
		"AutocompleteAnnotation",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.AnythingOfType("*annotation.AnnotationAutocompleteParameters"),
	).Return(nil, fmt.Errorf("annotation service unavailable"))

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	result, err := resolver.ListStockSuggestions(
		context.Background(),
		"ax4",
		nil,
		nil,
	)
	require.Error(err)
	require.ErrorContains(err, "annotation service unavailable")
	require.NotNil(result, "should return empty suggestion list on error")
	mockedStockClient.AssertExpectations(t)
	mockedAnnoClient.AssertExpectations(t)
}

func TestListStockSuggestionsInvalidAnnotationFieldEnum(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := new(clients.StockServiceClient)
	mockedStockClient.On(
		"AutocompleteStock",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.AnythingOfType("*stock.StockAutocompleteParameters"),
	).Return(mocks.MockStockSuggestionCollection(), nil)

	mockedAnnoClient := new(clients.TaggedAnnotationServiceClient)
	mockedAnnoClient.On(
		"AutocompleteAnnotation",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.AnythingOfType("*annotation.AnnotationAutocompleteParameters"),
	).Return(&anno.AnnotationSuggestionCollection{
		Data: []*anno.AnnotationSuggestion{
			{
				EntryId:     "DBS0236126",
				Field:       anno.AnnotationSearchField(99),
				DisplayText: "bad",
				Score:       1.0,
			},
		},
		Meta: &anno.Meta{Limit: 5, Total: 1},
	}, nil)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	_, err := resolver.ListStockSuggestions(
		context.Background(),
		"ax4",
		nil,
		nil,
	)
	require.Error(err)
	require.ErrorContains(err, "unexpected annotation search field")
	mockedStockClient.AssertExpectations(t)
	mockedAnnoClient.AssertExpectations(t)
}

func TestListStockSuggestionsConcurrentFetch(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	stockEntered := make(chan struct{})
	annoEntered := make(chan struct{})

	mockedStockClient := new(clients.StockServiceClient)
	mockedStockClient.On(
		"AutocompleteStock",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.AnythingOfType("*stock.StockAutocompleteParameters"),
	).Run(func(_ mock.Arguments) {
		close(stockEntered)
		<-annoEntered
	}).Return(mocks.MockStockSuggestionCollection(), nil)

	mockedAnnoClient := new(clients.TaggedAnnotationServiceClient)
	mockedAnnoClient.On(
		"AutocompleteAnnotation",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.AnythingOfType("*annotation.AnnotationAutocompleteParameters"),
	).Run(func(_ mock.Arguments) {
		close(annoEntered)
		<-stockEntered
	}).Return(mocks.MockAnnotationSuggestionCollection(), nil)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	result, err := resolver.ListStockSuggestions(
		context.Background(),
		"ax4",
		nil,
		nil,
	)
	require.NoError(err)
	require.Len(result.Suggestions, 4)
	mockedStockClient.AssertExpectations(t)
	mockedAnnoClient.AssertExpectations(t)
}

func TestSearchStocks(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	resolver := &QueryResolver{
		Registry: &mocks.MockRegistry{},
		Logger:   mocks.TestLogger(),
	}

	result, err := resolver.SearchStocks(
		context.Background(),
		"csa",
		nil,
		nil,
	)
	require.NoError(err)
	require.Equal(4, result.TotalCount, "should match merged unique count")
	require.NotNil(result.Limit, "should resolve limit from meta")
	require.Equal(50, *result.Limit, "should match effective limit from meta")
	require.Len(result.Results, 4, "should have four merged results")

	first := result.Results[0]
	require.Equal(mocks.MockStockSearchID, first.ID, "should keep highest score for duplicate DBS0236126")
	require.Equal(
		models.StockEntityTypeStrain,
		first.Entity,
		"should map strain entity",
	)
	require.Equal(
		models.StockSearchFieldEnumStockID,
		first.Field,
		"should map stock_id field",
	)
	require.Equal(mocks.MockStockSearchID, first.DisplayText, "should match display text")
	require.Equal(1.0, first.Score, "should match score")
	require.Equal(
		"GWDI strain carrying a csA deletion.",
		first.Summary,
		"should match summary",
	)
	require.Equal("AX4", first.StrainLabel, "should match strain label")

	second := result.Results[1]
	require.Equal("DBS0235559", second.ID, "should include enriched annotation result")
	require.Equal(models.StockEntityTypeStrain, second.Entity)
	require.Equal(models.StockSearchFieldEnumPhenotype, second.Field)
	require.Equal("delayed culmination", second.DisplayText)
	require.Equal(0.75, second.Score)
	require.Equal(
		"Catenin null mutant strain.",
		second.Summary,
		"should enrich summary from ListStrainsByIds",
	)
	require.Equal(
		"aarA-",
		second.StrainLabel,
		"should enrich strain label from ListStrainsByIds",
	)

	third := result.Results[2]
	require.Equal("DBS0236127", third.ID)
	require.Equal(
		models.StockSearchFieldEnumSummary,
		third.Field,
		"should map summary field",
	)

	plasmid := result.Results[3]
	require.Equal(
		models.StockEntityTypePlasmid,
		plasmid.Entity,
		"should map plasmid entity",
	)
	require.Empty(plasmid.StrainLabel, "plasmid results carry no strain label")
	require.Empty(plasmid.Summary, "should keep empty summary")
}

func TestSearchStocksPlasmidSkipsAnnotationSearch(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		SearchStock(
			mock.Anything,
			mock.MatchedBy(func(params *pb.StockSearchParameters) bool {
				attr := params.Data.Attributes
				return params.Data.Type == "stock" &&
					attr.Query == "pdv" &&
					attr.Limit == 25 &&
					attr.Entity == pb.StockEntity_STOCK_ENTITY_PLASMID
			}),
		).
		Return(mocks.MockStockSearchResultCollection(), nil)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}
	limit := 25
	entity := models.StockEntityTypePlasmid

	result, err := resolver.SearchStocks(
		context.Background(),
		"pdv",
		&limit,
		&entity,
	)
	require.NoError(err)
	require.Len(result.Results, 3)
	mockedAnnoClient.AssertNotCalled(t, "SearchAnnotation")
	mockedStockClient.AssertNotCalled(t, "ListStrainsByIds")
}

func TestSearchStocksShortQuerySkipsAnnotationSearch(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		SearchStock(
			mock.Anything,
			mock.MatchedBy(func(params *pb.StockSearchParameters) bool {
				return params.Data.Attributes.Query == "x "
			}),
		).
		Return(mocks.MockStockSearchResultCollection(), nil)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	result, err := resolver.SearchStocks(
		context.Background(),
		"x ",
		nil,
		nil,
	)
	require.NoError(err)
	require.Len(result.Results, 3, "should return stock-only hits")
	require.Equal(mocks.MockStockSearchID, result.Results[0].ID)
	require.Equal("DBS0236127", result.Results[1].ID)
	require.Equal("DBP0236126", result.Results[2].ID)
	mockedAnnoClient.AssertNotCalled(t, "SearchAnnotation")
	mockedStockClient.AssertNotCalled(t, "ListStrainsByIds")
}

func TestSearchStocksPassesExpectedParams(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		SearchStock(
			mock.Anything,
			mock.MatchedBy(func(params *pb.StockSearchParameters) bool {
				attr := params.Data.Attributes
				return params.Data.Type == "stock" &&
					attr.Query == "csa" &&
					attr.Limit == 25 &&
					attr.Entity == pb.StockEntity_STOCK_ENTITY_STRAIN
			}),
		).
		Return(mocks.MockStockSearchResultCollection(), nil)
	mockedStockClient.EXPECT().
		ListStrainsByIds(
			mock.Anything,
			&pb.StockIdList{Id: []string{mocks.MockAnnoSearchStrainID, mocks.MockStockSearchID}},
		).
		Return(mocks.MockSearchStrainList(), nil)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchAnnotation(
			mock.Anything,
			mock.MatchedBy(func(params *anno.AnnotationSearchParameters) bool {
				attr := params.Data.Attributes
				return params.Data.Type == "annotation" &&
					attr.Query == "csa" &&
					attr.Limit == 25
			}),
		).
		Return(mocks.MockAnnotationSearchResultCollection(), nil)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}
	limit := 25
	entity := models.StockEntityTypeStrain

	result, err := resolver.SearchStocks(
		context.Background(),
		"csa",
		&limit,
		&entity,
	)
	require.NoError(err)
	require.Len(result.Results, 4)
	require.Equal(25, *result.Limit)
}

func TestSearchStocksDefaultsToBothEntities(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		SearchStock(
			mock.Anything,
			mock.MatchedBy(func(params *pb.StockSearchParameters) bool {
				attr := params.Data.Attributes
				return attr.Query == "csa" &&
					attr.Limit == 0 &&
					attr.Entity == pb.StockEntity_STOCK_ENTITY_UNSPECIFIED
			}),
		).
		Return(mocks.MockStockSearchResultCollection(), nil)
	mockedStockClient.EXPECT().
		ListStrainsByIds(
			mock.Anything,
			&pb.StockIdList{Id: []string{mocks.MockAnnoSearchStrainID, mocks.MockStockSearchID}},
		).
		Return(mocks.MockSearchStrainList(), nil)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchAnnotation(
			mock.Anything,
			mock.MatchedBy(func(params *anno.AnnotationSearchParameters) bool {
				attr := params.Data.Attributes
				return params.Data.Type == "annotation" &&
					attr.Query == "csa" &&
					attr.Limit == 0
			}),
		).
		Return(mocks.MockAnnotationSearchResultCollection(), nil)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}
	entity := models.StockEntityTypeAll

	result, err := resolver.SearchStocks(
		context.Background(),
		"csa",
		nil,
		&entity,
	)
	require.NoError(err)
	require.Len(result.Results, 4)
}

func TestSearchStocksClampsAnnotationLimitAboveFifty(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		SearchStock(
			mock.Anything,
			mock.MatchedBy(func(params *pb.StockSearchParameters) bool {
				attr := params.Data.Attributes
				return attr.Query == "csa" &&
					attr.Limit == 75 &&
					attr.Entity == pb.StockEntity_STOCK_ENTITY_UNSPECIFIED
			}),
		).
		Return(mocks.MockStockSearchResultCollection(), nil)
	mockedStockClient.EXPECT().
		ListStrainsByIds(
			mock.Anything,
			&pb.StockIdList{Id: []string{mocks.MockAnnoSearchStrainID, mocks.MockStockSearchID}},
		).
		Return(mocks.MockSearchStrainList(), nil)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchAnnotation(
			mock.Anything,
			mock.MatchedBy(func(params *anno.AnnotationSearchParameters) bool {
				attr := params.Data.Attributes
				return attr.Query == "csa" && attr.Limit == 50
			}),
		).
		Return(mocks.MockAnnotationSearchResultCollection(), nil)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}
	limit := 75

	result, err := resolver.SearchStocks(
		context.Background(),
		"csa",
		&limit,
		nil,
	)
	require.NoError(err)
	require.Equal(75, *result.Limit)
	require.Len(result.Results, 4)
}

func TestSearchStocksEmptyAnnotationSkipsListStrainsByIds(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		SearchStock(
			mock.Anything,
			mock.AnythingOfType("*stock.StockSearchParameters"),
		).
		Return(mocks.MockStockSearchResultCollection(), nil)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchAnnotation(
			mock.Anything,
			mock.AnythingOfType("*annotation.AnnotationSearchParameters"),
		).
		Return(&anno.AnnotationSearchResultCollection{
			Data: []*anno.AnnotationSearchResult{},
			Meta: &anno.Meta{Limit: 50, Total: 0},
		}, nil)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	result, err := resolver.SearchStocks(
		context.Background(),
		"csa",
		nil,
		nil,
	)
	require.NoError(err)
	require.Len(result.Results, 3)
	mockedStockClient.AssertNotCalled(t, "ListStrainsByIds")
}

func TestSearchStocksDropsStaleAnnotationStrainID(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		SearchStock(
			mock.Anything,
			mock.AnythingOfType("*stock.StockSearchParameters"),
		).
		Return(&pb.StockSearchResultCollection{
			Data: []*pb.StockSearchResult{},
			Meta: &pb.Meta{Limit: 50, Total: 0},
		}, nil)
	mockedStockClient.EXPECT().
		ListStrainsByIds(
			mock.Anything,
			&pb.StockIdList{Id: []string{mocks.MockAnnoSearchStrainID, "DBS9999999"}},
		).
		Return(&pb.StrainList{
			Data: []*pb.StrainList_Data{
				{
					Type: "strain",
					Id:   mocks.MockAnnoSearchStrainID,
					Attributes: &pb.StrainAttributes{
						Summary: "Catenin null mutant strain.",
						Label:   "aarA-",
					},
				},
			},
		}, nil)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchAnnotation(
			mock.Anything,
			mock.AnythingOfType("*annotation.AnnotationSearchParameters"),
		).
		Return(&anno.AnnotationSearchResultCollection{
			Data: []*anno.AnnotationSearchResult{
				{
					EntryId:     "DBS9999999",
					Field:       anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_PHENOTYPE,
					DisplayText: "delayed culmination",
					Score:       900.0,
				},
				{
					EntryId:     mocks.MockAnnoSearchStrainID,
					Field:       anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_SYNONYM,
					DisplayText: "catenin null",
					Score:       500.0,
				},
			},
			Meta: &anno.Meta{Limit: 50, Total: 2},
		}, nil)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	result, err := resolver.SearchStocks(
		context.Background(),
		"ca",
		nil,
		nil,
	)
	require.NoError(err)
	require.Len(result.Results, 1)
	require.Equal(mocks.MockAnnoSearchStrainID, result.Results[0].ID)
	require.Equal(1, result.TotalCount)
}

func TestSearchStocksAllStaleAnnotationIDsNotFound(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		SearchStock(
			mock.Anything,
			mock.AnythingOfType("*stock.StockSearchParameters"),
		).
		Return(mocks.MockStockSearchResultCollection(), nil)
	mockedStockClient.EXPECT().
		ListStrainsByIds(
			mock.Anything,
			&pb.StockIdList{Id: []string{"DBS9999999"}},
		).
		Return(nil, status.Error(codes.NotFound, "could not find any strains"))

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchAnnotation(
			mock.Anything,
			mock.AnythingOfType("*annotation.AnnotationSearchParameters"),
		).
		Return(&anno.AnnotationSearchResultCollection{
			Data: []*anno.AnnotationSearchResult{
				{
					EntryId:     "DBS9999999",
					Field:       anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_PHENOTYPE,
					DisplayText: "delayed culmination",
					Score:       900.0,
				},
			},
			Meta: &anno.Meta{Limit: 50, Total: 1},
		}, nil)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	result, err := resolver.SearchStocks(
		context.Background(),
		"csa",
		nil,
		nil,
	)
	require.NoError(err)
	require.Len(result.Results, 3)
	require.Equal(3, result.TotalCount)
}

func TestSearchStocksError(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		SearchStock(
			mock.Anything,
			mock.AnythingOfType("*stock.StockSearchParameters"),
		).
		Return(nil, fmt.Errorf("stock service unavailable"))

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	result, err := resolver.SearchStocks(
		context.Background(),
		"csa",
		nil,
		nil,
	)
	require.Error(err)
	require.ErrorContains(err, "stock service unavailable")
	require.ErrorContains(err, "search stocks for query")
	require.NotNil(result, "should return empty search list on error")
	mockedAnnoClient.AssertNotCalled(t, "SearchAnnotation")
}

func TestSearchStocksAnnotationError(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		SearchStock(
			mock.Anything,
			mock.AnythingOfType("*stock.StockSearchParameters"),
		).
		Return(mocks.MockStockSearchResultCollection(), nil)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchAnnotation(
			mock.Anything,
			mock.AnythingOfType("*annotation.AnnotationSearchParameters"),
		).
		Return(nil, fmt.Errorf("annotation service unavailable"))

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	result, err := resolver.SearchStocks(
		context.Background(),
		"csa",
		nil,
		nil,
	)
	require.Error(err)
	require.ErrorContains(err, "annotation service unavailable")
	require.ErrorContains(err, "search annotations for query")
	require.NotNil(result)
	mockedStockClient.AssertNotCalled(t, "ListStrainsByIds")
}

func TestSearchStocksListStrainsByIdsError(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		SearchStock(
			mock.Anything,
			mock.AnythingOfType("*stock.StockSearchParameters"),
		).
		Return(mocks.MockStockSearchResultCollection(), nil)
	mockedStockClient.EXPECT().
		ListStrainsByIds(
			mock.Anything,
			mock.AnythingOfType("*stock.StockIdList"),
		).
		Return(nil, fmt.Errorf("strain batch lookup failed"))

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchAnnotation(
			mock.Anything,
			mock.AnythingOfType("*annotation.AnnotationSearchParameters"),
		).
		Return(mocks.MockAnnotationSearchResultCollection(), nil)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	result, err := resolver.SearchStocks(
		context.Background(),
		"csa",
		nil,
		nil,
	)
	require.Error(err)
	require.ErrorContains(err, "strain batch lookup failed")
	require.ErrorContains(err, "fetch strain metadata for query")
	require.NotNil(result)
}

func TestSearchStocksInvalidEntityEnum(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		SearchStock(
			mock.Anything,
			mock.AnythingOfType("*stock.StockSearchParameters"),
		).
		Return(&pb.StockSearchResultCollection{
			Data: []*pb.StockSearchResult{
				{
					Id:          mocks.MockStockSearchID,
					Entity:      pb.StockEntity(99),
					Field:       pb.StockSearchField_STOCK_SEARCH_FIELD_STOCK_ID,
					DisplayText: mocks.MockStockSearchID,
					Score:       1,
				},
			},
			Meta: &pb.Meta{Limit: 50, Total: 1},
		}, nil)
	mockedStockClient.EXPECT().
		ListStrainsByIds(
			mock.Anything,
			mock.AnythingOfType("*stock.StockIdList"),
		).
		Return(mocks.MockSearchStrainList(), nil)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchAnnotation(
			mock.Anything,
			mock.AnythingOfType("*annotation.AnnotationSearchParameters"),
		).
		Return(mocks.MockAnnotationSearchResultCollection(), nil)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	_, err := resolver.SearchStocks(
		context.Background(),
		"csa",
		nil,
		nil,
	)
	require.Error(err)
	require.ErrorContains(err, "unexpected stock entity")
}

func TestSearchStocksInvalidFieldEnum(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		SearchStock(
			mock.Anything,
			mock.AnythingOfType("*stock.StockSearchParameters"),
		).
		Return(&pb.StockSearchResultCollection{
			Data: []*pb.StockSearchResult{
				{
					Id:          mocks.MockStockSearchID,
					Entity:      pb.StockEntity_STOCK_ENTITY_STRAIN,
					Field:       pb.StockSearchField(99),
					DisplayText: mocks.MockStockSearchID,
					Score:       1,
				},
			},
			Meta: &pb.Meta{Limit: 50, Total: 1},
		}, nil)
	mockedStockClient.EXPECT().
		ListStrainsByIds(
			mock.Anything,
			mock.AnythingOfType("*stock.StockIdList"),
		).
		Return(mocks.MockSearchStrainList(), nil)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchAnnotation(
			mock.Anything,
			mock.AnythingOfType("*annotation.AnnotationSearchParameters"),
		).
		Return(mocks.MockAnnotationSearchResultCollection(), nil)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	_, err := resolver.SearchStocks(
		context.Background(),
		"csa",
		nil,
		nil,
	)
	require.Error(err)
	require.ErrorContains(err, "unexpected search field")
}

func TestSearchStocksInvalidAnnotationFieldEnum(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mockedStockClient := clients.NewStockServiceClient(t)
	mockedStockClient.EXPECT().
		SearchStock(
			mock.Anything,
			mock.AnythingOfType("*stock.StockSearchParameters"),
		).
		Return(mocks.MockStockSearchResultCollection(), nil)
	mockedStockClient.EXPECT().
		ListStrainsByIds(
			mock.Anything,
			mock.AnythingOfType("*stock.StockIdList"),
		).
		Return(mocks.MockSearchStrainList(), nil)

	mockedAnnoClient := clients.NewTaggedAnnotationServiceClient(t)
	mockedAnnoClient.EXPECT().
		SearchAnnotation(
			mock.Anything,
			mock.AnythingOfType("*annotation.AnnotationSearchParameters"),
		).
		Return(&anno.AnnotationSearchResultCollection{
			Data: []*anno.AnnotationSearchResult{
				{
					EntryId:     mocks.MockAnnoSearchStrainID,
					Field:       anno.AnnotationSearchField(99),
					DisplayText: "bad enum",
					Score:       1.0,
				},
			},
			Meta: &anno.Meta{Limit: 50, Total: 1},
		}, nil)

	reg := &bacterialStrainRegistry{
		MockRegistry: &mocks.MockRegistry{ConnMap: nil},
		stockClient:  mockedStockClient,
		annoClient:   mockedAnnoClient,
	}
	resolver := &QueryResolver{Registry: reg, Logger: mocks.TestLogger()}

	_, err := resolver.SearchStocks(
		context.Background(),
		"csa",
		nil,
		nil,
	)
	require.Error(err)
	require.ErrorContains(err, "unexpected annotation search field")
}

func TestConvertAnnotationSuggestionItem(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	synItem := &anno.AnnotationSuggestion{
		EntryId:     "DBS0235559",
		Field:       anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_SYNONYM,
		DisplayText: "catenin null",
		Score:       1001.2,
	}
	synRes := convertAnnotationSuggestionItem(synItem)
	require.True(E.IsRight(synRes))
	synModel := F.Pipe1(
		synRes,
		E.GetOrElse(F.Constant1[error](&models.StockSuggestion{})),
	)
	require.Equal("DBS0235559", synModel.ID)
	require.Equal(models.StockEntityTypeStrain, synModel.Entity)
	require.Equal(models.StockSearchFieldEnumSynonym, synModel.Field)
	require.Equal("catenin null", synModel.DisplayText)
	require.Equal(1001.2, synModel.Score)

	genoItem := &anno.AnnotationSuggestion{
		EntryId:     "DBS0236126",
		Field:       anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_GENOTYPE,
		DisplayText: "axeA1,axeB1",
		Score:       0.8,
	}
	genoRes := convertAnnotationSuggestionItem(genoItem)
	require.True(E.IsRight(genoRes))
	genoModel := F.Pipe1(
		genoRes,
		E.GetOrElse(F.Constant1[error](&models.StockSuggestion{})),
	)
	require.Equal("DBS0236126", genoModel.ID)
	require.Equal(models.StockEntityTypeStrain, genoModel.Entity)
	require.Equal(models.StockSearchFieldEnumGenotype, genoModel.Field)

	badItem := &anno.AnnotationSuggestion{
		EntryId:     "DBS0236126",
		Field:       anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_PHENOTYPE,
		DisplayText: "delayed culmination",
		Score:       0.5,
	}
	badRes := convertAnnotationSuggestionItem(badItem)
	require.True(E.IsLeft(badRes))
}

func TestMergeSuggestionsDeduplicatesByIDKeepingHighestScore(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	stockItems := []*models.StockSuggestion{
		{
			ID:          "DBS0236126",
			Entity:      models.StockEntityTypeStrain,
			Field:       models.StockSearchFieldEnumLabel,
			DisplayText: "AX4",
			Score:       0.5,
		},
		{
			ID:          "DBP0236126",
			Entity:      models.StockEntityTypePlasmid,
			Field:       models.StockSearchFieldEnumName,
			DisplayText: "pDV4",
			Score:       0.25,
		},
	}
	annoItems := []*models.StockSuggestion{
		{
			ID:          "DBS0236126",
			Entity:      models.StockEntityTypeStrain,
			Field:       models.StockSearchFieldEnumSynonym,
			DisplayText: "axenic-4",
			Score:       1001.0,
		},
	}

	merged := mergeSuggestionsFromState(listStockSuggestionsContext{
		stockSuggestions: stockItems,
		annoSuggestions:  annoItems,
		effectiveLimit:   5,
	})
	require.Len(merged, 2)
	require.Equal("DBS0236126", merged[0].ID)
	require.Equal(models.StockSearchFieldEnumSynonym, merged[0].Field)
	require.Equal(1001.0, merged[0].Score)
	require.Equal("DBP0236126", merged[1].ID)
}

func TestMergeSuggestionsTieBreakingAndLimitTruncation(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	stockItems := []*models.StockSuggestion{
		{
			ID:          "DBS000002",
			Entity:      models.StockEntityTypeStrain,
			Field:       models.StockSearchFieldEnumLabel,
			DisplayText: "beta",
			Score:       1.0,
		},
	}
	annoItems := []*models.StockSuggestion{
		{
			ID:          "DBS000001",
			Entity:      models.StockEntityTypeStrain,
			Field:       models.StockSearchFieldEnumSynonym,
			DisplayText: "alpha",
			Score:       1.0,
		},
		{
			ID:          "DBS000003",
			Entity:      models.StockEntityTypeStrain,
			Field:       models.StockSearchFieldEnumGenotype,
			DisplayText: "gamma",
			Score:       0.5,
		},
	}

	merged := mergeSuggestionsFromState(listStockSuggestionsContext{
		stockSuggestions: stockItems,
		annoSuggestions:  annoItems,
		effectiveLimit:   2,
	})
	require.Len(merged, 2)
	require.Equal("DBS000001", merged[0].ID, "equal score breaks tie by ID ASC")
	require.Equal("DBS000002", merged[1].ID)

	empty := mergeSuggestionsFromState(listStockSuggestionsContext{
		effectiveLimit: 5,
	})
	require.Empty(empty)
}

func TestConvertAnnotationSearchResultAllTenFields(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	cases := []struct {
		protoField anno.AnnotationSearchField
		gqlField   models.StockSearchFieldEnum
		display    string
	}{
		{
			protoField: anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_SYNONYM,
			gqlField:   models.StockSearchFieldEnumSynonym,
			display:    "catenin null",
		},
		{
			protoField: anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_GENOTYPE,
			gqlField:   models.StockSearchFieldEnumGenotype,
			display:    "axeA1,axeB1",
		},
		{
			protoField: anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_SYSTEMATIC_NAME,
			gqlField:   models.StockSearchFieldEnumSystematicName,
			display:    "AX4",
		},
		{
			protoField: anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_GENETIC_MODIFICATION,
			gqlField:   models.StockSearchFieldEnumGeneticModification,
			display:    "endogenous deletion",
		},
		{
			protoField: anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_MUTAGENESIS_METHOD,
			gqlField:   models.StockSearchFieldEnumMutagenesisMethod,
			display:    "Homologous Recombination",
		},
		{
			protoField: anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_PHENOTYPE,
			gqlField:   models.StockSearchFieldEnumPhenotype,
			display:    "delayed culmination",
		},
		{
			protoField: anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_ENVIRONMENT,
			gqlField:   models.StockSearchFieldEnumEnvironment,
			display:    "in the presence of caffeine",
		},
		{
			protoField: anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_ASSAY,
			gqlField:   models.StockSearchFieldEnumAssay,
			display:    "confocal microscopy",
		},
		{
			protoField: anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_CHARACTERISTIC,
			gqlField:   models.StockSearchFieldEnumCharacteristic,
			display:    "axenic",
		},
		{
			protoField: anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_PUBLIC_NOTE,
			gqlField:   models.StockSearchFieldEnumPublicNote,
			display:    "fruiting body note",
		},
	}

	metaByID := map[string]strainSearchMetadata{
		"DBS0235559": {
			summary:     "Catenin null mutant strain.",
			strainLabel: "aarA-",
		},
	}

	for _, tc := range cases {
		res := convertAnnoSearchResultsFromState(listStockSearchContext{
			annoCollection: &anno.AnnotationSearchResultCollection{
				Data: []*anno.AnnotationSearchResult{
					{
						EntryId:     "DBS0235559",
						Field:       tc.protoField,
						DisplayText: tc.display,
						Score:       504.2,
					},
				},
			},
			strainMetaByID: metaByID,
		})
		require.True(E.IsRight(res))
		items := F.Pipe1(
			res,
			E.GetOrElse(F.Constant1[error]([]*models.StockSearchResult(nil))),
		)
		require.Len(items, 1)
		require.Equal("DBS0235559", items[0].ID)
		require.Equal(models.StockEntityTypeStrain, items[0].Entity)
		require.Equal(tc.gqlField, items[0].Field)
		require.Equal(tc.display, items[0].DisplayText)
		require.Equal(504.2, items[0].Score)
		require.Equal("Catenin null mutant strain.", items[0].Summary)
		require.Equal("aarA-", items[0].StrainLabel)
	}

	badRes := validateAnnotationSearchItem(&anno.AnnotationSearchResult{
		EntryId:     "DBS0235559",
		Field:       anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_UNSPECIFIED,
		DisplayText: "unknown",
		Score:       0.1,
	})
	require.True(E.IsLeft(badRes))
}

func TestConvertAnnoSearchResultsDropsStaleStrainID(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	res := convertAnnoSearchResultsFromState(listStockSearchContext{
		annoCollection: &anno.AnnotationSearchResultCollection{
			Data: []*anno.AnnotationSearchResult{
				{
					EntryId:     "DBS9999999",
					Field:       anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_PHENOTYPE,
					DisplayText: "delayed culmination",
					Score:       500.0,
				},
				{
					EntryId:     "DBS0235559",
					Field:       anno.AnnotationSearchField_ANNOTATION_SEARCH_FIELD_SYNONYM,
					DisplayText: "catenin null",
					Score:       1000.0,
				},
			},
		},
		strainMetaByID: map[string]strainSearchMetadata{
			"DBS0235559": {
				summary:     "Catenin null mutant strain.",
				strainLabel: "aarA-",
			},
		},
	})
	require.True(E.IsRight(res))
	items := F.Pipe1(
		res,
		E.GetOrElse(F.Constant1[error]([]*models.StockSearchResult(nil))),
	)
	require.Len(items, 1)
	require.Equal("DBS0235559", items[0].ID)
}

func TestResolveAnnotationSearchLimitClampsAtFifty(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	nilLimit := resolveAnnotationSearchLimit(nil)
	require.Equal(int64(0), nilLimit)

	withinCap := 25
	withinLimit := resolveAnnotationSearchLimit(&withinCap)
	require.Equal(int64(25), withinLimit)

	atCap := 50
	atCapLimit := resolveAnnotationSearchLimit(&atCap)
	require.Equal(int64(50), atCapLimit)

	aboveCap := 75
	clampedAnnoLimit := resolveAnnotationSearchLimit(&aboveCap)
	stockLimit := resolveSearchLimit(&aboveCap)
	require.Equal(int64(50), clampedAnnoLimit)
	require.Equal(int64(75), stockLimit)
}

func TestMergeSearchResultsDeduplicatesByIDKeepingHighestScore(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	stockItems := []*models.StockSearchResult{
		{
			ID:          "DBS0236126",
			Entity:      models.StockEntityTypeStrain,
			Field:       models.StockSearchFieldEnumLabel,
			DisplayText: "AX4",
			Score:       0.5,
			Summary:     "GWDI strain carrying a csA deletion.",
			StrainLabel: "AX4",
		},
		{
			ID:          "DBP0236126",
			Entity:      models.StockEntityTypePlasmid,
			Field:       models.StockSearchFieldEnumName,
			DisplayText: "pDV4",
			Score:       0.25,
		},
	}
	annoItems := []*models.StockSearchResult{
		{
			ID:          "DBS0236126",
			Entity:      models.StockEntityTypeStrain,
			Field:       models.StockSearchFieldEnumPhenotype,
			DisplayText: "delayed culmination",
			Score:       504.2,
			Summary:     "GWDI strain carrying a csA deletion.",
			StrainLabel: "AX4",
		},
	}

	merged := mergeSearchResultsFromState(listStockSearchContext{
		stockResults:   stockItems,
		annoResults:    annoItems,
		effectiveLimit: 50,
	})
	require.Len(merged, 2)
	require.Equal("DBS0236126", merged[0].ID)
	require.Equal(models.StockSearchFieldEnumPhenotype, merged[0].Field)
	require.Equal(504.2, merged[0].Score)
	require.Equal("GWDI strain carrying a csA deletion.", merged[0].Summary)
	require.Equal("AX4", merged[0].StrainLabel)
	require.Equal("DBP0236126", merged[1].ID)
}

func TestMergeSearchResultsTieBreakingAndLimitTruncation(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	stockItems := []*models.StockSearchResult{
		{
			ID:          "DBS000002",
			Entity:      models.StockEntityTypeStrain,
			Field:       models.StockSearchFieldEnumLabel,
			DisplayText: "beta",
			Score:       1.0,
		},
	}
	annoItems := []*models.StockSearchResult{
		{
			ID:          "DBS000001",
			Entity:      models.StockEntityTypeStrain,
			Field:       models.StockSearchFieldEnumSynonym,
			DisplayText: "alpha",
			Score:       1.0,
		},
		{
			ID:          "DBS000003",
			Entity:      models.StockEntityTypeStrain,
			Field:       models.StockSearchFieldEnumGenotype,
			DisplayText: "gamma",
			Score:       0.5,
		},
	}

	merged := mergeSearchResultsFromState(listStockSearchContext{
		stockResults:   stockItems,
		annoResults:    annoItems,
		effectiveLimit: 2,
	})
	require.Len(merged, 2)
	require.Equal("DBS000001", merged[0].ID, "equal score breaks tie by ID ASC")
	require.Equal("DBS000002", merged[1].ID)

	empty := mergeSearchResultsFromState(listStockSearchContext{
		effectiveLimit: 50,
	})
	require.Empty(empty)
}
