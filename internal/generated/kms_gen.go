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
	"github.com/basaltic-sh/sdk-go/kms"

	"github.com/basaltic-sh/cli/internal/cli"
)

func init() { cli.RegisterService(newKmsCommand) }

// newKmsCommand builds `basaltic kms`.
func newKmsCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kms",
		Short: "Encryption keys",
		Long:  "Encryption keys.\n\nThis is a regional service: it acts in the region from --region, the\nBASALTIC_REGION environment variable, or the profile.",
	}
	cmd.AddCommand(newKmsKeyCommand(state))
	return cmd
}

// kmsClient builds the service client, resolving credentials on first use.
func kmsClient(state *cli.State) (*kms.Client, error) {
	cfg, err := state.SDK()
	if err != nil {
		return nil, err
	}
	return kms.New(cfg), nil
}

// newKmsKeyCommand builds `basaltic kms key`.
func newKmsKeyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "key",
		Short:   "Keys",
		Aliases: []string{"keys"},
	}
	cmd.AddCommand(newKmsKeyListCommand(state))
	cmd.AddCommand(newKmsKeyGetCommand(state))
	cmd.AddCommand(newKmsKeyCreateCommand(state))
	cmd.AddCommand(newKmsKeyUpdateCommand(state))
	cmd.AddCommand(newKmsKeyCancelDeletionCommand(state))
	cmd.AddCommand(newKmsKeyDecryptCommand(state))
	cmd.AddCommand(newKmsKeyDisableCommand(state))
	cmd.AddCommand(newKmsKeyEnableCommand(state))
	cmd.AddCommand(newKmsKeyEncryptCommand(state))
	cmd.AddCommand(newKmsKeyGenerateDataKeyCommand(state))
	cmd.AddCommand(newKmsKeyScheduleDeletionCommand(state))
	cmd.AddCommand(newKmsKeySignCommand(state))
	cmd.AddCommand(newKmsKeyVerifyCommand(state))
	return cmd
}

