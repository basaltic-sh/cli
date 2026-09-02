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
	cmd.AddCommand(newTelemetryLogCommand(state))
	cmd.AddCommand(newTelemetryLogGroupCommand(state))
	cmd.AddCommand(newTelemetryMetricCommand(state))
	cmd.AddCommand(newTelemetrySpanCommand(state))
	cmd.AddCommand(newTelemetryTraceCommand(state))
	cmd.AddCommand(newTelemetryTraceSettingsCommand(state))
	return cmd
}

// telemetryClient builds the service client, resolving credentials on first use.
func telemetryClient(state *cli.State) (*telemetry.Client, error) {
	cfg, err := state.SDK()
	if err != nil {
		return nil, err
	}
	return telemetry.New(cfg), nil
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
			c, err := telemetryClient(state)
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
	f.StringVar(&params.LogGroup, "log-group", "", "Filter by parent log group name")
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
			c, err := telemetryClient(state)
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if logsFlag != "" {
				if err := json.Unmarshal([]byte(logsFlag), &body.Logs); err != nil {
					return fmt.Errorf("--logs: %w", err)
				}
			}
			out, err := c.IngestLogs(cmd.Context(), &body)
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
	_ = cmd.MarkFlagRequired("logs")
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
			c, err := telemetryClient(state)
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
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Exact-match name lookup (single-result shortcut)")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newTelemetryLogGroupGetCommand builds `basaltic telemetry log-group get`.
func newTelemetryLogGroupGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a log group by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetLogGroup(cmd.Context(), args[0])
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
	var kmsKeyCrnFlag string
	var retentionDaysFlag int
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a log group",
		Args:  cobra.ExactArgs(0),
		Long:  "Create a log group.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
			}
			if cmd.Flags().Changed("kms-key-crn") {
				body.KMSKeyCRN = &kmsKeyCrnFlag
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
	f.StringVar(&kmsKeyCrnFlag, "kms-key-crn", "", "Kms key crn")
	f.StringVar(&body.Name, "name", "", "Name")
	_ = cmd.MarkFlagRequired("name")
	f.IntVar(&retentionDaysFlag, "retention-days", 0, "1..3650, or omit for never expire")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newTelemetryLogGroupUpdateCommand builds `basaltic telemetry log-group update`.
func newTelemetryLogGroupUpdateCommand(state *cli.State) *cobra.Command {
	var body telemetry.UpdateLogGroupRequest
	var bodyFile string
	var clearRetentionFlag bool
	var descriptionFlag string
	var kmsKeyCrnFlag string
	var retentionDaysFlag int
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a log group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("clear-retention") {
				body.ClearRetention = &clearRetentionFlag
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
			}
			if cmd.Flags().Changed("kms-key-crn") {
				body.KMSKeyCRN = &kmsKeyCrnFlag
			}
			if cmd.Flags().Changed("retention-days") {
				body.RetentionDays = &retentionDaysFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateLogGroup(cmd.Context(), args[0], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.BoolVar(&clearRetentionFlag, "clear-retention", false, "When true, sets retention to never-expire (ignores retention_days)")
	f.StringVar(&descriptionFlag, "description", "", "Description")
	f.StringVar(&kmsKeyCrnFlag, "kms-key-crn", "", "Pass an empty string to disassociate")
	f.IntVar(&retentionDaysFlag, "retention-days", 0, "1..3650; pass clear_retention to switch to never expire")
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
			c, err := telemetryClient(state)
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
	cmd := &cobra.Command{
		Use:   "list-names",
		Short: "List the distinct metric names emitted in a time window",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state)
			if err != nil {
				return err
			}
			out, err := c.ListMetricNames(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.End, "end", "", "End")
	_ = cmd.MarkFlagRequired("end")
	f.StringVar(&params.Start, "start", "", "Start")
	_ = cmd.MarkFlagRequired("start")
	return cmd
}

// newTelemetryMetricListNamesBodyCommand builds `basaltic telemetry metric list-names-body`.
func newTelemetryMetricListNamesBodyCommand(state *cli.State) *cobra.Command {
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "list-names-body",
		Short: "List the distinct metric names emitted in a time window (form body)",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state)
			if err != nil {
				return err
			}
			reader, closeBody, err := openBody(bodyFile)
			if err != nil {
				return err
			}
			defer closeBody()
			out, err := c.ListMetricNamesPost(cmd.Context(), reader)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "file", "f", "-", "File to send as the request body, or - for stdin.")
	return cmd
}

// newTelemetryMetricListSeriesCommand builds `basaltic telemetry metric list-series`.
func newTelemetryMetricListSeriesCommand(state *cli.State) *cobra.Command {
	var params telemetry.ListMetricSeriesParams
	cmd := &cobra.Command{
		Use:   "list-series",
		Short: "List distinct label sets for a metric",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state)
			if err != nil {
				return err
			}
			out, err := c.ListMetricSeries(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.End, "end", "", "End")
	_ = cmd.MarkFlagRequired("end")
	f.StringVar(&params.Metric, "metric", "", "Metric")
	_ = cmd.MarkFlagRequired("metric")
	f.StringVar(&params.Start, "start", "", "Start")
	_ = cmd.MarkFlagRequired("start")
	return cmd
}

// newTelemetryMetricListSeriesBodyCommand builds `basaltic telemetry metric list-series-body`.
func newTelemetryMetricListSeriesBodyCommand(state *cli.State) *cobra.Command {
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "list-series-body",
		Short: "List distinct label sets for a metric (form body)",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state)
			if err != nil {
				return err
			}
			reader, closeBody, err := openBody(bodyFile)
			if err != nil {
				return err
			}
			defer closeBody()
			out, err := c.ListMetricSeriesPost(cmd.Context(), reader)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "file", "f", "-", "File to send as the request body, or - for stdin.")
	return cmd
}

