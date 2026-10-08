package auth

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/basaltic-sh/cli/internal/config"
	_ "github.com/basaltic-sh/sdk-go/catalog"
)

func TestPublicResolutionDoesNotRequireCredentials(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	t.Setenv(config.EnvConfigFile, path)
	t.Setenv("BASALTIC_API_KEY", "")
	t.Setenv("BASALTIC_ACCESS_TOKEN", "")
	if err := os.WriteFile(path, []byte("current_profile: default\nprofiles:\n  default:\n    domain: example.test\n    account_id: private-account\n"), 0600); err != nil {
		t.Fatal(err)
	}
	resolved, err := ResolvePublic(Options{})
	if err != nil {
		t.Fatal(err)
	}
	if resolved.Config.AccountID != "" {
		t.Fatal("public discovery retained account scope")
	}
	if _, err := resolved.Config.TokenSource.Token(context.Background()); err == nil {
		t.Fatal("public config can authenticate")
	}
	endpoint, err := resolved.Config.EndpointResolver.ResolveEndpoint("catalog", "")
	if err != nil {
		t.Fatal(err)
	}
	if endpoint != "https://catalog.example.test" {
		t.Fatalf("profile endpoint lost: %s", endpoint)
	}
}
