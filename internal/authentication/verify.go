package authentication

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dictyBase/graphql-server/internal/app/middleware"
)

var (
	// contentCreatorRole  = []string{"content-writer", "content-admin"}
	// contentEditorRole   = []string{"content-editor", "content-admin"}
	// contentDeleteRole   = []string{"content-admin", "content-remover"}
	contentCreatorScope = "write:content"
	contentEditorScope  = "edit:content"
	contentDeleteScope  = "delete:content"
	userReadRole        = "user-user"
)

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

func CheckCreateContent(ctx context.Context) error {
	scopes, err := checkTokenClaims(ctx, "scopes")
	if err != nil {
		return err
	}
	err = checkScope(scopes, contentCreatorScope)
	if err != nil {
		return err
	}

	return nil
}

func CheckDeleteContent(ctx context.Context) error {
	scopes, err := checkTokenClaims(ctx, "scopes")
	if err != nil {
		return err
	}
	err = checkScope(scopes, contentDeleteScope)
	if err != nil {
		return err
	}

	return nil
}

func CheckUpdateContent(ctx context.Context) error {
	scopes, err := checkTokenClaims(ctx, "scopes")
	if err != nil {
		return err
	}
	err = checkScope(scopes, contentEditorScope)
	if err != nil {
		return err
	}

	return nil
}

func checkTokenClaims(
	ctx context.Context,
	scopeSlot string,
) (string, error) {
	token := middleware.TokenFromContext(ctx)
	claims := token.PrivateClaims()
	keys := make([]string, 0)
	for k := range claims {
		keys = append(keys, k)
	}
	if _, ok := claims[scopeSlot]; !ok {
		return "", fmt.Errorf(
			"query without claim %s not allowed => given claims %s",
			scopeSlot, keys,
		)
	}
	scopes := fmt.Sprintf("%v", claims[scopeSlot])
	return scopes, nil
}

/* func checkRole(roles string, expectedRoles []string) error {
	rolesOk := false
	for _, rls := range expectedRoles {
		if strings.Contains(roles, rls) {
			rolesOk = true
			break
		}
	}
	if !rolesOk {
		return errors.New("query without proper role is not allowed")
	}
	return nil
} */

func checkScope(scopes string, expectedScope string) error {
	if !strings.Contains(scopes, expectedScope) {
		return errors.New("query without proper scope is not allowed")
	}
	return nil
}
