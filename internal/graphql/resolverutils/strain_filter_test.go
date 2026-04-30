package resolverutils

import (
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	"github.com/dictyBase/graphql-server/internal/graphql/models"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/stretchr/testify/require"
)

func TestBacterialAnnotationFilter(t *testing.T) {
	t.Parallel()
	f := BacterialAnnotationFilter()
	require.Contains(t, f, registry.StrainCharOnto)
	require.Contains(t, f, registry.BacterialFoodSourceTag)
	require.NotContains(t, f, ",")
	require.NotContains(t, f, registry.BacterialStrainTag)
}

func TestStrainFilterToQueryFP_NilFilter(t *testing.T) {
	t.Parallel()
	q, err := E.UnwrapError(StrainFilterToQueryFP(nil)())
	require.NoError(t, err)
	require.Empty(t, q)
}

func TestStrainFilterToQueryFP_AllStrainTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filter   models.StrainListFilter
		expected string
	}{
		{
			name:     "StrainTypeAll",
			filter:   models.StrainListFilter{StrainType: models.StrainTypeAll},
			expected: "ontology==" + registry.DictyStrainPropOntology + ";tag==" + registry.GeneralStrainTag + ",tag==" + registry.GwdiStrainTag,
		},
		{
			name:     "StrainTypeRegular",
			filter:   models.StrainListFilter{StrainType: models.StrainTypeRegular},
			expected: "ontology==" + registry.DictyStrainPropOntology + ";tag==" + registry.GeneralStrainTag,
		},
		{
			name:     "StrainTypeGwdi",
			filter:   models.StrainListFilter{StrainType: models.StrainTypeGwdi},
			expected: "ontology==" + registry.DictyStrainPropOntology + ";tag==" + registry.GwdiStrainTag,
		},
		{
			name:     "StrainTypeBacterial",
			filter:   models.StrainListFilter{StrainType: models.StrainTypeBacterial},
			expected: "ontology==" + registry.StrainCharOnto + ";tag==" + registry.BacterialFoodSourceTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q, err := E.UnwrapError(StrainFilterToQueryFP(&tt.filter)())
			require.NoError(t, err)
			require.Equal(t, tt.expected, q)
		})
	}
}

func TestStrainFilterToQueryFP_AllExcludesBacterial(t *testing.T) {
	t.Parallel()
	q, err := E.UnwrapError(StrainFilterToQueryFP(&models.StrainListFilter{
		StrainType: models.StrainTypeAll,
	})())
	require.NoError(t, err)
	require.Contains(t, q, registry.GeneralStrainTag)
	require.Contains(t, q, registry.GwdiStrainTag)
	require.NotContains(t, q, registry.BacterialStrainTag)
	require.NotContains(t, q, registry.BacterialFoodSourceTag)
}

func TestStrainFilterToQueryFP_WithLabelAndSummary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filter   models.StrainListFilter
		contains []string
	}{
		{
			name: "label with bacterial",
			filter: models.StrainListFilter{
				Label:      strPtr("coli"),
				StrainType: models.StrainTypeBacterial,
			},
			contains: []string{
				"label=~coli",
				registry.BacterialFoodSourceTag,
			},
		},
		{
			name: "summary with regular",
			filter: models.StrainListFilter{
				Summary:    strPtr("knockout"),
				StrainType: models.StrainTypeRegular,
			},
			contains: []string{
				"summary=~knockout",
				registry.GeneralStrainTag,
			},
		},
		{
			name: "label and summary with gwdi",
			filter: models.StrainListFilter{
				Label:      strPtr("axe"),
				Summary:    strPtr("mutant"),
				StrainType: models.StrainTypeGwdi,
			},
			contains: []string{
				"label=~axe",
				"summary=~mutant",
				registry.GwdiStrainTag,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q, err := E.UnwrapError(StrainFilterToQueryFP(&tt.filter)())
			require.NoError(t, err)
			for _, c := range tt.contains {
				require.Contains(t, q, c)
			}
		})
	}
}

func TestStrainFilterToQueryFP_UnsupportedFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		filter      models.StrainListFilter
		errContains string
	}{
		{
			name: "ID not supported",
			filter: models.StrainListFilter{
				ID:         strPtr("DBS123456"),
				StrainType: models.StrainTypeAll,
			},
			errContains: "id filter is not yet supported",
		},
		{
			name: "InStock not supported",
			filter: models.StrainListFilter{
				InStock:    boolPtr(true),
				StrainType: models.StrainTypeAll,
			},
			errContains: "in_stock filter is not yet supported",
		},
		{
			name: "both ID and InStock not supported",
			filter: models.StrainListFilter{
				ID:         strPtr("DBS123456"),
				InStock:    boolPtr(false),
				StrainType: models.StrainTypeBacterial,
			},
			errContains: "id filter is not yet supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := E.UnwrapError(StrainFilterToQueryFP(&tt.filter)())
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.errContains)
		})
	}
}

func TestStrainFilterToQueryFP_InvalidStrainType(t *testing.T) {
	t.Parallel()
	_, err := E.UnwrapError(StrainFilterToQueryFP(&models.StrainListFilter{
		StrainType: models.StrainType("INVALID"),
	})())
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid strain type")
}

func strPtr(s string) *string { return &s }

func boolPtr(b bool) *bool { return &b }
