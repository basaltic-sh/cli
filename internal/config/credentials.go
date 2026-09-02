package config

import (
	"context"

	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	basaltic "github.com/basaltic-sh/sdk-go"
	"gopkg.in/yaml.v3"
)

// Cached access tokens live in their own file, apart from config.yaml.
//
// config.yaml is something people copy between machines and check into
// dotfile repositories. A token must not travel with it.
//
// The cache exists at all only because each CLI run is a fresh process: the
// SDK caches in memory, which serves a long-running program perfectly and
// buys a command-line tool nothing. Without this file every command would
// begin with a round trip to IAM.

// EnvCredentialsFile overrides the credentials file location.
const EnvCredentialsFile = "BASALTIC_CREDENTIALS_FILE"

// refreshSkew re-exchanges this long before expiry, so a token cannot lapse
// midway through a command and surface as an authorization failure.
const refreshSkew = 5 * time.Minute

type credentialsFile struct {
	Tokens map[string]cachedToken `yaml:"tokens,omitempty"`
}

type cachedToken struct {
	// AccessKeyID is part of the key so that rotating a credential
	// invalidates the cache. Without it, replacing a key would appear to do
	// nothing for up to an hour.
	AccessKeyID string    `yaml:"access_key_id"`
	AccessToken string    `yaml:"access_token"`
	ExpiresAt   time.Time `yaml:"expires_at"`
}

// CredentialsPath returns the credentials file location.
func CredentialsPath() (string, error) {
	if p := os.Getenv(EnvCredentialsFile); p != "" {
		return p, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "basaltic", "credentials.yaml"), nil
}

// FileTokenSource is a [basaltic.TokenSource] that caches on disk between CLI
// invocations, delegating to the SDK when the cache has nothing usable.
type FileTokenSource struct {
	Profile string
	Inner   *basaltic.ClientCredentialsSource

	mu sync.Mutex
}

// Token returns a cached token, exchanging one when the cache is empty or
// close to expiry.
func (s *FileTokenSource) Token(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	path, err := CredentialsPath()
	if err != nil {
		// No cache location is not fatal. Exchanging every time is slower,
		// but the command the user asked for still runs.
		return s.Inner.Token(ctx)
	}
	file := loadCredentials(path)
	if tok, ok := file.lookup(s.Profile, s.Inner.AccessKeyID); ok {
		return tok, nil
	}

	token, err := s.Inner.Token(ctx)
	if err != nil {
		return "", err
	}
	file.store(s.Profile, s.Inner.AccessKeyID, token, s.Inner.ExpiresAt())
	if err := saveCredentials(path, file); err != nil {
		// A cache we could not write costs one exchange next time. It must
		// not fail the command.
		fmt.Fprintf(os.Stderr, "warning: could not cache the access token: %v\n", err)
	}
	return token, nil
}

// Invalidate drops both the in-memory and on-disk copies.
//
// The SDK calls this when the platform refuses a token that has not yet
// expired — a revoked session. Without clearing the FILE too, the next
// invocation would read the same dead token back off disk.
func (s *FileTokenSource) Invalidate() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Inner.Invalidate()
	s.forget()
}

// Forget clears this profile's cached token. Used by `auth logout`.
func (s *FileTokenSource) Forget() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.forget()
}

func (s *FileTokenSource) forget() {
	path, err := CredentialsPath()
	if err != nil {
		return
	}
	file := loadCredentials(path)
	if file.Tokens == nil {
		return
	}
	delete(file.Tokens, s.Profile)
	_ = saveCredentials(path, file)
}

// ForgetProfileToken clears a profile's cached token without needing a
// configured credential, which `auth logout` must be able to do: a token can
// outlive the key that minted it, so clearing the key alone would leave a
// "logged out" CLI authenticating for up to an hour.
func ForgetProfileToken(profile string) error {
	path, err := CredentialsPath()
	if err != nil {
		return err
	}
	file := loadCredentials(path)
	if _, ok := file.Tokens[profile]; !ok {
		return nil
	}
	delete(file.Tokens, profile)
	return saveCredentials(path, file)
}

// CachedTokenExpiry reports when a profile's cached token expires, for
// `auth status`.
func CachedTokenExpiry(profile, accessKeyID string) (time.Time, bool) {
	path, err := CredentialsPath()
	if err != nil {
		return time.Time{}, false
	}
	file := loadCredentials(path)
	t, ok := file.Tokens[profile]
	if !ok || t.AccessKeyID != accessKeyID || t.AccessToken == "" {
		return time.Time{}, false
	}
	return t.ExpiresAt, true
}

func (f *credentialsFile) lookup(profile, accessKeyID string) (string, bool) {
	t, ok := f.Tokens[profile]
	if !ok || t.AccessToken == "" {
		return "", false
	}
	if t.AccessKeyID != accessKeyID {
		return "", false
	}
	if time.Now().After(t.ExpiresAt.Add(-refreshSkew)) {
		return "", false
	}
	return t.AccessToken, true
}

func (f *credentialsFile) store(profile, accessKeyID, token string, expiresAt time.Time) {
	if f.Tokens == nil {
		f.Tokens = map[string]cachedToken{}
	}
	f.Tokens[profile] = cachedToken{AccessKeyID: accessKeyID, AccessToken: token, ExpiresAt: expiresAt}
}

// loadCredentials never fails: an unreadable or corrupt cache is an empty
// cache, which costs one exchange rather than breaking every command.
func loadCredentials(path string) *credentialsFile {
	f := &credentialsFile{Tokens: map[string]cachedToken{}}
	data, err := os.ReadFile(path)
	if err != nil {
		return f
	}
	if err := yaml.Unmarshal(data, f); err != nil {
		return &credentialsFile{Tokens: map[string]cachedToken{}}
	}
	if f.Tokens == nil {
		f.Tokens = map[string]cachedToken{}
	}
	return f
}

func saveCredentials(path string, f *credentialsFile) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := yaml.Marshal(f)
	if err != nil {
		return err
	}
	if err := writeFileAtomic(path, data, 0o600); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}
