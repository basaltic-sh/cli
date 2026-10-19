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
	"github.com/basaltic-sh/sdk-go/registry"

	"github.com/basaltic-sh/cli/internal/cli"
)

func init() { cli.RegisterService(newRegistryCommand) }

// newRegistryCommand builds `basaltic registry`.
func newRegistryCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registry",
		Short: "KMS bindings use `kms_key`: a UUID, account-scoped name, or complete `crn:kms:<region>:<account>:key/<name>` in the authenticated account and current region",
		Long:  "KMS bindings use `kms_key`: a UUID, account-scoped name, or complete `crn:kms:<region>:<account>:key/<name>` in the authenticated account and current region.\n\nThis is a regional service: it acts in the region from --region, the\nBASALTIC_REGION environment variable, or the profile.",
	}
	cmd.AddCommand(newRegistryAuthorizationTokenCommand(state))
	cmd.AddCommand(newRegistryEncryptionConfigurationCommand(state))
	cmd.AddCommand(newRegistryImageCommand(state))
	cmd.AddCommand(newRegistryPullThroughCacheRuleCommand(state))
	cmd.AddCommand(newRegistryRepositoryCommand(state))
	return cmd
}

// registryClient builds the service client, resolving credentials on first use.
func registryClient(state *cli.State, path string) (*registry.Client, error) {
	cfg, err := state.ServiceSDK("registry", path)
	if err != nil {
		return nil, err
	}
	return registry.New(cfg), nil
}

// newRegistryAuthorizationTokenCommand builds `basaltic registry authorization-token`.
func newRegistryAuthorizationTokenCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "authorization-token",
		Short:   "Authorization tokens",
		Aliases: []string{"authorization-tokens"},
	}
	cmd.AddCommand(newRegistryAuthorizationTokenCreateCommand(state))
	return cmd
}

// newRegistryAuthorizationTokenCreateCommand builds `basaltic registry authorization-token create`.
func newRegistryAuthorizationTokenCreateCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Get a docker login credential",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/authorization-token")
			if err != nil {
				return err
			}
			out, err := c.GetAuthorizationToken(cmd.Context())
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

// newRegistryEncryptionConfigurationCommand builds `basaltic registry encryption-configuration`.
func newRegistryEncryptionConfigurationCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "encryption-configuration",
		Short:   "Encryption configurations",
		Aliases: []string{"encryption-configurations"},
	}
	cmd.AddCommand(newRegistryEncryptionConfigurationListCommand(state))
	cmd.AddCommand(newRegistryEncryptionConfigurationDeleteCommand(state))
	cmd.AddCommand(newRegistryEncryptionConfigurationGetRetainedDataCommand(state))
	cmd.AddCommand(newRegistryEncryptionConfigurationSetCommand(state))
	return cmd
}

// newRegistryEncryptionConfigurationListCommand builds `basaltic registry encryption-configuration list`.
func newRegistryEncryptionConfigurationListCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Get the encryption configuration",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/encryption-configuration")
			if err != nil {
				return err
			}
			out, err := c.GetEncryptionConfiguration(cmd.Context())
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

// newRegistryEncryptionConfigurationDeleteCommand builds `basaltic registry encryption-configuration delete`.
func newRegistryEncryptionConfigurationDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete encryption configuration",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/encryption-configuration")
			if err != nil {
				return err
			}
			if err := c.DeleteEncryptionConfiguration(cmd.Context()); err != nil {
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

// newRegistryEncryptionConfigurationGetRetainedDataCommand builds `basaltic registry encryption-configuration get-retained-data`.
func newRegistryEncryptionConfigurationGetRetainedDataCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-retained-data",
		Short: "Check retained registry data presence",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/encryption-configuration/retained-data")
			if err != nil {
				return err
			}
			out, err := c.GetRetainedRegistryPresence(cmd.Context())
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

// newRegistryEncryptionConfigurationSetCommand builds `basaltic registry encryption-configuration set`.
func newRegistryEncryptionConfigurationSetCommand(state *cli.State) *cobra.Command {
	var body registry.EncryptionConfigurationRequest
	var bodyFile string
	var kmsKeyFlag string
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set the encryption configuration",
		Args:  cobra.ExactArgs(0),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"encryption_type": "encryption-type",
				"kms_key":         "kms-key",
			}, []string{"encryption_type"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/encryption-configuration")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("kms-key") {
				body.KMSKey = &kmsKeyFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.PutEncryptionConfiguration(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.EncryptionType, "encryption-type", "", "Required (one of: platform_managed, kms)")
	f.StringVar(&kmsKeyFlag, "kms-key", "", "KMS key UUID, CRN or account-scoped name")
	return cmd
}

// newRegistryImageCommand builds `basaltic registry image`.
func newRegistryImageCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "image",
		Short:   "Images",
		Aliases: []string{"images"},
	}
	cmd.AddCommand(newRegistryImageListCommand(state))
	cmd.AddCommand(newRegistryImageGetCommand(state))
	cmd.AddCommand(newRegistryImageBatchDeleteCommand(state))
	return cmd
}

