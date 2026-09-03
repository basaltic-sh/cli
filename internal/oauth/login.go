// Package oauth runs the interactive login: the browser half of OAuth 2.0
// authorization code with PKCE, over a loopback redirect (RFC 8252).
//
// This is how a PERSON gets a token. A service account exchanges its key pair
// at the token endpoint and never comes here; the difference matters because a
// user token can create an organization, accept an invitation and switch org,
// and a service account token cannot.
//
// The CLI is a PUBLIC client — it ships to laptops, so any secret compiled
// into it would be readable by everyone holding the binary, and a secret
// everybody has authenticates nobody. Two things stand in for one:
//
//   - The authorization code is delivered only to a loopback address, so it
//     reaches a process on this machine or nowhere.
//   - PKCE ties redemption to a verifier that never crosses the network until
//     the moment it is spent, so a code seen in transit is not enough.
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
	"net"
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

// loginTimeout bounds the whole flow. It is generous because a person may have
// to sign in and pass a second factor, and short enough that a forgotten tab
// does not leave a listener bound indefinitely.
const loginTimeout = 5 * time.Minute

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

// Login runs the whole interactive flow and returns the tokens.
//
// notify is called with the authorization URL once it is known. It must PRINT
// the URL as well as attempting to open a browser: an SSH session has no
// browser to open, and a flow that silently waits on a page nobody can see is
// indistinguishable from a hang.
func Login(ctx context.Context, client *http.Client, m *Metadata, notify func(url string)) (*Tokens, error) {
	verifier, challenge, err := newPKCE()
	if err != nil {
		return nil, err
	}
	state, err := randomString(32)
	if err != nil {
		return nil, err
	}

	// Bind the listener BEFORE building the URL: the redirect must name the
	// port the OS actually gave us, and asking for one we have not bound is
	// how a race with another process becomes an unexplained failure.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("listen on loopback: %w", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	redirectURI := fmt.Sprintf("http://127.0.0.1:%d/callback", port)

	authURL, err := buildAuthorizeURL(m.AuthorizationEndpoint, redirectURI, challenge, state)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, loginTimeout)
	defer cancel()

	result := make(chan callbackResult, 1)
	srv := &http.Server{Handler: callbackHandler(state, result)}
	go func() { _ = srv.Serve(listener) }()
	defer func() {
		shutdownCtx, done := context.WithTimeout(context.Background(), 2*time.Second)
		defer done()
		_ = srv.Shutdown(shutdownCtx)
	}()

	notify(authURL)
	_ = openBrowser(authURL)

	select {
	case <-ctx.Done():
		return nil, errors.New("timed out waiting for the browser. Run the command again")
	case got := <-result:
		if got.err != nil {
			return nil, got.err
		}
		return exchange(ctx, client, m.TokenEndpoint, got.code, verifier, redirectURI)
	}
}

type callbackResult struct {
	code string
	err  error
}

// callbackHandler serves the one request the browser makes on return.
func callbackHandler(wantState string, out chan<- callbackResult) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		// The state comparison is what stops a code from somewhere else being
		// fed to this listener. It is the client's half of the flow and no
		// server check substitutes for it.
		if q.Get("state") != wantState {
			page(w, http.StatusBadRequest, "Sign-in failed",
				"The response did not match this sign-in attempt. Run the command again.")
			out <- callbackResult{err: errors.New("the redirect did not match this sign-in attempt (state mismatch)")}
			return
		}
		if e := q.Get("error"); e != "" {
			desc := q.Get("error_description")
			if desc == "" {
				desc = e
			}
			page(w, http.StatusOK, "Sign-in cancelled", desc)
			out <- callbackResult{err: fmt.Errorf("authorization was refused: %s", desc)}
			return
		}
		code := q.Get("code")
		if code == "" {
			page(w, http.StatusBadRequest, "Sign-in failed", "No authorization code was returned.")
			out <- callbackResult{err: errors.New("no authorization code in the redirect")}
			return
		}
		page(w, http.StatusOK, "Signed in", "You can close this tab and return to your terminal.")
		out <- callbackResult{code: code}
	})
	// Anything else on this listener is not part of the flow. Answer plainly
	// rather than 404-ing into a browser's error page.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page(w, http.StatusNotFound, "Nothing here", "This is the Basaltic CLI waiting for a sign-in redirect.")
	})
	return mux
}

func page(w http.ResponseWriter, status int, title, body string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// The URL of this page carries the authorization code, so it must not be
	// stored by the browser any longer than it takes to render.
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	fmt.Fprintf(w, `<!doctype html><meta charset="utf-8"><title>%s</title>`+
		`<body style="font:14px system-ui;margin:4rem auto;max-width:30rem;color:#111">`+
		`<h1 style="font-size:1.1rem">%s</h1><p>%s</p></body>`, title, title, body)
}

// exchange redeems the code for tokens.
func exchange(ctx context.Context, client *http.Client, tokenEndpoint, code, verifier, redirectURI string) (*Tokens, error) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"code_verifier": {verifier},
		"redirect_uri":  {redirectURI},
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
				return nil, errors.New("the sign-in could not be completed; it may have expired or already been used. Run the command again")
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

func buildAuthorizeURL(endpoint, redirectURI, challenge, state string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("authorization endpoint is not a valid URL: %w", err)
	}
	q := u.Query()
	q.Set("client_id", ClientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("code_challenge", challenge)
	q.Set("code_challenge_method", "S256")
	q.Set("state", state)
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

// openBrowser makes a best effort. Its error is deliberately ignored by the
// caller: the URL has already been printed, so a machine with no browser is
// inconvenienced rather than blocked.
func openBrowser(u string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", u).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", u).Start()
	default:
		return exec.Command("xdg-open", u).Start()
	}
}
