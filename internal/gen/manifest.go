package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// The shape of the SDK's api.json. Only the fields the CLI needs are
// declared; the manifest carries more.

type manifest struct {
	Module   string    `json:"module"`
	Version  string    `json:"version"`
	Services []service `json:"services"`
}

type service struct {
	Name             string      `json:"name"`
	Package          string      `json:"package"`
	Title            string      `json:"title"`
	Description      string      `json:"description"`
	EndpointTemplate string      `json:"endpoint_template"`
	Regional         bool        `json:"regional"`
	Operations       []operation `json:"operations"`
}

type operation struct {
	ID         string  `json:"id"`
	GoName     string  `json:"go_name"`
	Method     string  `json:"method"`
	Path       string  `json:"path"`
	Summary    string  `json:"summary"`
	Resource   string  `json:"resource"`
	Verb       string  `json:"verb"`
	PathParams []param `json:"path_params"`
	ParamsType string  `json:"params_type"`
	Params     []field `json:"params"`
	BodyType   string  `json:"body_type"`
	BodyKind   string  `json:"body_kind"`
	BodyFields []field `json:"body_fields"`
	ResultKind string  `json:"result_kind"`
	ResultType string  `json:"result_type"`
	ItemType   string  `json:"item_type"`
	Paginated  bool    `json:"paginated"`
	Idempotent bool    `json:"idempotent"`
}

type param struct {
	Wire   string `json:"wire"`
	GoName string `json:"go_name"`
	Doc    string `json:"doc"`
}

type field struct {
	Wire     string   `json:"wire"`
	GoName   string   `json:"go_name"`
	GoType   string   `json:"go_type"`
	FlagKind string   `json:"flag_kind"`
	Pointer  bool     `json:"pointer"`
	Required bool     `json:"required"`
	Enum     []string `json:"enum"`
	Doc      string   `json:"doc"`
}

func loadManifest(path string) (*manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading the SDK manifest: %w\n\n"+
			"It is api.json at the root of the sdk-go checkout, written when the SDK\n"+
			"is generated. Pass -manifest /path/to/sdk-go/api.json", err)
	}
	var m manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	if len(m.Services) == 0 {
		return nil, fmt.Errorf("%s describes no services", path)
	}
	return &m, nil
}

// flagName renders a wire name as a command-line flag.
func flagName(wire string) string {
	return strings.ReplaceAll(strings.ReplaceAll(wire, "_", "-"), "[]", "")
}
