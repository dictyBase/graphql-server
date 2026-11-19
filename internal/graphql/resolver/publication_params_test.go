package resolver

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestPublicationFetchesReceiveParams fails when a resolver calls the
// publication fetch without the NCBI identity values from the registry.
func TestPublicationFetchesReceiveParams(t *testing.T) {
	t.Parallel()
	src, err := os.ReadFile("publication.go")
	require.NoError(t, err)
	text := string(src)
	calls := strings.Count(text, "fetch.FetchPublication(")
	require.Equal(t, 2, calls, "expect the single and the related entries")
	require.Equal(
		t,
		calls,
		strings.Count(text, "fetch.PublicationParamsFromRegistry("),
	)
}