// newRegistryImageListCommand builds `basaltic registry image list`.
func newRegistryImageListCommand(state *cli.State) *cobra.Command {
	var params registry.DescribeImagesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list <repository-id>",
		Short: "List images",
		Args:  cobra.ExactArgs(1),
		Long:  "List images.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/repositories/{repository_id}/images")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.DescribeImagesAll(cmd.Context(), args[0], &params))
			}
			page, err := c.DescribeImages(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Maximum items to return")
	f.StringVar(&params.Marker, "marker", "", "Id of the last item on the previous page")
	f.StringVar(&params.TagPrefix, "tag-prefix", "", "Return only images carrying a tag with this prefix")
	f.StringVar(&params.TagStatus, "tag-status", "", "Restrict to one population")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newRegistryImageGetCommand builds `basaltic registry image get`.
func newRegistryImageGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <repository-id> <image-reference>",
		Short: "Describe an image",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/repositories/{repository_id}/images/{image_reference}")
			if err != nil {
				return err
			}
			out, err := c.DescribeImage(cmd.Context(), args[0], args[1])
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

// newRegistryImageBatchDeleteCommand builds `basaltic registry image batch-delete`.
func newRegistryImageBatchDeleteCommand(state *cli.State) *cobra.Command {
	var body registry.BatchDeleteImageRequest
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "batch-delete <repository-id>",
		Short: "Delete images",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"digests": "digests",
				"tags":    "tags",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/repositories/{repository_id}/images/batch-delete")
			if err != nil {
				return err
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.BatchDeleteImage(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringSliceVar(&body.Digests, "digests", nil, "Removes the manifest AND every tag pointing at it")
	f.StringSliceVar(&body.Tags, "tags", nil, "Removes only the tag")
	return cmd
}

// newRegistryPullThroughCacheRuleCommand builds `basaltic registry pull-through-cache-rule`.
func newRegistryPullThroughCacheRuleCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pull-through-cache-rule",
		Short:   "Pull through cache rules",
		Aliases: []string{"pull-through-cache-rules"},
	}
	cmd.AddCommand(newRegistryPullThroughCacheRuleListCommand(state))
	cmd.AddCommand(newRegistryPullThroughCacheRuleGetCommand(state))
	cmd.AddCommand(newRegistryPullThroughCacheRuleCreateCommand(state))
	cmd.AddCommand(newRegistryPullThroughCacheRuleUpdateCommand(state))
	cmd.AddCommand(newRegistryPullThroughCacheRuleDeleteCommand(state))
	return cmd
}

// newRegistryPullThroughCacheRuleListCommand builds `basaltic registry pull-through-cache-rule list`.
func newRegistryPullThroughCacheRuleListCommand(state *cli.State) *cobra.Command {
	var params registry.DescribePullThroughCacheRulesParams
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List pull-through cache rules",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/pull-through-cache-rules")
			if err != nil {
				return err
			}
			page, err := c.DescribePullThroughCacheRules(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact pull-through-cache-rule CRN, intersected with name")
	f.StringVar(&params.Name, "name", "", "Exact, account-scoped immutable repository name or cache-rule prefix")
	return cmd
}

// newRegistryPullThroughCacheRuleGetCommand builds `basaltic registry pull-through-cache-rule get`.
func newRegistryPullThroughCacheRuleGetCommand(state *cli.State) *cobra.Command {
	var scope registry.DescribePullThroughCacheRulesParams
	cmd := &cobra.Command{
		Use:   "get <ref>",
		Short: "Get a pull-through cache rule",
		Args:  cobra.ExactArgs(1),
		Long:  "Get a pull-through cache rule.\n\n<ref> is its id, its CRN or its name. It is read by its syntax alone, the way the\nplatform reads it: a crn: value is a CRN, the 36-character UUID form is an\nid, anything else is a name. A miss under one reading is not retried\nunder another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/pull-through-cache-rules/{prefix}")
			if err != nil {
				return err
			}
			out, err := c.DescribePullThroughCacheRuleByReference(cmd.Context(), args[0], &scope)
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

// newRegistryPullThroughCacheRuleCreateCommand builds `basaltic registry pull-through-cache-rule create`.
func newRegistryPullThroughCacheRuleCreateCommand(state *cli.State) *cobra.Command {
	var body registry.CreatePullThroughCacheRuleRequest
	var bodyFile string
	var upstreamPrefixFlag string
	var upstreamSecretFlag string
	var upstreamUsernameFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a pull-through cache rule",
		Args:  cobra.ExactArgs(0),
		Long:  "Create a pull-through cache rule.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"prefix":            "prefix",
				"upstream_prefix":   "upstream-prefix",
				"upstream_registry": "upstream-registry",
				"upstream_secret":   "upstream-secret",
				"upstream_username": "upstream-username",
			}, []string{"prefix", "upstream_registry"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/pull-through-cache-rules")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("upstream-prefix") {
				body.UpstreamPrefix = &upstreamPrefixFlag
			}
			if cmd.Flags().Changed("upstream-secret") {
				body.UpstreamSecret = &upstreamSecretFlag
			}
			if cmd.Flags().Changed("upstream-username") {
				body.UpstreamUsername = &upstreamUsernameFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreatePullThroughCacheRule(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&body.Prefix, "prefix", "", "One or two lowercase repository-name components")
	f.StringVar(&upstreamPrefixFlag, "upstream-prefix", "", "Optional path prepended upstream")
	f.StringVar(&body.UpstreamRegistry, "upstream-registry", "", "A hostname with an optional port")
	f.StringVar(&upstreamSecretFlag, "upstream-secret", "", "The password or token")
	f.StringVar(&upstreamUsernameFlag, "upstream-username", "", "Optional")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newRegistryPullThroughCacheRuleUpdateCommand builds `basaltic registry pull-through-cache-rule update`.
func newRegistryPullThroughCacheRuleUpdateCommand(state *cli.State) *cobra.Command {
	var body registry.UpdatePullThroughCacheRuleRequest
	var bodyFile string
	var upstreamPrefixFlag string
	var upstreamRegistryFlag string
	var upstreamSecretFlag string
	var upstreamUsernameFlag string
	cmd := &cobra.Command{
		Use:   "update <prefix>",
		Short: "Update a pull-through cache rule",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"upstream_prefix":   "upstream-prefix",
				"upstream_registry": "upstream-registry",
				"upstream_secret":   "upstream-secret",
				"upstream_username": "upstream-username",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/pull-through-cache-rules/{prefix}")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("upstream-prefix") {
				body.UpstreamPrefix = &upstreamPrefixFlag
			}
			if cmd.Flags().Changed("upstream-registry") {
				body.UpstreamRegistry = &upstreamRegistryFlag
			}
			if cmd.Flags().Changed("upstream-secret") {
				body.UpstreamSecret = &upstreamSecretFlag
			}
			if cmd.Flags().Changed("upstream-username") {
				body.UpstreamUsername = &upstreamUsernameFlag
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.UpdatePullThroughCacheRule(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&upstreamPrefixFlag, "upstream-prefix", "", "Upstream prefix")
	f.StringVar(&upstreamRegistryFlag, "upstream-registry", "", "Upstream registry")
	f.StringVar(&upstreamSecretFlag, "upstream-secret", "", "Upstream secret")
	f.StringVar(&upstreamUsernameFlag, "upstream-username", "", "Upstream username")
	return cmd
}

// newRegistryPullThroughCacheRuleDeleteCommand builds `basaltic registry pull-through-cache-rule delete`.
func newRegistryPullThroughCacheRuleDeleteCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <prefix>",
		Short: "Delete a pull-through cache rule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/pull-through-cache-rules/{prefix}")
			if err != nil {
				return err
			}
			if err := c.DeletePullThroughCacheRule(cmd.Context(), args[0]); err != nil {
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

// newRegistryRepositoryCommand builds `basaltic registry repository`.
func newRegistryRepositoryCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "repository",
		Short:   "Repositories",
		Aliases: []string{"repositories"},
	}
	cmd.AddCommand(newRegistryRepositoryListCommand(state))
	cmd.AddCommand(newRegistryRepositoryGetCommand(state))
	cmd.AddCommand(newRegistryRepositoryCreateCommand(state))
	cmd.AddCommand(newRegistryRepositoryUpdateCommand(state))
	cmd.AddCommand(newRegistryRepositoryDeleteCommand(state))
	cmd.AddCommand(newRegistryRepositoryDeleteLifecyclePolicyCommand(state))
	cmd.AddCommand(newRegistryRepositoryDeletePolicyCommand(state))
	cmd.AddCommand(newRegistryRepositoryGetLifecyclePolicyCommand(state))
	cmd.AddCommand(newRegistryRepositoryGetPolicyCommand(state))
	cmd.AddCommand(newRegistryRepositoryListTagsCommand(state))
	cmd.AddCommand(newRegistryRepositorySetLifecyclePolicyCommand(state))
	cmd.AddCommand(newRegistryRepositorySetPolicyCommand(state))
	cmd.AddCommand(newRegistryRepositoryStartLifecyclePolicyCommand(state))
	return cmd
}

// newRegistryRepositoryListCommand builds `basaltic registry repository list`.
func newRegistryRepositoryListCommand(state *cli.State) *cobra.Command {
	var params registry.DescribeRepositoriesParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List repositories",
		Args:  cobra.ExactArgs(0),
		Long:  "List repositories.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/repositories")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.DescribeRepositoriesAll(cmd.Context(), &params))
			}
			page, err := c.DescribeRepositories(cmd.Context(), &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVar(&params.CRN, "crn", "", "Exact repository CRN, intersected with name before pagination")
	f.IntVar(&params.Limit, "limit", 0, "Maximum items to return")
	f.StringVar(&params.Marker, "marker", "", "Id of the last item on the previous page")
	f.StringVar(&params.Name, "name", "", "Exact, account-scoped immutable repository name or cache-rule prefix")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newRegistryRepositoryGetCommand builds `basaltic registry repository get`.
func newRegistryRepositoryGetCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <repository-id>",
		Short: "Describe a repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/repositories/{repository_id}")
			if err != nil {
				return err
			}
			out, err := c.DescribeRepository(cmd.Context(), args[0])
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

