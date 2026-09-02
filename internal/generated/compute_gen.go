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
	"github.com/basaltic-sh/sdk-go/compute"

	"github.com/basaltic-sh/cli/internal/cli"
)

func init() { cli.RegisterService(newComputeCommand) }

// newComputeCommand builds `basaltic compute`.
func newComputeCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compute",
		Short: "Instances, images, flavors, keypairs and pools",
		Long:  "Instances, images, flavors, keypairs and pools.\n\nThis is a regional service: it acts in the region from --region, the\nBASALTIC_REGION environment variable, or the profile.",
	}
	cmd.AddCommand(newComputeFlavorCommand(state))
	cmd.AddCommand(newComputeImageCommand(state))
	cmd.AddCommand(newComputeInstanceCommand(state))
	cmd.AddCommand(newComputeInstancePoolCommand(state))
	cmd.AddCommand(newComputeKeypairCommand(state))
	return cmd
}

// computeClient builds the service client, resolving credentials on first use.
func computeClient(state *cli.State) (*compute.Client, error) {
	cfg, err := state.SDK()
	if err != nil {
		return nil, err
	}
	return compute.New(cfg), nil
}

// newComputeFlavorCommand builds `basaltic compute flavor`.
func newComputeFlavorCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "flavor",
		Short:   "Flavors",
		Aliases: []string{"flavors"},
	}
	cmd.AddCommand(newComputeFlavorListCommand(state))
	cmd.AddCommand(newComputeFlavorGetCommand(state))
	return cmd
}

// newComputeFlavorListCommand builds `basaltic compute flavor list`.
func newComputeFlavorListCommand(state *cli.State) *cobra.Command {
	var params compute.ListFlavorsParams
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List flavors",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListFlavors(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.Family, "family", "", "Filter by product family")
	return cmd
}

// newComputeFlavorGetCommand builds `basaltic compute flavor get`.
func newComputeFlavorGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <flavor-id>",
		Short: "Get flavor",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetFlavor(cmd.Context(), args[0])
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

// newComputeImageCommand builds `basaltic compute image`.
func newComputeImageCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "image",
		Short:   "Images",
		Aliases: []string{"images"},
	}
	cmd.AddCommand(newComputeImageListCommand(state))
	cmd.AddCommand(newComputeImageGetCommand(state))
	cmd.AddCommand(newComputeImageCreateCommand(state))
	cmd.AddCommand(newComputeImageUpdateCommand(state))
	cmd.AddCommand(newComputeImageDeleteCommand(state))
	return cmd
}

