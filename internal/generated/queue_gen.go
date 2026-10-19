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
	"github.com/basaltic-sh/sdk-go/queue"

	"github.com/basaltic-sh/cli/internal/cli"
)

func init() { cli.RegisterService(newQueueCommand) }

// newQueueCommand builds `basaltic queue`.
func newQueueCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "queue",
		Short: "KMS bindings use `kms_key`: a UUID, account-scoped name, or complete `crn:kms:<region>:<account>:key/<name>` in the authenticated account and current region",
		Long:  "KMS bindings use `kms_key`: a UUID, account-scoped name, or complete `crn:kms:<region>:<account>:key/<name>` in the authenticated account and current region.\n\nThis is a regional service: it acts in the region from --region, the\nBASALTIC_REGION environment variable, or the profile.",
	}
	cmd.AddCommand(newQueueMessageCommand(state))
	cmd.AddCommand(newQueueQueueCommand(state))
	return cmd
}

// queueClient builds the service client, resolving credentials on first use.
func queueClient(state *cli.State, path string) (*queue.Client, error) {
	cfg, err := state.ServiceSDK("queue", path)
	if err != nil {
		return nil, err
	}
	return queue.New(cfg), nil
}

// newQueueMessageCommand builds `basaltic queue message`.
func newQueueMessageCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "message",
		Short:   "Messages",
		Aliases: []string{"messages"},
	}
	cmd.AddCommand(newQueueMessageDeleteCommand(state))
	cmd.AddCommand(newQueueMessageChangeVisibilityCommand(state))
	cmd.AddCommand(newQueueMessageChangeVisibilityBatchCommand(state))
	cmd.AddCommand(newQueueMessageDeleteBatchCommand(state))
	cmd.AddCommand(newQueueMessagePeekCommand(state))
	cmd.AddCommand(newQueueMessageReceiveCommand(state))
	cmd.AddCommand(newQueueMessageSendCommand(state))
	cmd.AddCommand(newQueueMessageSendBatchCommand(state))
	return cmd
}

