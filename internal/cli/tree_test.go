package cli_test

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/basaltic-sh/cli/internal/cli"
	_ "github.com/basaltic-sh/cli/internal/commands"
	_ "github.com/basaltic-sh/cli/internal/generated"
)

// The command tree is generated from the SDK manifest, so these tests check
// properties the manifest cannot state rather than re-asserting its contents.
// A test that looked up "compute instance list" in the manifest and then found
// it in the tree would only be checking that the generator ran.

func rootCommand(t *testing.T) *cobra.Command {
	t.Helper()
	return cli.NewRootCommand(&cli.State{})
}

// walk visits every command, depth first.
func walk(cmd *cobra.Command, fn func(path string, c *cobra.Command)) {
	var rec func(prefix string, c *cobra.Command)
	rec = func(prefix string, c *cobra.Command) {
		path := strings.TrimSpace(prefix + " " + c.Name())
		fn(path, c)
		for _, sub := range c.Commands() {
			rec(path, sub)
		}
	}
	rec("", cmd)
}

func TestEveryCommandIsWellFormed(t *testing.T) {
	root := rootCommand(t)

	leaves, groups := 0, 0
	walk(root, func(path string, c *cobra.Command) {
		if c.Name() == "help" || c.Name() == "completion" || c.Parent() != nil && c.Parent().Name() == "completion" {
			return
		}
		if c.Short == "" {
			t.Errorf("%s: no short description, so it is blank in --help", path)
		}
		if strings.HasSuffix(c.Short, ".") {
			t.Errorf("%s: short description ends in a period, which cobra does not expect: %q", path, c.Short)
		}
		if c.Runnable() {
			leaves++
			// A leaf without Args accepts any number of positional arguments
			// and silently ignores the extras, so a typo in an id would run
			// against the wrong resource.
			if c.Args == nil {
				t.Errorf("%s: runnable but declares no Args validator", path)
			}
			return
		}
		groups++
		if len(c.Commands()) == 0 {
			t.Errorf("%s: neither runnable nor a group with children", path)
		}
	})

	// A guard on the walk: a wiring change that stopped registering the
	// generated tree would otherwise leave every assertion above vacuous.
	if leaves < 350 {
		t.Errorf("found %d runnable commands, want the whole generated surface (359 plus the built-ins)", leaves)
	}
	if groups < 60 {
		t.Errorf("found %d command groups, want one per service and resource", groups)
	}
}

func TestNoDuplicateNamesAmongSiblings(t *testing.T) {
	root := rootCommand(t)
	walk(root, func(path string, c *cobra.Command) {
		seen := map[string]string{}
		for _, sub := range c.Commands() {
			for _, name := range append([]string{sub.Name()}, sub.Aliases...) {
				if prev, dup := seen[name]; dup {
					t.Errorf("%s: %q is both %s and %s — one of them is unreachable", path, name, prev, sub.Name())
				}
				seen[name] = sub.Name()
			}
		}
	})
}

func TestNoFlagCollidesWithAGlobalFlag(t *testing.T) {
	root := rootCommand(t)
	global := map[string]bool{}
	root.PersistentFlags().VisitAll(func(f *pflag.Flag) { global[f.Name] = true })

	walk(root, func(path string, c *cobra.Command) {
		if !c.Runnable() {
			return
		}
		c.Flags().VisitAll(func(f *pflag.Flag) {
			// A local flag shadowing a global one means --region on a command
			// would set something else entirely.
			if global[f.Name] && c.LocalFlags().Lookup(f.Name) != nil && c.InheritedFlags().Lookup(f.Name) == nil {
				t.Errorf("%s: flag --%s shadows the global flag of the same name", path, f.Name)
			}
		})
	})
}

// The tree's shape is the product decision: service, then resource, then verb.
func TestTreeIsServiceResourceVerb(t *testing.T) {
	root := rootCommand(t)

	services := map[string]bool{}
	for _, c := range root.Commands() {
		if c.GroupID == "services" {
			services[c.Name()] = true
		}
	}
	for _, want := range []string{
		"audit", "billing", "certificate", "compute", "dns", "iam",
		"kms", "loadbalancer", "network", "quota", "secrets", "storage", "telemetry",
	} {
		if !services[want] {
			t.Errorf("no top-level command for the %s service", want)
		}
	}
	if len(services) != 13 {
		t.Errorf("found %d services at the top level, want 13", len(services))
	}

	// Nothing that is not part of the release should be reachable.
	for _, gone := range []string{"database", "registry", "queue", "notifications", "email", "topics", "repositories"} {
		if services[gone] {
			t.Errorf("%q is not part of the release but has a command", gone)
		}
	}

	// Verbs live at depth three, under a resource, not directly under a
	// service. That is what distinguishes this tree from the tag-derived one
	// it replaced.
	compute, _, err := root.Find([]string{"compute", "instance", "list"})
	if err != nil || compute.Name() != "list" {
		t.Fatalf("`compute instance list` did not resolve: %v", err)
	}
	if !compute.Runnable() {
		t.Error("`compute instance list` is not runnable")
	}
}

func TestResourcesAcceptTheirPluralAsAnAlias(t *testing.T) {
	root := rootCommand(t)
	for _, args := range [][]string{
		{"compute", "instances", "list"},
		{"network", "vpcs", "list"},
		{"storage", "volumes", "list"},
	} {
		cmd, _, err := root.Find(args)
		if err != nil || !cmd.Runnable() {
			t.Errorf("%v did not resolve through the plural alias: %v", args, err)
		}
	}
}

func TestServiceAliases(t *testing.T) {
	root := rootCommand(t)
	for alias, want := range map[string]string{"lb": "loadbalancer", "net": "network", "cert": "certificate"} {
		cmd, _, err := root.Find([]string{alias})
		if err != nil {
			t.Errorf("alias %q did not resolve: %v", alias, err)
			continue
		}
		if cmd.Name() != want {
			t.Errorf("alias %q resolved to %q, want %q", alias, cmd.Name(), want)
		}
	}
}

func TestPaginatedListsOfferAll(t *testing.T) {
	root := rootCommand(t)
	cmd, _, err := root.Find([]string{"compute", "instance", "list"})
	if err != nil {
		t.Fatal(err)
	}
	if cmd.Flags().Lookup("all") == nil {
		t.Error("a paginated list has no --all flag, so walking every page needs manual markers")
	}
	if cmd.Flags().Lookup("marker") == nil {
		t.Error("a paginated list has no --marker flag")
	}
}

func TestCreatesOfferIdempotencyAndAFileBody(t *testing.T) {
	root := rootCommand(t)
	cmd, _, err := root.Find([]string{"compute", "instance", "create"})
	if err != nil {
		t.Fatal(err)
	}
	if cmd.Flags().Lookup("idempotency-key") == nil {
		t.Error("create has no --idempotency-key, so a timed-out create cannot be retried safely")
	}
	if cmd.Flags().Lookup("from-file") == nil {
		t.Error("create has no --from-file, leaving nested request fields reachable only as JSON strings")
	}
	if f := cmd.Flags().Lookup("name"); f == nil {
		t.Error("create has no --name")
	}
}
