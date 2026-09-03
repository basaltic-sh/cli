// Command gen writes the CLI's command tree from the SDK's api.json manifest.
//
// The manifest describes the SDK's Go surface — every method name, its
// parameters and their types — so the CLI does not re-derive any of it from
// the OpenAPI specifications. Two implementations of the naming rules would
// drift, and this way the generator needs nothing private.
//
//	make generate SDK=../sdks/sdk-go
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// handWritten operations get a command written by hand instead of a generated
// one, and are skipped here.
//
// The serial console is the only one. It is a WebSocket upgrade carrying a raw
// tty, which needs a terminal in raw mode, an escape key and bidirectional
// binary frames — none of which a generated request/response command can
// express. The generated form would compile and then not work.
var handWritten = map[string]bool{
	"startSerialConsole": true,
}

// excluded services are not part of the release.
var excluded = map[string]bool{
	"registry":      true,
	"queue":         true,
	"notifications": true,
	"email":         true,
	"database":      true,
}

func main() {
	manifestPath := flag.String("manifest", "", "path to the SDK's api.json")
	outDir := flag.String("out", "../generated", "directory to write the command tree into")
	flag.Parse()

	if *manifestPath == "" {
		fmt.Fprintln(os.Stderr, "gen: -manifest is required: point it at the sdk-go checkout's api.json")
		flag.Usage()
		os.Exit(2)
	}
	if err := run(*manifestPath, *outDir); err != nil {
		fmt.Fprintln(os.Stderr, "gen:", err)
		os.Exit(1)
	}
}

func run(manifestPath, outDir string) error {
	m, err := loadManifest(manifestPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	// Remove previously generated files, so a service dropped from the
	// manifest does not linger as a command nothing can reach.
	old, _ := filepath.Glob(filepath.Join(outDir, "*_gen.go"))
	for _, f := range old {
		if err := os.Remove(f); err != nil {
			return err
		}
	}

	totalOps, totalCmds := 0, 0
	for _, svc := range m.Services {
		if excluded[svc.Name] {
			continue
		}
		groups := group(svc)
		if err := emitService(outDir, svc, groups); err != nil {
			return err
		}
		n := 0
		for _, g := range groups {
			n += len(g.Ops)
		}
		totalOps += n
		totalCmds += n + len(groups) + 1
		fmt.Printf("%-14s %3d commands in %2d resources\n", svc.Name, n, len(groups))
	}
	fmt.Printf("%-14s %3d commands, %d in the tree\n", "total", totalOps, totalCmds)
	return nil
}

// group buckets a service's operations by resource, in a stable order.
func group(svc service) []resourceGroup {
	byName := map[string][]operation{}
	for _, op := range svc.Operations {
		if handWritten[op.ID] {
			continue
		}
		byName[op.Resource] = append(byName[op.Resource], op)
	}
	names := make([]string, 0, len(byName))
	for n := range byName {
		names = append(names, n)
	}
	sort.Strings(names)

	out := make([]resourceGroup, 0, len(names))
	for _, n := range names {
		ops := byName[n]
		sort.Slice(ops, func(i, j int) bool { return verbOrder(ops[i].Verb) < verbOrder(ops[j].Verb) })
		out = append(out, resourceGroup{Name: n, Ops: ops})
	}
	return out
}

// verbOrder puts the common verbs first in `--help`, then everything else
// alphabetically. A reader scanning for "list" should not have to hunt past
// "attach-certificate".
func verbOrder(verb string) string {
	switch verb {
	case "list":
		return "0" + verb
	case "get":
		return "1" + verb
	case "create":
		return "2" + verb
	case "update":
		return "3" + verb
	case "delete":
		return "4" + verb
	}
	return "5" + verb
}
