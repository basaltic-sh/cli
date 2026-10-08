// Code generated from the Basaltic SDK manifest (api.json). DO NOT EDIT.
//
// Regenerate with:
//
//	make generate SDK=/path/to/sdk-go

package generated

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/basaltic-sh/sdk-go/billing"

	"github.com/basaltic-sh/cli/internal/cli"
)

func init() { cli.RegisterService(newBillingCommand) }

// newBillingCommand builds `basaltic billing`.
func newBillingCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "billing",
		Short: "Invoices, credits, payments and prices",
	}
	cmd.AddCommand(newBillingCreditCommand(state))
	cmd.AddCommand(newBillingFiscalInvoiceCommand(state))
	cmd.AddCommand(newBillingInvoiceCommand(state))
	cmd.AddCommand(newBillingPaymentCommand(state))
	cmd.AddCommand(newBillingPriceCommand(state))
	cmd.AddCommand(newBillingProfileCommand(state))
	cmd.AddCommand(newBillingTransactionCommand(state))
	cmd.AddCommand(newBillingUsageCommand(state))
	return cmd
}

// billingClient builds the service client, resolving credentials on first use.
func billingClient(state *cli.State, path string) (*billing.Client, error) {
	cfg, err := state.ServiceSDK("billing", path)
	if err != nil {
		return nil, err
	}
	return billing.New(cfg), nil
}

// newBillingCreditCommand builds `basaltic billing credit`.
func newBillingCreditCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "credit",
		Short:   "Credits",
		Aliases: []string{"credits"},
	}
	cmd.AddCommand(newBillingCreditListCommand(state))
	return cmd
}

// newBillingCreditListCommand builds `basaltic billing credit list`.
func newBillingCreditListCommand(state *cli.State) *cobra.Command {
	var params billing.ListCreditsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List credit grants",
		Args:  cobra.ExactArgs(0),
		Long:  "List credit grants.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := billingClient(state, "/v1/credits")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListCreditsAll(cmd.Context(), &params))
			}
			page, err := c.ListCredits(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact organization-scoped credit CRN")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newBillingFiscalInvoiceCommand builds `basaltic billing fiscal-invoice`.
func newBillingFiscalInvoiceCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fiscal-invoice",
		Short:   "Fiscal invoices",
		Aliases: []string{"fiscal-invoices"},
	}
	cmd.AddCommand(newBillingFiscalInvoiceListCommand(state))
	cmd.AddCommand(newBillingFiscalInvoiceGetXmlCommand(state))
	return cmd
}

// newBillingFiscalInvoiceListCommand builds `basaltic billing fiscal-invoice list`.
func newBillingFiscalInvoiceListCommand(state *cli.State) *cobra.Command {
	var params billing.ListFiscalInvoicesParams
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List fiscal invoice issuance and delivery status",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := billingClient(state, "/v1/fiscal-invoices")
			if err != nil {
				return err
			}
			out, err := c.ListFiscalInvoices(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.Invoice, "invoice", "", "Canonical billing invoice CRN")
	return cmd
}

// newBillingFiscalInvoiceGetXmlCommand builds `basaltic billing fiscal-invoice get-xml`.
func newBillingFiscalInvoiceGetXmlCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-xml <document-id>",
		Short: "Download issued NFS-e XML",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := billingClient(state, "/v1/fiscal-invoices/{document_id}/xml")
			if err != nil {
				return err
			}
			stream, err := c.GetFiscalInvoiceXml(cmd.Context(), args[0])
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

// newBillingInvoiceCommand builds `basaltic billing invoice`.
func newBillingInvoiceCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "invoice",
		Short:   "Invoices",
		Aliases: []string{"invoices"},
	}
	cmd.AddCommand(newBillingInvoiceListCommand(state))
	cmd.AddCommand(newBillingInvoiceGetCommand(state))
	cmd.AddCommand(newBillingInvoiceGetPdfCommand(state))
	return cmd
}

// newBillingInvoiceListCommand builds `basaltic billing invoice list`.
func newBillingInvoiceListCommand(state *cli.State) *cobra.Command {
	var params billing.ListInvoicesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List invoices",
		Args:  cobra.ExactArgs(0),
		Long:  "List invoices.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := billingClient(state, "/v1/invoices")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListInvoicesAll(cmd.Context(), &params))
			}
			page, err := c.ListInvoices(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact organization-scoped invoice CRN")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newBillingInvoiceGetCommand builds `basaltic billing invoice get`.
func newBillingInvoiceGetCommand(state *cli.State) *cobra.Command {
	var scope billing.ListInvoicesParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Get an invoice with its line items",
		Args:  cobra.ExactArgs(1),
		Long:  "Get an invoice with its line items.\n\n<ref> is its id or its CRN. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.\n\nThis resource has no name.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := billingClient(state, "/v1/invoices/{invoice_id}")
			if err != nil {
				return err
			}
			out, err := c.GetInvoiceByReference(cmd.Context(), args[0], &scope)
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

// newBillingInvoiceGetPdfCommand builds `basaltic billing invoice get-pdf`.
func newBillingInvoiceGetPdfCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-pdf <invoice-id>",
		Short: "Download an invoice as a PDF statement",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := billingClient(state, "/v1/invoices/{invoice_id}/pdf")
			if err != nil {
				return err
			}
			stream, err := c.GetInvoicePDF(cmd.Context(), args[0])
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

// newBillingPaymentCommand builds `basaltic billing payment`.
func newBillingPaymentCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "payment",
		Short:   "Payments",
		Aliases: []string{"payments"},
	}
	cmd.AddCommand(newBillingPaymentListCommand(state))
	return cmd
}

