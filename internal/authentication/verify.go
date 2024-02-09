package authentication

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dictyBase/graphql-server/internal/app/middleware"
)

var (
	userReadRole = "user-user"
)

func ValidateContent(
	ctx context.Context,
	scopeSlot string,
	expectedScope string,
) error {
	token := middleware.TokenFromContext(ctx)
	claims := token.PrivateClaims()
	if _, ok := claims[scopeSlot]; !ok {
		return fmt.Errorf(
			"query without claim %s not allowed",
			scopeSlot,
		)
	}
	scopes := fmt.Sprintf("%v", claims[scopeSlot])
	if !strings.Contains(scopes, expectedScope) {
		return fmt.Errorf(
			"given scope %s does not matches expected scope %s",
			expectedScope,
			scopes,
		)
	}

	return nil
}

func CheckReadUser(ctx context.Context) error {
	token := middleware.TokenFromContext(ctx)
	claims := token.PrivateClaims()
	if _, ok := claims["roles"]; !ok {
		return errors.New(
			"query without claim roles is not allowed",
		)
	}
	roles := fmt.Sprintf("%v", claims["roles"])
	if !strings.Contains(roles, userReadRole) {
		return errors.New("query without user-user roles not allowed")
	}
	return nil
}
