package commands

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/spf13/cobra"

	"github.com/basaltic-sh/cli/internal/auth"
	"github.com/basaltic-sh/cli/internal/cli"
	"github.com/basaltic-sh/cli/internal/selfupdate"
)

func init() { cli.RegisterBuiltin(newUpgradeCommand) }

func newUpgradeCommand(state *cli.State) *cobra.Command {
	var checkOnly bool
	var pinned string
	var force bool

	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Install the latest version of the CLI",
		Long: "Install the latest version of the CLI, replacing the running binary.\n\n" +
			"The download is verified against the checksum published with the release\n" +
			"before anything is written, and the binary is replaced by a rename, so an\n" +
			"interrupted upgrade leaves either the old version or the new one.\n\n" +
			"The CLI never upgrades itself. It will mention when a newer version exists;\n" +
			"installing it is always this command. Set " + selfupdate.EnvDisable + "=1 to\n" +
			"silence the notice.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			out := state.Printer().Out
			hc := &http.Client{Timeout: 5 * time.Minute}

			ctx, cancel := context.WithTimeout(cmd.Context(), 10*time.Minute)
			defer cancel()

			var rel *selfupdate.Release
			var err error
			if pinned != "" {
				rel, err = selfupdate.AtVersion(ctx, hc, pinned)
			} else {
				rel, err = selfupdate.Latest(ctx, hc)
			}
			if err != nil {
				return err
			}

			current := auth.Version
			switch {
			case pinned != "":
				// An explicit version is an instruction, not a suggestion:
				// downgrading to reproduce a bug is a real reason to run this.
			case !selfupdate.IsNewer(current, rel.Version()):
				if !force {
					fmt.Fprintf(out, "Already on the latest version (%s).\n", current)
					return nil
				}
			}

			if checkOnly {
				fmt.Fprintf(out, "Installed: %s\nAvailable: %s\n%s\n", current, rel.Version(), rel.HTMLURL)
				return nil
			}

			execPath, err := selfupdate.ResolvePath()
			if err != nil {
				return err
			}
			// Checked before downloading: there is no point spending a
			// multi-megabyte transfer to discover the file cannot be replaced.
			if err := selfupdate.CheckWritable(execPath); err != nil {
				return err
			}

			fmt.Fprintf(out, "Upgrading %s → %s for %s/%s\n", current, rel.Version(), runtime.GOOS, runtime.GOARCH)
			if err := selfupdate.Install(ctx, hc, rel, execPath); err != nil {
				return err
			}
			fmt.Fprintf(out, "Installed %s to %s\n", rel.Version(), execPath)
			if rel.HTMLURL != "" {
				fmt.Fprintf(out, "Release notes: %s\n", rel.HTMLURL)
			}
			return nil
		},
	}

	f := cmd.Flags()
	f.BoolVar(&checkOnly, "check", false, "report whether a newer version exists without installing it")
	f.StringVar(&pinned, "version", "", "install this exact version, including an older one")
	f.BoolVar(&force, "force", false, "reinstall even when already on the latest version")
	return cmd
}
