// Code generated from the Basaltic SDK manifest (api.json). DO NOT EDIT.
//
// Regenerate with:
//
//	make generate SDK=/path/to/sdk-go

package generated

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	basaltic "github.com/basaltic-sh/sdk-go"
	"github.com/basaltic-sh/sdk-go/telemetry"

	"github.com/basaltic-sh/cli/internal/cli"
)

func init() { cli.RegisterService(newTelemetryCommand) }

// newTelemetryCommand builds `basaltic telemetry`.
func newTelemetryCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "telemetry",
		Short:   "Logs, metrics and traces",
		Aliases: []string{"tel"},
		Long:    "Logs, metrics and traces.\n\nThis is a regional service: it acts in the region from --region, the\nBASALTIC_REGION environment variable, or the profile.",
	}
	cmd.AddCommand(newTelemetryAlertCommand(state))
	cmd.AddCommand(newTelemetryLogCommand(state))
	cmd.AddCommand(newTelemetryLogGroupCommand(state))
	cmd.AddCommand(newTelemetryMetricCommand(state))
	cmd.AddCommand(newTelemetryRetainedDataCommand(state))
	cmd.AddCommand(newTelemetrySpanCommand(state))
	cmd.AddCommand(newTelemetryTraceCommand(state))
	cmd.AddCommand(newTelemetryTraceSettingsCommand(state))
	return cmd
}

// telemetryClient builds the service client, resolving credentials on first use.
func telemetryClient(state *cli.State, path string) (*telemetry.Client, error) {
	cfg, err := state.ServiceSDK("telemetry", path)
	if err != nil {
		return nil, err
	}
	return telemetry.New(cfg), nil
}

// newTelemetryAlertCommand builds `basaltic telemetry alert`.
func newTelemetryAlertCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "alert",
		Short:   "Alerts",
		Aliases: []string{"alerts"},
	}
	cmd.AddCommand(newTelemetryAlertListCommand(state))
	cmd.AddCommand(newTelemetryAlertGetCommand(state))
	cmd.AddCommand(newTelemetryAlertCreateCommand(state))
	cmd.AddCommand(newTelemetryAlertUpdateCommand(state))
	cmd.AddCommand(newTelemetryAlertDeleteCommand(state))
	return cmd
}

// newTelemetryAlertListCommand builds `basaltic telemetry alert list`.
func newTelemetryAlertListCommand(state *cli.State) *cobra.Command {
	var params telemetry.ListAlertsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List metric alerts",
		Args:  cobra.ExactArgs(0),
		Long:  "List metric alerts.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/alerts")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListAlertsAll(cmd.Context(), &params))
			}
			page, err := c.ListAlerts(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact canonical alert CRN, intersected with name and state")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Exact name filter; an explicitly empty name matches nothing")
	f.StringVar(&params.State, "state", "", "Filter effective state")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newTelemetryAlertGetCommand builds `basaltic telemetry alert get`.
func newTelemetryAlertGetCommand(state *cli.State) *cobra.Command {
	var scope telemetry.ListAlertsParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Get a metric alert",
		Args:  cobra.ExactArgs(1),
		Long:  "Get a metric alert.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/alerts/{id}")
			if err != nil {
				return err
			}
			out, err := c.GetAlertByReference(cmd.Context(), args[0], &scope)
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

