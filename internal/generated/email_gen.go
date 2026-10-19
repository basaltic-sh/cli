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
	"github.com/basaltic-sh/sdk-go/email"

	"github.com/basaltic-sh/cli/internal/cli"
)

func init() { cli.RegisterService(newEmailCommand) }

// newEmailCommand builds `basaltic email`.
func newEmailCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "email",
		Short: "Regional API for sending and receiving email, identities and DNS records, templates, configuration sets, IP pools, recipient rules, message content and attachments",
		Long:  "Regional API for sending and receiving email, identities and DNS records, templates, configuration sets, IP pools, recipient rules, message content and attachments.\n\nThis is a regional service: it acts in the region from --region, the\nBASALTIC_REGION environment variable, or the profile.",
	}
	cmd.AddCommand(newEmailConfigurationSetCommand(state))
	cmd.AddCommand(newEmailDedicatedIpCommand(state))
	cmd.AddCommand(newEmailIdentityCommand(state))
	cmd.AddCommand(newEmailIpPoolCommand(state))
	cmd.AddCommand(newEmailReceivedMessageCommand(state))
	cmd.AddCommand(newEmailRecipientRuleCommand(state))
	cmd.AddCommand(newEmailRejectionCommand(state))
	cmd.AddCommand(newEmailSentMessageCommand(state))
	cmd.AddCommand(newEmailSuppressionCommand(state))
	cmd.AddCommand(newEmailTemplateCommand(state))
	return cmd
}

// emailClient builds the service client, resolving credentials on first use.
func emailClient(state *cli.State, path string) (*email.Client, error) {
	cfg, err := state.ServiceSDK("email", path)
	if err != nil {
		return nil, err
	}
	return email.New(cfg), nil
}

// newEmailConfigurationSetCommand builds `basaltic email configuration-set`.
func newEmailConfigurationSetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "configuration-set",
		Short:   "Configuration sets",
		Aliases: []string{"configuration-sets"},
	}
	cmd.AddCommand(newEmailConfigurationSetListCommand(state))
	cmd.AddCommand(newEmailConfigurationSetGetCommand(state))
	cmd.AddCommand(newEmailConfigurationSetCreateCommand(state))
	cmd.AddCommand(newEmailConfigurationSetUpdateCommand(state))
	cmd.AddCommand(newEmailConfigurationSetDeleteCommand(state))
	cmd.AddCommand(newEmailConfigurationSetAddSuppressionOverrideCommand(state))
	cmd.AddCommand(newEmailConfigurationSetRemoveSuppressionOverrideCommand(state))
	return cmd
}

