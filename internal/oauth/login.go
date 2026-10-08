// Package oauth runs the interactive login: the browser half of OAuth 2.0
// authorization code with PKCE.
//
// This is how a PERSON gets a token. A service account exchanges its key pair
// at the token endpoint and never comes here; the difference matters because a
// user token can create an organization, accept an invitation and switch org,
// and a service account token cannot.
//
// The code comes back OUT OF BAND: the console shows it, the user pastes it
// into the terminal. There is no loopback listener and no redirect, which is
// what makes this work on a machine you only reach over SSH — the browser is
// on the laptop, the CLI is on the box, and nothing has to bridge them.
//
// The CLI is a PUBLIC client — it ships to laptops, so any secret compiled
// into it would be readable by everyone holding the binary, and a secret
// everybody has authenticates nobody. PKCE stands in for one: the verifier
// ties redemption to the process that started the flow and does not cross the
// network until the moment it is spent, so a code read over someone's shoulder
// is not enough to spend it.
package oauth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// ClientID is the registered public client for this CLI. It matches
// OAUTH_CLI_CLIENT_ID on the platform side; there is no client registration
// API, because the set of first-party clients is small, fixed and
// security-relevant.
const ClientID = "basaltic-cli"

// RedirectURIOOB is the out-of-band pseudo-redirect (the long-standing name
// for "display the code instead of delivering it"). The server accepts nothing
// else for this client, and sends it in the authorize request so the console
// knows to show the code rather than navigate.
const RedirectURIOOB = "urn:ietf:wg:oauth:2.0:oob"

// Metadata is the subset of the RFC 8414 authorization-server document this
// flow needs.
type Metadata struct {
	Issuer                        string   `json:"issuer"`
	AuthorizationEndpoint         string   `json:"authorization_endpoint"`
	TokenEndpoint                 string   `json:"token_endpoint"`
	GrantTypesSupported           []string `json:"grant_types_supported"`
	CodeChallengeMethodsSupported []string `json:"code_challenge_methods_supported"`
}

// Tokens is what a completed login yields.
type Tokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// DiscoverMetadata reads the authorization-server document from an IAM base
// URL.
//
// The authorization endpoint is DISCOVERED rather than constructed. The
// console lives at a different host in every deployment, and RFC 8414 exists
// precisely so a client does not have to be told. It also means a region that
// has not enabled interactive login says so — the field is absent — instead of
// sending someone to a URL that 404s.
func DiscoverMetadata(ctx context.Context, client *http.Client, iamBaseURL string) (*Metadata, error) {
	u := strings.TrimRight(iamBaseURL, "/") + "/.well-known/oauth-authorization-server"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("read authorization server metadata: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("authorization server metadata: HTTP %d", resp.StatusCode)
	}
	var m Metadata
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&m); err != nil {
		return nil, fmt.Errorf("parse authorization server metadata: %w", err)
	}
	if m.TokenEndpoint == "" {
		return nil, errors.New("authorization server metadata has no token_endpoint")
	}
	return &m, nil
}

// SupportsInteractiveLogin reports whether this deployment advertises the code
// grant. A deployment that does not is not broken — it simply has no console
// configured to approve at, and the user should be told that rather than shown
// a failed redirect.
func (m *Metadata) SupportsInteractiveLogin() bool {
	if m.AuthorizationEndpoint == "" {
		return false
	}
	for _, g := range m.GrantTypesSupported {
		if g == "authorization_code" {
			return true
		}
	}
	return false
}

// Request is one sign-in attempt: a URL to approve at, and the verifier that
// will redeem whatever code comes back from it.
//
// The verifier is unexported and never printed. It is the only thing that
// distinguishes this process from anyone else who has seen the code, so a
// Request must not outlive the terminal session that made it.
type Request struct {
	// AuthorizationURL is where the person approves. Give it to them whether
	// or not a browser opens: an SSH session has none, and a flow that waits
	// silently on a page nobody can see is indistinguishable from a hang.
	AuthorizationURL string

	verifier string
}

