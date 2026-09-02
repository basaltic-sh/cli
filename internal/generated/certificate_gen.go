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
	"github.com/basaltic-sh/sdk-go/certificate"

	"github.com/basaltic-sh/cli/internal/cli"
)

func init() { cli.RegisterService(newCertificateCommand) }

// newCertificateCommand builds `basaltic certificate`.
func newCertificateCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "certificate",
		Short:   "TLS certificates",
		Aliases: []string{"cert"},
		Long:    "TLS certificates.\n\nThis is a regional service: it acts in the region from --region, the\nBASALTIC_REGION environment variable, or the profile.",
	}
	cmd.AddCommand(newCertificateCertificateListCommand(state))
	cmd.AddCommand(newCertificateCertificateGetCommand(state))
	cmd.AddCommand(newCertificateCertificateCreateCommand(state))
	cmd.AddCommand(newCertificateCertificateDeleteCommand(state))
	cmd.AddCommand(newCertificateCertificateGetMaterialCommand(state))
	cmd.AddCommand(newCertificateCertificateRevokeCommand(state))
	return cmd
}

// certificateClient builds the service client, resolving credentials on first use.
func certificateClient(state *cli.State) (*certificate.Client, error) {
	cfg, err := state.SDK()
	if err != nil {
		return nil, err
	}
	return certificate.New(cfg), nil
}

// newCertificateCertificateListCommand builds `basaltic certificate certificate list`.
func newCertificateCertificateListCommand(state *cli.State) *cobra.Command {
	var params certificate.ListCertificatesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List certificates",
		Args:  cobra.ExactArgs(0),
		Long:  "List certificates.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := certificateClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListCertificatesAll(cmd.Context(), &params))
			}
			page, err := c.ListCertificates(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Limit")
	f.StringVar(&params.Marker, "marker", "", "Resume token — the last certificate id from the previous page")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newCertificateCertificateGetCommand builds `basaltic certificate certificate get`.
func newCertificateCertificateGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <certificate-id>",
		Short: "Get certificate",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := certificateClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetCertificate(cmd.Context(), args[0])
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

// newCertificateCertificateCreateCommand builds `basaltic certificate certificate create`.
func newCertificateCertificateCreateCommand(state *cli.State) *cobra.Command {
	var body certificate.CertificateIssueRequest
	var bodyFile string
	var certificatePemFlag string
	var chainPemFlag string
	var keyAlgorithmFlag string
	var privateKeyPemFlag string
	var sourceFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create certificate",
		Args:  cobra.ExactArgs(0),
		Long:  "Create certificate.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := certificateClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("certificate-pem") {
				body.CertificatePEM = &certificatePemFlag
			}
			if cmd.Flags().Changed("chain-pem") {
				body.ChainPEM = &chainPemFlag
			}
			if cmd.Flags().Changed("key-algorithm") {
				body.KeyAlgorithm = (*certificate.CertificateKeyAlgorithm)(&keyAlgorithmFlag)
			}
			if cmd.Flags().Changed("private-key-pem") {
				body.PrivateKeyPEM = &privateKeyPemFlag
			}
			if cmd.Flags().Changed("source") {
				body.Source = (*certificate.CertificateSource)(&sourceFlag)
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
			out, err := c.CreateCertificate(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&certificatePemFlag, "certificate-pem", "", "PEM-encoded leaf certificate")
	f.StringVar(&chainPemFlag, "chain-pem", "", "PEM-encoded intermediate chain (optional when source=uploaded)")
	f.StringSliceVar(&body.Domains, "domains", nil, "Capped at 100 to stay inside the certificate authority's per-order limits")
	_ = cmd.MarkFlagRequired("domains")
	f.StringVar(&keyAlgorithmFlag, "key-algorithm", "", "Key algorithm (one of: ecdsa-p256, ecdsa-p384, rsa-2048, rsa-4096)")
	f.StringVar(&body.Name, "name", "", "Unique per account")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&privateKeyPemFlag, "private-key-pem", "", "PEM-encoded private key")
	f.StringVar(&sourceFlag, "source", "", "Defaults to \"acme\" — issued by the platform CA (one of: acme, uploaded)")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newCertificateCertificateDeleteCommand builds `basaltic certificate certificate delete`.
func newCertificateCertificateDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <certificate-id>",
		Short: "Delete certificate",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := certificateClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteCertificate(cmd.Context(), args[0]); err != nil {
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

// newCertificateCertificateGetMaterialCommand builds `basaltic certificate certificate get-material`.
func newCertificateCertificateGetMaterialCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-material <certificate-id>",
		Short: "Fetch certificate material (leaf, chain, private key)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := certificateClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetCertificateMaterial(cmd.Context(), args[0])
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

// newCertificateCertificateRevokeCommand builds `basaltic certificate certificate revoke`.
func newCertificateCertificateRevokeCommand(state *cli.State) *cobra.Command {
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "revoke <certificate-id>",
		Short: "Revoke certificate",
		Args:  cobra.ExactArgs(1),
		Long:  "Revoke certificate.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := certificateClient(state)
			if err != nil {
				return err
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.RevokeCertificate(cmd.Context(), args[0], reqOpts...)
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