// newTelemetryAlertCreateCommand builds `basaltic telemetry alert create`.
func newTelemetryAlertCreateCommand(state *cli.State) *cobra.Command {
	var body telemetry.CreateAlertRequest
	var bodyFile string
	var enabledFlag bool
	var expressionFlag string
	var silencedUntilFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a metric alert",
		Args:  cobra.ExactArgs(0),
		Long:  "Create a metric alert.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"enabled":        "enabled",
				"expression":     "expression",
				"for_seconds":    "for-seconds",
				"name":           "name",
				"operator":       "operator",
				"silenced_until": "silenced-until",
				"threshold":      "threshold",
				"topic":          "topic",
				"window_seconds": "window-seconds",
			}, []string{"expression", "for_seconds", "name", "operator", "threshold", "topic", "window_seconds"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/alerts")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("enabled") {
				body.Enabled = &enabledFlag
			}
			if expressionFlag != "" {
				if err := json.Unmarshal([]byte(expressionFlag), &body.Expression); err != nil {
					return fmt.Errorf("--expression: %w", err)
				}
			}
			if silencedUntilFlag != "" {
				parsed, err := parseTime(silencedUntilFlag)
				if err != nil {
					return fmt.Errorf("--silenced-until: %w", err)
				}
				body.SilencedUntil = &parsed
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateAlert(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.BoolVar(&enabledFlag, "enabled", false, "Enabled")
	f.StringVar(&expressionFlag, "expression", "", "Expression (JSON)")
	f.IntVar(&body.ForSeconds, "for-seconds", 0, "Continuous breach duration before firing; 0 fires immediately")
	f.StringVar(&body.Name, "name", "", "Immutable account-unique name")
	f.StringVar(&body.Operator, "operator", "", "Operator (one of: gt, gte, lt, lte)")
	f.StringVar(&silencedUntilFlag, "silenced-until", "", "Mute notifications until this instant (at most 30 days ahead) (RFC 3339)")
	f.Float64Var(&body.Threshold, "threshold", 0, "Threshold")
	f.StringVar(&body.Topic, "topic", "", "Required same-account, same-region notification topic reference: UUID, canonical CRN or exact name, classified by syntax")
	f.IntVar(&body.WindowSeconds, "window-seconds", 0, "Whole-minute window")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newTelemetryAlertUpdateCommand builds `basaltic telemetry alert update`.
func newTelemetryAlertUpdateCommand(state *cli.State) *cobra.Command {
	var body telemetry.UpdateAlertRequest
	var bodyFile string
	var enabledFlag bool
	var expressionFlag string
	var forSecondsFlag int
	var operatorFlag string
	var silencedUntilFlag string
	var thresholdFlag float64
	var topicFlag string
	var windowSecondsFlag int
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a metric alert",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"enabled":        "enabled",
				"expression":     "expression",
				"for_seconds":    "for-seconds",
				"operator":       "operator",
				"silenced_until": "silenced-until",
				"threshold":      "threshold",
				"topic":          "topic",
				"window_seconds": "window-seconds",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/alerts/{id}")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("enabled") {
				body.Enabled = &enabledFlag
			}
			if expressionFlag != "" {
				if err := json.Unmarshal([]byte(expressionFlag), &body.Expression); err != nil {
					return fmt.Errorf("--expression: %w", err)
				}
			}
			if cmd.Flags().Changed("for-seconds") {
				body.ForSeconds = &forSecondsFlag
			}
			if cmd.Flags().Changed("operator") {
				body.Operator = &operatorFlag
			}
			if silencedUntilFlag != "" {
				parsed, err := parseTime(silencedUntilFlag)
				if err != nil {
					return fmt.Errorf("--silenced-until: %w", err)
				}
				body.SilencedUntil = &parsed
			}
			if cmd.Flags().Changed("threshold") {
				body.Threshold = &thresholdFlag
			}
			if cmd.Flags().Changed("topic") {
				body.Topic = &topicFlag
			}
			if cmd.Flags().Changed("window-seconds") {
				body.WindowSeconds = &windowSecondsFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.UpdateAlert(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.BoolVar(&enabledFlag, "enabled", false, "Enabled")
	f.StringVar(&expressionFlag, "expression", "", "Expression (JSON)")
	f.IntVar(&forSecondsFlag, "for-seconds", 0, "Continuous breach duration before firing; 0 fires immediately")
	f.StringVar(&operatorFlag, "operator", "", "Operator (one of: gt, gte, lt, lte)")
	f.StringVar(&silencedUntilFlag, "silenced-until", "", "Mute notifications until this instant (at most 30 days ahead) (RFC 3339)")
	f.Float64Var(&thresholdFlag, "threshold", 0, "Threshold")
	f.StringVar(&topicFlag, "topic", "", "Required same-account, same-region notification topic reference: UUID, canonical CRN or exact name, classified by syntax")
	f.IntVar(&windowSecondsFlag, "window-seconds", 0, "Whole-minute window")
	return cmd
}

// newTelemetryAlertDeleteCommand builds `basaltic telemetry alert delete`.
func newTelemetryAlertDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a metric alert",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/alerts/{id}")
			if err != nil {
				return err
			}
			if err := c.DeleteAlert(cmd.Context(), args[0]); err != nil {
				return err
			}
			state.Printer().Done("Deleted.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newTelemetryLogCommand builds `basaltic telemetry log`.
func newTelemetryLogCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "log",
		Short:   "Logs",
		Aliases: []string{"logs"},
	}
	cmd.AddCommand(newTelemetryLogListCommand(state))
	cmd.AddCommand(newTelemetryLogGetCommand(state))
	cmd.AddCommand(newTelemetryLogIngestCommand(state))
	return cmd
}

// newTelemetryLogListCommand builds `basaltic telemetry log list`.
func newTelemetryLogListCommand(state *cli.State) *cobra.Command {
	var params telemetry.SearchLogsParams
	var fromFlag string
	var toFlag string
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Search log records",
		Args:  cobra.ExactArgs(0),
		Long:  "Search log records.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/logs")
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
				return state.Printer().Iter(c.SearchLogsAll(cmd.Context(), &params))
			}
			page, err := c.SearchLogs(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&fromFlag, "from", "", "Lower bound on timestamp (inclusive) (RFC 3339)")
	_ = cmd.MarkFlagRequired("from")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.LogGroup, "log-group", "", "Filter by log-group UUID, CRN or exact name in the authenticated account")
	f.StringVar(&params.LogStream, "log-stream", "", "Filter by stream name within the group")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.MinSeverity, "min-severity", "", "Filter to records with severity >= this band One of: \"TRACE\", \"DEBUG\", \"INFO\", \"WARN\", \"ERROR\", \"FATAL\"")
	f.StringVar(&params.Q, "q", "", "Case-insensitive substring on the log body")
	f.StringVar(&params.Region, "emitter-region", "", "Filter by emitter region")
	f.StringVar(&toFlag, "to", "", "Upper bound on timestamp (exclusive); must be after from and within 31d (RFC 3339)")
	_ = cmd.MarkFlagRequired("to")
	f.StringVar(&params.TraceID, "trace-id", "", "Filter to records carrying this W3C trace id")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newTelemetryLogGetCommand builds `basaltic telemetry log get`.
func newTelemetryLogGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <log-id>",
		Short: "Get a single log record by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/logs/{log_id}")
			if err != nil {
				return err
			}
			out, err := c.GetLog(cmd.Context(), args[0])
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

// newTelemetryLogIngestCommand builds `basaltic telemetry log ingest`.
func newTelemetryLogIngestCommand(state *cli.State) *cobra.Command {
	var body telemetry.IngestRequest
	var bodyFile string
	var logsFlag string
	cmd := &cobra.Command{
		Use:   "ingest",
		Short: "Ingest a batch of log records",
		Args:  cobra.ExactArgs(0),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"logs": "logs",
			}, []string{"logs"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/logs")
			if err != nil {
				return err
			}
			if logsFlag != "" {
				if err := json.Unmarshal([]byte(logsFlag), &body.Logs); err != nil {
					return fmt.Errorf("--logs: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.IngestLogs(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&logsFlag, "logs", "", "Logs (JSON)")
	return cmd
}

// newTelemetryLogGroupCommand builds `basaltic telemetry log-group`.
func newTelemetryLogGroupCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "log-group",
		Short:   "Log groups",
		Aliases: []string{"log-groups"},
	}
	cmd.AddCommand(newTelemetryLogGroupListCommand(state))
	cmd.AddCommand(newTelemetryLogGroupGetCommand(state))
	cmd.AddCommand(newTelemetryLogGroupCreateCommand(state))
	cmd.AddCommand(newTelemetryLogGroupUpdateCommand(state))
	cmd.AddCommand(newTelemetryLogGroupDeleteCommand(state))
	return cmd
}

// newTelemetryLogGroupListCommand builds `basaltic telemetry log-group list`.
func newTelemetryLogGroupListCommand(state *cli.State) *cobra.Command {
	var params telemetry.ListLogGroupsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List log groups (or look up one by name)",
		Args:  cobra.ExactArgs(0),
		Long:  "List log groups (or look up one by name).\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/log-groups")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListLogGroupsAll(cmd.Context(), &params))
			}
			page, err := c.ListLogGroups(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact telemetry/log-group CRN filter in the authenticated account and serving region")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Exact name filter, intersected with crn and pagination")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newTelemetryLogGroupGetCommand builds `basaltic telemetry log-group get`.
