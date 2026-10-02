package authentication

import (
	"errors"
	"testing"
)

// TestNoopClientImplementsLogtoClient ensures the no-op client can stand in
// for the real Logto client in the registry.
func TestNoopClientImplementsLogtoClient(t *testing.T) {
	var _ LogtoClient = NoopClient{}
}

func TestNoopClientReturnsAuthDisabledError(t *testing.T) {
	client := NoopClient{}

	t.Run("AccessToken", func(t *testing.T) {
		token, err := client.AccessToken()
		if !errors.Is(err, ErrAuthDisabled) {
			t.Fatalf("expected ErrAuthDisabled, got %v", err)
		}
		if token != nil {
			t.Fatalf("expected nil token, got %+v", token)
		}
	})

	t.Run("User", func(t *testing.T) {
		resp, err := client.User("123")
		if !errors.Is(err, ErrAuthDisabled) {
			t.Fatalf("expected ErrAuthDisabled, got %v", err)
		}
		if resp != nil {
			t.Fatalf("expected nil response, got %+v", resp)
		}
	})

	t.Run("UserWithEmail", func(t *testing.T) {
		resp, err := client.UserWithEmail("foo@bar.com")
		if !errors.Is(err, ErrAuthDisabled) {
			t.Fatalf("expected ErrAuthDisabled, got %v", err)
		}
		if resp != nil {
			t.Fatalf("expected nil response, got %+v", resp)
		}
	})

	t.Run("CheckUserWithUserName", func(t *testing.T) {
		ok, id, err := client.CheckUserWithUserName("foo")
		if !errors.Is(err, ErrAuthDisabled) {
			t.Fatalf("expected ErrAuthDisabled, got %v", err)
		}
		if ok || id != "" {
			t.Fatalf("expected zero values, got ok=%v id=%q", ok, id)
		}
	})

	t.Run("CheckUser", func(t *testing.T) {
		ok, id, err := client.CheckUser("foo@bar.com")
		if !errors.Is(err, ErrAuthDisabled) {
			t.Fatalf("expected ErrAuthDisabled, got %v", err)
		}
		if ok || id != "" {
			t.Fatalf("expected zero values, got ok=%v id=%q", ok, id)
		}
	})

	t.Run("AddCustomUserInformation", func(t *testing.T) {
		err := client.AddCustomUserInformation("123", "456", nil)
		if !errors.Is(err, ErrAuthDisabled) {
			t.Fatalf("expected ErrAuthDisabled, got %v", err)
		}
	})

	t.Run("Roles", func(t *testing.T) {
		roles, err := client.Roles("123")
		if !errors.Is(err, ErrAuthDisabled) {
			t.Fatalf("expected ErrAuthDisabled, got %v", err)
		}
		if roles != nil {
			t.Fatalf("expected nil roles, got %+v", roles)
		}
	})

	t.Run("Permissions", func(t *testing.T) {
		perms, err := client.Permissions("123")
		if !errors.Is(err, ErrAuthDisabled) {
			t.Fatalf("expected ErrAuthDisabled, got %v", err)
		}
		if perms != nil {
			t.Fatalf("expected nil permissions, got %+v", perms)
		}
	})

	t.Run("CreateUser", func(t *testing.T) {
		id, err := client.CreateUser("123", nil)
		if !errors.Is(err, ErrAuthDisabled) {
			t.Fatalf("expected ErrAuthDisabled, got %v", err)
		}
		if id != "" {
			t.Fatalf("expected empty id, got %q", id)
		}
	})
}
