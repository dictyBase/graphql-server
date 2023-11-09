package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type contextKey string

// String output the details of context key
func (c contextKey) String() string {
	return "context key " + string(c)
}

var AuthContextKey = contextKey("jwtToken")

type JWTAuth struct {
	Set      jwk.Set
	Audience string
	Issuer   string
}

func NewJWTAuth(url, audience, issuer string) (*JWTAuth, error) {
	set, err := jwk.Fetch(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("error fetching JWK resource: %w", err)
	}

	return &JWTAuth{Set: set, Audience: audience, Issuer: issuer}, nil
}

func (mdw *JWTAuth) validateToken(req *http.Request) (jwt.Token, error) {
	if len(req.Header.Get("Authorization")) <= 0 {
		return nil, nil
	}
	token, err := jwt.ParseRequest(
		req,
		jwt.WithHeaderKey("Authorization"),
		jwt.WithKeySet(mdw.Set),
	)
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %s", err)
	}
	options := []jwt.ValidateOption{
		jwt.WithAudience(mdw.Audience),
		jwt.WithIssuer(mdw.Issuer),
	}
	if err := jwt.Validate(token, options...); err != nil {
		return nil, fmt.Errorf("error validating token: %s", err)
	}
	return token, nil
}

func (mdw *JWTAuth) JwtHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wrt http.ResponseWriter, req *http.Request) {
		token, err := mdw.validateToken(req)
		if err != nil {
			http.Error(wrt, err.Error(), http.StatusUnauthorized)
			return
		}
		newCtx := context.WithValue(req.Context(), AuthContextKey, token)
		next.ServeHTTP(wrt, req.WithContext(newCtx))
	})
}

func TokenFromContext(ctx context.Context) jwt.Token {
	return ctx.Value(AuthContextKey).(jwt.Token)
}
