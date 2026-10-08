package selfupdate

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// maxArchiveSize bounds a download. The real archives are a few megabytes; a
// cap stops a wrong URL filling a disk.
const maxArchiveSize = 200 << 20

// Install replaces the running binary with the one from rel.
//
// The order matters and is the whole safety argument: download, verify the
// checksum, extract, write beside the target, then rename over it. Nothing
// touches the installed path until the bytes have been checked, and the final
// step is a rename, which is atomic — an interrupted upgrade leaves either the
// old binary or the new one, never half of either.
func Install(ctx context.Context, hc *http.Client, rel *Release, execPath string) error {
	archive, checksums, err := rel.AssetFor(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return err
	}

	want, err := checksumFor(ctx, hc, checksums, archive.Name)
	if err != nil {
		return err
	}

	dir := filepath.Dir(execPath)
	tmp, err := os.CreateTemp(dir, ".basaltic-download-*")
	if err != nil {
		return fmt.Errorf("cannot stage the download next to %s: %w", execPath, err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	sum, err := download(ctx, hc, archive.BrowserDownloadURL, tmp)
	tmp.Close()
	if err != nil {
		return err
	}
	if sum != want {
		return fmt.Errorf(
			"the download does not match the published checksum.\n"+
				"  expected %s\n  got      %s\n"+
				"Nothing was installed. This is worth reporting rather than retrying", want, sum)
	}

	binary, err := extractBinary(tmpName, archive.Name)
	if err != nil {
		return err
	}
	defer os.Remove(binary)

	return replace(binary, execPath)
}

// checksumFor downloads checksums.txt and returns the hash for one file.
func checksumFor(ctx context.Context, hc *http.Client, checksums *Asset, name string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, checksums.BrowserDownloadURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := hc.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetching checksums: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetching checksums: http %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}
	// Lines are "<hex>  <filename>", the sha256sum format goreleaser writes.
	for _, line := range strings.Split(string(body), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[1] == name {
			return strings.ToLower(fields[0]), nil
		}
	}
	return "", fmt.Errorf("checksums.txt does not list %s", name)
}

// download copies a URL to w, returning the hex sha256 of what was written.
func download(ctx context.Context, hc *http.Client, url string, w io.Writer) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := hc.Do(req)
	if err != nil {
		return "", fmt.Errorf("downloading %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("downloading %s: http %d", url, resp.StatusCode)
	}
	h := sha256.New()
	if _, err := io.Copy(io.MultiWriter(w, h), io.LimitReader(resp.Body, maxArchiveSize)); err != nil {
		return "", fmt.Errorf("downloading %s: %w", url, err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// extractBinary pulls the basaltic executable out of the archive and returns
// the path it was written to.
func extractBinary(archivePath, archiveName string) (string, error) {
	want := "basaltic"
	if runtime.GOOS == "windows" {
		want = "basaltic.exe"
	}
	if strings.HasSuffix(archiveName, ".zip") {
		return extractFromZip(archivePath, want)
	}
	return extractFromTarGz(archivePath, want)
}

func extractFromTarGz(path, want string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return "", fmt.Errorf("the download is not a gzip archive: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", fmt.Errorf("reading the archive: %w", err)
		}
		// Match on the base name only, and never join an archive-supplied
		// path: a crafted archive must not be able to write outside the
		// staging directory.
		if filepath.Base(hdr.Name) != want || hdr.Typeflag != tar.TypeReg {
			continue
		}
		return writeStaged(tr, filepath.Dir(path))
	}
	return "", fmt.Errorf("the archive does not contain %s", want)
}

func extractFromZip(path, want string) (string, error) {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return "", fmt.Errorf("the download is not a zip archive: %w", err)
	}
	defer zr.Close()
	for _, f := range zr.File {
		if filepath.Base(f.Name) != want || f.FileInfo().IsDir() {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()
		return writeStaged(rc, filepath.Dir(path))
	}
	return "", fmt.Errorf("the archive does not contain %s", want)
}

func writeStaged(r io.Reader, dir string) (string, error) {
	out, err := os.CreateTemp(dir, ".basaltic-new-*")
	if err != nil {
		return "", err
	}
	defer out.Close()
	if _, err := io.Copy(out, io.LimitReader(r, maxArchiveSize)); err != nil {
		os.Remove(out.Name())
		return "", err
	}
	if err := out.Chmod(0o755); err != nil {
		os.Remove(out.Name())
		return "", err
	}
	return out.Name(), nil
}

// replace swaps the new binary into place.
//
// On Unix a rename over a running executable is fine — the kernel holds the
// old inode open for processes already using it. Windows refuses, so the old
// file is moved aside first and left for the next run to clean up.
func replace(newPath, execPath string) error {
	if runtime.GOOS == "windows" {
		old := execPath + ".old"
		_ = os.Remove(old)
		if err := os.Rename(execPath, old); err != nil {
			return fmt.Errorf("cannot move the running binary aside: %w", err)
		}
		if err := os.Rename(newPath, execPath); err != nil {
			// Put it back rather than leaving nothing installed.
			_ = os.Rename(old, execPath)
			return fmt.Errorf("installing the new binary: %w", err)
		}
		return nil
	}
	if err := os.Rename(newPath, execPath); err != nil {
		return fmt.Errorf("installing to %s: %w", execPath, err)
	}
	return nil
}

// ResolvePath returns the real path of the running binary, following symlinks.
//
// Following matters: a CLI on PATH is often a symlink into a version
// directory, and replacing the link rather than its target would leave the
// old binary installed and the link no longer pointing at a managed file.
func ResolvePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("cannot locate the running binary: %w", err)
	}
	resolved, err := filepath.EvalSymlinks(exe)
	if err != nil {
		return exe, nil
	}
	return resolved, nil
}

// CheckWritable reports why an in-place upgrade cannot proceed, or nil.
//
// The two cases worth catching before downloading anything: a binary owned by
// a package manager, where upgrading behind its back leaves it convinced the
// old version is installed; and a directory the user cannot write, where the
// only useful answer is to say which one.
func CheckWritable(execPath string) error {
	if mgr := packageManager(execPath); mgr != "" {
		return fmt.Errorf(
			"this binary was installed by %s, which owns the file at %s.\n"+
				"Upgrade through %s instead — replacing it here would leave %s\n"+
				"believing the old version is still installed", mgr, execPath, mgr, mgr)
	}
	dir := filepath.Dir(execPath)
	probe, err := os.CreateTemp(dir, ".basaltic-write-probe-*")
	if err != nil {
		return fmt.Errorf(
			"cannot write to %s, where the binary lives.\n"+
				"Either re-run with the necessary permissions, or install somewhere you own:\n"+
				"    curl -fsSL https://get.basaltic.sh/cli | BASALTIC_INSTALL_DIR=$HOME/.local/bin sh", dir)
	}
	probe.Close()
	os.Remove(probe.Name())
	return nil
}

// packageManager names the manager that owns a path, or "".
func packageManager(execPath string) string {
	switch {
	case strings.Contains(execPath, "/Cellar/"), strings.Contains(execPath, "/homebrew/"):
		return "Homebrew"
	case strings.HasPrefix(execPath, "/nix/store/"):
		return "Nix"
	case strings.HasPrefix(execPath, "/snap/"):
		return "snap"
	case strings.HasPrefix(execPath, "/usr/bin/"), strings.HasPrefix(execPath, "/usr/lib/"):
		// The system package directories. /usr/local is deliberately not
		// here: that is where install.sh puts things.
		return "the system package manager"
	}
	return ""
}
