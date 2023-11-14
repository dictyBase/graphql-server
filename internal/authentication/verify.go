package authentication

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dictyBase/graphql-server/internal/app/middleware"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var (
	contentCreatorRole  = []string{"content-writer", "content-admin"}
	contentEditorRole   = []string{"content-editor", "content-admin"}
	contentCreatorScope = "write:content"
	contentEditorScope  = "edit:content"
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

func CheckCreateContent(ctx context.Context) error {
	roles, scopes, err := checkTokenClaims(ctx, "roles", "scopes")
	if err != nil {
		return err
	}

	err = checkRole(roles, contentCreatorRole)
	if err != nil {
		return err
	}

	err = checkScope(scopes, contentCreatorScope)
	if err != nil {
		return err
	}

	return nil
}


func CheckUpdateContent(ctx context.Context) error {
	roles, scopes, err := checkTokenClaims(ctx, "roles", "scopes")
	if err != nil {
		return err
	}
	
	err = checkRole(roles, contentEditorRole)
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
	roleSlot string,
	scopeSlot string,
) (string, string, error) {
	token, err := HasToken(ctx)
	if err != nil {
		return "", "", err
	}
	claims := token.PrivateClaims()
	for _, clm := range []string{roleSlot, scopeSlot} {
		if _, ok := claims[clm]; !ok {
			return "", "", fmt.Errorf("query without claim %s not allowed", clm)
		}
	}
	roles := fmt.Sprintf("%v", claims[roleSlot])
	scopes := fmt.Sprintf("%v", claims[scopeSlot])

	return roles, scopes, nil
}

func checkRole(roles string, expectedRoles []string) error {
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
}

func checkScope(scopes string, expectedScope string) error {
	if !strings.Contains(scopes, expectedScope) {
		return errors.New("query without proper scope is not allowed")
	}
	return nil
}
