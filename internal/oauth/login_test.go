package oauth

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
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

// The state comparison is the CLIENT's half of the flow, and nothing on the
// server substitutes for it: the platform will happily issue a code for a
// legitimate authorization, and this is what stops a code from a DIFFERENT
// authorization being fed to this listener.
func TestCallbackRefusesMismatchedState(t *testing.T) {
	result := make(chan callbackResult, 1)
	h := callbackHandler("the-real-state", result)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/callback?code=stolen&state=some-other-state", nil)
	h.ServeHTTP(rec, req)

	got := <-result
	if got.err == nil {
		t.Fatal("a mismatched state must be refused")
	}
	if got.code != "" {
		t.Errorf("a refused callback must not yield a code, got %q", got.code)
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestCallbackAcceptsMatchingState(t *testing.T) {
	result := make(chan callbackResult, 1)
	h := callbackHandler("s", result)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/callback?code=abc123&state=s", nil))

	got := <-result
	if got.err != nil {
		t.Fatalf("unexpected error: %v", got.err)
	}
	if got.code != "abc123" {
		t.Errorf("code = %q, want abc123", got.code)
	}
	// The page renders the code in its URL, so it must not be stored.
	if cc := rec.Header().Get("Cache-Control"); cc != "no-store" {
		t.Errorf("Cache-Control = %q, want no-store", cc)
	}
}

// A refusal on the consent page comes back as ?error=..., not as a transport
// failure. Reporting it as "no code returned" would send the user looking for
// a bug instead of telling them they pressed Cancel.
func TestCallbackSurfacesAuthorizationError(t *testing.T) {
	result := make(chan callbackResult, 1)
	h := callbackHandler("s", result)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet,
		"/callback?state=s&error=access_denied&error_description=You+cancelled", nil))

	got := <-result
	if got.err == nil {
		t.Fatal("an authorization error must be reported")
	}
	if !strings.Contains(got.err.Error(), "You cancelled") {
		t.Errorf("the description should reach the user, got %q", got.err)
	}
}

func TestBuildAuthorizeURL(t *testing.T) {
	raw, err := buildAuthorizeURL("https://console.example/cli/authorize",
		"http://127.0.0.1:53682/callback", "the-challenge", "the-state")
	if err != nil {
		t.Fatalf("buildAuthorizeURL: %v", err)
	}
	for _, want := range []string{
		"client_id=" + ClientID,
		"code_challenge=the-challenge",
		"code_challenge_method=S256",
		"state=the-state",
		"redirect_uri=http%3A%2F%2F127.0.0.1%3A53682%2Fcallback",
	} {
		if !strings.Contains(raw, want) {
			t.Errorf("authorize URL is missing %s\ngot %s", want, raw)
		}
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
