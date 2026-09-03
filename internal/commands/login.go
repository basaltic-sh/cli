package commands

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"github.com/basaltic-sh/cli/internal/cli"
	"github.com/basaltic-sh/cli/internal/config"
	"github.com/basaltic-sh/cli/internal/oauth"
	basaltic "github.com/basaltic-sh/sdk-go"
)

func init() { cli.RegisterBuiltin(newLoginCommand) }

// `basaltic login` is the top-level spelling, because that is what the word
// means in every other CLI a person will have used. `basaltic auth login` is
// the same command; auth/ is where the rest of the credential verbs live and
// splitting the pair across two places would be worse than repeating it.
func newLoginCommand(state *cli.State) *cobra.Command {
	cmd := runLoginCommand(state)
	cmd.Use = "login"
	return cmd
}

// runLoginCommand builds the interactive sign-in.
//
// WHY THIS EXISTS ALONGSIDE THE KEY PAIR. A service account cannot create an
// organization, accept an invitation, or switch organization — those routes
// require a user principal. Before this, the CLI could only ever be a service
// account, so those operations had no CLI at all.
func runLoginCommand(state *cli.State) *cobra.Command {
	var profileName, region, accountID string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Sign in as yourself through a browser",
		Long: "Sign in as yourself. Opens a browser, asks you to approve this CLI and pick an\n" +
			"organization, and stores a short-lived token that acts as you.\n\n" +
			"Use this when you are a person at a terminal. For a program that has to run\n" +
			"without you — CI, a cron job — store a service account key instead with\n" +
			"`basaltic auth login --api-key ACCESS_KEY_ID:SECRET`.\n\n" +
			"The token and its refresh token go to credentials.yaml at mode 0600, never to\n" +
			"config.yaml, so a config file stays safe to copy between machines.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// The key-pair path keeps working under this verb. It is
			// unambiguous — you either passed a key or you did not — and it is
			// the spelling already in the docs and in people's CI scripts.
			if apiKey, _ := cmd.Flags().GetString("api-key"); apiKey != "" {
				return storeAPIKeyProfile(state, cmd, profileName, region, accountID)
			}

			file, err := config.Load()
			if err != nil {
				return err
			}
			name := profileName
			if name == "" {
				name, _ = file.Resolve(profileFlag(cmd))
			}
			profile := file.Profiles[name]
			if region != "" {
				profile.Region = region
			}
			if profile.Region == "" {
				profile.Region = regionFlag(cmd)
			}
			if profile.Region == "" {
				return fmt.Errorf(
					"profile %q has no region.\n"+
						"Pass --set-region, or set one with `basaltic config set %s.region <region>`", name, name)
			}

			client := &http.Client{Timeout: 30 * time.Second}
			if profile.Insecure {
				client = insecureHTTPClient()
			}

			iamURL, err := iamEndpoint(profile)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(cmd.Context(), 6*time.Minute)
			defer cancel()

			meta, err := oauth.DiscoverMetadata(ctx, client, iamURL)
			if err != nil {
				return err
			}
			// Discovered, not assumed. A deployment with no console configured
			// to approve at is not broken — it simply cannot serve this flow,
			// and saying so beats sending someone to a URL that 404s.
			if !meta.SupportsInteractiveLogin() {
				return fmt.Errorf(
					"%s does not offer interactive sign-in.\n"+
						"Use a service account key: `basaltic auth login --api-key ACCESS_KEY_ID:SECRET`", iamURL)
			}

			out := state.Printer().Out
			toks, err := oauth.Login(ctx, client, meta, func(u string) {
				// Printed BEFORE the browser is attempted, and always. An SSH
				// session has no browser to open, and a flow that waits
				// silently on a page nobody can see looks like a hang.
				fmt.Fprintf(out, "Open this URL to sign in:\n\n  %s\n\nWaiting…\n", u)
			})
			if err != nil {
				return err
			}

			if err := config.StoreSession(name, toks.AccessToken, toks.RefreshToken, meta.TokenEndpoint, toks.ExpiresAt); err != nil {
				return err
			}
			// Persist the profile so the next command resolves the same region
			// without the flag, and so a brand-new profile exists at all.
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

			path, _ := config.CredentialsPath()
			fmt.Fprintf(out, "\nSigned in. Session stored in %s (profile %q).\n", path, name)
			fmt.Fprintln(out, "Run `basaltic auth status` to see who you are.")
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&profileName, "profile-name", "", "profile to write (default: the current default profile)")
	f.StringVar(&region, "set-region", "", "region to store in the profile")
	f.StringVar(&accountID, "set-account-id", "", "account to store in the profile (--api-key only)")
	return cmd
}

// iamEndpoint resolves the IAM base URL the same way every other call does, so
// a profile pointing at a development rig signs in against that rig.
func iamEndpoint(profile config.Profile) (string, error) {
	opts := []basaltic.Option{basaltic.WithRegion(profile.Region)}
	if profile.Domain != "" {
		opts = append(opts, basaltic.WithDomain(profile.Domain))
	}
	for service, endpoint := range profile.Endpoints {
		opts = append(opts, basaltic.WithServiceEndpoint(service, endpoint))
	}
	// A token source is not needed to resolve an endpoint, but NewConfig wants
	// a complete configuration; the anonymous source is never called.
	opts = append(opts, basaltic.WithAccessToken("unused-for-endpoint-resolution"))
	cfg, err := basaltic.NewConfig(context.Background(), opts...)
	if err != nil {
		return "", err
	}
	return cfg.EndpointResolver.ResolveEndpoint("iam", profile.Region)
}

func regionFlag(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString("region")
	return v
}

// insecureHTTPClient disables TLS verification, for development rigs with
// self-signed certificates. It mirrors what the profile already does for every
// other call, so the sign-in path does not become the one place a rig fails.
func insecureHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}
