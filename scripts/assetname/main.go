// Command assetname prints the archive names the CLI's updater expects, so a
// shell script can compare them against what goreleaser actually built.
package main

import (
	"fmt"
	"os"

	"github.com/basaltic-sh/cli/internal/selfupdate"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: assetname <version>")
		os.Exit(2)
	}
	version := os.Args[1]
	for _, p := range []struct{ os, arch string }{
		{"linux", "amd64"}, {"linux", "arm64"},
		{"darwin", "amd64"}, {"darwin", "arm64"},
		{"windows", "amd64"},
	} {
		fmt.Println(selfupdate.AssetName(version, p.os, p.arch))
	}
}
