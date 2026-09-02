// The command generator is a module of its own, so the CLI it emits does not
// carry its dependencies. It reads the SDK's api.json manifest and needs
// nothing else — in particular not a checkout of the OpenAPI specification
// repository, which is private.
module github.com/basaltic-sh/cli/internal/gen

go 1.23.0
