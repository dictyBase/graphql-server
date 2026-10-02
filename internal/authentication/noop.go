package authentication

import "errors"

// ErrAuthDisabled is returned by every NoopClient method. It surfaces as a
// regular GraphQL error for resolvers that need the Logto management API,
// which is the intended behavior when the server runs with auth disabled.
var ErrAuthDisabled = errors.New("authentication is disabled")

// NoopClient is a stub LogtoClient registered in place of the real Logto
// management client when the graphql-server runs with --auth-enabled=false.
// It satisfies the LogtoClient interface so that the registry and resolvers
// need no conditional logic; calls simply fail with ErrAuthDisabled.
type NoopClient struct{}

// compile-time check that NoopClient implements LogtoClient
var _ LogtoClient = NoopClient{}

func (NoopClient) AccessToken() (*AccessTokenResp, error) {
	return nil, ErrAuthDisabled
}

func (NoopClient) CheckUserWithUserName(string) (bool, string, error) {
	return false, "", ErrAuthDisabled
}

func (NoopClient) UserWithEmail(string) (*UserResp, error) {
	return nil, ErrAuthDisabled
}

func (NoopClient) CheckUser(string) (bool, string, error) {
	return false, "", ErrAuthDisabled
}

func (NoopClient) AddCustomUserInformation(
	string,
	string,
	*APIUsersPatchCustomData,
) error {
	return ErrAuthDisabled
}

func (NoopClient) User(string) (*UserResp, error) {
	return nil, ErrAuthDisabled
}

func (NoopClient) Roles(string) ([]*RoleResp, error) {
	return nil, ErrAuthDisabled
}

func (NoopClient) Permissions(string) ([]*PermissionResp, error) {
	return nil, ErrAuthDisabled
}

func (NoopClient) CreateUser(string, *APIUsersPostReq) (string, error) {
	return "", ErrAuthDisabled
}
