package authentication

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dictyBase/graphql-server/internal/app/middleware"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func HasToken(ctx context.Context) (jwt.Token, error) {
	token := middleware.TokenFromContext(ctx)
	if token == nil {
		return nil, errors.New("no jwt token found")
	}
	return token, nil
}

func CheckReadUser(ctx context.Context) error {
	token, err := HasToken(ctx)
	if err != nil {
		return err
	}
	claims := token.PrivateClaims()
	if _, ok := claims["roles"]; !ok {
		return errors.New(
			"query without claim roles is not allowed",
		)
	}
	roles := fmt.Sprintf("%v", claims["roles"])
	if !strings.Contains(roles, "user-user") {
		return errors.New("query without user-user roles not allowed")
	}
	return nil
}