// newKmsKeyListCommand builds `basaltic kms key list`.
func newKmsKeyListCommand(state *cli.State) *cobra.Command {
	var params kms.ListKeysParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List KMS keys",
		Args:  cobra.ExactArgs(0),
		Long:  "List KMS keys.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := kmsClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListKeysAll(cmd.Context(), &params))
			}
			page, err := c.ListKeys(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Resume token — the last key id from the previous page")
	f.StringVar(&params.Name, "name", "", "Optional substring filter on key name")
	f.StringVar((*string)(&params.State), "state", "", "Optional state filter (enabled / disabled / pending_deletion) (one of: enabled, disabled, pending_deletion)")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newKmsKeyGetCommand builds `basaltic kms key get`.
func newKmsKeyGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <key-id>",
		Short: "Get a KMS key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := kmsClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetKey(cmd.Context(), args[0])
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

// newKmsKeyCreateCommand builds `basaltic kms key create`.
func newKmsKeyCreateCommand(state *cli.State) *cobra.Command {
	var body kms.CreateKeyRequest
	var bodyFile string
	var descriptionFlag string
	var keyUsageFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a KMS key",
		Args:  cobra.ExactArgs(0),
		Long:  "Create a KMS key.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := kmsClient(state)
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
			if cmd.Flags().Changed("key-usage") {
				body.KeyUsage = (*kms.KeyUsage)(&keyUsageFlag)
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
			out, err := c.CreateKey(cmd.Context(), &body, reqOpts...)
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
	f.StringVar((*string)(&body.KeySpec), "key-spec", "", "Key spec (one of: aes-256, rsa-2048, rsa-4096, ecdsa-p256)")
	_ = cmd.MarkFlagRequired("key-spec")
	f.StringVar(&keyUsageFlag, "key-usage", "", "Required for RSA specs (both encrypt_decrypt and sign_verify are valid) (one of: encrypt_decrypt, sign_verify)")
	f.StringVar(&body.Name, "name", "", "Unique per account")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newKmsKeyUpdateCommand builds `basaltic kms key update`.
func newKmsKeyUpdateCommand(state *cli.State) *cobra.Command {
	var body kms.UpdateKeyRequest
	var bodyFile string
	var descriptionFlag string
	var nameFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <key-id>",
		Short: "Update key metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := kmsClient(state)
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
			out, err := c.UpdateKey(cmd.Context(), args[0], &body)
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

// newKmsKeyCancelDeletionCommand builds `basaltic kms key cancel-deletion`.
func newKmsKeyCancelDeletionCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-deletion <key-id>",
		Short: "Cancel a scheduled deletion",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := kmsClient(state)
			if err != nil {
				return err
			}
			out, err := c.CancelKeyDeletion(cmd.Context(), args[0])
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

// newKmsKeyDecryptCommand builds `basaltic kms key decrypt`.
func newKmsKeyDecryptCommand(state *cli.State) *cobra.Command {
	var body kms.DecryptRequest
	var bodyFile string
	var aadFlag string
	var ciphertextFlag string
	cmd := &cobra.Command{
		Use:   "decrypt <key-id>",
		Short: "Decrypt a ciphertext",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := kmsClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if aadFlag != "" {
				body.Aad = []byte(aadFlag)
			}
			if ciphertextFlag != "" {
				body.Ciphertext = []byte(ciphertextFlag)
			}
			out, err := c.Decrypt(cmd.Context(), args[0], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&aadFlag, "aad", "", "Optional base64-encoded AAD")
	f.StringVar(&ciphertextFlag, "ciphertext", "", "Base64-encoded ciphertext produced by Encrypt")
	_ = cmd.MarkFlagRequired("ciphertext")
	return cmd
}

// newKmsKeyDisableCommand builds `basaltic kms key disable`.
func newKmsKeyDisableCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable <key-id>",
		Short: "Disable a key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := kmsClient(state)
			if err != nil {
				return err
			}
			out, err := c.DisableKey(cmd.Context(), args[0])
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

// newKmsKeyEnableCommand builds `basaltic kms key enable`.
func newKmsKeyEnableCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable <key-id>",
		Short: "Enable a disabled key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := kmsClient(state)
			if err != nil {
				return err
			}
			out, err := c.EnableKey(cmd.Context(), args[0])
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

// newKmsKeyEncryptCommand builds `basaltic kms key encrypt`.
func newKmsKeyEncryptCommand(state *cli.State) *cobra.Command {
	var body kms.EncryptRequest
	var bodyFile string
	var aadFlag string
	var plaintextFlag string
	cmd := &cobra.Command{
		Use:   "encrypt <key-id>",
		Short: "Encrypt a payload",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := kmsClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if aadFlag != "" {
				body.Aad = []byte(aadFlag)
			}
			if plaintextFlag != "" {
				body.Plaintext = []byte(plaintextFlag)
			}
			out, err := c.Encrypt(cmd.Context(), args[0], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&aadFlag, "aad", "", "Optional base64-encoded additional authenticated data (AES-GCM AEAD)")
	f.StringVar(&plaintextFlag, "plaintext", "", "Base64-encoded plaintext")
	_ = cmd.MarkFlagRequired("plaintext")
	return cmd
}

// newKmsKeyGenerateDataKeyCommand builds `basaltic kms key generate-data-key`.
func newKmsKeyGenerateDataKeyCommand(state *cli.State) *cobra.Command {
	var body kms.GenerateDataKeyRequest
	var bodyFile string
	var numberOfBytesFlag int
	cmd := &cobra.Command{
		Use:   "generate-data-key <key-id>",
		Short: "Generate a fresh data key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := kmsClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("number-of-bytes") {
				body.NumberOfBytes = &numberOfBytesFlag
			}
			out, err := c.GenerateDataKey(cmd.Context(), args[0], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.IntVar(&numberOfBytesFlag, "number-of-bytes", 0, "Size of the generated data key in bytes: 16 (AES-128), 32 (AES-256, the default) or 64 (HMAC-SHA512) (one of: 16, 32, 64)")
	return cmd
}

// newKmsKeyScheduleDeletionCommand builds `basaltic kms key schedule-deletion`.
func newKmsKeyScheduleDeletionCommand(state *cli.State) *cobra.Command {
	var body kms.ScheduleKeyDeletionRequest
	var bodyFile string
	var pendingWindowInDaysFlag int
	cmd := &cobra.Command{
		Use:   "schedule-deletion <key-id>",
		Short: "Schedule key for deletion",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := kmsClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("pending-window-in-days") {
				body.PendingWindowInDays = &pendingWindowInDaysFlag
			}
			out, err := c.ScheduleKeyDeletion(cmd.Context(), args[0], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.IntVar(&pendingWindowInDaysFlag, "pending-window-in-days", 0, "How long the key sits in pending_deletion before it is hard-deleted")
	return cmd
}

// newKmsKeySignCommand builds `basaltic kms key sign`.
func newKmsKeySignCommand(state *cli.State) *cobra.Command {
	var body kms.SignRequest
	var bodyFile string
	var messageFlag string
	var signingAlgorithmFlag string
	cmd := &cobra.Command{
		Use:   "sign <key-id>",
		Short: "Sign a message",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := kmsClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if messageFlag != "" {
				body.Message = []byte(messageFlag)
			}
			if cmd.Flags().Changed("signing-algorithm") {
				body.SigningAlgorithm = (*kms.SigningAlgorithm)(&signingAlgorithmFlag)
			}
			out, err := c.Sign(cmd.Context(), args[0], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&messageFlag, "message", "", "Base64-encoded message to sign")
	_ = cmd.MarkFlagRequired("message")
	f.StringVar(&signingAlgorithmFlag, "signing-algorithm", "", "Signing algorithm (one of: RSASSA_PSS_SHA_256, RSASSA_PKCS1_V1_5_SHA_256, ECDSA_SHA_256)")
	return cmd
}

// newKmsKeyVerifyCommand builds `basaltic kms key verify`.
func newKmsKeyVerifyCommand(state *cli.State) *cobra.Command {
	var body kms.VerifyRequest
	var bodyFile string
	var messageFlag string
	var signatureFlag string
	var signingAlgorithmFlag string
	cmd := &cobra.Command{
		Use:   "verify <key-id>",
		Short: "Verify a signature",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := kmsClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if messageFlag != "" {
				body.Message = []byte(messageFlag)
			}
			if signatureFlag != "" {
				body.Signature = []byte(signatureFlag)
			}
			if cmd.Flags().Changed("signing-algorithm") {
				body.SigningAlgorithm = (*kms.SigningAlgorithm)(&signingAlgorithmFlag)
			}
			out, err := c.Verify(cmd.Context(), args[0], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&messageFlag, "message", "", "Base64-encoded original message")
	_ = cmd.MarkFlagRequired("message")
	f.StringVar(&signatureFlag, "signature", "", "Base64-encoded signature produced by Sign")
	_ = cmd.MarkFlagRequired("signature")
	f.StringVar(&signingAlgorithmFlag, "signing-algorithm", "", "Signing algorithm (one of: RSASSA_PSS_SHA_256, RSASSA_PKCS1_V1_5_SHA_256, ECDSA_SHA_256)")
	return cmd
}
