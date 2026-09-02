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
	"github.com/basaltic-sh/sdk-go/database"

	"github.com/basaltic-sh/cli/internal/cli"
)

func init() { cli.RegisterService(newDatabaseCommand) }

// newDatabaseCommand builds `basaltic database`.
func newDatabaseCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "database",
		Short:   "Managed database clusters",
		Aliases: []string{"db"},
		Long:    "Managed database clusters.\n\nThis is a regional service: it acts in the region from --region, the\nBASALTIC_REGION environment variable, or the profile.",
	}
	cmd.AddCommand(newDatabaseClusterCommand(state))
	cmd.AddCommand(newDatabaseDatabaseCommand(state))
	cmd.AddCommand(newDatabaseEngineCommand(state))
	cmd.AddCommand(newDatabaseParameterGroupCommand(state))
	cmd.AddCommand(newDatabaseUserCommand(state))
	return cmd
}

// databaseClient builds the service client, resolving credentials on first use.
func databaseClient(state *cli.State) (*database.Client, error) {
	cfg, err := state.SDK()
	if err != nil {
		return nil, err
	}
	return database.New(cfg), nil
}

// newDatabaseClusterCommand builds `basaltic database cluster`.
func newDatabaseClusterCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cluster",
		Short:   "Clusters",
		Aliases: []string{"clusters"},
	}
	cmd.AddCommand(newDatabaseClusterListCommand(state))
	cmd.AddCommand(newDatabaseClusterGetCommand(state))
	cmd.AddCommand(newDatabaseClusterCreateCommand(state))
	cmd.AddCommand(newDatabaseClusterUpdateCommand(state))
	cmd.AddCommand(newDatabaseClusterDeleteCommand(state))
	cmd.AddCommand(newDatabaseClusterAddReplicaCommand(state))
	cmd.AddCommand(newDatabaseClusterConvertToHaCommand(state))
	cmd.AddCommand(newDatabaseClusterFailoverCommand(state))
	cmd.AddCommand(newDatabaseClusterGetBackupCommand(state))
	cmd.AddCommand(newDatabaseClusterListBackupsCommand(state))
	cmd.AddCommand(newDatabaseClusterRemoveReplicaCommand(state))
	cmd.AddCommand(newDatabaseClusterRequestBackupCommand(state))
	cmd.AddCommand(newDatabaseClusterRestoreCommand(state))
	return cmd
}

// newDatabaseClusterListCommand builds `basaltic database cluster list`.
func newDatabaseClusterListCommand(state *cli.State) *cobra.Command {
	var params database.ListClustersParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List database clusters",
		Args:  cobra.ExactArgs(0),
		Long:  "List database clusters.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListClustersAll(cmd.Context(), &params))
			}
			page, err := c.ListClusters(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.EngineType, "engine-type", "", "One of: \"postgres\", \"valkey\"")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Name")
	f.StringVar(&params.Status, "status", "", "Status")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newDatabaseClusterGetCommand builds `basaltic database cluster get`.
func newDatabaseClusterGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <cluster-id>",
		Short: "Get a database cluster",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetCluster(cmd.Context(), args[0])
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

