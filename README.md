# basaltic

The command-line interface for the [Basaltic](https://basaltic.sh) cloud
platform.

```bash
go install github.com/basaltic-sh/cli@latest
```

## Getting started

```bash
basaltic auth login --api-key ACCESS_KEY_ID:SECRET --set-region sa-saopaulo-1
basaltic compute instance list
```

`auth login` stores the credential in a profile and verifies it before
reporting success — a credential that is saved but refused fails later and
somewhere else, which is worse than not saving it.

## How commands are organised

```
basaltic <service> <resource> <verb> [flags]
```

```bash
basaltic compute instance list
basaltic compute instance start i-abc123
basaltic network vpc create --name prod --cidr-v4 10.0.0.0/16
basaltic storage volume list
basaltic loadbalancer listener create --load-balancer-id lb-1 --port 443
basaltic database cluster failover c-1
```

Positional arguments are the resource's identifiers, in the order the API's
path takes them. Everything else is a flag.

Longer service names have short aliases, and every resource accepts its
plural:

```bash
basaltic lb listener list        # loadbalancer
basaltic net vpc list            # network
basaltic db cluster list         # database
basaltic compute instances list  # same as `instance`
```

A service with a single resource of the same name drops the middle word:
`basaltic certificate list`, not `basaltic certificate certificate list`.

### Why this shape

It is the API's own structure. The alternative groupings were considered and
rejected: AWS's `aws ec2 describe-instances` mirrors API action names and
inherits their inconsistency, and kubectl's verb-first `get pods` works
because Kubernetes is uniform CRUD over one object model. This platform is
not — `start`, `resize`, `failover`, `rotate`, `reinstall` and `attach` are a
third of the surface, and verb-first would put the rarest word first and turn
the top level into a list of twenty-five verbs.

## Output

`--output text` (the default) prints tables and key/value blocks. `--output
json` and `--output yaml` print the API's own field names, for piping.

```bash
basaltic compute instance list -o json | jq -r '.[].id'
basaltic compute instance list --no-headers | awk '{print $1}'
```

## Paging

List commands return one page and say when there is more. `--all` walks every
page:

```bash
basaltic compute instance list --limit 50
basaltic compute instance list --all
```

## Creating things

Scalar fields are flags. Anything structured takes JSON, because a nested
object has no honest flat representation:

```bash
basaltic compute instance create \
  --name web-01 \
  --flavor-id f-1 \
  --image-id debian-13 \
  --networks '[{"subnet_id":"sn-1","assign_public_ip":true}]'
```

For anything more than that, `--from-file` takes the whole body as JSON or
YAML, and flags override individual fields of it:

```bash
basaltic compute instance create --from-file instance.yaml --name web-02
```

### Retrying a create

A create that times out may already have succeeded, so the CLI will not
repeat one on its own. Give it a key and it becomes safe to retry — the
platform returns the original outcome instead of building a second resource:

```bash
basaltic compute instance create --name web-01 --flavor-id f-1 \
  --idempotency-key "$(uuidgen)"
```

## Profiles

Configuration lives in `~/.config/basaltic/config.yaml`:

```yaml
default_profile: production
profiles:
  production:
    region: sa-saopaulo-1
    api_key: ACCESS_KEY_ID:SECRET
    account_id: acme
```

```bash
basaltic config list
basaltic config use production
basaltic config set region sa-saopaulo-1
basaltic --profile staging compute instance list
```

Cached access tokens are kept in a **separate** file,
`~/.config/basaltic/credentials.yaml`. `config.yaml` is the file people copy
between machines and check into dotfile repositories; a token must not travel
with it.

Environment variables override the profile: `BASALTIC_API_KEY`,
`BASALTIC_REGION`, `BASALTIC_ACCOUNT_ID`, `BASALTIC_PROFILE`. A flag overrides
both.

## Accounts

`--account-id` selects which account a command acts on. Without it, commands
act on the credential's own account, and another account's resources answer
"not found" rather than "not permitted" — an account you are not acting in is
one whose resources you cannot see.

## Generated code

`internal/generated` is built from the SDK's `api.json`, which describes the
SDK's Go surface: every method, its parameters and their types. The CLI does
not read the OpenAPI specifications and does not embed a copy of them.

```bash
make generate SDK=/path/to/sdk-go
```

Generated files are committed, so building the CLI needs nothing but this
repository. Do not edit them; change the SDK, regenerate, and commit the
result. `make check-generated` fails when the two have drifted.

## License

Apache 2.0. See [LICENSE](LICENSE).
