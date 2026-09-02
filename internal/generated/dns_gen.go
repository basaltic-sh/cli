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
	"github.com/basaltic-sh/sdk-go/dns"

	"github.com/basaltic-sh/cli/internal/cli"
)

func init() { cli.RegisterService(newDnsCommand) }

// newDnsCommand builds `basaltic dns`.
func newDnsCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dns",
		Short: "DNS zones and records",
	}
	cmd.AddCommand(newDnsRecordCommand(state))
	cmd.AddCommand(newDnsZoneCommand(state))
	return cmd
}

// dnsClient builds the service client, resolving credentials on first use.
func dnsClient(state *cli.State) (*dns.Client, error) {
	cfg, err := state.SDK()
	if err != nil {
		return nil, err
	}
	return dns.New(cfg), nil
}

// newDnsRecordCommand builds `basaltic dns record`.
func newDnsRecordCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "record",
		Short:   "Records",
		Aliases: []string{"records"},
	}
	cmd.AddCommand(newDnsRecordListCommand(state))
	cmd.AddCommand(newDnsRecordGetCommand(state))
	cmd.AddCommand(newDnsRecordCreateCommand(state))
	cmd.AddCommand(newDnsRecordUpdateCommand(state))
	cmd.AddCommand(newDnsRecordDeleteCommand(state))
	return cmd
}

