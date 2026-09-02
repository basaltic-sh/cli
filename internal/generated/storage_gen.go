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
	"github.com/basaltic-sh/sdk-go/storage"

	"github.com/basaltic-sh/cli/internal/cli"
)

func init() { cli.RegisterService(newStorageCommand) }

// newStorageCommand builds `basaltic storage`.
func newStorageCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "storage",
		Short: "Volumes, snapshots, buckets and objects",
		Long:  "Volumes, snapshots, buckets and objects.\n\nThis is a regional service: it acts in the region from --region, the\nBASALTIC_REGION environment variable, or the profile.",
	}
	cmd.AddCommand(newStorageBucketCommand(state))
	cmd.AddCommand(newStorageMultipartUploadCommand(state))
	cmd.AddCommand(newStorageObjectCommand(state))
	cmd.AddCommand(newStorageSnapshotCommand(state))
	cmd.AddCommand(newStorageSnapshotPolicyCommand(state))
	cmd.AddCommand(newStorageVolumeCommand(state))
	cmd.AddCommand(newStorageVolumeTypeCommand(state))
	return cmd
}

// storageClient builds the service client, resolving credentials on first use.
func storageClient(state *cli.State) (*storage.Client, error) {
	cfg, err := state.SDK()
	if err != nil {
		return nil, err
	}
	return storage.New(cfg), nil
}

// newStorageBucketCommand builds `basaltic storage bucket`.
func newStorageBucketCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bucket",
		Short:   "Buckets",
		Aliases: []string{"buckets"},
	}
	cmd.AddCommand(newStorageBucketListCommand(state))
	cmd.AddCommand(newStorageBucketCreateCommand(state))
	cmd.AddCommand(newStorageBucketDeleteCommand(state))
	cmd.AddCommand(newStorageBucketDeleteCorsCommand(state))
	cmd.AddCommand(newStorageBucketDeleteEncryptionCommand(state))
	cmd.AddCommand(newStorageBucketDeleteLifecycleCommand(state))
	cmd.AddCommand(newStorageBucketDeleteObjectLockCommand(state))
	cmd.AddCommand(newStorageBucketDeletePolicyCommand(state))
	cmd.AddCommand(newStorageBucketDeleteTaggingCommand(state))
	cmd.AddCommand(newStorageBucketGetCorsCommand(state))
	cmd.AddCommand(newStorageBucketGetEncryptionCommand(state))
	cmd.AddCommand(newStorageBucketGetLifecycleCommand(state))
	cmd.AddCommand(newStorageBucketGetObjectLockCommand(state))
	cmd.AddCommand(newStorageBucketGetPolicyCommand(state))
	cmd.AddCommand(newStorageBucketGetTaggingCommand(state))
	cmd.AddCommand(newStorageBucketGetVersioningCommand(state))
	cmd.AddCommand(newStorageBucketHeadCommand(state))
	cmd.AddCommand(newStorageBucketListObjectVersionsCommand(state))
	cmd.AddCommand(newStorageBucketRestoreCommand(state))
	cmd.AddCommand(newStorageBucketSetCorsCommand(state))
	cmd.AddCommand(newStorageBucketSetDeletionProtectionCommand(state))
	cmd.AddCommand(newStorageBucketSetEncryptionCommand(state))
	cmd.AddCommand(newStorageBucketSetLifecycleCommand(state))
	cmd.AddCommand(newStorageBucketSetObjectLockCommand(state))
	cmd.AddCommand(newStorageBucketSetPolicyCommand(state))
	cmd.AddCommand(newStorageBucketSetTaggingCommand(state))
	cmd.AddCommand(newStorageBucketSetVersioningCommand(state))
	return cmd
}

