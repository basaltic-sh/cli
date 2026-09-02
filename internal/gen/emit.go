package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const header = `// Code generated from the Basaltic SDK manifest (api.json). DO NOT EDIT.
//
// Regenerate with:
//
//	make generate SDK=/path/to/sdk-go
`

// emitService writes one service's command tree.
func emitService(outDir string, svc service, resources []resourceGroup) error {
	var b strings.Builder

	pkg := svc.Package
	fmt.Fprintf(&b, "func init() { cli.RegisterService(new%sCommand) }\n\n", exported(svc.Name))

	// Service command.
	short := serviceShort[svc.Name]
	if short == "" {
		short = firstSentence(svc.Description)
	}
	fmt.Fprintf(&b, "// new%sCommand builds `basaltic %s`.\n", exported(svc.Name), svc.Name)
	fmt.Fprintf(&b, "func new%sCommand(state *cli.State) *cobra.Command {\n", exported(svc.Name))
	fmt.Fprintf(&b, "\tcmd := &cobra.Command{\n\t\tUse:   %q,\n\t\tShort: %q,\n", svc.Name, short)
	if aliases := serviceAliases[svc.Name]; len(aliases) > 0 {
		fmt.Fprintf(&b, "\t\tAliases: []string{%s},\n", quoteList(aliases))
	}
	if svc.Regional {
		b.WriteString("\t\tLong: " + quote(short+".\n\nThis is a regional service: it acts in the region from --region, the\nBASALTIC_REGION environment variable, or the profile.") + ",\n")
	}
	b.WriteString("\t}\n")
	// A service with a single resource of the same name would otherwise read
	// as `basaltic certificate certificate list`. Its verbs hang directly off
	// the service instead.
	if collapsed(svc, resources) {
		for _, op := range resources[0].Ops {
			fmt.Fprintf(&b, "\tcmd.AddCommand(new%s%s%sCommand(state))\n",
				exported(svc.Name), exported(resources[0].Name), exported(op.Verb))
		}
	} else {
		for _, r := range resources {
			fmt.Fprintf(&b, "\tcmd.AddCommand(new%s%sCommand(state))\n", exported(svc.Name), exported(r.Name))
		}
	}
	b.WriteString("\treturn cmd\n}\n\n")

	// Client helper.
	fmt.Fprintf(&b, "// %sClient builds the service client, resolving credentials on first use.\n", unexported(svc.Name))
	fmt.Fprintf(&b, "func %sClient(state *cli.State) (*%s.Client, error) {\n", unexported(svc.Name), pkg)
	b.WriteString("\tcfg, err := state.SDK()\n\tif err != nil {\n\t\treturn nil, err\n\t}\n")
	fmt.Fprintf(&b, "\treturn %s.New(cfg), nil\n}\n\n", pkg)

	for _, r := range resources {
		if collapsed(svc, resources) {
			// The resource command itself is not emitted; only its verbs.
			for _, op := range r.Ops {
				if err := emitOperation(&b, svc, r, op); err != nil {
					return err
				}
			}
			continue
		}
		if err := emitResource(&b, svc, r); err != nil {
			return err
		}
	}

	src := renderFile(svc, b.String())
	formatted, err := format.Source(src)
	if err != nil {
		broken := filepath.Join(outDir, svc.Name+"_gen.go.broken")
		_ = os.WriteFile(broken, src, 0o644)
		return fmt.Errorf("%s: %w (unformatted source at %s)", svc.Name, err, broken)
	}
	return os.WriteFile(filepath.Join(outDir, svc.Name+"_gen.go"), formatted, 0o644)
}

func renderFile(svc service, body string) []byte {
	var b bytes.Buffer
	b.WriteString(header)
	b.WriteString("\npackage generated\n\n")

	imports := []string{}
	code := stripComments(body)
	for marker, imp := range map[string]string{
		"json.Unmarshal": "encoding/json",
		"fmt.":           "fmt",
		"os.":            "os",
		"strings.":       "strings",
		"io.Reader":      "io",
	} {
		if strings.Contains(code, marker) {
			imports = append(imports, imp)
		}
	}
	sort.Strings(imports)

	b.WriteString("import (\n")
	for _, imp := range imports {
		fmt.Fprintf(&b, "\t%q\n", imp)
	}
	if len(imports) > 0 {
		b.WriteString("\n")
	}
	b.WriteString("\t\"github.com/spf13/cobra\"\n\n")
	if strings.Contains(code, "basaltic.") {
		b.WriteString("\tbasaltic \"github.com/basaltic-sh/sdk-go\"\n")
	}
	fmt.Fprintf(&b, "\t%q\n", "github.com/basaltic-sh/sdk-go/"+svc.Package)
	b.WriteString("\n\t\"github.com/basaltic-sh/cli/internal/cli\"\n")
	b.WriteString(")\n\n")
	b.WriteString(body)
	return b.Bytes()
}

func stripComments(s string) string {
	var b strings.Builder
	for _, line := range strings.Split(s, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			continue
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}

// collapsed reports whether a service's single resource shares its name, in
// which case the resource level adds a word and no information.
func collapsed(svc service, resources []resourceGroup) bool {
	return len(resources) == 1 && resources[0].Name == singular(svc.Name)
}

// singular trims a plural service name so "secrets" matches the resource
// "secret".
func singular(s string) string {
	if strings.HasSuffix(s, "ss") || !strings.HasSuffix(s, "s") {
		return s
	}
	return strings.TrimSuffix(s, "s")
}

// resourceGroup is one noun and the operations under it.
type resourceGroup struct {
	Name string
	Ops  []operation
}

func emitResource(b *strings.Builder, svc service, r resourceGroup) error {
	fn := fmt.Sprintf("new%s%sCommand", exported(svc.Name), exported(r.Name))
	fmt.Fprintf(b, "// %s builds `basaltic %s %s`.\n", fn, svc.Name, r.Name)
	fmt.Fprintf(b, "func %s(state *cli.State) *cobra.Command {\n", fn)
	fmt.Fprintf(b, "\tcmd := &cobra.Command{\n\t\tUse:     %q,\n\t\tShort:   %q,\n\t\tAliases: []string{%q},\n\t}\n",
		r.Name, title(plural(r.Name)), plural(r.Name))
	for _, op := range r.Ops {
		fmt.Fprintf(b, "\tcmd.AddCommand(new%s%s%sCommand(state))\n",
			exported(svc.Name), exported(r.Name), exported(op.Verb))
	}
	b.WriteString("\treturn cmd\n}\n\n")

	for _, op := range r.Ops {
		if err := emitOperation(b, svc, r, op); err != nil {
			return err
		}
	}
	return nil
}