func newTelemetryLogGroupGetCommand(state *cli.State) *cobra.Command {
	var scope telemetry.ListLogGroupsParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Get a log group by id",
		Args:  cobra.ExactArgs(1),
		Long:  "Get a log group by id.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/log-groups/{id}")
			if err != nil {
				return err
			}
			out, err := c.GetLogGroupByReference(cmd.Context(), args[0], &scope)
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

// newTelemetryLogGroupCreateCommand builds `basaltic telemetry log-group create`.
func newTelemetryLogGroupCreateCommand(state *cli.State) *cobra.Command {
	var body telemetry.CreateLogGroupRequest
	var bodyFile string
	var descriptionFlag string
	var kmsKeyFlag string
	var retentionDaysFlag int
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a log group",
		Args:  cobra.ExactArgs(0),
		Long:  "Create a log group.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"description":    "description",
				"kms_key":        "kms-key",
				"name":           "name",
				"retention_days": "retention-days",
				"tags":           "tags",
			}, []string{"name"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/log-groups")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
			}
			if cmd.Flags().Changed("kms-key") {
				body.KMSKey = &kmsKeyFlag
			}
			if cmd.Flags().Changed("retention-days") {
				body.RetentionDays = &retentionDaysFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateLogGroup(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&descriptionFlag, "description", "", "Description")
	f.StringVar(&kmsKeyFlag, "kms-key", "", "KMS key UUID, CRN or exact name in the authenticated account and serving region")
	f.StringVar(&body.Name, "name", "", "Resource names must not start with the literal crn: prefix or be UUIDs (canonical, compact, braced, or urn:uuid: forms, in either case)")
	f.IntVar(&retentionDaysFlag, "retention-days", 0, "Retention in days (one of: 1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1096, 1827, 2192, 2557, 2922, 3288, 3653, <nil>)")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newTelemetryLogGroupUpdateCommand builds `basaltic telemetry log-group update`.
func newTelemetryLogGroupUpdateCommand(state *cli.State) *cobra.Command {
	var body telemetry.UpdateLogGroupRequest
	var bodyFile string
	var descriptionFlag string
	var kmsKeyFlag string
	var retentionDaysFlag int
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a log group",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"description":    "description",
				"kms_key":        "kms-key",
				"retention_days": "retention-days",
				"tags":           "tags",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/log-groups/{id}")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
			}
			if cmd.Flags().Changed("kms-key") {
				body.KMSKey = &kmsKeyFlag
			}
			if cmd.Flags().Changed("retention-days") {
				body.RetentionDays = &retentionDaysFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.UpdateLogGroup(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&descriptionFlag, "description", "", "Description")
	f.StringVar(&kmsKeyFlag, "kms-key", "", "KMS key UUID, CRN or exact name in this account and region")
	f.IntVar(&retentionDaysFlag, "retention-days", 0, "Omit to preserve retention (one of: 1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1096, 1827, 2192, 2557, 2922, 3288, 3653, <nil>)")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newTelemetryLogGroupDeleteCommand builds `basaltic telemetry log-group delete`.
func newTelemetryLogGroupDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a log group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/log-groups/{id}")
			if err != nil {
				return err
			}
			if err := c.DeleteLogGroup(cmd.Context(), args[0]); err != nil {
				return err
			}
			state.Printer().Done("Deleted.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newTelemetryMetricCommand builds `basaltic telemetry metric`.
func newTelemetryMetricCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "metric",
		Short:   "Metrics",
		Aliases: []string{"metrics"},
	}
	cmd.AddCommand(newTelemetryMetricListNamesCommand(state))
	cmd.AddCommand(newTelemetryMetricListNamesBodyCommand(state))
	cmd.AddCommand(newTelemetryMetricListSeriesCommand(state))
	cmd.AddCommand(newTelemetryMetricListSeriesBodyCommand(state))
	cmd.AddCommand(newTelemetryMetricQueryCommand(state))
	cmd.AddCommand(newTelemetryMetricQueryBodyCommand(state))
	cmd.AddCommand(newTelemetryMetricQueryRangeCommand(state))
	cmd.AddCommand(newTelemetryMetricQueryRangeBodyCommand(state))
	cmd.AddCommand(newTelemetryMetricWriteCommand(state))
	return cmd
}

// newTelemetryMetricListNamesCommand builds `basaltic telemetry metric list-names`.
func newTelemetryMetricListNamesCommand(state *cli.State) *cobra.Command {
	var params telemetry.ListMetricNamesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list-names",
		Short: "List metric names",
		Args:  cobra.ExactArgs(0),
		Long:  "List metric names.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/metrics/names")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListMetricNamesAll(cmd.Context(), &params))
			}
			page, err := c.ListMetricNames(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.End, "end", "", "End")
	_ = cmd.MarkFlagRequired("end")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Start, "start", "", "Start")
	_ = cmd.MarkFlagRequired("start")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newTelemetryMetricListNamesBodyCommand builds `basaltic telemetry metric list-names-body`.