// newRegistryRepositoryCreateCommand builds `basaltic registry repository create`.
func newRegistryRepositoryCreateCommand(state *cli.State) *cobra.Command {
	var body registry.CreateRepositoryRequest
	var bodyFile string
	var deletionProtectionFlag bool
	var tagMutabilityFlag string
	var tagsFlag string
	var idempotencyKey string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a repository",
		Args:  cobra.ExactArgs(0),
		Long:  "Create a repository.\n\nPass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"deletion_protection": "deletion-protection",
				"name":                "name",
				"tag_mutability":      "tag-mutability",
				"tags":                "tags",
			}, []string{"name"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/repositories")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("deletion-protection") {
				body.DeletionProtection = &deletionProtectionFlag
			}
			if cmd.Flags().Changed("tag-mutability") {
				body.TagMutability = &tagMutabilityFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			if idempotencyKey != "" {
				reqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))
			}
			out, err := c.CreateRepository(cmd.Context(), &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.BoolVar(&deletionProtectionFlag, "deletion-protection", false, "Deletion protection")
	f.StringVar(&body.Name, "name", "", "Resource names must not start with the literal crn: prefix or be UUIDs (canonical, compact, braced, or urn:uuid: forms, in either case)")
	f.StringVar(&tagMutabilityFlag, "tag-mutability", "", "Tag mutability (one of: mutable, immutable)")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")
	return cmd
}