// newBillingPaymentListCommand builds `basaltic billing payment list`.
func newBillingPaymentListCommand(state *cli.State) *cobra.Command {
	var params billing.ListPaymentsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List invoice payments",
		Args:  cobra.ExactArgs(0),
		Long:  "List invoice payments.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := billingClient(state, "/v1/payments")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListPaymentsAll(cmd.Context(), &params))
			}
			page, err := c.ListPayments(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact organization-scoped payment CRN")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newBillingPriceCommand builds `basaltic billing price`.
func newBillingPriceCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "price",
		Short:   "Prices",
		Aliases: []string{"prices"},
	}
	cmd.AddCommand(newBillingPriceListCommand(state))
	return cmd
}

// newBillingPriceListCommand builds `basaltic billing price list`.
func newBillingPriceListCommand(state *cli.State) *cobra.Command {
	var params billing.ListPricesParams
	var atFlag string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List catalog prices",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := billingClient(state, "/v1/prices")
			if err != nil {
				return err
			}
			if atFlag != "" {
				parsed, err := parseTime(atFlag)
				if err != nil {
					return fmt.Errorf("--at: %w", err)
				}
				params.At = parsed
			}
			out, err := c.ListPrices(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&atFlag, "at", "", "Read the catalog as of this instant instead of now, for showing a historical price (RFC 3339)")
	f.StringVar(&params.Family, "family", "", "Only SKUs whose metadata.family matches — how the managed products are separated from the general compute flavors")
	f.StringVar(&params.ResourceType, "resource-type", "", "Resource type")
	f.StringVar(&params.Service, "service", "", "Only SKUs billed by this service")
	f.StringVar(&params.Sku, "sku", "", "Exactly one SKU")
	return cmd
}

// newBillingProfileCommand builds `basaltic billing profile`.
func newBillingProfileCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "profile",
		Short:   "Profiles",
		Aliases: []string{"profiles"},
	}
	cmd.AddCommand(newBillingProfileListCommand(state))
	cmd.AddCommand(newBillingProfileUpdateCommand(state))
	return cmd
}

// newBillingProfileListCommand builds `basaltic billing profile list`.
func newBillingProfileListCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Read the organization billing profile",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := billingClient(state, "/v1/profile")
			if err != nil {
				return err
			}
			out, err := c.GetBillingProfile(cmd.Context())
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

// newBillingProfileUpdateCommand builds `basaltic billing profile update`.
func newBillingProfileUpdateCommand(state *cli.State) *cobra.Command {
	var body billing.BillingProfile
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Save organization billing details",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := billingClient(state, "/v1/profile")
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			out, err := c.UpdateBillingProfile(cmd.Context(), &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.City, "city", "", "City")
	f.StringVar(&body.CompanyName, "company-name", "", "Full legal name of the individual or company")
	f.StringVar(&body.Complement, "complement", "", "Complement")
	f.StringVar(&body.Country, "country", "", "ISO 3166-1 alpha-2 country code")
	f.StringVar(&body.CustomerType, "customer-type", "", "Customer type (one of: , individual, company)")
	f.StringVar(&body.Email, "email", "", "Billing email for fiscal invoice delivery")
	f.StringVar(&body.ForeignTaxID, "foreign-tax-id", "", "Foreign identifier; not validated as a Brazilian document")
	f.StringSliceVar(&body.MissingFields, "missing-fields", nil, "Missing fields")
	f.StringVar(&body.MunicipalityCode, "municipality-code", "", "Seven-digit IBGE municipality code, required for a Brazilian recipient")
	f.StringVar(&body.Neighborhood, "neighborhood", "", "Neighborhood")
	f.StringVar(&body.NoTaxIDReason, "no-tax-id-reason", "", "Required for a foreign recipient without a tax identifier")
	f.StringVar(&body.Phone, "phone", "", "Phone")
	f.StringVar(&body.PostalCode, "postal-code", "", "Eight-digit CEP for Brazil; optional international postal code abroad")
	f.BoolVar(&body.Ready, "ready", false, "Ready")
	f.StringVar(&body.State, "state", "", "Two-letter UF for Brazil; free-form state/province abroad")
	f.StringVar(&body.StreetName, "street-name", "", "Street name")
	f.StringVar(&body.StreetNumber, "street-number", "", "Street number")
	f.StringVar(&body.TaxID, "tax-id", "", "CPF for a Brazilian individual or CNPJ for a Brazilian company")
	return cmd
}

// newBillingTransactionCommand builds `basaltic billing transaction`.
func newBillingTransactionCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "transaction",
		Short:   "Transactions",
		Aliases: []string{"transactions"},
	}
	cmd.AddCommand(newBillingTransactionListCommand(state))
	return cmd
}

// newBillingTransactionListCommand builds `basaltic billing transaction list`.
func newBillingTransactionListCommand(state *cli.State) *cobra.Command {
	var params billing.ListTransactionsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List ledger transactions",
		Args:  cobra.ExactArgs(0),
		Long:  "List ledger transactions.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := billingClient(state, "/v1/transactions")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListTransactionsAll(cmd.Context(), &params))
			}
			page, err := c.ListTransactions(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact organization-scoped transaction CRN")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newBillingUsageCommand builds `basaltic billing usage`.
func newBillingUsageCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "usage",
		Short:   "Usages",
		Aliases: []string{"usages"},
	}
	cmd.AddCommand(newBillingUsageListCommand(state))
	return cmd
}

// newBillingUsageListCommand builds `basaltic billing usage list`.
func newBillingUsageListCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Get month-to-date usage total",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := billingClient(state, "/v1/usage")
			if err != nil {
				return err
			}
			out, err := c.GetCurrentUsage(cmd.Context())
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
