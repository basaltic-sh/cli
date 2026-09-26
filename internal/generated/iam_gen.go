// Code generated from the Basaltic SDK manifest (api.json). DO NOT EDIT.
//
// Regenerate with:
//
//	make generate SDK=/path/to/sdk-go

package generated

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	basaltic "github.com/basaltic-sh/sdk-go"
	"github.com/basaltic-sh/sdk-go/iam"

	"github.com/basaltic-sh/cli/internal/cli"
)

func init() { cli.RegisterService(newIamCommand) }

// newIamCommand builds `basaltic iam`.
func newIamCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "iam",
		Short: "Account roles, policies, service accounts and sessions",
	}
	cmd.AddCommand(newIamOauthCommand(state))
	cmd.AddCommand(newIamPolicyCommand(state))
	cmd.AddCommand(newIamRegionCommand(state))
	cmd.AddCommand(newIamRoleCommand(state))
	cmd.AddCommand(newIamServiceAccountCommand(state))
	cmd.AddCommand(newIamStsSessionCommand(state))
	cmd.AddCommand(newIamTokenCommand(state))
	return cmd
}

// iamClient builds the service client, resolving credentials on first use.
func iamClient(state *cli.State) (*iam.Client, error) {
	cfg, err := state.SDK()
	if err != nil {
		return nil, err
	}
	return iam.New(cfg), nil
}

// newIamOauthCommand builds `basaltic iam oauth`.
func newIamOauthCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "oauth",
		Short:   "Oauths",
		Aliases: []string{"oauths"},
	}
	cmd.AddCommand(newIamOauthAuthorizeCommand(state))
	return cmd
}

// newIamOauthAuthorizeCommand builds `basaltic iam oauth authorize`.
func newIamOauthAuthorizeCommand(state *cli.State) *cobra.Command {
	var body iam.OAuthAuthorizeRequest
	var bodyFile string
	var stateFlag string
	cmd := &cobra.Command{
		Use:   "authorize",
		Short: "Approve a CLI login and issue an authorization code",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("state") {
				body.State = &stateFlag
			}
			out, err := c.AuthorizeOAuthClient(cmd.Context(), &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.ClientID, "client-id", "", "The registered client being approved")
	_ = cmd.MarkFlagRequired("client-id")
	f.StringVar(&body.CodeChallenge, "code-challenge", "", "Base64url SHA-256 of the client's PKCE verifier, without padding")
	_ = cmd.MarkFlagRequired("code-challenge")
	f.StringVar(&body.CodeChallengeMethod, "code-challenge-method", "", "S256 only (one of: S256)")
	_ = cmd.MarkFlagRequired("code-challenge-method")
	f.StringVar((*string)(&body.Organization), "organization", "", "Organization")
	_ = cmd.MarkFlagRequired("organization")
	f.StringVar(&body.RedirectURI, "redirect-uri", "", "For the CLI this must be urn:ietf:wg:oauth:2.0:oob — the out-of-band pseudo-redirect, meaning the code is DISPLAYED rather than delivered anywhere")
	_ = cmd.MarkFlagRequired("redirect-uri")
	f.StringVar(&stateFlag, "state", "", "Opaque value echoed back on the redirect, unchanged")
	return cmd
}

// newIamPolicyCommand builds `basaltic iam policy`.
func newIamPolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "policy",
		Short:   "Policies",
		Aliases: []string{"policies"},
	}
	cmd.AddCommand(newIamPolicyListCommand(state))
	cmd.AddCommand(newIamPolicyGetCommand(state))
	cmd.AddCommand(newIamPolicyCreateCommand(state))
	cmd.AddCommand(newIamPolicyUpdateCommand(state))
	cmd.AddCommand(newIamPolicyDeleteCommand(state))
	cmd.AddCommand(newIamPolicyListRolesCommand(state))
	cmd.AddCommand(newIamPolicyListServiceAccountsCommand(state))
	return cmd
}

// newIamPolicyListCommand builds `basaltic iam policy list`.
func newIamPolicyListCommand(state *cli.State) *cobra.Command {
	var params iam.ListPoliciesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List policies",
		Args:  cobra.ExactArgs(0),
		Long:  "List policies.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListPoliciesAll(cmd.Context(), &params))
			}
			page, err := c.ListPolicies(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact returned CRN, combined with name using AND before pagination")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Exact resource name, combined with crn using AND before pagination")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newIamPolicyGetCommand builds `basaltic iam policy get`.
