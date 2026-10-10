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

func TestAutoscalingFlagsPreserveWireValues(t *testing.T) {
	for _, tc := range []struct {
		name  string
		flags []string
		want  string
	}{
		{"bounds only", []string{"--min-count=0", "--max-count=3"}, `{"min_count":0,"max_count":3}`},
		{"explicit desired", []string{"--desired-count=0"}, `{"desired_count":0}`},
		{"compatibility alias", []string{"--replica-count=2"}, `{"replica_count":2}`},
		{"CPU and custom metrics", []string{`--autoscaling={"enabled":false,"drain_seconds":0,"metrics":[{"source":"cpu","target_type":"utilization","target_value":60},{"source":"telemetry","target_type":"average_value","target_value":100,"name":"requests","labels":{"service":"web"},"sample_aggregation":"rate"}]}`}, `{"autoscaling":{"enabled":false,"drain_seconds":0,"metrics":[{"source":"cpu","target_type":"utilization","target_value":60},{"source":"telemetry","target_type":"average_value","target_value":100,"name":"requests","labels":{"service":"web"},"sample_aggregation":"rate"}]}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			seen := make(chan map[string]any, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var body map[string]any
				if r.Method != "PATCH" {
					t.Errorf("unexpected method %s", r.Method)
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
				}
				seen <- body
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"load_balancer":{"id":"lb"}}`)
			}))
			defer server.Close()
			config := filepath.Join(t.TempDir(), "config.yaml")
			if err := os.WriteFile(config, []byte("profiles:\n  default:\n    region: test-region\n    endpoints:\n      loadbalancer: "+server.URL+"\n"), 0600); err != nil {
				t.Fatal(err)
			}
			args := []string{"-test.run=^TestLifecycleCLIProcess$", "--", "loadbalancer", "load-balancer", "update", "lb", "--output=json"}
			cmd := exec.Command(os.Args[0], append(args, tc.flags...)...)
			cmd.Env = append(os.Environ(), "BASALTIC_LIFECYCLE_TEST_PROCESS=1", "BASALTIC_CONFIG_FILE="+config, "BASALTIC_PROFILE=default", "BASALTIC_ACCESS_TOKEN=test-only", "BASALTIC_NO_UPDATE_CHECK=1")
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("CLI: %s (%v)", out, err)
			}
			var want map[string]any
			if err := json.Unmarshal([]byte(tc.want), &want); err != nil {
				t.Fatal(err)
			}
			select {
			case got := <-seen:
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("got %#v, want %#v", got, want)
				}
			default:
				t.Fatal("no request")
			}
		})
	}
}
