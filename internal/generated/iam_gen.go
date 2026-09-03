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
		Short: "Identity, access, accounts and organizations",
	}
	cmd.AddCommand(newIamAccountCommand(state))
	cmd.AddCommand(newIamGroupCommand(state))
	cmd.AddCommand(newIamInvitationCommand(state))
	cmd.AddCommand(newIamOauthCommand(state))
	cmd.AddCommand(newIamOrganizationCommand(state))
	cmd.AddCommand(newIamPolicyCommand(state))
	cmd.AddCommand(newIamRegionCommand(state))
	cmd.AddCommand(newIamRoleCommand(state))
	cmd.AddCommand(newIamServiceAccountCommand(state))
	cmd.AddCommand(newIamStsSessionCommand(state))
	cmd.AddCommand(newIamTokenCommand(state))
	cmd.AddCommand(newIamUserCommand(state))
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

// newIamAccountCommand builds `basaltic iam account`.
func newIamAccountCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account",
		Short:   "Accounts",
		Aliases: []string{"accounts"},
	}
	cmd.AddCommand(newIamAccountListCommand(state))
	cmd.AddCommand(newIamAccountGetCommand(state))
	cmd.AddCommand(newIamAccountCreateCommand(state))
	cmd.AddCommand(newIamAccountUpdateCommand(state))
	cmd.AddCommand(newIamAccountDeleteCommand(state))
	return cmd
}

