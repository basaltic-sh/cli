package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/basaltic-sh/cli/internal/cli"
	"github.com/basaltic-sh/cli/internal/config"
)

func init() { cli.RegisterBuiltin(newConfigCommand) }

func newConfigCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Read and modify the CLI's configuration",
		Long: "Read and modify ~/.config/basaltic/config.yaml.\n\n" +
			"Cached tokens live in a separate file and are managed by `basaltic auth`.",
	}
	cmd.AddCommand(
		newConfigListCommand(state),
		newConfigGetCommand(state),
		newConfigSetCommand(state),
		newConfigUseCommand(state),
		newConfigPathCommand(state),
	)
	return cmd
}

func newConfigListCommand(state *cli.State) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List the configured profiles",
		Aliases: []string{"profiles"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := config.Load()
			if err != nil {
				return err
			}
			if len(file.Profiles) == 0 {
				fmt.Fprintln(state.Printer().Out, "No profiles configured. Run `basaltic auth login`.")
				return nil
			}
			names := make([]string, 0, len(file.Profiles))
			for n := range file.Profiles {
				names = append(names, n)
			}
			sort.Strings(names)

			type row struct {
				Name      string `json:"name"`
				Default   bool   `json:"default"`
				Region    string `json:"region,omitempty"`
				AccountID string `json:"account_id,omitempty"`
				HasKey    bool   `json:"has_credential"`
			}
			rows := make([]row, 0, len(names))
			for _, n := range names {
				p := file.Profiles[n]
				rows = append(rows, row{n, n == file.DefaultProfile, p.Region, p.AccountID, p.APIKey != ""})
			}
			return state.Printer().Iter(sliceIter(rows))
		},
	}
}

func newConfigGetCommand(state *cli.State) *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Read one setting",
		Long:  "Read one setting. Keys are region, api_key, account_id, domain and insecure.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := config.Load()
			if err != nil {
				return err
			}
			_, profile := file.Resolve(profileFlag(cmd))
			v, ok := profileField(profile, args[0])
			if !ok {
				return unknownKey(args[0])
			}
			fmt.Fprintln(state.Printer().Out, v)
			return nil
		},
	}
}

func newConfigSetCommand(state *cli.State) *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Write one setting on the active profile",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := config.Load()
			if err != nil {
				return err
			}
			name, profile := file.Resolve(profileFlag(cmd))
			key, value := args[0], args[1]

			switch strings.ToLower(key) {
			case "region":
				profile.Region = value
			case "api_key", "api-key":
				if _, _, err := (config.Profile{APIKey: value}).Credentials(); err != nil {
					return err
				}
				profile.APIKey = value
				// The cached token was minted by the previous key.
				if err := config.ForgetProfileToken(name); err != nil {
					return err
				}
			case "account_id", "account-id":
				profile.AccountID = value
			case "domain":
				profile.Domain = value
			case "insecure":
				profile.Insecure = value == "true" || value == "1" || value == "yes"
			default:
				return unknownKey(key)
			}

			if file.Profiles == nil {
				file.Profiles = map[string]config.Profile{}
			}
			file.Profiles[name] = profile
			if file.DefaultProfile == "" {
				file.DefaultProfile = name
			}
			if err := config.Save(file); err != nil {
				return err
			}
			fmt.Fprintf(state.Printer().Out, "Set %s on profile %q.\n", key, name)
			return nil
		},
	}
}

func newConfigUseCommand(state *cli.State) *cobra.Command {
	return &cobra.Command{
		Use:   "use <profile>",
		Short: "Make a profile the default",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := config.Load()
			if err != nil {
				return err
			}
			if _, ok := file.Profiles[args[0]]; !ok {
				return fmt.Errorf("no profile named %q; `basaltic config list` shows what exists", args[0])
			}
			file.DefaultProfile = args[0]
			if err := config.Save(file); err != nil {
				return err
			}
			fmt.Fprintf(state.Printer().Out, "Default profile is now %q.\n", args[0])
			return nil
		},
	}
}

func newConfigPathCommand(state *cli.State) *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print where the configuration lives",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, err := config.Path()
			if err != nil {
				return err
			}
			credPath, err := config.CredentialsPath()
			if err != nil {
				return err
			}
			fmt.Fprintf(state.Printer().Out, "config:      %s\ncredentials: %s\n", cfgPath, credPath)
			return nil
		},
	}
}

func profileField(p config.Profile, key string) (string, bool) {
	switch strings.ToLower(key) {
	case "region":
		return p.Region, true
	case "api_key", "api-key":
		return mask(strings.SplitN(p.APIKey, ":", 2)[0]), true
	case "account_id", "account-id":
		return p.AccountID, true
	case "domain":
		return p.Domain, true
	case "insecure":
		return fmt.Sprint(p.Insecure), true
	}
	return "", false
}

func unknownKey(key string) error {
	return fmt.Errorf("unknown setting %q: use region, api_key, account_id, domain or insecure", key)
}

// sliceIter adapts a slice to the iterator the printer renders, so the
// built-in commands produce the same tables the generated ones do.
func sliceIter[T any](items []T) func(func(T, error) bool) {
	return func(yield func(T, error) bool) {
		for _, it := range items {
			if !yield(it, nil) {
				return
			}
		}
	}
}
