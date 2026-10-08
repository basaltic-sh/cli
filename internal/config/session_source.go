package config

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"
)

// Refresher renews a session against its token endpoint. It is a function
// rather than an import so this package does not depend on the oauth package,
// which depends on this one for storage.
type Refresher func(ctx context.Context, client *http.Client, tokenEndpoint, refreshToken string) (accessToken, newRefreshToken string, expiresAt time.Time, err error)

// SessionTokenSource is a [basaltic.TokenSource] backed by an interactive
// login rather than a key pair.
//
// The difference from FileTokenSource is not just where the token comes from.
// A service-account source can always mint another token, because the key pair
// that produces one is sitting in the config file — so an expired cache costs
// a round trip and nothing else. This source cannot: when the refresh token is
// gone or refused, the only way forward is a person at a browser. That is why
// a failure here says "run basaltic login" instead of retrying.
type SessionTokenSource struct {
	Profile string
	Client  *http.Client
	Refresh Refresher

	mu sync.Mutex
}

// ErrSessionExpired means the stored session can no longer be renewed and a
// new interactive login is required.
var ErrSessionExpired = errors.New("your session has expired — run `basaltic login`")

// Token returns the stored access token, refreshing it when it is close to
// expiry.
func (s *SessionTokenSource) Token(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	access, refresh, endpoint, expires, ok := LookupSession(s.Profile)
	if !ok {
		return "", ErrSessionExpired
	}
	if SessionFresh(expires) {
		return access, nil
	}
	if refresh == "" || endpoint == "" {
		return "", ErrSessionExpired
	}

	newAccess, newRefresh, newExpiry, err := s.Refresh(ctx, s.Client, endpoint, refresh)
	if err != nil {
		// A refusal here is not a transport problem to retry — the session is
		// over. Clearing the stored copy stops every later command repeating
		// the same doomed refresh.
		_, _, _ = ForgetSession(s.Profile)
		return "", ErrSessionExpired
	}

	// Store BEFORE returning. The refresh token rotated, and the one we just
	// spent is already dead on the server; losing the replacement here would
	// end the session on the next command with nothing to explain why.
	if err := StoreSession(s.Profile, newAccess, newRefresh, endpoint, newExpiry); err != nil {
		return "", err
	}
	return newAccess, nil
}

// Invalidate drops the stored session.
//
// The SDK calls this when the platform refuses a token that has not expired —
// a revoked session. There is nothing to fall back to, so the next command
// asks for a login rather than silently trying again.
func (s *SessionTokenSource) Invalidate() {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, _, _ = ForgetSession(s.Profile)
}