// newIamAccountListCommand builds `basaltic iam account list`.
func newIamAccountListCommand(state *cli.State) *cobra.Command {
	var params iam.ListAccountsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List accounts",
		Args:  cobra.ExactArgs(0),
		Long:  "List accounts.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListAccountsAll(cmd.Context(), &params))
			}
			page, err := c.ListAccounts(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newIamAccountGetCommand builds `basaltic iam account get`.
func newIamAccountGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <account-id>",
		Short: "Get account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetAccount(cmd.Context(), args[0])
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

// newIamAccountCreateCommand builds `basaltic iam account create`.
func newIamAccountCreateCommand(state *cli.State) *cobra.Command {
	var body iam.CreateAccountRequest
	var bodyFile string
	var descriptionFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create account",
		Args:  cobra.ExactArgs(0),
		Long:  "Create account.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
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
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateAccount(cmd.Context(), &body, reqOpts...)
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
	f.StringVar(&body.Handle, "handle", "", "Handle")
	_ = cmd.MarkFlagRequired("handle")
	f.StringVar(&body.Name, "name", "", "Name")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newIamAccountUpdateCommand builds `basaltic iam account update`.
func newIamAccountUpdateCommand(state *cli.State) *cobra.Command {
	var body iam.UpdateAccountRequest
	var bodyFile string
	var descriptionFlag string
	var nameFlag string
	cmd := &cobra.Command{
		Use:   "update <account-id>",
		Short: "Update account",
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
			if cmd.Flags().Changed("name") {
				body.Name = &nameFlag
			}
			out, err := c.UpdateAccount(cmd.Context(), args[0], &body)
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
	f.StringVar(&nameFlag, "name", "", "Name")
	return cmd
}

// newIamAccountDeleteCommand builds `basaltic iam account delete`.
func newIamAccountDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <account-id>",
		Short: "Delete account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteAccount(cmd.Context(), args[0]); err != nil {
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

// newIamGroupCommand builds `basaltic iam group`.
func newIamGroupCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "group",
		Short:   "Groups",
		Aliases: []string{"groups"},
	}
	cmd.AddCommand(newIamGroupListCommand(state))
	cmd.AddCommand(newIamGroupGetCommand(state))
	cmd.AddCommand(newIamGroupCreateCommand(state))
	cmd.AddCommand(newIamGroupUpdateCommand(state))
	cmd.AddCommand(newIamGroupDeleteCommand(state))
	cmd.AddCommand(newIamGroupAttachPolicyCommand(state))
	cmd.AddCommand(newIamGroupDeleteInlinePolicyCommand(state))
	cmd.AddCommand(newIamGroupDetachPolicyCommand(state))
	cmd.AddCommand(newIamGroupGetInlinePolicyCommand(state))
	cmd.AddCommand(newIamGroupListInlinePoliciesCommand(state))
	cmd.AddCommand(newIamGroupListPoliciesCommand(state))
	cmd.AddCommand(newIamGroupListServiceAccountsCommand(state))
	cmd.AddCommand(newIamGroupListUsersCommand(state))
	cmd.AddCommand(newIamGroupSetInlinePolicyCommand(state))
	return cmd
}

// newIamGroupListCommand builds `basaltic iam group list`.
func newIamGroupListCommand(state *cli.State) *cobra.Command {
	var params iam.ListGroupsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List groups",
		Args:  cobra.ExactArgs(0),
		Long:  "List groups.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListGroupsAll(cmd.Context(), &params))
			}
			page, err := c.ListGroups(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Filter by name (exact match or prefix with *)")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newIamGroupGetCommand builds `basaltic iam group get`.
func newIamGroupGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <group-id>",
		Short: "Get group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetGroup(cmd.Context(), args[0])
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

// newIamGroupCreateCommand builds `basaltic iam group create`.
func newIamGroupCreateCommand(state *cli.State) *cobra.Command {
	var body iam.GroupCreateRequest
	var bodyFile string
	var descriptionFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create group",
		Args:  cobra.ExactArgs(0),
		Long:  "Create group.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
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
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateGroup(cmd.Context(), &body, reqOpts...)
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
	f.StringVar(&body.Name, "name", "", "Name")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newIamGroupUpdateCommand builds `basaltic iam group update`.
func newIamGroupUpdateCommand(state *cli.State) *cobra.Command {
	var body iam.GroupUpdateRequest
	var bodyFile string
	var descriptionFlag string
	var nameFlag string
	cmd := &cobra.Command{
		Use:   "update <group-id>",
		Short: "Update group",
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
			if cmd.Flags().Changed("name") {
				body.Name = &nameFlag
			}
			out, err := c.UpdateGroup(cmd.Context(), args[0], &body)
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
	f.StringVar(&nameFlag, "name", "", "Name")
	return cmd
}

// newIamGroupDeleteCommand builds `basaltic iam group delete`.
func newIamGroupDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <group-id>",
		Short: "Delete group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteGroup(cmd.Context(), args[0]); err != nil {
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

// newIamGroupAttachPolicyCommand builds `basaltic iam group attach-policy`.
func newIamGroupAttachPolicyCommand(state *cli.State) *cobra.Command {
	var body iam.PolicyAttachRequest
	var bodyFile string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "attach-policy <group-id>",
		Short: "Attach policy to group",
		Args:  cobra.ExactArgs(1),
		Long:  "Attach policy to group.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
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
			if err := c.AttachGroupPolicy(cmd.Context(), args[0], &body, reqOpts...); err != nil {
				return err
			}
			state.Printer().Done("Attach policy requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.PolicyID, "policy-id", "", "Policy id")
	_ = cmd.MarkFlagRequired("policy-id")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newIamGroupDeleteInlinePolicyCommand builds `basaltic iam group delete-inline-policy`.
func newIamGroupDeleteInlinePolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-inline-policy <group-id> <policy-name>",
		Short: "Delete a group's inline policy by name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteGroupInlinePolicy(cmd.Context(), args[0], args[1]); err != nil {
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

// newIamGroupDetachPolicyCommand builds `basaltic iam group detach-policy`.
func newIamGroupDetachPolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach-policy <group-id> <policy-id>",
		Short: "Detach policy from group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.DetachGroupPolicy(cmd.Context(), args[0], args[1]); err != nil {
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

// newIamGroupGetInlinePolicyCommand builds `basaltic iam group get-inline-policy`.
func newIamGroupGetInlinePolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-inline-policy <group-id> <policy-name>",
		Short: "Get a group's inline policy by name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetGroupInlinePolicy(cmd.Context(), args[0], args[1])
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

// newIamGroupListInlinePoliciesCommand builds `basaltic iam group list-inline-policies`.
func newIamGroupListInlinePoliciesCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-inline-policies <group-id>",
		Short: "List a group's inline policies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListGroupInlinePolicies(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamGroupListPoliciesCommand builds `basaltic iam group list-policies`.
func newIamGroupListPoliciesCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-policies <group-id>",
		Short: "List group policies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListGroupPolicies(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamGroupListServiceAccountsCommand builds `basaltic iam group list-service-accounts`.
func newIamGroupListServiceAccountsCommand(state *cli.State) *cobra.Command {
	var params iam.ListGroupServiceAccountsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list-service-accounts <group-id>",
		Short: "List group service accounts",
		Args:  cobra.ExactArgs(1),
		Long:  "List group service accounts.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListGroupServiceAccountsAll(cmd.Context(), args[0], &params))
			}
			page, err := c.ListGroupServiceAccounts(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newIamGroupListUsersCommand builds `basaltic iam group list-users`.
func newIamGroupListUsersCommand(state *cli.State) *cobra.Command {
	var params iam.ListGroupUsersParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list-users <group-id>",
		Short: "List group users",
		Args:  cobra.ExactArgs(1),
		Long:  "List group users.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListGroupUsersAll(cmd.Context(), args[0], &params))
			}
			page, err := c.ListGroupUsers(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newIamGroupSetInlinePolicyCommand builds `basaltic iam group set-inline-policy`.
func newIamGroupSetInlinePolicyCommand(state *cli.State) *cobra.Command {
	var body iam.PutInlinePolicyRequest
	var bodyFile string
	var documentFlag string
	cmd := &cobra.Command{
		Use:   "set-inline-policy <group-id> <policy-name>",
		Short: "Create or replace a group's inline policy",
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
			out, err := c.PutGroupInlinePolicy(cmd.Context(), args[0], args[1], &body)
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

// newIamInvitationCommand builds `basaltic iam invitation`.
func newIamInvitationCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "invitation",
		Short:   "Invitations",
		Aliases: []string{"invitations"},
	}
	cmd.AddCommand(newIamInvitationListCommand(state))
	cmd.AddCommand(newIamInvitationGetCommand(state))
	cmd.AddCommand(newIamInvitationCancelCommand(state))
	return cmd
}

// newIamInvitationListCommand builds `basaltic iam invitation list`.
func newIamInvitationListCommand(state *cli.State) *cobra.Command {
	var params iam.ListInvitationsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List invitations",
		Args:  cobra.ExactArgs(0),
		Long:  "List invitations.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListInvitationsAll(cmd.Context(), &params))
			}
			page, err := c.ListInvitations(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newIamInvitationGetCommand builds `basaltic iam invitation get`.
func newIamInvitationGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <invitation-id>",
		Short: "Get invitation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetInvitation(cmd.Context(), args[0])
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

// newIamInvitationCancelCommand builds `basaltic iam invitation cancel`.
func newIamInvitationCancelCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel <invitation-id>",
		Short: "Cancel invitation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.CancelInvitation(cmd.Context(), args[0]); err != nil {
				return err
			}
			state.Printer().Done("Cancel requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
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
	f.StringVar(&body.OrganizationID, "organization-id", "", "Which organization the resulting session is scoped to")
	_ = cmd.MarkFlagRequired("organization-id")
	f.StringVar(&body.RedirectURI, "redirect-uri", "", "Where to deliver the code")
	_ = cmd.MarkFlagRequired("redirect-uri")
	f.StringVar(&stateFlag, "state", "", "Opaque value echoed back on the redirect, unchanged")
	return cmd
}

// newIamOrganizationCommand builds `basaltic iam organization`.
func newIamOrganizationCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "organization",
		Short:   "Organizations",
		Aliases: []string{"organizations"},
	}
	cmd.AddCommand(newIamOrganizationListCommand(state))
	cmd.AddCommand(newIamOrganizationGetCommand(state))
	cmd.AddCommand(newIamOrganizationUpdateCommand(state))
	cmd.AddCommand(newIamOrganizationDeleteCommand(state))
	return cmd
}

// newIamOrganizationListCommand builds `basaltic iam organization list`.
func newIamOrganizationListCommand(state *cli.State) *cobra.Command {
	var params iam.ListOrganizationsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List organizations",
		Args:  cobra.ExactArgs(0),
		Long:  "List organizations.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListOrganizationsAll(cmd.Context(), &params))
			}
			page, err := c.ListOrganizations(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newIamOrganizationGetCommand builds `basaltic iam organization get`.
func newIamOrganizationGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <organization-id>",
		Short: "Get organization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetOrganization(cmd.Context(), args[0])
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

// newIamOrganizationUpdateCommand builds `basaltic iam organization update`.
func newIamOrganizationUpdateCommand(state *cli.State) *cobra.Command {
	var body iam.OrganizationUpdateRequest
	var bodyFile string
	var descriptionFlag string
	var nameFlag string
	cmd := &cobra.Command{
		Use:   "update <organization-id>",
		Short: "Update organization",
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
			if cmd.Flags().Changed("name") {
				body.Name = &nameFlag
			}
			out, err := c.UpdateOrganization(cmd.Context(), args[0], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.CaptchaToken, "captcha-token", "", "Google reCAPTCHA token for bot protection")
	_ = cmd.MarkFlagRequired("captcha-token")
	f.StringVar(&descriptionFlag, "description", "", "Description")
	f.StringVar(&nameFlag, "name", "", "Name")
	return cmd
}

// newIamOrganizationDeleteCommand builds `basaltic iam organization delete`.
func newIamOrganizationDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <organization-id>",
		Short: "Delete organization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteOrganization(cmd.Context(), args[0]); err != nil {
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
	cmd.AddCommand(newIamPolicyListGroupsCommand(state))
	cmd.AddCommand(newIamPolicyListRolesCommand(state))
	cmd.AddCommand(newIamPolicyListServiceAccountsCommand(state))
	cmd.AddCommand(newIamPolicyListUsersCommand(state))
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
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Filter by name (exact match or prefix with *)")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newIamPolicyGetCommand builds `basaltic iam policy get`.
func newIamPolicyGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <policy-id>",
		Short: "Get policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetPolicy(cmd.Context(), args[0])
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
	f.StringVar(&body.Name, "name", "", "Name")
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
	var nameFlag string
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
			if cmd.Flags().Changed("name") {
				body.Name = &nameFlag
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
	f.StringVar(&nameFlag, "name", "", "Name")
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

// newIamPolicyListGroupsCommand builds `basaltic iam policy list-groups`.
func newIamPolicyListGroupsCommand(state *cli.State) *cobra.Command {
	var params iam.ListPolicyGroupsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list-groups <policy-id>",
		Short: "List groups with policy",
		Args:  cobra.ExactArgs(1),
		Long:  "List groups with policy.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListPolicyGroupsAll(cmd.Context(), args[0], &params))
			}
			page, err := c.ListPolicyGroups(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
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
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
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
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newIamPolicyListUsersCommand builds `basaltic iam policy list-users`.
func newIamPolicyListUsersCommand(state *cli.State) *cobra.Command {
	var params iam.ListPolicyUsersParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list-users <policy-id>",
		Short: "List users with policy",
		Args:  cobra.ExactArgs(1),
		Long:  "List users with policy.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListPolicyUsersAll(cmd.Context(), args[0], &params))
			}
			page, err := c.ListPolicyUsers(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
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
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List regions",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			out, err := c.ListRegions(cmd.Context())
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
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Filter by name (exact match or prefix with *)")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newIamRoleGetCommand builds `basaltic iam role get`.
func newIamRoleGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <role-id>",
		Short: "Get role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetRole(cmd.Context(), args[0])
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
	f.StringVar(&body.Name, "name", "", "Name")
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
	var nameFlag string
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
			if cmd.Flags().Changed("name") {
				body.Name = &nameFlag
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
	f.StringVar(&nameFlag, "name", "", "Name")
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
	f.StringVar(&body.RoleID, "role-id", "", "Role id")
	_ = cmd.MarkFlagRequired("role-id")
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
	f.StringVar(&body.AccountID, "target-account-id", "", "The account the resulting credentials act in — the ownership scope stamped on the session, the same scope a signed request selects with X-Account-Id")
	_ = cmd.MarkFlagRequired("target-account-id")
	f.IntVar(&durationSecondsFlag, "duration-seconds", 0, "Credential validity duration (15 min to 12 hours)")
	f.StringVar(&body.RoleID, "role-id", "", "The role to assume")
	_ = cmd.MarkFlagRequired("role-id")
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
	f.StringVar(&body.PolicyID, "policy-id", "", "Policy id")
	_ = cmd.MarkFlagRequired("policy-id")
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
	cmd := &cobra.Command{
		Use:   "list-inline-policies <role-id>",
		Short: "List a role's inline policies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListRoleInlinePolicies(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamRoleListPoliciesCommand builds `basaltic iam role list-policies`.
func newIamRoleListPoliciesCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-policies <role-id>",
		Short: "List role policies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListRolePolicies(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
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
	f.StringVar(&body.PolicyID, "policy-id", "", "Policy id")
	_ = cmd.MarkFlagRequired("policy-id")
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
	cmd.AddCommand(newIamServiceAccountAddGroupCommand(state))
	cmd.AddCommand(newIamServiceAccountAttachPolicyCommand(state))
	cmd.AddCommand(newIamServiceAccountCreateCredentialCommand(state))
	cmd.AddCommand(newIamServiceAccountDeleteCredentialCommand(state))
	cmd.AddCommand(newIamServiceAccountDeleteInlinePolicyCommand(state))
	cmd.AddCommand(newIamServiceAccountDetachPolicyCommand(state))
	cmd.AddCommand(newIamServiceAccountGetInlinePolicyCommand(state))
	cmd.AddCommand(newIamServiceAccountGetPermissionBoundaryCommand(state))
	cmd.AddCommand(newIamServiceAccountListCredentialsCommand(state))
	cmd.AddCommand(newIamServiceAccountListGroupsCommand(state))
	cmd.AddCommand(newIamServiceAccountListInlinePoliciesCommand(state))
	cmd.AddCommand(newIamServiceAccountListPoliciesCommand(state))
	cmd.AddCommand(newIamServiceAccountRemoveGroupCommand(state))
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
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Filter by name (exact match or prefix with *)")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newIamServiceAccountGetCommand builds `basaltic iam service-account get`.
func newIamServiceAccountGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <service-account-id>",
		Short: "Get service account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetServiceAccount(cmd.Context(), args[0])
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
	f.StringVar(&body.Name, "name", "", "Name")
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

// newIamServiceAccountAddGroupCommand builds `basaltic iam service-account add-group`.
func newIamServiceAccountAddGroupCommand(state *cli.State) *cobra.Command {
	var body iam.ServiceAccountGroupAddRequest
	var bodyFile string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "add-group <service-account-id>",
		Short: "Add service account to group",
		Args:  cobra.ExactArgs(1),
		Long:  "Add service account to group.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
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
			if err := c.AddServiceAccountToGroup(cmd.Context(), args[0], &body, reqOpts...); err != nil {
				return err
			}
			state.Printer().Done("Add group requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.GroupID, "group-id", "", "Group id")
	_ = cmd.MarkFlagRequired("group-id")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
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
	f.StringVar(&body.PolicyID, "policy-id", "", "Policy id")
	_ = cmd.MarkFlagRequired("policy-id")
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
	cmd := &cobra.Command{
		Use:   "list-credentials <service-account-id>",
		Short: "List credentials",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListServiceAccountCredentials(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamServiceAccountListGroupsCommand builds `basaltic iam service-account list-groups`.
func newIamServiceAccountListGroupsCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-groups <service-account-id>",
		Short: "List service account groups",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListServiceAccountGroups(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamServiceAccountListInlinePoliciesCommand builds `basaltic iam service-account list-inline-policies`.
func newIamServiceAccountListInlinePoliciesCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-inline-policies <service-account-id>",
		Short: "List a service account's inline policies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListServiceAccountInlinePolicies(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamServiceAccountListPoliciesCommand builds `basaltic iam service-account list-policies`.
func newIamServiceAccountListPoliciesCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-policies <service-account-id>",
		Short: "List service account policies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListServiceAccountPolicies(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamServiceAccountRemoveGroupCommand builds `basaltic iam service-account remove-group`.
func newIamServiceAccountRemoveGroupCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-group <service-account-id> <group-id>",
		Short: "Remove service account from group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.RemoveServiceAccountFromGroup(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			state.Printer().Done("Remove group requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
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
	f.StringVar(&body.PolicyID, "policy-id", "", "Policy id")
	_ = cmd.MarkFlagRequired("policy-id")
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
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.PrincipalID, "principal-id", "", "Filter by principal ID (user or service account)")
	f.StringVar(&params.PrincipalType, "principal-type", "", "Filter by principal type")
	f.StringVar(&params.RoleID, "role-id", "", "Filter by role ID")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newIamStsSessionGetCommand builds `basaltic iam sts-session get`.
func newIamStsSessionGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <session-id>",
		Short: "Get STS session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetSTSSession(cmd.Context(), args[0])
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
	f.StringVar(&redirectUriFlag, "redirect-uri", "", "The same redirect_uri the code was issued for")
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

// newIamUserCommand builds `basaltic iam user`.
func newIamUserCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "user",
		Short:   "Users",
		Aliases: []string{"users"},
	}
	cmd.AddCommand(newIamUserListCommand(state))
	cmd.AddCommand(newIamUserGetCommand(state))
	cmd.AddCommand(newIamUserAddCommand(state))
	cmd.AddCommand(newIamUserAddGroupCommand(state))
	cmd.AddCommand(newIamUserAttachPolicyCommand(state))
	cmd.AddCommand(newIamUserDeleteInlinePolicyCommand(state))
	cmd.AddCommand(newIamUserDetachPolicyCommand(state))
	cmd.AddCommand(newIamUserGetInlinePolicyCommand(state))
	cmd.AddCommand(newIamUserGetPermissionBoundaryCommand(state))
	cmd.AddCommand(newIamUserListGroupsCommand(state))
	cmd.AddCommand(newIamUserListInlinePoliciesCommand(state))
	cmd.AddCommand(newIamUserListPoliciesCommand(state))
	cmd.AddCommand(newIamUserRemoveCommand(state))
	cmd.AddCommand(newIamUserRemoveGroupCommand(state))
	cmd.AddCommand(newIamUserRemovePermissionBoundaryCommand(state))
	cmd.AddCommand(newIamUserSetInlinePolicyCommand(state))
	cmd.AddCommand(newIamUserSetPermissionBoundaryCommand(state))
	return cmd
}

// newIamUserListCommand builds `basaltic iam user list`.
func newIamUserListCommand(state *cli.State) *cobra.Command {
	var params iam.ListUsersParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List users",
		Args:  cobra.ExactArgs(0),
		Long:  "List users.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListUsersAll(cmd.Context(), &params))
			}
			page, err := c.ListUsers(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Filter by name (exact match or prefix with *)")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newIamUserGetCommand builds `basaltic iam user get`.
func newIamUserGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <user-id>",
		Short: "Get user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetUser(cmd.Context(), args[0])
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

// newIamUserAddCommand builds `basaltic iam user add`.
func newIamUserAddCommand(state *cli.State) *cobra.Command {
	var body iam.UserAddRequest
	var bodyFile string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add user to organization",
		Args:  cobra.ExactArgs(0),
		Long:  "Add user to organization.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
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
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.AddUser(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.Email, "email", "", "Email of the user to add")
	_ = cmd.MarkFlagRequired("email")
	f.StringSliceVar(&body.GroupIDs, "group-ids", nil, "IDs of groups to add the user to")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newIamUserAddGroupCommand builds `basaltic iam user add-group`.
func newIamUserAddGroupCommand(state *cli.State) *cobra.Command {
	var body iam.UserGroupAddRequest
	var bodyFile string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "add-group <user-id>",
		Short: "Add user to group",
		Args:  cobra.ExactArgs(1),
		Long:  "Add user to group.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
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
			if err := c.AddUserToGroup(cmd.Context(), args[0], &body, reqOpts...); err != nil {
				return err
			}
			state.Printer().Done("Add group requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.GroupID, "group-id", "", "Group id")
	_ = cmd.MarkFlagRequired("group-id")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newIamUserAttachPolicyCommand builds `basaltic iam user attach-policy`.
func newIamUserAttachPolicyCommand(state *cli.State) *cobra.Command {
	var body iam.PolicyAttachRequest
	var bodyFile string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "attach-policy <user-id>",
		Short: "Attach policy to user",
		Args:  cobra.ExactArgs(1),
		Long:  "Attach policy to user.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
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
			if err := c.AttachUserPolicy(cmd.Context(), args[0], &body, reqOpts...); err != nil {
				return err
			}
			state.Printer().Done("Attach policy requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.PolicyID, "policy-id", "", "Policy id")
	_ = cmd.MarkFlagRequired("policy-id")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newIamUserDeleteInlinePolicyCommand builds `basaltic iam user delete-inline-policy`.
func newIamUserDeleteInlinePolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-inline-policy <user-id> <policy-name>",
		Short: "Delete a user's inline policy by name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.DeleteUserInlinePolicy(cmd.Context(), args[0], args[1]); err != nil {
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

// newIamUserDetachPolicyCommand builds `basaltic iam user detach-policy`.
func newIamUserDetachPolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach-policy <user-id> <policy-id>",
		Short: "Detach policy from user",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.DetachUserPolicy(cmd.Context(), args[0], args[1]); err != nil {
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

// newIamUserGetInlinePolicyCommand builds `basaltic iam user get-inline-policy`.
func newIamUserGetInlinePolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-inline-policy <user-id> <policy-name>",
		Short: "Get a user's inline policy by name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetUserInlinePolicy(cmd.Context(), args[0], args[1])
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

// newIamUserGetPermissionBoundaryCommand builds `basaltic iam user get-permission-boundary`.
func newIamUserGetPermissionBoundaryCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-permission-boundary <user-id>",
		Short: "Get a user's permission boundary",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetUserPermissionBoundary(cmd.Context(), args[0])
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

// newIamUserListGroupsCommand builds `basaltic iam user list-groups`.
func newIamUserListGroupsCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-groups <user-id>",
		Short: "List user groups",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListUserGroups(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamUserListInlinePoliciesCommand builds `basaltic iam user list-inline-policies`.
func newIamUserListInlinePoliciesCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-inline-policies <user-id>",
		Short: "List a user's inline policies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListUserInlinePolicies(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamUserListPoliciesCommand builds `basaltic iam user list-policies`.
func newIamUserListPoliciesCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-policies <user-id>",
		Short: "List user policies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListUserPolicies(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamUserRemoveCommand builds `basaltic iam user remove`.
func newIamUserRemoveCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <user-id>",
		Short: "Remove user from organization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.RemoveUser(cmd.Context(), args[0]); err != nil {
				return err
			}
			state.Printer().Done("Remove requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamUserRemoveGroupCommand builds `basaltic iam user remove-group`.
func newIamUserRemoveGroupCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-group <user-id> <group-id>",
		Short: "Remove user from group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.RemoveUserFromGroup(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			state.Printer().Done("Remove group requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newIamUserRemovePermissionBoundaryCommand builds `basaltic iam user remove-permission-boundary`.
func newIamUserRemovePermissionBoundaryCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-permission-boundary <user-id>",
		Short: "Remove a user's permission boundary",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := iamClient(state)
			if err != nil {
				return err
			}
			if err := c.RemoveUserPermissionBoundary(cmd.Context(), args[0]); err != nil {
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

// newIamUserSetInlinePolicyCommand builds `basaltic iam user set-inline-policy`.
func newIamUserSetInlinePolicyCommand(state *cli.State) *cobra.Command {
	var body iam.PutInlinePolicyRequest
	var bodyFile string
	var documentFlag string
	cmd := &cobra.Command{
		Use:   "set-inline-policy <user-id> <policy-name>",
		Short: "Create or replace a user's inline policy",
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
			out, err := c.PutUserInlinePolicy(cmd.Context(), args[0], args[1], &body)
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

// newIamUserSetPermissionBoundaryCommand builds `basaltic iam user set-permission-boundary`.
func newIamUserSetPermissionBoundaryCommand(state *cli.State) *cobra.Command {
	var body iam.SetBoundaryRequest
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "set-permission-boundary <user-id>",
		Short: "Set a user's permission boundary",
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
			if err := c.SetUserPermissionBoundary(cmd.Context(), args[0], &body); err != nil {
				return err
			}
			state.Printer().Done("Set permission boundary requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.PolicyID, "policy-id", "", "Policy id")
	_ = cmd.MarkFlagRequired("policy-id")
	return cmd
}
