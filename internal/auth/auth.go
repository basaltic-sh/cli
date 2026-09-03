// Package auth turns a profile plus flags and environment into the SDK
// configuration every command is built from.
package auth

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/basaltic-sh/cli/internal/config"
	"github.com/basaltic-sh/cli/internal/oauth"
	basaltic "github.com/basaltic-sh/sdk-go"
)

// Options are the settings a command can override from the command line.
type Options struct {
	Profile   string
	APIKey    string // "ACCESS_KEY_ID:SECRET"
	Region    string
	AccountID string
	Insecure  bool
}

// Resolved is everything a command needs, plus what `auth status` reports.
type Resolved struct {
	Config      *basaltic.Config
	Profile     string
	Region      string
	AccountID   string
	AccessKeyID string
	TokenSource *config.FileTokenSource

	// UserSession is set when the profile authenticates as a PERSON — an
	// interactive login rather than a service-account key pair. The two are
	// mutually exclusive, and which one is in play changes what the caller can
	// do: only a user can create an organization, accept an invitation or
	// switch org.
	UserSession bool
	SessionTTL  time.Time
}

// Resolve builds the SDK configuration.
//
// Precedence, highest first: an explicit flag, the environment, the profile.
// The SDK reads the same environment variables, but the CLI resolves them
// itself so that `auth status` can say where each value came from.
func Resolve(opts Options) (*Resolved, error) {
	file, err := config.Load()
	if err != nil {
		return nil, err
	}
	name, profile := file.Resolve(opts.Profile)

	apiKey := firstNonEmpty(opts.APIKey, os.Getenv("BASALTIC_API_KEY"), profile.APIKey)
	region := firstNonEmpty(opts.Region, os.Getenv(basaltic.EnvRegion), profile.Region)
	accountID := firstNonEmpty(opts.AccountID, os.Getenv(basaltic.EnvAccountID), profile.AccountID)

	sdkOpts := []basaltic.Option{
		basaltic.WithRegion(region),
		basaltic.WithUserAgent("basaltic-cli/" + Version),
	}
	if accountID != "" {
		sdkOpts = append(sdkOpts, basaltic.WithAccountID(accountID))
	}
	if profile.Domain != "" {
		sdkOpts = append(sdkOpts, basaltic.WithDomain(profile.Domain))
	}
	for service, endpoint := range profile.Endpoints {
		sdkOpts = append(sdkOpts, basaltic.WithServiceEndpoint(service, endpoint))
	}
	if profile.Insecure || opts.Insecure {
		sdkOpts = append(sdkOpts, basaltic.WithHTTPClient(insecureClient()))
	}

	res := &Resolved{Profile: name, Region: region, AccountID: accountID}

	// A token supplied directly wins: something upstream already holds one,
	// and the CLI has no key pair to exchange.
	if token := os.Getenv(basaltic.EnvAccessToken); token != "" {
		sdkOpts = append(sdkOpts, basaltic.WithAccessToken(token))
		cfg, err := basaltic.NewConfig(context.Background(), sdkOpts...)
		if err != nil {
			return nil, err
		}
		res.Config = cfg
		return res, nil
	}

	// An interactive login is used only when no key pair is configured. A
	// profile with both is a profile someone has deliberately pointed at a
	// service account, and silently preferring their personal session would
	// run commands as the wrong principal.
	if strings.TrimSpace(apiKey) == "" {
		if _, _, _, expires, ok := config.LookupSession(name); ok {
			ts := &config.SessionTokenSource{
				Profile: name,
				Client:  httpClient(profile.Insecure || opts.Insecure),
				Refresh: refreshSession,
			}
			cfg, err := basaltic.NewConfig(context.Background(), append(sdkOpts, basaltic.WithTokenSource(ts))...)
			if err != nil {
				return nil, err
			}
			res.Config = cfg
			res.UserSession = true
			res.SessionTTL = expires
			return res, nil
		}
		return nil, fmt.Errorf(
			"no credentials for profile %q.\n"+
				"Run `basaltic login` to sign in as yourself, set BASALTIC_API_KEY, or pass --api-key ACCESS_KEY_ID:SECRET", name)
	}
	keyID, secret, err := config.Profile{APIKey: apiKey}.Credentials()
	if err != nil {
		return nil, fmt.Errorf("profile %q: %w", name, err)
	}
	res.AccessKeyID = keyID

	// Build the SDK's own source first, then wrap it in the disk cache. The
	// SDK owns the exchange; the CLI owns only the fact that each invocation
	// is a fresh process.
	base, err := basaltic.NewConfig(context.Background(), append(sdkOpts, basaltic.WithClientCredentials(keyID, secret))...)
	if err != nil {
		return nil, err
	}
	inner, ok := base.TokenSource.(*basaltic.ClientCredentialsSource)
	if !ok {
		res.Config = base
		return res, nil
	}
	ts := &config.FileTokenSource{Profile: name, Inner: inner}
	base.TokenSource = ts
	res.TokenSource = ts
	res.Config = base
	return res, nil
}

// insecureClient disables TLS verification, for development rigs with
// self-signed certificates. It applies to every host the profile talks to.
func insecureClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// refreshSession adapts the oauth package's Refresh to the config package's
// Refresher, which is a function type precisely so config does not import
// oauth — oauth already imports config for storage.
func refreshSession(ctx context.Context, client *http.Client, tokenEndpoint, refreshToken string) (string, string, time.Time, error) {
	toks, err := oauth.Refresh(ctx, client, tokenEndpoint, refreshToken)
	if err != nil {
		return "", "", time.Time{}, err
	}
	return toks.AccessToken, toks.RefreshToken, toks.ExpiresAt, nil
}

// httpClient returns the client the login and refresh calls use. It mirrors
// the SDK's TLS handling so a development rig with a self-signed certificate
// behaves the same on the auth path as everywhere else.
func httpClient(insecure bool) *http.Client {
	if insecure {
		return insecureClient()
	}
	return &http.Client{Timeout: 30 * time.Second}
}
