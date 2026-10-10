package output

import (
	"encoding/json"
	"github.com/basaltic-sh/sdk-go/storage"
	"reflect"
	"testing"
)

func TestSnapshotUsageJSONPreservesNullAndZero(t *testing.T) {
	for _, raw := range []string{
		`{"state":"unknown","scope":"volume_lineage","billable":false,"measured_at":null,"lineage_retained_bytes":null}`,
		`{"state":"stale","scope":"volume_lineage","billable":false,"measured_at":"2026-10-10T12:30:00Z","lineage_retained_bytes":null}`,
		`{"state":"measured","scope":"volume_lineage","billable":false,"measured_at":"2026-10-10T12:30:00Z","lineage_retained_bytes":0}`,
		`{"state":"measured","scope":"volume_lineage","billable":false,"measured_at":"2026-10-10T12:30:00Z","lineage_retained_bytes":4096}`,
	} {
		var snapshot storage.Snapshot
		if err := json.Unmarshal([]byte(`{"snapshot_usage":`+raw+`}`), &snapshot); err != nil {
			t.Fatal(err)
		}
		p, out, _ := newPrinter(JSON)
		if err := p.Value(&snapshot); err != nil {
			t.Fatal(err)
		}
		var got map[string]any
		var want any
		if err := json.Unmarshal(out.Bytes(), &got); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal([]byte(raw), &want); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(want, got["snapshot_usage"]) {
			t.Fatalf("CLI changed measurement: want %s got %s", raw, out.String())
		}
	}
}
