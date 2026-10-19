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
	"github.com/basaltic-sh/sdk-go/notifications"

	"github.com/basaltic-sh/cli/internal/cli"
)

func init() { cli.RegisterService(newNotificationsCommand) }

// newNotificationsCommand builds `basaltic notifications`.
func newNotificationsCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "notifications",
		Short: "Publish/subscribe notification topics with email and webhook delivery",
		Long:  "Publish/subscribe notification topics with email and webhook delivery.\n\nThis is a regional service: it acts in the region from --region, the\nBASALTIC_REGION environment variable, or the profile.",
	}
	cmd.AddCommand(newNotificationsSubscriptionCommand(state))
	cmd.AddCommand(newNotificationsTopicCommand(state))
	return cmd
}

// notificationsClient builds the service client, resolving credentials on first use.
func notificationsClient(state *cli.State, path string) (*notifications.Client, error) {
	cfg, err := state.ServiceSDK("notifications", path)
	if err != nil {
		return nil, err
	}
	return notifications.New(cfg), nil
}

// newNotificationsSubscriptionCommand builds `basaltic notifications subscription`.
func newNotificationsSubscriptionCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "subscription",
		Short:   "Subscriptions",
		Aliases: []string{"subscriptions"},
	}
	cmd.AddCommand(newNotificationsSubscriptionListCommand(state))
	cmd.AddCommand(newNotificationsSubscriptionGetCommand(state))
	cmd.AddCommand(newNotificationsSubscriptionUpdateCommand(state))
	cmd.AddCommand(newNotificationsSubscriptionGetSigningSecretCommand(state))
	cmd.AddCommand(newNotificationsSubscriptionResendConfirmationCommand(state))
	cmd.AddCommand(newNotificationsSubscriptionRotateSigningSecretCommand(state))
	cmd.AddCommand(newNotificationsSubscriptionSubscribeCommand(state))
	cmd.AddCommand(newNotificationsSubscriptionUnsubscribeCommand(state))
	return cmd
}

