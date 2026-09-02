// Code generated from the Basaltic SDK manifest (api.json). DO NOT EDIT.
//
// Regenerate with:
//
//	make generate SDK=/path/to/sdk-go

package generated

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/basaltic-sh/sdk-go/audit"

	"github.com/basaltic-sh/cli/internal/cli"
)

func init() { cli.RegisterService(newAuditCommand) }

// newAuditCommand builds `basaltic audit`.
func newAuditCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Audit logs",
	}
	cmd.AddCommand(newAuditLogCommand(state))
	return cmd
}

// auditClient builds the service client, resolving credentials on first use.
func auditClient(state *cli.State) (*audit.Client, error) {
	cfg, err := state.SDK()
	if err != nil {
		return nil, err
	}
	return audit.New(cfg), nil
}

// newAuditLogCommand builds `basaltic audit log`.
func newAuditLogCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "log",
		Short:   "Logs",
		Aliases: []string{"logs"},
	}
	cmd.AddCommand(newAuditLogListCommand(state))
	cmd.AddCommand(newAuditLogGetCommand(state))
	return cmd
}

// newAuditLogListCommand builds `basaltic audit log list`.
func newAuditLogListCommand(state *cli.State) *cobra.Command {
	var params audit.ListAuditLogsParams
	var fromFlag string
	var toFlag string
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List audit logs",
		Args:  cobra.ExactArgs(0),
		Long:  "List audit logs.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auditClient(state)
			if err != nil {
				return err
			}
			if fromFlag != "" {
				parsed, err := parseTime(fromFlag)
				if err != nil {
					return fmt.Errorf("--from: %w", err)
				}
				params.From = parsed
			}
			if toFlag != "" {
				parsed, err := parseTime(toFlag)
				if err != nil {
					return fmt.Errorf("--to: %w", err)
				}
				params.To = parsed
			}
			if fetchAll {
				return state.Printer().Iter(c.ListAuditLogsAll(cmd.Context(), &params))
			}
			page, err := c.ListAuditLogs(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.Action, "action", "", "Filter by action (exact match or prefix with wildcard, e.g., \"iam.*\")")
	f.StringVar(&params.ActorID, "actor-id", "", "Filter by actor ID (user or service account)")
	f.StringVar(&params.ActorType, "actor-type", "", "Filter by actor type One of: \"user\", \"service_account\", \"system\"")
	f.StringVar(&fromFlag, "from", "", "Filter logs from this timestamp (inclusive) (RFC 3339)")
	f.StringVar(&params.IPAddress, "ip-address", "", "Filter by IP address")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.ResourceID, "resource-id", "", "Filter by resource ID")
	f.StringVar(&params.ResourceType, "resource-type", "", "Filter by resource type")
	f.StringVar(&params.Status, "status", "", "Filter by status One of: \"success\", \"failure\", \"denied\"")
	f.StringVar(&toFlag, "to", "", "Filter logs until this timestamp (exclusive) (RFC 3339)")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newAuditLogGetCommand builds `basaltic audit log get`.
func newAuditLogGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <log-id>",
		Short: "Get audit log entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auditClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetAuditLog(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}
