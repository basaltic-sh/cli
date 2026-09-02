package auth

// Version is the CLI's version, reported by `basaltic version` and sent in
// the User-Agent. Set at build time with:
//
//	go build -ldflags "-X github.com/basaltic-sh/cli/internal/auth.Version=1.2.3"
var Version = "dev"
