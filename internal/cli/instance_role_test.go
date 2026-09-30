package cli_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestInstanceIAMRoleFlagsPreserveExplicitDetach(t *testing.T) {
	for _, tc := range []struct {
		name    string
		flags   []string
		want    string
		present bool
	}{
		{"omitted", []string{"--description=updated"}, "", false},
		{"attach", []string{"--iam-role=worker"}, "worker", true},
		{"detach", []string{"--iam-role="}, "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			seen := make(chan map[string]any, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PATCH" {
					t.Errorf("unexpected method %s", r.Method)
				}
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
				}
				seen <- body
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"instance":{"id":"vm"}}`)
			}))
			defer server.Close()
			config := filepath.Join(t.TempDir(), "config.yaml")
			if err := os.WriteFile(config, []byte("profiles:\n  default:\n    region: test-region\n    endpoints:\n      compute: "+server.URL+"\n"), 0600); err != nil {
				t.Fatal(err)
			}
			args := []string{"-test.run=^TestLifecycleCLIProcess$", "--", "compute", "instance", "update", "vm", "--output=json"}
			args = append(args, tc.flags...)
			cmd := exec.Command(os.Args[0], args...)
			cmd.Env = append(os.Environ(), "BASALTIC_LIFECYCLE_TEST_PROCESS=1", "BASALTIC_CONFIG_FILE="+config, "BASALTIC_PROFILE=default", "BASALTIC_ACCESS_TOKEN=test-only", "BASALTIC_NO_UPDATE_CHECK=1")
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("CLI: %v %s", err, out)
			}
			body := <-seen
			got, ok := body["iam_role"]
			if ok != tc.present || (tc.present && got != tc.want) {
				t.Fatalf("iam_role=%v, present=%v", got, ok)
			}
		})
	}
}
