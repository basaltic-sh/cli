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
	cmd.AddCommand(newBillingInvoiceCommand(state))
	cmd.AddCommand(newBillingPaymentCommand(state))
	cmd.AddCommand(newBillingPriceCommand(state))
	cmd.AddCommand(newBillingTransactionCommand(state))
	cmd.AddCommand(newBillingUsageCommand(state))
	return cmd
}

// billingClient builds the service client, resolving credentials on first use.
func billingClient(state *cli.State) (*billing.Client, error) {
	cfg, err := state.SDK()
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
			c, err := billingClient(state)
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
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
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
			c, err := billingClient(state)
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
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newBillingInvoiceGetCommand builds `basaltic billing invoice get`.
func newBillingInvoiceGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <invoice-id>",
		Short: "Get an invoice with its line items",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := billingClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetInvoice(cmd.Context(), args[0])
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
			c, err := billingClient(state)
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
			c, err := billingClient(state)
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
			c, err := billingClient(state)
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
			c, err := billingClient(state)
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
			c, err := billingClient(state)
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
