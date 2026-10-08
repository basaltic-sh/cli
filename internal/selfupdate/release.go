// Package selfupdate finds and installs newer releases of the CLI.
//
// Releases are GitHub releases built by goreleaser. The asset naming here is a
// contract with .goreleaser.yaml and with install.sh: all three construct the
// same file names, and changing one without the others breaks upgrades
// silently for everyone who has already installed.
package selfupdate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

// Repo is the release source.
const (
	Owner = "basaltic-sh"
	Repo  = "cli"
)

// apiTimeout bounds a release lookup. Short: this runs on the way to doing
// what the user actually asked for.
const apiTimeout = 10 * time.Second

// Release is a published version.
type Release struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	HTMLURL     string    `json:"html_url"`
	Draft       bool      `json:"draft"`
	Prerelease  bool      `json:"prerelease"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []Asset   `json:"assets"`
}

// Asset is one downloadable file on a release.
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// Version is the release's version without the leading "v".
func (r *Release) Version() string { return strings.TrimPrefix(r.TagName, "v") }

// AssetFor returns the archive for a platform, and the checksums file.
func (r *Release) AssetFor(goos, goarch string) (archive, checksums *Asset, err error) {
	want := AssetName(r.Version(), goos, goarch)
	for i := range r.Assets {
		switch r.Assets[i].Name {
		case want:
			archive = &r.Assets[i]
		case "checksums.txt":
			checksums = &r.Assets[i]
		}
	}
	if archive == nil {
		return nil, nil, fmt.Errorf(
			"release %s has no build for %s/%s (looked for %s)", r.TagName, goos, goarch, want)
	}
	if checksums == nil {
		// Refusing here rather than installing unverified bytes: an upgrade
		// command that silently skips verification is worse than one that
		// fails.
		return nil, nil, fmt.Errorf("release %s publishes no checksums.txt, so the download cannot be verified", r.TagName)
	}
	return archive, checksums, nil
}

// AssetName builds the archive file name for a platform.
//
// Must match `archives.name_template` in .goreleaser.yaml and the equivalent
// construction in install.sh.
func AssetName(version, goos, goarch string) string {
	ext := "tar.gz"
	if goos == "windows" {
		ext = "zip"
	}
	return fmt.Sprintf("basaltic_%s_%s_%s.%s", version, goos, goarch, ext)
}

// Latest returns the most recent published release.
func Latest(ctx context.Context, hc *http.Client) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", Owner, Repo)
	return fetchRelease(ctx, hc, url)
}

// AtVersion returns a specific release by tag.
func AtVersion(ctx context.Context, hc *http.Client, version string) (*Release, error) {
	tag := version
	if !strings.HasPrefix(tag, "v") {
		tag = "v" + tag
	}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", Owner, Repo, tag)
	return fetchRelease(ctx, hc, url)
}

func fetchRelease(ctx context.Context, hc *http.Client, url string) (*Release, error) {
	if hc == nil {
		hc = &http.Client{Timeout: apiTimeout}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("checking for releases: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20))

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		// GitHub answers 404 rather than 403 for a private repository, so
		// "not found" is the expected reply while this one is private and
		// says nothing about whether a release exists.
		return nil, fmt.Errorf(
			"no releases found at github.com/%s/%s.\n"+
				"If the repository is still private, GitHub reports it as not found and\n"+
				"upgrading this way will not work yet — install from a build instead", Owner, Repo)
	case http.StatusForbidden:
		if resp.Header.Get("X-RateLimit-Remaining") == "0" {
			return nil, fmt.Errorf("GitHub's API rate limit is exhausted for this IP; try again later")
		}
		return nil, fmt.Errorf("GitHub refused the release lookup (http 403)")
	default:
		return nil, fmt.Errorf("checking for releases: http %d", resp.StatusCode)
	}

	var rel Release
	if err := json.Unmarshal(body, &rel); err != nil {
		return nil, fmt.Errorf("checking for releases: malformed response: %w", err)
	}
	if rel.TagName == "" {
		return nil, fmt.Errorf("checking for releases: the response names no version")
	}
	return &rel, nil
}

// IsNewer reports whether candidate is a later version than current.
//
// A development build — anything that is not valid semver, which includes the
// "dev" default — is never considered older, so a locally built binary is not
// nagged to "upgrade" to a release it may well be ahead of.
func IsNewer(current, candidate string) bool {
	cur := canonical(current)
	cand := canonical(candidate)
	if cur == "" || cand == "" {
		return false
	}
	return semver.Compare(cand, cur) > 0
}

func canonical(v string) string {
	if v == "" {
		return ""
	}
	if !strings.HasPrefix(v, "v") {
		v = "v" + v
	}
	if !semver.IsValid(v) {
		return ""
	}
	return semver.Canonical(v)
}
