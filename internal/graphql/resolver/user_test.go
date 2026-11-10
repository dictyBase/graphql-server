package resolver

import (
	"context"
	"testing"

	"github.com/dictyBase/graphql-server/internal/authentication"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks"
	"github.com/dictyBase/graphql-server/internal/registry"
	"github.com/stretchr/testify/require"
)

func TestQueryUserAndUserByEmailAuthDisabled(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	reg := registry.NewRegistry()
	reg.AddAuthClient(registry.AUTH, authentication.NoopClient{})
	qrs := &QueryResolver{
		Registry: reg,
		Logger:   mocks.TestLogger(),
	}

	u1, err := qrs.User(context.Background(), "123")
	req.NoError(err)
	req.NotNil(u1)
	req.Equal("Dicty", u1.GetData().GetAttributes().GetFirstName())
	req.Equal("Stock Center", u1.GetData().GetAttributes().GetLastName())

	u2, err := qrs.UserByEmail(context.Background(), "user@example.org")
	req.NoError(err)
	req.NotNil(u2)
	req.Equal("Dicty", u2.GetData().GetAttributes().GetFirstName())
	req.Equal("Stock Center", u2.GetData().GetAttributes().GetLastName())
}
