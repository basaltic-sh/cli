// Code generated from the Basaltic SDK manifest (api.json). DO NOT EDIT.
//
// Regenerate with:
//
//	make generate SDK=/path/to/sdk-go

package generated

import (
	"github.com/spf13/cobra"

	"github.com/basaltic-sh/sdk-go/catalog"

	"github.com/basaltic-sh/cli/internal/cli"
)

func init() { cli.RegisterService(newCatalogCommand) }

// newCatalogCommand builds `basaltic catalog`.
func newCatalogCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "catalog",
		Short: "Public global platform discovery",
	}
	cmd.AddCommand(newCatalogRegionCommand(state))
	return cmd
}

// catalogClient builds the public service client without credentials.
func catalogClient(state *cli.State, path string) (*catalog.Client, error) {
	cfg, err := state.PublicSDK()
	if err != nil {
		return nil, err
	}
	return catalog.New(cfg), nil
}

// newCatalogRegionCommand builds `basaltic catalog region`.
func newCatalogRegionCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "region",
		Short:   "Regions",
		Aliases: []string{"regions"},
	}
	cmd.AddCommand(newCatalogRegionListCommand(state))
	cmd.AddCommand(newCatalogRegionGetCommand(state))
	return cmd
}

// newCatalogRegionListCommand builds `basaltic catalog region list`.
func newCatalogRegionListCommand(state *cli.State) *cobra.Command {
	var params catalog.ListRegionsParams
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List regions",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := catalogClient(state, "/v1/regions")
			if err != nil {
				return err
			}
			out, err := c.ListRegions(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact returned CRN, combined with name using AND before pagination")
	f.StringVar(&params.Name, "name", "", "Exact resource name, combined with crn using AND before pagination")
	return cmd
}

// newCatalogRegionGetCommand builds `basaltic catalog region get`.
func newCatalogRegionGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <code>",
		Short: "Get a region",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := catalogClient(state, "/v1/regions/{code}")
			if err != nil {
				return err
			}
			out, err := c.GetRegion(cmd.Context(), args[0])
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