func newTelemetryMetricListNamesBodyCommand(state *cli.State) *cobra.Command {
	var body telemetry.ListMetricNamesPostRequest
	var bodyFile string
	var limitFlag int
	var markerFlag string
	cmd := &cobra.Command{
		Use:   "list-names-body",
		Short: "List metric names (JSON)",
		Args:  cobra.ExactArgs(0),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"end":    "end",
				"limit":  "limit",
				"marker": "marker",
				"start":  "start",
			}, []string{"end", "start"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/metrics/names")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("limit") {
				body.Limit = &limitFlag
			}
			if cmd.Flags().Changed("marker") {
				body.Marker = &markerFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			page, err := c.ListMetricNamesPost(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.End, "end", "", "End")
	f.IntVar(&limitFlag, "limit", 0, "Limit")
	f.StringVar(&markerFlag, "marker", "", "Marker")
	f.StringVar(&body.Start, "start", "", "Start")
	return cmd
}

// newTelemetryMetricListSeriesCommand builds `basaltic telemetry metric list-series`.
func newTelemetryMetricListSeriesCommand(state *cli.State) *cobra.Command {
	var params telemetry.ListMetricSeriesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list-series",
		Short: "List metric series",
		Args:  cobra.ExactArgs(0),
		Long:  "List metric series.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/metrics/series")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListMetricSeriesAll(cmd.Context(), &params))
			}
			page, err := c.ListMetricSeries(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.End, "end", "", "End")
	_ = cmd.MarkFlagRequired("end")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Metric, "metric", "", "Metric")
	_ = cmd.MarkFlagRequired("metric")
	f.StringVar(&params.Start, "start", "", "Start")
	_ = cmd.MarkFlagRequired("start")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newTelemetryMetricListSeriesBodyCommand builds `basaltic telemetry metric list-series-body`.
