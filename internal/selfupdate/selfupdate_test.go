package selfupdate

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// The asset name is a contract between three places: .goreleaser.yaml's
// name_template, this package, and install.sh. They are checked against a
// real snapshot build in CI; this pins the Go side.
func TestAssetName(t *testing.T) {
	tests := []struct {
		version, goos, goarch, want string
	}{
		{"1.2.3", "linux", "amd64", "basaltic_1.2.3_linux_amd64.tar.gz"},
		{"1.2.3", "darwin", "arm64", "basaltic_1.2.3_darwin_arm64.tar.gz"},
		{"0.0.1-next", "linux", "arm64", "basaltic_0.0.1-next_linux_arm64.tar.gz"},
		// Windows ships a zip, which the extractor branches on.
		{"1.2.3", "windows", "amd64", "basaltic_1.2.3_windows_amd64.zip"},
	}
	for _, tc := range tests {
		if got := AssetName(tc.version, tc.goos, tc.goarch); got != tc.want {
			t.Errorf("AssetName(%q, %q, %q) = %q, want %q", tc.version, tc.goos, tc.goarch, got, tc.want)
		}
	}
}

func TestIsNewer(t *testing.T) {
	tests := []struct {
		current, candidate string
		want               bool
	}{
		{"1.0.0", "1.0.1", true},
		{"1.0.0", "1.1.0", true},
		{"1.0.0", "2.0.0", true},
		{"1.0.1", "1.0.0", false},
		{"1.0.0", "1.0.0", false},
		{"v1.0.0", "v1.0.1", true},
		// A prerelease is older than the release it precedes.
		{"1.0.0-rc1", "1.0.0", true},
		{"1.0.0", "1.0.0-rc1", false},
		// A development build is never told it is out of date: it may well be
		// ahead of the latest release.
		{"dev", "1.0.0", false},
		{"1.0.0", "dev", false},
		{"", "1.0.0", false},
		{"0.0.1-next", "1.0.0", true},
	}
	for _, tc := range tests {
		if got := IsNewer(tc.current, tc.candidate); got != tc.want {
			t.Errorf("IsNewer(%q, %q) = %v, want %v", tc.current, tc.candidate, got, tc.want)
		}
	}
}

func TestAssetForRejectsAReleaseWithoutChecksums(t *testing.T) {
	// Installing unverified bytes is the thing this package exists to avoid,
	// so a release with no checksums.txt must fail rather than proceed.
	rel := &Release{
		TagName: "v1.0.0",
		Assets:  []Asset{{Name: AssetName("1.0.0", "linux", "amd64")}},
	}
	if _, _, err := rel.AssetFor("linux", "amd64"); err == nil {
		t.Fatal("AssetFor succeeded with no checksums.txt")
	} else if !strings.Contains(err.Error(), "checksums") {
		t.Errorf("error does not mention checksums: %v", err)
	}
}

func TestAssetForRejectsAnUnbuiltPlatform(t *testing.T) {
	rel := &Release{
		TagName: "v1.0.0",
		Assets:  []Asset{{Name: "checksums.txt"}},
	}
	_, _, err := rel.AssetFor("plan9", "mips")
	if err == nil {
		t.Fatal("AssetFor succeeded for a platform with no build")
	}
	if !strings.Contains(err.Error(), "plan9") {
		t.Errorf("error does not name the platform: %v", err)
	}
}

// A private repository answers 404, which is the state this repo is in today.
// The message has to say so, or it reads as "there are no releases".
func TestPrivateRepositoryExplainsItself(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := fetchRelease(context.Background(), srv.Client(), srv.URL)
	if err == nil {
		t.Fatal("a 404 was not reported as an error")
	}
	if !strings.Contains(err.Error(), "private") {
		t.Errorf("the 404 message does not mention the private-repository case: %v", err)
	}
}