func newIamPolicyGetCommand(state *cli.State) *cobra.Command {
	var scope iam.ListPoliciesParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Get policy",
		Args:  cobra.ExactArgs(1),
		Long:  "Get policy.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetPolicyByReference(cmd.Context(), args[0], &scope)
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

// newIamPolicyCreateCommand builds `basaltic iam policy create`.
func newIamPolicyCreateCommand(state *cli.State) *cobra.Command {
	var body iam.PolicyCreateRequest
	var bodyFile string
	var descriptionFlag string
	var documentFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create policy",
		Args:  cobra.ExactArgs(0),
		Long:  "Create policy.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
			}
			if documentFlag != "" {
				if err := json.Unmarshal([]byte(documentFlag), &body.Document); err != nil {
					return fmt.Errorf("--document: %w", err)
				}
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreatePolicy(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&descriptionFlag, "description", "", "Description")
	f.StringVar(&documentFlag, "document", "", "Document (JSON)")
	_ = cmd.MarkFlagRequired("document")
	f.StringVar(&body.Name, "name", "", "Resource names must not start with the literal crn: prefix or be UUIDs (canonical, compact, braced, or urn:uuid: forms, in either case)")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newIamPolicyUpdateCommand builds `basaltic iam policy update`.
func newIamPolicyUpdateCommand(state *cli.State) *cobra.Command {
	var body iam.PolicyUpdateRequest
	var bodyFile string
	var descriptionFlag string
	var documentFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <policy-id>",
		Short: "Update policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
			}
			if documentFlag != "" {
				if err := json.Unmarshal([]byte(documentFlag), &body.Document); err != nil {
					return fmt.Errorf("--document: %w", err)
				}
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdatePolicy(cmd.Context(), args[0], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&descriptionFlag, "description", "", "Description")
	f.StringVar(&documentFlag, "document", "", "Document (JSON)")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newIamPolicyDeleteCommand builds `basaltic iam policy delete`.
func newIamPolicyDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <policy-id>",
		Short: "Delete policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.DeletePolicy(cmd.Context(), args[0]); err != nil {
				return err
			}
			state.Printer().Done("Deleted.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamPolicyListRolesCommand builds `basaltic iam policy list-roles`.
func newIamPolicyListRolesCommand(state *cli.State) *cobra.Command {
	var params iam.ListPolicyRolesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list-roles <policy-id>",
		Short: "List roles with policy",
		Args:  cobra.ExactArgs(1),
		Long:  "List roles with policy.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListPolicyRolesAll(cmd.Context(), args[0], &params))
			}
			page, err := c.ListPolicyRoles(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact returned CRN, combined with name using AND before pagination")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Exact resource name, combined with crn using AND before pagination")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newIamPolicyListServiceAccountsCommand builds `basaltic iam policy list-service-accounts`.
func newIamPolicyListServiceAccountsCommand(state *cli.State) *cobra.Command {
	var params iam.ListPolicyServiceAccountsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list-service-accounts <policy-id>",
		Short: "List service accounts with policy",
		Args:  cobra.ExactArgs(1),
		Long:  "List service accounts with policy.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListPolicyServiceAccountsAll(cmd.Context(), args[0], &params))
			}
			page, err := c.ListPolicyServiceAccounts(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact returned CRN, combined with name using AND before pagination")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Exact resource name, combined with crn using AND before pagination")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newIamRegionCommand builds `basaltic iam region`.
func newIamRegionCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "region",
		Short:   "Regions",
		Aliases: []string{"regions"},
	}
	cmd.AddCommand(newIamRegionListCommand(state))
	return cmd
}

// newIamRegionListCommand builds `basaltic iam region list`.
func newIamRegionListCommand(state *cli.State) *cobra.Command {
	var params iam.ListRegionsParams
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List regions",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
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

// newIamRoleCommand builds `basaltic iam role`.
func newIamRoleCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "role",
		Short:   "Roles",
		Aliases: []string{"roles"},
	}
	cmd.AddCommand(newIamRoleListCommand(state))
	cmd.AddCommand(newIamRoleGetCommand(state))
	cmd.AddCommand(newIamRoleCreateCommand(state))
	cmd.AddCommand(newIamRoleUpdateCommand(state))
	cmd.AddCommand(newIamRoleDeleteCommand(state))
	cmd.AddCommand(newIamRoleAssumeCommand(state))
	cmd.AddCommand(newIamRoleAssumeWithWebIdentityCommand(state))
	cmd.AddCommand(newIamRoleAttachPolicyCommand(state))
	cmd.AddCommand(newIamRoleDeleteInlinePolicyCommand(state))
	cmd.AddCommand(newIamRoleDetachPolicyCommand(state))
	cmd.AddCommand(newIamRoleGetInlinePolicyCommand(state))
	cmd.AddCommand(newIamRoleGetPermissionBoundaryCommand(state))
	cmd.AddCommand(newIamRoleListInlinePoliciesCommand(state))
	cmd.AddCommand(newIamRoleListPoliciesCommand(state))
	cmd.AddCommand(newIamRoleRemovePermissionBoundaryCommand(state))
	cmd.AddCommand(newIamRoleSetInlinePolicyCommand(state))
	cmd.AddCommand(newIamRoleSetPermissionBoundaryCommand(state))
	return cmd
}

