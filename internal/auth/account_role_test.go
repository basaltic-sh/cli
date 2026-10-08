package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	basaltic "github.com/basaltic-sh/sdk-go"
	"github.com/basaltic-sh/sdk-go/workspace"
)

func TestAccountOperationScope(t *testing.T) {
	for _, tc := range []struct {
		service, path string
		account       bool
	}{
		{"iam", "/v1/auth/ssh-keys", false}, {"iam", "/v1/auth/linux-identity", false},
		{"iam", "/v1/assume-role", false}, {"iam", "/v1/assume-role-with-web-identity", false},
		{"iam", "/v1/service-accounts/id/ssh-keys", true}, {"iam", "/v1/roles", true},
		{"compute", "/v1/instances", true}, {"network", "/v1/vpcs", true},
		{"workspace", "/v1/account-roles", false}, {"billing", "/v1/invoices", false}, {"audit", "/v1/audit-logs", false}, {"quota", "/v1/quotas", false},
	} {
		if got := AccountScoped(tc.service, tc.path); got != tc.account {
			t.Errorf("%s %s: %v", tc.service, tc.path, got)
		}
	}
}
func TestSoleAccountRoleIsAssumedWithoutReplacingPersonalSession(t *testing.T) {
	assignments := []workspace.AccountRole{{AccountID: "account-uuid", AccountHandle: "tenant", RoleID: "role-uuid", RoleCRN: "crn:iam::tenant:role/Operator"}}
	reads, mints := 0, 0
	wrongAccount := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer personal" {
			t.Errorf("role resolution used wrong source identity")
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/account-roles":
			reads++
			_ = json.NewEncoder(w).Encode(map[string]any{"account_roles": assignments})
		case "/v1/assume-role":
			mints++
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body["role"] != "crn:iam::tenant:role/Operator" {
				t.Errorf("unexpected role: %v", body["role"])
			}
			account := "account-uuid"
			if wrongAccount {
				account = "other"
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "account-token", "role_id": "role-uuid", "account_id": account, "account_handle": "tenant", "expiration": time.Now().Add(15 * time.Minute).Format(time.RFC3339)})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	cfg, err := basaltic.NewConfig(context.Background(), basaltic.WithAccessToken("personal"), basaltic.WithAccountID("tenant"), basaltic.WithServiceEndpoint("iam", server.URL), basaltic.WithServiceEndpoint("workspace", server.URL))
	if err != nil {
		t.Fatal(err)
	}
	source := &AccountRoleSource{Base: cfg}
	for i := 0; i < 2; i++ {
		token, err := source.Token(context.Background())
		if err != nil || token != "account-token" {
			t.Fatalf("token: %q %v", token, err)
		}
	}
	if reads != 1 || mints != 1 {
		t.Fatalf("expected one invocation-local session, got %d/%d", reads, mints)
	}
	source.Invalidate()
	token, err := cfg.TokenSource.Token(context.Background())
	if err != nil || token != "personal" {
		t.Fatalf("personal session was replaced: %q %v", token, err)
	}
	assignments = append(assignments, workspace.AccountRole{AccountHandle: "tenant", RoleID: "other-role"})
	if _, err = source.Token(context.Background()); err == nil {
		t.Fatal("ambiguous account role silently chosen")
	}
	if mints != 1 {
		t.Fatal("ambiguous role minted a session")
	}
	assignments = nil
	if _, err = source.Token(context.Background()); err == nil {
		t.Fatal("missing assignment silently allowed")
	}
	assignments = []workspace.AccountRole{{AccountID: "account-uuid", AccountHandle: "tenant", RoleID: "role-uuid", RoleCRN: "crn:iam::tenant:role/Operator"}}
	wrongAccount = true
	if _, err = source.Token(context.Background()); err == nil {
		t.Fatal("wrong-account session accepted")
	}
}
