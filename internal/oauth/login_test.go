package oauth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// The challenge must be exactly what the server recomputes from the verifier:
// base64url of the SHA-256, with no padding. A padded or standard-alphabet
// encoding fails at redemption with invalid_grant, which the server
// deliberately does not explain — so getting this wrong would surface as an
// unexplainable login failure.
func TestPKCEChallengeMatchesServerComputation(t *testing.T) {
	verifier, challenge, err := newPKCE()
	if err != nil {
		t.Fatalf("newPKCE: %v", err)
	}

	sum := sha256.Sum256([]byte(verifier))
	want := base64.RawURLEncoding.EncodeToString(sum[:])
	if challenge != want {
		t.Errorf("challenge = %q, want %q", challenge, want)
	}
	if strings.ContainsAny(challenge, "+/=") {
		t.Errorf("challenge %q is not base64url without padding", challenge)
	}
	// RFC 7636 requires a verifier of 43..128 characters.
	if len(verifier) < 43 || len(verifier) > 128 {
		t.Errorf("verifier length %d is outside the RFC 7636 range", len(verifier))
	}
}

func TestPKCEIsFreshEachTime(t *testing.T) {
	a, _, err := newPKCE()
	if err != nil {
		t.Fatal(err)
	}
	b, _, err := newPKCE()
	if err != nil {
		t.Fatal(err)
	}
	if a == b {
		t.Fatal("two logins produced the same verifier; the randomness is not being read")
	}
}

// The authorize URL is the whole request: there is no redirect to carry
// anything back, so what is not in this URL is not in the flow.
func TestBuildAuthorizeURL(t *testing.T) {
	raw, err := buildAuthorizeURL("https://console.example/cli/authorize", "the-challenge")
	if err != nil {
		t.Fatalf("buildAuthorizeURL: %v", err)
	}
	for _, want := range []string{
		"client_id=" + ClientID,
		"code_challenge=the-challenge",
		"code_challenge_method=S256",
		"redirect_uri=urn%3Aietf%3Awg%3Aoauth%3A2.0%3Aoob",
	} {
		if !strings.Contains(raw, want) {
			t.Errorf("authorize URL is missing %s\ngot %s", want, raw)
		}
	}
	// State binds a redirect to the request that caused it. There is no
	// redirect, so sending one would be cargo cult: nothing on either side
	// would compare it.
	if strings.Contains(raw, "state=") {
		t.Errorf("the out-of-band flow has no redirect to protect with state\ngot %s", raw)
	}
	// An authorize URL the console already parameterises must keep its own
	// query, not have it replaced.
	raw, err = buildAuthorizeURL("https://console.example/cli/authorize?ref=docs", "c")
	if err != nil {
		t.Fatalf("buildAuthorizeURL: %v", err)
	}
	if !strings.Contains(raw, "ref=docs") {
		t.Errorf("an existing query parameter was dropped: %s", raw)
	}
}

// The code arrives through a clipboard and a terminal, either of which adds a
// newline. Sending it untrimmed would fail with invalid_grant, which the
// server deliberately does not explain — so the user would be told their
// correct paste was wrong.
func TestRedeemTrimsThePastedCode(t *testing.T) {
	var got url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		got = r.PostForm
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"access_token":"at","refresh_token":"rt","token_type":"Bearer","expires_in":3600}`)
	}))
	defer srv.Close()

	req, err := NewRequest(&Metadata{AuthorizationEndpoint: "https://console.example/cli/authorize"})
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	toks, err := req.Redeem(context.Background(), srv.Client(), srv.URL, "  the-code\n")
	if err != nil {
		t.Fatalf("Redeem: %v", err)
	}
	if toks.AccessToken != "at" || toks.RefreshToken != "rt" {
		t.Errorf("tokens = %+v", toks)
	}
	if got.Get("code") != "the-code" {
		t.Errorf("code = %q, want the-code", got.Get("code"))
	}
	// The redirect the server re-checks the code against (RFC 6749 4.1.3).
	if got.Get("redirect_uri") != RedirectURIOOB {
		t.Errorf("redirect_uri = %q, want %q", got.Get("redirect_uri"), RedirectURIOOB)
	}
	// The verifier is what proves this is the process that started the flow.
	// A request that reached the endpoint without one would be a client with
	// no authentication at all.
	sum := sha256.Sum256([]byte(got.Get("code_verifier")))
	if !strings.Contains(req.AuthorizationURL,
		"code_challenge="+base64.RawURLEncoding.EncodeToString(sum[:])) {
		t.Error("the verifier sent does not match the challenge that was approved")
	}
}

// An empty paste must not be spent against the token endpoint. It cannot
// succeed, and the answer it would come back with — invalid_grant — reads as
// "your code was wrong" rather than "you pressed enter".
func TestRedeemRefusesAnEmptyCode(t *testing.T) {
	req, err := NewRequest(&Metadata{AuthorizationEndpoint: "https://console.example/x"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := req.Redeem(context.Background(), http.DefaultClient, "http://127.0.0.1:1", "   \n"); err == nil {
		t.Fatal("an empty code must be refused before it is sent")
	}
}

// A deployment with no authorization endpoint cannot serve this flow. Saying
// so beats opening a browser at an empty URL.
func TestSupportsInteractiveLogin(t *testing.T) {
	cases := []struct {
		name string
		m    Metadata
		want bool
	}{
		{"configured", Metadata{AuthorizationEndpoint: "https://c/x", GrantTypesSupported: []string{"client_credentials", "authorization_code"}}, true},
		{"no endpoint", Metadata{GrantTypesSupported: []string{"authorization_code"}}, false},
		{"endpoint but grant not advertised", Metadata{AuthorizationEndpoint: "https://c/x", GrantTypesSupported: []string{"client_credentials"}}, false},
		{"empty", Metadata{}, false},
	}
	for _, tc := range cases {
		if got := tc.m.SupportsInteractiveLogin(); got != tc.want {
			t.Errorf("%s: got %v want %v", tc.name, got, tc.want)
		}
	}
}
