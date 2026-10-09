package cli_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstanceUserDataWireEncoding(t *testing.T) {
	cloudConfig := "#cloud-config\nbootcmd:\n  - echo boot\nusers:\n  - name: customer\nruncmd:\n  - echo run\n"
	encoded := base64.StdEncoding.EncodeToString([]byte(cloudConfig))
	for _, tc := range []struct {
		name, value string
		invalid     bool
	}{
		{"cloud config", encoded, false},
		{"invalid base64", "not base64!", true},
		{"explicit empty override", "", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			seen := make(chan map[string]any, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" || r.URL.Path != "/v1/instances" {
					t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
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
			dir := t.TempDir()
			config := filepath.Join(dir, "config.yaml")
			if err := os.WriteFile(config, []byte("profiles:\n  default:\n    region: test-region\n    endpoints:\n      compute: "+server.URL+"\n"), 0600); err != nil {
				t.Fatal(err)
			}
			bodyFile := filepath.Join(dir, "body.json")
			if err := os.WriteFile(bodyFile, []byte(`{"user_data":"`+encoded+`"}`), 0600); err != nil {
				t.Fatal(err)
			}
			args := []string{"-test.run=^TestLifecycleCLIProcess$", "--", "compute", "instance", "create", "--name=test", "--flavor=small", "--networks=[]", "--from-file=" + bodyFile, "--user-data=" + tc.value, "--output=json"}
			cmd := exec.Command(os.Args[0], args...)
			cmd.Env = append(os.Environ(), "BASALTIC_LIFECYCLE_TEST_PROCESS=1", "BASALTIC_CONFIG_FILE="+config, "BASALTIC_PROFILE=default", "BASALTIC_ACCESS_TOKEN=test-only", "BASALTIC_NO_UPDATE_CHECK=1")
			out, err := cmd.CombinedOutput()
			if tc.invalid {
				if err == nil || !strings.Contains(string(out), "--user-data: invalid base64") {
					t.Fatalf("invalid input accepted: %s (%v)", out, err)
				}
				select {
				case <-seen:
					t.Fatal("invalid input sent to API")
				default:
				}
				return
			}
			if err != nil {
				t.Fatalf("CLI: %v %s", err, out)
			}
			select {
			case body := <-seen:
				got, _ := body["user_data"].(string)
				if got != tc.value {
					t.Fatalf("wire user_data = %q, want %q", got, tc.value)
				}
				if tc.value != "" {
					decoded, err := base64.StdEncoding.DecodeString(got)
					if err != nil || string(decoded) != cloudConfig {
						t.Fatalf("one decode must recover cloud config: %q (%v)", decoded, err)
					}
				}
			default:
				t.Fatal("no create request")
			}
		})
	}
}