// NewRequest starts a sign-in.
func NewRequest(m *Metadata) (*Request, error) {
	verifier, challenge, err := newPKCE()
	if err != nil {
		return nil, err
	}
	u, err := buildAuthorizeURL(m.AuthorizationEndpoint, challenge)
	if err != nil {
		return nil, err
	}
	return &Request{AuthorizationURL: u, verifier: verifier}, nil
}

// Redeem exchanges a pasted code for tokens.
//
// Surrounding whitespace is trimmed because the code arrives through a
// clipboard and a terminal, both of which add a newline or a stray space
// without being asked. Nothing else is normalised: the code is base64, so
// changing its case or dropping a character would turn a copy-paste slip into
// "invalid_grant", which the server deliberately does not explain.
func (r *Request) Redeem(ctx context.Context, client *http.Client, tokenEndpoint, code string) (*Tokens, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, errors.New("no code was entered")
	}
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"code_verifier": {r.verifier},
		"redirect_uri":  {RedirectURIOOB},
		"client_id":     {ClientID},
	}
	return post(ctx, client, tokenEndpoint, form)
}

// Refresh renews a session. The refresh token ROTATES on every use, so the
// caller must store the one that comes back.
func Refresh(ctx context.Context, client *http.Client, tokenEndpoint, refreshToken string) (*Tokens, error) {
	form := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
		"client_id":     {ClientID},
	}
	return post(ctx, client, tokenEndpoint, form)
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type oauthError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func post(ctx context.Context, client *http.Client, endpoint string, form url.Values) (*Tokens, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token endpoint: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))

	if resp.StatusCode != http.StatusOK {
		var e oauthError
		if json.Unmarshal(body, &e) == nil && e.Error != "" {
			// invalid_grant covers an expired, spent or mismatched code, and
			// the server deliberately does not say which — telling them apart
			// would tell a holder of a stolen code which half they got wrong.
			if e.Error == "invalid_grant" {
				return nil, errors.New("that code was not accepted; it may have expired, been used already, or been mistyped. Run the command again")
			}
			if e.ErrorDescription != "" {
				return nil, fmt.Errorf("%s: %s", e.Error, e.ErrorDescription)
			}
			return nil, errors.New(e.Error)
		}
		return nil, fmt.Errorf("token endpoint: HTTP %d", resp.StatusCode)
	}

	var tr tokenResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return nil, fmt.Errorf("parse token response: %w", err)
	}
	if tr.AccessToken == "" {
		return nil, errors.New("the token endpoint returned no access token")
	}
	return &Tokens{
		AccessToken:  tr.AccessToken,
		RefreshToken: tr.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second),
	}, nil
}

// buildAuthorizeURL points the browser at the consent page.
//
// No state parameter. State exists to bind a REDIRECT to the request that
// caused it, and there is no redirect here — the code is read off a page by
// the person who asked for it. PKCE is what binds the code to this process,
// and it does so whether or not anybody echoes a nonce.
func buildAuthorizeURL(endpoint, challenge string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("authorization endpoint is not a valid URL: %w", err)
	}
	q := u.Query()
	q.Set("client_id", ClientID)
	q.Set("redirect_uri", RedirectURIOOB)
	q.Set("code_challenge", challenge)
	q.Set("code_challenge_method", "S256")
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// newPKCE returns a verifier and its S256 challenge (RFC 7636).
func newPKCE() (verifier, challenge string, err error) {
	verifier, err = randomString(64)
	if err != nil {
		return "", "", err
	}
	sum := sha256.Sum256([]byte(verifier))
	return verifier, base64.RawURLEncoding.EncodeToString(sum[:]), nil
}

// randomString returns n bytes of randomness, base64url encoded without
// padding — the character set RFC 7636 allows for a verifier.
func randomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("read randomness: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// OpenBrowser makes a best effort. Its error is meant to be ignored: the URL
// has already been printed, so a machine with no browser — which is most of
// the machines this flow now has to work on — is inconvenienced rather than
// blocked.
func OpenBrowser(u string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", u).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", u).Start()
	default:
		return exec.Command("xdg-open", u).Start()
	}
}