// newRegistryRepositoryUpdateCommand builds `basaltic registry repository update`.
func newRegistryRepositoryUpdateCommand(state *cli.State) *cobra.Command {
	var body registry.UpdateRepositoryRequest
	var bodyFile string
	var deletionProtectionFlag bool
	var tagMutabilityFlag string
	var tagsFlag string
	cmd := &cobra.Command{
		Use:   "update <repository-id>",
		Short: "Update repository settings",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"deletion_protection": "deletion-protection",
				"tag_mutability":      "tag-mutability",
				"tags":                "tags",
			}, []string{})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/repositories/{repository_id}")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("deletion-protection") {
				body.DeletionProtection = &deletionProtectionFlag
			}
			if cmd.Flags().Changed("tag-mutability") {
				body.TagMutability = &tagMutabilityFlag
			}
			if tagsFlag != "" {
				if err := json.Unmarshal([]byte(tagsFlag), &body.Tags); err != nil {
					return fmt.Errorf("--tags: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.UpdateRepository(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.BoolVar(&deletionProtectionFlag, "deletion-protection", false, "Deletion protection")
	f.StringVar(&tagMutabilityFlag, "tag-mutability", "", "Tag mutability (one of: mutable, immutable)")
	f.StringVar(&tagsFlag, "tags", "", "Tags (JSON)")
	return cmd
}

// newRegistryRepositoryDeleteCommand builds `basaltic registry repository delete`.
func newRegistryRepositoryDeleteCommand(state *cli.State) *cobra.Command {
	var params registry.DeleteRepositoryParams
	var forceFlag bool
	cmd := &cobra.Command{
		Use:   "delete <repository-id>",
		Short: "Delete a repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/repositories/{repository_id}")
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("force") {
				params.Force = &forceFlag
			}
			if err := c.DeleteRepository(cmd.Context(), args[0], &params); err != nil {
				return err
			}
			state.Printer().Done("Deleted.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	f.BoolVar(&forceFlag, "force", false, "Delete the repository even though it still contains images")
	return cmd
}

// newRegistryRepositoryDeleteLifecyclePolicyCommand builds `basaltic registry repository delete-lifecycle-policy`.
func newRegistryRepositoryDeleteLifecyclePolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-lifecycle-policy <repository-id>",
		Short: "Delete the lifecycle policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/repositories/{repository_id}/lifecycle-policy")
			if err != nil {
				return err
			}
			if err := c.DeleteLifecyclePolicy(cmd.Context(), args[0]); err != nil {
				return err
			}
			state.Printer().Done("Delete lifecycle policy requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newRegistryRepositoryDeletePolicyCommand builds `basaltic registry repository delete-policy`.
func newRegistryRepositoryDeletePolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-policy <repository-id>",
		Short: "Delete the repository policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/repositories/{repository_id}/policy")
			if err != nil {
				return err
			}
			if err := c.DeleteRepositoryPolicy(cmd.Context(), args[0]); err != nil {
				return err
			}
			state.Printer().Done("Delete policy requested.")
			return nil
		},
	}
	f := cmd.Flags()
	_ = f
	return cmd
}