func newTelemetryMetricListSeriesBodyCommand(state *cli.State) *cobra.Command {
	var body telemetry.ListMetricSeriesPostRequest
	var bodyFile string
	var limitFlag int
	var markerFlag string
	cmd := &cobra.Command{
		Use:   "list-series-body",
		Short: "List metric series (JSON)",
		Args:  cobra.ExactArgs(0),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"end":    "end",
				"limit":  "limit",
				"marker": "marker",
				"metric": "metric",
				"start":  "start",
			}, []string{"end", "metric", "start"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/metrics/series")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("limit") {
				body.Limit = &limitFlag
			}
			if cmd.Flags().Changed("marker") {
				body.Marker = &markerFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			page, err := c.ListMetricSeriesPost(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.End, "end", "", "End")
	f.IntVar(&limitFlag, "limit", 0, "Limit")
	f.StringVar(&markerFlag, "marker", "", "Marker")
	f.StringVar(&body.Metric, "metric", "", "Metric")
	f.StringVar(&body.Start, "start", "", "Start")
	return cmd
}

// newTelemetryMetricQueryCommand builds `basaltic telemetry metric query`.
func newTelemetryMetricQueryCommand(state *cli.State) *cobra.Command {
	var params telemetry.QueryMetricsInstantParams
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Run an instant metric query",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/metrics/query")
			if err != nil {
				return err
			}
			out, err := c.QueryMetricsInstant(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.Agg, "agg", "", "Aggregation One of: \"avg\", \"sum\", \"min\", \"max\", \"count\", \"last\", \"rate\", \"increase\"")
	f.StringSliceVar(&params.By, "by", nil, "Group-by label names")
	f.StringSliceVar(&params.Match, "match", nil, "Label matchers (e.g")
	f.StringVar(&params.Metric, "metric", "", "Metric name")
	_ = cmd.MarkFlagRequired("metric")
	f.StringVar(&params.Step, "step", "", "Lookback window duration (e.g")
	f.StringVar(&params.Time, "time", "", "Evaluation timestamp (RFC 3339 or unix seconds)")
	return cmd
}

// newTelemetryMetricQueryBodyCommand builds `basaltic telemetry metric query-body`.
func newTelemetryMetricQueryBodyCommand(state *cli.State) *cobra.Command {
	var body telemetry.QueryMetricsInstantPostRequest
	var bodyFile string
	var aggFlag string
	var stepFlag string
	var timeFlag string
	cmd := &cobra.Command{
		Use:   "query-body",
		Short: "Run an instant metric query (JSON)",
		Args:  cobra.ExactArgs(0),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"agg":    "agg",
				"by":     "by",
				"match":  "match",
				"metric": "metric",
				"step":   "step",
				"time":   "time",
			}, []string{"metric"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/metrics/query")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("agg") {
				body.Agg = &aggFlag
			}
			if cmd.Flags().Changed("step") {
				body.Step = &stepFlag
			}
			if cmd.Flags().Changed("time") {
				body.Time = &timeFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.QueryMetricsInstantPost(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&aggFlag, "agg", "", "Agg (one of: avg, sum, min, max, count, last, rate, increase)")
	f.StringSliceVar(&body.By, "by", nil, "By")
	f.StringSliceVar(&body.Match, "match", nil, "Match")
	f.StringVar(&body.Metric, "metric", "", "Metric")
	f.StringVar(&stepFlag, "step", "", "Step")
	f.StringVar(&timeFlag, "time", "", "Time")
	return cmd
}

// newTelemetryMetricQueryRangeCommand builds `basaltic telemetry metric query-range`.
func newTelemetryMetricQueryRangeCommand(state *cli.State) *cobra.Command {
	var params telemetry.QueryMetricsRangeParams
	cmd := &cobra.Command{
		Use:   "query-range",
		Short: "Run a range metric query",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/metrics/query-range")
			if err != nil {
				return err
			}
			out, err := c.QueryMetricsRange(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.Agg, "agg", "", "One of: \"avg\", \"sum\", \"min\", \"max\", \"count\", \"last\", \"rate\", \"increase\"")
	f.StringSliceVar(&params.By, "by", nil, "By")
	f.StringVar(&params.End, "end", "", "End")
	_ = cmd.MarkFlagRequired("end")
	f.StringSliceVar(&params.Match, "match", nil, "Match")
	f.StringVar(&params.Metric, "metric", "", "Metric")
	_ = cmd.MarkFlagRequired("metric")
	f.StringVar(&params.Start, "start", "", "Start")
	_ = cmd.MarkFlagRequired("start")
	f.StringVar(&params.Step, "step", "", "Bucket width duration (e.g")
	return cmd
}

// newTelemetryMetricQueryRangeBodyCommand builds `basaltic telemetry metric query-range-body`.
func newTelemetryMetricQueryRangeBodyCommand(state *cli.State) *cobra.Command {
	var body telemetry.QueryMetricsRangePostRequest
	var bodyFile string
	var aggFlag string
	var stepFlag string
	cmd := &cobra.Command{
		Use:   "query-range-body",
		Short: "Run a range metric query (JSON)",
		Args:  cobra.ExactArgs(0),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"agg":    "agg",
				"by":     "by",
				"end":    "end",
				"match":  "match",
				"metric": "metric",
				"start":  "start",
				"step":   "step",
			}, []string{"end", "metric", "start"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/metrics/query-range")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("agg") {
				body.Agg = &aggFlag
			}
			if cmd.Flags().Changed("step") {
				body.Step = &stepFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.QueryMetricsRangePost(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&aggFlag, "agg", "", "Agg (one of: avg, sum, min, max, count, last, rate, increase)")
	f.StringSliceVar(&body.By, "by", nil, "By")
	f.StringVar(&body.End, "end", "", "End")
	f.StringSliceVar(&body.Match, "match", nil, "Match")
	f.StringVar(&body.Metric, "metric", "", "Metric")
	f.StringVar(&body.Start, "start", "", "Start")
	f.StringVar(&stepFlag, "step", "", "Step")
	return cmd
}

// newTelemetryMetricWriteCommand builds `basaltic telemetry metric write`.
func newTelemetryMetricWriteCommand(state *cli.State) *cobra.Command {
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "write",
		Short: "Prometheus remote_write ingest",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/metrics/write")
			if err != nil {
				return err
			}
			reader, closeBody, err := openBody(bodyFile)
			if err != nil {
				return err
			}
			defer closeBody()
			if err := c.WriteMetrics(cmd.Context(), reader); err != nil {
				return err
			}
			state.Printer().Done("Write requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "file", "f", "-", "File to send as the request body, or - for stdin.")
	return cmd
}

// newTelemetryRetainedDataCommand builds `basaltic telemetry retained-data`.
func newTelemetryRetainedDataCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "retained-data",
		Short:   "Retained datas",
		Aliases: []string{"retained-datas"},
	}
	cmd.AddCommand(newTelemetryRetainedDataGetCommand(state))
	return cmd
}

// newTelemetryRetainedDataGetCommand builds `basaltic telemetry retained-data get`.
func newTelemetryRetainedDataGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Check retained telemetry presence",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/retained-data")
			if err != nil {
				return err
			}
			out, err := c.GetRetainedTelemetryPresence(cmd.Context())
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

// newTelemetrySpanCommand builds `basaltic telemetry span`.
func newTelemetrySpanCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "span",
		Short:   "Spans",
		Aliases: []string{"spans"},
	}
	cmd.AddCommand(newTelemetrySpanIngestCommand(state))
	return cmd
}

