package server

import (
	"flag"
	"testing"

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
