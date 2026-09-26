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
	"github.com/basaltic-sh/sdk-go/workspace"

	"github.com/basaltic-sh/cli/internal/cli"
)

func init() { cli.RegisterService(newWorkspaceCommand) }

// newWorkspaceCommand builds `basaltic workspace`.
func newWorkspaceCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "Organizations, accounts, people, groups and organization policies",
	}
	cmd.AddCommand(newWorkspaceAccountCommand(state))
	cmd.AddCommand(newWorkspaceAccountRoleCommand(state))
	cmd.AddCommand(newWorkspaceGroupCommand(state))
	cmd.AddCommand(newWorkspaceInvitationCommand(state))
	cmd.AddCommand(newWorkspaceOrganizationCommand(state))
	cmd.AddCommand(newWorkspacePolicyCommand(state))
	cmd.AddCommand(newWorkspaceRoleCommand(state))
	cmd.AddCommand(newWorkspaceServiceAccountCommand(state))
	cmd.AddCommand(newWorkspaceUserCommand(state))
	return cmd
}

// workspaceClient builds the service client, resolving credentials on first use.
func workspaceClient(state *cli.State) (*workspace.Client, error) {
	cfg, err := state.SDK()
	if err != nil {
		return nil, err
	}
	return workspace.New(cfg), nil
}

// newWorkspaceAccountCommand builds `basaltic workspace account`.
func newWorkspaceAccountCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account",
		Short:   "Accounts",
		Aliases: []string{"accounts"},
	}
	cmd.AddCommand(newWorkspaceAccountListCommand(state))
	cmd.AddCommand(newWorkspaceAccountGetCommand(state))
	cmd.AddCommand(newWorkspaceAccountCreateCommand(state))
	cmd.AddCommand(newWorkspaceAccountUpdateCommand(state))
	cmd.AddCommand(newWorkspaceAccountDeleteCommand(state))
	cmd.AddCommand(newWorkspaceAccountAssignRoleAssignmentCommand(state))
	cmd.AddCommand(newWorkspaceAccountGetResourceCommand(state))
	cmd.AddCommand(newWorkspaceAccountListRoleAssignmentsCommand(state))
	cmd.AddCommand(newWorkspaceAccountRemoveRoleAssignmentCommand(state))
	return cmd
}

