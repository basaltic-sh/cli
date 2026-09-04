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
func loadbalancerClient(state *cli.State) (*loadbalancer.Client, error) {
	cfg, err := state.SDK()
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
	cmd := &cobra.Command{
		Use:   "list <id>",
		Short: "List this load balancer's listeners",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListListeners(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newLoadbalancerListenerGetCommand builds `basaltic loadbalancer listener get`.
func newLoadbalancerListenerGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id> <listener-id>",
		Short: "Get a listener",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetListener(cmd.Context(), args[0], args[1])
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
	var certificateCrnFlag string
	var certificatePemFlag string
	var certificatesFlag string
	var chainPemFlag string
	var defaultTargetGroupIdFlag string
	var exposureFlag string
	var privateKeyPemFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create <id>",
		Short: "Create a listener on this load balancer",
		Args:  cobra.ExactArgs(1),
		Long:  "Create a listener on this load balancer.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("certificate-crn") {
				body.CertificateCRN = &certificateCrnFlag
			}
			if cmd.Flags().Changed("certificate-pem") {
				body.CertificatePEM = &certificatePemFlag
			}
			if certificatesFlag != "" {
				if err := json.Unmarshal([]byte(certificatesFlag), &body.Certificates); err != nil {
					return fmt.Errorf("--certificates: %w", err)
				}
			}
			if cmd.Flags().Changed("chain-pem") {
				body.ChainPEM = &chainPemFlag
			}
			if cmd.Flags().Changed("default-target-group-id") {
				body.DefaultTargetGroupID = &defaultTargetGroupIdFlag
			}
			if cmd.Flags().Changed("exposure") {
				body.Exposure = &exposureFlag
			}
			if cmd.Flags().Changed("private-key-pem") {
				body.PrivateKeyPEM = &privateKeyPemFlag
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
	f.StringVar(&certificateCrnFlag, "certificate-crn", "", "Deprecated — use certificates[]")
	f.StringVar(&certificatePemFlag, "certificate-pem", "", "Deprecated — use certificates[]")
	f.StringVar(&certificatesFlag, "certificates", "", "Certificates (JSON)")
	f.StringVar(&chainPemFlag, "chain-pem", "", "Deprecated — use certificates[]")
	f.StringVar(&defaultTargetGroupIdFlag, "default-target-group-id", "", "Default target group id")
	f.StringVar(&exposureFlag, "exposure", "", "Which LB addresses are bound (one of: public_only, private_only, both)")
	f.IntVar(&body.Port, "port", 0, "Port")
	_ = cmd.MarkFlagRequired("port")
	f.StringVar(&privateKeyPemFlag, "private-key-pem", "", "Deprecated — use certificates[]")
	f.StringVar(&body.Protocol, "protocol", "", "Protocol (one of: http, https, tcp, udp)")
	_ = cmd.MarkFlagRequired("protocol")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newLoadbalancerListenerUpdateCommand builds `basaltic loadbalancer listener update`.
func newLoadbalancerListenerUpdateCommand(state *cli.State) *cobra.Command {
	var body loadbalancer.UpdateListenerRequest
	var bodyFile string
	var certificateCrnFlag string
	var certificatePemFlag string
	var chainPemFlag string
	var clearDefaultTargetGroupFlag bool
	var defaultTargetGroupIdFlag string
	var exposureFlag string
	var privateKeyPemFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <id> <listener-id>",
		Short: "Patch a listener (rotate cert, change default target group)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("certificate-crn") {
				body.CertificateCRN = &certificateCrnFlag
			}
			if cmd.Flags().Changed("certificate-pem") {
				body.CertificatePEM = &certificatePemFlag
			}
			if cmd.Flags().Changed("chain-pem") {
				body.ChainPEM = &chainPemFlag
			}
			if cmd.Flags().Changed("clear-default-target-group") {
				body.ClearDefaultTargetGroup = &clearDefaultTargetGroupFlag
			}
			if cmd.Flags().Changed("default-target-group-id") {
				body.DefaultTargetGroupID = &defaultTargetGroupIdFlag
			}
			if cmd.Flags().Changed("exposure") {
				body.Exposure = &exposureFlag
			}
			if cmd.Flags().Changed("private-key-pem") {
				body.PrivateKeyPEM = &privateKeyPemFlag
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
	f.StringVar(&certificateCrnFlag, "certificate-crn", "", "Certificate crn")
	f.StringVar(&certificatePemFlag, "certificate-pem", "", "Certificate pem")
	f.StringVar(&chainPemFlag, "chain-pem", "", "Chain pem")
	f.BoolVar(&clearDefaultTargetGroupFlag, "clear-default-target-group", false, "Clear default target group")
	f.StringVar(&defaultTargetGroupIdFlag, "default-target-group-id", "", "Default target group id")
	f.StringVar(&exposureFlag, "exposure", "", "Mutate which addresses are bound (one of: public_only, private_only, both)")
	f.StringVar(&privateKeyPemFlag, "private-key-pem", "", "Private key pem")
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
			c, err := loadbalancerClient(state)
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
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
	f.StringVar(&body.CertificateCRN, "certificate-crn", "", "The certificate to serve, by CRN")
	_ = cmd.MarkFlagRequired("certificate-crn")
	f.BoolVar(&isDefaultFlag, "is-default", false, "When true, demote whatever's currently default and promote this cert in the same transaction")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newLoadbalancerListenerDetachCertificateCommand builds `basaltic loadbalancer listener detach-certificate`.
func newLoadbalancerListenerDetachCertificateCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach-certificate <id> <listener-id> <certificate-crn>",
		Short: "Detach a certificate from an HTTPS listener",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state)
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
			c, err := loadbalancerClient(state)
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
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Status, "status", "", "One of: \"provisioning\", \"active\", \"error\", \"deleting\"")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newLoadbalancerLoadBalancerGetCommand builds `basaltic loadbalancer load-balancer get`.
func newLoadbalancerLoadBalancerGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a load balancer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetLoadBalancer(cmd.Context(), args[0])
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
	var floatingIpidFlag string
	var replicaCountFlag int
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a load balancer",
		Args:  cobra.ExactArgs(0),
		Long:  "Create a load balancer.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("floating-ip-id") {
				body.FloatingIPID = &floatingIpidFlag
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
	f.StringVar(&body.FlavorID, "flavor-id", "", "Compute flavor for each LB instance")
	_ = cmd.MarkFlagRequired("flavor-id")
	f.StringVar(&floatingIpidFlag, "floating-ip-id", "", "Optional FIP attached on create for public exposure")
	f.StringSliceVar(&body.KeyNames, "key-names", nil, "Platform-operator break-glass only")
	f.StringVar(&body.Name, "name", "", "1..127 chars of [A-Za-z0-9._-]")
	_ = cmd.MarkFlagRequired("name")
	f.IntVar(&replicaCountFlag, "replica-count", 0, "Number of LB compute instances")
	f.StringSliceVar(&body.SecurityGroupIDs, "security-group-ids", nil, "Security groups attached to every replica NIC (AWS ALB shape)")
	_ = cmd.MarkFlagRequired("security-group-ids")
	f.StringVar(&body.SubnetID, "subnet-id", "", "Subnet the LB instances attach to")
	_ = cmd.MarkFlagRequired("subnet-id")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&body.Type, "type", "", "Type (one of: application, network)")
	_ = cmd.MarkFlagRequired("type")
	f.StringVar(&body.VPCID, "vpc-id", "", "VPC the LB will live in")
	_ = cmd.MarkFlagRequired("vpc-id")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newLoadbalancerLoadBalancerUpdateCommand builds `basaltic loadbalancer load-balancer update`.
func newLoadbalancerLoadBalancerUpdateCommand(state *cli.State) *cobra.Command {
	var body loadbalancer.UpdateLoadBalancerRequest
	var bodyFile string
	var flavorIdFlag string
	var nameFlag string
	var replicaCountFlag int
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Rename, scale, or resize a load balancer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("flavor-id") {
				body.FlavorID = &flavorIdFlag
			}
			if cmd.Flags().Changed("name") {
				body.Name = &nameFlag
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
	f.StringVar(&flavorIdFlag, "flavor-id", "", "Resize each replica to a different compute flavor")
	f.StringVar(&nameFlag, "name", "", "Name")
	f.IntVar(&replicaCountFlag, "replica-count", 0, "Resize the set of load balancer instances")
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
			c, err := loadbalancerClient(state)
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
	cmd := &cobra.Command{
		Use:   "list-replicas <id>",
		Short: "List the LB's instance replicas with live health",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListLoadBalancerReplicas(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
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
	cmd.AddCommand(newLoadbalancerRuleDeleteOrphanedCommand(state))
	return cmd
}

// newLoadbalancerRuleListCommand builds `basaltic loadbalancer rule list`.
func newLoadbalancerRuleListCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <id> <listener-id>",
		Short: "List this listener's rules",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListRules(cmd.Context(), args[0], args[1])
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newLoadbalancerRuleGetCommand builds `basaltic loadbalancer rule get`.
func newLoadbalancerRuleGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id> <listener-id> <rule-id>",
		Short: "Get a routing rule",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetRule(cmd.Context(), args[0], args[1], args[2])
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
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
	_ = cmd.MarkFlagRequired("conditions")
	f.IntVar(&body.Priority, "priority", 0, "Priority")
	_ = cmd.MarkFlagRequired("priority")
	f.StringVar(&body.TargetGroupID, "target-group-id", "", "Target group id")
	_ = cmd.MarkFlagRequired("target-group-id")
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
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
	_ = cmd.MarkFlagRequired("conditions")
	f.IntVar(&body.Priority, "priority", 0, "Priority")
	_ = cmd.MarkFlagRequired("priority")
	f.StringVar(&body.TargetGroupID, "target-group-id", "", "Target group id")
	_ = cmd.MarkFlagRequired("target-group-id")
	return cmd
}

// newLoadbalancerRuleDeleteCommand builds `basaltic loadbalancer rule delete`.
func newLoadbalancerRuleDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <id> <listener-id> <rule-id>",
		Short: "Delete a routing rule",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state)
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

// newLoadbalancerRuleDeleteOrphanedCommand builds `basaltic loadbalancer rule delete-orphaned`.
func newLoadbalancerRuleDeleteOrphanedCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-orphaned <id> <rule-id>",
		Short: "Delete a rule (superseded)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteRule(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			state.Printer().Done("Delete orphaned requested.")
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
			c, err := loadbalancerClient(state)
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
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Protocol, "protocol", "", "One of: \"http\", \"https\", \"tcp\", \"udp\"")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newLoadbalancerTargetGroupGetCommand builds `basaltic loadbalancer target-group get`.
func newLoadbalancerTargetGroupGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a target group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetTargetGroup(cmd.Context(), args[0])
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
	var instancePoolIdFlag string
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if healthCheckFlag != "" {
				if err := json.Unmarshal([]byte(healthCheckFlag), &body.HealthCheck); err != nil {
					return fmt.Errorf("--health-check: %w", err)
				}
			}
			if cmd.Flags().Changed("instance-pool-id") {
				body.InstancePoolID = &instancePoolIdFlag
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
	f.StringVar(&instancePoolIdFlag, "instance-pool-id", "", "Compute instance pool to draw backends from")
	f.StringVar(&body.Name, "name", "", "Name")
	_ = cmd.MarkFlagRequired("name")
	f.IntVar(&body.Port, "port", 0, "Port")
	_ = cmd.MarkFlagRequired("port")
	f.StringVar(&body.Protocol, "protocol", "", "Protocol (one of: http, https, tcp, udp)")
	_ = cmd.MarkFlagRequired("protocol")
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
	var nameFlag string
	var proxyProtocolFlag bool
	var sessionAffinityFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Rename a target group or retune its health check, framing, or stickiness",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if healthCheckFlag != "" {
				if err := json.Unmarshal([]byte(healthCheckFlag), &body.HealthCheck); err != nil {
					return fmt.Errorf("--health-check: %w", err)
				}
			}
			if cmd.Flags().Changed("name") {
				body.Name = &nameFlag
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
	f.StringVar(&nameFlag, "name", "", "Name")
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
			c, err := loadbalancerClient(state)
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
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
	f.StringVar(&body.TargetRef, "target-ref", "", "Must match the group's target_type: an IP address for ip, a compute instance id for instance")
	_ = cmd.MarkFlagRequired("target-ref")
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
			c, err := loadbalancerClient(state)
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
			c, err := loadbalancerClient(state)
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
	cmd := &cobra.Command{
		Use:   "list-targets <id>",
		Short: "List targets in this group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := loadbalancerClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListTargets(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}
