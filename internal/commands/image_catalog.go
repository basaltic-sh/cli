package commands

import (
	"strings"

	"github.com/basaltic-sh/cli/internal/output"
	"github.com/basaltic-sh/sdk-go/compute"
)

func init() { output.Present(compute.ImageCatalogCategory{}, imageCatalogRow) }

// Keep image names visible in text output: the generic renderer omits nested
// arrays. JSON and YAML retain the complete categorized API response.
func imageCatalogRow(value any) []output.Field {
	var category *compute.ImageCatalogCategory
	switch v := value.(type) {
	case compute.ImageCatalogCategory:
		category = &v
	case *compute.ImageCatalogCategory:
		category = v
	}
	names := make([]string, 0, len(category.Images))
	for _, image := range category.Images {
		names = append(names, image.Name+" ("+image.Architecture+")")
	}
	return []output.Field{
		{Key: "category", Value: category.Name},
		{Key: "images", Value: strings.Join(names, ", ")},
	}
}