// newComputeImageListCommand builds `basaltic compute image list`.
func newComputeImageListCommand(state *cli.State) *cobra.Command {
	var params compute.ListImagesParams
	var allVersionsFlag bool
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List images",
		Args:  cobra.ExactArgs(0),
		Long:  "List images.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("all-versions") {
				params.AllVersions = &allVersionsFlag
			}
			if fetchAll {
				return state.Printer().Iter(c.ListImagesAll(cmd.Context(), &params))
			}
			page, err := c.ListImages(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.BoolVar(&allVersionsFlag, "all-versions", false, "Include builds a newer version has superseded")
	f.StringVar(&params.Architecture, "architecture", "", "Architecture")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Substring match on name")
	f.StringVar(&params.OS, "os", "", "Os")
	f.StringVar(&params.Status, "status", "", "One of: \"pending\", \"importing\", \"active\", \"error\", \"hidden\"")
	f.StringVar(&params.Visibility, "visibility", "", "One of: \"public\", \"private\"")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newComputeImageGetCommand builds `basaltic compute image get`.
func newComputeImageGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <image-id>",
		Short: "Get an image",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetImage(cmd.Context(), args[0])
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

// newComputeImageCreateCommand builds `basaltic compute image create`.
func newComputeImageCreateCommand(state *cli.State) *cobra.Command {
	var body compute.ImageCreateRequest
	var bodyFile string
	var architectureFlag string
	var attributesFlag string
	var currentFlag bool
	var descriptionFlag string
	var eolDateFlag string
	var minDiskGbFlag int
	var minRammbFlag int
	var osFlag string
	var osVersionFlag string
	var tagsFlag string
	var versionFlag string
	var visibilityFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Import an image from an object URL",
		Args:  cobra.ExactArgs(0),
		Long:  "Import an image from an object URL.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("architecture") {
				body.Architecture = &architectureFlag
			}
			if attributesFlag != "" {
				if err := json.Unmarshal([]byte(attributesFlag), &body.Attributes); err != nil {
					return fmt.Errorf("--attributes: %w", err)
				}
			}
			if cmd.Flags().Changed("current") {
				body.Current = &currentFlag
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
			}
			if cmd.Flags().Changed("eol-date") {
				body.EOLDate = &eolDateFlag
			}
			if cmd.Flags().Changed("min-disk-gb") {
				body.MinDiskGB = &minDiskGbFlag
			}
			if cmd.Flags().Changed("min-ram-mb") {
				body.MinRAMMB = &minRammbFlag
			}
			if cmd.Flags().Changed("os") {
				body.OS = &osFlag
			}
			if cmd.Flags().Changed("os-version") {
				body.OSVersion = &osVersionFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			if cmd.Flags().Changed("version") {
				body.Version = &versionFlag
			}
			if cmd.Flags().Changed("visibility") {
				body.Visibility = &visibilityFlag
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateImage(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&architectureFlag, "architecture", "", "Architecture")
	f.StringVar(&attributesFlag, "attributes", "", "Attributes (JSON)")
	f.BoolVar(&currentFlag, "current", false, "Make this the current version for its (name, architecture) once active")
	f.StringVar(&descriptionFlag, "description", "", "Description")
	f.StringVar(&eolDateFlag, "eol-date", "", "The day this release stops receiving free security updates")
	f.StringVar(&body.Format, "format", "", "Source disk format at source_url; converted to the raw base on import (one of: qcow2, raw, vmdk, vhd, vhdx, vdi)")
	_ = cmd.MarkFlagRequired("format")
	f.IntVar(&minDiskGbFlag, "min-disk-gb", 0, "Min disk gb")
	f.IntVar(&minRammbFlag, "min-ram-mb", 0, "Min ram mb")
	f.StringVar(&body.Name, "name", "", "Movable tag name (e.g")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&osFlag, "os", "", "Os")
	f.StringVar(&osVersionFlag, "os-version", "", "Os version")
	f.StringVar(&body.SourceURL, "source-url", "", "Presigned https GET URL to the disk in an object store you control")
	_ = cmd.MarkFlagRequired("source-url")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&versionFlag, "version", "", "Identifies this build within name, and must be unique there — re-publishing a version that a tag already carries is a 409")
	f.StringVar(&visibilityFlag, "visibility", "", "Visibility (one of: public, private)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newComputeImageUpdateCommand builds `basaltic compute image update`.
func newComputeImageUpdateCommand(state *cli.State) *cobra.Command {
	var body compute.ImageUpdateRequest
	var bodyFile string
	var attributesFlag string
	var currentFlag bool
	var descriptionFlag string
	var eolDateFlag string
	var nameFlag string
	var tagsFlag string
	var visibilityFlag string
	cmd := &cobra.Command{
		Use:   "update <image-id>",
		Short: "Update an image's metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if attributesFlag != "" {
				if err := json.Unmarshal([]byte(attributesFlag), &body.Attributes); err != nil {
					return fmt.Errorf("--attributes: %w", err)
				}
			}
			if cmd.Flags().Changed("current") {
				body.Current = &currentFlag
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
			}
			if cmd.Flags().Changed("eol-date") {
				body.EOLDate = &eolDateFlag
			}
			if cmd.Flags().Changed("name") {
				body.Name = &nameFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			if cmd.Flags().Changed("visibility") {
				body.Visibility = &visibilityFlag
			}
			out, err := c.UpdateImage(cmd.Context(), args[0], &body)
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
	f.BoolVar(&currentFlag, "current", false, "Switch the resolve-by-name pointer for this image's name")
	f.StringVar(&descriptionFlag, "description", "", "Description")
	f.StringVar(&eolDateFlag, "eol-date", "", "Set the release's end-of-life date")
	f.StringVar(&nameFlag, "name", "", "Name")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&visibilityFlag, "visibility", "", "Visibility (one of: public, private)")
	return cmd
}

// newComputeImageDeleteCommand builds `basaltic compute image delete`.
func newComputeImageDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <image-id>",
		Short: "Delete (hide) an image",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			out, err := c.DeleteImage(cmd.Context(), args[0])
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

// newComputeInstanceCommand builds `basaltic compute instance`.
func newComputeInstanceCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "instance",
		Short:   "Instances",
		Aliases: []string{"instances"},
	}
	cmd.AddCommand(newComputeInstanceListCommand(state))
	cmd.AddCommand(newComputeInstanceGetCommand(state))
	cmd.AddCommand(newComputeInstanceCreateCommand(state))
	cmd.AddCommand(newComputeInstanceUpdateCommand(state))
	cmd.AddCommand(newComputeInstanceDeleteCommand(state))
	cmd.AddCommand(newComputeInstanceAttachNicCommand(state))
	cmd.AddCommand(newComputeInstanceAttachVolumeCommand(state))
	cmd.AddCommand(newComputeInstanceConsoleOutputCommand(state))
	cmd.AddCommand(newComputeInstanceConsoleScreenshotCommand(state))
	cmd.AddCommand(newComputeInstanceConsoleTicketCommand(state))
	cmd.AddCommand(newComputeInstanceDetachNicCommand(state))
	cmd.AddCommand(newComputeInstanceDetachVolumeCommand(state))
	cmd.AddCommand(newComputeInstanceListNicsCommand(state))
	cmd.AddCommand(newComputeInstanceListVolumesCommand(state))
	cmd.AddCommand(newComputeInstanceRebootCommand(state))
	cmd.AddCommand(newComputeInstanceReinstallCommand(state))
	cmd.AddCommand(newComputeInstanceResizeCommand(state))
	cmd.AddCommand(newComputeInstanceStartCommand(state))
	cmd.AddCommand(newComputeInstanceStopCommand(state))
	cmd.AddCommand(newComputeInstanceUpdateVolumeCommand(state))
	return cmd
}