// newNotificationsSubscriptionListCommand builds `basaltic notifications subscription list`.
func newNotificationsSubscriptionListCommand(state *cli.State) *cobra.Command {
	var params notifications.ListSubscriptionsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list <topic-id>",
		Short: "List subscriptions",
		Args:  cobra.ExactArgs(1),
		Long:  "List subscriptions.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := notificationsClient(state, "/v1/topics/{topic_id}/subscriptions")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListSubscriptionsAll(cmd.Context(), args[0], &params))
			}
			page, err := c.ListSubscriptions(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact nested topic/name/subscription/UUID CRN; intersects the parent topic and status")
	f.IntVar(&params.Limit, "limit", 0, "Maximum items to return")
	f.StringVar(&params.Marker, "marker", "", "Id of the last item on the previous page")
	f.StringVar(&params.Status, "status", "", "One of: \"pending_confirmation\", \"confirmed\"")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newNotificationsSubscriptionGetCommand builds `basaltic notifications subscription get`.
func newNotificationsSubscriptionGetCommand(state *cli.State) *cobra.Command {
	var scope notifications.ListSubscriptionsParams
	cmd := &cobra.Command{
		Use:   "get <topic-id> <ref>",
		Short: "Describe a subscription",
		Args:  cobra.ExactArgs(2),
		Long:  "Describe a subscription.\n\n<ref> is its id or its CRN. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.\n\nThis resource has no name.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := notificationsClient(state, "/v1/topics/{topic_id}/subscriptions/{subscription_id}")
			if err != nil {
				return err
			}
			out, err := c.DescribeSubscriptionByReference(cmd.Context(), args[0], args[1], &scope)
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

// newNotificationsSubscriptionUpdateCommand builds `basaltic notifications subscription update`.
func newNotificationsSubscriptionUpdateCommand(state *cli.State) *cobra.Command {
	var body notifications.UpdateSubscriptionRequest
	var bodyFile string
	var deliveryPolicyFlag string
	var filterPolicyFlag string
	var rawMessageDeliveryFlag bool
	var redrivePolicyFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <topic-id> <subscription-id>",
		Short: "Update mutable subscription settings",
		Args:  cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"delivery_policy":      "delivery-policy",
				"filter_policy":        "filter-policy",
				"raw_message_delivery": "raw-message-delivery",
				"redrive_policy":       "redrive-policy",
				"tags":                 "tags",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := notificationsClient(state, "/v1/topics/{topic_id}/subscriptions/{subscription_id}")
			if err != nil {
				return err
			}
			if deliveryPolicyFlag != "" {
				if err := json.Unmarshal([]byte(deliveryPolicyFlag), &body.DeliveryPolicy); err != nil {
					return fmt.Errorf("--delivery-policy: %w", err)
				}
			}
			if filterPolicyFlag != "" {
				if err := json.Unmarshal([]byte(filterPolicyFlag), &body.FilterPolicy); err != nil {
					return fmt.Errorf("--filter-policy: %w", err)
				}
			}
			if cmd.Flags().Changed("raw-message-delivery") {
				body.RawMessageDelivery = &rawMessageDeliveryFlag
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
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.UpdateSubscription(cmd.Context(), args[0], args[1], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&deliveryPolicyFlag, "delivery-policy", "", "Delivery policy (JSON)")
	f.StringVar(&filterPolicyFlag, "filter-policy", "", "Filter policy (JSON)")
	f.BoolVar(&rawMessageDeliveryFlag, "raw-message-delivery", false, "Raw message delivery")
	f.StringVar(&redrivePolicyFlag, "redrive-policy", "", "Omitted leaves the dead-letter queue unchanged (JSON)")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newNotificationsSubscriptionGetSigningSecretCommand builds `basaltic notifications subscription get-signing-secret`.
func newNotificationsSubscriptionGetSigningSecretCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-signing-secret <topic-id> <subscription-id>",
		Short: "Read the webhook signing secret",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := notificationsClient(state, "/v1/topics/{topic_id}/subscriptions/{subscription_id}/signing-secret")
			if err != nil {
				return err
			}
			out, err := c.GetSubscriptionSigningSecret(cmd.Context(), args[0], args[1])
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

// newNotificationsSubscriptionResendConfirmationCommand builds `basaltic notifications subscription resend-confirmation`.
func newNotificationsSubscriptionResendConfirmationCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resend-confirmation <topic-id> <subscription-id>",
		Short: "Re-send the confirmation handshake",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := notificationsClient(state, "/v1/topics/{topic_id}/subscriptions/{subscription_id}/resend-confirmation")
			if err != nil {
				return err
			}
			out, err := c.ResendConfirmation(cmd.Context(), args[0], args[1])
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

// newNotificationsSubscriptionRotateSigningSecretCommand builds `basaltic notifications subscription rotate-signing-secret`.
func newNotificationsSubscriptionRotateSigningSecretCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rotate-signing-secret <topic-id> <subscription-id>",
		Short: "Rotate the webhook signing secret",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := notificationsClient(state, "/v1/topics/{topic_id}/subscriptions/{subscription_id}/signing-secret/rotate")
			if err != nil {
				return err
			}
			out, err := c.RotateSubscriptionSigningSecret(cmd.Context(), args[0], args[1])
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

// newNotificationsSubscriptionSubscribeCommand builds `basaltic notifications subscription subscribe`.
func newNotificationsSubscriptionSubscribeCommand(state *cli.State) *cobra.Command {
	var body notifications.SubscribeRequest
	var bodyFile string
	var deliveryPolicyFlag string
	var filterPolicyFlag string
	var rawMessageDeliveryFlag bool
	var redrivePolicyFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "subscribe <topic-id>",
		Short: "Subscribe an endpoint to a topic",
		Args:  cobra.ExactArgs(1),
		Long:  "Subscribe an endpoint to a topic.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"delivery_policy":      "delivery-policy",
				"endpoint":             "endpoint",
				"filter_policy":        "filter-policy",
				"protocol":             "protocol",
				"raw_message_delivery": "raw-message-delivery",
				"redrive_policy":       "redrive-policy",
				"tags":                 "tags",
			}, []string{"endpoint", "protocol"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := notificationsClient(state, "/v1/topics/{topic_id}/subscriptions")
			if err != nil {
				return err
			}
			if deliveryPolicyFlag != "" {
				if err := json.Unmarshal([]byte(deliveryPolicyFlag), &body.DeliveryPolicy); err != nil {
					return fmt.Errorf("--delivery-policy: %w", err)
				}
			}
			if filterPolicyFlag != "" {
				if err := json.Unmarshal([]byte(filterPolicyFlag), &body.FilterPolicy); err != nil {
					return fmt.Errorf("--filter-policy: %w", err)
				}
			}
			if cmd.Flags().Changed("raw-message-delivery") {
				body.RawMessageDelivery = &rawMessageDeliveryFlag
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
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.Subscribe(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&deliveryPolicyFlag, "delivery-policy", "", "Delivery policy (JSON)")
	f.StringVar(&body.Endpoint, "endpoint", "", "An email address; an http/https URL whose scheme matches the protocol and resolves to a publicly-routable address; or, for queue, a queue UUID, CRN or exact name in the subscriber account and serving region")
	f.StringVar(&filterPolicyFlag, "filter-policy", "", "Filter policy (JSON)")
	f.StringVar(&body.Protocol, "protocol", "", "Protocol (one of: email, https, http, queue)")
	f.BoolVar(&rawMessageDeliveryFlag, "raw-message-delivery", false, "http/https only")
	f.StringVar(&redrivePolicyFlag, "redrive-policy", "", "Redrive policy (JSON)")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newNotificationsSubscriptionUnsubscribeCommand builds `basaltic notifications subscription unsubscribe`.
func newNotificationsSubscriptionUnsubscribeCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unsubscribe <topic-id> <subscription-id>",
		Short: "Unsubscribe",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := notificationsClient(state, "/v1/topics/{topic_id}/subscriptions/{subscription_id}")
			if err != nil {
				return err
			}
			if err := c.Unsubscribe(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			state.Printer().Done("Unsubscribe requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newNotificationsTopicCommand builds `basaltic notifications topic`.
func newNotificationsTopicCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "topic",
		Short:   "Topics",
		Aliases: []string{"topics"},
	}
	cmd.AddCommand(newNotificationsTopicListCommand(state))
	cmd.AddCommand(newNotificationsTopicGetCommand(state))
	cmd.AddCommand(newNotificationsTopicCreateCommand(state))
	cmd.AddCommand(newNotificationsTopicUpdateCommand(state))
	cmd.AddCommand(newNotificationsTopicDeleteCommand(state))
	cmd.AddCommand(newNotificationsTopicDeletePolicyCommand(state))
	cmd.AddCommand(newNotificationsTopicGetMessageCommand(state))
	cmd.AddCommand(newNotificationsTopicGetPolicyCommand(state))
	cmd.AddCommand(newNotificationsTopicListDeliveriesCommand(state))
	cmd.AddCommand(newNotificationsTopicListMessagesCommand(state))
	cmd.AddCommand(newNotificationsTopicPublishCommand(state))
	cmd.AddCommand(newNotificationsTopicPublishBatchCommand(state))
	cmd.AddCommand(newNotificationsTopicSetPolicyCommand(state))
	return cmd
}

// newNotificationsTopicListCommand builds `basaltic notifications topic list`.
func newNotificationsTopicListCommand(state *cli.State) *cobra.Command {
	var params notifications.ListTopicsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List topics",
		Args:  cobra.ExactArgs(0),
		Long:  "List topics.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := notificationsClient(state, "/v1/topics")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListTopicsAll(cmd.Context(), &params))
			}
			page, err := c.ListTopics(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact notifications/topic CRN; foreign or mismatched CRNs return an empty page")
	f.IntVar(&params.Limit, "limit", 0, "Maximum items to return")
	f.StringVar(&params.Marker, "marker", "", "Id of the last item on the previous page")
	f.StringVar(&params.Name, "name", "", "Exact immutable topic name in the caller account; intersects crn")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newNotificationsTopicGetCommand builds `basaltic notifications topic get`.
func newNotificationsTopicGetCommand(state *cli.State) *cobra.Command {
	var scope notifications.ListTopicsParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Describe a topic",
		Args:  cobra.ExactArgs(1),
		Long:  "Describe a topic.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := notificationsClient(state, "/v1/topics/{topic_id}")
			if err != nil {
				return err
			}
			out, err := c.DescribeTopicByReference(cmd.Context(), args[0], &scope)
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

// newNotificationsTopicCreateCommand builds `basaltic notifications topic create`.
func newNotificationsTopicCreateCommand(state *cli.State) *cobra.Command {
	var body notifications.CreateTopicRequest
	var bodyFile string
	var deliveryPolicyFlag string
	var displayNameFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a topic",
		Args:  cobra.ExactArgs(0),
		Long:  "Create a topic.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"delivery_policy": "delivery-policy",
				"display_name":    "display-name",
				"name":            "name",
				"tags":            "tags",
			}, []string{"name"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := notificationsClient(state, "/v1/topics")
			if err != nil {
				return err
			}
			if deliveryPolicyFlag != "" {
				if err := json.Unmarshal([]byte(deliveryPolicyFlag), &body.DeliveryPolicy); err != nil {
					return fmt.Errorf("--delivery-policy: %w", err)
				}
			}
			if cmd.Flags().Changed("display-name") {
				body.DisplayName = &displayNameFlag
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
			out, err := c.CreateTopic(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&deliveryPolicyFlag, "delivery-policy", "", "Delivery policy (JSON)")
	f.StringVar(&displayNameFlag, "display-name", "", "Display name")
	f.StringVar(&body.Name, "name", "", "1-256 characters of letters, digits, dot, dash, or underscore")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newNotificationsTopicUpdateCommand builds `basaltic notifications topic update`.
func newNotificationsTopicUpdateCommand(state *cli.State) *cobra.Command {
	var body notifications.UpdateTopicRequest
	var bodyFile string
	var deliveryPolicyFlag string
	var displayNameFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <topic-id>",
		Short: "Update mutable topic metadata",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"delivery_policy": "delivery-policy",
				"display_name":    "display-name",
				"tags":            "tags",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := notificationsClient(state, "/v1/topics/{topic_id}")
			if err != nil {
				return err
			}
			if deliveryPolicyFlag != "" {
				if err := json.Unmarshal([]byte(deliveryPolicyFlag), &body.DeliveryPolicy); err != nil {
					return fmt.Errorf("--delivery-policy: %w", err)
				}
			}
			if cmd.Flags().Changed("display-name") {
				body.DisplayName = &displayNameFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.UpdateTopic(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&deliveryPolicyFlag, "delivery-policy", "", "Delivery policy (JSON)")
	f.StringVar(&displayNameFlag, "display-name", "", "Display name")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newNotificationsTopicDeleteCommand builds `basaltic notifications topic delete`.
func newNotificationsTopicDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <topic-id>",
		Short: "Delete a topic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := notificationsClient(state, "/v1/topics/{topic_id}")
			if err != nil {
				return err
			}
			if err := c.DeleteTopic(cmd.Context(), args[0]); err != nil {
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

// newNotificationsTopicDeletePolicyCommand builds `basaltic notifications topic delete-policy`.
func newNotificationsTopicDeletePolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-policy <topic-id>",
		Short: "Delete the topic's resource policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := notificationsClient(state, "/v1/topics/{topic_id}/policy")
			if err != nil {
				return err
			}
			if err := c.DeleteTopicPolicy(cmd.Context(), args[0]); err != nil {
				return err
			}
			state.Printer().Done("Delete policy requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newNotificationsTopicGetMessageCommand builds `basaltic notifications topic get-message`.
func newNotificationsTopicGetMessageCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-message <topic-id> <message-id>",
		Short: "Describe a published notification",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := notificationsClient(state, "/v1/topics/{topic_id}/messages/{message_id}")
			if err != nil {
				return err
			}
			out, err := c.DescribeMessage(cmd.Context(), args[0], args[1])
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

// newNotificationsTopicGetPolicyCommand builds `basaltic notifications topic get-policy`.
func newNotificationsTopicGetPolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-policy <topic-id>",
		Short: "Read the topic's resource policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := notificationsClient(state, "/v1/topics/{topic_id}/policy")
			if err != nil {
				return err
			}
			out, err := c.GetTopicPolicy(cmd.Context(), args[0])
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

// newNotificationsTopicListDeliveriesCommand builds `basaltic notifications topic list-deliveries`.
func newNotificationsTopicListDeliveriesCommand(state *cli.State) *cobra.Command {
	var params notifications.ListDeliveriesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list-deliveries <topic-id>",
		Short: "List delivery records",
		Args:  cobra.ExactArgs(1),
		Long:  "List delivery records.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := notificationsClient(state, "/v1/topics/{topic_id}/deliveries")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListDeliveriesAll(cmd.Context(), args[0], &params))
			}
			page, err := c.ListDeliveries(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact message/UUID/delivery/UUID identity")
	f.IntVar(&params.Limit, "limit", 0, "Maximum items to return")
	f.StringVar(&params.Marker, "marker", "", "Id of the last item on the previous page")
	f.StringVar(&params.Message, "message", "", "Restrict to one published message")
	f.StringVar(&params.Status, "status", "", "One of: \"pending\", \"delivered\", \"failed\"")
	f.StringVar(&params.Subscription, "subscription", "", "Restrict to one subscription by UUID or nested topic/name/subscription/UUID CRN")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newNotificationsTopicListMessagesCommand builds `basaltic notifications topic list-messages`.
func newNotificationsTopicListMessagesCommand(state *cli.State) *cobra.Command {
	var params notifications.ListMessagesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list-messages <topic-id>",
		Short: "List messages",
		Args:  cobra.ExactArgs(1),
		Long:  "List messages.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := notificationsClient(state, "/v1/topics/{topic_id}/messages")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListMessagesAll(cmd.Context(), args[0], &params))
			}
			page, err := c.ListMessages(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact message CRN")
	f.StringVar(&params.FanoutState, "fanout-state", "", "One of: \"pending\", \"complete\"")
	f.IntVar(&params.Limit, "limit", 0, "Maximum items to return")
	f.StringVar(&params.Marker, "marker", "", "Id of the last item on the previous page")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newNotificationsTopicPublishCommand builds `basaltic notifications topic publish`.
func newNotificationsTopicPublishCommand(state *cli.State) *cobra.Command {
	var body notifications.PublishRequest
	var bodyFile string
	var attributesFlag string
	var messageStructureFlag string
	var subjectFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "publish <topic-id>",
		Short: "Publish a notification",
		Args:  cobra.ExactArgs(1),
		Long:  "Publish a notification.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"attributes":        "attributes",
				"message":           "message",
				"message_structure": "message-structure",
				"subject":           "subject",
			}, []string{"message"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := notificationsClient(state, "/v1/topics/{topic_id}/publish")
			if err != nil {
				return err
			}
			if attributesFlag != "" {
				if err := json.Unmarshal([]byte(attributesFlag), &body.Attributes); err != nil {
					return fmt.Errorf("--attributes: %w", err)
				}
			}
			if cmd.Flags().Changed("message-structure") {
				body.MessageStructure = &messageStructureFlag
			}
			if cmd.Flags().Changed("subject") {
				body.Subject = &subjectFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.Publish(cmd.Context(), args[0], &body, reqOpts...)
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
	f.StringVar(&body.Message, "message", "", "The payload")
	f.StringVar(&messageStructureFlag, "message-structure", "", "Omit to deliver message verbatim to every protocol (one of: , json)")
	f.StringVar(&subjectFlag, "subject", "", "Used as the email subject")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newNotificationsTopicPublishBatchCommand builds `basaltic notifications topic publish-batch`.
func newNotificationsTopicPublishBatchCommand(state *cli.State) *cobra.Command {
	var body notifications.PublishBatchRequest
	var bodyFile string
	var entriesFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "publish-batch <topic-id>",
		Short: "Publish up to 10 notifications",
		Args:  cobra.ExactArgs(1),
		Long:  "Publish up to 10 notifications.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"entries": "entries",
			}, []string{"entries"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := notificationsClient(state, "/v1/topics/{topic_id}/publish-batch")
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
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.PublishBatch(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&entriesFlag, "entries", "", "Entries (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newNotificationsTopicSetPolicyCommand builds `basaltic notifications topic set-policy`.
func newNotificationsTopicSetPolicyCommand(state *cli.State) *cobra.Command {
	var body notifications.PutTopicPolicyRequest
	var bodyFile string
	var documentFlag string
	cmd := &cobra.Command{
		Use:   "set-policy <topic-id>",
		Short: "Set the topic's resource policy",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"document": "document",
			}, []string{"document"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := notificationsClient(state, "/v1/topics/{topic_id}/policy")
			if err != nil {
				return err
			}
			if documentFlag != "" {
				if err := json.Unmarshal([]byte(documentFlag), &body.Document); err != nil {
					return fmt.Errorf("--document: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.PutTopicPolicy(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&documentFlag, "document", "", "Document (JSON)")
	return cmd
}