// newTelemetrySpanIngestCommand builds `basaltic telemetry span ingest`.
func newTelemetrySpanIngestCommand(state *cli.State) *cobra.Command {
	var body telemetry.IngestSpansRequest
	var bodyFile string
	var spansFlag string
	cmd := &cobra.Command{
		Use:   "ingest",
		Short: "Ingest a batch of trace spans",
		Args:  cobra.ExactArgs(0),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"spans": "spans",
			}, []string{"spans"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/spans")
			if err != nil {
				return err
			}
			if spansFlag != "" {
				if err := json.Unmarshal([]byte(spansFlag), &body.Spans); err != nil {
					return fmt.Errorf("--spans: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.IngestSpans(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&spansFlag, "spans", "", "Spans (JSON)")
	return cmd
}

// newTelemetryTraceCommand builds `basaltic telemetry trace`.
func newTelemetryTraceCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "trace",
		Short:   "Traces",
		Aliases: []string{"traces"},
	}
	cmd.AddCommand(newTelemetryTraceListCommand(state))
	cmd.AddCommand(newTelemetryTraceGetCommand(state))
	cmd.AddCommand(newTelemetryTraceListSpansCommand(state))
	return cmd
}

// newTelemetryTraceListCommand builds `basaltic telemetry trace list`.
func newTelemetryTraceListCommand(state *cli.State) *cobra.Command {
	var params telemetry.SearchTracesParams
	var fromFlag string
	var toFlag string
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List traces",
		Args:  cobra.ExactArgs(0),
		Long:  "List traces.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/traces")
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
				return state.Printer().Iter(c.SearchTracesAll(cmd.Context(), &params))
			}
			page, err := c.SearchTraces(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&fromFlag, "from", "", "Lower bound on span start_time (inclusive) (RFC 3339)")
	_ = cmd.MarkFlagRequired("from")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.Float64Var(&params.MinDurationMs, "min-duration-ms", 0, "Filter to traces whose end-to-end duration is at least this many ms")
	f.StringVar(&params.Operation, "operation", "", "Filter by span name (operation)")
	f.StringVar(&params.Service, "service", "", "Filter by service.name")
	f.StringVar(&params.StatusCode, "status-code", "", "Filter by status code One of: \"UNSET\", \"OK\", \"ERROR\"")
	f.StringVar(&toFlag, "to", "", "Upper bound on span start_time (exclusive); must be after from and within 31d (RFC 3339)")
	_ = cmd.MarkFlagRequired("to")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newTelemetryTraceGetCommand builds `basaltic telemetry trace get`.
func newTelemetryTraceGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <trace-id>",
		Short: "Get a trace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/traces/{trace_id}")
			if err != nil {
				return err
			}
			out, err := c.GetTrace(cmd.Context(), args[0])
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

