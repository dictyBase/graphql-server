package stock

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestEveryStockPubsStateCarriesPublicationParams fails when a resolver
// builds a stockPubsState without the NCBI identity values. Without them the
// fetch falls back to unauthenticated NCBI requests.
func TestEveryStockPubsStateCarriesPublicationParams(t *testing.T) {
	t.Parallel()
	files, err := filepath.Glob("*.go")
	require.NoError(t, err)
	built := 0
	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		src, err := os.ReadFile(name)
		require.NoError(t, err)
		text := string(src)
		states := strings.Count(text, "stockPubsState{},")
		withParams := strings.Count(
			text,
			"stockPubsParamsLens.Set(fetch.PublicationParamsFromRegistry(",
		)
		require.Equal(t, states, withParams, "%s", name)
		built += states
	}
	require.Equal(t, 3, built, "expect the strain, plasmid and phenotype entries")
}
