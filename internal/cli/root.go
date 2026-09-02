// Package cli assembles the command tree and the context every command runs
// with.
package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	basaltic "github.com/basaltic-sh/sdk-go"
	"github.com/spf13/cobra"

	"github.com/basaltic-sh/cli/internal/auth"
	"github.com/basaltic-sh/cli/internal/output"
)

// contextKey keys the per-invocation state stashed on the command context.
type contextKey struct{}

// State is what every command needs: how to reach the platform, and how to
// print what comes back.
type State struct {
	opts     auth.Options
	printer  *output.Printer
	resolved *auth.Resolved
}

// Printer returns the configured renderer.
func (s *State) Printer() *output.Printer { return s.printer }

// SDK resolves credentials on first use.
//
// Deferred rather than done at startup so that `basaltic --help`, `version`
// and `config set` all work on a machine with no credentials at all.
func (s *State) SDK() (*basaltic.Config, error) {
	if s.resolved == nil {
		r, err := auth.Resolve(s.opts)
		if err != nil {
			return nil, err
		}
		s.resolved = r
	}
	return s.resolved.Config, nil
}

// Resolved exposes the full resolution for `auth status`.
func (s *State) Resolved() (*auth.Resolved, error) {
	if _, err := s.SDK(); err != nil {
		return nil, err
	}
	return s.resolved, nil
}

// FromContext retrieves the state a command runs with.
func FromContext(ctx context.Context) *State {
	s, _ := ctx.Value(contextKey{}).(*State)
	return s
}

// Execute builds the command tree and runs it.
func Execute() int {
	state := &State{printer: &output.Printer{Out: os.Stdout, Err: os.Stderr, Format: output.Text}}
	root := NewRootCommand(state)

	ctx := context.WithValue(context.Background(), contextKey{}, state)
	if err := root.ExecuteContext(ctx); err != nil {
		printError(os.Stderr, err)
		return 1
	}
	return 0
}

// NewRootCommand assembles the tree.
func NewRootCommand(state *State) *cobra.Command {
	var outputFormat string

	root := &cobra.Command{
		Use:   "basaltic",
		Short: "Command-line interface for the Basaltic cloud platform",
		Long: "Command-line interface for the Basaltic cloud platform.\n\n" +
			"Commands are grouped by service, then by resource:\n\n" +
			"    basaltic <service> <resource> <verb> [flags]\n\n" +
			"For example:\n\n" +
			"    basaltic compute instance list\n" +
			"    basaltic network vpc create --name prod --cidr 10.0.0.0/16\n" +
			"    basaltic storage volume attach vol-1 --instance-id i-1",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			f, err := output.ParseFormat(outputFormat)
			if err != nil {
				return err
			}
			state.printer.Format = f
			noHeaders, _ := cmd.Flags().GetBool("no-headers")
			state.printer.NoHeaders = noHeaders
			return nil
		},
	}

	pf := root.PersistentFlags()
	pf.StringVarP(&state.opts.Profile, "profile", "p", "", "configuration profile to use")
	pf.StringVar(&state.opts.APIKey, "api-key", "", "credential as ACCESS_KEY_ID:SECRET (overrides the profile)")
	pf.StringVar(&state.opts.Region, "region", "", "region for regional services (overrides the profile)")
	pf.StringVar(&state.opts.AccountID, "account-id", "", "account to act on, sent as X-Account-Id")
	pf.StringVarP(&outputFormat, "output", "o", "text", "output format: text, json or yaml")
	pf.Bool("no-headers", false, "omit the header row in text tables")
	pf.BoolVar(&state.opts.Insecure, "insecure", false, "skip TLS verification (development rigs only)")

	root.AddGroup(
		&cobra.Group{ID: "services", Title: "Services:"},
		&cobra.Group{ID: "cli", Title: "CLI:"},
	)

	for _, add := range serviceCommands {
		cmd := add(state)
		cmd.GroupID = "services"
		root.AddCommand(cmd)
	}
	for _, add := range builtinCommands {
		cmd := add(state)
		cmd.GroupID = "cli"
		root.AddCommand(cmd)
	}
	for _, g := range grafts {
		parent, _, err := root.Find(g.path)
		if err != nil || parent == nil {
			// A hand-written command whose parent the generator no longer
			// emits would otherwise vanish silently.
			panic(fmt.Sprintf("cannot graft onto %v: %v", g.path, err))
		}
		parent.AddCommand(g.add(state))
	}
	return root
}

// serviceCommands is populated by the generated packages, and builtinCommands
// by the hand-written ones, each through Register at init time.
var (
	serviceCommands []func(*State) *cobra.Command
	builtinCommands []func(*State) *cobra.Command
)

// RegisterService adds a generated service command tree.
func RegisterService(add func(*State) *cobra.Command) {
	serviceCommands = append(serviceCommands, add)
}

// RegisterBuiltin adds a hand-written top-level command.
func RegisterBuiltin(add func(*State) *cobra.Command) {
	builtinCommands = append(builtinCommands, add)
}

// graft is a hand-written command that belongs inside the generated tree.
type graft struct {
	path []string
	add  func(*State) *cobra.Command
}

var grafts []graft

// RegisterAt adds a hand-written command underneath a generated one, for the
// few operations a generated request/response command cannot express — the
// serial console, which is a WebSocket carrying a raw tty.
//
// The generator skips those operations, so nothing is being replaced here: the
// slot is empty and this fills it.
func RegisterAt(path []string, add func(*State) *cobra.Command) {
	grafts = append(grafts, graft{path: path, add: add})
}

// printError renders a failure the way the person reading it needs.
//
// The platform's own errors already say what went wrong; what the CLI adds is
// the next step, which differs sharply between classes that look alike on the
// surface.
func printError(w *os.File, err error) {
	fmt.Fprintln(w, "Error:", err)

	switch {
	case basaltic.IsUnauthorized(err):
		fmt.Fprintln(w, "\nThe credential was refused. Check `basaltic auth status`, or run `basaltic auth login`.")
	case basaltic.IsAccessDenied(err):
		fmt.Fprintln(w, "\nThe credential is valid but policy does not allow this. It needs a policy change, not a new key.")
	case basaltic.IsQuotaExceeded(err):
		fmt.Fprintln(w, "\nThis is a quota limit, not a rate limit. Retrying will not clear it; raising the quota will.")
	case basaltic.IsNotFound(err):
		fmt.Fprintln(w, "\nThe resource does not exist, or is not visible to this account.")
		fmt.Fprintln(w, "If it belongs to another account, pass --account-id.")
	case basaltic.IsRateLimited(err):
		fmt.Fprintln(w, "\nThrottled. The CLI already retried; wait before trying again.")
	}

	var authErr *basaltic.AuthError
	if errors.As(err, &authErr) && authErr.Code == "invalid_client" {
		fmt.Fprintln(w, "\nRun `basaltic auth login` to store a working credential.")
	}
}
