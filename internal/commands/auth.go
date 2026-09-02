// Package commands holds the CLI's own commands — the ones that are about the
// CLI rather than about the platform.
package commands

import (
	"context"
	"fmt"

	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/basaltic-sh/cli/internal/auth"
	"github.com/basaltic-sh/cli/internal/cli"
	"github.com/basaltic-sh/cli/internal/config"
)

func init() { cli.RegisterBuiltin(newAuthCommand) }

func newAuthCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate the CLI and inspect its credentials",
	}
	cmd.AddCommand(newAuthLoginCommand(state), newAuthStatusCommand(state), newAuthLogoutCommand(state))
	return cmd
}

func newAuthLoginCommand(state *cli.State) *cobra.Command {
	var profileName, region, accountID string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Store a service account credential in a profile",
		Long: "Store a service account's access key pair in a profile, and verify it.\n\n" +
			"The key pair stays the one long-lived credential. The CLI exchanges it for\n" +
			"a short-lived token on each run and caches that token separately, in\n" +
			"credentials.yaml, so it never travels with a config file you copy.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Read the CLI-wide --api-key rather than declaring another of
			// the same name, which would shadow it.
			apiKey, _ := cmd.Flags().GetString("api-key")
			if apiKey == "" {
				return fmt.Errorf("--api-key is required, written as ACCESS_KEY_ID:SECRET_ACCESS_KEY")
			}
			if _, _, err := (config.Profile{APIKey: apiKey}).Credentials(); err != nil {
				return err
			}

			file, err := config.Load()
			if err != nil {
				return err
			}
			name := profileName
			if name == "" {
				name = file.DefaultProfile
			}
			if name == "" {
				name = "default"
			}

			p := file.Profiles[name]
			p.APIKey = apiKey
			if region != "" {
				p.Region = region
			}
			if accountID != "" {
				p.AccountID = accountID
			}
			file.Profiles[name] = p
			if file.DefaultProfile == "" {
				file.DefaultProfile = name
			}
			if err := config.Save(file); err != nil {
				return err
			}

			// A stored credential replaces whatever the profile had, so any
			// token cached against the old key is dead. Clearing it here
			// avoids a confusing hour where login "worked" but every command
			// still presents the previous identity.
			if err := config.ForgetProfileToken(name); err != nil {
				return err
			}

			// Verify before reporting success: a credential that is stored
			// but refused is worse than one that was never stored, because
			// the failure surfaces later and somewhere else.
			resolved, err := auth.Resolve(auth.Options{Profile: name})
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
			defer cancel()
			if _, err := resolved.Config.TokenSource.Token(ctx); err != nil {
				return fmt.Errorf("the credential was stored in profile %q but was refused: %w", name, err)
			}

			path, _ := config.Path()
			fmt.Fprintf(state.Printer().Out, "Signed in. Profile %q saved to %s\n", name, path)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&profileName, "profile-name", "", "profile to write (default: the current default profile)")
	f.StringVar(&region, "set-region", "", "region to store in the profile")
	f.StringVar(&accountID, "set-account-id", "", "account to store in the profile")
	return cmd
}

func newAuthStatusCommand(state *cli.State) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show which credential, region and account are in effect",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			resolved, err := state.Resolved()
			if err != nil {
				return err
			}
			out := state.Printer().Out

			fmt.Fprintf(out, "Profile:     %s\n", resolved.Profile)
			fmt.Fprintf(out, "Region:      %s\n", orNone(resolved.Region))
			fmt.Fprintf(out, "Account:     %s\n", orDefault(resolved.AccountID, "(the credential's own)"))
			fmt.Fprintf(out, "Access key:  %s\n", orDefault(mask(resolved.AccessKeyID), "(a token supplied directly)"))

			if resolved.AccessKeyID != "" {
				if expiry, ok := config.CachedTokenExpiry(resolved.Profile, resolved.AccessKeyID); ok {
					fmt.Fprintf(out, "Token:       cached, expires in %s\n", until(expiry))
				} else {
					fmt.Fprintln(out, "Token:       none cached; the next command will exchange one")
				}
			}

			ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
			defer cancel()
			if _, err := resolved.Config.TokenSource.Token(ctx); err != nil {
				fmt.Fprintf(out, "Status:      NOT working — %v\n", err)
				return err
			}
			fmt.Fprintln(out, "Status:      working")
			return nil
		},
	}
}

func newAuthLogoutCommand(state *cli.State) *cobra.Command {
	var forget bool
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Discard the cached token, and optionally the stored credential",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := config.Load()
			if err != nil {
				return err
			}
			name, profile := file.Resolve(profileFlag(cmd))

			// Clear the token first and unconditionally. A token outlives the
			// key that minted it, so clearing only the key would leave a
			// "logged out" CLI authenticating for up to an hour.
			if err := config.ForgetProfileToken(name); err != nil {
				return err
			}
			out := state.Printer().Out
			fmt.Fprintf(out, "Cleared the cached token for profile %q.\n", name)

			if forget {
				profile.APIKey = ""
				file.Profiles[name] = profile
				if err := config.Save(file); err != nil {
					return err
				}
				fmt.Fprintf(out, "Removed the stored credential from profile %q.\n", name)
			} else {
				fmt.Fprintln(out, "The stored credential is unchanged; pass --forget-credential to remove it too.")
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&forget, "forget-credential", false, "also remove the api_key from the profile")
	return cmd
}

func profileFlag(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString("profile")
	return v
}

func mask(key string) string {
	if key == "" {
		return ""
	}
	if len(key) <= 8 {
		return strings.Repeat("*", len(key))
	}
	return key[:6] + strings.Repeat("*", len(key)-10) + key[len(key)-4:]
}

func until(t time.Time) string {
	d := time.Until(t)
	if d <= 0 {
		return "already (it will be re-exchanged)"
	}
	return d.Round(time.Minute).String()
}

func orNone(s string) string { return orDefault(s, "(not set)") }

func orDefault(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}
