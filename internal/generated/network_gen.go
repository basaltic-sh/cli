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
	"github.com/basaltic-sh/sdk-go/network"

	"github.com/basaltic-sh/cli/internal/cli"
)

func init() { cli.RegisterService(newNetworkCommand) }

// newNetworkCommand builds `basaltic network`.
func newNetworkCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "network",
		Short:   "VPCs, subnets, gateways, routes and security groups",
		Aliases: []string{"net"},
		Long:    "VPCs, subnets, gateways, routes and security groups.\n\nThis is a regional service: it acts in the region from --region, the\nBASALTIC_REGION environment variable, or the profile.",
	}
	cmd.AddCommand(newNetworkEgressOnlyGatewayCommand(state))
	cmd.AddCommand(newNetworkFloatingIpCommand(state))
	cmd.AddCommand(newNetworkInterfaceCommand(state))
	cmd.AddCommand(newNetworkInternetGatewayCommand(state))
	cmd.AddCommand(newNetworkNatGatewayCommand(state))
	cmd.AddCommand(newNetworkRouteCommand(state))
	cmd.AddCommand(newNetworkRouteTableCommand(state))
	cmd.AddCommand(newNetworkSecurityGroupCommand(state))
	cmd.AddCommand(newNetworkSecurityGroupRuleCommand(state))
	cmd.AddCommand(newNetworkSubnetCommand(state))
	cmd.AddCommand(newNetworkVpcCommand(state))
	return cmd
}

// networkClient builds the service client, resolving credentials on first use.
func networkClient(state *cli.State) (*network.Client, error) {
	cfg, err := state.SDK()
	if err != nil {
		return nil, err
	}
	return network.New(cfg), nil
}

// newNetworkEgressOnlyGatewayCommand builds `basaltic network egress-only-gateway`.
func newNetworkEgressOnlyGatewayCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "egress-only-gateway",
		Short:   "Egress only gateways",
		Aliases: []string{"egress-only-gateways"},
	}
	cmd.AddCommand(newNetworkEgressOnlyGatewayListCommand(state))
	cmd.AddCommand(newNetworkEgressOnlyGatewayGetCommand(state))
	cmd.AddCommand(newNetworkEgressOnlyGatewayCreateCommand(state))
	cmd.AddCommand(newNetworkEgressOnlyGatewayUpdateCommand(state))
	cmd.AddCommand(newNetworkEgressOnlyGatewayDeleteCommand(state))
	return cmd
}