// The whole point of the package: bytes that do not match the published
// checksum must never reach the install path.
func TestInstallRefusesAMismatchedChecksum(t *testing.T) {
	dir := t.TempDir()
	archive, _ := buildArchive(t, dir, "the real binary")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "checksums.txt"):
			// A hash of something else entirely.
			fmt.Fprintf(w, "%s  %s\n", strings.Repeat("0", 64), filepath.Base(archive))
		default:
			http.ServeFile(w, r, archive)
		}
	}))
	defer srv.Close()

	target := filepath.Join(dir, "basaltic")
	if err := os.WriteFile(target, []byte("original"), 0o755); err != nil {
		t.Fatal(err)
	}

	rel := releaseServing(srv.URL, filepath.Base(archive))
	err := Install(context.Background(), srv.Client(), rel, target)
	if err == nil {
		t.Fatal("Install accepted a mismatched checksum")
	}
	if !strings.Contains(err.Error(), "checksum") {
		t.Errorf("error does not mention the checksum: %v", err)
	}
	// And the installed binary is untouched.
	got, _ := os.ReadFile(target)
	if string(got) != "original" {
		t.Errorf("the existing binary was modified despite the failure: %q", got)
	}
}

func TestInstallReplacesTheBinary(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the replace path differs on Windows and is exercised there")
	}
	dir := t.TempDir()
	archive, sum := buildArchive(t, dir, "the new binary")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "checksums.txt") {
			fmt.Fprintf(w, "%s  %s\n", sum, filepath.Base(archive))
			return
		}
		http.ServeFile(w, r, archive)
	}))
	defer srv.Close()

	target := filepath.Join(dir, "basaltic")
	if err := os.WriteFile(target, []byte("old binary"), 0o755); err != nil {
		t.Fatal(err)
	}

	rel := releaseServing(srv.URL, filepath.Base(archive))
	if err := Install(context.Background(), srv.Client(), rel, target); err != nil {
		t.Fatalf("Install: %v", err)
	}
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "the new binary" {
		t.Errorf("binary is %q, want the new one", got)
	}
	info, err := os.Stat(target)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0o111 == 0 {
		t.Errorf("the installed binary is not executable: %v", info.Mode())
	}
	// Nothing left staged in the install directory.
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".basaltic-") {
			t.Errorf("left a staging file behind: %s", e.Name())
		}
	}
}

func TestCheckWritableNamesThePackageManager(t *testing.T) {
	for path, want := range map[string]string{
		"/opt/homebrew/Cellar/basaltic/1.0.0/bin/basaltic": "Homebrew",
		"/nix/store/abc-basaltic/bin/basaltic":             "Nix",
		"/snap/basaltic/current/bin/basaltic":              "snap",
		"/usr/bin/basaltic":                                "the system package manager",
	} {
		err := CheckWritable(path)
		if err == nil {
			t.Errorf("CheckWritable(%q) allowed an upgrade in place", path)
			continue
		}
		if !strings.Contains(err.Error(), want) {
			t.Errorf("CheckWritable(%q) = %v, want it to name %q", path, err, want)
		}
	}
	// /usr/local/bin is where install.sh puts things, so it must NOT be
	// treated as package-managed.
	if mgr := packageManager("/usr/local/bin/basaltic"); mgr != "" {
		t.Errorf("/usr/local/bin was treated as owned by %q", mgr)
	}
}

// buildArchive writes a tar.gz containing a basaltic binary with the given
// contents, and returns its path and sha256.
func buildArchive(t *testing.T, dir, content string) (string, string) {
	t.Helper()
	path := filepath.Join(dir, AssetName("9.9.9", runtime.GOOS, runtime.GOARCH))
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)
	if err := tw.WriteHeader(&tar.Header{
		Name: "basaltic", Mode: 0o755, Size: int64(len(content)), Typeflag: tar.TypeReg,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tw.Close()
	gz.Close()
	f.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	h := sha256.Sum256(data)
	return path, hex.EncodeToString(h[:])
}

func releaseServing(base, archiveName string) *Release {
	return &Release{
		TagName: "v9.9.9",
		Assets: []Asset{
			{Name: archiveName, BrowserDownloadURL: base + "/" + archiveName},
			{Name: "checksums.txt", BrowserDownloadURL: base + "/checksums.txt"},
		},
	}
}
