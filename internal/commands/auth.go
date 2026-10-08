// Package commands holds the CLI's own commands — the ones that are about the
// CLI rather than about the platform.
package commands

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

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
	cmd.AddCommand(runLoginCommand(state), newAuthStatusCommand(state), newAuthLogoutCommand(state))
	return cmd
}

// storeAPIKeyProfile is the service-account half of `login`: it stores a key
// pair in a profile and verifies it.
//
// It stayed under the same verb rather than moving to a new one. Whether you
// passed --api-key is unambiguous, the spelling is already in people's CI
// scripts and in the docs, and moving it would break those to buy nothing but
// tidiness.
func storeAPIKeyProfile(state *cli.State, cmd *cobra.Command, profileName, region, accountID string) error {
	{
		{
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
			fmt.Fprintf(state.Printer().Out, "Stored. Profile %q saved to %s\n", name, path)
			return nil
		}
	}
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
			if resolved.UserSession {
				// Naming the principal KIND is the point of this line. A
				// service account cannot create an organization, accept an
				// invitation or switch org, and someone hitting that wall
				// needs to know which of the two they are before the error
				// makes sense.
				fmt.Fprintln(out, "Signed in:   as yourself (browser sign-in)")
				switch {
				case time.Until(resolved.SessionTTL) <= 0:
					fmt.Fprintln(out, "Session:     expired; the next command will renew it")
				default:
					fmt.Fprintf(out, "Session:     expires in %s (renewed automatically)\n", until(resolved.SessionTTL))
				}
			} else {
				fmt.Fprintf(out, "Access key:  %s\n", orDefault(mask(resolved.AccessKeyID), "(a token supplied directly)"))
			}

			if resolved.AccessKeyID != "" {
				expiry, ok := config.CachedTokenExpiry(resolved.Profile, resolved.AccessKeyID)
				switch {
				case !ok:
					fmt.Fprintln(out, "Token:       none cached; the next command will exchange one")
				case time.Until(expiry) <= 0:
					fmt.Fprintln(out, "Token:       cached but expired; the next command will exchange one")
				default:
					fmt.Fprintf(out, "Token:       cached, expires in %s\n", until(expiry))
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

			// A browser sign-in is a SERVER-SIDE session, so deleting the local
			// copy is not logging out — it is losing the receipt. Revoke first,
			// then delete; and delete even when the revoke fails, because a
			// user who typed logout must not be left holding a live credential
			// because the network was down.
			if access, _, endpoint, _, ok := config.LookupSession(name); ok {
				if err := revokeSession(cmd.Context(), endpoint, access); err != nil {
					fmt.Fprintf(out, "Warning: could not revoke the session server-side (%v).\n", err)
					fmt.Fprintln(out, "It will expire on its own; revoke it from the console to be certain.")
				} else {
					fmt.Fprintln(out, "Revoked the browser sign-in session.")
				}
				if _, _, err := config.ForgetSession(name); err != nil {
					return err
				}
				fmt.Fprintf(out, "Signed out of profile %q.\n", name)
			}

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

// until renders a remaining lifetime. Only called for a token that has not
// expired; the expired case reads badly as a duration and is worded
// separately by the caller.
func until(t time.Time) string {
	d := time.Until(t)
	if d < time.Minute {
		return "under a minute"
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

// revokeSession ends a browser sign-in server-side (RFC 7009).
//
// The access token is both the credential and the subject: revoking it ends
// the session behind it, so everything that session issued stops working on
// the next request rather than at expiry.
func revokeSession(ctx context.Context, tokenEndpoint, accessToken string) error {
	if tokenEndpoint == "" || accessToken == "" {
		return errors.New("no session endpoint recorded")
	}
	endpoint := strings.TrimSuffix(tokenEndpoint, "/token") + "/revoke"
	form := url.Values{"token": {accessToken}, "token_type_hint": {"access_token"}}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// RFC 7009 says a revocation endpoint answers 200 whether or not anything
	// was revoked, so anything else is a real failure.
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return nil
}
