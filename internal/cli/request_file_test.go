package cli_test

import (
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

func TestRequiredRequestFileFields(t *testing.T) {
	for _, tc := range []struct {
		name, input         string
		stdin               bool
		flags               []string
		wantName, wantError string
	}{
		{name: "JSON file", input: `{"name":"file-name","flavor":"small","networks":[]}`, wantName: "file-name"},
		{name: "YAML stdin", input: "name: file-name\nflavor: small\nnetworks: []\n", stdin: true, wantName: "file-name"},
		{name: "flags override file", input: `{"name":"file-name","flavor":"small","networks":[]}`, flags: []string{"--name=flag-name"}, wantName: "flag-name"},
		{name: "missing field", input: `{"name":"file-name","networks":[]}`, wantError: "required field(s) flavor not set"},
		{name: "null field", input: `{"name":null,"flavor":"small","networks":[]}`, wantError: "required field(s) name not set"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			seen := make(chan map[string]any, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			input := filepath.Join(dir, "body.json")
			if err := os.WriteFile(input, []byte(tc.input), 0600); err != nil {
				t.Fatal(err)
			}
			if tc.stdin {
				input = "-"
			}
			args := []string{"-test.run=^TestLifecycleCLIProcess$", "--", "compute", "instance", "create", "--from-file=" + input, "--output=json"}
			args = append(args, tc.flags...)
			cmd := exec.Command(os.Args[0], args...)
			cmd.Stdin = strings.NewReader(tc.input)
			cmd.Env = append(os.Environ(), "BASALTIC_LIFECYCLE_TEST_PROCESS=1", "BASALTIC_CONFIG_FILE="+config, "BASALTIC_PROFILE=default", "BASALTIC_ACCESS_TOKEN=test-only", "BASALTIC_NO_UPDATE_CHECK=1")
			out, err := cmd.CombinedOutput()
			if tc.wantError != "" {
				if err == nil || !strings.Contains(string(out), tc.wantError) {
					t.Fatalf("wanted %q, got %s (%v)", tc.wantError, out, err)
				}
				select {
				case <-seen:
					t.Fatal("invalid request sent")
				default:
				}
				return
			}
			if err != nil {
				t.Fatalf("CLI: %s (%v)", out, err)
			}
			select {
			case body := <-seen:
				if body["name"] != tc.wantName || body["flavor"] != "small" {
					t.Fatalf("wrong request: %#v", body)
				}
			default:
				t.Fatal("no request")
			}
		})
	}
}