// newDatabaseClusterCreateCommand builds `basaltic database cluster create`.
func newDatabaseClusterCreateCommand(state *cli.State) *cobra.Command {
	var body database.CreateClusterRequest
	var bodyFile string
	var adminUserFlag string
	var assignPublicIpFlag bool
	var defaultDatabaseFlag string
	var descriptionFlag string
	var engineVersionFlag string
	var instanceCountFlag int
	var metadataFlag string
	var networksFlag string
	var parameterGroupIdFlag string
	var restoreFromFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a database cluster",
		Args:  cobra.ExactArgs(0),
		Long:  "Create a database cluster.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("admin-user") {
				body.AdminUser = &adminUserFlag
			}
			if cmd.Flags().Changed("assign-public-ip") {
				body.AssignPublicIP = &assignPublicIpFlag
			}
			if cmd.Flags().Changed("default-database") {
				body.DefaultDatabase = &defaultDatabaseFlag
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
			}
			if cmd.Flags().Changed("engine-version") {
				body.EngineVersion = &engineVersionFlag
			}
			if cmd.Flags().Changed("instance-count") {
				body.InstanceCount = &instanceCountFlag
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
			if cmd.Flags().Changed("parameter-group-id") {
				body.ParameterGroupID = &parameterGroupIdFlag
			}
			if restoreFromFlag != "" {
				if err := json.Unmarshal([]byte(restoreFromFlag), &body.RestoreFrom); err != nil {
					return fmt.Errorf("--restore-from: %w", err)
				}
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
			out, err := c.CreateCluster(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&adminUserFlag, "admin-user", "", "Optional; defaults to \"admin\"")
	f.BoolVar(&assignPublicIpFlag, "assign-public-ip", false, "Endpoint exposure")
	f.StringVar(&defaultDatabaseFlag, "default-database", "", "Optional; defaults to \"default\"")
	f.StringVar(&descriptionFlag, "description", "", "Description")
	f.StringVar(&body.EngineType, "engine-type", "", "Engine type (one of: postgres, valkey)")
	_ = cmd.MarkFlagRequired("engine-type")
	f.StringVar(&engineVersionFlag, "engine-version", "", "e.g")
	f.StringVar(&body.FlavorID, "flavor-id", "", "Compute flavor each cluster member runs on")
	_ = cmd.MarkFlagRequired("flavor-id")
	f.IntVar(&instanceCountFlag, "instance-count", 0, "Node count")
	f.StringSliceVar(&body.KeyNames, "key-names", nil, "Platform-operator break-glass only")
	f.StringVar(&metadataFlag, "metadata", "", "Metadata (JSON)")
	f.StringVar(&body.Name, "name", "", "Name")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&networksFlag, "networks", "", "At least one network interface is required (JSON)")
	_ = cmd.MarkFlagRequired("networks")
	f.StringVar(&parameterGroupIdFlag, "parameter-group-id", "", "Binds the cluster to a parameter group whose settings its members boot with and converge to")
	f.StringVar(&restoreFromFlag, "restore-from", "", "Restore from (JSON)")
	f.IntVar(&body.StorageGB, "storage-gb", 0, "Storage gb")
	_ = cmd.MarkFlagRequired("storage-gb")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newDatabaseClusterUpdateCommand builds `basaltic database cluster update`.
func newDatabaseClusterUpdateCommand(state *cli.State) *cobra.Command {
	var body database.UpdateClusterRequest
	var bodyFile string
	var descriptionFlag string
	var parameterGroupIdFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <cluster-id>",
		Short: "Update a database cluster (description + tags only)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
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
			if cmd.Flags().Changed("parameter-group-id") {
				body.ParameterGroupID = &parameterGroupIdFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateCluster(cmd.Context(), args[0], &body)
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
	f.StringVar(&parameterGroupIdFlag, "parameter-group-id", "", "Re-binds the cluster to a different parameter group; its members converge on the new settings")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newDatabaseClusterDeleteCommand builds `basaltic database cluster delete`.
func newDatabaseClusterDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <cluster-id>",
		Short: "Delete a database cluster",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteCluster(cmd.Context(), args[0]); err != nil {
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

// newDatabaseClusterAddReplicaCommand builds `basaltic database cluster add-replica`.
func newDatabaseClusterAddReplicaCommand(state *cli.State) *cobra.Command {
	var body database.AddReplicaRequest
	var bodyFile string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "add-replica <cluster-id>",
		Short: "Add a read replica to the cluster",
		Args:  cobra.ExactArgs(1),
		Long:  "Add a read replica to the cluster.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
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
			out, err := c.AddReplica(cmd.Context(), args[0], body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newDatabaseClusterConvertToHaCommand builds `basaltic database cluster convert-to-ha`.
func newDatabaseClusterConvertToHaCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert-to-ha <cluster-id>",
		Short: "Convert a single-node postgres cluster to Patroni-managed HA",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			out, err := c.ConvertClusterToHA(cmd.Context(), args[0])
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

// newDatabaseClusterFailoverCommand builds `basaltic database cluster failover`.
func newDatabaseClusterFailoverCommand(state *cli.State) *cobra.Command {
	var body database.FailoverRequest
	var bodyFile string
	var targetMemberFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "failover <cluster-id>",
		Short: "Trigger a planned HA switchover (Patroni or Valkey Sentinel)",
		Args:  cobra.ExactArgs(1),
		Long:  "Trigger a planned HA switchover (Patroni or Valkey Sentinel).\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("target-member") {
				body.TargetMember = &targetMemberFlag
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.FailoverCluster(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&targetMemberFlag, "target-member", "", "Optional member name to promote")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newDatabaseClusterGetBackupCommand builds `basaltic database cluster get-backup`.
func newDatabaseClusterGetBackupCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-backup <cluster-id> <backup-id>",
		Short: "Get a backup",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetBackup(cmd.Context(), args[0], args[1])
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

// newDatabaseClusterListBackupsCommand builds `basaltic database cluster list-backups`.
func newDatabaseClusterListBackupsCommand(state *cli.State) *cobra.Command {
	var params database.ListBackupsParams
	cmd := &cobra.Command{
		Use:   "list-backups <cluster-id>",
		Short: "List a cluster's backups",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListBackups(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.Kind, "kind", "", "Filter by backup kind")
	f.StringVar(&params.Status, "status", "", "Status")
	return cmd
}

// newDatabaseClusterRemoveReplicaCommand builds `basaltic database cluster remove-replica`.
func newDatabaseClusterRemoveReplicaCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-replica <cluster-id> <instance-id>",
		Short: "Remove a read replica from the cluster",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			if err := c.RemoveReplica(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			state.Printer().Done("Remove replica requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newDatabaseClusterRequestBackupCommand builds `basaltic database cluster request-backup`.
func newDatabaseClusterRequestBackupCommand(state *cli.State) *cobra.Command {
	var body database.RequestBackupRequest
	var bodyFile string
	var kindFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "request-backup <cluster-id>",
		Short: "Request a manual backup",
		Args:  cobra.ExactArgs(1),
		Long:  "Request a manual backup.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("kind") {
				body.Kind = &kindFlag
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.RequestBackup(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&kindFlag, "kind", "", "Backup kind (one of: base, incremental, full, incr)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newDatabaseClusterRestoreCommand builds `basaltic database cluster restore`.
func newDatabaseClusterRestoreCommand(state *cli.State) *cobra.Command {
	var body database.RestoreClusterRequest
	var bodyFile string
	var recoveryTargetTimeFlag string
	cmd := &cobra.Command{
		Use:   "restore <cluster-id>",
		Short: "Restore a backup into this cluster, in place",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if recoveryTargetTimeFlag != "" {
				parsed, err := parseTime(recoveryTargetTimeFlag)
				if err != nil {
					return fmt.Errorf("--recovery-target-time: %w", err)
				}
				body.RecoveryTargetTime = &parsed
			}
			out, err := c.RestoreCluster(cmd.Context(), args[0], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.BackupID, "backup-id", "", "The backup to restore")
	_ = cmd.MarkFlagRequired("backup-id")
	f.StringVar(&body.Confirm, "confirm", "", "Must equal the cluster's name")
	_ = cmd.MarkFlagRequired("confirm")
	f.StringVar(&recoveryTargetTimeFlag, "recovery-target-time", "", "RFC3339 timestamp for point-in-time recovery — WAL is replayed forward from the backup until this time, then stops (RFC 3339)")
	return cmd
}

// newDatabaseDatabaseCommand builds `basaltic database database`.
func newDatabaseDatabaseCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "database",
		Short:   "Databases",
		Aliases: []string{"databases"},
	}
	cmd.AddCommand(newDatabaseDatabaseListCommand(state))
	cmd.AddCommand(newDatabaseDatabaseGetCommand(state))
	cmd.AddCommand(newDatabaseDatabaseCreateCommand(state))
	cmd.AddCommand(newDatabaseDatabaseDeleteCommand(state))
	return cmd
}

// newDatabaseDatabaseListCommand builds `basaltic database database list`.
func newDatabaseDatabaseListCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <cluster-id>",
		Short: "List a cluster's logical databases",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListDatabases(cmd.Context(), args[0])
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

// newDatabaseDatabaseGetCommand builds `basaltic database database get`.
func newDatabaseDatabaseGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <cluster-id> <database-id>",
		Short: "Get a logical database",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetDatabase(cmd.Context(), args[0], args[1])
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

// newDatabaseDatabaseCreateCommand builds `basaltic database database create`.
func newDatabaseDatabaseCreateCommand(state *cli.State) *cobra.Command {
	var body database.CreateDatabaseRequest
	var bodyFile string
	var characterSetFlag string
	var collationFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create <cluster-id>",
		Short: "Create a logical database",
		Args:  cobra.ExactArgs(1),
		Long:  "Create a logical database.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("character-set") {
				body.CharacterSet = &characterSetFlag
			}
			if cmd.Flags().Changed("collation") {
				body.Collation = &collationFlag
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateDatabase(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&characterSetFlag, "character-set", "", "Character set")
	f.StringVar(&collationFlag, "collation", "", "Collation")
	f.StringVar(&body.Name, "name", "", "Name")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newDatabaseDatabaseDeleteCommand builds `basaltic database database delete`.
func newDatabaseDatabaseDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <cluster-id> <database-id>",
		Short: "Delete a logical database",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteDatabase(cmd.Context(), args[0], args[1]); err != nil {
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

// newDatabaseEngineCommand builds `basaltic database engine`.
func newDatabaseEngineCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "engine",
		Short:   "Engines",
		Aliases: []string{"engines"},
	}
	cmd.AddCommand(newDatabaseEngineListCommand(state))
	cmd.AddCommand(newDatabaseEngineListParametersCommand(state))
	return cmd
}

// newDatabaseEngineListCommand builds `basaltic database engine list`.
func newDatabaseEngineListCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List supported database engines",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListEngines(cmd.Context())
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

// newDatabaseEngineListParametersCommand builds `basaltic database engine list-parameters`.
func newDatabaseEngineListParametersCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-parameters <engine-type>",
		Short: "List an engine's tunable parameters",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			out, err := c.ListEngineParameters(cmd.Context(), args[0])
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

// newDatabaseParameterGroupCommand builds `basaltic database parameter-group`.
func newDatabaseParameterGroupCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "parameter-group",
		Short:   "Parameter groups",
		Aliases: []string{"parameter-groups"},
	}
	cmd.AddCommand(newDatabaseParameterGroupListCommand(state))
	cmd.AddCommand(newDatabaseParameterGroupGetCommand(state))
	cmd.AddCommand(newDatabaseParameterGroupCreateCommand(state))
	cmd.AddCommand(newDatabaseParameterGroupUpdateCommand(state))
	cmd.AddCommand(newDatabaseParameterGroupDeleteCommand(state))
	return cmd
}

// newDatabaseParameterGroupListCommand builds `basaltic database parameter-group list`.
func newDatabaseParameterGroupListCommand(state *cli.State) *cobra.Command {
	var params database.ListParameterGroupsParams
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List parameter groups",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListParameterGroups(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.EngineType, "engine-type", "", "One of: \"postgres\"")
	return cmd
}

// newDatabaseParameterGroupGetCommand builds `basaltic database parameter-group get`.
func newDatabaseParameterGroupGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <parameter-group-id>",
		Short: "Get a parameter group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetParameterGroup(cmd.Context(), args[0])
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

// newDatabaseParameterGroupCreateCommand builds `basaltic database parameter-group create`.
func newDatabaseParameterGroupCreateCommand(state *cli.State) *cobra.Command {
	var body database.CreateParameterGroupRequest
	var bodyFile string
	var descriptionFlag string
	var engineVersionFlag string
	var paramsFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a parameter group",
		Args:  cobra.ExactArgs(0),
		Long:  "Create a parameter group.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
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
			if cmd.Flags().Changed("engine-version") {
				body.EngineVersion = &engineVersionFlag
			}
			if paramsFlag != "" {
				if err := json.Unmarshal([]byte(paramsFlag), &body.Params); err != nil {
					return fmt.Errorf("--params: %w", err)
				}
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
			out, err := c.CreateParameterGroup(cmd.Context(), &body, reqOpts...)
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
	f.StringVar(&body.EngineType, "engine-type", "", "Engine type (one of: postgres, valkey)")
	_ = cmd.MarkFlagRequired("engine-type")
	f.StringVar(&engineVersionFlag, "engine-version", "", "Defaults to the engine current default when omitted")
	f.StringVar(&body.Name, "name", "", "Name")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&paramsFlag, "params", "", "Engine settings this group applies (JSON)")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newDatabaseParameterGroupUpdateCommand builds `basaltic database parameter-group update`.
func newDatabaseParameterGroupUpdateCommand(state *cli.State) *cobra.Command {
	var body database.UpdateParameterGroupRequest
	var bodyFile string
	var descriptionFlag string
	var paramsFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <parameter-group-id>",
		Short: "Update a parameter group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
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
			if paramsFlag != "" {
				if err := json.Unmarshal([]byte(paramsFlag), &body.Params); err != nil {
					return fmt.Errorf("--params: %w", err)
				}
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateParameterGroup(cmd.Context(), args[0], &body)
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
	f.StringVar(&paramsFlag, "params", "", "Params (JSON)")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newDatabaseParameterGroupDeleteCommand builds `basaltic database parameter-group delete`.
func newDatabaseParameterGroupDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <parameter-group-id>",
		Short: "Delete a parameter group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteParameterGroup(cmd.Context(), args[0]); err != nil {
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

// newDatabaseUserCommand builds `basaltic database user`.
func newDatabaseUserCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "user",
		Short:   "Users",
		Aliases: []string{"users"},
	}
	cmd.AddCommand(newDatabaseUserListCommand(state))
	cmd.AddCommand(newDatabaseUserGetCommand(state))
	cmd.AddCommand(newDatabaseUserCreateCommand(state))
	cmd.AddCommand(newDatabaseUserDeleteCommand(state))
	cmd.AddCommand(newDatabaseUserRotatePasswordCommand(state))
	return cmd
}

// newDatabaseUserListCommand builds `basaltic database user list`.
func newDatabaseUserListCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <cluster-id>",
		Short: "List a cluster's database users",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListDBUsers(cmd.Context(), args[0])
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

// newDatabaseUserGetCommand builds `basaltic database user get`.
func newDatabaseUserGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <cluster-id> <user-id>",
		Short: "Get a database user",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetDBUser(cmd.Context(), args[0], args[1])
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

// newDatabaseUserCreateCommand builds `basaltic database user create`.
func newDatabaseUserCreateCommand(state *cli.State) *cobra.Command {
	var body database.CreateDBUserRequest
	var bodyFile string
	var permissionsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create <cluster-id>",
		Short: "Create a database user",
		Args:  cobra.ExactArgs(1),
		Long:  "Create a database user.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if permissionsFlag != "" {
				if err := json.Unmarshal([]byte(permissionsFlag), &body.Permissions); err != nil {
					return fmt.Errorf("--permissions: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateDBUser(cmd.Context(), args[0], &body, reqOpts...)
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
	f.StringVar(&permissionsFlag, "permissions", "", "Engine-specific grant/permission map (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newDatabaseUserDeleteCommand builds `basaltic database user delete`.
func newDatabaseUserDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <cluster-id> <user-id>",
		Short: "Delete a database user",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteDBUser(cmd.Context(), args[0], args[1]); err != nil {
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

// newDatabaseUserRotatePasswordCommand builds `basaltic database user rotate-password`.
func newDatabaseUserRotatePasswordCommand(state *cli.State) *cobra.Command {
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "rotate-password <cluster-id> <user-id>",
		Short: "Rotate a database user's password",
		Args:  cobra.ExactArgs(2),
		Long:  "Rotate a database user's password.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := databaseClient(state)
			if err != nil {
				return err
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.RotateDBUserPassword(cmd.Context(), args[0], args[1], reqOpts...)
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
