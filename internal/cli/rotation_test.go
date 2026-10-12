package cli_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

func TestKMSRotationWireContract(t *testing.T) {
	for _, tc := range []struct {
		name, method, path, token, body, response string
		status                                    int
		args                                      []string
	}{
		{name: "disable preserves explicit false", method: "PUT", path: "/v1/keys/key-1/rotation", body: `{"enabled":false,"rotation_period_days":90}`, response: `{"rotation":{"enabled":false,"rotation_period_days":90,"next_rotation_at":null}}`, status: 200, args: []string{"set-rotation", "key-1", "--enabled=false", "--rotation-period-days=90"}},
		{name: "accepted rotation retains retry token", method: "POST", path: "/v1/keys/key-1/rotate", token: "rotation-test-token", response: `{"rotation":{"id":"rotation-1","key_id":"key-1","trigger":"on_demand","status":"pending","previous_version":null,"key_version":null,"requested_at":"2026-10-12T00:00:00Z","started_at":null,"completed_at":null,"error_code":null}}`, status: 202, args: []string{"rotate", "key-1", "--idempotency-key=rotation-test-token"}},
		{name: "historical public key", method: "GET", path: "/v1/keys/key-1/public-key?key_version=2", response: `{"key_id":"key-1","key_version":2,"key_spec":"ecdsa-p256","public_key":"test"}`, status: 200, args: []string{"public-key", "key-1", "--key-version=2"}},
		{name: "versioned verification", method: "POST", path: "/v1/keys/key-1/verify", body: `{"message":"aGk=","signature":"c2ln","key_version":2}`, response: `{"signature_valid":true}`, status: 200, args: []string{"verify", "key-1", "--message=aGk=", "--signature=c2ln", "--key-version=2"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			seen := make(chan bool, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tc.method || r.URL.RequestURI() != tc.path {
					t.Errorf("request %s %s", r.Method, r.URL.RequestURI())
				}
				if got := r.Header.Get("Idempotency-Key"); got != tc.token {
					t.Errorf("token %q", got)
				}
				if tc.body != "" {
					var got, want map[string]any
					if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
						t.Error(err)
					}
					if err := json.Unmarshal([]byte(tc.body), &want); err != nil {
						t.Error(err)
					}
					if !reflect.DeepEqual(got, want) {
						t.Errorf("body %#v; want %#v", got, want)
					}
				}
				seen <- true
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				fmt.Fprint(w, tc.response)
			}))
			defer server.Close()
			config := filepath.Join(t.TempDir(), "config.yaml")
			if err := os.WriteFile(config, []byte("profiles:\n  default:\n    region: test-region\n    endpoints:\n      kms: "+server.URL+"\n"), 0600); err != nil {
				t.Fatal(err)
			}
			args := append([]string{"-test.run=^TestLifecycleCLIProcess$", "--", "kms", "key"}, tc.args...)
			args = append(args, "--output=json")
			cmd := exec.Command(os.Args[0], args...)
			cmd.Env = append(os.Environ(), "BASALTIC_LIFECYCLE_TEST_PROCESS=1", "BASALTIC_CONFIG_FILE="+config, "BASALTIC_PROFILE=default", "BASALTIC_ACCESS_TOKEN=test-only", "BASALTIC_NO_UPDATE_CHECK=1")
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("CLI: %s (%v)", out, err)
			}
			var output any
			if err := json.Unmarshal(out, &output); err != nil {
				t.Fatalf("Invalid JSON output %s: %v", out, err)
			}
			select {
			case <-seen:
			default:
				t.Fatal("no request")
			}
		})
	}
}