// newRegistryRepositoryGetLifecyclePolicyCommand builds `basaltic registry repository get-lifecycle-policy`.
func newRegistryRepositoryGetLifecyclePolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-lifecycle-policy <repository-id>",
		Short: "Get the lifecycle policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/repositories/{repository_id}/lifecycle-policy")
			if err != nil {
				return err
			}
			out, err := c.GetLifecyclePolicy(cmd.Context(), args[0])
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

// newRegistryRepositoryGetPolicyCommand builds `basaltic registry repository get-policy`.
func newRegistryRepositoryGetPolicyCommand(state *cli.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-policy <repository-id>",
		Short: "Get the repository policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/repositories/{repository_id}/policy")
			if err != nil {
				return err
			}
			out, err := c.GetRepositoryPolicy(cmd.Context(), args[0])
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

// newRegistryRepositoryListTagsCommand builds `basaltic registry repository list-tags`.
func newRegistryRepositoryListTagsCommand(state *cli.State) *cobra.Command {
	var params registry.ListTagsParams
	var fetchAll bool
	cmd := &cobra.Command{
		Use:   "list-tags <repository-id>",
		Short: "List tags",
		Args:  cobra.ExactArgs(1),
		Long:  "List tags.\n\nReturns one page. Pass --all to walk every page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/repositories/{repository_id}/tags")
			if err != nil {
				return err
			}
			if fetchAll {
				return state.Printer().Iter(c.ListTagsAll(cmd.Context(), args[0], &params))
			}
			page, err := c.ListTags(cmd.Context(), args[0], &params)
			if err != nil {
				return err
			}
			return state.Printer().Page(page)
		},
	}
	f := cmd.Flags()
	_ = f
	f.IntVar(&params.Limit, "limit", 0, "Maximum items to return")
	f.StringVar(&params.Marker, "marker", "", "Name of the last tag on the previous page")
	f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")
	return cmd
}