// newDnsRecordListCommand builds `basaltic dns record list`.
func newDnsRecordListCommand(state *cli.State) *cobra.Command {
	var params dns.ListRecordsParams
	var includeManagedFlag bool
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list <zone-id>",
		Short: "List records",
		Args:  cobra.ExactArgs(1),
		Long:  "List records.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := dnsClient(state)
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("include-managed") {
				params.IncludeManaged = &includeManagedFlag
			}
			if fetchAll {
				return state.Printer().Iter(c.ListRecordsAll(cmd.Context(), args[0], &params))
			}
			page, err := c.ListRecords(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.BoolVar(&includeManagedFlag, "include-managed", false, "Include the platform-stamped rows (SOA, apex NS, and the DNSSEC set) alongside your own")
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Resume token — the last record id from the previous page")
	f.StringVar(&params.Name, "name", "", "Substring match on the record name")
	f.StringVar(&params.Type, "type", "", "Exact record type to filter by (e.g")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newDnsRecordGetCommand builds `basaltic dns record get`.
func newDnsRecordGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <zone-id> <record-id>",
		Short: "Get record",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := dnsClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetRecord(cmd.Context(), args[0], args[1])
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

// newDnsRecordCreateCommand builds `basaltic dns record create`.
func newDnsRecordCreateCommand(state *cli.State) *cobra.Command {
	var body dns.RecordCreateRequest
	var bodyFile string
	var ttlFlag int
	var valuesFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create <zone-id>",
		Short: "Create record",
		Args:  cobra.ExactArgs(1),
		Long:  "Create record.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := dnsClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("ttl") {
				body.TTL = &ttlFlag
			}
			if valuesFlag != "" {
				if err := json.Unmarshal([]byte(valuesFlag), &body.Values); err != nil {
					return fmt.Errorf("--values: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateRecord(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.Name, "name", "", "Record name (FQDN)")
	_ = cmd.MarkFlagRequired("name")
	f.IntVar(&ttlFlag, "ttl", 0, "Ttl")
	f.StringVar((*string)(&body.Type), "type", "", "Type (one of: A, AAAA, AFSDB, APL, CAA, CERT, CNAME, CSYNC, DHCID, DNAME, EUI48, EUI64, HINFO, HTTPS, IPSECKEY, KX, L32, L64, LOC, LP, MX, NAPTR, NID, NS, OPENPGPKEY, PTR, RKEY, RP, SMIMEA, SPF, SRV, SSHFP, SVCB, TLSA, TXT, URI)")
	_ = cmd.MarkFlagRequired("type")
	f.StringVar(&valuesFlag, "values", "", "Values (JSON)")
	_ = cmd.MarkFlagRequired("values")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newDnsRecordUpdateCommand builds `basaltic dns record update`.
func newDnsRecordUpdateCommand(state *cli.State) *cobra.Command {
	var body dns.RecordUpdateRequest
	var bodyFile string
	var ttlFlag int
	var valuesFlag string
	cmd := &cobra.Command{
		Use:   "update <zone-id> <record-id>",
		Short: "Update record",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := dnsClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("ttl") {
				body.TTL = &ttlFlag
			}
			if valuesFlag != "" {
				if err := json.Unmarshal([]byte(valuesFlag), &body.Values); err != nil {
					return fmt.Errorf("--values: %w", err)
				}
			}
			out, err := c.UpdateRecord(cmd.Context(), args[0], args[1], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.IntVar(&ttlFlag, "ttl", 0, "Ttl")
	f.StringVar(&valuesFlag, "values", "", "Values (JSON)")
	return cmd
}

// newDnsRecordDeleteCommand builds `basaltic dns record delete`.
func newDnsRecordDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <zone-id> <record-id>",
		Short: "Delete record",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := dnsClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteRecord(cmd.Context(), args[0], args[1]); err != nil {
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

// newDnsZoneCommand builds `basaltic dns zone`.
func newDnsZoneCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "zone",
		Short:   "Zones",
		Aliases: []string{"zones"},
	}
	cmd.AddCommand(newDnsZoneListCommand(state))
	cmd.AddCommand(newDnsZoneGetCommand(state))
	cmd.AddCommand(newDnsZoneCreateCommand(state))
	cmd.AddCommand(newDnsZoneUpdateCommand(state))
	cmd.AddCommand(newDnsZoneDeleteCommand(state))
	cmd.AddCommand(newDnsZoneAssociateVpcAssociationCommand(state))
	cmd.AddCommand(newDnsZoneDissociateVpcAssociationCommand(state))
	cmd.AddCommand(newDnsZoneExportCommand(state))
	cmd.AddCommand(newDnsZoneGetRecordImportCommand(state))
	cmd.AddCommand(newDnsZoneImportCommand(state))
	cmd.AddCommand(newDnsZoneListVpcAssociationsCommand(state))
	cmd.AddCommand(newDnsZoneVerifyOwnershipCommand(state))
	return cmd
}

// newDnsZoneListCommand builds `basaltic dns zone list`.
func newDnsZoneListCommand(state *cli.State) *cobra.Command {
	var params dns.ListZonesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List zones",
		Args:  cobra.ExactArgs(0),
		Long:  "List zones.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := dnsClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListZonesAll(cmd.Context(), &params))
			}
			page, err := c.ListZones(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Resume token — the last zone id from the previous page")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newDnsZoneGetCommand builds `basaltic dns zone get`.
func newDnsZoneGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <zone-id>",
		Short: "Get zone",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := dnsClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetZone(cmd.Context(), args[0])
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

// newDnsZoneCreateCommand builds `basaltic dns zone create`.
func newDnsZoneCreateCommand(state *cli.State) *cobra.Command {
	var body dns.ZoneCreateRequest
	var bodyFile string
	var descriptionFlag string
	var dnssecFlag bool
	var importExistingRecordsFlag bool
	var tagsFlag string
	var visibilityFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create zone",
		Args:  cobra.ExactArgs(0),
		Long:  "Create zone.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := dnsClient(state)
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
			if cmd.Flags().Changed("dnssec") {
				body.DNSSEC = &dnssecFlag
			}
			if cmd.Flags().Changed("import-existing-records") {
				body.ImportExistingRecords = &importExistingRecordsFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			if cmd.Flags().Changed("visibility") {
				body.Visibility = &visibilityFlag
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateZone(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&descriptionFlag, "description", "", "Free-form note stored with the zone")
	f.BoolVar(&dnssecFlag, "dnssec", false, "Sign the zone with DNSSEC")
	f.BoolVar(&importExistingRecordsFlag, "import-existing-records", false, "Read the domain's records from the nameservers that serve it TODAY and copy them into this zone, before you move the delegation here")
	f.StringVar(&body.Name, "name", "", "Zone FQDN")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&visibilityFlag, "visibility", "", "private restricts the zone to the VPCs named in vpc_ids and requires at least one; public (the default) rejects vpc_ids outright rather than ignoring them (one of: public, private)")
	f.StringSliceVar(&body.VPCIDs, "vpc-ids", nil, "VPCs the zone resolves in")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newDnsZoneUpdateCommand builds `basaltic dns zone update`.
func newDnsZoneUpdateCommand(state *cli.State) *cobra.Command {
	var body dns.ZoneUpdateRequest
	var bodyFile string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <zone-id>",
		Short: "Update zone",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := dnsClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateZone(cmd.Context(), args[0], &body)
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

// newDnsZoneDeleteCommand builds `basaltic dns zone delete`.
func newDnsZoneDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <zone-id>",
		Short: "Delete zone",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := dnsClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteZone(cmd.Context(), args[0]); err != nil {
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

// newDnsZoneAssociateVpcAssociationCommand builds `basaltic dns zone associate-vpc-association`.
func newDnsZoneAssociateVpcAssociationCommand(state *cli.State) *cobra.Command {
	var body dns.VPCAssociationRequest
	var bodyFile string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "associate-vpc-association <zone-id>",
		Short: "Associate a VPC with a private zone",
		Args:  cobra.ExactArgs(1),
		Long:  "Associate a VPC with a private zone.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := dnsClient(state)
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
			out, err := c.AssociateZoneVPC(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.VPCID, "vpc-id", "", "VPC to associate with this private zone")
	_ = cmd.MarkFlagRequired("vpc-id")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newDnsZoneDissociateVpcAssociationCommand builds `basaltic dns zone dissociate-vpc-association`.
func newDnsZoneDissociateVpcAssociationCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dissociate-vpc-association <zone-id> <vpc-id>",
		Short: "Dissociate a VPC from a private zone",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := dnsClient(state)
			if err != nil {
				return err
			}
			if err := c.DissociateZoneVPC(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			state.Printer().Done("Dissociate vpc association requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newDnsZoneExportCommand builds `basaltic dns zone export`.
func newDnsZoneExportCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export <zone-id>",
		Short: "Export the zone as a zone file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := dnsClient(state)
			if err != nil {
				return err
			}
			stream, err := c.ExportZoneFile(cmd.Context(), args[0])
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

// newDnsZoneGetRecordImportCommand builds `basaltic dns zone get-record-import`.
func newDnsZoneGetRecordImportCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-record-import <zone-id>",
		Short: "Get the record-import outcome",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := dnsClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetZoneRecordImport(cmd.Context(), args[0])
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

// newDnsZoneImportCommand builds `basaltic dns zone import`.
func newDnsZoneImportCommand(state *cli.State) *cobra.Command {
	var body dns.ZoneImportRequest
	var bodyFile string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "import <zone-id>",
		Short: "Import a zone file",
		Args:  cobra.ExactArgs(1),
		Long:  "Import a zone file.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := dnsClient(state)
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
			out, err := c.ImportZoneFile(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.ZoneFile, "zone-file", "", "The zone file, as text")
	_ = cmd.MarkFlagRequired("zone-file")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newDnsZoneListVpcAssociationsCommand builds `basaltic dns zone list-vpc-associations`.
func newDnsZoneListVpcAssociationsCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-vpc-associations <zone-id>",
		Short: "List VPC associations",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := dnsClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListZoneVPCAssociations(cmd.Context(), args[0])
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

// newDnsZoneVerifyOwnershipCommand builds `basaltic dns zone verify-ownership`.
func newDnsZoneVerifyOwnershipCommand(state *cli.State) *cobra.Command {
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "verify-ownership <zone-id>",
		Short: "Verify zone ownership",
		Args:  cobra.ExactArgs(1),
		Long:  "Verify zone ownership.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := dnsClient(state)
			if err != nil {
				return err
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.VerifyZoneOwnership(cmd.Context(), args[0], reqOpts...)
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