// newQueueMessageDeleteCommand builds `basaltic queue message delete`.
func newQueueMessageDeleteCommand(state *cli.State) *cobra.Command {
	var params queue.DeleteMessageParams
	cmd := &cobra.Command{
		Use:   "delete <queue-id>",
		Short: "Delete a leased message",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := queueClient(state, "/v1/queues/{queue_id}/messages")
			if err != nil {
				return err
			}
			if err := c.DeleteMessage(cmd.Context(), args[0], &params); err != nil {
				return err
			}
			state.Printer().Done("Deleted.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.ReceiptHandle, "receipt-handle", "", "The handle returned by the receive that leased this message")
	_ = cmd.MarkFlagRequired("receipt-handle")
	return cmd
}

// newQueueMessageChangeVisibilityCommand builds `basaltic queue message change-visibility`.
func newQueueMessageChangeVisibilityCommand(state *cli.State) *cobra.Command {
	var body queue.ChangeVisibilityRequest
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "change-visibility <queue-id>",
		Short: "Extend or shorten a lease",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"receipt_handle":             "receipt-handle",
				"visibility_timeout_seconds": "visibility-timeout-seconds",
			}, []string{"receipt_handle", "visibility_timeout_seconds"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := queueClient(state, "/v1/queues/{queue_id}/messages/visibility")
			if err != nil {
				return err
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			if err := c.ChangeMessageVisibility(cmd.Context(), args[0], &body, reqOpts...); err != nil {
				return err
			}
			state.Printer().Done("Change visibility requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.ReceiptHandle, "receipt-handle", "", "Receipt handle")
	f.IntVar(&body.VisibilityTimeoutSeconds, "visibility-timeout-seconds", 0, "Measured from now, not from the original lease")
	return cmd
}

// newQueueMessageChangeVisibilityBatchCommand builds `basaltic queue message change-visibility-batch`.
func newQueueMessageChangeVisibilityBatchCommand(state *cli.State) *cobra.Command {
	var body queue.ChangeVisibilityBatchRequest
	var bodyFile string
	var entriesFlag string
	cmd := &cobra.Command{
		Use:   "change-visibility-batch <queue-id>",
		Short: "Extend or shorten up to ten leases",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"entries": "entries",
			}, []string{"entries"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := queueClient(state, "/v1/queues/{queue_id}/messages/visibility-batch")
			if err != nil {
				return err
			}
			if entriesFlag != "" {
				if err := json.Unmarshal([]byte(entriesFlag), &body.Entries); err != nil {
					return fmt.Errorf("--entries: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			page, err := c.ChangeMessageVisibilityBatch(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&entriesFlag, "entries", "", "Entries (JSON)")
	return cmd
}

// newQueueMessageDeleteBatchCommand builds `basaltic queue message delete-batch`.
func newQueueMessageDeleteBatchCommand(state *cli.State) *cobra.Command {
	var body queue.DeleteMessageBatchRequest
	var bodyFile string
	var entriesFlag string
	cmd := &cobra.Command{
		Use:   "delete-batch <queue-id>",
		Short: "Delete up to ten leased messages",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"entries": "entries",
			}, []string{"entries"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := queueClient(state, "/v1/queues/{queue_id}/messages/delete-batch")
			if err != nil {
				return err
			}
			if entriesFlag != "" {
				if err := json.Unmarshal([]byte(entriesFlag), &body.Entries); err != nil {
					return fmt.Errorf("--entries: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			page, err := c.DeleteMessageBatch(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&entriesFlag, "entries", "", "Entries (JSON)")
	return cmd
}

// newQueueMessagePeekCommand builds `basaltic queue message peek`.
func newQueueMessagePeekCommand(state *cli.State) *cobra.Command {
	var params queue.PeekMessagesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "peek <queue-id>",
		Short: "Peek at a queue's messages",
		Args:  cobra.ExactArgs(1),
		Long:  "Peek at a queue's messages.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := queueClient(state, "/v1/queues/{queue_id}/messages")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.PeekMessagesAll(cmd.Context(), args[0], &params))
			}
			page, err := c.PeekMessages(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact message CRN under this queue")
	f.IntVar(&params.Limit, "limit", 0, "Maximum items to return")
	f.StringVar(&params.Marker, "marker", "", "Id of the last item on the previous page")
	f.StringVar(&params.Visibility, "visibility", "", "Restrict the peek to one population")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newQueueMessageReceiveCommand builds `basaltic queue message receive`.
func newQueueMessageReceiveCommand(state *cli.State) *cobra.Command {
	var body queue.ReceiveRequest
	var bodyFile string
	var maxMessagesFlag int
	var visibilityTimeoutSecondsFlag int
	var waitTimeSecondsFlag int
	cmd := &cobra.Command{
		Use:   "receive <queue-id>",
		Short: "Lease messages",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"max_messages":               "max-messages",
				"visibility_timeout_seconds": "visibility-timeout-seconds",
				"wait_time_seconds":          "wait-time-seconds",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := queueClient(state, "/v1/queues/{queue_id}/messages/receive")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("max-messages") {
				body.MaxMessages = &maxMessagesFlag
			}
			if cmd.Flags().Changed("visibility-timeout-seconds") {
				body.VisibilityTimeoutSeconds = &visibilityTimeoutSecondsFlag
			}
			if cmd.Flags().Changed("wait-time-seconds") {
				body.WaitTimeSeconds = &waitTimeSecondsFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			page, err := c.ReceiveMessages(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.IntVar(&maxMessagesFlag, "max-messages", 0, "How many messages to lease")
	f.IntVar(&visibilityTimeoutSecondsFlag, "visibility-timeout-seconds", 0, "Overrides the queue's default for this lease only")
	f.IntVar(&waitTimeSecondsFlag, "wait-time-seconds", 0, "Block up to this long for a message rather than returning empty")
	return cmd
}

// newQueueMessageSendCommand builds `basaltic queue message send`.
func newQueueMessageSendCommand(state *cli.State) *cobra.Command {
	var body queue.SendMessageRequest
	var bodyFile string
	var attributesFlag string
	var delaySecondsFlag int
	var messageDeduplicationIdFlag string
	var messageGroupIdFlag string
	cmd := &cobra.Command{
		Use:   "send <queue-id>",
		Short: "Send a message",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"attributes":               "attributes",
				"body":                     "body",
				"delay_seconds":            "delay-seconds",
				"message_deduplication_id": "message-deduplication-id",
				"message_group_id":         "message-group-id",
			}, []string{"body"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := queueClient(state, "/v1/queues/{queue_id}/messages")
			if err != nil {
				return err
			}
			if attributesFlag != "" {
				if err := json.Unmarshal([]byte(attributesFlag), &body.Attributes); err != nil {
					return fmt.Errorf("--attributes: %w", err)
				}
			}
			if cmd.Flags().Changed("delay-seconds") {
				body.DelaySeconds = &delaySecondsFlag
			}
			if cmd.Flags().Changed("message-deduplication-id") {
				body.MessageDeduplicationID = &messageDeduplicationIdFlag
			}
			if cmd.Flags().Changed("message-group-id") {
				body.MessageGroupID = &messageGroupIdFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.SendMessage(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&attributesFlag, "attributes", "", "Attributes (JSON)")
	f.StringVar(&body.Body, "body", "", "The payload")
	f.IntVar(&delaySecondsFlag, "delay-seconds", 0, "Overrides the queue's default for this message")
	f.StringVar(&messageDeduplicationIdFlag, "message-deduplication-id", "", "FIFO only")
	f.StringVar(&messageGroupIdFlag, "message-group-id", "", "Required on FIFO queues, rejected on standard ones")
	return cmd
}

// newQueueMessageSendBatchCommand builds `basaltic queue message send-batch`.
func newQueueMessageSendBatchCommand(state *cli.State) *cobra.Command {
	var body queue.SendMessageBatchRequest
	var bodyFile string
	var entriesFlag string
	cmd := &cobra.Command{
		Use:   "send-batch <queue-id>",
		Short: "Send up to ten messages",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"entries": "entries",
			}, []string{"entries"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := queueClient(state, "/v1/queues/{queue_id}/messages/batch")
			if err != nil {
				return err
			}
			if entriesFlag != "" {
				if err := json.Unmarshal([]byte(entriesFlag), &body.Entries); err != nil {
					return fmt.Errorf("--entries: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			page, err := c.SendMessageBatch(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&entriesFlag, "entries", "", "Entries (JSON)")
	return cmd
}

// newQueueQueueCommand builds `basaltic queue queue`.
func newQueueQueueCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "queue",
		Short:   "Queues",
		Aliases: []string{"queues"},
	}
	cmd.AddCommand(newQueueQueueListCommand(state))
	cmd.AddCommand(newQueueQueueGetCommand(state))
	cmd.AddCommand(newQueueQueueCreateCommand(state))
	cmd.AddCommand(newQueueQueueUpdateCommand(state))
	cmd.AddCommand(newQueueQueueDeleteCommand(state))
	cmd.AddCommand(newQueueQueuePurgeCommand(state))
	return cmd
}

// newQueueQueueListCommand builds `basaltic queue queue list`.
func newQueueQueueListCommand(state *cli.State) *cobra.Command {
	var params queue.ListQueuesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List queues",
		Args:  cobra.ExactArgs(0),
		Long:  "List queues.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := queueClient(state, "/v1/queues")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListQueuesAll(cmd.Context(), &params))
			}
			page, err := c.ListQueues(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact queue CRN")
	f.IntVar(&params.Limit, "limit", 0, "Maximum items to return")
	f.StringVar(&params.Marker, "marker", "", "Id of the last item on the previous page")
	f.StringVar(&params.Name, "name", "", "Exact account-scoped queue name")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newQueueQueueGetCommand builds `basaltic queue queue get`.
func newQueueQueueGetCommand(state *cli.State) *cobra.Command {
	var scope queue.ListQueuesParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Describe a queue",
		Args:  cobra.ExactArgs(1),
		Long:  "Describe a queue.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := queueClient(state, "/v1/queues/{queue_id}")
			if err != nil {
				return err
			}
			out, err := c.DescribeQueueByReference(cmd.Context(), args[0], &scope)
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

// newQueueQueueCreateCommand builds `basaltic queue queue create`.
func newQueueQueueCreateCommand(state *cli.State) *cobra.Command {
	var body queue.CreateQueueRequest
	var bodyFile string
	var contentBasedDeduplicationFlag bool
	var delaySecondsFlag int
	var kmsDataKeyReusePeriodSecondsFlag int
	var kmsKeyFlag string
	var maxMessageBytesFlag int
	var messageRetentionSecondsFlag int
	var receiveWaitTimeSecondsFlag int
	var redrivePolicyFlag string
	var tagsFlag string
	var visibilityTimeoutSecondsFlag int
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a queue",
		Args:  cobra.ExactArgs(0),
		Long:  "Create a queue.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"content_based_deduplication":       "content-based-deduplication",
				"delay_seconds":                     "delay-seconds",
				"kms_data_key_reuse_period_seconds": "kms-data-key-reuse-period-seconds",
				"kms_key":                           "kms-key",
				"max_message_bytes":                 "max-message-bytes",
				"message_retention_seconds":         "message-retention-seconds",
				"name":                              "name",
				"receive_wait_time_seconds":         "receive-wait-time-seconds",
				"redrive_policy":                    "redrive-policy",
				"tags":                              "tags",
				"visibility_timeout_seconds":        "visibility-timeout-seconds",
			}, []string{"name"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := queueClient(state, "/v1/queues")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("content-based-deduplication") {
				body.ContentBasedDeduplication = &contentBasedDeduplicationFlag
			}
			if cmd.Flags().Changed("delay-seconds") {
				body.DelaySeconds = &delaySecondsFlag
			}
			if cmd.Flags().Changed("kms-data-key-reuse-period-seconds") {
				body.KMSDataKeyReusePeriodSeconds = &kmsDataKeyReusePeriodSecondsFlag
			}
			if cmd.Flags().Changed("kms-key") {
				body.KMSKey = &kmsKeyFlag
			}
			if cmd.Flags().Changed("max-message-bytes") {
				body.MaxMessageBytes = &maxMessageBytesFlag
			}
			if cmd.Flags().Changed("message-retention-seconds") {
				body.MessageRetentionSeconds = &messageRetentionSecondsFlag
			}
			if cmd.Flags().Changed("receive-wait-time-seconds") {
				body.ReceiveWaitTimeSeconds = &receiveWaitTimeSecondsFlag
			}
			if redrivePolicyFlag != "" {
				if err := json.Unmarshal([]byte(redrivePolicyFlag), &body.RedrivePolicy); err != nil {
					return fmt.Errorf("--redrive-policy: %w", err)
				}
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			if cmd.Flags().Changed("visibility-timeout-seconds") {
				body.VisibilityTimeoutSeconds = &visibilityTimeoutSecondsFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateQueue(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.BoolVar(&contentBasedDeduplicationFlag, "content-based-deduplication", false, "FIFO queues only")
	f.IntVar(&delaySecondsFlag, "delay-seconds", 0, "Delay seconds")
	f.IntVar(&kmsDataKeyReusePeriodSecondsFlag, "kms-data-key-reuse-period-seconds", 0, "Kms data key reuse period seconds")
	f.StringVar(&kmsKeyFlag, "kms-key", "", "KMS key UUID, CRN or account-scoped name")
	f.IntVar(&maxMessageBytesFlag, "max-message-bytes", 0, "Max message bytes")
	f.IntVar(&messageRetentionSecondsFlag, "message-retention-seconds", 0, "Message retention seconds")
	f.StringVar(&body.Name, "name", "", "Letters, digits, dot, dash, or underscore, starting with a letter or digit")
	f.IntVar(&receiveWaitTimeSecondsFlag, "receive-wait-time-seconds", 0, "Receive wait time seconds")
	f.StringVar(&redrivePolicyFlag, "redrive-policy", "", "Redrive policy (JSON)")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.IntVar(&visibilityTimeoutSecondsFlag, "visibility-timeout-seconds", 0, "Visibility timeout seconds")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newQueueQueueUpdateCommand builds `basaltic queue queue update`.
func newQueueQueueUpdateCommand(state *cli.State) *cobra.Command {
	var body queue.UpdateQueueRequest
	var bodyFile string
	var delaySecondsFlag int
	var kmsDataKeyReusePeriodSecondsFlag int
	var kmsKeyFlag string
	var maxMessageBytesFlag int
	var messageRetentionSecondsFlag int
	var receiveWaitTimeSecondsFlag int
	var redrivePolicyFlag string
	var tagsFlag string
	var visibilityTimeoutSecondsFlag int
	cmd := &cobra.Command{
		Use:   "update <queue-id>",
		Short: "Update queue attributes",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"delay_seconds":                     "delay-seconds",
				"kms_data_key_reuse_period_seconds": "kms-data-key-reuse-period-seconds",
				"kms_key":                           "kms-key",
				"max_message_bytes":                 "max-message-bytes",
				"message_retention_seconds":         "message-retention-seconds",
				"receive_wait_time_seconds":         "receive-wait-time-seconds",
				"redrive_policy":                    "redrive-policy",
				"tags":                              "tags",
				"visibility_timeout_seconds":        "visibility-timeout-seconds",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := queueClient(state, "/v1/queues/{queue_id}")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("delay-seconds") {
				body.DelaySeconds = &delaySecondsFlag
			}
			if cmd.Flags().Changed("kms-data-key-reuse-period-seconds") {
				body.KMSDataKeyReusePeriodSeconds = &kmsDataKeyReusePeriodSecondsFlag
			}
			if cmd.Flags().Changed("kms-key") {
				body.KMSKey = &kmsKeyFlag
			}
			if cmd.Flags().Changed("max-message-bytes") {
				body.MaxMessageBytes = &maxMessageBytesFlag
			}
			if cmd.Flags().Changed("message-retention-seconds") {
				body.MessageRetentionSeconds = &messageRetentionSecondsFlag
			}
			if cmd.Flags().Changed("receive-wait-time-seconds") {
				body.ReceiveWaitTimeSeconds = &receiveWaitTimeSecondsFlag
			}
			if redrivePolicyFlag != "" {
				if err := json.Unmarshal([]byte(redrivePolicyFlag), &body.RedrivePolicy); err != nil {
					return fmt.Errorf("--redrive-policy: %w", err)
				}
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			if cmd.Flags().Changed("visibility-timeout-seconds") {
				body.VisibilityTimeoutSeconds = &visibilityTimeoutSecondsFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.UpdateQueue(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.IntVar(&delaySecondsFlag, "delay-seconds", 0, "Delay seconds")
	f.IntVar(&kmsDataKeyReusePeriodSecondsFlag, "kms-data-key-reuse-period-seconds", 0, "Kms data key reuse period seconds")
	f.StringVar(&kmsKeyFlag, "kms-key", "", "KMS key UUID, CRN or account-scoped name")
	f.IntVar(&maxMessageBytesFlag, "max-message-bytes", 0, "Max message bytes")
	f.IntVar(&messageRetentionSecondsFlag, "message-retention-seconds", 0, "Message retention seconds")
	f.IntVar(&receiveWaitTimeSecondsFlag, "receive-wait-time-seconds", 0, "Receive wait time seconds")
	f.StringVar(&redrivePolicyFlag, "redrive-policy", "", "Omit to keep the current policy; null removes it (JSON)")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.IntVar(&visibilityTimeoutSecondsFlag, "visibility-timeout-seconds", 0, "Visibility timeout seconds")
	return cmd
}

// newQueueQueueDeleteCommand builds `basaltic queue queue delete`.
func newQueueQueueDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <queue-id>",
		Short: "Delete a queue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := queueClient(state, "/v1/queues/{queue_id}")
			if err != nil {
				return err
			}
			if err := c.DeleteQueue(cmd.Context(), args[0]); err != nil {
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

// newQueueQueuePurgeCommand builds `basaltic queue queue purge`.
func newQueueQueuePurgeCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purge <queue-id>",
		Short: "Delete every message on a queue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := queueClient(state, "/v1/queues/{queue_id}/purge")
			if err != nil {
				return err
			}
			out, err := c.PurgeQueue(cmd.Context(), args[0])
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