// newTelemetryMetricQueryCommand builds `basaltic telemetry metric query`.
func newTelemetryMetricQueryCommand(state *cli.State) *cobra.Command {
	var params telemetry.QueryMetricsInstantParams
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Instant structured metric query",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state)
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
	_ = cmd.MarkFlagRequired("agg")
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
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "query-body",
		Short: "Instant structured metric query (form body)",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state)
			if err != nil {
				return err
			}
			reader, closeBody, err := openBody(bodyFile)
			if err != nil {
				return err
			}
			defer closeBody()
			out, err := c.QueryMetricsInstantPost(cmd.Context(), reader)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "file", "f", "-", "File to send as the request body, or - for stdin.")
	return cmd
}

// newTelemetryMetricQueryRangeCommand builds `basaltic telemetry metric query-range`.
func newTelemetryMetricQueryRangeCommand(state *cli.State) *cobra.Command {
	var params telemetry.QueryMetricsRangeParams
	cmd := &cobra.Command{
		Use:   "query-range",
		Short: "Range structured metric query",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state)
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
	_ = cmd.MarkFlagRequired("agg")
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
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "query-range-body",
		Short: "Range structured metric query (form body)",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state)
			if err != nil {
				return err
			}
			reader, closeBody, err := openBody(bodyFile)
			if err != nil {
				return err
			}
			defer closeBody()
			out, err := c.QueryMetricsRangePost(cmd.Context(), reader)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "file", "f", "-", "File to send as the request body, or - for stdin.")
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
			c, err := telemetryClient(state)
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if spansFlag != "" {
				if err := json.Unmarshal([]byte(spansFlag), &body.Spans); err != nil {
					return fmt.Errorf("--spans: %w", err)
				}
			}
			out, err := c.IngestSpans(cmd.Context(), &body)
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
	_ = cmd.MarkFlagRequired("spans")
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
			c, err := telemetryClient(state)
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
		Short: "Get all spans for a trace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state)
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

// newTelemetryTraceSettingsCommand builds `basaltic telemetry trace-settings`.
func newTelemetryTraceSettingsCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "trace-settings",
		Short:   "Trace settingses",
		Aliases: []string{"trace-settingses"},
	}
	cmd.AddCommand(newTelemetryTraceSettingsGetCommand(state))
	cmd.AddCommand(newTelemetryTraceSettingsSetCommand(state))
	return cmd
}

// newTelemetryTraceSettingsGetCommand builds `basaltic telemetry trace-settings get`.
func newTelemetryTraceSettingsGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get the caller account's trace settings",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state)
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

// newTelemetryTraceSettingsSetCommand builds `basaltic telemetry trace-settings set`.
func newTelemetryTraceSettingsSetCommand(state *cli.State) *cobra.Command {
	var body telemetry.UpdateTraceSettingsRequest
	var bodyFile string
	var clearRetentionFlag bool
	var kmsKeyCrnFlag string
	var retentionDaysFlag int
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Update the caller account's trace settings",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := telemetryClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("clear-retention") {
				body.ClearRetention = &clearRetentionFlag
			}
			if cmd.Flags().Changed("kms-key-crn") {
				body.KMSKeyCRN = &kmsKeyCrnFlag
			}
			if cmd.Flags().Changed("retention-days") {
				body.RetentionDays = &retentionDaysFlag
			}
			out, err := c.PutTraceSettings(cmd.Context(), &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.BoolVar(&clearRetentionFlag, "clear-retention", false, "When true, sets retention to never-expire (ignores retention_days)")
	f.StringVar(&kmsKeyCrnFlag, "kms-key-crn", "", "Pass an empty string to disassociate")
	f.IntVar(&retentionDaysFlag, "retention-days", 0, "1..3650; pass clear_retention=true to switch to never-expire")
	return cmd
}
