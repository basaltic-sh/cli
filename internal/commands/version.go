package commands

import (
	"runtime"

	basaltic "github.com/basaltic-sh/sdk-go"
	"github.com/spf13/cobra"

	"github.com/basaltic-sh/cli/internal/auth"
	"github.com/basaltic-sh/cli/internal/cli"
)

func init() { cli.RegisterBuiltin(newVersionCommand) }

func newVersionCommand(state *cli.State) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the CLI and SDK versions",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return state.Printer().Value(struct {
				CLI  string `json:"cli"`
				SDK  string `json:"sdk"`
				Go   string `json:"go"`
				OS   string `json:"os"`
				Arch string `json:"arch"`
			}{auth.Version, basaltic.Version, runtime.Version(), runtime.GOOS, runtime.GOARCH})
		},
	}
}
