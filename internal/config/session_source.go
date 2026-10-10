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

	mu         sync.Mutex
	lastAccess string
}

// ErrSessionExpired means the stored session can no longer be renewed and a
// new interactive login is required.
var ErrSessionExpired = errors.New("your session has expired — run `basaltic login`")

// Token returns the stored access token, refreshing it when it is close to
// expiry.
func (s *SessionTokenSource) Token(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	path, err := CredentialsPath()
	if err != nil {
		return "", err
	}
	unlock, err := lockCredentials(ctx, path)
	if err != nil {
		return "", err
	}
	defer unlock()
	file := loadCredentials(path)
	session, ok := file.Sessions[s.Profile]
	if !ok || session.AccessToken == "" {
		return "", ErrSessionExpired
	}
	if SessionFresh(session.ExpiresAt) {
		s.lastAccess = session.AccessToken
		return session.AccessToken, nil
	}
	if session.RefreshToken == "" || session.TokenEndpoint == "" {
		return "", ErrSessionExpired
	}
	newAccess, newRefresh, newExpiry, err := s.Refresh(ctx, s.Client, session.TokenEndpoint, session.RefreshToken)
	if err != nil {
		// Only an explicit OAuth revocation invalidates a login. Network failures
		// and cancellation leave the refresh token available for a later attempt.
		if errors.Is(err, ErrSessionExpired) {
			delete(file.Sessions, s.Profile)
			if saveErr := saveCredentials(path, file); saveErr != nil {
				return "", saveErr
			}
		}
		return "", err
	}
	file.Sessions[s.Profile] = userSession{AccessToken: newAccess, RefreshToken: newRefresh, ExpiresAt: newExpiry, TokenEndpoint: session.TokenEndpoint}
	if err := saveCredentials(path, file); err != nil {
		return "", err
	}
	s.lastAccess = newAccess
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
	path, err := CredentialsPath()
	if err != nil {
		return
	}
	_ = updateCredentials(context.Background(), path, func(file *credentialsFile) error {
		// A request carrying an older token may fail after another process has
		// renewed the session. It must not erase that replacement.
		if s.lastAccess != "" && file.Sessions[s.Profile].AccessToken == s.lastAccess {
			delete(file.Sessions, s.Profile)
		}
		return nil
	})
}
