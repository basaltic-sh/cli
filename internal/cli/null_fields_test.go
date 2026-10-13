package cli_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

func TestPatchFilePreservesNullAndFlagOverrides(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		flags      []string
		want       map[string]any
	}{
		{name: "json clear", body: `{"retention_days":null,"kms_key":null}`, want: map[string]any{"retention_days": nil, "kms_key": nil}},
		{name: "yaml clear", body: "retention_days: null\n", want: map[string]any{"retention_days": nil}},
		{name: "flag overrides null", body: `{"retention_days":null}`, flags: []string{"--retention-days", "7"}, want: map[string]any{"retention_days": float64(7)}},
		{name: "omitted", body: `{}`, want: map[string]any{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var got map[string]any
				if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
					t.Error(err)
				}
				if r.Method != "PATCH" || !reflect.DeepEqual(got, tc.want) {
					t.Errorf("%s body = %#v; want %#v", r.Method, got, tc.want)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"log_group":{"id":"group"}}`))
			}))
			defer server.Close()
			dir := t.TempDir()
			config := filepath.Join(dir, "config.yaml")
			body := filepath.Join(dir, "patch.json")
			if err := os.WriteFile(config, []byte("profiles:\n  default:\n    region: test-region\n    endpoints:\n      telemetry: "+server.URL+"\n"), 0600); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(body, []byte(tc.body), 0600); err != nil {
				t.Fatal(err)
			}
			args := []string{"-test.run=^TestLifecycleCLIProcess$", "--", "telemetry", "log-group", "update", "group", "--from-file", body}
			args = append(args, tc.flags...)
			cmd := exec.Command(os.Args[0], args...)
			cmd.Env = append(os.Environ(), "BASALTIC_LIFECYCLE_TEST_PROCESS=1", "BASALTIC_CONFIG_FILE="+config, "BASALTIC_PROFILE=default", "BASALTIC_ACCESS_TOKEN=test-only", "BASALTIC_NO_UPDATE_CHECK=1")
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("%s: %v", out, err)
			}
		})
	}
}
