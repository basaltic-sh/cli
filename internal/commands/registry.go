package commands

import (
	"fmt"

	"github.com/basaltic-sh/cli/internal/cli"
	"github.com/basaltic-sh/sdk-go/registry"
	"github.com/spf13/cobra"
)

func init() {
	cli.RegisterAt([]string{"registry"}, registryLoginPassword)
}

// Keep stdout suitable for docker login --password-stdin.
func registryLoginPassword(state *cli.State) *cobra.Command {
	return &cobra.Command{
		Use:   "get-login-password",
		Short: "Print a temporary container registry login password",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := state.ServiceSDK("registry", "/v1/authorization-token")
			if err != nil {
				return err
			}
			credential, err := registry.New(cfg).GetAuthorizationToken(cmd.Context())
			if err != nil {
				return err
			}
			if credential == nil || credential.Password == "" {
				return fmt.Errorf("registry returned no login password")
			}
			_, err = fmt.Fprintln(state.Printer().Out, credential.Password)
			return err
		},
	}
}
