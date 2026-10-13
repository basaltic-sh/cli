package output

import (
	"encoding/json"
	"github.com/basaltic-sh/sdk-go/storage"
	"reflect"
	"testing"
)

func TestSnapshotSizeJSONPreservesZeroAndCapacity(t *testing.T) {
	for _, raw := range []string{
		`0`, `4096`, `10737418240`,
	} {
		var snapshot storage.Snapshot
		if err := json.Unmarshal([]byte(`{"size_bytes":`+raw+`}`), &snapshot); err != nil {
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
		if !reflect.DeepEqual(want, got["size_bytes"]) {
			t.Fatalf("CLI changed snapshot size: want %s got %s", raw, out.String())
		}
	}
}
