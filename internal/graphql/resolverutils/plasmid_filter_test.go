package resolverutils

import (
	"testing"

	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/stretchr/testify/require"
)

func TestPlasmidFilterToQueryNilFilter(t *testing.T) {
	t.Parallel()
	query, err := PlasmidFilterToQuery(nil)
	require.NoError(t, err)
	require.Equal(t, "", query)
}

func TestPlasmidFilterToQuerySummaryOnly(t *testing.T) {
	t.Parallel()
	summary := "GoldenBraid"
	query, err := PlasmidFilterToQuery(&models.PlasmidListFilter{
		Summary:     &summary,
		PlasmidType: models.PlasmidTypeRegular,
	})
	require.NoError(t, err)
	require.Equal(t, "summary=~GoldenBraid;ontology==dicty_plasmid_keyword;tag==regular plasmid", query)
}

func TestPlasmidFilterToQueryNameOnly(t *testing.T) {
	t.Parallel()
	name := "pTest"
	query, err := PlasmidFilterToQuery(&models.PlasmidListFilter{
		Name:        &name,
		PlasmidType: models.PlasmidTypeAll,
	})
	require.NoError(t, err)
	require.Equal(t, "plasmid_name===pTest;ontology==dicty_plasmid_keyword;tag==regular plasmid,tag==golden braid", query)
}

func TestPlasmidFilterToQueryCombinedSummaryAndName(t *testing.T) {
	t.Parallel()
	summary := "test"
	name := "pTest"
	query, err := PlasmidFilterToQuery(&models.PlasmidListFilter{
		Summary:     &summary,
		Name:        &name,
		PlasmidType: models.PlasmidTypeGoldenBraid,
	})
	require.NoError(t, err)
	require.Equal(t, "summary=~test;plasmid_name===pTest;ontology==dicty_plasmid_keyword;tag==golden braid", query)
}

func TestPlasmidFilterToQueryPlasmidTypeAll(t *testing.T) {
	t.Parallel()
	query, err := PlasmidFilterToQuery(&models.PlasmidListFilter{
		PlasmidType: models.PlasmidTypeAll,
	})
	require.NoError(t, err)
	require.Equal(t, "ontology==dicty_plasmid_keyword;tag==regular plasmid,tag==golden braid", query)
}

func TestPlasmidFilterToQueryPlasmidTypeRegular(t *testing.T) {
	t.Parallel()
	query, err := PlasmidFilterToQuery(&models.PlasmidListFilter{
		PlasmidType: models.PlasmidTypeRegular,
	})
	require.NoError(t, err)
	require.Equal(t, "ontology==dicty_plasmid_keyword;tag==regular plasmid", query)
}

func TestPlasmidFilterToQueryPlasmidTypeGoldenBraid(t *testing.T) {
	t.Parallel()
	query, err := PlasmidFilterToQuery(&models.PlasmidListFilter{
		PlasmidType: models.PlasmidTypeGoldenBraid,
	})
	require.NoError(t, err)
	require.Equal(t, "ontology==dicty_plasmid_keyword;tag==golden braid", query)
}

func TestPlasmidFilterToQueryInvalidPlasmidType(t *testing.T) {
	t.Parallel()
	query, err := PlasmidFilterToQuery(&models.PlasmidListFilter{
		PlasmidType: models.PlasmidType("INVALID"),
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid plasmid type")
	require.Equal(t, "", query)
}

func TestPlasmidFilterToQueryUnsupportedInStock(t *testing.T) {
	t.Parallel()
	inStock := true
	_, err := PlasmidFilterToQuery(&models.PlasmidListFilter{
		InStock:     &inStock,
		PlasmidType: models.PlasmidTypeRegular,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "in_stock filter is not yet supported")
}

func TestPlasmidFilterToQueryUnsupportedID(t *testing.T) {
	t.Parallel()
	id := "DBP123456"
	_, err := PlasmidFilterToQuery(&models.PlasmidListFilter{
		ID:          &id,
		PlasmidType: models.PlasmidTypeRegular,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "id filter is not yet supported")
}

func TestPlasmidTypeQueryAllTags(t *testing.T) {
	t.Parallel()
	query, err := plasmidTypeQuery(&models.PlasmidListFilter{
		PlasmidType: models.PlasmidTypeAll,
	})
	require.NoError(t, err)
	require.Contains(t, query, registry.DictyPlasmidPropOntology)
	require.Contains(t, query, registry.RegularPlasmidTag)
	require.Contains(t, query, registry.GoldenBraidPlasmidTag)
}
