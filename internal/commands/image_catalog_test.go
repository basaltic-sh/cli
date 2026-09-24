package commands

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/basaltic-sh/cli/internal/output"
	basaltic "github.com/basaltic-sh/sdk-go"
	"github.com/basaltic-sh/sdk-go/compute"
)

func TestImageCatalogOutputShowsNamesAndPreservesStructuredImages(t *testing.T) {
	page := &basaltic.Page[compute.ImageCatalogCategory]{Items: []compute.ImageCatalogCategory{{
		Name: "platform", Images: []*compute.CatalogImage{{ID: "image", Name: "debian-13", Architecture: "amd64"}},
	}}}
	var out bytes.Buffer
	printer := output.Printer{Out: &out, Format: output.Text}
	if err := printer.Page(page); err != nil {
		t.Fatal(err)
	}
	for _, value := range []string{"platform", "debian-13", "amd64"} {
		if !strings.Contains(out.String(), value) {
			t.Fatalf("table omits %s: %s", value, out.String())
		}
	}
	out.Reset()
	printer.Format = output.JSON
	if err := printer.Page(page); err != nil {
		t.Fatal(err)
	}
	var categories []compute.ImageCatalogCategory
	if err := json.Unmarshal(out.Bytes(), &categories); err != nil {
		t.Fatal(err)
	}
	if len(categories) != 1 || len(categories[0].Images) != 1 || categories[0].Images[0].ID != "image" {
		t.Fatalf("structured output lost images: %s", out.String())
	}
}
