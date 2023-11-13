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

func CheckReadUser(ctx context.Context) (bool, error) {
	token, err := HasToken(ctx)
	if err != nil {
		return false, err
	}
	claims := token.PrivateClaims()
	for _, clm := range []string{"roles", "scopes"} {
		if _, ok := claims[clm]; !ok {
			return false, fmt.Errorf(
				"query without claim %s is not allowed",
				clm,
			)
		}
	}
	roles := fmt.Sprintf("%v", claims["roles"])
	if !strings.Contains(roles, "user-user") {
		return false, errors.New("query without user-user roles not allowed")
	}
	scopes := fmt.Sprintf("%v",claims["scopes"])
	if !strings.Contains(scopes,"read:user") {
		return false, errors.New("query without read:user scope not allowed")
	}
	return true, nil
}
