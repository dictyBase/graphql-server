package resolverutils

import (
	"testing"

	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/stretchr/testify/require"
)

func TestPlasmidFilterToQueryNilFilter(t *testing.T) {
	t.Parallel()
	query, err := PlasmidFilterToQuery(nil)
	require.NoError(t, err)
	require.Equal(t, "", query)
}

func TestPlasmidFilterToQuerySummaryOnlyWithAllType(t *testing.T) {
	t.Parallel()
	summary := "GoldenBraid"
	query, err := PlasmidFilterToQuery(&models.PlasmidListFilter{
		Summary:     &summary,
		PlasmidType: models.PlasmidTypeAll,
	})
	require.NoError(t, err)
	require.Equal(t, "summary=~GoldenBraid", query)
}

func TestPlasmidFilterToQueryNameOnlyWithAllType(t *testing.T) {
	t.Parallel()
	name := "pTest"
	query, err := PlasmidFilterToQuery(&models.PlasmidListFilter{
		Name:        &name,
		PlasmidType: models.PlasmidTypeAll,
	})
	require.NoError(t, err)
	require.Equal(t, "plasmid_name===pTest", query)
}

func TestPlasmidFilterToQueryCombinedSummaryAndNameWithAllType(t *testing.T) {
	t.Parallel()
	summary := "test"
	name := "pTest"
	query, err := PlasmidFilterToQuery(&models.PlasmidListFilter{
		Summary:     &summary,
		Name:        &name,
		PlasmidType: models.PlasmidTypeAll,
	})
	require.NoError(t, err)
	require.Equal(t, "summary=~test;plasmid_name===pTest", query)
}

func TestPlasmidFilterToQueryPlasmidTypeAllOnly(t *testing.T) {
	t.Parallel()
	query, err := PlasmidFilterToQuery(&models.PlasmidListFilter{
		PlasmidType: models.PlasmidTypeAll,
	})
	require.NoError(t, err)
	require.Equal(t, "", query)
}

func TestPlasmidFilterToQueryPlasmidTypeRegularUnverified(t *testing.T) {
	t.Parallel()
	_, err := PlasmidFilterToQuery(&models.PlasmidListFilter{
		PlasmidType: models.PlasmidTypeRegular,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "plasmid_type filter is not yet verified")
}

func TestPlasmidFilterToQueryPlasmidTypeGoldenBraidUnverified(t *testing.T) {
	t.Parallel()
	_, err := PlasmidFilterToQuery(&models.PlasmidListFilter{
		PlasmidType: models.PlasmidTypeGoldenBraid,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "plasmid_type filter is not yet verified")
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
		PlasmidType: models.PlasmidTypeAll,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "in_stock filter is not yet supported")
}

func TestPlasmidFilterToQueryUnsupportedID(t *testing.T) {
	t.Parallel()
	id := "DBP123456"
	_, err := PlasmidFilterToQuery(&models.PlasmidListFilter{
		ID:          &id,
		PlasmidType: models.PlasmidTypeAll,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "id filter is not yet supported")
}