// newWorkspaceAccountListCommand builds `basaltic workspace account list`.
func newWorkspaceAccountListCommand(state *cli.State) *cobra.Command {
	var params workspace.ListAccountsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List accounts",
		Args:  cobra.ExactArgs(0),
		Long:  "List accounts.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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
	f.StringVar(&params.CRN, "crn", "", "Exact returned CRN, combined with name using AND before pagination")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Exact resource name, combined with crn using AND before pagination")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newWorkspaceAccountGetCommand builds `basaltic workspace account get`.
func newWorkspaceAccountGetCommand(state *cli.State) *cobra.Command {
	var scope workspace.ListAccountsParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Get account",
		Args:  cobra.ExactArgs(1),
		Long:  "Get account.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetAccountByReference(cmd.Context(), args[0], &scope)
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

// newWorkspaceAccountCreateCommand builds `basaltic workspace account create`.
func newWorkspaceAccountCreateCommand(state *cli.State) *cobra.Command {
	var body workspace.CreateAccountRequest
	var bodyFile string
	var descriptionFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create account",
		Args:  cobra.ExactArgs(0),
		Long:  "Create account.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceAccountUpdateCommand builds `basaltic workspace account update`.
func newWorkspaceAccountUpdateCommand(state *cli.State) *cobra.Command {
	var body workspace.UpdateAccountRequest
	var bodyFile string
	var descriptionFlag string
	var nameFlag string
	cmd := &cobra.Command{
		Use:   "update <account-id>",
		Short: "Update account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceAccountDeleteCommand builds `basaltic workspace account delete`.
func newWorkspaceAccountDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <account-id>",
		Short: "Delete account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceAccountAssignRoleAssignmentCommand builds `basaltic workspace account assign-role-assignment`.
func newWorkspaceAccountAssignRoleAssignmentCommand(state *cli.State) *cobra.Command {
	var body workspace.AccountRoleAssignmentCreateRequest
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "assign-role-assignment <account-id>",
		Short: "Assign account role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			out, err := c.AssignAccountRole(cmd.Context(), args[0], &body)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.PrincipalID, "principal-id", "", "Immutable UUID of a user or users-only group in this organization")
	_ = cmd.MarkFlagRequired("principal-id")
	f.StringVar(&body.PrincipalType, "principal-type", "", "Principal type (one of: user, group)")
	_ = cmd.MarkFlagRequired("principal-type")
	f.StringVar(&body.RoleID, "role-id", "", "Immutable UUID of a role owned by the target account")
	_ = cmd.MarkFlagRequired("role-id")
	return cmd
}

// newWorkspaceAccountGetResourceCommand builds `basaltic workspace account get-resource`.
func newWorkspaceAccountGetResourceCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-resource <account-id>",
		Short: "Check account resource presence",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetAccountResources(cmd.Context(), args[0])
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

// newWorkspaceAccountListRoleAssignmentsCommand builds `basaltic workspace account list-role-assignments`.
func newWorkspaceAccountListRoleAssignmentsCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-role-assignments <account-id>",
		Short: "List account role assignments",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListAccountRoleAssignments(cmd.Context(), args[0])
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

// newWorkspaceAccountRemoveRoleAssignmentCommand builds `basaltic workspace account remove-role-assignment`.
func newWorkspaceAccountRemoveRoleAssignmentCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-role-assignment <account-id> <assignment-id>",
		Short: "Remove account role assignment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
			if err != nil {
				return err
			}
			if err := c.RemoveAccountRoleAssignment(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			state.Printer().Done("Remove role assignment requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newWorkspaceAccountRoleCommand builds `basaltic workspace account-role`.
func newWorkspaceAccountRoleCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account-role",
		Short:   "Account roles",
		Aliases: []string{"account-roles"},
	}
	cmd.AddCommand(newWorkspaceAccountRoleListCommand(state))
	return cmd
}

// newWorkspaceAccountRoleListCommand builds `basaltic workspace account-role list`.
func newWorkspaceAccountRoleListCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available account roles",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListAccountRoles(cmd.Context())
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

// newWorkspaceGroupCommand builds `basaltic workspace group`.
func newWorkspaceGroupCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "group",
		Short:   "Groups",
		Aliases: []string{"groups"},
	}
	cmd.AddCommand(newWorkspaceGroupListCommand(state))
	cmd.AddCommand(newWorkspaceGroupGetCommand(state))
	cmd.AddCommand(newWorkspaceGroupCreateCommand(state))
	cmd.AddCommand(newWorkspaceGroupUpdateCommand(state))
	cmd.AddCommand(newWorkspaceGroupDeleteCommand(state))
	cmd.AddCommand(newWorkspaceGroupAttachPolicyCommand(state))
	cmd.AddCommand(newWorkspaceGroupDeleteInlinePolicyCommand(state))
	cmd.AddCommand(newWorkspaceGroupDetachPolicyCommand(state))
	cmd.AddCommand(newWorkspaceGroupGetInlinePolicyCommand(state))
	cmd.AddCommand(newWorkspaceGroupListInlinePoliciesCommand(state))
	cmd.AddCommand(newWorkspaceGroupListPoliciesCommand(state))
	cmd.AddCommand(newWorkspaceGroupListUsersCommand(state))
	cmd.AddCommand(newWorkspaceGroupSetInlinePolicyCommand(state))
	return cmd
}

// newWorkspaceGroupListCommand builds `basaltic workspace group list`.
func newWorkspaceGroupListCommand(state *cli.State) *cobra.Command {
	var params workspace.ListGroupsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List groups",
		Args:  cobra.ExactArgs(0),
		Long:  "List groups.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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
	f.StringVar(&params.CRN, "crn", "", "Exact returned CRN, combined with name using AND before pagination")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Exact resource name, combined with crn using AND before pagination")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newWorkspaceGroupGetCommand builds `basaltic workspace group get`.
func newWorkspaceGroupGetCommand(state *cli.State) *cobra.Command {
	var scope workspace.ListGroupsParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Get group",
		Args:  cobra.ExactArgs(1),
		Long:  "Get group.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetGroupByReference(cmd.Context(), args[0], &scope)
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

// newWorkspaceGroupCreateCommand builds `basaltic workspace group create`.
func newWorkspaceGroupCreateCommand(state *cli.State) *cobra.Command {
	var body workspace.GroupCreateRequest
	var bodyFile string
	var descriptionFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create group",
		Args:  cobra.ExactArgs(0),
		Long:  "Create group.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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
	f.StringVar(&body.Name, "name", "", "Resource names must not start with the literal crn: prefix or be UUIDs (canonical, compact, braced, or urn:uuid: forms, in either case)")
	_ = cmd.MarkFlagRequired("name")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newWorkspaceGroupUpdateCommand builds `basaltic workspace group update`.
func newWorkspaceGroupUpdateCommand(state *cli.State) *cobra.Command {
	var body workspace.GroupUpdateRequest
	var bodyFile string
	var descriptionFlag string
	cmd := &cobra.Command{
		Use:   "update <group-id>",
		Short: "Update group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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
	return cmd
}

// newWorkspaceGroupDeleteCommand builds `basaltic workspace group delete`.
func newWorkspaceGroupDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <group-id>",
		Short: "Delete group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceGroupAttachPolicyCommand builds `basaltic workspace group attach-policy`.
func newWorkspaceGroupAttachPolicyCommand(state *cli.State) *cobra.Command {
	var body workspace.PolicyAttachRequest
	var bodyFile string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "attach-policy <group-id>",
		Short: "Attach policy to group",
		Args:  cobra.ExactArgs(1),
		Long:  "Attach policy to group.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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
	f.StringVar((*string)(&body.Policy), "policy", "", "Policy")
	_ = cmd.MarkFlagRequired("policy")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newWorkspaceGroupDeleteInlinePolicyCommand builds `basaltic workspace group delete-inline-policy`.
func newWorkspaceGroupDeleteInlinePolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-inline-policy <group-id> <policy-name>",
		Short: "Delete a group's inline policy by name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceGroupDetachPolicyCommand builds `basaltic workspace group detach-policy`.
func newWorkspaceGroupDetachPolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach-policy <group-id> <policy-id>",
		Short: "Detach policy from group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceGroupGetInlinePolicyCommand builds `basaltic workspace group get-inline-policy`.
func newWorkspaceGroupGetInlinePolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-inline-policy <group-id> <policy-name>",
		Short: "Get a group's inline policy by name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceGroupListInlinePoliciesCommand builds `basaltic workspace group list-inline-policies`.
func newWorkspaceGroupListInlinePoliciesCommand(state *cli.State) *cobra.Command {
	var params workspace.ListGroupInlinePoliciesParams
	cmd := &cobra.Command{
		Use:   "list-inline-policies <group-id>",
		Short: "List a group's inline policies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListGroupInlinePolicies(cmd.Context(), args[0], &params)
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

// newWorkspaceGroupListPoliciesCommand builds `basaltic workspace group list-policies`.
func newWorkspaceGroupListPoliciesCommand(state *cli.State) *cobra.Command {
	var params workspace.ListGroupPoliciesParams
	cmd := &cobra.Command{
		Use:   "list-policies <group-id>",
		Short: "List group policies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListGroupPolicies(cmd.Context(), args[0], &params)
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

// newWorkspaceGroupListUsersCommand builds `basaltic workspace group list-users`.
func newWorkspaceGroupListUsersCommand(state *cli.State) *cobra.Command {
	var params workspace.ListGroupUsersParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list-users <group-id>",
		Short: "List group users",
		Args:  cobra.ExactArgs(1),
		Long:  "List group users.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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
	f.StringVar(&params.CRN, "crn", "", "Exact returned CRN, combined with name using AND before pagination")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Exact resource name, combined with crn using AND before pagination")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newWorkspaceGroupSetInlinePolicyCommand builds `basaltic workspace group set-inline-policy`.
func newWorkspaceGroupSetInlinePolicyCommand(state *cli.State) *cobra.Command {
	var body workspace.PutInlinePolicyRequest
	var bodyFile string
	var documentFlag string
	cmd := &cobra.Command{
		Use:   "set-inline-policy <group-id> <policy-name>",
		Short: "Create or replace a group's inline policy",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceInvitationCommand builds `basaltic workspace invitation`.
func newWorkspaceInvitationCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "invitation",
		Short:   "Invitations",
		Aliases: []string{"invitations"},
	}
	cmd.AddCommand(newWorkspaceInvitationListCommand(state))
	cmd.AddCommand(newWorkspaceInvitationGetCommand(state))
	cmd.AddCommand(newWorkspaceInvitationCancelCommand(state))
	return cmd
}

// newWorkspaceInvitationListCommand builds `basaltic workspace invitation list`.
func newWorkspaceInvitationListCommand(state *cli.State) *cobra.Command {
	var params workspace.ListInvitationsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List invitations",
		Args:  cobra.ExactArgs(0),
		Long:  "List invitations.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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
	f.StringVar(&params.CRN, "crn", "", "Exact returned CRN, combined with name using AND before pagination")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Exact resource name, combined with crn using AND before pagination")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newWorkspaceInvitationGetCommand builds `basaltic workspace invitation get`.
func newWorkspaceInvitationGetCommand(state *cli.State) *cobra.Command {
	var scope workspace.ListInvitationsParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Get invitation",
		Args:  cobra.ExactArgs(1),
		Long:  "Get invitation.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetInvitationByReference(cmd.Context(), args[0], &scope)
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

// newWorkspaceInvitationCancelCommand builds `basaltic workspace invitation cancel`.
func newWorkspaceInvitationCancelCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel <invitation-id>",
		Short: "Cancel invitation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceOrganizationCommand builds `basaltic workspace organization`.
func newWorkspaceOrganizationCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "organization",
		Short:   "Organizations",
		Aliases: []string{"organizations"},
	}
	cmd.AddCommand(newWorkspaceOrganizationListCommand(state))
	cmd.AddCommand(newWorkspaceOrganizationGetCommand(state))
	cmd.AddCommand(newWorkspaceOrganizationUpdateCommand(state))
	cmd.AddCommand(newWorkspaceOrganizationDeleteCommand(state))
	return cmd
}

// newWorkspaceOrganizationListCommand builds `basaltic workspace organization list`.
func newWorkspaceOrganizationListCommand(state *cli.State) *cobra.Command {
	var params workspace.ListOrganizationsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List organizations",
		Args:  cobra.ExactArgs(0),
		Long:  "List organizations.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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
	f.StringVar(&params.CRN, "crn", "", "Exact returned CRN, combined with name using AND before pagination")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Exact resource name, combined with crn using AND before pagination")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newWorkspaceOrganizationGetCommand builds `basaltic workspace organization get`.
func newWorkspaceOrganizationGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <organization-id>",
		Short: "Get organization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceOrganizationUpdateCommand builds `basaltic workspace organization update`.
func newWorkspaceOrganizationUpdateCommand(state *cli.State) *cobra.Command {
	var body workspace.OrganizationUpdateRequest
	var bodyFile string
	var descriptionFlag string
	var nameFlag string
	cmd := &cobra.Command{
		Use:   "update <organization-id>",
		Short: "Update organization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceOrganizationDeleteCommand builds `basaltic workspace organization delete`.
func newWorkspaceOrganizationDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <organization-id>",
		Short: "Delete organization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspacePolicyCommand builds `basaltic workspace policy`.
func newWorkspacePolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "policy",
		Short:   "Policies",
		Aliases: []string{"policies"},
	}
	cmd.AddCommand(newWorkspacePolicyListCommand(state))
	cmd.AddCommand(newWorkspacePolicyGetCommand(state))
	cmd.AddCommand(newWorkspacePolicyCreateCommand(state))
	cmd.AddCommand(newWorkspacePolicyUpdateCommand(state))
	cmd.AddCommand(newWorkspacePolicyDeleteCommand(state))
	cmd.AddCommand(newWorkspacePolicyListGroupsCommand(state))
	cmd.AddCommand(newWorkspacePolicyListRolesCommand(state))
	cmd.AddCommand(newWorkspacePolicyListServiceAccountsCommand(state))
	cmd.AddCommand(newWorkspacePolicyListUsersCommand(state))
	return cmd
}

// newWorkspacePolicyListCommand builds `basaltic workspace policy list`.
func newWorkspacePolicyListCommand(state *cli.State) *cobra.Command {
	var params workspace.ListPoliciesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List policies",
		Args:  cobra.ExactArgs(0),
		Long:  "List policies.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspacePolicyGetCommand builds `basaltic workspace policy get`.
func newWorkspacePolicyGetCommand(state *cli.State) *cobra.Command {
	var scope workspace.ListPoliciesParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Get policy",
		Args:  cobra.ExactArgs(1),
		Long:  "Get policy.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspacePolicyCreateCommand builds `basaltic workspace policy create`.
func newWorkspacePolicyCreateCommand(state *cli.State) *cobra.Command {
	var body workspace.PolicyCreateRequest
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
			c, err := workspaceClient(state)
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

// newWorkspacePolicyUpdateCommand builds `basaltic workspace policy update`.
func newWorkspacePolicyUpdateCommand(state *cli.State) *cobra.Command {
	var body workspace.PolicyUpdateRequest
	var bodyFile string
	var descriptionFlag string
	var documentFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <policy-id>",
		Short: "Update policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspacePolicyDeleteCommand builds `basaltic workspace policy delete`.
func newWorkspacePolicyDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <policy-id>",
		Short: "Delete policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspacePolicyListGroupsCommand builds `basaltic workspace policy list-groups`.
func newWorkspacePolicyListGroupsCommand(state *cli.State) *cobra.Command {
	var params workspace.ListPolicyGroupsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list-groups <policy-id>",
		Short: "List groups with policy",
		Args:  cobra.ExactArgs(1),
		Long:  "List groups with policy.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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
	f.StringVar(&params.CRN, "crn", "", "Exact returned CRN, combined with name using AND before pagination")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Exact resource name, combined with crn using AND before pagination")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newWorkspacePolicyListRolesCommand builds `basaltic workspace policy list-roles`.
func newWorkspacePolicyListRolesCommand(state *cli.State) *cobra.Command {
	var params workspace.ListPolicyRolesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list-roles <policy-id>",
		Short: "List roles with policy",
		Args:  cobra.ExactArgs(1),
		Long:  "List roles with policy.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspacePolicyListServiceAccountsCommand builds `basaltic workspace policy list-service-accounts`.
func newWorkspacePolicyListServiceAccountsCommand(state *cli.State) *cobra.Command {
	var params workspace.ListPolicyServiceAccountsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list-service-accounts <policy-id>",
		Short: "List service accounts with policy",
		Args:  cobra.ExactArgs(1),
		Long:  "List service accounts with policy.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspacePolicyListUsersCommand builds `basaltic workspace policy list-users`.
func newWorkspacePolicyListUsersCommand(state *cli.State) *cobra.Command {
	var params workspace.ListPolicyUsersParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list-users <policy-id>",
		Short: "List users with policy",
		Args:  cobra.ExactArgs(1),
		Long:  "List users with policy.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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
	f.StringVar(&params.CRN, "crn", "", "Exact returned CRN, combined with name using AND before pagination")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Exact resource name, combined with crn using AND before pagination")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newWorkspaceRoleCommand builds `basaltic workspace role`.
func newWorkspaceRoleCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "role",
		Short:   "Roles",
		Aliases: []string{"roles"},
	}
	cmd.AddCommand(newWorkspaceRoleAttachPolicyCommand(state))
	cmd.AddCommand(newWorkspaceRoleDetachPolicyCommand(state))
	cmd.AddCommand(newWorkspaceRoleListPoliciesCommand(state))
	return cmd
}

// newWorkspaceRoleAttachPolicyCommand builds `basaltic workspace role attach-policy`.
func newWorkspaceRoleAttachPolicyCommand(state *cli.State) *cobra.Command {
	var body workspace.OrganizationPolicyAttachRequest
	var bodyFile string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "attach-policy <role-id>",
		Short: "Attach policy to role",
		Args:  cobra.ExactArgs(1),
		Long:  "Attach policy to role.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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
	f.StringVar(&body.PolicyID, "policy-id", "", "Immutable UUID of the organization policy to attach")
	_ = cmd.MarkFlagRequired("policy-id")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newWorkspaceRoleDetachPolicyCommand builds `basaltic workspace role detach-policy`.
func newWorkspaceRoleDetachPolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach-policy <role-id> <policy-id>",
		Short: "Detach policy from role",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceRoleListPoliciesCommand builds `basaltic workspace role list-policies`.
func newWorkspaceRoleListPoliciesCommand(state *cli.State) *cobra.Command {
	var params workspace.ListRolePoliciesParams
	cmd := &cobra.Command{
		Use:   "list-policies <role-id>",
		Short: "List role policies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceServiceAccountCommand builds `basaltic workspace service-account`.
func newWorkspaceServiceAccountCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "service-account",
		Short:   "Service accounts",
		Aliases: []string{"service-accounts"},
	}
	cmd.AddCommand(newWorkspaceServiceAccountAttachPolicyCommand(state))
	cmd.AddCommand(newWorkspaceServiceAccountDetachPolicyCommand(state))
	cmd.AddCommand(newWorkspaceServiceAccountListPoliciesCommand(state))
	return cmd
}

// newWorkspaceServiceAccountAttachPolicyCommand builds `basaltic workspace service-account attach-policy`.
func newWorkspaceServiceAccountAttachPolicyCommand(state *cli.State) *cobra.Command {
	var body workspace.OrganizationPolicyAttachRequest
	var bodyFile string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "attach-policy <service-account-id>",
		Short: "Attach policy to service account",
		Args:  cobra.ExactArgs(1),
		Long:  "Attach policy to service account.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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
	f.StringVar(&body.PolicyID, "policy-id", "", "Immutable UUID of the organization policy to attach")
	_ = cmd.MarkFlagRequired("policy-id")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newWorkspaceServiceAccountDetachPolicyCommand builds `basaltic workspace service-account detach-policy`.
func newWorkspaceServiceAccountDetachPolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach-policy <service-account-id> <policy-id>",
		Short: "Detach policy from service account",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceServiceAccountListPoliciesCommand builds `basaltic workspace service-account list-policies`.
func newWorkspaceServiceAccountListPoliciesCommand(state *cli.State) *cobra.Command {
	var params workspace.ListServiceAccountPoliciesParams
	cmd := &cobra.Command{
		Use:   "list-policies <service-account-id>",
		Short: "List service account policies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceUserCommand builds `basaltic workspace user`.
func newWorkspaceUserCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "user",
		Short:   "Users",
		Aliases: []string{"users"},
	}
	cmd.AddCommand(newWorkspaceUserListCommand(state))
	cmd.AddCommand(newWorkspaceUserGetCommand(state))
	cmd.AddCommand(newWorkspaceUserAddCommand(state))
	cmd.AddCommand(newWorkspaceUserAddGroupCommand(state))
	cmd.AddCommand(newWorkspaceUserAttachPolicyCommand(state))
	cmd.AddCommand(newWorkspaceUserDeleteInlinePolicyCommand(state))
	cmd.AddCommand(newWorkspaceUserDetachPolicyCommand(state))
	cmd.AddCommand(newWorkspaceUserGetInlinePolicyCommand(state))
	cmd.AddCommand(newWorkspaceUserGetPermissionBoundaryCommand(state))
	cmd.AddCommand(newWorkspaceUserListGroupsCommand(state))
	cmd.AddCommand(newWorkspaceUserListInlinePoliciesCommand(state))
	cmd.AddCommand(newWorkspaceUserListPoliciesCommand(state))
	cmd.AddCommand(newWorkspaceUserRemoveCommand(state))
	cmd.AddCommand(newWorkspaceUserRemoveGroupCommand(state))
	cmd.AddCommand(newWorkspaceUserRemovePermissionBoundaryCommand(state))
	cmd.AddCommand(newWorkspaceUserSetInlinePolicyCommand(state))
	cmd.AddCommand(newWorkspaceUserSetPermissionBoundaryCommand(state))
	return cmd
}

// newWorkspaceUserListCommand builds `basaltic workspace user list`.
func newWorkspaceUserListCommand(state *cli.State) *cobra.Command {
	var params workspace.ListUsersParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List users",
		Args:  cobra.ExactArgs(0),
		Long:  "List users.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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
	f.StringVar(&params.CRN, "crn", "", "Exact returned CRN, combined with name using AND before pagination")
	f.IntVar(&params.Limit, "limit", 0, "Maximum number of items to return")
	f.StringVar(&params.Marker, "marker", "", "Opaque pagination cursor")
	f.StringVar(&params.Name, "name", "", "Exact resource name, combined with crn using AND before pagination")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newWorkspaceUserGetCommand builds `basaltic workspace user get`.
func newWorkspaceUserGetCommand(state *cli.State) *cobra.Command {
	var scope workspace.ListUsersParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Get user",
		Args:  cobra.ExactArgs(1),
		Long:  "Get user.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
			if err != nil {
				return err
			}
			out, err := c.GetUserByReference(cmd.Context(), args[0], &scope)
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

// newWorkspaceUserAddCommand builds `basaltic workspace user add`.
func newWorkspaceUserAddCommand(state *cli.State) *cobra.Command {
	var body workspace.UserAddRequest
	var bodyFile string
	var groupsFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add user to organization",
		Args:  cobra.ExactArgs(0),
		Long:  "Add user to organization.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
			if err != nil {
				return err
			}
			if bodyFile != "" {
				if err := loadBody(bodyFile, &body); err != nil {
					return err
				}
			}
			if groupsFlag != "" {
				if err := json.Unmarshal([]byte(groupsFlag), &body.Groups); err != nil {
					return fmt.Errorf("--groups: %w", err)
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
	f.StringVar(&groupsFlag, "groups", "", "Groups to assign when the invitation is accepted (JSON)")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newWorkspaceUserAddGroupCommand builds `basaltic workspace user add-group`.
func newWorkspaceUserAddGroupCommand(state *cli.State) *cobra.Command {
	var body workspace.UserGroupAddRequest
	var bodyFile string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "add-group <user-id>",
		Short: "Add user to group",
		Args:  cobra.ExactArgs(1),
		Long:  "Add user to group.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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
	f.StringVar((*string)(&body.Group), "group", "", "Group")
	_ = cmd.MarkFlagRequired("group")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newWorkspaceUserAttachPolicyCommand builds `basaltic workspace user attach-policy`.
func newWorkspaceUserAttachPolicyCommand(state *cli.State) *cobra.Command {
	var body workspace.PolicyAttachRequest
	var bodyFile string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "attach-policy <user-id>",
		Short: "Attach policy to user",
		Args:  cobra.ExactArgs(1),
		Long:  "Attach policy to user.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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
	f.StringVar((*string)(&body.Policy), "policy", "", "Policy")
	_ = cmd.MarkFlagRequired("policy")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newWorkspaceUserDeleteInlinePolicyCommand builds `basaltic workspace user delete-inline-policy`.
func newWorkspaceUserDeleteInlinePolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-inline-policy <user-id> <policy-name>",
		Short: "Delete a user's inline policy by name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceUserDetachPolicyCommand builds `basaltic workspace user detach-policy`.
func newWorkspaceUserDetachPolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach-policy <user-id> <policy-id>",
		Short: "Detach policy from user",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceUserGetInlinePolicyCommand builds `basaltic workspace user get-inline-policy`.
func newWorkspaceUserGetInlinePolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-inline-policy <user-id> <policy-name>",
		Short: "Get a user's inline policy by name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceUserGetPermissionBoundaryCommand builds `basaltic workspace user get-permission-boundary`.
func newWorkspaceUserGetPermissionBoundaryCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-permission-boundary <user-id>",
		Short: "Get a user's permission boundary",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceUserListGroupsCommand builds `basaltic workspace user list-groups`.
func newWorkspaceUserListGroupsCommand(state *cli.State) *cobra.Command {
	var params workspace.ListUserGroupsParams
	cmd := &cobra.Command{
		Use:   "list-groups <user-id>",
		Short: "List user groups",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListUserGroups(cmd.Context(), args[0], &params)
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

// newWorkspaceUserListInlinePoliciesCommand builds `basaltic workspace user list-inline-policies`.
func newWorkspaceUserListInlinePoliciesCommand(state *cli.State) *cobra.Command {
	var params workspace.ListUserInlinePoliciesParams
	cmd := &cobra.Command{
		Use:   "list-inline-policies <user-id>",
		Short: "List a user's inline policies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListUserInlinePolicies(cmd.Context(), args[0], &params)
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

// newWorkspaceUserListPoliciesCommand builds `basaltic workspace user list-policies`.
func newWorkspaceUserListPoliciesCommand(state *cli.State) *cobra.Command {
	var params workspace.ListUserPoliciesParams
	cmd := &cobra.Command{
		Use:   "list-policies <user-id>",
		Short: "List user policies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
			if err != nil {
				return err
			}
			page, err := c.ListUserPolicies(cmd.Context(), args[0], &params)
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

// newWorkspaceUserRemoveCommand builds `basaltic workspace user remove`.
func newWorkspaceUserRemoveCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <user-id>",
		Short: "Remove user from organization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceUserRemoveGroupCommand builds `basaltic workspace user remove-group`.
func newWorkspaceUserRemoveGroupCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-group <user-id> <group-id>",
		Short: "Remove user from group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceUserRemovePermissionBoundaryCommand builds `basaltic workspace user remove-permission-boundary`.
func newWorkspaceUserRemovePermissionBoundaryCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-permission-boundary <user-id>",
		Short: "Remove a user's permission boundary",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceUserSetInlinePolicyCommand builds `basaltic workspace user set-inline-policy`.
func newWorkspaceUserSetInlinePolicyCommand(state *cli.State) *cobra.Command {
	var body workspace.PutInlinePolicyRequest
	var bodyFile string
	var documentFlag string
	cmd := &cobra.Command{
		Use:   "set-inline-policy <user-id> <policy-name>",
		Short: "Create or replace a user's inline policy",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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

// newWorkspaceUserSetPermissionBoundaryCommand builds `basaltic workspace user set-permission-boundary`.
func newWorkspaceUserSetPermissionBoundaryCommand(state *cli.State) *cobra.Command {
	var body workspace.SetBoundaryRequest
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "set-permission-boundary <user-id>",
		Short: "Set a user's permission boundary",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := workspaceClient(state)
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
	f.StringVar((*string)(&body.Policy), "policy", "", "Policy")
	_ = cmd.MarkFlagRequired("policy")
	return cmd
}