// newStorageBucketListCommand builds `basaltic storage bucket list`.
func newStorageBucketListCommand(state *cli.State) *cobra.Command {
	var params storage.ListBucketsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List buckets",
		Args:  cobra.ExactArgs(0),
		Long:  "List buckets.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListBucketsAll(cmd.Context(), &params))
			}
			page, err := c.ListBuckets(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Resume token — the last bucket name from the previous page")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newStorageBucketCreateCommand builds `basaltic storage bucket create`.
func newStorageBucketCreateCommand(state *cli.State) *cobra.Command {
	var body storage.CreateBucketRequest
	var bodyFile string
	var objectLockEnabledFlag bool
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create bucket",
		Args:  cobra.ExactArgs(0),
		Long:  "Create bucket.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("object-lock-enabled") {
				body.ObjectLockEnabled = &objectLockEnabledFlag
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateBucket(cmd.Context(), &body, reqOpts...)
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
	f.BoolVar(&objectLockEnabledFlag, "object-lock-enabled", false, "When true, enables S3 Object Lock on the bucket at creation time and turns versioning on")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newStorageBucketDeleteCommand builds `basaltic storage bucket delete`.
func newStorageBucketDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <bucket>",
		Short: "Delete bucket",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteBucket(cmd.Context(), args[0]); err != nil {
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

// newStorageBucketDeleteCorsCommand builds `basaltic storage bucket delete-cors`.
func newStorageBucketDeleteCorsCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-cors <bucket>",
		Short: "Delete bucket CORS configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteBucketCORS(cmd.Context(), args[0]); err != nil {
				return err
			}
			state.Printer().Done("Delete cors requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newStorageBucketDeleteEncryptionCommand builds `basaltic storage bucket delete-encryption`.
func newStorageBucketDeleteEncryptionCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-encryption <bucket>",
		Short: "Delete bucket encryption configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteBucketEncryption(cmd.Context(), args[0]); err != nil {
				return err
			}
			state.Printer().Done("Delete encryption requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newStorageBucketDeleteLifecycleCommand builds `basaltic storage bucket delete-lifecycle`.
func newStorageBucketDeleteLifecycleCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-lifecycle <bucket>",
		Short: "Delete bucket lifecycle configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteBucketLifecycle(cmd.Context(), args[0]); err != nil {
				return err
			}
			state.Printer().Done("Delete lifecycle requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newStorageBucketDeleteObjectLockCommand builds `basaltic storage bucket delete-object-lock`.
func newStorageBucketDeleteObjectLockCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-object-lock <bucket>",
		Short: "Delete bucket object-lock configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteBucketObjectLock(cmd.Context(), args[0]); err != nil {
				return err
			}
			state.Printer().Done("Delete object lock requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newStorageBucketDeletePolicyCommand builds `basaltic storage bucket delete-policy`.
func newStorageBucketDeletePolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-policy <bucket>",
		Short: "Delete bucket policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteBucketPolicy(cmd.Context(), args[0]); err != nil {
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

// newStorageBucketDeleteTaggingCommand builds `basaltic storage bucket delete-tagging`.
func newStorageBucketDeleteTaggingCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-tagging <bucket>",
		Short: "Delete bucket tag set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteBucketTagging(cmd.Context(), args[0]); err != nil {
				return err
			}
			state.Printer().Done("Delete tagging requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newStorageBucketGetCorsCommand builds `basaltic storage bucket get-cors`.
func newStorageBucketGetCorsCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-cors <bucket>",
		Short: "Get bucket CORS configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetBucketCORS(cmd.Context(), args[0])
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

// newStorageBucketGetEncryptionCommand builds `basaltic storage bucket get-encryption`.
func newStorageBucketGetEncryptionCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-encryption <bucket>",
		Short: "Get bucket encryption configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetBucketEncryption(cmd.Context(), args[0])
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

// newStorageBucketGetLifecycleCommand builds `basaltic storage bucket get-lifecycle`.
func newStorageBucketGetLifecycleCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-lifecycle <bucket>",
		Short: "Get bucket lifecycle configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetBucketLifecycle(cmd.Context(), args[0])
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

// newStorageBucketGetObjectLockCommand builds `basaltic storage bucket get-object-lock`.
func newStorageBucketGetObjectLockCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-object-lock <bucket>",
		Short: "Get bucket object-lock configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetBucketObjectLock(cmd.Context(), args[0])
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

// newStorageBucketGetPolicyCommand builds `basaltic storage bucket get-policy`.
func newStorageBucketGetPolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-policy <bucket>",
		Short: "Get bucket policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetBucketPolicy(cmd.Context(), args[0])
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

// newStorageBucketGetTaggingCommand builds `basaltic storage bucket get-tagging`.
func newStorageBucketGetTaggingCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-tagging <bucket>",
		Short: "Get bucket tag set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetBucketTagging(cmd.Context(), args[0])
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

// newStorageBucketGetVersioningCommand builds `basaltic storage bucket get-versioning`.
func newStorageBucketGetVersioningCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-versioning <bucket>",
		Short: "Get bucket versioning state",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetBucketVersioning(cmd.Context(), args[0])
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

// newStorageBucketHeadCommand builds `basaltic storage bucket head`.
func newStorageBucketHeadCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "head <bucket>",
		Short: "Head bucket",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if err := c.HeadBucket(cmd.Context(), args[0]); err != nil {
				return err
			}
			state.Printer().Done("Head requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newStorageBucketListObjectVersionsCommand builds `basaltic storage bucket list-object-versions`.
func newStorageBucketListObjectVersionsCommand(state *cli.State) *cobra.Command {
	var params storage.ListObjectVersionsParams
	cmd := &cobra.Command{
		Use:   "list-object-versions <bucket>",
		Short: "List object versions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			out, err := c.ListObjectVersions(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.KeyMarker, "key-marker", "", "Key marker")
	f.IntVar(&params.MaxKeys, "max-keys", 0, "Max keys")
	f.StringVar(&params.Prefix, "prefix", "", "Prefix")
	f.StringVar(&params.VersionIDMarker, "version-id-marker", "", "Version id marker")
	return cmd
}

// newStorageBucketRestoreCommand builds `basaltic storage bucket restore`.
func newStorageBucketRestoreCommand(state *cli.State) *cobra.Command {
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "restore <bucket>",
		Short: "Restore a bucket pending deletion",
		Args:  cobra.ExactArgs(1),
		Long:  "Restore a bucket pending deletion.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			if err := c.RestoreBucket(cmd.Context(), args[0], reqOpts...); err != nil {
				return err
			}
			state.Printer().Done("Restore requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newStorageBucketSetCorsCommand builds `basaltic storage bucket set-cors`.
func newStorageBucketSetCorsCommand(state *cli.State) *cobra.Command {
	var body storage.PutBucketCORSRequest
	var bodyFile string
	var corsFlag string
	cmd := &cobra.Command{
		Use:   "set-cors <bucket>",
		Short: "Put bucket CORS configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if corsFlag != "" {
				if err := json.Unmarshal([]byte(corsFlag), &body.CORS); err != nil {
					return fmt.Errorf("--cors: %w", err)
				}
			}
			if err := c.PutBucketCORS(cmd.Context(), args[0], &body); err != nil {
				return err
			}
			state.Printer().Done("Set cors requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&corsFlag, "cors", "", "Cors (JSON)")
	_ = cmd.MarkFlagRequired("cors")
	return cmd
}

// newStorageBucketSetDeletionProtectionCommand builds `basaltic storage bucket set-deletion-protection`.
func newStorageBucketSetDeletionProtectionCommand(state *cli.State) *cobra.Command {
	var body storage.PutBucketDeletionProtectionRequest
	var bodyFile string
	var recoveryDaysFlag int
	cmd := &cobra.Command{
		Use:   "set-deletion-protection <bucket>",
		Short: "Set bucket deletion protection",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("recovery-days") {
				body.RecoveryDays = &recoveryDaysFlag
			}
			if err := c.PutBucketDeletionProtection(cmd.Context(), args[0], &body); err != nil {
				return err
			}
			state.Printer().Done("Set deletion protection requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.BoolVar(&body.Enabled, "enabled", false, "When true, DeleteBucket schedules deletion instead of removing immediately")
	_ = cmd.MarkFlagRequired("enabled")
	f.IntVar(&recoveryDaysFlag, "recovery-days", 0, "Recovery window when enabling protection")
	return cmd
}

// newStorageBucketSetEncryptionCommand builds `basaltic storage bucket set-encryption`.
func newStorageBucketSetEncryptionCommand(state *cli.State) *cobra.Command {
	var body storage.PutBucketEncryptionRequest
	var bodyFile string
	var encryptionFlag string
	cmd := &cobra.Command{
		Use:   "set-encryption <bucket>",
		Short: "Put bucket encryption configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if encryptionFlag != "" {
				if err := json.Unmarshal([]byte(encryptionFlag), &body.Encryption); err != nil {
					return fmt.Errorf("--encryption: %w", err)
				}
			}
			if err := c.PutBucketEncryption(cmd.Context(), args[0], &body); err != nil {
				return err
			}
			state.Printer().Done("Set encryption requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&encryptionFlag, "encryption", "", "Encryption (JSON)")
	_ = cmd.MarkFlagRequired("encryption")
	return cmd
}

// newStorageBucketSetLifecycleCommand builds `basaltic storage bucket set-lifecycle`.
func newStorageBucketSetLifecycleCommand(state *cli.State) *cobra.Command {
	var body storage.PutBucketLifecycleRequest
	var bodyFile string
	var lifecycleFlag string
	cmd := &cobra.Command{
		Use:   "set-lifecycle <bucket>",
		Short: "Put bucket lifecycle configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if lifecycleFlag != "" {
				if err := json.Unmarshal([]byte(lifecycleFlag), &body.Lifecycle); err != nil {
					return fmt.Errorf("--lifecycle: %w", err)
				}
			}
			if err := c.PutBucketLifecycle(cmd.Context(), args[0], &body); err != nil {
				return err
			}
			state.Printer().Done("Set lifecycle requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&lifecycleFlag, "lifecycle", "", "Lifecycle (JSON)")
	_ = cmd.MarkFlagRequired("lifecycle")
	return cmd
}

// newStorageBucketSetObjectLockCommand builds `basaltic storage bucket set-object-lock`.
func newStorageBucketSetObjectLockCommand(state *cli.State) *cobra.Command {
	var body storage.PutBucketObjectLockRequest
	var bodyFile string
	var objectLockFlag string
	cmd := &cobra.Command{
		Use:   "set-object-lock <bucket>",
		Short: "Put bucket object-lock configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if objectLockFlag != "" {
				if err := json.Unmarshal([]byte(objectLockFlag), &body.ObjectLock); err != nil {
					return fmt.Errorf("--object-lock: %w", err)
				}
			}
			if err := c.PutBucketObjectLock(cmd.Context(), args[0], &body); err != nil {
				return err
			}
			state.Printer().Done("Set object lock requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&objectLockFlag, "object-lock", "", "Object lock (JSON)")
	_ = cmd.MarkFlagRequired("object-lock")
	return cmd
}

// newStorageBucketSetPolicyCommand builds `basaltic storage bucket set-policy`.
func newStorageBucketSetPolicyCommand(state *cli.State) *cobra.Command {
	var body storage.PutBucketPolicyRequest
	var bodyFile string
	var documentFlag string
	cmd := &cobra.Command{
		Use:   "set-policy <bucket>",
		Short: "Put bucket policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if documentFlag != "" {
				if err := json.Unmarshal([]byte(documentFlag), &body.Document); err != nil {
					return fmt.Errorf("--document: %w", err)
				}
			}
			if err := c.PutBucketPolicy(cmd.Context(), args[0], &body); err != nil {
				return err
			}
			state.Printer().Done("Set policy requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&documentFlag, "document", "", "Document (JSON)")
	_ = cmd.MarkFlagRequired("document")
	return cmd
}

// newStorageBucketSetTaggingCommand builds `basaltic storage bucket set-tagging`.
func newStorageBucketSetTaggingCommand(state *cli.State) *cobra.Command {
	var body storage.PutBucketTaggingRequest
	var bodyFile string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "set-tagging <bucket>",
		Short: "Put bucket tag set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
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
			if err := c.PutBucketTagging(cmd.Context(), args[0], &body); err != nil {
				return err
			}
			state.Printer().Done("Set tagging requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	_ = cmd.MarkFlagRequired("tags")
	return cmd
}

// newStorageBucketSetVersioningCommand builds `basaltic storage bucket set-versioning`.
func newStorageBucketSetVersioningCommand(state *cli.State) *cobra.Command {
	var body storage.PutBucketVersioningRequest
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "set-versioning <bucket>",
		Short: "Set bucket versioning state",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if err := c.PutBucketVersioning(cmd.Context(), args[0], &body); err != nil {
				return err
			}
			state.Printer().Done("Set versioning requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.Status, "status", "", "Status (one of: enabled, suspended)")
	_ = cmd.MarkFlagRequired("status")
	return cmd
}

// newStorageMultipartUploadCommand builds `basaltic storage multipart-upload`.
func newStorageMultipartUploadCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "multipart-upload",
		Short:   "Multipart uploads",
		Aliases: []string{"multipart-uploads"},
	}
	cmd.AddCommand(newStorageMultipartUploadListCommand(state))
	cmd.AddCommand(newStorageMultipartUploadAbortCommand(state))
	cmd.AddCommand(newStorageMultipartUploadCompleteCommand(state))
	cmd.AddCommand(newStorageMultipartUploadInitiateCommand(state))
	cmd.AddCommand(newStorageMultipartUploadListPartsCommand(state))
	cmd.AddCommand(newStorageMultipartUploadUploadPartCommand(state))
	return cmd
}

// newStorageMultipartUploadListCommand builds `basaltic storage multipart-upload list`.
func newStorageMultipartUploadListCommand(state *cli.State) *cobra.Command {
	var params storage.ListMultipartUploadsParams
	cmd := &cobra.Command{
		Use:   "list <bucket>",
		Short: "List in-flight multipart uploads",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListMultipartUploads(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.MaxUploads, "max-uploads", 0, "Max uploads")
	f.StringVar(&params.Prefix, "prefix", "", "Prefix")
	return cmd
}

// newStorageMultipartUploadAbortCommand builds `basaltic storage multipart-upload abort`.
func newStorageMultipartUploadAbortCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "abort <bucket> <upload-id>",
		Short: "Abort a multipart upload",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if err := c.AbortMultipartUpload(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			state.Printer().Done("Abort requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newStorageMultipartUploadCompleteCommand builds `basaltic storage multipart-upload complete`.
func newStorageMultipartUploadCompleteCommand(state *cli.State) *cobra.Command {
	var body storage.CompleteMultipartUploadRequest
	var bodyFile string
	var partsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "complete <bucket> <upload-id>",
		Short: "Complete a multipart upload",
		Args:  cobra.ExactArgs(2),
		Long:  "Complete a multipart upload.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if partsFlag != "" {
				if err := json.Unmarshal([]byte(partsFlag), &body.Parts); err != nil {
					return fmt.Errorf("--parts: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CompleteMultipartUpload(cmd.Context(), args[0], args[1], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&partsFlag, "parts", "", "Every part the assembled object is made of, in ascending part_number order (JSON)")
	_ = cmd.MarkFlagRequired("parts")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newStorageMultipartUploadInitiateCommand builds `basaltic storage multipart-upload initiate`.
func newStorageMultipartUploadInitiateCommand(state *cli.State) *cobra.Command {
	var body storage.InitiateMultipartUploadRequest
	var bodyFile string
	var contentTypeFlag string
	var metadataFlag string
	var storageClassFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "initiate <bucket>",
		Short: "Initiate a multipart upload",
		Args:  cobra.ExactArgs(1),
		Long:  "Initiate a multipart upload.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("content-type") {
				body.ContentType = &contentTypeFlag
			}
			if metadataFlag != "" {
				if err := json.Unmarshal([]byte(metadataFlag), &body.Metadata); err != nil {
					return fmt.Errorf("--metadata: %w", err)
				}
			}
			if cmd.Flags().Changed("storage-class") {
				body.StorageClass = &storageClassFlag
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.InitiateMultipartUpload(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&contentTypeFlag, "content-type", "", "Content type")
	f.StringVar(&body.Key, "key", "", "Key")
	_ = cmd.MarkFlagRequired("key")
	f.StringVar(&metadataFlag, "metadata", "", "Metadata (JSON)")
	f.StringVar(&storageClassFlag, "storage-class", "", "Storage class")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newStorageMultipartUploadListPartsCommand builds `basaltic storage multipart-upload list-parts`.
func newStorageMultipartUploadListPartsCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-parts <bucket> <upload-id>",
		Short: "List uploaded parts",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			out, err := c.ListParts(cmd.Context(), args[0], args[1])
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

// newStorageMultipartUploadUploadPartCommand builds `basaltic storage multipart-upload upload-part`.
func newStorageMultipartUploadUploadPartCommand(state *cli.State) *cobra.Command {
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "upload-part <bucket> <upload-id> <part-number>",
		Short: "Upload a part",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			reader, closeBody, err := openBody(bodyFile)
			if err != nil {
				return err
			}
			defer closeBody()
			out, err := c.UploadPart(cmd.Context(), args[0], args[1], args[2], reader)
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

// newStorageObjectCommand builds `basaltic storage object`.
func newStorageObjectCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "object",
		Short:   "Objects",
		Aliases: []string{"objects"},
	}
	cmd.AddCommand(newStorageObjectListCommand(state))
	cmd.AddCommand(newStorageObjectGetCommand(state))
	cmd.AddCommand(newStorageObjectDeleteCommand(state))
	cmd.AddCommand(newStorageObjectHeadCommand(state))
	cmd.AddCommand(newStorageObjectPutCommand(state))
	return cmd
}

// newStorageObjectListCommand builds `basaltic storage object list`.
func newStorageObjectListCommand(state *cli.State) *cobra.Command {
	var params storage.ListObjectsParams
	cmd := &cobra.Command{
		Use:   "list <bucket>",
		Short: "List objects",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			out, err := c.ListObjects(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.Delimiter, "delimiter", "", "Delimiter")
	f.StringVar(&params.Marker, "marker", "", "Marker")
	f.IntVar(&params.MaxKeys, "max-keys", 0, "Max keys")
	f.StringVar(&params.Prefix, "prefix", "", "Prefix")
	return cmd
}

// newStorageObjectGetCommand builds `basaltic storage object get`.
func newStorageObjectGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <bucket> <key>",
		Short: "Download object",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			stream, err := c.GetObject(cmd.Context(), args[0], args[1])
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

// newStorageObjectDeleteCommand builds `basaltic storage object delete`.
func newStorageObjectDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <bucket> <key>",
		Short: "Delete object",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteObject(cmd.Context(), args[0], args[1]); err != nil {
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

// newStorageObjectHeadCommand builds `basaltic storage object head`.
func newStorageObjectHeadCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "head <bucket> <key>",
		Short: "Head object",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if err := c.HeadObject(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			state.Printer().Done("Head requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newStorageObjectPutCommand builds `basaltic storage object put`.
func newStorageObjectPutCommand(state *cli.State) *cobra.Command {
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "put <bucket> <key>",
		Short: "Upload object",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			reader, closeBody, err := openBody(bodyFile)
			if err != nil {
				return err
			}
			defer closeBody()
			out, err := c.PutObject(cmd.Context(), args[0], args[1], reader)
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

// newStorageSnapshotCommand builds `basaltic storage snapshot`.
func newStorageSnapshotCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "snapshot",
		Short:   "Snapshots",
		Aliases: []string{"snapshots"},
	}
	cmd.AddCommand(newStorageSnapshotListCommand(state))
	cmd.AddCommand(newStorageSnapshotGetCommand(state))
	cmd.AddCommand(newStorageSnapshotCreateCommand(state))
	cmd.AddCommand(newStorageSnapshotUpdateCommand(state))
	cmd.AddCommand(newStorageSnapshotDeleteCommand(state))
	return cmd
}

// newStorageSnapshotListCommand builds `basaltic storage snapshot list`.
func newStorageSnapshotListCommand(state *cli.State) *cobra.Command {
	var params storage.ListSnapshotsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List snapshots",
		Args:  cobra.ExactArgs(0),
		Long:  "List snapshots.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListSnapshotsAll(cmd.Context(), &params))
			}
			page, err := c.ListSnapshots(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Resume token — the last snapshot id from the previous page")
	f.StringVar(&params.Name, "name", "", "Case-insensitive substring match on the snapshot name")
	f.StringVar((*string)(&params.Status), "status", "", "Status (one of: creating, available, deleting, error)")
	f.StringVar(&params.VolumeID, "volume-id", "", "Narrow the listing to snapshots of one volume")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newStorageSnapshotGetCommand builds `basaltic storage snapshot get`.
func newStorageSnapshotGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <snapshot-id>",
		Short: "Get snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetSnapshot(cmd.Context(), args[0])
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

// newStorageSnapshotCreateCommand builds `basaltic storage snapshot create`.
func newStorageSnapshotCreateCommand(state *cli.State) *cobra.Command {
	var body storage.SnapshotCreateRequest
	var bodyFile string
	var descriptionFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create snapshot",
		Args:  cobra.ExactArgs(0),
		Long:  "Create snapshot.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
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
			out, err := c.CreateSnapshot(cmd.Context(), &body, reqOpts...)
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
	f.StringVar(&body.VolumeID, "volume-id", "", "Volume id")
	_ = cmd.MarkFlagRequired("volume-id")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newStorageSnapshotUpdateCommand builds `basaltic storage snapshot update`.
func newStorageSnapshotUpdateCommand(state *cli.State) *cobra.Command {
	var body storage.SnapshotUpdateRequest
	var bodyFile string
	var descriptionFlag string
	var nameFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <snapshot-id>",
		Short: "Update snapshot metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
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
			if cmd.Flags().Changed("name") {
				body.Name = &nameFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateSnapshot(cmd.Context(), args[0], &body)
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
	f.StringVar(&nameFlag, "name", "", "Name")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newStorageSnapshotDeleteCommand builds `basaltic storage snapshot delete`.
func newStorageSnapshotDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <snapshot-id>",
		Short: "Delete snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteSnapshot(cmd.Context(), args[0]); err != nil {
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

// newStorageSnapshotPolicyCommand builds `basaltic storage snapshot-policy`.
func newStorageSnapshotPolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "snapshot-policy",
		Short:   "Snapshot policies",
		Aliases: []string{"snapshot-policies"},
	}
	cmd.AddCommand(newStorageSnapshotPolicyListCommand(state))
	cmd.AddCommand(newStorageSnapshotPolicyGetCommand(state))
	cmd.AddCommand(newStorageSnapshotPolicyCreateCommand(state))
	cmd.AddCommand(newStorageSnapshotPolicyUpdateCommand(state))
	cmd.AddCommand(newStorageSnapshotPolicyDeleteCommand(state))
	return cmd
}

// newStorageSnapshotPolicyListCommand builds `basaltic storage snapshot-policy list`.
func newStorageSnapshotPolicyListCommand(state *cli.State) *cobra.Command {
	var params storage.ListSnapshotPoliciesParams
	var enabledFlag bool
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List snapshot policies",
		Args:  cobra.ExactArgs(0),
		Long:  "List snapshot policies.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("enabled") {
				params.Enabled = &enabledFlag
			}
			if fetchAll {
				return state.Printer().Iter(c.ListSnapshotPoliciesAll(cmd.Context(), &params))
			}
			page, err := c.ListSnapshotPolicies(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.BoolVar(&enabledFlag, "enabled", false, "Narrow to enabled (or paused) policies")
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Resume token — the last policy id from the previous page")
	f.StringVar(&params.Name, "name", "", "Case-insensitive substring match on the policy name")
	f.StringVar(&params.VolumeID, "volume-id", "", "Narrow the listing to the policy attached to one volume")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newStorageSnapshotPolicyGetCommand builds `basaltic storage snapshot-policy get`.
func newStorageSnapshotPolicyGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <policy-id>",
		Short: "Get snapshot policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetSnapshotPolicy(cmd.Context(), args[0])
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

// newStorageSnapshotPolicyCreateCommand builds `basaltic storage snapshot-policy create`.
func newStorageSnapshotPolicyCreateCommand(state *cli.State) *cobra.Command {
	var body storage.SnapshotPolicyCreateRequest
	var bodyFile string
	var descriptionFlag string
	var enabledFlag bool
	var retentionDaysFlag int
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create snapshot policy",
		Args:  cobra.ExactArgs(0),
		Long:  "Create snapshot policy.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
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
			if cmd.Flags().Changed("enabled") {
				body.Enabled = &enabledFlag
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
			out, err := c.CreateSnapshotPolicy(cmd.Context(), &body, reqOpts...)
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
	f.BoolVar(&enabledFlag, "enabled", false, "Defaults to true")
	f.IntVar(&body.IntervalMinutes, "interval-minutes", 0, "Interval minutes")
	_ = cmd.MarkFlagRequired("interval-minutes")
	f.StringVar(&body.Name, "name", "", "Unique within the account — it names the policy in its CRN")
	_ = cmd.MarkFlagRequired("name")
	f.IntVar(&body.RetentionCount, "retention-count", 0, "Retention count")
	_ = cmd.MarkFlagRequired("retention-count")
	f.IntVar(&retentionDaysFlag, "retention-days", 0, "Retention days")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&body.VolumeID, "volume-id", "", "Volume id")
	_ = cmd.MarkFlagRequired("volume-id")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newStorageSnapshotPolicyUpdateCommand builds `basaltic storage snapshot-policy update`.
func newStorageSnapshotPolicyUpdateCommand(state *cli.State) *cobra.Command {
	var body storage.SnapshotPolicyUpdateRequest
	var bodyFile string
	var descriptionFlag string
	var enabledFlag bool
	var intervalMinutesFlag int
	var nameFlag string
	var retentionCountFlag int
	var retentionDaysFlag int
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <policy-id>",
		Short: "Update snapshot policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
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
			if cmd.Flags().Changed("enabled") {
				body.Enabled = &enabledFlag
			}
			if cmd.Flags().Changed("interval-minutes") {
				body.IntervalMinutes = &intervalMinutesFlag
			}
			if cmd.Flags().Changed("name") {
				body.Name = &nameFlag
			}
			if cmd.Flags().Changed("retention-count") {
				body.RetentionCount = &retentionCountFlag
			}
			if cmd.Flags().Changed("retention-days") {
				body.RetentionDays = &retentionDaysFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateSnapshotPolicy(cmd.Context(), args[0], &body)
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
	f.BoolVar(&enabledFlag, "enabled", false, "false pauses the policy, true resumes it")
	f.IntVar(&intervalMinutesFlag, "interval-minutes", 0, "Interval minutes")
	f.StringVar(&nameFlag, "name", "", "Name")
	f.IntVar(&retentionCountFlag, "retention-count", 0, "Retention count")
	f.IntVar(&retentionDaysFlag, "retention-days", 0, "Retention days")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newStorageSnapshotPolicyDeleteCommand builds `basaltic storage snapshot-policy delete`.
func newStorageSnapshotPolicyDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <policy-id>",
		Short: "Delete snapshot policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteSnapshotPolicy(cmd.Context(), args[0]); err != nil {
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

// newStorageVolumeCommand builds `basaltic storage volume`.
func newStorageVolumeCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "volume",
		Short:   "Volumes",
		Aliases: []string{"volumes"},
	}
	cmd.AddCommand(newStorageVolumeListCommand(state))
	cmd.AddCommand(newStorageVolumeGetCommand(state))
	cmd.AddCommand(newStorageVolumeCreateCommand(state))
	cmd.AddCommand(newStorageVolumeUpdateCommand(state))
	cmd.AddCommand(newStorageVolumeDeleteCommand(state))
	cmd.AddCommand(newStorageVolumeExtendCommand(state))
	return cmd
}

// newStorageVolumeListCommand builds `basaltic storage volume list`.
func newStorageVolumeListCommand(state *cli.State) *cobra.Command {
	var params storage.ListVolumesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List volumes",
		Args:  cobra.ExactArgs(0),
		Long:  "List volumes.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListVolumesAll(cmd.Context(), &params))
			}
			page, err := c.ListVolumes(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Resume token — the last volume id from the previous page")
	f.StringVar(&params.Name, "name", "", "Case-insensitive substring match on the volume name")
	f.StringVar((*string)(&params.Status), "status", "", "Status (one of: creating, available, in_use, extending, deleting, error)")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newStorageVolumeGetCommand builds `basaltic storage volume get`.
func newStorageVolumeGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <volume-id>",
		Short: "Get volume",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetVolume(cmd.Context(), args[0])
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

// newStorageVolumeCreateCommand builds `basaltic storage volume create`.
func newStorageVolumeCreateCommand(state *cli.State) *cobra.Command {
	var body storage.VolumeCreateRequest
	var bodyFile string
	var bootableFlag bool
	var descriptionFlag string
	var sourceImageIdFlag string
	var sourceSnapshotIdFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create volume",
		Args:  cobra.ExactArgs(0),
		Long:  "Create volume.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("bootable") {
				body.Bootable = &bootableFlag
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
			}
			if cmd.Flags().Changed("source-image-id") {
				body.SourceImageID = &sourceImageIdFlag
			}
			if cmd.Flags().Changed("source-snapshot-id") {
				body.SourceSnapshotID = &sourceSnapshotIdFlag
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
			out, err := c.CreateVolume(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.BoolVar(&bootableFlag, "bootable", false, "Bootable")
	f.StringVar(&descriptionFlag, "description", "", "Description")
	f.StringVar(&body.Name, "name", "", "Name")
	_ = cmd.MarkFlagRequired("name")
	f.IntVar(&body.SizeGB, "size-gb", 0, "Size gb")
	_ = cmd.MarkFlagRequired("size-gb")
	f.StringVar(&sourceImageIdFlag, "source-image-id", "", "Reserved: when supplied, the volume will be provisioned as a clone of the referenced image's base snapshot")
	f.StringVar(&sourceSnapshotIdFlag, "source-snapshot-id", "", "Clone the new volume from an existing snapshot (restore)")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar((*string)(&body.VolumeType), "volume-type", "", "Volume type (one of: ssd, nvme)")
	_ = cmd.MarkFlagRequired("volume-type")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newStorageVolumeUpdateCommand builds `basaltic storage volume update`.
func newStorageVolumeUpdateCommand(state *cli.State) *cobra.Command {
	var body storage.VolumeUpdateRequest
	var bodyFile string
	var descriptionFlag string
	var nameFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <volume-id>",
		Short: "Update volume metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
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
			if cmd.Flags().Changed("name") {
				body.Name = &nameFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateVolume(cmd.Context(), args[0], &body)
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
	f.StringVar(&nameFlag, "name", "", "Name")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newStorageVolumeDeleteCommand builds `basaltic storage volume delete`.
func newStorageVolumeDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <volume-id>",
		Short: "Delete volume",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteVolume(cmd.Context(), args[0]); err != nil {
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

// newStorageVolumeExtendCommand builds `basaltic storage volume extend`.
func newStorageVolumeExtendCommand(state *cli.State) *cobra.Command {
	var body storage.VolumeExtendRequest
	var bodyFile string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "extend <volume-id>",
		Short: "Extend volume",
		Args:  cobra.ExactArgs(1),
		Long:  "Extend volume.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
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
			out, err := c.ExtendVolume(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.IntVar(&body.NewSizeGB, "new-size-gb", 0, "New size in GB")
	_ = cmd.MarkFlagRequired("new-size-gb")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newStorageVolumeTypeCommand builds `basaltic storage volume-type`.
func newStorageVolumeTypeCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "volume-type",
		Short:   "Volume types",
		Aliases: []string{"volume-types"},
	}
	cmd.AddCommand(newStorageVolumeTypeListCommand(state))
	return cmd
}

// newStorageVolumeTypeListCommand builds `basaltic storage volume-type list`.
func newStorageVolumeTypeListCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List volume types",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := storageClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListVolumeTypes(cmd.Context())
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