// newEmailConfigurationSetListCommand builds `basaltic email configuration-set list`.
func newEmailConfigurationSetListCommand(state *cli.State) *cobra.Command {
	var params email.ListConfigurationSetsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List configuration sets",
		Args:  cobra.ExactArgs(0),
		Long:  "List configuration sets.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/configuration-sets")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListConfigurationSetsAll(cmd.Context(), &params))
			}
			page, err := c.ListConfigurationSets(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact CRN filter")
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Marker")
	f.StringVar(&params.Name, "name", "", "Exact immutable account-scoped name, combined with all other filters")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newEmailConfigurationSetGetCommand builds `basaltic email configuration-set get`.
func newEmailConfigurationSetGetCommand(state *cli.State) *cobra.Command {
	var scope email.ListConfigurationSetsParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Describe a configuration set",
		Args:  cobra.ExactArgs(1),
		Long:  "Describe a configuration set.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/configuration-sets/{configuration_set_id}")
			if err != nil {
				return err
			}
			out, err := c.DescribeConfigurationSetByReference(cmd.Context(), args[0], &scope)
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

// newEmailConfigurationSetCreateCommand builds `basaltic email configuration-set create`.
func newEmailConfigurationSetCreateCommand(state *cli.State) *cobra.Command {
	var body email.CreateConfigurationSetRequest
	var bodyFile string
	var ipPoolFlag string
	var sendingEnabledFlag bool
	var tagsFlag string
	var trackingClickFlag bool
	var trackingOpenFlag bool
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a configuration set",
		Args:  cobra.ExactArgs(0),
		Long:  "Create a configuration set.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"ip_pool":               "ip-pool",
				"name":                  "name",
				"sending_enabled":       "sending-enabled",
				"suppression_overrides": "suppression-overrides",
				"tags":                  "tags",
				"tracking_click":        "tracking-click",
				"tracking_open":         "tracking-open",
			}, []string{"name"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/configuration-sets")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("ip-pool") {
				body.IPPool = &ipPoolFlag
			}
			if cmd.Flags().Changed("sending-enabled") {
				body.SendingEnabled = &sendingEnabledFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			if cmd.Flags().Changed("tracking-click") {
				body.TrackingClick = &trackingClickFlag
			}
			if cmd.Flags().Changed("tracking-open") {
				body.TrackingOpen = &trackingOpenFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateConfigurationSet(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&ipPoolFlag, "ip-pool", "", "UUID, CRN or account-unique name")
	f.StringVar(&body.Name, "name", "", "Resource names must not start with the literal crn: prefix or be UUIDs (canonical, compact, braced, or urn:uuid: forms, in either case)")
	f.BoolVar(&sendingEnabledFlag, "sending-enabled", false, "Sending enabled")
	f.StringSliceVar(&body.SuppressionOverrides, "suppression-overrides", nil, "Suppression overrides")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.BoolVar(&trackingClickFlag, "tracking-click", false, "Tracking click")
	f.BoolVar(&trackingOpenFlag, "tracking-open", false, "Tracking open")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newEmailConfigurationSetUpdateCommand builds `basaltic email configuration-set update`.
func newEmailConfigurationSetUpdateCommand(state *cli.State) *cobra.Command {
	var body email.UpdateConfigurationSetRequest
	var bodyFile string
	var ipPoolFlag string
	var sendingEnabledFlag bool
	var tagsFlag string
	var trackingClickFlag bool
	var trackingOpenFlag bool
	cmd := &cobra.Command{
		Use:   "update <configuration-set-id>",
		Short: "Update a configuration set",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"ip_pool":               "ip-pool",
				"sending_enabled":       "sending-enabled",
				"suppression_overrides": "suppression-overrides",
				"tags":                  "tags",
				"tracking_click":        "tracking-click",
				"tracking_open":         "tracking-open",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/configuration-sets/{configuration_set_id}")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("ip-pool") {
				body.IPPool = &ipPoolFlag
			}
			if cmd.Flags().Changed("sending-enabled") {
				body.SendingEnabled = &sendingEnabledFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			if cmd.Flags().Changed("tracking-click") {
				body.TrackingClick = &trackingClickFlag
			}
			if cmd.Flags().Changed("tracking-open") {
				body.TrackingOpen = &trackingOpenFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.UpdateConfigurationSet(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&ipPoolFlag, "ip-pool", "", "UUID, CRN or account-unique name")
	f.BoolVar(&sendingEnabledFlag, "sending-enabled", false, "Sending enabled")
	f.StringSliceVar(&body.SuppressionOverrides, "suppression-overrides", nil, "Suppression overrides")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.BoolVar(&trackingClickFlag, "tracking-click", false, "Tracking click")
	f.BoolVar(&trackingOpenFlag, "tracking-open", false, "Tracking open")
	return cmd
}

// newEmailConfigurationSetDeleteCommand builds `basaltic email configuration-set delete`.
func newEmailConfigurationSetDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <configuration-set-id>",
		Short: "Delete a configuration set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/configuration-sets/{configuration_set_id}")
			if err != nil {
				return err
			}
			if err := c.DeleteConfigurationSet(cmd.Context(), args[0]); err != nil {
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

// newEmailConfigurationSetAddSuppressionOverrideCommand builds `basaltic email configuration-set add-suppression-override`.
func newEmailConfigurationSetAddSuppressionOverrideCommand(state *cli.State) *cobra.Command {
	var body email.SuppressionOverrideRequest
	var bodyFile string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "add-suppression-override <configuration-set-id>",
		Short: "Add a suppression override address",
		Args:  cobra.ExactArgs(1),
		Long:  "Add a suppression override address.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"address": "address",
			}, []string{"address"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/configuration-sets/{configuration_set_id}/suppression-overrides")
			if err != nil {
				return err
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.AddSuppressionOverride(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.Address, "address", "", "Address")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newEmailConfigurationSetRemoveSuppressionOverrideCommand builds `basaltic email configuration-set remove-suppression-override`.
func newEmailConfigurationSetRemoveSuppressionOverrideCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-suppression-override <configuration-set-id> <address>",
		Short: "Remove a suppression override address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/configuration-sets/{configuration_set_id}/suppression-overrides/{address}")
			if err != nil {
				return err
			}
			out, err := c.RemoveSuppressionOverride(cmd.Context(), args[0], args[1])
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

// newEmailDedicatedIpCommand builds `basaltic email dedicated-ip`.
func newEmailDedicatedIpCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dedicated-ip",
		Short:   "Dedicated ips",
		Aliases: []string{"dedicated-ips"},
	}
	cmd.AddCommand(newEmailDedicatedIpListCommand(state))
	cmd.AddCommand(newEmailDedicatedIpGetCommand(state))
	cmd.AddCommand(newEmailDedicatedIpUpdateCommand(state))
	cmd.AddCommand(newEmailDedicatedIpReleaseCommand(state))
	cmd.AddCommand(newEmailDedicatedIpRequestCommand(state))
	return cmd
}

// newEmailDedicatedIpListCommand builds `basaltic email dedicated-ip list`.
func newEmailDedicatedIpListCommand(state *cli.State) *cobra.Command {
	var params email.ListDedicatedIPsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List dedicated IPs",
		Args:  cobra.ExactArgs(0),
		Long:  "List dedicated IPs.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/dedicated-ips")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListDedicatedIPsAll(cmd.Context(), &params))
			}
			page, err := c.ListDedicatedIPs(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact dedicated-IP CRN, filtered before pagination")
	f.StringVar(&params.IPPool, "ip-pool", "", "Ip pool")
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Marker")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newEmailDedicatedIpGetCommand builds `basaltic email dedicated-ip get`.
func newEmailDedicatedIpGetCommand(state *cli.State) *cobra.Command {
	var scope email.ListDedicatedIPsParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Describe a dedicated IP",
		Args:  cobra.ExactArgs(1),
		Long:  "Describe a dedicated IP.\n\n<ref> is its id or its CRN. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.\n\nThis resource has no name.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/dedicated-ips/{dedicated_ip_id}")
			if err != nil {
				return err
			}
			out, err := c.DescribeDedicatedIPByReference(cmd.Context(), args[0], &scope)
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

// newEmailDedicatedIpUpdateCommand builds `basaltic email dedicated-ip update`.
func newEmailDedicatedIpUpdateCommand(state *cli.State) *cobra.Command {
	var body email.UpdateDedicatedIPRequest
	var bodyFile string
	var ipPoolFlag string
	var warmupStatusFlag string
	cmd := &cobra.Command{
		Use:   "update <dedicated-ip-id>",
		Short: "Update a dedicated IP",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"ip_pool":       "ip-pool",
				"warmup_status": "warmup-status",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/dedicated-ips/{dedicated_ip_id}")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("ip-pool") {
				body.IPPool = &ipPoolFlag
			}
			if cmd.Flags().Changed("warmup-status") {
				body.WarmupStatus = &warmupStatusFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.UpdateDedicatedIP(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&ipPoolFlag, "ip-pool", "", "UUID, CRN or account-unique pool name")
	f.StringVar(&warmupStatusFlag, "warmup-status", "", "Warmup status (one of: new, warming, warm)")
	return cmd
}

// newEmailDedicatedIpReleaseCommand builds `basaltic email dedicated-ip release`.
func newEmailDedicatedIpReleaseCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release <dedicated-ip-id>",
		Short: "Release a dedicated IP",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/dedicated-ips/{dedicated_ip_id}")
			if err != nil {
				return err
			}
			if err := c.ReleaseDedicatedIP(cmd.Context(), args[0]); err != nil {
				return err
			}
			state.Printer().Done("Release requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newEmailDedicatedIpRequestCommand builds `basaltic email dedicated-ip request`.
func newEmailDedicatedIpRequestCommand(state *cli.State) *cobra.Command {
	var body email.RequestDedicatedIPRequest
	var bodyFile string
	var ipPoolFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "request",
		Short: "Request a dedicated IP",
		Args:  cobra.ExactArgs(0),
		Long:  "Request a dedicated IP.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"ip_pool": "ip-pool",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/dedicated-ips")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("ip-pool") {
				body.IPPool = &ipPoolFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.RequestDedicatedIP(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&ipPoolFlag, "ip-pool", "", "UUID, CRN or account-unique pool name; empty leaves unassigned")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newEmailIdentityCommand builds `basaltic email identity`.
func newEmailIdentityCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "identity",
		Short:   "Identities",
		Aliases: []string{"identities"},
	}
	cmd.AddCommand(newEmailIdentityListCommand(state))
	cmd.AddCommand(newEmailIdentityGetCommand(state))
	cmd.AddCommand(newEmailIdentityCreateCommand(state))
	cmd.AddCommand(newEmailIdentityUpdateCommand(state))
	cmd.AddCommand(newEmailIdentityDeleteCommand(state))
	cmd.AddCommand(newEmailIdentityBeginDkimCommand(state))
	cmd.AddCommand(newEmailIdentityCancelDkimCommand(state))
	cmd.AddCommand(newEmailIdentityFinalizeDkimCommand(state))
	cmd.AddCommand(newEmailIdentityGetReceivingCommand(state))
	cmd.AddCommand(newEmailIdentityListDmarcReportsCommand(state))
	cmd.AddCommand(newEmailIdentitySetBimiCommand(state))
	cmd.AddCommand(newEmailIdentitySetMailFromDomainCommand(state))
	cmd.AddCommand(newEmailIdentityUpdateReceivingCommand(state))
	cmd.AddCommand(newEmailIdentityVerifyCommand(state))
	cmd.AddCommand(newEmailIdentityVerifyBimiCommand(state))
	cmd.AddCommand(newEmailIdentityVerifyMailFromDomainCommand(state))
	return cmd
}

// newEmailIdentityListCommand builds `basaltic email identity list`.
func newEmailIdentityListCommand(state *cli.State) *cobra.Command {
	var params email.ListIdentitiesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List sender identities",
		Args:  cobra.ExactArgs(0),
		Long:  "List sender identities.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/identities")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListIdentitiesAll(cmd.Context(), &params))
			}
			page, err := c.ListIdentities(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact CRN filter")
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Marker")
	f.StringVar(&params.Name, "name", "", "Exact immutable account-scoped name, combined with all other filters")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newEmailIdentityGetCommand builds `basaltic email identity get`.
func newEmailIdentityGetCommand(state *cli.State) *cobra.Command {
	var scope email.ListIdentitiesParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Describe an identity",
		Args:  cobra.ExactArgs(1),
		Long:  "Describe an identity.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/identities/{identity_id}")
			if err != nil {
				return err
			}
			out, err := c.DescribeIdentityByReference(cmd.Context(), args[0], &scope)
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

// newEmailIdentityCreateCommand builds `basaltic email identity create`.
func newEmailIdentityCreateCommand(state *cli.State) *cobra.Command {
	var body email.CreateIdentityRequest
	var bodyFile string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a sender identity",
		Args:  cobra.ExactArgs(0),
		Long:  "Create a sender identity.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"name": "name",
				"tags": "tags",
			}, []string{"name"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/identities")
			if err != nil {
				return err
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
			out, err := c.CreateIdentity(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.Name, "name", "", "Sender domain (FQDN, lowercase)")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newEmailIdentityUpdateCommand builds `basaltic email identity update`.
func newEmailIdentityUpdateCommand(state *cli.State) *cobra.Command {
	var body email.UpdateIdentityRequest
	var bodyFile string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <identity-id>",
		Short: "Update identity",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"tags": "tags",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/identities/{identity_id}")
			if err != nil {
				return err
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.UpdateIdentity(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newEmailIdentityDeleteCommand builds `basaltic email identity delete`.
func newEmailIdentityDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <identity-id>",
		Short: "Delete an identity",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/identities/{identity_id}")
			if err != nil {
				return err
			}
			if err := c.DeleteIdentity(cmd.Context(), args[0]); err != nil {
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

// newEmailIdentityBeginDkimCommand builds `basaltic email identity begin-dkim`.
func newEmailIdentityBeginDkimCommand(state *cli.State) *cobra.Command {
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "begin-dkim <identity-id>",
		Short: "Begin DKIM key rotation",
		Args:  cobra.ExactArgs(1),
		Long:  "Begin DKIM key rotation.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/identities/{identity_id}/dkim/rotate")
			if err != nil {
				return err
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.BeginDkimRotation(cmd.Context(), args[0], reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newEmailIdentityCancelDkimCommand builds `basaltic email identity cancel-dkim`.
func newEmailIdentityCancelDkimCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-dkim <identity-id>",
		Short: "Cancel in-progress DKIM rotation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/identities/{identity_id}/dkim/rotate")
			if err != nil {
				return err
			}
			out, err := c.CancelDkimRotation(cmd.Context(), args[0])
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

// newEmailIdentityFinalizeDkimCommand builds `basaltic email identity finalize-dkim`.
func newEmailIdentityFinalizeDkimCommand(state *cli.State) *cobra.Command {
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "finalize-dkim <identity-id>",
		Short: "Finalize DKIM key rotation",
		Args:  cobra.ExactArgs(1),
		Long:  "Finalize DKIM key rotation.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/identities/{identity_id}/dkim/rotate/finalize")
			if err != nil {
				return err
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.FinalizeDkimRotation(cmd.Context(), args[0], reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newEmailIdentityGetReceivingCommand builds `basaltic email identity get-receiving`.
func newEmailIdentityGetReceivingCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-receiving <identity-id>",
		Short: "Get receiving configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/identities/{identity_id}/receiving")
			if err != nil {
				return err
			}
			out, err := c.GetReceivingConfiguration(cmd.Context(), args[0])
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

// newEmailIdentityListDmarcReportsCommand builds `basaltic email identity list-dmarc-reports`.
func newEmailIdentityListDmarcReportsCommand(state *cli.State) *cobra.Command {
	var params email.ListDmarcReportsParams
	var fromFlag string
	var toFlag string
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list-dmarc-reports <identity-id>",
		Short: "List DMARC report records",
		Args:  cobra.ExactArgs(1),
		Long:  "List DMARC report records.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/identities/{identity_id}/dmarc-reports")
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
				return state.Printer().Iter(c.ListDmarcReportsAll(cmd.Context(), args[0], &params))
			}
			page, err := c.ListDmarcReports(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&fromFlag, "from", "", "From (RFC 3339)")
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Marker")
	f.StringVar(&toFlag, "to", "", "To (RFC 3339)")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newEmailIdentitySetBimiCommand builds `basaltic email identity set-bimi`.
func newEmailIdentitySetBimiCommand(state *cli.State) *cobra.Command {
	var body email.SetBimiRequest
	var bodyFile string
	var logoUrlFlag string
	var vmcUrlFlag string
	cmd := &cobra.Command{
		Use:   "set-bimi <identity-id>",
		Short: "Set BIMI logo / VMC URLs",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"logo_url": "logo-url",
				"vmc_url":  "vmc-url",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/identities/{identity_id}/bimi")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("logo-url") {
				body.LogoURL = &logoUrlFlag
			}
			if cmd.Flags().Changed("vmc-url") {
				body.VmcURL = &vmcUrlFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.SetBimi(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&logoUrlFlag, "logo-url", "", "HTTPS SVG Tiny PS URL; null clears BIMI")
	f.StringVar(&vmcUrlFlag, "vmc-url", "", "Optional HTTPS VMC URL; omission preserves it and null removes the certificate")
	return cmd
}

// newEmailIdentitySetMailFromDomainCommand builds `basaltic email identity set-mail-from-domain`.
func newEmailIdentitySetMailFromDomainCommand(state *cli.State) *cobra.Command {
	var body email.SetMailFromDomainRequest
	var bodyFile string
	var domainFlag string
	cmd := &cobra.Command{
		Use:   "set-mail-from-domain <identity-id>",
		Short: "Set custom MAIL FROM domain",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"domain": "domain",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/identities/{identity_id}/mail-from-domain")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("domain") {
				body.Domain = &domainFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.SetMailFromDomain(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&domainFlag, "domain", "", "Mail-from domain; null clears the binding")
	return cmd
}

// newEmailIdentityUpdateReceivingCommand builds `basaltic email identity update-receiving`.
func newEmailIdentityUpdateReceivingCommand(state *cli.State) *cobra.Command {
	var body email.UpdateReceivingConfigurationRequest
	var bodyFile string
	var attachmentsBucketFlag string
	var enabledFlag bool
	cmd := &cobra.Command{
		Use:   "update-receiving <identity-id>",
		Short: "Update receiving configuration",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"attachments_bucket": "attachments-bucket",
				"enabled":            "enabled",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/identities/{identity_id}/receiving")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("attachments-bucket") {
				body.AttachmentsBucket = &attachmentsBucketFlag
			}
			if cmd.Flags().Changed("enabled") {
				body.Enabled = &enabledFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			if err := c.UpdateReceivingConfiguration(cmd.Context(), args[0], &body, reqOpts...); err != nil {
				return err
			}
			state.Printer().Done("Update receiving requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&attachmentsBucketFlag, "attachments-bucket", "", "Same-account, same-region bucket name, UUID or CRN")
	f.BoolVar(&enabledFlag, "enabled", false, "Enabled")
	return cmd
}

// newEmailIdentityVerifyCommand builds `basaltic email identity verify`.
func newEmailIdentityVerifyCommand(state *cli.State) *cobra.Command {
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "verify <identity-id>",
		Short: "Trigger DNS verification",
		Args:  cobra.ExactArgs(1),
		Long:  "Trigger DNS verification.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/identities/{identity_id}/verify")
			if err != nil {
				return err
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.VerifyIdentity(cmd.Context(), args[0], reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newEmailIdentityVerifyBimiCommand builds `basaltic email identity verify-bimi`.
func newEmailIdentityVerifyBimiCommand(state *cli.State) *cobra.Command {
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "verify-bimi <identity-id>",
		Short: "Verify BIMI logo (+ optional VMC)",
		Args:  cobra.ExactArgs(1),
		Long:  "Verify BIMI logo (+ optional VMC).\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/identities/{identity_id}/bimi/verify")
			if err != nil {
				return err
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.VerifyBimi(cmd.Context(), args[0], reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newEmailIdentityVerifyMailFromDomainCommand builds `basaltic email identity verify-mail-from-domain`.
func newEmailIdentityVerifyMailFromDomainCommand(state *cli.State) *cobra.Command {
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "verify-mail-from-domain <identity-id>",
		Short: "Verify custom MAIL FROM domain MX",
		Args:  cobra.ExactArgs(1),
		Long:  "Verify custom MAIL FROM domain MX.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/identities/{identity_id}/mail-from-domain/verify")
			if err != nil {
				return err
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.VerifyMailFromDomain(cmd.Context(), args[0], reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newEmailIpPoolCommand builds `basaltic email ip-pool`.
func newEmailIpPoolCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ip-pool",
		Short:   "Ip pools",
		Aliases: []string{"ip-pools"},
	}
	cmd.AddCommand(newEmailIpPoolListCommand(state))
	cmd.AddCommand(newEmailIpPoolGetCommand(state))
	cmd.AddCommand(newEmailIpPoolCreateCommand(state))
	cmd.AddCommand(newEmailIpPoolUpdateCommand(state))
	cmd.AddCommand(newEmailIpPoolDeleteCommand(state))
	return cmd
}

// newEmailIpPoolListCommand builds `basaltic email ip-pool list`.
func newEmailIpPoolListCommand(state *cli.State) *cobra.Command {
	var params email.ListIPPoolsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List IP pools",
		Args:  cobra.ExactArgs(0),
		Long:  "List IP pools.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/ip-pools")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListIPPoolsAll(cmd.Context(), &params))
			}
			page, err := c.ListIPPools(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact CRN filter")
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Marker")
	f.StringVar(&params.Name, "name", "", "Exact immutable account-scoped name, combined with all other filters")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newEmailIpPoolGetCommand builds `basaltic email ip-pool get`.
func newEmailIpPoolGetCommand(state *cli.State) *cobra.Command {
	var scope email.ListIPPoolsParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Describe an IP pool",
		Args:  cobra.ExactArgs(1),
		Long:  "Describe an IP pool.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/ip-pools/{ip_pool_id}")
			if err != nil {
				return err
			}
			out, err := c.DescribeIPPoolByReference(cmd.Context(), args[0], &scope)
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

// newEmailIpPoolCreateCommand builds `basaltic email ip-pool create`.
func newEmailIpPoolCreateCommand(state *cli.State) *cobra.Command {
	var body email.CreateIPPoolRequest
	var bodyFile string
	var descriptionFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an IP pool",
		Args:  cobra.ExactArgs(0),
		Long:  "Create an IP pool.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"description": "description",
				"name":        "name",
			}, []string{"name"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/ip-pools")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateIPPool(cmd.Context(), &body, reqOpts...)
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
	f.StringVar(&body.Name, "name", "", "Resource names must not start with the literal crn: prefix or be UUIDs (canonical, compact, braced, or urn:uuid: forms, in either case)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newEmailIpPoolUpdateCommand builds `basaltic email ip-pool update`.
func newEmailIpPoolUpdateCommand(state *cli.State) *cobra.Command {
	var body email.UpdateIPPoolRequest
	var bodyFile string
	var descriptionFlag string
	cmd := &cobra.Command{
		Use:   "update <ip-pool-id>",
		Short: "Update an IP pool",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"description": "description",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/ip-pools/{ip_pool_id}")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.UpdateIPPool(cmd.Context(), args[0], &body, reqOpts...)
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
	return cmd
}

// newEmailIpPoolDeleteCommand builds `basaltic email ip-pool delete`.
func newEmailIpPoolDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <ip-pool-id>",
		Short: "Delete an IP pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/ip-pools/{ip_pool_id}")
			if err != nil {
				return err
			}
			if err := c.DeleteIPPool(cmd.Context(), args[0]); err != nil {
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

// newEmailReceivedMessageCommand builds `basaltic email received-message`.
func newEmailReceivedMessageCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "received-message",
		Short:   "Received messages",
		Aliases: []string{"received-messages"},
	}
	cmd.AddCommand(newEmailReceivedMessageListCommand(state))
	cmd.AddCommand(newEmailReceivedMessageGetCommand(state))
	cmd.AddCommand(newEmailReceivedMessageDownloadAttachmentCommand(state))
	cmd.AddCommand(newEmailReceivedMessageDownloadRawCommand(state))
	cmd.AddCommand(newEmailReceivedMessageGetContentCommand(state))
	return cmd
}

// newEmailReceivedMessageListCommand builds `basaltic email received-message list`.
func newEmailReceivedMessageListCommand(state *cli.State) *cobra.Command {
	var params email.ListReceivedMessagesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List received messages",
		Args:  cobra.ExactArgs(0),
		Long:  "List received messages.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/received-messages")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListReceivedMessagesAll(cmd.Context(), &params))
			}
			page, err := c.ListReceivedMessages(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.Identity, "identity", "", "Receiving identity UUID, CRN or name")
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Marker")
	f.StringVar(&params.Recipient, "recipient", "", "Recipient")
	f.StringVar(&params.ReviewReason, "review-reason", "", "One of: \"clear\", \"spam\", \"sender_unverified\", \"malware\"")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newEmailReceivedMessageGetCommand builds `basaltic email received-message get`.
func newEmailReceivedMessageGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <message-id>",
		Short: "Get received message metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/received-messages/{message_id}")
			if err != nil {
				return err
			}
			out, err := c.DescribeReceivedMessage(cmd.Context(), args[0])
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

// newEmailReceivedMessageDownloadAttachmentCommand builds `basaltic email received-message download-attachment`.
func newEmailReceivedMessageDownloadAttachmentCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download-attachment <message-id> <attachment-id>",
		Short: "Download received attachment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/received-messages/{message_id}/attachments/{attachment_id}")
			if err != nil {
				return err
			}
			stream, err := c.DownloadReceivedAttachment(cmd.Context(), args[0], args[1])
			if err != nil {
				return err
			}
			return state.Printer().Stream(stream)
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newEmailReceivedMessageDownloadRawCommand builds `basaltic email received-message download-raw`.
func newEmailReceivedMessageDownloadRawCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download-raw <message-id>",
		Short: "Download received message",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/received-messages/{message_id}/raw")
			if err != nil {
				return err
			}
			stream, err := c.DownloadReceivedMessage(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.Printer().Stream(stream)
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newEmailReceivedMessageGetContentCommand builds `basaltic email received-message get-content`.
func newEmailReceivedMessageGetContentCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-content <message-id>",
		Short: "Get received message content",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/received-messages/{message_id}/content")
			if err != nil {
				return err
			}
			out, err := c.GetReceivedMessageContent(cmd.Context(), args[0])
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

// newEmailRecipientRuleCommand builds `basaltic email recipient-rule`.
func newEmailRecipientRuleCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "recipient-rule",
		Short:   "Recipient rules",
		Aliases: []string{"recipient-rules"},
	}
	cmd.AddCommand(newEmailRecipientRuleListCommand(state))
	cmd.AddCommand(newEmailRecipientRuleGetCommand(state))
	cmd.AddCommand(newEmailRecipientRuleCreateCommand(state))
	cmd.AddCommand(newEmailRecipientRuleUpdateCommand(state))
	cmd.AddCommand(newEmailRecipientRuleDeleteCommand(state))
	return cmd
}

// newEmailRecipientRuleListCommand builds `basaltic email recipient-rule list`.
func newEmailRecipientRuleListCommand(state *cli.State) *cobra.Command {
	var params email.ListRecipientRulesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list <identity-id>",
		Short: "List recipient rules",
		Args:  cobra.ExactArgs(1),
		Long:  "List recipient rules.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/identities/{identity_id}/recipient-rules")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListRecipientRulesAll(cmd.Context(), args[0], &params))
			}
			page, err := c.ListRecipientRules(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Marker")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newEmailRecipientRuleGetCommand builds `basaltic email recipient-rule get`.
func newEmailRecipientRuleGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <identity-id> <rule-id>",
		Short: "Get recipient rule",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/identities/{identity_id}/recipient-rules/{rule_id}")
			if err != nil {
				return err
			}
			out, err := c.DescribeRecipientRule(cmd.Context(), args[0], args[1])
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

// newEmailRecipientRuleCreateCommand builds `basaltic email recipient-rule create`.
func newEmailRecipientRuleCreateCommand(state *cli.State) *cobra.Command {
	var body email.CreateRecipientRuleRequest
	var bodyFile string
	var enabledFlag bool
	var localPartFlag string
	cmd := &cobra.Command{
		Use:   "create <identity-id>",
		Short: "Create recipient rule",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"enabled":    "enabled",
				"local_part": "local-part",
				"match_type": "match-type",
				"topic":      "topic",
			}, []string{"match_type", "topic"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/identities/{identity_id}/recipient-rules")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("enabled") {
				body.Enabled = &enabledFlag
			}
			if cmd.Flags().Changed("local-part") {
				body.LocalPart = &localPartFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.CreateRecipientRule(cmd.Context(), args[0], &body, reqOpts...)
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
	f.StringVar(&localPartFlag, "local-part", "", "Required for exact and prefix; empty or omitted for all")
	f.StringVar(&body.MatchType, "match-type", "", "Exact local part, local-part prefix, or catch-all for this identity (one of: exact, prefix, all)")
	f.StringVar(&body.Topic, "topic", "", "Notifications topic name, UUID or CRN")
	return cmd
}

// newEmailRecipientRuleUpdateCommand builds `basaltic email recipient-rule update`.
func newEmailRecipientRuleUpdateCommand(state *cli.State) *cobra.Command {
	var body email.UpdateRecipientRuleRequest
	var bodyFile string
	var enabledFlag bool
	var localPartFlag string
	var matchTypeFlag string
	var topicFlag string
	cmd := &cobra.Command{
		Use:   "update <identity-id> <rule-id>",
		Short: "Update recipient rule",
		Args:  cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"enabled":    "enabled",
				"local_part": "local-part",
				"match_type": "match-type",
				"topic":      "topic",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/identities/{identity_id}/recipient-rules/{rule_id}")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("enabled") {
				body.Enabled = &enabledFlag
			}
			if cmd.Flags().Changed("local-part") {
				body.LocalPart = &localPartFlag
			}
			if cmd.Flags().Changed("match-type") {
				body.MatchType = &matchTypeFlag
			}
			if cmd.Flags().Changed("topic") {
				body.Topic = &topicFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.UpdateRecipientRule(cmd.Context(), args[0], args[1], &body, reqOpts...)
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
	f.StringVar(&localPartFlag, "local-part", "", "Required for exact and prefix; empty or omitted for all")
	f.StringVar(&matchTypeFlag, "match-type", "", "Exact local part, local-part prefix, or catch-all for this identity (one of: exact, prefix, all)")
	f.StringVar(&topicFlag, "topic", "", "Notifications topic name, UUID or CRN")
	return cmd
}

// newEmailRecipientRuleDeleteCommand builds `basaltic email recipient-rule delete`.
func newEmailRecipientRuleDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <identity-id> <rule-id>",
		Short: "Delete recipient rule",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/identities/{identity_id}/recipient-rules/{rule_id}")
			if err != nil {
				return err
			}
			if err := c.DeleteRecipientRule(cmd.Context(), args[0], args[1]); err != nil {
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

// newEmailRejectionCommand builds `basaltic email rejection`.
func newEmailRejectionCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rejection",
		Short:   "Rejections",
		Aliases: []string{"rejections"},
	}
	cmd.AddCommand(newEmailRejectionListCommand(state))
	return cmd
}

// newEmailRejectionListCommand builds `basaltic email rejection list`.
func newEmailRejectionListCommand(state *cli.State) *cobra.Command {
	var params email.ListRejectionsParams
	var beforeFlag string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List API-side send rejections",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/rejections")
			if err != nil {
				return err
			}
			if beforeFlag != "" {
				parsed, err := parseTime(beforeFlag)
				if err != nil {
					return fmt.Errorf("--before: %w", err)
				}
				params.Before = parsed
			}
			page, err := c.ListRejections(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&beforeFlag, "before", "", "RFC 3339 timestamp — return only rejections older than this (RFC 3339)")
	f.StringVar(&params.CRN, "crn", "", "Exact rejection event CRN")
	f.StringVar(&params.Identity, "identity", "", "Identity UUID, CRN or account-unique name; missing or foreign references return empty results")
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Reason, "reason", "", "Narrow to one rejection reason")
	return cmd
}

// newEmailSentMessageCommand builds `basaltic email sent-message`.
func newEmailSentMessageCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sent-message",
		Short:   "Sent messages",
		Aliases: []string{"sent-messages"},
	}
	cmd.AddCommand(newEmailSentMessageListCommand(state))
	cmd.AddCommand(newEmailSentMessageGetCommand(state))
	cmd.AddCommand(newEmailSentMessageDownloadAttachmentCommand(state))
	cmd.AddCommand(newEmailSentMessageDownloadRawCommand(state))
	cmd.AddCommand(newEmailSentMessageGetContentCommand(state))
	cmd.AddCommand(newEmailSentMessageGetEngagementCommand(state))
	cmd.AddCommand(newEmailSentMessageSendCommand(state))
	cmd.AddCommand(newEmailSentMessageSendBatchCommand(state))
	return cmd
}

// newEmailSentMessageListCommand builds `basaltic email sent-message list`.
func newEmailSentMessageListCommand(state *cli.State) *cobra.Command {
	var params email.ListSentMessagesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List sent messages",
		Args:  cobra.ExactArgs(0),
		Long:  "List sent messages.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/sent-messages")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListSentMessagesAll(cmd.Context(), &params))
			}
			page, err := c.ListSentMessages(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact CRN filter")
	f.StringVar(&params.Identity, "identity", "", "Identity UUID, CRN or account-unique name; missing or foreign references return empty results")
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Marker")
	f.StringVar(&params.Recipient, "recipient", "", "Recipient")
	f.StringVar(&params.Status, "status", "", "Status")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newEmailSentMessageGetCommand builds `basaltic email sent-message get`.
func newEmailSentMessageGetCommand(state *cli.State) *cobra.Command {
	var scope email.ListSentMessagesParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Describe a sent message",
		Args:  cobra.ExactArgs(1),
		Long:  "Describe a sent message.\n\n<ref> is its id or its CRN. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.\n\nThis resource has no name.\n\nA name is unique only within its parent: pass --identity with a name, or the\nlookup can match more than one and is refused.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/sent-messages/{message_id}")
			if err != nil {
				return err
			}
			out, err := c.DescribeSentMessageByReference(cmd.Context(), args[0], &scope)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&scope.Identity, "identity", "", "Identity UUID, CRN or account-unique name; missing or foreign references return empty results")
	return cmd
}

// newEmailSentMessageDownloadAttachmentCommand builds `basaltic email sent-message download-attachment`.
func newEmailSentMessageDownloadAttachmentCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download-attachment <message-id> <attachment-id>",
		Short: "Download sent attachment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/sent-messages/{message_id}/attachments/{attachment_id}")
			if err != nil {
				return err
			}
			stream, err := c.DownloadSentAttachment(cmd.Context(), args[0], args[1])
			if err != nil {
				return err
			}
			return state.Printer().Stream(stream)
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newEmailSentMessageDownloadRawCommand builds `basaltic email sent-message download-raw`.
func newEmailSentMessageDownloadRawCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download-raw <message-id>",
		Short: "Download sent message",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/sent-messages/{message_id}/raw")
			if err != nil {
				return err
			}
			stream, err := c.DownloadSentMessage(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.Printer().Stream(stream)
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newEmailSentMessageGetContentCommand builds `basaltic email sent-message get-content`.
func newEmailSentMessageGetContentCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-content <message-id>",
		Short: "Get sent message content",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/sent-messages/{message_id}/content")
			if err != nil {
				return err
			}
			out, err := c.GetSentMessageContent(cmd.Context(), args[0])
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

// newEmailSentMessageGetEngagementCommand builds `basaltic email sent-message get-engagement`.
func newEmailSentMessageGetEngagementCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-engagement <message-id>",
		Short: "Sent message engagement",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/sent-messages/{message_id}/engagement")
			if err != nil {
				return err
			}
			out, err := c.GetSentMessageEngagement(cmd.Context(), args[0])
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

// newEmailSentMessageSendCommand builds `basaltic email sent-message send`.
func newEmailSentMessageSendCommand(state *cli.State) *cobra.Command {
	var body email.SendMessageRequest
	var bodyFile string
	var attachmentsFlag string
	var attachmentsBucketFlag string
	var configurationSetFlag string
	var contentFlag string
	var subjectFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send a message",
		Args:  cobra.ExactArgs(0),
		Long:  "Send a message.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"attachments":        "attachments",
				"attachments_bucket": "attachments-bucket",
				"configuration_set":  "configuration-set",
				"content":            "content",
				"from":               "from",
				"identity":           "identity",
				"subject":            "subject",
				"to":                 "to",
			}, []string{"content", "from", "identity", "to"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/sent-messages")
			if err != nil {
				return err
			}
			if attachmentsFlag != "" {
				if err := json.Unmarshal([]byte(attachmentsFlag), &body.Attachments); err != nil {
					return fmt.Errorf("--attachments: %w", err)
				}
			}
			if cmd.Flags().Changed("attachments-bucket") {
				body.AttachmentsBucket = &attachmentsBucketFlag
			}
			if cmd.Flags().Changed("configuration-set") {
				body.ConfigurationSet = &configurationSetFlag
			}
			if contentFlag != "" {
				if err := json.Unmarshal([]byte(contentFlag), &body.Content); err != nil {
					return fmt.Errorf("--content: %w", err)
				}
			}
			if cmd.Flags().Changed("subject") {
				body.Subject = &subjectFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.SendMessage(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&attachmentsFlag, "attachments", "", "Bucket objects to attach (JSON)")
	f.StringVar(&attachmentsBucketFlag, "attachments-bucket", "", "Bucket UUID, CRN or name for attachments embedded in raw MIME")
	f.StringVar(&configurationSetFlag, "configuration-set", "", "UUID, CRN or account-unique name")
	f.StringVar(&contentFlag, "content", "", "Content (JSON)")
	f.StringVar(&body.From, "from", "", "From")
	f.StringVar(&body.Identity, "identity", "", "UUID, CRN or account-unique name")
	f.StringVar(&subjectFlag, "subject", "", "Subject for text/HTML content")
	f.StringSliceVar(&body.To, "to", nil, "To")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newEmailSentMessageSendBatchCommand builds `basaltic email sent-message send-batch`.
func newEmailSentMessageSendBatchCommand(state *cli.State) *cobra.Command {
	var body email.BatchSendRequest
	var bodyFile string
	var messagesFlag string
	cmd := &cobra.Command{
		Use:   "send-batch",
		Short: "Send message batch",
		Args:  cobra.ExactArgs(0),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"messages": "messages",
			}, []string{"messages"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/sent-messages/batch")
			if err != nil {
				return err
			}
			if messagesFlag != "" {
				if err := json.Unmarshal([]byte(messagesFlag), &body.Messages); err != nil {
					return fmt.Errorf("--messages: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			page, err := c.SendMessageBatch(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&messagesFlag, "messages", "", "Messages (JSON)")
	return cmd
}

// newEmailSuppressionCommand builds `basaltic email suppression`.
func newEmailSuppressionCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "suppression",
		Short:   "Suppressions",
		Aliases: []string{"suppressions"},
	}
	cmd.AddCommand(newEmailSuppressionListCommand(state))
	cmd.AddCommand(newEmailSuppressionGetCommand(state))
	cmd.AddCommand(newEmailSuppressionCreateCommand(state))
	cmd.AddCommand(newEmailSuppressionUpdateCommand(state))
	cmd.AddCommand(newEmailSuppressionDeleteCommand(state))
	return cmd
}

// newEmailSuppressionListCommand builds `basaltic email suppression list`.
func newEmailSuppressionListCommand(state *cli.State) *cobra.Command {
	var params email.ListSuppressionsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List suppression entries",
		Args:  cobra.ExactArgs(0),
		Long:  "List suppression entries.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/suppressions")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListSuppressionsAll(cmd.Context(), &params))
			}
			page, err := c.ListSuppressions(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact CRN filter")
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Marker")
	f.StringVar(&params.Reason, "reason", "", "Return only entries put on the list this way")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newEmailSuppressionGetCommand builds `basaltic email suppression get`.
func newEmailSuppressionGetCommand(state *cli.State) *cobra.Command {
	var scope email.ListSuppressionsParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Get suppression",
		Args:  cobra.ExactArgs(1),
		Long:  "Get suppression.\n\n<ref> is its id or its CRN. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.\n\nThis resource has no name.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/suppressions/{suppression_id}")
			if err != nil {
				return err
			}
			out, err := c.DescribeSuppressionByReference(cmd.Context(), args[0], &scope)
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

// newEmailSuppressionCreateCommand builds `basaltic email suppression create`.
func newEmailSuppressionCreateCommand(state *cli.State) *cobra.Command {
	var body email.CreateSuppressionRequest
	var bodyFile string
	var descriptionFlag string
	var expiresAtFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Add a suppression entry",
		Args:  cobra.ExactArgs(0),
		Long:  "Add a suppression entry.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"address":     "address",
				"description": "description",
				"expires_at":  "expires-at",
			}, []string{"address"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/suppressions")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
			}
			if expiresAtFlag != "" {
				parsed, err := parseTime(expiresAtFlag)
				if err != nil {
					return fmt.Errorf("--expires-at: %w", err)
				}
				body.ExpiresAt = &parsed
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateSuppression(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.Address, "address", "", "Address")
	f.StringVar(&descriptionFlag, "description", "", "Description")
	f.StringVar(&expiresAtFlag, "expires-at", "", "Expires at (RFC 3339)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newEmailSuppressionUpdateCommand builds `basaltic email suppression update`.
func newEmailSuppressionUpdateCommand(state *cli.State) *cobra.Command {
	var body email.UpdateSuppressionRequest
	var bodyFile string
	var descriptionFlag string
	var expiresAtFlag string
	cmd := &cobra.Command{
		Use:   "update <suppression-id>",
		Short: "Update suppression",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"description": "description",
				"expires_at":  "expires-at",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/suppressions/{suppression_id}")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
			}
			if expiresAtFlag != "" {
				parsed, err := parseTime(expiresAtFlag)
				if err != nil {
					return fmt.Errorf("--expires-at: %w", err)
				}
				body.ExpiresAt = &parsed
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.UpdateSuppression(cmd.Context(), args[0], &body, reqOpts...)
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
	f.StringVar(&expiresAtFlag, "expires-at", "", "Null makes the suppression permanent; omission preserves its expiry (RFC 3339)")
	return cmd
}

// newEmailSuppressionDeleteCommand builds `basaltic email suppression delete`.
func newEmailSuppressionDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <suppression-id>",
		Short: "Delete suppression",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/suppressions/{suppression_id}")
			if err != nil {
				return err
			}
			if err := c.DeleteSuppression(cmd.Context(), args[0]); err != nil {
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

// newEmailTemplateCommand builds `basaltic email template`.
func newEmailTemplateCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "template",
		Short:   "Templates",
		Aliases: []string{"templates"},
	}
	cmd.AddCommand(newEmailTemplateListCommand(state))
	cmd.AddCommand(newEmailTemplateGetCommand(state))
	cmd.AddCommand(newEmailTemplateCreateCommand(state))
	cmd.AddCommand(newEmailTemplateUpdateCommand(state))
	cmd.AddCommand(newEmailTemplateDeleteCommand(state))
	return cmd
}

// newEmailTemplateListCommand builds `basaltic email template list`.
func newEmailTemplateListCommand(state *cli.State) *cobra.Command {
	var params email.ListTemplatesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List templates",
		Args:  cobra.ExactArgs(0),
		Long:  "List templates.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/templates")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListTemplatesAll(cmd.Context(), &params))
			}
			page, err := c.ListTemplates(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact CRN filter")
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Marker")
	f.StringVar(&params.Name, "name", "", "Exact immutable account-scoped name, combined with all other filters")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newEmailTemplateGetCommand builds `basaltic email template get`.
func newEmailTemplateGetCommand(state *cli.State) *cobra.Command {
	var scope email.ListTemplatesParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Describe a template",
		Args:  cobra.ExactArgs(1),
		Long:  "Describe a template.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/templates/{template_id}")
			if err != nil {
				return err
			}
			out, err := c.DescribeTemplateByReference(cmd.Context(), args[0], &scope)
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

// newEmailTemplateCreateCommand builds `basaltic email template create`.
func newEmailTemplateCreateCommand(state *cli.State) *cobra.Command {
	var body email.CreateTemplateRequest
	var bodyFile string
	var bodyTypeFlag string
	var subjectTemplateFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a template",
		Args:  cobra.ExactArgs(0),
		Long:  "Create a template.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"body_template":    "body-template",
				"body_type":        "body-type",
				"name":             "name",
				"subject_template": "subject-template",
			}, []string{"body_template", "name"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/templates")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("body-type") {
				body.BodyType = &bodyTypeFlag
			}
			if cmd.Flags().Changed("subject-template") {
				body.SubjectTemplate = &subjectTemplateFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateTemplate(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.BodyTemplate, "body-template", "", "Body template")
	f.StringVar(&bodyTypeFlag, "body-type", "", "Body type (one of: text, html)")
	f.StringVar(&body.Name, "name", "", "Resource names must not start with the literal crn: prefix or be UUIDs (canonical, compact, braced, or urn:uuid: forms, in either case)")
	f.StringVar(&subjectTemplateFlag, "subject-template", "", "Subject template")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newEmailTemplateUpdateCommand builds `basaltic email template update`.
func newEmailTemplateUpdateCommand(state *cli.State) *cobra.Command {
	var body email.UpdateTemplateRequest
	var bodyFile string
	var bodyTemplateFlag string
	var bodyTypeFlag string
	var subjectTemplateFlag string
	cmd := &cobra.Command{
		Use:   "update <template-id>",
		Short: "Update a template",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"body_template":    "body-template",
				"body_type":        "body-type",
				"subject_template": "subject-template",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/templates/{template_id}")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("body-template") {
				body.BodyTemplate = &bodyTemplateFlag
			}
			if cmd.Flags().Changed("body-type") {
				body.BodyType = &bodyTypeFlag
			}
			if cmd.Flags().Changed("subject-template") {
				body.SubjectTemplate = &subjectTemplateFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.UpdateTemplate(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&bodyTemplateFlag, "body-template", "", "Body template")
	f.StringVar(&bodyTypeFlag, "body-type", "", "Body type (one of: text, html)")
	f.StringVar(&subjectTemplateFlag, "subject-template", "", "Subject template")
	return cmd
}

// newEmailTemplateDeleteCommand builds `basaltic email template delete`.
func newEmailTemplateDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <template-id>",
		Short: "Delete a template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := emailClient(state, "/v1/templates/{template_id}")
			if err != nil {
				return err
			}
			if err := c.DeleteTemplate(cmd.Context(), args[0]); err != nil {
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
