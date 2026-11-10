package resolverutils

import (
	"context"
	"errors"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/dictyBase/fp-go-loom/loomcheck"
	"github.com/dictyBase/fp-go-loom/pipelinecheck"
	"github.com/dictyBase/graphql-server/internal/authentication"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks"
	"github.com/dictyBase/graphql-server/internal/graphql/mocks/clients"
	"github.com/stretchr/testify/require"
)

func TestUserLookupPipelineCheck(t *testing.T) {
	t.Parallel()
	cfg := pipelinecheck.Config{
		Roots: []string{"user_lookup.go"},
		IsEntrypoint: func(name string) bool {
			return name == "UserByEmail" ||
				name == "StockDepositor" ||
				name == "UserByID" ||
				name == "QueryUserByEmail"
		},
		RequirePointFreeSeed:      true,
		RequirePointFreeBranching: true,
		RequireTrivialDelegate:    true,
	}
	pipelinecheck.Require(t, cfg)
	pipelinecheck.RequireNoNonRawTryCatchCallback(t, cfg)
	pipelinecheck.RequireNoIfErrInTryCatch(t, cfg)
	pipelinecheck.RequireBareIf(t, cfg)
	pipelinecheck.RequireTrivialDelegate(t, cfg)
	pipelinecheck.RequireNoImperativeLoop(t, cfg)
	pipelinecheck.RequireNoImperativeBranch(t, cfg)
	pipelinecheck.RequireNoDirectApply(t, cfg)
	pipelinecheck.RequireNoPopulatedStruct(t, cfg)
	pipelinecheck.RequireNoLocalPipeClosure(t, cfg)
	pipelinecheck.RequireNoRedundantFold(t, cfg)
	pipelinecheck.RequireSafePrintf(t, cfg)
	pipelinecheck.RequireNoHandRolledFileJoinInArrayMap(t, cfg)
}

func TestUserLookupLoomCheck(t *testing.T) {
	t.Parallel()
	loomcheck.Require(t, loomcheck.Config{
		Roots:               []string{"user_lookup.go"},
		StrictAllPredicates: true,
	})
}

func TestUserByEmailAuthDisabledReturnsFakeUser(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	res := F.Pipe5(
		UserLookupState{},
		LookupCtxLens.Set(context.Background()),
		LookupQueryLens.Set("admin@dictybase.org"),
		LookupClientLens.Set(authentication.NoopClient{}),
		LookupLoggerLens.Set(mocks.TestLogger()),
		UserByEmail,
	)
	req.NoError(res.F1)
	req.NotNil(res.F2)
	req.Equal(int64(2375), res.F2.GetData().GetId())
	req.Equal("Dicty", res.F2.GetData().GetAttributes().GetFirstName())
	req.Equal("Stock Center", res.F2.GetData().GetAttributes().GetLastName())
	req.Equal(
		"dictystocks@northwestern.edu",
		res.F2.GetData().GetAttributes().GetEmail(),
	)
}

func TestStockDepositorEmptyReturnsFakeUser(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	res := F.Pipe5(
		UserLookupState{},
		LookupCtxLens.Set(context.Background()),
		LookupQueryLens.Set(""),
		LookupClientLens.Set(authentication.NoopClient{}),
		LookupLoggerLens.Set(mocks.TestLogger()),
		StockDepositor,
	)
	req.NoError(res.F1)
	req.NotNil(res.F2)
	req.Equal("Dicty", res.F2.GetData().GetAttributes().GetFirstName())
	req.Equal("Stock Center", res.F2.GetData().GetAttributes().GetLastName())
}