// newTelemetryTraceListSpansCommand builds `basaltic telemetry trace list-spans`.
func newTelemetryTraceListSpansCommand(state *cli.State) *cobra.Command {
	var params telemetry.ListTraceSpansParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list-spans <trace-id>",
		Short: "List spans in a trace",
		Args:  cobra.ExactArgs(1),
		Long:  "List spans in a trace.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/traces/{trace_id}/spans")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListTraceSpansAll(cmd.Context(), args[0], &params))
			}
			page, err := c.ListTraceSpans(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newTelemetryTraceSettingsCommand builds `basaltic telemetry trace-settings`.
func newTelemetryTraceSettingsCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "trace-settings",
		Short:   "Trace settingses",
		Aliases: []string{"trace-settingses"},
	}
	cmd.AddCommand(newTelemetryTraceSettingsGetCommand(state))
	cmd.AddCommand(newTelemetryTraceSettingsUpdateCommand(state))
	cmd.AddCommand(newTelemetryTraceSettingsDeleteCommand(state))
	return cmd
}

// newTelemetryTraceSettingsGetCommand builds `basaltic telemetry trace-settings get`.
func newTelemetryTraceSettingsGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get the caller account's trace settings",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/trace-settings")
			if err != nil {
				return err
			}
			out, err := c.GetTraceSettings(cmd.Context())
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

// newTelemetryTraceSettingsUpdateCommand builds `basaltic telemetry trace-settings update`.
func newTelemetryTraceSettingsUpdateCommand(state *cli.State) *cobra.Command {
	var body telemetry.UpdateTraceSettingsRequest
	var bodyFile string
	var kmsKeyFlag string
	var retentionDaysFlag int
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update the caller account's trace settings",
		Args:  cobra.ExactArgs(0),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"kms_key":        "kms-key",
				"retention_days": "retention-days",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/trace-settings")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("kms-key") {
				body.KMSKey = &kmsKeyFlag
			}
			if cmd.Flags().Changed("retention-days") {
				body.RetentionDays = &retentionDaysFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.UpdateTraceSettings(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&kmsKeyFlag, "kms-key", "", "KMS key UUID, CRN or exact name in this account and region")
	f.IntVar(&retentionDaysFlag, "retention-days", 0, "Omit to preserve retention (one of: 1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1096, 1827, 2192, 2557, 2922, 3288, 3653, <nil>)")
	return cmd
}

// newTelemetryTraceSettingsDeleteCommand builds `basaltic telemetry trace-settings delete`.
func newTelemetryTraceSettingsDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete trace settings",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state, "/v1/trace-settings")
			if err != nil {
				return err
			}
			if err := c.DeleteTraceSettings(cmd.Context()); err != nil {
				return err
			}
			state.Printer().Done("Deleted.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}