// newComputeInstanceListCommand builds `basaltic compute instance list`.
func newComputeInstanceListCommand(state *cli.State) *cobra.Command {
	var params compute.ListInstancesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List instances",
		Args:  cobra.ExactArgs(0),
		Long:  "List instances.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListInstancesAll(cmd.Context(), &params))
			}
			page, err := c.ListInstances(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.FlavorID, "flavor-id", "", "Filter by flavor ID")
	f.StringVar(&params.ImageID, "image-id", "", "Filter by image ID")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Filter by name (exact match or prefix with *)")
	f.StringVar((*string)(&params.VMState), "vm-state", "", "Filter by lifecycle state (one of: pending, building, running, stopping, stopped, rebooting, deleting, deleted, error)")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newComputeInstanceGetCommand builds `basaltic compute instance get`.
func newComputeInstanceGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <instance-id>",
		Short: "Get instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetInstance(cmd.Context(), args[0])
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

// newComputeInstanceCreateCommand builds `basaltic compute instance create`.
func newComputeInstanceCreateCommand(state *cli.State) *cobra.Command {
	var body compute.InstanceCreateRequest
	var bodyFile string
	var assignPublicIpFlag bool
	var bootVolumeSizeGbFlag int
	var bootVolumeTypeFlag string
	var dataVolumesFlag string
	var descriptionFlag string
	var iamRoleIdFlag string
	var imageIdFlag string
	var metadataFlag string
	var networksFlag string
	var tagsFlag string
	var userDataFlag string
	var volumesFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create instance",
		Args:  cobra.ExactArgs(0),
		Long:  "Create instance.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("assign-public-ip") {
				body.AssignPublicIP = &assignPublicIpFlag
			}
			if cmd.Flags().Changed("boot-volume-size-gb") {
				body.BootVolumeSizeGB = &bootVolumeSizeGbFlag
			}
			if cmd.Flags().Changed("boot-volume-type") {
				body.BootVolumeType = &bootVolumeTypeFlag
			}
			if dataVolumesFlag != "" {
				if err := json.Unmarshal([]byte(dataVolumesFlag), &body.DataVolumes); err != nil {
					return fmt.Errorf("--data-volumes: %w", err)
				}
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
			}
			if cmd.Flags().Changed("iam-role-id") {
				body.IAMRoleID = &iamRoleIdFlag
			}
			if cmd.Flags().Changed("image-id") {
				body.ImageID = &imageIdFlag
			}
			if metadataFlag != "" {
				if err := json.Unmarshal([]byte(metadataFlag), &body.Metadata); err != nil {
					return fmt.Errorf("--metadata: %w", err)
				}
			}
			if networksFlag != "" {
				if err := json.Unmarshal([]byte(networksFlag), &body.Networks); err != nil {
					return fmt.Errorf("--networks: %w", err)
				}
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			if userDataFlag != "" {
				body.UserData = []byte(userDataFlag)
			}
			if volumesFlag != "" {
				if err := json.Unmarshal([]byte(volumesFlag), &body.Volumes); err != nil {
					return fmt.Errorf("--volumes: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateInstance(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.BoolVar(&assignPublicIpFlag, "assign-public-ip", false, "The older, instance-wide spelling of networks[0].assign_public_ip, and it means the primary NIC — the only interface it could ever have addressed")
	f.IntVar(&bootVolumeSizeGbFlag, "boot-volume-size-gb", 0, "Boot disk size cloned from the image; omitted = the image's min_disk_gb")
	f.StringVar(&bootVolumeTypeFlag, "boot-volume-type", "", "Boot disk tier; omitted = the region default (one of: ssd, nvme)")
	f.StringVar(&dataVolumesFlag, "data-volumes", "", "Blank data volumes created and bound with the instance (JSON)")
	f.StringVar(&descriptionFlag, "description", "", "Description")
	f.StringVar(&body.FlavorID, "flavor-id", "", "Flavor ID")
	_ = cmd.MarkFlagRequired("flavor-id")
	f.StringVar(&iamRoleIdFlag, "iam-role-id", "", "Attach this IAM role to the instance")
	f.StringVar(&imageIdFlag, "image-id", "", "Image to clone the boot disk from (required if not booting from volume)")
	f.StringSliceVar(&body.KeyNames, "key-names", nil, "SSH keypair names to authorize on the instance")
	f.StringVar(&metadataFlag, "metadata", "", "Metadata (JSON)")
	f.StringVar(&body.Name, "name", "", "Name")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&networksFlag, "networks", "", "Networks to attach (JSON)")
	f.StringSliceVar(&body.SecurityGroups, "security-groups", nil, "Security group names or IDs")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&userDataFlag, "user-data", "", "Base64-encoded user data (cloud-init)")
	f.StringVar(&volumesFlag, "volumes", "", "Volume attachments for boot from volume (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newComputeInstanceUpdateCommand builds `basaltic compute instance update`.
func newComputeInstanceUpdateCommand(state *cli.State) *cobra.Command {
	var body compute.InstanceUpdateRequest
	var bodyFile string
	var descriptionFlag string
	var metadataFlag string
	var nameFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <instance-id>",
		Short: "Update instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
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
			if metadataFlag != "" {
				if err := json.Unmarshal([]byte(metadataFlag), &body.Metadata); err != nil {
					return fmt.Errorf("--metadata: %w", err)
				}
			}
			if cmd.Flags().Changed("name") {
				body.Name = &nameFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateInstance(cmd.Context(), args[0], &body)
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
	f.StringVar(&metadataFlag, "metadata", "", "Metadata (JSON)")
	f.StringVar(&nameFlag, "name", "", "Name")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newComputeInstanceDeleteCommand builds `basaltic compute instance delete`.
func newComputeInstanceDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <instance-id>",
		Short: "Delete instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteInstance(cmd.Context(), args[0]); err != nil {
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

// newComputeInstanceAttachNicCommand builds `basaltic compute instance attach-nic`.
func newComputeInstanceAttachNicCommand(state *cli.State) *cobra.Command {
	var body compute.AttachInstanceNICRequest
	var bodyFile string
	var interfaceIdFlag string
	var ipAddressFlag string
	var macFlag string
	var subnetIdFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "attach-nic <instance-id>",
		Short: "Attach a NIC to a running instance",
		Args:  cobra.ExactArgs(1),
		Long:  "Attach a NIC to a running instance.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
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
			if cmd.Flags().Changed("ip-address") {
				body.IPAddress = &ipAddressFlag
			}
			if cmd.Flags().Changed("mac") {
				body.MAC = &macFlag
			}
			if cmd.Flags().Changed("subnet-id") {
				body.SubnetID = &subnetIdFlag
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.AttachInstanceNIC(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&interfaceIdFlag, "interface-id", "", "Existing standalone interface to attach")
	f.StringVar(&ipAddressFlag, "ip-address", "", "Ip address")
	f.StringVar(&macFlag, "mac", "", "Mac")
	f.StringSliceVar(&body.SecurityGroupIDs, "security-group-ids", nil, "Security group ids")
	f.StringVar(&subnetIdFlag, "subnet-id", "", "Subnet id")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newComputeInstanceAttachVolumeCommand builds `basaltic compute instance attach-volume`.
func newComputeInstanceAttachVolumeCommand(state *cli.State) *cobra.Command {
	var body compute.AttachInstanceVolumeRequest
	var bodyFile string
	var deviceFlag string
	var fstypeFlag string
	var mountPathFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "attach-volume <instance-id>",
		Short: "Attach a data volume to an instance",
		Args:  cobra.ExactArgs(1),
		Long:  "Attach a data volume to an instance.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("device") {
				body.Device = &deviceFlag
			}
			if cmd.Flags().Changed("fstype") {
				body.Fstype = &fstypeFlag
			}
			if cmd.Flags().Changed("mount-path") {
				body.MountPath = &mountPathFlag
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.AttachInstanceVolume(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&deviceFlag, "device", "", "Optional device-name override; auto-picks the next free slot (vdb/vdc/…) when omitted")
	f.StringVar(&fstypeFlag, "fstype", "", "Filesystem the in-guest agent formats the disk with, and only when mount_path is set and the disk is blank (one of: ext4, xfs)")
	f.StringVar(&mountPathFlag, "mount-path", "", "When set, the in-guest agent formats the disk (only if blank) and mounts it at this path")
	f.StringVar(&body.VolumeID, "volume-id", "", "Volume id")
	_ = cmd.MarkFlagRequired("volume-id")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newComputeInstanceConsoleOutputCommand builds `basaltic compute instance console-output`.
func newComputeInstanceConsoleOutputCommand(state *cli.State) *cobra.Command {
	var params compute.GetConsoleOutputParams
	cmd := &cobra.Command{
		Use:   "console-output <instance-id>",
		Short: "Get the instance's serial console output",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetConsoleOutput(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.MaxBytes, "max-bytes", 0, "Return at most this many bytes from the END of the transcript")
	return cmd
}

// newComputeInstanceConsoleScreenshotCommand builds `basaltic compute instance console-screenshot`.
func newComputeInstanceConsoleScreenshotCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "console-screenshot <instance-id>",
		Short: "Capture the instance's display",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			stream, err := c.GetConsoleScreenshot(cmd.Context(), args[0])
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

// newComputeInstanceConsoleTicketCommand builds `basaltic compute instance console-ticket`.
func newComputeInstanceConsoleTicketCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "console-ticket <instance-id>",
		Short: "Mint a ticket for the serial console",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			out, err := c.CreateSerialConsoleTicket(cmd.Context(), args[0])
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

// newComputeInstanceDetachNicCommand builds `basaltic compute instance detach-nic`.
func newComputeInstanceDetachNicCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach-nic <instance-id> <interface-id>",
		Short: "Detach a NIC from a running instance",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			if err := c.DetachInstanceNIC(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			state.Printer().Done("Detach nic requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newComputeInstanceDetachVolumeCommand builds `basaltic compute instance detach-volume`.
func newComputeInstanceDetachVolumeCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach-volume <instance-id> <volume-id>",
		Short: "Detach a data volume from an instance",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			if err := c.DetachInstanceVolume(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			state.Printer().Done("Detach volume requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newComputeInstanceListNicsCommand builds `basaltic compute instance list-nics`.
func newComputeInstanceListNicsCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-nics <instance-id>",
		Short: "List the instance's network interfaces",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			out, err := c.ListInstanceNiCs(cmd.Context(), args[0])
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

// newComputeInstanceListVolumesCommand builds `basaltic compute instance list-volumes`.
func newComputeInstanceListVolumesCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-volumes <instance-id>",
		Short: "List the instance's attached volumes",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			out, err := c.ListInstanceVolumes(cmd.Context(), args[0])
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

// newComputeInstanceRebootCommand builds `basaltic compute instance reboot`.
func newComputeInstanceRebootCommand(state *cli.State) *cobra.Command {
	var body compute.InstanceRebootRequest
	var bodyFile string
	var hardFlag bool
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "reboot <instance-id>",
		Short: "Reboot instance",
		Args:  cobra.ExactArgs(1),
		Long:  "Reboot instance.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("hard") {
				body.Hard = &hardFlag
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			if err := c.RebootInstance(cmd.Context(), args[0], &body, reqOpts...); err != nil {
				return err
			}
			state.Printer().Done("Reboot requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.BoolVar(&hardFlag, "hard", false, "Force a power cycle (destroy + start, equivalent to a reset button) instead of the default ACPI graceful reboot the guest can act on")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newComputeInstanceReinstallCommand builds `basaltic compute instance reinstall`.
func newComputeInstanceReinstallCommand(state *cli.State) *cobra.Command {
	var body compute.ReinstallInstanceRequest
	var bodyFile string
	var imageIdFlag string
	var sizeGbFlag int
	var volumeTypeFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "reinstall <instance-id>",
		Short: "Reinstall instance",
		Args:  cobra.ExactArgs(1),
		Long:  "Reinstall instance.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("image-id") {
				body.ImageID = &imageIdFlag
			}
			if cmd.Flags().Changed("size-gb") {
				body.SizeGB = &sizeGbFlag
			}
			if cmd.Flags().Changed("volume-type") {
				body.VolumeType = &volumeTypeFlag
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			if err := c.ReinstallInstance(cmd.Context(), args[0], &body, reqOpts...); err != nil {
				return err
			}
			state.Printer().Done("Reinstall requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&imageIdFlag, "image-id", "", "Replacement image")
	f.IntVar(&sizeGbFlag, "size-gb", 0, "Replacement boot disk size; omitted = the image's min_disk_gb")
	f.StringVar(&volumeTypeFlag, "volume-type", "", "Replacement boot disk tier; omitted = the region default (one of: ssd, nvme)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newComputeInstanceResizeCommand builds `basaltic compute instance resize`.
func newComputeInstanceResizeCommand(state *cli.State) *cobra.Command {
	var body compute.ResizeInstanceRequest
	var bodyFile string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "resize <instance-id>",
		Short: "Resize instance",
		Args:  cobra.ExactArgs(1),
		Long:  "Resize instance.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
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
			if err := c.ResizeInstance(cmd.Context(), args[0], &body, reqOpts...); err != nil {
				return err
			}
			state.Printer().Done("Resize requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.FlavorID, "flavor-id", "", "The target flavor to resize to")
	_ = cmd.MarkFlagRequired("flavor-id")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newComputeInstanceStartCommand builds `basaltic compute instance start`.
func newComputeInstanceStartCommand(state *cli.State) *cobra.Command {
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "start <instance-id>",
		Short: "Start instance",
		Args:  cobra.ExactArgs(1),
		Long:  "Start instance.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			if err := c.StartInstance(cmd.Context(), args[0], reqOpts...); err != nil {
				return err
			}
			state.Printer().Done("Start requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newComputeInstanceStopCommand builds `basaltic compute instance stop`.
func newComputeInstanceStopCommand(state *cli.State) *cobra.Command {
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "stop <instance-id>",
		Short: "Stop instance",
		Args:  cobra.ExactArgs(1),
		Long:  "Stop instance.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			if err := c.StopInstance(cmd.Context(), args[0], reqOpts...); err != nil {
				return err
			}
			state.Printer().Done("Stop requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newComputeInstanceUpdateVolumeCommand builds `basaltic compute instance update-volume`.
func newComputeInstanceUpdateVolumeCommand(state *cli.State) *cobra.Command {
	var body compute.UpdateInstanceVolumeAttachmentRequest
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "update-volume <instance-id> <volume-id>",
		Short: "Update a volume attachment's settings",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if err := c.UpdateInstanceVolumeAttachment(cmd.Context(), args[0], args[1], &body); err != nil {
				return err
			}
			state.Printer().Done("Update volume requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.BoolVar(&body.DeleteOnTermination, "delete-on-termination", false, "Delete on termination")
	_ = cmd.MarkFlagRequired("delete-on-termination")
	return cmd
}

// newComputeInstancePoolCommand builds `basaltic compute instance-pool`.
func newComputeInstancePoolCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "instance-pool",
		Short:   "Instance pools",
		Aliases: []string{"instance-pools"},
	}
	cmd.AddCommand(newComputeInstancePoolListCommand(state))
	cmd.AddCommand(newComputeInstancePoolGetCommand(state))
	cmd.AddCommand(newComputeInstancePoolCreateCommand(state))
	cmd.AddCommand(newComputeInstancePoolUpdateCommand(state))
	cmd.AddCommand(newComputeInstancePoolDeleteCommand(state))
	cmd.AddCommand(newComputeInstancePoolAttachFloatingIpCommand(state))
	cmd.AddCommand(newComputeInstancePoolDetachFloatingIpCommand(state))
	cmd.AddCommand(newComputeInstancePoolListFloatingIpsCommand(state))
	cmd.AddCommand(newComputeInstancePoolListInstancesCommand(state))
	cmd.AddCommand(newComputeInstancePoolRefreshCommand(state))
	return cmd
}

// newComputeInstancePoolListCommand builds `basaltic compute instance-pool list`.
func newComputeInstancePoolListCommand(state *cli.State) *cobra.Command {
	var params compute.ListInstancePoolsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List instance pools",
		Args:  cobra.ExactArgs(0),
		Long:  "List instance pools.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListInstancePoolsAll(cmd.Context(), &params))
			}
			page, err := c.ListInstancePools(cmd.Context(), &params)
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

// newComputeInstancePoolGetCommand builds `basaltic compute instance-pool get`.
func newComputeInstancePoolGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <pool-id>",
		Short: "Get an instance pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetInstancePool(cmd.Context(), args[0])
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

// newComputeInstancePoolCreateCommand builds `basaltic compute instance-pool create`.
func newComputeInstancePoolCreateCommand(state *cli.State) *cobra.Command {
	var body compute.InstancePoolCreateRequest
	var bodyFile string
	var assignPublicIpFlag bool
	var bootVolumeSizeGbFlag int
	var bootVolumeTypeFlag string
	var dataVolumesFlag string
	var descriptionFlag string
	var desiredCountFlag int
	var extraNiCsFlag string
	var flavorIdFlag string
	var iamRoleIdFlag string
	var imageIdFlag string
	var maxCountFlag int
	var metadataFlag string
	var minCountFlag int
	var subnetIdFlag string
	var tagsFlag string
	var templateFlag string
	var userDataFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an instance pool",
		Args:  cobra.ExactArgs(0),
		Long:  "Create an instance pool.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("assign-public-ip") {
				body.AssignPublicIP = &assignPublicIpFlag
			}
			if cmd.Flags().Changed("boot-volume-size-gb") {
				body.BootVolumeSizeGB = &bootVolumeSizeGbFlag
			}
			if cmd.Flags().Changed("boot-volume-type") {
				body.BootVolumeType = &bootVolumeTypeFlag
			}
			if dataVolumesFlag != "" {
				if err := json.Unmarshal([]byte(dataVolumesFlag), &body.DataVolumes); err != nil {
					return fmt.Errorf("--data-volumes: %w", err)
				}
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
			}
			if cmd.Flags().Changed("desired-count") {
				body.DesiredCount = &desiredCountFlag
			}
			if extraNiCsFlag != "" {
				if err := json.Unmarshal([]byte(extraNiCsFlag), &body.ExtraNICs); err != nil {
					return fmt.Errorf("--extra-nics: %w", err)
				}
			}
			if cmd.Flags().Changed("flavor-id") {
				body.FlavorID = &flavorIdFlag
			}
			if cmd.Flags().Changed("iam-role-id") {
				body.IAMRoleID = &iamRoleIdFlag
			}
			if cmd.Flags().Changed("image-id") {
				body.ImageID = &imageIdFlag
			}
			if cmd.Flags().Changed("max-count") {
				body.MaxCount = &maxCountFlag
			}
			if metadataFlag != "" {
				if err := json.Unmarshal([]byte(metadataFlag), &body.Metadata); err != nil {
					return fmt.Errorf("--metadata: %w", err)
				}
			}
			if cmd.Flags().Changed("min-count") {
				body.MinCount = &minCountFlag
			}
			if cmd.Flags().Changed("subnet-id") {
				body.SubnetID = &subnetIdFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			if templateFlag != "" {
				if err := json.Unmarshal([]byte(templateFlag), &body.Template); err != nil {
					return fmt.Errorf("--template: %w", err)
				}
			}
			if userDataFlag != "" {
				body.UserData = []byte(userDataFlag)
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateInstancePool(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.BoolVar(&assignPublicIpFlag, "assign-public-ip", false, "Superseded by template.assign_public_ip")
	f.IntVar(&bootVolumeSizeGbFlag, "boot-volume-size-gb", 0, "Superseded by template.boot_volume_size_gb")
	f.StringVar(&bootVolumeTypeFlag, "boot-volume-type", "", "Superseded by template.boot_volume_type (one of: ssd, nvme)")
	f.StringVar(&dataVolumesFlag, "data-volumes", "", "Superseded by template.data_volumes (JSON)")
	f.StringVar(&descriptionFlag, "description", "", "Description")
	f.IntVar(&desiredCountFlag, "desired-count", 0, "Desired count")
	f.StringVar(&extraNiCsFlag, "extra-nics", "", "Superseded by template.networks[1:] (JSON)")
	f.StringVar(&flavorIdFlag, "flavor-id", "", "Superseded by template.flavor_id")
	f.StringVar(&iamRoleIdFlag, "iam-role-id", "", "Superseded by template.iam_role_id")
	f.StringVar(&imageIdFlag, "image-id", "", "Superseded by template.image_id")
	f.StringSliceVar(&body.KeypairNames, "keypair-names", nil, "Superseded by template.key_names")
	f.IntVar(&maxCountFlag, "max-count", 0, "Max count")
	f.StringVar(&metadataFlag, "metadata", "", "Superseded by template.metadata (JSON)")
	f.IntVar(&minCountFlag, "min-count", 0, "Min count")
	f.StringVar(&body.Name, "name", "", "Name")
	_ = cmd.MarkFlagRequired("name")
	f.StringSliceVar(&body.SecurityGroupIDs, "security-group-ids", nil, "Superseded by template.networks[0].security_group_ids")
	f.StringVar(&subnetIdFlag, "subnet-id", "", "Superseded by template.networks[0].subnet_id")
	f.StringVar(&tagsFlag, "tags", "", "Labels on the pool resource, for IAM conditions (basalt:RequestTag/<key> here, basalt:ResourceTag/<key> on later operations) and cost attribution (JSON)")
	f.StringVar(&templateFlag, "template", "", "Template (JSON)")
	f.StringVar(&userDataFlag, "user-data", "", "Superseded by template.user_data")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newComputeInstancePoolUpdateCommand builds `basaltic compute instance-pool update`.
func newComputeInstancePoolUpdateCommand(state *cli.State) *cobra.Command {
	var body compute.InstancePoolUpdateRequest
	var bodyFile string
	var desiredCountFlag int
	var tagsFlag string
	var templateFlag string
	cmd := &cobra.Command{
		Use:   "update <pool-id>",
		Short: "Update an instance pool's size, tags or launch template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("desired-count") {
				body.DesiredCount = &desiredCountFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			if templateFlag != "" {
				if err := json.Unmarshal([]byte(templateFlag), &body.Template); err != nil {
					return fmt.Errorf("--template: %w", err)
				}
			}
			out, err := c.UpdateInstancePool(cmd.Context(), args[0], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.IntVar(&desiredCountFlag, "desired-count", 0, "New target size, bounded by the pool's min_count/max_count and the hard platform cap of 100")
	f.StringVar(&tagsFlag, "tags", "", "REPLACES the pool's labels: the map you send becomes the whole set, an empty object clears them, and omitting the field leaves them alone (JSON)")
	f.StringVar(&templateFlag, "template", "", "Replaces the launch config WHOLESALE — the object you send is what the pool launches next, and anything you leave out is cleared rather than kept (JSON)")
	return cmd
}

// newComputeInstancePoolDeleteCommand builds `basaltic compute instance-pool delete`.
func newComputeInstancePoolDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <pool-id>",
		Short: "Delete an instance pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteInstancePool(cmd.Context(), args[0]); err != nil {
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

// newComputeInstancePoolAttachFloatingIpCommand builds `basaltic compute instance-pool attach-floating-ip`.
func newComputeInstancePoolAttachFloatingIpCommand(state *cli.State) *cobra.Command {
	var body compute.InstancePoolFloatingIPAttachRequest
	var bodyFile string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "attach-floating-ip <pool-id>",
		Short: "Give the pool a shared public address",
		Args:  cobra.ExactArgs(1),
		Long:  "Give the pool a shared public address.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
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
			out, err := c.AttachInstancePoolFloatingIP(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.FloatingIPID, "floating-ip-id", "", "An already-allocated floating IP of yours, currently attached to nothing")
	_ = cmd.MarkFlagRequired("floating-ip-id")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newComputeInstancePoolDetachFloatingIpCommand builds `basaltic compute instance-pool detach-floating-ip`.
func newComputeInstancePoolDetachFloatingIpCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach-floating-ip <pool-id> <floating-ip-id>",
		Short: "Take a shared address off the pool",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			if err := c.DetachInstancePoolFloatingIP(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			state.Printer().Done("Detach floating ip requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newComputeInstancePoolListFloatingIpsCommand builds `basaltic compute instance-pool list-floating-ips`.
func newComputeInstancePoolListFloatingIpsCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-floating-ips <pool-id>",
		Short: "List the pool's shared public addresses",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListInstancePoolFloatingIPs(cmd.Context(), args[0])
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

// newComputeInstancePoolListInstancesCommand builds `basaltic compute instance-pool list-instances`.
func newComputeInstancePoolListInstancesCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-instances <pool-id>",
		Short: "List a pool's instance bindings",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListPoolInstances(cmd.Context(), args[0])
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

// newComputeInstancePoolRefreshCommand builds `basaltic compute instance-pool refresh`.
func newComputeInstancePoolRefreshCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refresh <pool-id>",
		Short: "Roll every member onto the pool's current launch template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			out, err := c.RefreshInstancePool(cmd.Context(), args[0])
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

// newComputeKeypairCommand builds `basaltic compute keypair`.
func newComputeKeypairCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keypair",
		Short:   "Keypairs",
		Aliases: []string{"keypairs"},
	}
	cmd.AddCommand(newComputeKeypairListCommand(state))
	cmd.AddCommand(newComputeKeypairGetCommand(state))
	cmd.AddCommand(newComputeKeypairCreateCommand(state))
	cmd.AddCommand(newComputeKeypairDeleteCommand(state))
	return cmd
}

// newComputeKeypairListCommand builds `basaltic compute keypair list`.
func newComputeKeypairListCommand(state *cli.State) *cobra.Command {
	var params compute.ListKeypairsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List keypairs",
		Args:  cobra.ExactArgs(0),
		Long:  "List keypairs.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListKeypairsAll(cmd.Context(), &params))
			}
			page, err := c.ListKeypairs(cmd.Context(), &params)
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

// newComputeKeypairGetCommand builds `basaltic compute keypair get`.
func newComputeKeypairGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <keypair-id>",
		Short: "Get keypair",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetKeypair(cmd.Context(), args[0])
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

// newComputeKeypairCreateCommand builds `basaltic compute keypair create`.
func newComputeKeypairCreateCommand(state *cli.State) *cobra.Command {
	var body compute.KeypairCreateRequest
	var bodyFile string
	var publicKeyFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create keypair",
		Args:  cobra.ExactArgs(0),
		Long:  "Create keypair.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("public-key") {
				body.PublicKey = &publicKeyFlag
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
			out, err := c.CreateKeypair(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.Name, "name", "", "Name")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&publicKeyFlag, "public-key", "", "SSH public key (if not provided, a new keypair will be generated)")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newComputeKeypairDeleteCommand builds `basaltic compute keypair delete`.
func newComputeKeypairDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <keypair-id>",
		Short: "Delete keypair",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := computeClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteKeypair(cmd.Context(), args[0]); err != nil {
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