func TestStockDepositorAuthDisabledRedactsEmailAndDoesNotReturnFakeUser(
	t *testing.T,
) {
	t.Parallel()
	tests := []struct {
		name      string
		depositor string
		firstName string
		lastName  string
	}{
		{
			name:      "single token email local part",
			depositor: "gadi@bcm.edu",
			firstName: "Gadi",
			lastName:  "",
		},
		{
			name:      "dotted email local part",
			depositor: "john.doe@northwestern.edu",
			firstName: "John",
			lastName:  "Doe",
		},
		{
			name:      "multi-part email local part",
			depositor: "mary_jane-watson@school.edu",
			firstName: "Mary",
			lastName:  "Jane Watson",
		},
		{
			name:      "local part with digits",
			depositor: "kbania99@seinfeld.org",
			firstName: "Kbania",
			lastName:  "",
		},
		{
			name:      "only symbols before domain",
			depositor: "123@example.org",
			firstName: "External",
			lastName:  "Depositor",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := require.New(t)
			res := F.Pipe5(
				UserLookupState{},
				LookupCtxLens.Set(context.Background()),
				LookupQueryLens.Set(tc.depositor),
				LookupClientLens.Set(authentication.NoopClient{}),
				LookupLoggerLens.Set(mocks.TestLogger()),
				StockDepositor,
			)
			req.NoError(res.F1)
			req.NotNil(res.F2)
			attrs := res.F2.GetData().GetAttributes()
			req.Equal(tc.firstName, attrs.GetFirstName())
			req.Equal(tc.lastName, attrs.GetLastName())
			req.Empty(attrs.GetEmail())
			req.NotContains(attrs.GetFirstName(), "@")
			req.NotContains(attrs.GetLastName(), "@")
			req.NotEqual("Dicty", attrs.GetFirstName())
			req.NotEqual("Stock Center", attrs.GetLastName())
		})
	}
}

func TestStockDepositorAuthEnabledSuccess(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	mockAuth := new(clients.LogtoClient)
	mockAuth.On("UserWithEmail", "kenny@bania.com").Return(
		&authentication.UserResp{
			ID:           "user42",
			Username:     "Kenny",
			Name:         "Bania",
			PrimaryEmail: "kenny@bania.com",
		},
		nil,
	)

	res := F.Pipe5(
		UserLookupState{},
		LookupCtxLens.Set(context.Background()),
		LookupQueryLens.Set("kenny@bania.com"),
		LookupClientLens.Set(mockAuth),
		LookupLoggerLens.Set(mocks.TestLogger()),
		StockDepositor,
	)
	req.NoError(res.F1)
	req.NotNil(res.F2)
	req.Equal(int64(42), res.F2.GetData().GetId())
	req.Equal("Kenny", res.F2.GetData().GetAttributes().GetFirstName())
	req.Equal("Bania", res.F2.GetData().GetAttributes().GetLastName())
	req.Equal("kenny@bania.com", res.F2.GetData().GetAttributes().GetEmail())
}

func TestStockDepositorNonAuthErrorPropagates(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	mockAuth := new(clients.LogtoClient)
	mockAuth.On("UserWithEmail", "missing@bania.com").Return(
		nil,
		errors.New("user not found"),
	)

	res := F.Pipe5(
		UserLookupState{},
		LookupCtxLens.Set(context.Background()),
		LookupQueryLens.Set("missing@bania.com"),
		LookupClientLens.Set(mockAuth),
		LookupLoggerLens.Set(mocks.TestLogger()),
		StockDepositor,
	)
	req.Error(res.F1)
	req.Nil(res.F2)
	req.ErrorContains(res.F1, "unable to retreieve user")
}

func TestUserByIDAuthDisabledReturnsFakeUser(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	res := F.Pipe5(
		UserLookupState{},
		LookupCtxLens.Set(context.Background()),
		LookupQueryLens.Set("123"),
		LookupClientLens.Set(authentication.NoopClient{}),
		LookupLoggerLens.Set(mocks.TestLogger()),
		UserByID,
	)
	req.NoError(res.F1)
	req.NotNil(res.F2)
	req.Equal("Dicty", res.F2.GetData().GetAttributes().GetFirstName())
	req.Equal("Stock Center", res.F2.GetData().GetAttributes().GetLastName())
}