// newRegistryRepositorySetLifecyclePolicyCommand builds `basaltic registry repository set-lifecycle-policy`.
func newRegistryRepositorySetLifecyclePolicyCommand(state *cli.State) *cobra.Command {
	var body registry.LifecyclePolicyRequest
	var bodyFile string
	var rulesFlag string
	cmd := &cobra.Command{
		Use:   "set-lifecycle-policy <repository-id>",
		Short: "Set the lifecycle policy",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"rules": "rules",
			}, []string{"rules"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/repositories/{repository_id}/lifecycle-policy")
			if err != nil {
				return err
			}
			if rulesFlag != "" {
				if err := json.Unmarshal([]byte(rulesFlag), &body.Rules); err != nil {
					return fmt.Errorf("--rules: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.PutLifecyclePolicy(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&rulesFlag, "rules", "", "Rules (JSON)")
	return cmd
}

// newRegistryRepositorySetPolicyCommand builds `basaltic registry repository set-policy`.
func newRegistryRepositorySetPolicyCommand(state *cli.State) *cobra.Command {
	var body registry.RepositoryPolicyRequest
	var bodyFile string
	var documentFlag string
	cmd := &cobra.Command{
		Use:   "set-policy <repository-id>",
		Short: "Set the repository policy",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"document": "document",
			}, []string{"document"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/repositories/{repository_id}/policy")
			if err != nil {
				return err
			}
			if documentFlag != "" {
				if err := json.Unmarshal([]byte(documentFlag), &body.Document); err != nil {
					return fmt.Errorf("--document: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.SetRepositoryPolicy(cmd.Context(), args[0], &body, reqOpts...)
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
	return cmd
}

// newRegistryRepositoryStartLifecyclePolicyCommand builds `basaltic registry repository start-lifecycle-policy`.
func newRegistryRepositoryStartLifecyclePolicyCommand(state *cli.State) *cobra.Command {
	var body registry.LifecyclePolicyRequest
	var bodyFile string
	var rulesFlag string
	cmd := &cobra.Command{
		Use:   "start-lifecycle-policy <repository-id>",
		Short: "Preview a lifecycle policy",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepareBody(cmd, bodyFile, &body, map[string]string{
				"rules": "rules",
			}, []string{"rules"})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := registryClient(state, "/v1/repositories/{repository_id}/lifecycle-policy/preview")
			if err != nil {
				return err
			}
			if rulesFlag != "" {
				if err := json.Unmarshal([]byte(rulesFlag), &body.Rules); err != nil {
					return fmt.Errorf("--rules: %w", err)
				}
			}
			var reqOpts []basaltic.RequestOption
			reqOpts = append(reqOpts, bodyRequestOptions(cmd)...)
			out, err := c.StartLifecyclePolicyPreview(cmd.Context(), args[0], &body, reqOpts...)
			if err != nil {
				return err
			}
			return state.Printer().Value(out)
		},
	}
	f := cmd.Flags()
	_ = f
	f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")
	f.StringVar(&rulesFlag, "rules", "", "Rules (JSON)")
	return cmd
}
