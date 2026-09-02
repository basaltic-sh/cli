// Command basaltic is the command-line interface for the Basaltic cloud
// platform.
package main

import (
	"os"

	"github.com/basaltic-sh/cli/internal/cli"

	// Registering the command tree is a side effect of importing these: the
	// generated service commands and the CLI's own.
	_ "github.com/basaltic-sh/cli/internal/commands"
	_ "github.com/basaltic-sh/cli/internal/generated"
)

func main() {
	os.Exit(cli.Execute())
}
