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
	"github.com/basaltic-sh/sdk-go/loadbalancer"

	"github.com/basaltic-sh/cli/internal/cli"
)

func init() { cli.RegisterService(newLoadbalancerCommand) }

// newLoadbalancerCommand builds `basaltic loadbalancer`.
func newLoadbalancerCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "loadbalancer",
		Short:   "Load balancers, listeners, rules and target groups",
		Aliases: []string{"lb"},
		Long:    "Load balancers, listeners, rules and target groups.\n\nThis is a regional service: it acts in the region from --region, the\nBASALTIC_REGION environment variable, or the profile.",
	}
	cmd.AddCommand(newLoadbalancerListenerCommand(state))
	cmd.AddCommand(newLoadbalancerLoadBalancerCommand(state))
	cmd.AddCommand(newLoadbalancerRuleCommand(state))
	cmd.AddCommand(newLoadbalancerTargetGroupCommand(state))
	return cmd
}

// loadbalancerClient builds the service client, resolving credentials on first use.
func loadbalancerClient(state *cli.State, path string) (*loadbalancer.Client, error) {
	cfg, err := state.ServiceSDK("loadbalancer", path)
	if err != nil {
		return nil, err
	}
	return loadbalancer.New(cfg), nil
}

// newLoadbalancerListenerCommand builds `basaltic loadbalancer listener`.
func newLoadbalancerListenerCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "listener",
		Short:   "Listeners",
		Aliases: []string{"listeners"},
	}
	cmd.AddCommand(newLoadbalancerListenerListCommand(state))
	cmd.AddCommand(newLoadbalancerListenerGetCommand(state))
	cmd.AddCommand(newLoadbalancerListenerCreateCommand(state))
	cmd.AddCommand(newLoadbalancerListenerUpdateCommand(state))
	cmd.AddCommand(newLoadbalancerListenerDeleteCommand(state))
	cmd.AddCommand(newLoadbalancerListenerAttachCertificateCommand(state))
	cmd.AddCommand(newLoadbalancerListenerDetachCertificateCommand(state))
	return cmd
}

