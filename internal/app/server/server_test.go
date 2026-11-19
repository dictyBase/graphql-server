package server

import (
	"flag"
	"net/http"
	"testing"
	"time"

	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/urfave/cli"
)

func authTestContext(
	t *testing.T,
	authEnabled bool,
	values map[string]string,
) *cli.Context {
	t.Helper()
	fset := flag.NewFlagSet("auth-test", flag.ContinueOnError)
	for _, f := range []cli.Flag{
		cli.BoolFlag{Name: "auth-enabled"},
		cli.StringFlag{Name: "jwks-uri"},
		cli.StringFlag{Name: "jwt-issuer"},
		cli.StringFlag{Name: "jwt-audience"},
		cli.StringFlag{Name: "auth-api-endpoint"},
		cli.StringFlag{Name: "app-id"},
		cli.StringFlag{Name: "app-secret"},
	} {
		f.Apply(fset)
	}
	if authEnabled {
		if err := fset.Set("auth-enabled", "true"); err != nil {
			t.Fatalf("unexpected error setting auth-enabled: %s", err)
		}
	}
	for name, val := range values {
		if err := fset.Set(name, val); err != nil {
			t.Fatalf("unexpected error setting %s: %s", name, err)
		}
	}
	return cli.NewContext(nil, fset, nil)
}

func TestValidateAuthFlagsDisabled(t *testing.T) {
	// no auth flags set, auth disabled -> valid
	ctx := authTestContext(t, false, nil)
	if err := validateAuthFlags(ctx, false); err != nil {
		t.Fatalf("expected no error with auth disabled, got %s", err)
	}
}

func TestValidateAuthFlagsEnabledMissing(t *testing.T) {
	// auth enabled but no jwks-uri -> error naming the flag
	ctx := authTestContext(t, true, nil)
	err := validateAuthFlags(ctx, true)
	if err == nil {
		t.Fatal("expected error when auth enabled with missing flags")
	}
	if want := "--jwks-uri is required when auth is enabled"; err.Error() != want {
		t.Fatalf("expected %q, got %q", want, err.Error())
	}
}

func TestValidateAuthFlagsEnabledComplete(t *testing.T) {
	ctx := authTestContext(t, true, map[string]string{
		"jwks-uri":          "https://logto.example.com/oidc/jwks",
		"jwt-issuer":        "https://logto.example.com/oidc",
		"jwt-audience":      "https://api.example.com",
		"auth-api-endpoint": "https://logto.example.com",
		"app-id":            "app",
		"app-secret":        "secret",
	})
	if err := validateAuthFlags(ctx, true); err != nil {
		t.Fatalf("expected no error with all flags set, got %s", err)
	}
}

// TestNewGraphQLServerWriteTimeout pins the write timeout to the ingress
// proxy-read-timeout (nginx default 60s) so slow stock/annotation searches
// return GraphQL responses instead of 502s from a truncated server write.
func TestNewGraphQLServerWriteTimeout(t *testing.T) {
	srv := newGraphQLServer(http.NotFoundHandler())
	if want := 60 * time.Second; srv.WriteTimeout != want {
		t.Fatalf("expected WriteTimeout %s, got %s", want, srv.WriteTimeout)
	}
}

func ncbiTestContext(t *testing.T, values map[string]string) *cli.Context {
	t.Helper()
	fset := flag.NewFlagSet("ncbi-test", flag.ContinueOnError)
	for _, f := range []cli.Flag{
		cli.StringFlag{Name: "ncbi-api-key"},
		cli.StringFlag{Name: "ncbi-email"},
		cli.StringFlag{Name: "ncbi-tool"},
	} {
		f.Apply(fset)
	}
	for name, val := range values {
		if err := fset.Set(name, val); err != nil {
			t.Fatalf("unexpected error setting %s: %s", name, err)
		}
	}
	return cli.NewContext(nil, fset, nil)
}

func TestAddNcbiRecordsStoresFlagValues(t *testing.T) {
	nreg := registry.NewRegistry()
	addNcbiRecords(ncbiTestContext(t, map[string]string{
		"ncbi-api-key": "key-1",
		"ncbi-email":   "ops@example.org",
		"ncbi-tool":    "dictybase-graphql",
	}), nreg)
	if got := nreg.GetRecord(registry.NcbiAPIKey); got != "key-1" {
		t.Fatalf("expected the api key, got %q", got)
	}
	if got := nreg.GetRecord(registry.NcbiEmail); got != "ops@example.org" {
		t.Fatalf("expected the email, got %q", got)
	}
	if got := nreg.GetRecord(registry.NcbiTool); got != "dictybase-graphql" {
		t.Fatalf("expected the tool, got %q", got)
	}
}

func TestAddNcbiRecordsUnconfigured(t *testing.T) {
	nreg := registry.NewRegistry()
	addNcbiRecords(ncbiTestContext(t, nil), nreg)
	for _, key := range []string{
		registry.NcbiAPIKey, registry.NcbiEmail, registry.NcbiTool,
	} {
		if got := nreg.GetRecord(key); got != "" {
			t.Fatalf("expected an empty %s record, got %q", key, got)
		}
	}
}

func TestValidateNcbiFlags(t *testing.T) {
	for _, tc := range []struct {
		name   string
		values map[string]string
		want   string
	}{
		{name: "none set", values: nil},
		{name: "all set", values: map[string]string{
			"ncbi-api-key": "key-1",
			"ncbi-email":   "ops@example.org",
			"ncbi-tool":    "dictybase-graphql",
		}},
		{
			name:   "api key only",
			values: map[string]string{"ncbi-api-key": "key-1"},
			want:   "--ncbi-email, --ncbi-tool required when an NCBI flag is set",
		},
		{
			name: "tool missing",
			values: map[string]string{
				"ncbi-api-key": "key-1",
				"ncbi-email":   "ops@example.org",
			},
			want: "--ncbi-tool required when an NCBI flag is set",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := validateNcbiFlags(ncbiTestContext(t, tc.values))
			if tc.want == "" {
				if err != nil {
					t.Fatalf("expected no error, got %s", err)
				}
				return
			}
			if err == nil || err.Error() != tc.want {
				t.Fatalf("expected %q, got %v", tc.want, err)
			}
		})
	}
}
