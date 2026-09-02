// Code generated from the Basaltic SDK manifest (api.json). DO NOT EDIT.
//
// Regenerate with:
//
//	make generate SDK=/path/to/sdk-go

package generated

import (
	"github.com/spf13/cobra"

	"github.com/basaltic-sh/sdk-go/quota"

	"github.com/basaltic-sh/cli/internal/cli"
)

func init() { cli.RegisterService(newQuotaCommand) }

// newQuotaCommand builds `basaltic quota`.
func newQuotaCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quota",
		Short: "Account quotas",
	}
	cmd.AddCommand(newQuotaQuotaListCommand(state))
	return cmd
}

// quotaClient builds the service client, resolving credentials on first use.
func quotaClient(state *cli.State) (*quota.Client, error) {
	cfg, err := state.SDK()
	if err != nil {
		return nil, err
	}
	return quota.New(cfg), nil
}

// newQuotaQuotaListCommand builds `basaltic quota quota list`.
func newQuotaQuotaListCommand(state *cli.State) *cobra.Command {
	var params quota.ListQuotasParams
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List quotas",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := quotaClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListQuotas(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.Region, "for-region", "", "Region to filter regional quotas by")
	return cmd
}