// newLoadbalancerListenerListCommand builds `basaltic loadbalancer listener list`.
func newLoadbalancerListenerListCommand(state *cli.State) *cobra.Command {
	var params loadbalancer.ListListenersParams
	cmd := &cobra.Command{
		Use:   "list <id>",
		Short: "List this load balancer's listeners",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/load-balancers/{id}/listeners")
			if err != nil {
				return err
			}
			page, err := c.ListListeners(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact scoped CRN")
	f.StringVar(&params.Name, "name", "", "Exact immutable name")
	return cmd
}

// newLoadbalancerListenerGetCommand builds `basaltic loadbalancer listener get`.
func newLoadbalancerListenerGetCommand(state *cli.State) *cobra.Command {
	var scope loadbalancer.ListListenersParams
	cmd := &cobra.Command{
		Use:   "get <id> <ref>",
		Short: "Get a listener",
		Args:  cobra.ExactArgs(2),
		Long:  "Get a listener.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/load-balancers/{id}/listeners/{listener_id}")
			if err != nil {
				return err
			}
			out, err := c.GetListenerByReference(cmd.Context(), args[0], args[1], &scope)
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

// newLoadbalancerListenerCreateCommand builds `basaltic loadbalancer listener create`.
func newLoadbalancerListenerCreateCommand(state *cli.State) *cobra.Command {
	var body loadbalancer.CreateListenerRequest
	var bodyFile string
	var certificatesFlag string
	var defaultTargetGroupFlag string
	var exposureFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create <id>",
		Short: "Create a listener on this load balancer",
		Args:  cobra.ExactArgs(1),
		Long:  "Create a listener on this load balancer.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"certificates":         "certificates",
				"default_target_group": "default-target-group",
				"exposure":             "exposure",
				"port":                 "port",
				"protocol":             "protocol",
				"tags":                 "tags",
			}, []string{"port", "protocol"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/load-balancers/{id}/listeners")
			if err != nil {
				return err
			}
			if certificatesFlag != "" {
				if err := json.Unmarshal([]byte(certificatesFlag), &body.Certificates); err != nil {
					return fmt.Errorf("--certificates: %w", err)
				}
			}
			if cmd.Flags().Changed("default-target-group") {
				body.DefaultTargetGroup = &defaultTargetGroupFlag
			}
			if cmd.Flags().Changed("exposure") {
				body.Exposure = &exposureFlag
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
			out, err := c.CreateListener(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&certificatesFlag, "certificates", "", "Certificates (JSON)")
	f.StringVar(&defaultTargetGroupFlag, "default-target-group", "", "Default target group")
	f.StringVar(&exposureFlag, "exposure", "", "Which LB addresses are bound (one of: public_only, private_only, both)")
	f.IntVar(&body.Port, "port", 0, "Port")
	f.StringVar(&body.Protocol, "protocol", "", "Protocol (one of: http, https, tcp, udp)")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newLoadbalancerListenerUpdateCommand builds `basaltic loadbalancer listener update`.
func newLoadbalancerListenerUpdateCommand(state *cli.State) *cobra.Command {
	var body loadbalancer.UpdateListenerRequest
	var bodyFile string
	var certificateFlag string
	var clearDefaultTargetGroupFlag bool
	var defaultTargetGroupFlag string
	var exposureFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <id> <listener-id>",
		Short: "Patch a listener (rotate cert, change default target group)",
		Args:  cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"certificate":                "certificate",
				"clear_default_target_group": "clear-default-target-group",
				"default_target_group":       "default-target-group",
				"exposure":                   "exposure",
				"tags":                       "tags",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/load-balancers/{id}/listeners/{listener_id}")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("certificate") {
				body.Certificate = &certificateFlag
			}
			if cmd.Flags().Changed("clear-default-target-group") {
				body.ClearDefaultTargetGroup = &clearDefaultTargetGroupFlag
			}
			if cmd.Flags().Changed("default-target-group") {
				body.DefaultTargetGroup = &defaultTargetGroupFlag
			}
			if cmd.Flags().Changed("exposure") {
				body.Exposure = &exposureFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateListener(cmd.Context(), args[0], args[1], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&certificateFlag, "certificate", "", "Certificate")
	f.BoolVar(&clearDefaultTargetGroupFlag, "clear-default-target-group", false, "Clear default target group")
	f.StringVar(&defaultTargetGroupFlag, "default-target-group", "", "Default target group")
	f.StringVar(&exposureFlag, "exposure", "", "Mutate which addresses are bound (one of: public_only, private_only, both)")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newLoadbalancerListenerDeleteCommand builds `basaltic loadbalancer listener delete`.
func newLoadbalancerListenerDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <id> <listener-id>",
		Short: "Delete a listener",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/load-balancers/{id}/listeners/{listener_id}")
			if err != nil {
				return err
			}
			if err := c.DeleteListener(cmd.Context(), args[0], args[1]); err != nil {
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

// newLoadbalancerListenerAttachCertificateCommand builds `basaltic loadbalancer listener attach-certificate`.
func newLoadbalancerListenerAttachCertificateCommand(state *cli.State) *cobra.Command {
	var body loadbalancer.AttachListenerCertificateRequest
	var bodyFile string
	var isDefaultFlag bool
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "attach-certificate <id> <listener-id>",
		Short: "Attach an additional certificate to an HTTPS listener",
		Args:  cobra.ExactArgs(2),
		Long:  "Attach an additional certificate to an HTTPS listener.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"certificate": "certificate",
				"is_default":  "is-default",
			}, []string{"certificate"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/load-balancers/{id}/listeners/{listener_id}/certificates")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("is-default") {
				body.IsDefault = &isDefaultFlag
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.AttachListenerCertificate(cmd.Context(), args[0], args[1], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.Certificate, "certificate", "", "The certificate to serve, by CRN, UUID or exact account-scoped name")
	f.BoolVar(&isDefaultFlag, "is-default", false, "When true, demote whatever's currently default and promote this cert in the same transaction")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newLoadbalancerListenerDetachCertificateCommand builds `basaltic loadbalancer listener detach-certificate`.
func newLoadbalancerListenerDetachCertificateCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach-certificate <id> <listener-id> <certificate-id>",
		Short: "Detach a certificate from an HTTPS listener",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/load-balancers/{id}/listeners/{listener_id}/certificates/{certificate_id}")
			if err != nil {
				return err
			}
			if err := c.DetachListenerCertificate(cmd.Context(), args[0], args[1], args[2]); err != nil {
				return err
			}
			state.Printer().Done("Detach certificate requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newLoadbalancerLoadBalancerCommand builds `basaltic loadbalancer load-balancer`.
func newLoadbalancerLoadBalancerCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "load-balancer",
		Short:   "Load balancers",
		Aliases: []string{"load-balancers"},
	}
	cmd.AddCommand(newLoadbalancerLoadBalancerListCommand(state))
	cmd.AddCommand(newLoadbalancerLoadBalancerGetCommand(state))
	cmd.AddCommand(newLoadbalancerLoadBalancerCreateCommand(state))
	cmd.AddCommand(newLoadbalancerLoadBalancerUpdateCommand(state))
	cmd.AddCommand(newLoadbalancerLoadBalancerDeleteCommand(state))
	cmd.AddCommand(newLoadbalancerLoadBalancerListReplicasCommand(state))
	return cmd
}

// newLoadbalancerLoadBalancerListCommand builds `basaltic loadbalancer load-balancer list`.
func newLoadbalancerLoadBalancerListCommand(state *cli.State) *cobra.Command {
	var params loadbalancer.ListLoadBalancersParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List load balancers",
		Args:  cobra.ExactArgs(0),
		Long:  "List load balancers.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/load-balancers")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListLoadBalancersAll(cmd.Context(), &params))
			}
			page, err := c.ListLoadBalancers(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact scoped CRN")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Exact immutable name")
	f.StringVar(&params.Status, "status", "", "One of: \"provisioning\", \"active\", \"error\", \"deleting\"")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newLoadbalancerLoadBalancerGetCommand builds `basaltic loadbalancer load-balancer get`.
func newLoadbalancerLoadBalancerGetCommand(state *cli.State) *cobra.Command {
	var scope loadbalancer.ListLoadBalancersParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Get a load balancer",
		Args:  cobra.ExactArgs(1),
		Long:  "Get a load balancer.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/load-balancers/{id}")
			if err != nil {
				return err
			}
			out, err := c.GetLoadBalancerByReference(cmd.Context(), args[0], &scope)
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

// newLoadbalancerLoadBalancerCreateCommand builds `basaltic loadbalancer load-balancer create`.
func newLoadbalancerLoadBalancerCreateCommand(state *cli.State) *cobra.Command {
	var body loadbalancer.CreateLoadBalancerRequest
	var bodyFile string
	var autoscalingFlag string
	var desiredCountFlag int
	var floatingIpFlag string
	var maxCountFlag int
	var minCountFlag int
	var replicaCountFlag int
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a load balancer",
		Args:  cobra.ExactArgs(0),
		Long:  "Create a load balancer.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"autoscaling":     "autoscaling",
				"desired_count":   "desired-count",
				"flavor":          "flavor",
				"floating_ip":     "floating-ip",
				"floating_ips":    "floating-ips",
				"max_count":       "max-count",
				"min_count":       "min-count",
				"name":            "name",
				"replica_count":   "replica-count",
				"security_groups": "security-groups",
				"subnet":          "subnet",
				"tags":            "tags",
				"type":            "type",
				"vpc":             "vpc",
			}, []string{"flavor", "name", "security_groups", "subnet", "type", "vpc"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/load-balancers")
			if err != nil {
				return err
			}
			if autoscalingFlag != "" {
				if err := json.Unmarshal([]byte(autoscalingFlag), &body.Autoscaling); err != nil {
					return fmt.Errorf("--autoscaling: %w", err)
				}
			}
			if cmd.Flags().Changed("desired-count") {
				body.DesiredCount = &desiredCountFlag
			}
			if cmd.Flags().Changed("floating-ip") {
				body.FloatingIP = &floatingIpFlag
			}
			if cmd.Flags().Changed("max-count") {
				body.MaxCount = &maxCountFlag
			}
			if cmd.Flags().Changed("min-count") {
				body.MinCount = &minCountFlag
			}
			if cmd.Flags().Changed("replica-count") {
				body.ReplicaCount = &replicaCountFlag
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
			out, err := c.CreateLoadBalancer(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&autoscalingFlag, "autoscaling", "", "Autoscaling (JSON)")
	f.IntVar(&desiredCountFlag, "desired-count", 0, "Steady target within min_count and max_count")
	f.StringVar(&body.Flavor, "flavor", "", "Compute flavor for each LB instance")
	f.StringVar(&floatingIpFlag, "floating-ip", "", "Public IPv4 shorthand")
	f.StringSliceVar(&body.FloatingIPs, "floating-ips", nil, "Existing free floating IPs from this account and region, at most one per family and visibility (private/public, IPv4/IPv6)")
	f.IntVar(&maxCountFlag, "max-count", 0, "Upper capacity bound including rollout surge")
	f.IntVar(&minCountFlag, "min-count", 0, "Lower capacity bound")
	f.StringVar(&body.Name, "name", "", "1..127 chars of [A-Za-z0-9._-] Resource names must not start with the literal crn: prefix or be UUIDs (canonical, compact, braced, or urn:uuid: forms, in either case)")
	f.IntVar(&replicaCountFlag, "replica-count", 0, "Deprecated input alias of desired_count; send only one")
	f.StringSliceVar(&body.SecurityGroups, "security-groups", nil, "Security groups attached to every replica NIC (AWS ALB shape)")
	f.StringVar(&body.Subnet, "subnet", "", "Subnet the LB instances attach to")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&body.Type, "type", "", "Type (one of: application, network)")
	f.StringVar(&body.VPC, "vpc", "", "VPC the LB will live in")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newLoadbalancerLoadBalancerUpdateCommand builds `basaltic loadbalancer load-balancer update`.
func newLoadbalancerLoadBalancerUpdateCommand(state *cli.State) *cobra.Command {
	var body loadbalancer.UpdateLoadBalancerRequest
	var bodyFile string
	var autoscalingFlag string
	var desiredCountFlag int
	var flavorFlag string
	var maxCountFlag int
	var minCountFlag int
	var replicaCountFlag int
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Scale or resize a load balancer",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"autoscaling":   "autoscaling",
				"desired_count": "desired-count",
				"flavor":        "flavor",
				"max_count":     "max-count",
				"min_count":     "min-count",
				"replica_count": "replica-count",
				"tags":          "tags",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/load-balancers/{id}")
			if err != nil {
				return err
			}
			if autoscalingFlag != "" {
				if err := json.Unmarshal([]byte(autoscalingFlag), &body.Autoscaling); err != nil {
					return fmt.Errorf("--autoscaling: %w", err)
				}
			}
			if cmd.Flags().Changed("desired-count") {
				body.DesiredCount = &desiredCountFlag
			}
			if cmd.Flags().Changed("flavor") {
				body.Flavor = &flavorFlag
			}
			if cmd.Flags().Changed("max-count") {
				body.MaxCount = &maxCountFlag
			}
			if cmd.Flags().Changed("min-count") {
				body.MinCount = &minCountFlag
			}
			if cmd.Flags().Changed("replica-count") {
				body.ReplicaCount = &replicaCountFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateLoadBalancer(cmd.Context(), args[0], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&autoscalingFlag, "autoscaling", "", "Autoscaling (JSON)")
	f.IntVar(&desiredCountFlag, "desired-count", 0, "Steady target within min_count and max_count")
	f.StringVar(&flavorFlag, "flavor", "", "Resize each replica to a different compute flavor")
	f.IntVar(&maxCountFlag, "max-count", 0, "Upper capacity bound including rollout surge")
	f.IntVar(&minCountFlag, "min-count", 0, "Lower capacity bound")
	f.IntVar(&replicaCountFlag, "replica-count", 0, "Deprecated alias of desired_count; send only one")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newLoadbalancerLoadBalancerDeleteCommand builds `basaltic loadbalancer load-balancer delete`.
func newLoadbalancerLoadBalancerDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a load balancer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/load-balancers/{id}")
			if err != nil {
				return err
			}
			if err := c.DeleteLoadBalancer(cmd.Context(), args[0]); err != nil {
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

// newLoadbalancerLoadBalancerListReplicasCommand builds `basaltic loadbalancer load-balancer list-replicas`.
func newLoadbalancerLoadBalancerListReplicasCommand(state *cli.State) *cobra.Command {
	var params loadbalancer.ListLoadBalancerReplicasParams
	cmd := &cobra.Command{
		Use:   "list-replicas <id>",
		Short: "List the LB's instance replicas with live health",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/load-balancers/{id}/replicas")
			if err != nil {
				return err
			}
			page, err := c.ListLoadBalancerReplicas(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact scoped CRN")
	f.StringVar(&params.Name, "name", "", "Exact immutable name")
	return cmd
}

// newLoadbalancerRuleCommand builds `basaltic loadbalancer rule`.
func newLoadbalancerRuleCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rule",
		Short:   "Rules",
		Aliases: []string{"rules"},
	}
	cmd.AddCommand(newLoadbalancerRuleListCommand(state))
	cmd.AddCommand(newLoadbalancerRuleGetCommand(state))
	cmd.AddCommand(newLoadbalancerRuleCreateCommand(state))
	cmd.AddCommand(newLoadbalancerRuleUpdateCommand(state))
	cmd.AddCommand(newLoadbalancerRuleDeleteCommand(state))
	return cmd
}

// newLoadbalancerRuleListCommand builds `basaltic loadbalancer rule list`.
func newLoadbalancerRuleListCommand(state *cli.State) *cobra.Command {
	var params loadbalancer.ListRulesParams
	cmd := &cobra.Command{
		Use:   "list <id> <listener-id>",
		Short: "List this listener's rules",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/load-balancers/{id}/listeners/{listener_id}/rules")
			if err != nil {
				return err
			}
			page, err := c.ListRules(cmd.Context(), args[0], args[1], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact scoped CRN")
	f.StringVar(&params.Name, "name", "", "Exact immutable name")
	return cmd
}

// newLoadbalancerRuleGetCommand builds `basaltic loadbalancer rule get`.
func newLoadbalancerRuleGetCommand(state *cli.State) *cobra.Command {
	var scope loadbalancer.ListRulesParams
	cmd := &cobra.Command{
		Use:   "get <id> <listener-id> <ref>",
		Short: "Get a routing rule",
		Args:  cobra.ExactArgs(3),
		Long:  "Get a routing rule.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/load-balancers/{id}/listeners/{listener_id}/rules/{rule_id}")
			if err != nil {
				return err
			}
			out, err := c.GetRuleByReference(cmd.Context(), args[0], args[1], args[2], &scope)
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

// newLoadbalancerRuleCreateCommand builds `basaltic loadbalancer rule create`.
func newLoadbalancerRuleCreateCommand(state *cli.State) *cobra.Command {
	var body loadbalancer.CreateRuleRequest
	var bodyFile string
	var conditionsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create <id> <listener-id>",
		Short: "Create a routing rule on this listener (HTTP/HTTPS only)",
		Args:  cobra.ExactArgs(2),
		Long:  "Create a routing rule on this listener (HTTP/HTTPS only).\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"conditions":   "conditions",
				"priority":     "priority",
				"target_group": "target-group",
			}, []string{"conditions", "priority", "target_group"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/load-balancers/{id}/listeners/{listener_id}/rules")
			if err != nil {
				return err
			}
			if conditionsFlag != "" {
				if err := json.Unmarshal([]byte(conditionsFlag), &body.Conditions); err != nil {
					return fmt.Errorf("--conditions: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateRule(cmd.Context(), args[0], args[1], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&conditionsFlag, "conditions", "", "Conditions (JSON)")
	f.IntVar(&body.Priority, "priority", 0, "Priority")
	f.StringVar(&body.TargetGroup, "target-group", "", "Target group")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newLoadbalancerRuleUpdateCommand builds `basaltic loadbalancer rule update`.
func newLoadbalancerRuleUpdateCommand(state *cli.State) *cobra.Command {
	var body loadbalancer.UpdateRuleRequest
	var bodyFile string
	var conditionsFlag string
	cmd := &cobra.Command{
		Use:   "update <id> <listener-id> <rule-id>",
		Short: "Update a routing rule (full replace)",
		Args:  cobra.ExactArgs(3),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"conditions":   "conditions",
				"priority":     "priority",
				"target_group": "target-group",
			}, []string{"conditions", "priority", "target_group"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/load-balancers/{id}/listeners/{listener_id}/rules/{rule_id}")
			if err != nil {
				return err
			}
			if conditionsFlag != "" {
				if err := json.Unmarshal([]byte(conditionsFlag), &body.Conditions); err != nil {
					return fmt.Errorf("--conditions: %w", err)
				}
			}
			out, err := c.UpdateRule(cmd.Context(), args[0], args[1], args[2], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&conditionsFlag, "conditions", "", "Conditions (JSON)")
	f.IntVar(&body.Priority, "priority", 0, "Priority")
	f.StringVar(&body.TargetGroup, "target-group", "", "Target group")
	return cmd
}

// newLoadbalancerRuleDeleteCommand builds `basaltic loadbalancer rule delete`.
func newLoadbalancerRuleDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <id> <listener-id> <rule-id>",
		Short: "Delete a routing rule",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/load-balancers/{id}/listeners/{listener_id}/rules/{rule_id}")
			if err != nil {
				return err
			}
			if err := c.DeleteRuleInListener(cmd.Context(), args[0], args[1], args[2]); err != nil {
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

// newLoadbalancerTargetGroupCommand builds `basaltic loadbalancer target-group`.
func newLoadbalancerTargetGroupCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "target-group",
		Short:   "Target groups",
		Aliases: []string{"target-groups"},
	}
	cmd.AddCommand(newLoadbalancerTargetGroupListCommand(state))
	cmd.AddCommand(newLoadbalancerTargetGroupGetCommand(state))
	cmd.AddCommand(newLoadbalancerTargetGroupCreateCommand(state))
	cmd.AddCommand(newLoadbalancerTargetGroupUpdateCommand(state))
	cmd.AddCommand(newLoadbalancerTargetGroupDeleteCommand(state))
	cmd.AddCommand(newLoadbalancerTargetGroupAttachTargetCommand(state))
	cmd.AddCommand(newLoadbalancerTargetGroupDetachTargetCommand(state))
	cmd.AddCommand(newLoadbalancerTargetGroupGetTargetCommand(state))
	cmd.AddCommand(newLoadbalancerTargetGroupListTargetsCommand(state))
	return cmd
}

// newLoadbalancerTargetGroupListCommand builds `basaltic loadbalancer target-group list`.
func newLoadbalancerTargetGroupListCommand(state *cli.State) *cobra.Command {
	var params loadbalancer.ListTargetGroupsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List target groups",
		Args:  cobra.ExactArgs(0),
		Long:  "List target groups.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/target-groups")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListTargetGroupsAll(cmd.Context(), &params))
			}
			page, err := c.ListTargetGroups(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact scoped CRN")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Exact immutable name")
	f.StringVar(&params.Protocol, "protocol", "", "One of: \"http\", \"https\", \"tcp\", \"udp\"")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newLoadbalancerTargetGroupGetCommand builds `basaltic loadbalancer target-group get`.
func newLoadbalancerTargetGroupGetCommand(state *cli.State) *cobra.Command {
	var scope loadbalancer.ListTargetGroupsParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Get a target group",
		Args:  cobra.ExactArgs(1),
		Long:  "Get a target group.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/target-groups/{id}")
			if err != nil {
				return err
			}
			out, err := c.GetTargetGroupByReference(cmd.Context(), args[0], &scope)
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

// newLoadbalancerTargetGroupCreateCommand builds `basaltic loadbalancer target-group create`.
func newLoadbalancerTargetGroupCreateCommand(state *cli.State) *cobra.Command {
	var body loadbalancer.CreateTargetGroupRequest
	var bodyFile string
	var healthCheckFlag string
	var instancePoolFlag string
	var proxyProtocolFlag bool
	var sessionAffinityFlag string
	var tagsFlag string
	var targetModeFlag string
	var targetTypeFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a target group",
		Args:  cobra.ExactArgs(0),
		Long:  "Create a target group.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"health_check":     "health-check",
				"instance_pool":    "instance-pool",
				"name":             "name",
				"port":             "port",
				"protocol":         "protocol",
				"proxy_protocol":   "proxy-protocol",
				"session_affinity": "session-affinity",
				"tags":             "tags",
				"target_mode":      "target-mode",
				"target_type":      "target-type",
			}, []string{"name", "port", "protocol"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/target-groups")
			if err != nil {
				return err
			}
			if healthCheckFlag != "" {
				if err := json.Unmarshal([]byte(healthCheckFlag), &body.HealthCheck); err != nil {
					return fmt.Errorf("--health-check: %w", err)
				}
			}
			if cmd.Flags().Changed("instance-pool") {
				body.InstancePool = &instancePoolFlag
			}
			if cmd.Flags().Changed("proxy-protocol") {
				body.ProxyProtocol = &proxyProtocolFlag
			}
			if sessionAffinityFlag != "" {
				if err := json.Unmarshal([]byte(sessionAffinityFlag), &body.SessionAffinity); err != nil {
					return fmt.Errorf("--session-affinity: %w", err)
				}
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			if cmd.Flags().Changed("target-mode") {
				body.TargetMode = &targetModeFlag
			}
			if cmd.Flags().Changed("target-type") {
				body.TargetType = &targetTypeFlag
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateTargetGroup(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&healthCheckFlag, "health-check", "", "Health check (JSON)")
	f.StringVar(&instancePoolFlag, "instance-pool", "", "Compute instance pool to draw backends from")
	f.StringVar(&body.Name, "name", "", "Resource names must not start with the literal crn: prefix or be UUIDs (canonical, compact, braced, or urn:uuid: forms, in either case)")
	f.IntVar(&body.Port, "port", 0, "Port")
	f.StringVar(&body.Protocol, "protocol", "", "Protocol (one of: http, https, tcp, udp)")
	f.BoolVar(&proxyProtocolFlag, "proxy-protocol", false, "Proxy protocol")
	f.StringVar(&sessionAffinityFlag, "session-affinity", "", "Session affinity (JSON)")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&targetModeFlag, "target-mode", "", "static (the default) takes the backends you attach as targets (one of: static, pool)")
	f.StringVar(&targetTypeFlag, "target-type", "", "Target type (one of: ip, instance, function)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newLoadbalancerTargetGroupUpdateCommand builds `basaltic loadbalancer target-group update`.
func newLoadbalancerTargetGroupUpdateCommand(state *cli.State) *cobra.Command {
	var body loadbalancer.UpdateTargetGroupRequest
	var bodyFile string
	var healthCheckFlag string
	var proxyProtocolFlag bool
	var sessionAffinityFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update target group health checks, framing, or stickiness",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"health_check":     "health-check",
				"proxy_protocol":   "proxy-protocol",
				"session_affinity": "session-affinity",
				"tags":             "tags",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/target-groups/{id}")
			if err != nil {
				return err
			}
			if healthCheckFlag != "" {
				if err := json.Unmarshal([]byte(healthCheckFlag), &body.HealthCheck); err != nil {
					return fmt.Errorf("--health-check: %w", err)
				}
			}
			if cmd.Flags().Changed("proxy-protocol") {
				body.ProxyProtocol = &proxyProtocolFlag
			}
			if sessionAffinityFlag != "" {
				if err := json.Unmarshal([]byte(sessionAffinityFlag), &body.SessionAffinity); err != nil {
					return fmt.Errorf("--session-affinity: %w", err)
				}
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateTargetGroup(cmd.Context(), args[0], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&healthCheckFlag, "health-check", "", "Health check (JSON)")
	f.BoolVar(&proxyProtocolFlag, "proxy-protocol", false, "Toggle PROXY v2 framing on upstream connections")
	f.StringVar(&sessionAffinityFlag, "session-affinity", "", "Replaces the stickiness config (JSON)")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newLoadbalancerTargetGroupDeleteCommand builds `basaltic loadbalancer target-group delete`.
func newLoadbalancerTargetGroupDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a target group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/target-groups/{id}")
			if err != nil {
				return err
			}
			if err := c.DeleteTargetGroup(cmd.Context(), args[0]); err != nil {
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

// newLoadbalancerTargetGroupAttachTargetCommand builds `basaltic loadbalancer target-group attach-target`.
func newLoadbalancerTargetGroupAttachTargetCommand(state *cli.State) *cobra.Command {
	var body loadbalancer.AttachTargetRequest
	var bodyFile string
	var portFlag int
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "attach-target <id>",
		Short: "Attach a target to this group",
		Args:  cobra.ExactArgs(1),
		Long:  "Attach a target to this group.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"port":   "port",
				"target": "target",
			}, []string{"target"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/target-groups/{id}/targets")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("port") {
				body.Port = &portFlag
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.AttachTarget(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.IntVar(&portFlag, "port", 0, "Port")
	f.StringVar(&body.Target, "target", "", "Must match the group's target_type: an IP address for ip, a compute instance UUID, CRN or exact account-scoped name for instance")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newLoadbalancerTargetGroupDetachTargetCommand builds `basaltic loadbalancer target-group detach-target`.
func newLoadbalancerTargetGroupDetachTargetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach-target <id> <target-id>",
		Short: "Detach a target",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/target-groups/{id}/targets/{target_id}")
			if err != nil {
				return err
			}
			if err := c.DetachTarget(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			state.Printer().Done("Detach target requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newLoadbalancerTargetGroupGetTargetCommand builds `basaltic loadbalancer target-group get-target`.
func newLoadbalancerTargetGroupGetTargetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-target <id> <target-id>",
		Short: "Get a target",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/target-groups/{id}/targets/{target_id}")
			if err != nil {
				return err
			}
			out, err := c.GetTarget(cmd.Context(), args[0], args[1])
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

// newLoadbalancerTargetGroupListTargetsCommand builds `basaltic loadbalancer target-group list-targets`.
func newLoadbalancerTargetGroupListTargetsCommand(state *cli.State) *cobra.Command {
	var params loadbalancer.ListTargetsParams
	cmd := &cobra.Command{
		Use:   "list-targets <id>",
		Short: "List targets in this group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state, "/v1/target-groups/{id}/targets")
			if err != nil {
				return err
			}
			page, err := c.ListTargets(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact scoped CRN")
	f.StringVar(&params.Name, "name", "", "Exact immutable name")
	return cmd
}
