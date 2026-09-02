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
	"github.com/basaltic-sh/sdk-go/secrets"

	"github.com/basaltic-sh/cli/internal/cli"
)

func init() { cli.RegisterService(newSecretsCommand) }

// newSecretsCommand builds `basaltic secrets`.
func newSecretsCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "Secrets and their versions",
		Long:  "Secrets and their versions.\n\nThis is a regional service: it acts in the region from --region, the\nBASALTIC_REGION environment variable, or the profile.",
	}
	cmd.AddCommand(newSecretsSecretListCommand(state))
	cmd.AddCommand(newSecretsSecretGetCommand(state))
	cmd.AddCommand(newSecretsSecretCreateCommand(state))
	cmd.AddCommand(newSecretsSecretUpdateCommand(state))
	cmd.AddCommand(newSecretsSecretDeleteCommand(state))
	cmd.AddCommand(newSecretsSecretGetValueCommand(state))
	cmd.AddCommand(newSecretsSecretListVersionsCommand(state))
	cmd.AddCommand(newSecretsSecretRestoreCommand(state))
	cmd.AddCommand(newSecretsSecretSetValueCommand(state))
	return cmd
}

// secretsClient builds the service client, resolving credentials on first use.
func secretsClient(state *cli.State) (*secrets.Client, error) {
	cfg, err := state.SDK()
	if err != nil {
		return nil, err
	}
	return secrets.New(cfg), nil
}

// newSecretsSecretListCommand builds `basaltic secrets secret list`.
func newSecretsSecretListCommand(state *cli.State) *cobra.Command {
	var params secrets.ListSecretsParams
	var includeDeletedFlag bool
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List secrets",
		Args:  cobra.ExactArgs(0),
		Long:  "List secrets.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := secretsClient(state)
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("include-deleted") {
				params.IncludeDeleted = &includeDeletedFlag
			}
			if fetchAll {
				return state.Printer().Iter(c.ListSecretsAll(cmd.Context(), &params))
			}
			page, err := c.ListSecrets(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.BoolVar(&includeDeletedFlag, "include-deleted", false, "Include secrets in the recovery window")
	f.IntVar(&params.Limit, "limit", 0, "Maximum items to return")
	f.StringVar(&params.Marker, "marker", "", "Marker")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newSecretsSecretGetCommand builds `basaltic secrets secret get`.
func newSecretsSecretGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <secret-id>",
		Short: "Describe a secret (no value)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := secretsClient(state)
			if err != nil {
				return err
			}
			out, err := c.DescribeSecret(cmd.Context(), args[0])
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

// newSecretsSecretCreateCommand builds `basaltic secrets secret create`.
func newSecretsSecretCreateCommand(state *cli.State) *cobra.Command {
	var body secrets.CreateSecretRequest
	var bodyFile string
	var descriptionFlag string
	var kmsKeyIdFlag string
	var recoveryWindowSecondsFlag int
	var tagsFlag string
	var valueFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new secret with an initial value",
		Args:  cobra.ExactArgs(0),
		Long:  "Create a new secret with an initial value.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := secretsClient(state)
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
			if cmd.Flags().Changed("kms-key-id") {
				body.KMSKeyID = &kmsKeyIdFlag
			}
			if cmd.Flags().Changed("recovery-window-seconds") {
				body.RecoveryWindowSeconds = &recoveryWindowSecondsFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			if valueFlag != "" {
				body.Value = []byte(valueFlag)
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateSecret(cmd.Context(), &body, reqOpts...)
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
	f.StringVar(&kmsKeyIdFlag, "kms-key-id", "", "Customer-managed KMS key to encrypt this secret under")
	f.StringVar(&body.Name, "name", "", "Unique within the calling account")
	_ = cmd.MarkFlagRequired("name")
	f.IntVar(&recoveryWindowSecondsFlag, "recovery-window-seconds", 0, "Recovery window seconds")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&valueFlag, "value", "", "Base64 of the initial value bytes (1 byte - 64 KiB)")
	_ = cmd.MarkFlagRequired("value")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newSecretsSecretUpdateCommand builds `basaltic secrets secret update`.
func newSecretsSecretUpdateCommand(state *cli.State) *cobra.Command {
	var body secrets.UpdateSecretRequest
	var bodyFile string
	var descriptionFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <secret-id>",
		Short: "Update mutable metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := secretsClient(state)
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
			out, err := c.UpdateSecret(cmd.Context(), args[0], &body)
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

// newSecretsSecretDeleteCommand builds `basaltic secrets secret delete`.
func newSecretsSecretDeleteCommand(state *cli.State) *cobra.Command {
	var body secrets.DeleteSecretRequest
	var bodyFile string
	var recoveryWindowSecondsFlag int
	cmd := &cobra.Command{
		Use:   "delete <secret-id>",
		Short: "Schedule deletion (soft delete with recovery window)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := secretsClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("recovery-window-seconds") {
				body.RecoveryWindowSeconds = &recoveryWindowSecondsFlag
			}
			out, err := c.DeleteSecret(cmd.Context(), args[0], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.IntVar(&recoveryWindowSecondsFlag, "recovery-window-seconds", 0, "Override the secret's default window")
	return cmd
}

// newSecretsSecretGetValueCommand builds `basaltic secrets secret get-value`.
func newSecretsSecretGetValueCommand(state *cli.State) *cobra.Command {
	var params secrets.GetSecretValueParams
	cmd := &cobra.Command{
		Use:   "get-value <secret-id>",
		Short: "Read the current value (or a specific version)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := secretsClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetSecretValue(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Version, "version", 0, "Specific version to read")
	return cmd
}

// newSecretsSecretListVersionsCommand builds `basaltic secrets secret list-versions`.
func newSecretsSecretListVersionsCommand(state *cli.State) *cobra.Command {
	var params secrets.ListVersionsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list-versions <secret-id>",
		Short: "List versions",
		Args:  cobra.ExactArgs(1),
		Long:  "List versions.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := secretsClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListVersionsAll(cmd.Context(), args[0], &params))
			}
			page, err := c.ListVersions(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Maximum items to return")
	f.StringVar(&params.Marker, "marker", "", "Marker")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newSecretsSecretRestoreCommand builds `basaltic secrets secret restore`.
func newSecretsSecretRestoreCommand(state *cli.State) *cobra.Command {
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "restore <secret-id>",
		Short: "Restore a secret from the recovery window",
		Args:  cobra.ExactArgs(1),
		Long:  "Restore a secret from the recovery window.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := secretsClient(state)
			if err != nil {
				return err
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.RestoreSecret(cmd.Context(), args[0], reqOpts...)
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

// newSecretsSecretSetValueCommand builds `basaltic secrets secret set-value`.
func newSecretsSecretSetValueCommand(state *cli.State) *cobra.Command {
	var body secrets.PutSecretValueRequest
	var bodyFile string
	var valueFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "set-value <secret-id>",
		Short: "Store a new version (becomes current)",
		Args:  cobra.ExactArgs(1),
		Long:  "Store a new version (becomes current).\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := secretsClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if valueFlag != "" {
				body.Value = []byte(valueFlag)
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.PutSecretValue(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&valueFlag, "value", "", "Base64 of the new value bytes (1 byte - 64 KiB)")
	_ = cmd.MarkFlagRequired("value")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}
