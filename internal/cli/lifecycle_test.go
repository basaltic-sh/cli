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

	"github.com/basaltic-sh/cli/internal/cli"
)

func TestLifecycleCLIProcess(t *testing.T) {
	if os.Getenv("BASALTIC_LIFECYCLE_TEST_PROCESS") != "1" {
		return
	}
	for i, arg := range os.Args {
		if arg == "--" {
			os.Args = append([]string{"basaltic"}, os.Args[i+1:]...)
			os.Exit(cli.Execute())
		}
	}
	t.Fatal("missing command arguments")
}

func TestLifecycleRevisionRoundTrip(t *testing.T) {
	requests := make(chan string, 8)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests <- r.Method + " " + r.Header.Get("If-Match")
		if r.URL.Path != "/v1/buckets/example/lifecycle" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("ETag", `"revision-1"`)
			fmt.Fprint(w, `{"lifecycle":{"rules":[]},"revision":"revision-1"}`)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	config := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(config, []byte("profiles:\n  default:\n    region: test-region\n    endpoints:\n      storage: "+server.URL+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	run := func(args ...string) ([]byte, error) {
		cmd := exec.Command(os.Args[0], append([]string{"-test.run=^TestLifecycleCLIProcess$", "--", "storage", "bucket"}, args...)...)
		cmd.Env = append(os.Environ(), "BASALTIC_LIFECYCLE_TEST_PROCESS=1", "BASALTIC_CONFIG_FILE="+config, "BASALTIC_PROFILE=default", "BASALTIC_ACCESS_TOKEN=test-only", "BASALTIC_NO_UPDATE_CHECK=1")
		return cmd.CombinedOutput()
	}
	out, err := run("get-lifecycle", "example", "--output=json")
	if err != nil {
		t.Fatalf("get: %v: %s", err, out)
	}
	var result struct {
		Revision  string `json:"revision"`
		Lifecycle struct {
			Rules []json.RawMessage `json:"rules"`
		} `json:"lifecycle"`
	}
	if err := json.Unmarshal(out, &result); err != nil || result.Revision != "revision-1" || result.Lifecycle.Rules == nil {
		t.Fatalf("revision missing from output: %s (%v)", out, err)
	}
	if got := <-requests; got != "GET " {
		t.Fatal(got)
	}
	for _, verb := range []string{"set-lifecycle", "delete-lifecycle"} {
		args := []string{verb, "example"}
		method := "DELETE"
		if verb == "set-lifecycle" {
			args = append(args, "--lifecycle", `{"rules":[]}`)
			method = "PUT"
		}
		out, err := run(args...)
		if err == nil || !strings.Contains(string(out), "revision") {
			t.Fatalf("%s accepted missing revision: %s (%v)", verb, out, err)
		}
		select {
		case got := <-requests:
			t.Fatalf("unprotected request: %s", got)
		default:
		}
		out, err = run(append(args, "--revision", result.Revision)...)
		if err != nil {
			t.Fatalf("%s: %s (%v)", verb, out, err)
		}
		if got := <-requests; got != method+` "revision-1"` {
			t.Fatal(got)
		}
	}
}
