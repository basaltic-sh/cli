package auth

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	basaltic "github.com/basaltic-sh/sdk-go"
	"github.com/basaltic-sh/sdk-go/iam"
	"github.com/basaltic-sh/sdk-go/workspace"
)

// AccountScoped separates organization/personal operations from account IAM.
// An interactive user has exactly one effective role in a selected account.
// Explicit bearer tokens and service-account credentials keep their identity.
func AccountScoped(service, path string) bool {
	switch service {
	case "workspace", "audit", "billing", "quota", "catalog":
		return false
	case "iam":
		return !strings.HasPrefix(path, "/v1/auth/") && !strings.HasPrefix(path, "/v1/sts/") && !strings.HasPrefix(path, "/v1/assume-role") && !strings.HasPrefix(path, "/oauth/")
	default:
		return true
	}
}

// AccountRoleSource is invocation-local. The backend revalidates assignments
// and restrictions on every use of the resulting session. Invalidating that
// session must not erase the person's separate interactive login.
type AccountRoleSource struct {
	Base    *basaltic.Config
	mu      sync.Mutex
	token   string
	expires time.Time
}

func (s *AccountRoleSource) Token(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.token != "" && time.Until(s.expires) > time.Minute {
		return s.token, nil
	}
	if s.Base == nil || s.Base.AccountID == "" {
		return "", fmt.Errorf("select an account with --account-id before using account resources")
	}
	page, err := workspace.New(s.Base).ListAccountRoles(ctx)
	if err != nil {
		return "", fmt.Errorf("resolve your account role: %w", err)
	}
	roles := map[string]workspace.AccountRole{}
	for _, role := range page.Items {
		if role.AccountHandle == s.Base.AccountID || role.AccountID == s.Base.AccountID {
			roles[role.RoleID] = role
		}
	}
	if len(roles) != 1 {
		return "", fmt.Errorf("account %q requires exactly one assigned role; found %d", s.Base.AccountID, len(roles))
	}
	var role workspace.AccountRole
	for _, r := range roles {
		role = r
	}
	if role.RoleID == "" || role.RoleCRN == "" {
		return "", fmt.Errorf("the account role response is incomplete")
	}
	duration := 900
	session, err := iam.New(s.Base).AssumeRole(ctx, &iam.AssumeRoleRequest{Role: role.RoleCRN, DurationSeconds: &duration})
	if err != nil {
		return "", fmt.Errorf("assume your account role: %w", err)
	}
	if session.AccessToken == "" || !session.Expiration.After(time.Now()) || session.RoleID != role.RoleID || session.AccountID != role.AccountID || session.AccountHandle != role.AccountHandle {
		return "", fmt.Errorf("the returned role session does not match the selected account")
	}
	s.token, s.expires = session.AccessToken, session.Expiration
	return s.token, nil
}
func (s *AccountRoleSource) Invalidate() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = ""
	s.expires = time.Time{}
}