// newNetworkEgressOnlyGatewayListCommand builds `basaltic network egress-only-gateway list`.
func newNetworkEgressOnlyGatewayListCommand(state *cli.State) *cobra.Command {
	var params network.ListEgressOnlyGatewaysParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List egress-only gateways",
		Args:  cobra.ExactArgs(0),
		Long:  "List egress-only gateways.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListEgressOnlyGatewaysAll(cmd.Context(), &params))
			}
			page, err := c.ListEgressOnlyGateways(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Resume token — the last id from the previous page")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newNetworkEgressOnlyGatewayGetCommand builds `basaltic network egress-only-gateway get`.
func newNetworkEgressOnlyGatewayGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <egress-only-gateway-id>",
		Short: "Get egress-only gateway",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetEgressOnlyGateway(cmd.Context(), args[0])
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

// newNetworkEgressOnlyGatewayCreateCommand builds `basaltic network egress-only-gateway create`.
func newNetworkEgressOnlyGatewayCreateCommand(state *cli.State) *cobra.Command {
	var body network.EgressOnlyGatewayCreateRequest
	var bodyFile string
	var descriptionFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create egress-only gateway",
		Args:  cobra.ExactArgs(0),
		Long:  "Create egress-only gateway.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
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
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateEgressOnlyGateway(cmd.Context(), &body, reqOpts...)
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
	f.StringVar(&body.Name, "name", "", "Name")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&body.VPCID, "vpc-id", "", "The VPC must have an IPv6 CIDR (an egress-only gateway only routes v6)")
	_ = cmd.MarkFlagRequired("vpc-id")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newNetworkEgressOnlyGatewayUpdateCommand builds `basaltic network egress-only-gateway update`.
func newNetworkEgressOnlyGatewayUpdateCommand(state *cli.State) *cobra.Command {
	var body network.EgressOnlyGatewayUpdateRequest
	var bodyFile string
	var descriptionFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <egress-only-gateway-id>",
		Short: "Update egress-only gateway",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
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
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateEgressOnlyGateway(cmd.Context(), args[0], &body)
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
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newNetworkEgressOnlyGatewayDeleteCommand builds `basaltic network egress-only-gateway delete`.
func newNetworkEgressOnlyGatewayDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <egress-only-gateway-id>",
		Short: "Delete egress-only gateway",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteEgressOnlyGateway(cmd.Context(), args[0]); err != nil {
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

// newNetworkFloatingIpCommand builds `basaltic network floating-ip`.
func newNetworkFloatingIpCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "floating-ip",
		Short:   "Floating ips",
		Aliases: []string{"floating-ips"},
	}
	cmd.AddCommand(newNetworkFloatingIpListCommand(state))
	cmd.AddCommand(newNetworkFloatingIpGetCommand(state))
	cmd.AddCommand(newNetworkFloatingIpCreateCommand(state))
	cmd.AddCommand(newNetworkFloatingIpUpdateCommand(state))
	cmd.AddCommand(newNetworkFloatingIpDeleteCommand(state))
	cmd.AddCommand(newNetworkFloatingIpAttachCommand(state))
	cmd.AddCommand(newNetworkFloatingIpDetachCommand(state))
	return cmd
}

// newNetworkFloatingIpListCommand builds `basaltic network floating-ip list`.
func newNetworkFloatingIpListCommand(state *cli.State) *cobra.Command {
	var params network.ListFloatingIPsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List floating IPs",
		Args:  cobra.ExactArgs(0),
		Long:  "List floating IPs.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListFloatingIPsAll(cmd.Context(), &params))
			}
			page, err := c.ListFloatingIPs(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Resume token — the last id from the previous page")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newNetworkFloatingIpGetCommand builds `basaltic network floating-ip get`.
func newNetworkFloatingIpGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <floating-ip-id>",
		Short: "Get floating IP",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetFloatingIP(cmd.Context(), args[0])
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

// newNetworkFloatingIpCreateCommand builds `basaltic network floating-ip create`.
func newNetworkFloatingIpCreateCommand(state *cli.State) *cobra.Command {
	var body network.FloatingIPCreateRequest
	var bodyFile string
	var descriptionFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Allocate floating IP",
		Args:  cobra.ExactArgs(0),
		Long:  "Allocate floating IP.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
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
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateFloatingIP(cmd.Context(), &body, reqOpts...)
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
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newNetworkFloatingIpUpdateCommand builds `basaltic network floating-ip update`.
func newNetworkFloatingIpUpdateCommand(state *cli.State) *cobra.Command {
	var body network.FloatingIPUpdateRequest
	var bodyFile string
	var descriptionFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <floating-ip-id>",
		Short: "Update floating IP",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
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
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateFloatingIP(cmd.Context(), args[0], &body)
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
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newNetworkFloatingIpDeleteCommand builds `basaltic network floating-ip delete`.
func newNetworkFloatingIpDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <floating-ip-id>",
		Short: "Release floating IP",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteFloatingIP(cmd.Context(), args[0]); err != nil {
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

// newNetworkFloatingIpAttachCommand builds `basaltic network floating-ip attach`.
func newNetworkFloatingIpAttachCommand(state *cli.State) *cobra.Command {
	var body network.AttachFloatingIPRequest
	var bodyFile string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "attach <floating-ip-id>",
		Short: "Attach a floating IP to an interface",
		Args:  cobra.ExactArgs(1),
		Long:  "Attach a floating IP to an interface.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.AttachFloatingIP(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.InterfaceID, "interface-id", "", "Interface id")
	_ = cmd.MarkFlagRequired("interface-id")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newNetworkFloatingIpDetachCommand builds `basaltic network floating-ip detach`.
func newNetworkFloatingIpDetachCommand(state *cli.State) *cobra.Command {
	var body network.DetachFloatingIPRequest
	var bodyFile string
	var interfaceIdFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "detach <floating-ip-id>",
		Short: "Detach a floating IP",
		Args:  cobra.ExactArgs(1),
		Long:  "Detach a floating IP.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("interface-id") {
				body.InterfaceID = &interfaceIdFlag
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.DetachFloatingIP(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&interfaceIdFlag, "interface-id", "", "The member to remove")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newNetworkInterfaceCommand builds `basaltic network interface`.
func newNetworkInterfaceCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "interface",
		Short:   "Interfaces",
		Aliases: []string{"interfaces"},
	}
	cmd.AddCommand(newNetworkInterfaceListCommand(state))
	cmd.AddCommand(newNetworkInterfaceGetCommand(state))
	cmd.AddCommand(newNetworkInterfaceCreateCommand(state))
	cmd.AddCommand(newNetworkInterfaceUpdateCommand(state))
	cmd.AddCommand(newNetworkInterfaceDeleteCommand(state))
	cmd.AddCommand(newNetworkInterfaceListSecurityGroupsCommand(state))
	cmd.AddCommand(newNetworkInterfaceSetSecurityGroupCommand(state))
	return cmd
}

// newNetworkInterfaceListCommand builds `basaltic network interface list`.
func newNetworkInterfaceListCommand(state *cli.State) *cobra.Command {
	var params network.ListInterfacesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List interfaces",
		Args:  cobra.ExactArgs(0),
		Long:  "List interfaces.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListInterfacesAll(cmd.Context(), &params))
			}
			page, err := c.ListInterfaces(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Resume token — the last id from the previous page")
	f.StringVar(&params.SubnetID, "subnet-id", "", "Subnet id")
	f.StringVar(&params.VPCID, "vpc-id", "", "Vpc id")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newNetworkInterfaceGetCommand builds `basaltic network interface get`.
func newNetworkInterfaceGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <interface-id>",
		Short: "Get interface",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetInterface(cmd.Context(), args[0])
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

// newNetworkInterfaceCreateCommand builds `basaltic network interface create`.
func newNetworkInterfaceCreateCommand(state *cli.State) *cobra.Command {
	var body network.InterfaceCreateRequest
	var bodyFile string
	var descriptionFlag string
	var ipAddressFlag string
	var macFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create interface",
		Args:  cobra.ExactArgs(0),
		Long:  "Create interface.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
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
			if cmd.Flags().Changed("ip-address") {
				body.IPAddress = &ipAddressFlag
			}
			if cmd.Flags().Changed("mac") {
				body.MAC = &macFlag
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
			out, err := c.CreateInterface(cmd.Context(), &body, reqOpts...)
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
	f.StringVar(&ipAddressFlag, "ip-address", "", "Defaults to the next free address in the subnet")
	f.StringVar(&macFlag, "mac", "", "Defaults to a fresh locally-administered EUI-48")
	f.StringVar(&body.Name, "name", "", "Name")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&body.SubnetID, "subnet-id", "", "Subnet id")
	_ = cmd.MarkFlagRequired("subnet-id")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newNetworkInterfaceUpdateCommand builds `basaltic network interface update`.
func newNetworkInterfaceUpdateCommand(state *cli.State) *cobra.Command {
	var body network.InterfaceUpdateRequest
	var bodyFile string
	var descriptionFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <interface-id>",
		Short: "Update interface",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
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
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateInterface(cmd.Context(), args[0], &body)
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
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newNetworkInterfaceDeleteCommand builds `basaltic network interface delete`.
func newNetworkInterfaceDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <interface-id>",
		Short: "Delete interface",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteInterface(cmd.Context(), args[0]); err != nil {
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

// newNetworkInterfaceListSecurityGroupsCommand builds `basaltic network interface list-security-groups`.
func newNetworkInterfaceListSecurityGroupsCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-security-groups <interface-id>",
		Short: "List interface security-group membership",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListInterfaceSecurityGroups(cmd.Context(), args[0])
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

// newNetworkInterfaceSetSecurityGroupCommand builds `basaltic network interface set-security-group`.
func newNetworkInterfaceSetSecurityGroupCommand(state *cli.State) *cobra.Command {
	var body network.InterfaceSecurityGroupsRequest
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "set-security-group <interface-id>",
		Short: "Set interface security-group membership",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			page, err := c.SetInterfaceSecurityGroups(cmd.Context(), args[0], &body)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringSliceVar(&body.SecurityGroupIDs, "security-group-ids", nil, "Security group ids")
	_ = cmd.MarkFlagRequired("security-group-ids")
	return cmd
}

// newNetworkInternetGatewayCommand builds `basaltic network internet-gateway`.
func newNetworkInternetGatewayCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "internet-gateway",
		Short:   "Internet gateways",
		Aliases: []string{"internet-gateways"},
	}
	cmd.AddCommand(newNetworkInternetGatewayListCommand(state))
	cmd.AddCommand(newNetworkInternetGatewayGetCommand(state))
	cmd.AddCommand(newNetworkInternetGatewayCreateCommand(state))
	cmd.AddCommand(newNetworkInternetGatewayUpdateCommand(state))
	cmd.AddCommand(newNetworkInternetGatewayDeleteCommand(state))
	cmd.AddCommand(newNetworkInternetGatewayAttachCommand(state))
	cmd.AddCommand(newNetworkInternetGatewayDetachCommand(state))
	return cmd
}

// newNetworkInternetGatewayListCommand builds `basaltic network internet-gateway list`.
func newNetworkInternetGatewayListCommand(state *cli.State) *cobra.Command {
	var params network.ListInternetGatewaysParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List internet gateways",
		Args:  cobra.ExactArgs(0),
		Long:  "List internet gateways.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListInternetGatewaysAll(cmd.Context(), &params))
			}
			page, err := c.ListInternetGateways(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Resume token — the last id from the previous page")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newNetworkInternetGatewayGetCommand builds `basaltic network internet-gateway get`.
func newNetworkInternetGatewayGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <internet-gateway-id>",
		Short: "Get internet gateway",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetInternetGateway(cmd.Context(), args[0])
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

// newNetworkInternetGatewayCreateCommand builds `basaltic network internet-gateway create`.
func newNetworkInternetGatewayCreateCommand(state *cli.State) *cobra.Command {
	var body network.InternetGatewayCreateRequest
	var bodyFile string
	var descriptionFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create internet gateway",
		Args:  cobra.ExactArgs(0),
		Long:  "Create internet gateway.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
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
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateInternetGateway(cmd.Context(), &body, reqOpts...)
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
	f.StringVar(&body.Name, "name", "", "Name")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newNetworkInternetGatewayUpdateCommand builds `basaltic network internet-gateway update`.
func newNetworkInternetGatewayUpdateCommand(state *cli.State) *cobra.Command {
	var body network.InternetGatewayUpdateRequest
	var bodyFile string
	var descriptionFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <internet-gateway-id>",
		Short: "Update internet gateway",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
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
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateInternetGateway(cmd.Context(), args[0], &body)
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
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newNetworkInternetGatewayDeleteCommand builds `basaltic network internet-gateway delete`.
func newNetworkInternetGatewayDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <internet-gateway-id>",
		Short: "Delete internet gateway",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteInternetGateway(cmd.Context(), args[0]); err != nil {
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

// newNetworkInternetGatewayAttachCommand builds `basaltic network internet-gateway attach`.
func newNetworkInternetGatewayAttachCommand(state *cli.State) *cobra.Command {
	var body network.InternetGatewayAttachRequest
	var bodyFile string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "attach <internet-gateway-id>",
		Short: "Attach internet gateway to a VPC",
		Args:  cobra.ExactArgs(1),
		Long:  "Attach internet gateway to a VPC.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.AttachInternetGateway(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.VPCID, "vpc-id", "", "VPC to attach this IGW to")
	_ = cmd.MarkFlagRequired("vpc-id")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newNetworkInternetGatewayDetachCommand builds `basaltic network internet-gateway detach`.
func newNetworkInternetGatewayDetachCommand(state *cli.State) *cobra.Command {
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "detach <internet-gateway-id>",
		Short: "Detach internet gateway from its VPC",
		Args:  cobra.ExactArgs(1),
		Long:  "Detach internet gateway from its VPC.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.DetachInternetGateway(cmd.Context(), args[0], reqOpts...)
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

// newNetworkNatGatewayCommand builds `basaltic network nat-gateway`.
func newNetworkNatGatewayCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "nat-gateway",
		Short:   "Nat gateways",
		Aliases: []string{"nat-gateways"},
	}
	cmd.AddCommand(newNetworkNatGatewayListCommand(state))
	cmd.AddCommand(newNetworkNatGatewayGetCommand(state))
	cmd.AddCommand(newNetworkNatGatewayCreateCommand(state))
	cmd.AddCommand(newNetworkNatGatewayUpdateCommand(state))
	cmd.AddCommand(newNetworkNatGatewayDeleteCommand(state))
	return cmd
}

// newNetworkNatGatewayListCommand builds `basaltic network nat-gateway list`.
func newNetworkNatGatewayListCommand(state *cli.State) *cobra.Command {
	var params network.ListNATGatewaysParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List NAT gateways",
		Args:  cobra.ExactArgs(0),
		Long:  "List NAT gateways.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListNATGatewaysAll(cmd.Context(), &params))
			}
			page, err := c.ListNATGateways(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Resume token — the last id from the previous page")
	f.StringVar(&params.SubnetID, "subnet-id", "", "Subnet id")
	f.StringVar(&params.VPCID, "vpc-id", "", "Vpc id")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newNetworkNatGatewayGetCommand builds `basaltic network nat-gateway get`.
func newNetworkNatGatewayGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <nat-gateway-id>",
		Short: "Get NAT gateway",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetNATGateway(cmd.Context(), args[0])
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

// newNetworkNatGatewayCreateCommand builds `basaltic network nat-gateway create`.
func newNetworkNatGatewayCreateCommand(state *cli.State) *cobra.Command {
	var body network.NATGatewayCreateRequest
	var bodyFile string
	var descriptionFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create NAT gateway",
		Args:  cobra.ExactArgs(0),
		Long:  "Create NAT gateway.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
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
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateNATGateway(cmd.Context(), &body, reqOpts...)
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
	f.StringVar(&body.Name, "name", "", "Name")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&body.SubnetID, "subnet-id", "", "Subnet the NAT GW lives in")
	_ = cmd.MarkFlagRequired("subnet-id")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newNetworkNatGatewayUpdateCommand builds `basaltic network nat-gateway update`.
func newNetworkNatGatewayUpdateCommand(state *cli.State) *cobra.Command {
	var body network.NATGatewayUpdateRequest
	var bodyFile string
	var descriptionFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <nat-gateway-id>",
		Short: "Update NAT gateway",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
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
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateNATGateway(cmd.Context(), args[0], &body)
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
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newNetworkNatGatewayDeleteCommand builds `basaltic network nat-gateway delete`.
func newNetworkNatGatewayDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <nat-gateway-id>",
		Short: "Delete NAT gateway",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteNATGateway(cmd.Context(), args[0]); err != nil {
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

// newNetworkRouteCommand builds `basaltic network route`.
func newNetworkRouteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "route",
		Short:   "Routes",
		Aliases: []string{"routes"},
	}
	cmd.AddCommand(newNetworkRouteListCommand(state))
	cmd.AddCommand(newNetworkRouteGetCommand(state))
	cmd.AddCommand(newNetworkRouteCreateCommand(state))
	cmd.AddCommand(newNetworkRouteUpdateCommand(state))
	cmd.AddCommand(newNetworkRouteDeleteCommand(state))
	return cmd
}

// newNetworkRouteListCommand builds `basaltic network route list`.
func newNetworkRouteListCommand(state *cli.State) *cobra.Command {
	var params network.ListRoutesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list <route-table-id>",
		Short: "List routes",
		Args:  cobra.ExactArgs(1),
		Long:  "List routes.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListRoutesAll(cmd.Context(), args[0], &params))
			}
			page, err := c.ListRoutes(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Resume token — the last id from the previous page")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newNetworkRouteGetCommand builds `basaltic network route get`.
func newNetworkRouteGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <route-table-id> <route-id>",
		Short: "Get route",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetRoute(cmd.Context(), args[0], args[1])
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

// newNetworkRouteCreateCommand builds `basaltic network route create`.
func newNetworkRouteCreateCommand(state *cli.State) *cobra.Command {
	var body network.RouteCreateRequest
	var bodyFile string
	var descriptionFlag string
	var tagsFlag string
	var targetEgressOnlyGatewayIdFlag string
	var targetInternetGatewayIdFlag string
	var targetIpFlag string
	var targetNatGatewayIdFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create <route-table-id>",
		Short: "Create route",
		Args:  cobra.ExactArgs(1),
		Long:  "Create route.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
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
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			if cmd.Flags().Changed("target-egress-only-gateway-id") {
				body.TargetEgressOnlyGatewayID = &targetEgressOnlyGatewayIdFlag
			}
			if cmd.Flags().Changed("target-internet-gateway-id") {
				body.TargetInternetGatewayID = &targetInternetGatewayIdFlag
			}
			if cmd.Flags().Changed("target-ip") {
				body.TargetIP = &targetIpFlag
			}
			if cmd.Flags().Changed("target-nat-gateway-id") {
				body.TargetNATGatewayID = &targetNatGatewayIdFlag
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateRoute(cmd.Context(), args[0], &body, reqOpts...)
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
	f.StringVar(&body.Destination, "destination", "", "Destination")
	_ = cmd.MarkFlagRequired("destination")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&targetEgressOnlyGatewayIdFlag, "target-egress-only-gateway-id", "", "IPv6-only")
	f.StringVar(&targetInternetGatewayIdFlag, "target-internet-gateway-id", "", "Target internet gateway id")
	f.StringVar(&targetIpFlag, "target-ip", "", "Unicast next hop inside this VPC's CIDR (same IP family as destination)")
	f.StringVar(&targetNatGatewayIdFlag, "target-nat-gateway-id", "", "Target nat gateway id")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newNetworkRouteUpdateCommand builds `basaltic network route update`.
func newNetworkRouteUpdateCommand(state *cli.State) *cobra.Command {
	var body network.RouteUpdateRequest
	var bodyFile string
	var descriptionFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <route-table-id> <route-id>",
		Short: "Update route",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
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
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateRoute(cmd.Context(), args[0], args[1], &body)
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
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newNetworkRouteDeleteCommand builds `basaltic network route delete`.
func newNetworkRouteDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <route-table-id> <route-id>",
		Short: "Delete route",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteRoute(cmd.Context(), args[0], args[1]); err != nil {
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

// newNetworkRouteTableCommand builds `basaltic network route-table`.
func newNetworkRouteTableCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "route-table",
		Short:   "Route tables",
		Aliases: []string{"route-tables"},
	}
	cmd.AddCommand(newNetworkRouteTableListCommand(state))
	cmd.AddCommand(newNetworkRouteTableGetCommand(state))
	cmd.AddCommand(newNetworkRouteTableCreateCommand(state))
	cmd.AddCommand(newNetworkRouteTableUpdateCommand(state))
	cmd.AddCommand(newNetworkRouteTableDeleteCommand(state))
	return cmd
}

// newNetworkRouteTableListCommand builds `basaltic network route-table list`.
func newNetworkRouteTableListCommand(state *cli.State) *cobra.Command {
	var params network.ListRouteTablesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List route tables",
		Args:  cobra.ExactArgs(0),
		Long:  "List route tables.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListRouteTablesAll(cmd.Context(), &params))
			}
			page, err := c.ListRouteTables(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Resume token — the last id from the previous page")
	f.StringVar(&params.VPCID, "vpc-id", "", "Filter by VPC ID")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newNetworkRouteTableGetCommand builds `basaltic network route-table get`.
func newNetworkRouteTableGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <route-table-id>",
		Short: "Get route table",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetRouteTable(cmd.Context(), args[0])
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

// newNetworkRouteTableCreateCommand builds `basaltic network route-table create`.
func newNetworkRouteTableCreateCommand(state *cli.State) *cobra.Command {
	var body network.RouteTableCreateRequest
	var bodyFile string
	var descriptionFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create route table",
		Args:  cobra.ExactArgs(0),
		Long:  "Create route table.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
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
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateRouteTable(cmd.Context(), &body, reqOpts...)
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
	f.StringVar(&body.Name, "name", "", "1-63 chars, lowercase alphanumeric + hyphen")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&body.VPCID, "vpc-id", "", "VPC this table belongs to")
	_ = cmd.MarkFlagRequired("vpc-id")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newNetworkRouteTableUpdateCommand builds `basaltic network route-table update`.
func newNetworkRouteTableUpdateCommand(state *cli.State) *cobra.Command {
	var body network.RouteTableUpdateRequest
	var bodyFile string
	var descriptionFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <route-table-id>",
		Short: "Update route table",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
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
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateRouteTable(cmd.Context(), args[0], &body)
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
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newNetworkRouteTableDeleteCommand builds `basaltic network route-table delete`.
func newNetworkRouteTableDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <route-table-id>",
		Short: "Delete route table",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteRouteTable(cmd.Context(), args[0]); err != nil {
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

// newNetworkSecurityGroupCommand builds `basaltic network security-group`.
func newNetworkSecurityGroupCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "security-group",
		Short:   "Security groups",
		Aliases: []string{"security-groups"},
	}
	cmd.AddCommand(newNetworkSecurityGroupListCommand(state))
	cmd.AddCommand(newNetworkSecurityGroupGetCommand(state))
	cmd.AddCommand(newNetworkSecurityGroupCreateCommand(state))
	cmd.AddCommand(newNetworkSecurityGroupUpdateCommand(state))
	cmd.AddCommand(newNetworkSecurityGroupDeleteCommand(state))
	return cmd
}

// newNetworkSecurityGroupListCommand builds `basaltic network security-group list`.
func newNetworkSecurityGroupListCommand(state *cli.State) *cobra.Command {
	var params network.ListSecurityGroupsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List security groups",
		Args:  cobra.ExactArgs(0),
		Long:  "List security groups.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListSecurityGroupsAll(cmd.Context(), &params))
			}
			page, err := c.ListSecurityGroups(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Resume token — the last id from the previous page")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newNetworkSecurityGroupGetCommand builds `basaltic network security-group get`.
func newNetworkSecurityGroupGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <security-group-id>",
		Short: "Get security group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetSecurityGroup(cmd.Context(), args[0])
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

// newNetworkSecurityGroupCreateCommand builds `basaltic network security-group create`.
func newNetworkSecurityGroupCreateCommand(state *cli.State) *cobra.Command {
	var body network.SecurityGroupCreateRequest
	var bodyFile string
	var descriptionFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create security group",
		Args:  cobra.ExactArgs(0),
		Long:  "Create security group.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
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
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateSecurityGroup(cmd.Context(), &body, reqOpts...)
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
	f.StringVar(&body.Name, "name", "", "Name")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newNetworkSecurityGroupUpdateCommand builds `basaltic network security-group update`.
func newNetworkSecurityGroupUpdateCommand(state *cli.State) *cobra.Command {
	var body network.SecurityGroupUpdateRequest
	var bodyFile string
	var descriptionFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <security-group-id>",
		Short: "Update security group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
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
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateSecurityGroup(cmd.Context(), args[0], &body)
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
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newNetworkSecurityGroupDeleteCommand builds `basaltic network security-group delete`.
func newNetworkSecurityGroupDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <security-group-id>",
		Short: "Delete security group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteSecurityGroup(cmd.Context(), args[0]); err != nil {
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

// newNetworkSecurityGroupRuleCommand builds `basaltic network security-group-rule`.
func newNetworkSecurityGroupRuleCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "security-group-rule",
		Short:   "Security group rules",
		Aliases: []string{"security-group-rules"},
	}
	cmd.AddCommand(newNetworkSecurityGroupRuleListCommand(state))
	cmd.AddCommand(newNetworkSecurityGroupRuleGetCommand(state))
	cmd.AddCommand(newNetworkSecurityGroupRuleCreateCommand(state))
	cmd.AddCommand(newNetworkSecurityGroupRuleDeleteCommand(state))
	return cmd
}

// newNetworkSecurityGroupRuleListCommand builds `basaltic network security-group-rule list`.
func newNetworkSecurityGroupRuleListCommand(state *cli.State) *cobra.Command {
	var params network.ListSecurityGroupRulesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list <security-group-id>",
		Short: "List security group rules",
		Args:  cobra.ExactArgs(1),
		Long:  "List security group rules.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListSecurityGroupRulesAll(cmd.Context(), args[0], &params))
			}
			page, err := c.ListSecurityGroupRules(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Resume token — the last id from the previous page")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newNetworkSecurityGroupRuleGetCommand builds `basaltic network security-group-rule get`.
func newNetworkSecurityGroupRuleGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <security-group-id> <rule-id>",
		Short: "Get security group rule",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetSecurityGroupRule(cmd.Context(), args[0], args[1])
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

// newNetworkSecurityGroupRuleCreateCommand builds `basaltic network security-group-rule create`.
func newNetworkSecurityGroupRuleCreateCommand(state *cli.State) *cobra.Command {
	var body network.SecurityGroupRuleCreateRequest
	var bodyFile string
	var descriptionFlag string
	var ethertypeFlag string
	var portMaxFlag int
	var portMinFlag int
	var sourceCidrFlag string
	var sourceSecurityGroupIdFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create <security-group-id>",
		Short: "Create security group rule",
		Args:  cobra.ExactArgs(1),
		Long:  "Create security group rule.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
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
			if cmd.Flags().Changed("ethertype") {
				body.Ethertype = (*network.SecurityGroupRuleEthertype)(&ethertypeFlag)
			}
			if cmd.Flags().Changed("port-max") {
				body.PortMax = &portMaxFlag
			}
			if cmd.Flags().Changed("port-min") {
				body.PortMin = &portMinFlag
			}
			if cmd.Flags().Changed("source-cidr") {
				body.SourceCIDR = &sourceCidrFlag
			}
			if cmd.Flags().Changed("source-security-group-id") {
				body.SourceSecurityGroupID = &sourceSecurityGroupIdFlag
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateSecurityGroupRule(cmd.Context(), args[0], &body, reqOpts...)
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
	f.StringVar((*string)(&body.Direction), "direction", "", "Direction (one of: ingress, egress)")
	_ = cmd.MarkFlagRequired("direction")
	f.StringVar(&ethertypeFlag, "ethertype", "", "Ethertype (one of: ipv4, ipv6)")
	f.IntVar(&portMaxFlag, "port-max", 0, "Port max")
	f.IntVar(&portMinFlag, "port-min", 0, "Port min")
	f.StringVar((*string)(&body.Protocol), "protocol", "", "Protocol (one of: tcp, udp, icmp, all)")
	_ = cmd.MarkFlagRequired("protocol")
	f.StringVar(&sourceCidrFlag, "source-cidr", "", "Source cidr")
	f.StringVar(&sourceSecurityGroupIdFlag, "source-security-group-id", "", "Source security group id")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newNetworkSecurityGroupRuleDeleteCommand builds `basaltic network security-group-rule delete`.
func newNetworkSecurityGroupRuleDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <security-group-id> <rule-id>",
		Short: "Delete security group rule",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteSecurityGroupRule(cmd.Context(), args[0], args[1]); err != nil {
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

// newNetworkSubnetCommand builds `basaltic network subnet`.
func newNetworkSubnetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "subnet",
		Short:   "Subnets",
		Aliases: []string{"subnets"},
	}
	cmd.AddCommand(newNetworkSubnetListCommand(state))
	cmd.AddCommand(newNetworkSubnetGetCommand(state))
	cmd.AddCommand(newNetworkSubnetCreateCommand(state))
	cmd.AddCommand(newNetworkSubnetUpdateCommand(state))
	cmd.AddCommand(newNetworkSubnetDeleteCommand(state))
	return cmd
}

// newNetworkSubnetListCommand builds `basaltic network subnet list`.
func newNetworkSubnetListCommand(state *cli.State) *cobra.Command {
	var params network.ListSubnetsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List subnets",
		Args:  cobra.ExactArgs(0),
		Long:  "List subnets.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListSubnetsAll(cmd.Context(), &params))
			}
			page, err := c.ListSubnets(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Resume token — the last id from the previous page")
	f.StringVar(&params.VPCID, "vpc-id", "", "Filter by VPC ID")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newNetworkSubnetGetCommand builds `basaltic network subnet get`.
func newNetworkSubnetGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <subnet-id>",
		Short: "Get subnet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetSubnet(cmd.Context(), args[0])
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

// newNetworkSubnetCreateCommand builds `basaltic network subnet create`.
func newNetworkSubnetCreateCommand(state *cli.State) *cobra.Command {
	var body network.SubnetCreateRequest
	var bodyFile string
	var cidrv6Flag string
	var descriptionFlag string
	var gatewayIpFlag string
	var routeTableIdFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create subnet",
		Args:  cobra.ExactArgs(0),
		Long:  "Create subnet.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("cidr-v6") {
				body.CIDRV6 = &cidrv6Flag
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
			}
			if cmd.Flags().Changed("gateway-ip") {
				body.GatewayIP = &gatewayIpFlag
			}
			if cmd.Flags().Changed("route-table-id") {
				body.RouteTableID = &routeTableIdFlag
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
			out, err := c.CreateSubnet(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.CIDR, "cidr", "", "Cidr")
	_ = cmd.MarkFlagRequired("cidr")
	f.StringVar(&cidrv6Flag, "cidr-v6", "", "Makes the subnet dual-stack")
	f.StringVar(&descriptionFlag, "description", "", "Description")
	f.StringVar(&gatewayIpFlag, "gateway-ip", "", "Defaults to the first usable host in the CIDR")
	f.StringVar(&body.Name, "name", "", "Name")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&routeTableIdFlag, "route-table-id", "", "Defaults to the VPC's main route table")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&body.VPCID, "vpc-id", "", "Vpc id")
	_ = cmd.MarkFlagRequired("vpc-id")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newNetworkSubnetUpdateCommand builds `basaltic network subnet update`.
func newNetworkSubnetUpdateCommand(state *cli.State) *cobra.Command {
	var body network.SubnetUpdateRequest
	var bodyFile string
	var descriptionFlag string
	var routeTableIdFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <subnet-id>",
		Short: "Update subnet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
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
			if cmd.Flags().Changed("route-table-id") {
				body.RouteTableID = &routeTableIdFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateSubnet(cmd.Context(), args[0], &body)
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
	f.StringVar(&routeTableIdFlag, "route-table-id", "", "Re-associate the subnet with a different route table")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newNetworkSubnetDeleteCommand builds `basaltic network subnet delete`.
func newNetworkSubnetDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <subnet-id>",
		Short: "Delete subnet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteSubnet(cmd.Context(), args[0]); err != nil {
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

// newNetworkVpcCommand builds `basaltic network vpc`.
func newNetworkVpcCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vpc",
		Short:   "Vpcs",
		Aliases: []string{"vpcs"},
	}
	cmd.AddCommand(newNetworkVpcListCommand(state))
	cmd.AddCommand(newNetworkVpcGetCommand(state))
	cmd.AddCommand(newNetworkVpcCreateCommand(state))
	cmd.AddCommand(newNetworkVpcUpdateCommand(state))
	cmd.AddCommand(newNetworkVpcDeleteCommand(state))
	return cmd
}

// newNetworkVpcListCommand builds `basaltic network vpc list`.
func newNetworkVpcListCommand(state *cli.State) *cobra.Command {
	var params network.ListVPCsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List VPCs",
		Args:  cobra.ExactArgs(0),
		Long:  "List VPCs.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListVPCsAll(cmd.Context(), &params))
			}
			page, err := c.ListVPCs(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Resume token — the last id from the previous page")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newNetworkVpcGetCommand builds `basaltic network vpc get`.
func newNetworkVpcGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <vpc-id>",
		Short: "Get VPC",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetVPC(cmd.Context(), args[0])
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

// newNetworkVpcCreateCommand builds `basaltic network vpc create`.
func newNetworkVpcCreateCommand(state *cli.State) *cobra.Command {
	var body network.VPCCreateRequest
	var bodyFile string
	var assignIPv6cidrFlag bool
	var descriptionFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create VPC",
		Args:  cobra.ExactArgs(0),
		Long:  "Create VPC.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("assign-ipv6-cidr") {
				body.AssignIPv6CIDR = &assignIPv6cidrFlag
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
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
			out, err := c.CreateVPC(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.BoolVar(&assignIPv6cidrFlag, "assign-ipv6-cidr", false, "Request a globally-routable /60 delegated from the region's IPv6 pool — the only way a VPC gets an IPv6 prefix; the region must have IPv6 enabled")
	f.StringVar(&body.CIDRV4, "cidr-v4", "", "Must be private (RFC 1918): within 10.0.0.0/8, 172.16.0.0/12 or 192.168.0.0/16")
	_ = cmd.MarkFlagRequired("cidr-v4")
	f.StringVar(&descriptionFlag, "description", "", "Description")
	f.StringVar(&body.Name, "name", "", "Name")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newNetworkVpcUpdateCommand builds `basaltic network vpc update`.
func newNetworkVpcUpdateCommand(state *cli.State) *cobra.Command {
	var body network.VPCUpdateRequest
	var bodyFile string
	var descriptionFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <vpc-id>",
		Short: "Update VPC",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
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
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateVPC(cmd.Context(), args[0], &body)
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
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newNetworkVpcDeleteCommand builds `basaltic network vpc delete`.
func newNetworkVpcDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <vpc-id>",
		Short: "Delete VPC",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := networkClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteVPC(cmd.Context(), args[0]); err != nil {
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