// newIamRoleListCommand builds `basaltic iam role list`.
func newIamRoleListCommand(state *cli.State) *cobra.Command {
	var params iam.ListRolesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List roles",
		Args:  cobra.ExactArgs(0),
		Long:  "List roles.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListRolesAll(cmd.Context(), &params))
			}
			page, err := c.ListRoles(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact returned CRN, combined with name using AND before pagination")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Exact resource name, combined with crn using AND before pagination")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newIamRoleGetCommand builds `basaltic iam role get`.
func newIamRoleGetCommand(state *cli.State) *cobra.Command {
	var scope iam.ListRolesParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Get role",
		Args:  cobra.ExactArgs(1),
		Long:  "Get role.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetRoleByReference(cmd.Context(), args[0], &scope)
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

// newIamRoleCreateCommand builds `basaltic iam role create`.
func newIamRoleCreateCommand(state *cli.State) *cobra.Command {
	var body iam.RoleCreateRequest
	var bodyFile string
	var descriptionFlag string
	var tagsFlag string
	var trustPolicyFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create role",
		Args:  cobra.ExactArgs(0),
		Long:  "Create role.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			if trustPolicyFlag != "" {
				if err := json.Unmarshal([]byte(trustPolicyFlag), &body.TrustPolicy); err != nil {
					return fmt.Errorf("--trust-policy: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateRole(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&descriptionFlag, "description", "", "Description")
	f.StringVar(&body.Name, "name", "", "Resource names must not start with the literal crn: prefix or be UUIDs (canonical, compact, braced, or urn:uuid: forms, in either case)")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&trustPolicyFlag, "trust-policy", "", "Trust policy (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newIamRoleUpdateCommand builds `basaltic iam role update`.
func newIamRoleUpdateCommand(state *cli.State) *cobra.Command {
	var body iam.RoleUpdateRequest
	var bodyFile string
	var descriptionFlag string
	var tagsFlag string
	var trustPolicyFlag string
	cmd := &cobra.Command{
		Use:   "update <role-id>",
		Short: "Update role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			if trustPolicyFlag != "" {
				if err := json.Unmarshal([]byte(trustPolicyFlag), &body.TrustPolicy); err != nil {
					return fmt.Errorf("--trust-policy: %w", err)
				}
			}
			out, err := c.UpdateRole(cmd.Context(), args[0], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&descriptionFlag, "description", "", "Description")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&trustPolicyFlag, "trust-policy", "", "Trust policy (JSON)")
	return cmd
}

// newIamRoleDeleteCommand builds `basaltic iam role delete`.
func newIamRoleDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <role-id>",
		Short: "Delete role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteRole(cmd.Context(), args[0]); err != nil {
				return err
			}
			state.Printer().Done("Deleted.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamRoleAssumeCommand builds `basaltic iam role assume`.
func newIamRoleAssumeCommand(state *cli.State) *cobra.Command {
	var body iam.AssumeRoleRequest
	var bodyFile string
	var durationSecondsFlag int
	var policyFlag string
	cmd := &cobra.Command{
		Use:   "assume",
		Short: "Assume role",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("duration-seconds") {
				body.DurationSeconds = &durationSecondsFlag
			}
			if policyFlag != "" {
				if err := json.Unmarshal([]byte(policyFlag), &body.Policy); err != nil {
					return fmt.Errorf("--policy: %w", err)
				}
			}
			out, err := c.AssumeRole(cmd.Context(), &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.IntVar(&durationSecondsFlag, "duration-seconds", 0, "Credential validity duration (15 min to 12 hours)")
	f.StringVar(&policyFlag, "policy", "", "Policy (JSON)")
	f.StringVar((*string)(&body.Role), "role", "", "Role")
	_ = cmd.MarkFlagRequired("role")
	return cmd
}

// newIamRoleAssumeWithWebIdentityCommand builds `basaltic iam role assume-with-web-identity`.
func newIamRoleAssumeWithWebIdentityCommand(state *cli.State) *cobra.Command {
	var body iam.AssumeRoleWithWebIdentityRequest
	var bodyFile string
	var durationSecondsFlag int
	var sessionNameFlag string
	cmd := &cobra.Command{
		Use:   "assume-with-web-identity",
		Short: "Assume role with web identity",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("duration-seconds") {
				body.DurationSeconds = &durationSecondsFlag
			}
			if cmd.Flags().Changed("session-name") {
				body.SessionName = &sessionNameFlag
			}
			out, err := c.AssumeRoleWithWebIdentity(cmd.Context(), &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar((*string)(&body.Account), "account", "", "Account")
	_ = cmd.MarkFlagRequired("account")
	f.IntVar(&durationSecondsFlag, "duration-seconds", 0, "Credential validity duration (15 min to 12 hours)")
	f.StringVar((*string)(&body.Role), "role", "", "Role")
	_ = cmd.MarkFlagRequired("role")
	f.StringVar(&sessionNameFlag, "session-name", "", "A label recorded on the session and in the audit trail")
	f.StringVar(&body.WebIdentityToken, "web-identity-token", "", "The identity token to exchange, as a signed JWT")
	_ = cmd.MarkFlagRequired("web-identity-token")
	return cmd
}

// newIamRoleAttachPolicyCommand builds `basaltic iam role attach-policy`.
func newIamRoleAttachPolicyCommand(state *cli.State) *cobra.Command {
	var body iam.RolePolicyAttachRequest
	var bodyFile string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "attach-policy <role-id>",
		Short: "Attach policy to role",
		Args:  cobra.ExactArgs(1),
		Long:  "Attach policy to role.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			if err := c.AttachRolePolicy(cmd.Context(), args[0], &body, reqOpts...); err != nil {
				return err
			}
			state.Printer().Done("Attach policy requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar((*string)(&body.Policy), "policy", "", "Policy")
	_ = cmd.MarkFlagRequired("policy")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newIamRoleDeleteInlinePolicyCommand builds `basaltic iam role delete-inline-policy`.
func newIamRoleDeleteInlinePolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-inline-policy <role-id> <policy-name>",
		Short: "Delete a role's inline policy by name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteRoleInlinePolicy(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			state.Printer().Done("Delete inline policy requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamRoleDetachPolicyCommand builds `basaltic iam role detach-policy`.
func newIamRoleDetachPolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach-policy <role-id> <policy-id>",
		Short: "Detach policy from role",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.DetachRolePolicy(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			state.Printer().Done("Detach policy requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamRoleGetInlinePolicyCommand builds `basaltic iam role get-inline-policy`.
func newIamRoleGetInlinePolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-inline-policy <role-id> <policy-name>",
		Short: "Get a role's inline policy by name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetRoleInlinePolicy(cmd.Context(), args[0], args[1])
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

// newIamRoleGetPermissionBoundaryCommand builds `basaltic iam role get-permission-boundary`.
func newIamRoleGetPermissionBoundaryCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-permission-boundary <role-id>",
		Short: "Get a role's permission boundary",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetRolePermissionBoundary(cmd.Context(), args[0])
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

// newIamRoleListInlinePoliciesCommand builds `basaltic iam role list-inline-policies`.
func newIamRoleListInlinePoliciesCommand(state *cli.State) *cobra.Command {
	var params iam.ListRoleInlinePoliciesParams
	cmd := &cobra.Command{
		Use:   "list-inline-policies <role-id>",
		Short: "List a role's inline policies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListRoleInlinePolicies(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact returned CRN, combined with name using AND before pagination")
	f.StringVar(&params.Name, "name", "", "Exact resource name, combined with crn using AND before pagination")
	return cmd
}

// newIamRoleListPoliciesCommand builds `basaltic iam role list-policies`.
func newIamRoleListPoliciesCommand(state *cli.State) *cobra.Command {
	var params iam.ListRolePoliciesParams
	cmd := &cobra.Command{
		Use:   "list-policies <role-id>",
		Short: "List role policies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListRolePolicies(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact returned CRN, combined with name using AND before pagination")
	f.StringVar(&params.Name, "name", "", "Exact resource name, combined with crn using AND before pagination")
	return cmd
}

// newIamRoleRemovePermissionBoundaryCommand builds `basaltic iam role remove-permission-boundary`.
func newIamRoleRemovePermissionBoundaryCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-permission-boundary <role-id>",
		Short: "Remove a role's permission boundary",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.RemoveRolePermissionBoundary(cmd.Context(), args[0]); err != nil {
				return err
			}
			state.Printer().Done("Remove permission boundary requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamRoleSetInlinePolicyCommand builds `basaltic iam role set-inline-policy`.
func newIamRoleSetInlinePolicyCommand(state *cli.State) *cobra.Command {
	var body iam.PutInlinePolicyRequest
	var bodyFile string
	var documentFlag string
	cmd := &cobra.Command{
		Use:   "set-inline-policy <role-id> <policy-name>",
		Short: "Create or replace a role's inline policy",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if documentFlag != "" {
				if err := json.Unmarshal([]byte(documentFlag), &body.Document); err != nil {
					return fmt.Errorf("--document: %w", err)
				}
			}
			out, err := c.PutRoleInlinePolicy(cmd.Context(), args[0], args[1], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&documentFlag, "document", "", "Document (JSON)")
	_ = cmd.MarkFlagRequired("document")
	return cmd
}

// newIamRoleSetPermissionBoundaryCommand builds `basaltic iam role set-permission-boundary`.
func newIamRoleSetPermissionBoundaryCommand(state *cli.State) *cobra.Command {
	var body iam.SetBoundaryRequest
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "set-permission-boundary <role-id>",
		Short: "Set a role's permission boundary",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if err := c.SetRolePermissionBoundary(cmd.Context(), args[0], &body); err != nil {
				return err
			}
			state.Printer().Done("Set permission boundary requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar((*string)(&body.Policy), "policy", "", "Policy")
	_ = cmd.MarkFlagRequired("policy")
	return cmd
}

// newIamServiceAccountCommand builds `basaltic iam service-account`.
func newIamServiceAccountCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "service-account",
		Short:   "Service accounts",
		Aliases: []string{"service-accounts"},
	}
	cmd.AddCommand(newIamServiceAccountListCommand(state))
	cmd.AddCommand(newIamServiceAccountGetCommand(state))
	cmd.AddCommand(newIamServiceAccountCreateCommand(state))
	cmd.AddCommand(newIamServiceAccountUpdateCommand(state))
	cmd.AddCommand(newIamServiceAccountDeleteCommand(state))
	cmd.AddCommand(newIamServiceAccountAttachPolicyCommand(state))
	cmd.AddCommand(newIamServiceAccountCreateCredentialCommand(state))
	cmd.AddCommand(newIamServiceAccountDeleteCredentialCommand(state))
	cmd.AddCommand(newIamServiceAccountDeleteInlinePolicyCommand(state))
	cmd.AddCommand(newIamServiceAccountDetachPolicyCommand(state))
	cmd.AddCommand(newIamServiceAccountGetInlinePolicyCommand(state))
	cmd.AddCommand(newIamServiceAccountGetPermissionBoundaryCommand(state))
	cmd.AddCommand(newIamServiceAccountListCredentialsCommand(state))
	cmd.AddCommand(newIamServiceAccountListInlinePoliciesCommand(state))
	cmd.AddCommand(newIamServiceAccountListPoliciesCommand(state))
	cmd.AddCommand(newIamServiceAccountRemovePermissionBoundaryCommand(state))
	cmd.AddCommand(newIamServiceAccountSetInlinePolicyCommand(state))
	cmd.AddCommand(newIamServiceAccountSetPermissionBoundaryCommand(state))
	return cmd
}

// newIamServiceAccountListCommand builds `basaltic iam service-account list`.
func newIamServiceAccountListCommand(state *cli.State) *cobra.Command {
	var params iam.ListServiceAccountsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List service accounts",
		Args:  cobra.ExactArgs(0),
		Long:  "List service accounts.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListServiceAccountsAll(cmd.Context(), &params))
			}
			page, err := c.ListServiceAccounts(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact returned CRN, combined with name using AND before pagination")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Exact resource name, combined with crn using AND before pagination")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newIamServiceAccountGetCommand builds `basaltic iam service-account get`.
func newIamServiceAccountGetCommand(state *cli.State) *cobra.Command {
	var scope iam.ListServiceAccountsParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Get service account",
		Args:  cobra.ExactArgs(1),
		Long:  "Get service account.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetServiceAccountByReference(cmd.Context(), args[0], &scope)
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

// newIamServiceAccountCreateCommand builds `basaltic iam service-account create`.
func newIamServiceAccountCreateCommand(state *cli.State) *cobra.Command {
	var body iam.ServiceAccountCreateRequest
	var bodyFile string
	var descriptionFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create service account",
		Args:  cobra.ExactArgs(0),
		Long:  "Create service account.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateServiceAccount(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&descriptionFlag, "description", "", "Description")
	f.StringVar(&body.Name, "name", "", "Resource names must not start with the literal crn: prefix or be UUIDs (canonical, compact, braced, or urn:uuid: forms, in either case)")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newIamServiceAccountUpdateCommand builds `basaltic iam service-account update`.
func newIamServiceAccountUpdateCommand(state *cli.State) *cobra.Command {
	var body iam.ServiceAccountUpdateRequest
	var bodyFile string
	var descriptionFlag string
	var enabledFlag bool
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <service-account-id>",
		Short: "Update service account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("description") {
				body.Description = &descriptionFlag
			}
			if cmd.Flags().Changed("enabled") {
				body.Enabled = &enabledFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			out, err := c.UpdateServiceAccount(cmd.Context(), args[0], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&descriptionFlag, "description", "", "Description")
	f.BoolVar(&enabledFlag, "enabled", false, "Enabled")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newIamServiceAccountDeleteCommand builds `basaltic iam service-account delete`.
func newIamServiceAccountDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <service-account-id>",
		Short: "Delete service account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteServiceAccount(cmd.Context(), args[0]); err != nil {
				return err
			}
			state.Printer().Done("Deleted.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamServiceAccountAttachPolicyCommand builds `basaltic iam service-account attach-policy`.
func newIamServiceAccountAttachPolicyCommand(state *cli.State) *cobra.Command {
	var body iam.PolicyAttachRequest
	var bodyFile string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "attach-policy <service-account-id>",
		Short: "Attach policy to service account",
		Args:  cobra.ExactArgs(1),
		Long:  "Attach policy to service account.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			if err := c.AttachServiceAccountPolicy(cmd.Context(), args[0], &body, reqOpts...); err != nil {
				return err
			}
			state.Printer().Done("Attach policy requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar((*string)(&body.Policy), "policy", "", "Policy")
	_ = cmd.MarkFlagRequired("policy")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newIamServiceAccountCreateCredentialCommand builds `basaltic iam service-account create-credential`.
func newIamServiceAccountCreateCredentialCommand(state *cli.State) *cobra.Command {
	var body iam.CredentialCreateRequest
	var bodyFile string
	var expiresAtFlag string
	cmd := &cobra.Command{
		Use:   "create-credential <service-account-id>",
		Short: "Create credential",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if expiresAtFlag != "" {
				parsed, err := parseTime(expiresAtFlag)
				if err != nil {
					return fmt.Errorf("--expires-at: %w", err)
				}
				body.ExpiresAt = &parsed
			}
			out, err := c.CreateServiceAccountCredential(cmd.Context(), args[0], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&expiresAtFlag, "expires-at", "", "Optional expiration date (RFC 3339)")
	f.StringVar(&body.Name, "name", "", "Name")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

// newIamServiceAccountDeleteCredentialCommand builds `basaltic iam service-account delete-credential`.
func newIamServiceAccountDeleteCredentialCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-credential <service-account-id> <credential-id>",
		Short: "Delete credential",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteServiceAccountCredential(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			state.Printer().Done("Delete credential requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamServiceAccountDeleteInlinePolicyCommand builds `basaltic iam service-account delete-inline-policy`.
func newIamServiceAccountDeleteInlinePolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-inline-policy <service-account-id> <policy-name>",
		Short: "Delete a service account's inline policy by name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteServiceAccountInlinePolicy(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			state.Printer().Done("Delete inline policy requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamServiceAccountDetachPolicyCommand builds `basaltic iam service-account detach-policy`.
func newIamServiceAccountDetachPolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach-policy <service-account-id> <policy-id>",
		Short: "Detach policy from service account",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.DetachServiceAccountPolicy(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			state.Printer().Done("Detach policy requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamServiceAccountGetInlinePolicyCommand builds `basaltic iam service-account get-inline-policy`.
func newIamServiceAccountGetInlinePolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-inline-policy <service-account-id> <policy-name>",
		Short: "Get a service account's inline policy by name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetServiceAccountInlinePolicy(cmd.Context(), args[0], args[1])
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

// newIamServiceAccountGetPermissionBoundaryCommand builds `basaltic iam service-account get-permission-boundary`.
func newIamServiceAccountGetPermissionBoundaryCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-permission-boundary <service-account-id>",
		Short: "Get a service account's permission boundary",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetServiceAccountPermissionBoundary(cmd.Context(), args[0])
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

// newIamServiceAccountListCredentialsCommand builds `basaltic iam service-account list-credentials`.
func newIamServiceAccountListCredentialsCommand(state *cli.State) *cobra.Command {
	var params iam.ListServiceAccountCredentialsParams
	cmd := &cobra.Command{
		Use:   "list-credentials <service-account-id>",
		Short: "List credentials",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListServiceAccountCredentials(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact returned CRN, combined with name using AND before pagination")
	f.StringVar(&params.Name, "name", "", "Exact resource name, combined with crn using AND before pagination")
	return cmd
}

// newIamServiceAccountListInlinePoliciesCommand builds `basaltic iam service-account list-inline-policies`.
func newIamServiceAccountListInlinePoliciesCommand(state *cli.State) *cobra.Command {
	var params iam.ListServiceAccountInlinePoliciesParams
	cmd := &cobra.Command{
		Use:   "list-inline-policies <service-account-id>",
		Short: "List a service account's inline policies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListServiceAccountInlinePolicies(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact returned CRN, combined with name using AND before pagination")
	f.StringVar(&params.Name, "name", "", "Exact resource name, combined with crn using AND before pagination")
	return cmd
}

// newIamServiceAccountListPoliciesCommand builds `basaltic iam service-account list-policies`.
func newIamServiceAccountListPoliciesCommand(state *cli.State) *cobra.Command {
	var params iam.ListServiceAccountPoliciesParams
	cmd := &cobra.Command{
		Use:   "list-policies <service-account-id>",
		Short: "List service account policies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListServiceAccountPolicies(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact returned CRN, combined with name using AND before pagination")
	f.StringVar(&params.Name, "name", "", "Exact resource name, combined with crn using AND before pagination")
	return cmd
}

// newIamServiceAccountRemovePermissionBoundaryCommand builds `basaltic iam service-account remove-permission-boundary`.
func newIamServiceAccountRemovePermissionBoundaryCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-permission-boundary <service-account-id>",
		Short: "Remove a service account's permission boundary",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.RemoveServiceAccountPermissionBoundary(cmd.Context(), args[0]); err != nil {
				return err
			}
			state.Printer().Done("Remove permission boundary requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamServiceAccountSetInlinePolicyCommand builds `basaltic iam service-account set-inline-policy`.
func newIamServiceAccountSetInlinePolicyCommand(state *cli.State) *cobra.Command {
	var body iam.PutInlinePolicyRequest
	var bodyFile string
	var documentFlag string
	cmd := &cobra.Command{
		Use:   "set-inline-policy <service-account-id> <policy-name>",
		Short: "Create or replace a service account's inline policy",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if documentFlag != "" {
				if err := json.Unmarshal([]byte(documentFlag), &body.Document); err != nil {
					return fmt.Errorf("--document: %w", err)
				}
			}
			out, err := c.PutServiceAccountInlinePolicy(cmd.Context(), args[0], args[1], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&documentFlag, "document", "", "Document (JSON)")
	_ = cmd.MarkFlagRequired("document")
	return cmd
}

// newIamServiceAccountSetPermissionBoundaryCommand builds `basaltic iam service-account set-permission-boundary`.
func newIamServiceAccountSetPermissionBoundaryCommand(state *cli.State) *cobra.Command {
	var body iam.SetBoundaryRequest
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "set-permission-boundary <service-account-id>",
		Short: "Set a service account's permission boundary",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if err := c.SetServiceAccountPermissionBoundary(cmd.Context(), args[0], &body); err != nil {
				return err
			}
			state.Printer().Done("Set permission boundary requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar((*string)(&body.Policy), "policy", "", "Policy")
	_ = cmd.MarkFlagRequired("policy")
	return cmd
}

// newIamStsSessionCommand builds `basaltic iam sts-session`.
func newIamStsSessionCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sts-session",
		Short:   "Sts sessions",
		Aliases: []string{"sts-sessions"},
	}
	cmd.AddCommand(newIamStsSessionListCommand(state))
	cmd.AddCommand(newIamStsSessionGetCommand(state))
	cmd.AddCommand(newIamStsSessionRevokeCommand(state))
	return cmd
}

// newIamStsSessionListCommand builds `basaltic iam sts-session list`.
func newIamStsSessionListCommand(state *cli.State) *cobra.Command {
	var params iam.ListSTSSessionsParams
	var activeOnlyFlag bool
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List STS sessions",
		Args:  cobra.ExactArgs(0),
		Long:  "List STS sessions.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("active-only") {
				params.ActiveOnly = &activeOnlyFlag
			}
			if fetchAll {
				return state.Printer().Iter(c.ListSTSSessionsAll(cmd.Context(), &params))
			}
			page, err := c.ListSTSSessions(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.BoolVar(&activeOnlyFlag, "active-only", false, "Only show active (non-expired, non-revoked) sessions")
	f.StringVar(&params.CRN, "crn", "", "Exact returned CRN, combined with name using AND before pagination")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Exact resource name, combined with crn using AND before pagination")
	f.StringVar(&params.Principal, "principal", "", "Principal reference; principal_type is required")
	f.StringVar(&params.PrincipalType, "principal-type", "", "Filter by principal type")
	f.StringVar(&params.Role, "role", "", "Role UUID, immutable name, or account-qualified role CRN")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newIamStsSessionGetCommand builds `basaltic iam sts-session get`.
func newIamStsSessionGetCommand(state *cli.State) *cobra.Command {
	var scope iam.ListSTSSessionsParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Get STS session",
		Args:  cobra.ExactArgs(1),
		Long:  "Get STS session.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.\n\nA name is unique only within its parent: pass --principal or --role with a name, or the\nlookup can match more than one and is refused.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetSTSSessionByReference(cmd.Context(), args[0], &scope)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&scope.Principal, "principal", "", "Principal reference; principal_type is required")
	f.StringVar(&scope.Role, "role", "", "Role UUID, immutable name, or account-qualified role CRN")
	return cmd
}

// newIamStsSessionRevokeCommand builds `basaltic iam sts-session revoke`.
func newIamStsSessionRevokeCommand(state *cli.State) *cobra.Command {
	var body iam.RevokeSTSSessionRequest
	var bodyFile string
	var reasonFlag string
	cmd := &cobra.Command{
		Use:   "revoke <session-id>",
		Short: "Revoke STS session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("reason") {
				body.Reason = &reasonFlag
			}
			out, err := c.RevokeSTSSession(cmd.Context(), args[0], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&reasonFlag, "reason", "", "Reason for revoking the session")
	return cmd
}

// newIamTokenCommand builds `basaltic iam token`.
func newIamTokenCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token",
		Short:   "Tokens",
		Aliases: []string{"tokens"},
	}
	cmd.AddCommand(newIamTokenCreateCommand(state))
	cmd.AddCommand(newIamTokenRevokeCommand(state))
	return cmd
}

// newIamTokenCreateCommand builds `basaltic iam token create`.
func newIamTokenCreateCommand(state *cli.State) *cobra.Command {
	var body iam.OAuthTokenRequest
	var bodyFile string
	var clientIdFlag string
	var clientSecretFlag string
	var codeFlag string
	var codeVerifierFlag string
	var durationSecondsFlag int
	var redirectUriFlag string
	var refreshTokenFlag string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Exchange an access key for a bearer token",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("client-id") {
				body.ClientID = &clientIdFlag
			}
			if cmd.Flags().Changed("client-secret") {
				body.ClientSecret = &clientSecretFlag
			}
			if cmd.Flags().Changed("code") {
				body.Code = &codeFlag
			}
			if cmd.Flags().Changed("code-verifier") {
				body.CodeVerifier = &codeVerifierFlag
			}
			if cmd.Flags().Changed("duration-seconds") {
				body.DurationSeconds = &durationSecondsFlag
			}
			if cmd.Flags().Changed("redirect-uri") {
				body.RedirectURI = &redirectUriFlag
			}
			if cmd.Flags().Changed("refresh-token") {
				body.RefreshToken = &refreshTokenFlag
			}
			out, err := c.GetOAuthToken(cmd.Context(), &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&clientIdFlag, "client-id", "", "The access key id")
	f.StringVar(&clientSecretFlag, "client-secret", "", "The secret access key")
	f.StringVar(&codeFlag, "code", "", "The authorization code from the consent redirect")
	f.StringVar(&codeVerifierFlag, "code-verifier", "", "The PKCE verifier whose SHA-256 was sent as code_challenge when the flow started (RFC 7636)")
	f.IntVar(&durationSecondsFlag, "duration-seconds", 0, "Requested token lifetime")
	f.StringVar(&body.GrantType, "grant-type", "", "client_credentials is the one to use for a service account: it exchanges an access key pair for a token, and needs nothing else (one of: client_credentials, authorization_code, refresh_token)")
	_ = cmd.MarkFlagRequired("grant-type")
	f.StringVar(&redirectUriFlag, "redirect-uri", "", "The same redirect_uri the code was issued for — for the CLI, urn:ietf:wg:oauth:2.0:oob")
	f.StringVar(&refreshTokenFlag, "refresh-token", "", "refresh_token grant only")
	return cmd
}

// newIamTokenRevokeCommand builds `basaltic iam token revoke`.
func newIamTokenRevokeCommand(state *cli.State) *cobra.Command {
	var body iam.OAuthRevokeRequest
	var bodyFile string
	var tokenTypeHintFlag string
	cmd := &cobra.Command{
		Use:   "revoke",
		Short: "Revoke a bearer token",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("token-type-hint") {
				body.TokenTypeHint = &tokenTypeHintFlag
			}
			if err := c.RevokeOAuthToken(cmd.Context(), &body); err != nil {
				return err
			}
			state.Printer().Done("Revoked.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.Token, "token", "", "The access token to revoke")
	_ = cmd.MarkFlagRequired("token")
	f.StringVar(&tokenTypeHintFlag, "token-type-hint", "", "Accepted and ignored — the token identifies itself")
	return cmd
}
