#!/bin/sh
#
# Install the basaltic CLI.
#
#     curl -fsSL https://get.basaltic.sh/cli | sh
#
# Options, as environment variables:
#
#     BASALTIC_VERSION=1.2.3      install this version instead of the latest
#     BASALTIC_INSTALL_DIR=DIR    install here instead of the default
#
# The default is /usr/local/bin when it is writable, and ~/.local/bin
# otherwise. Nothing is installed until its checksum matches the one published
# with the release.
#
# This script is piped into a shell, so it is written to be read first:
# POSIX sh, no dependencies beyond curl or wget, and it says what it is doing.

set -eu

OWNER="basaltic-sh"
REPO="cli"
BINARY="basaltic"

# --- output -----------------------------------------------------------------

# Everything informational goes to stderr, so that piping the script's own
# output somewhere never mixes with what it says.
say() { printf '%s\n' "$*" >&2; }
die() { printf 'error: %s\n' "$*" >&2; exit 1; }

# --- platform ---------------------------------------------------------------

detect_os() {
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    case "$os" in
        linux)  echo linux ;;
        darwin) echo darwin ;;
        msys*|mingw*|cygwin*)
            die "Windows is not supported by this script. Download the .zip from
       https://github.com/$OWNER/$REPO/releases and put basaltic.exe on your PATH." ;;
        *) die "unsupported operating system: $os" ;;
    esac
}

detect_arch() {
    arch=$(uname -m)
    case "$arch" in
        x86_64|amd64)  echo amd64 ;;
        aarch64|arm64) echo arm64 ;;
        *) die "unsupported architecture: $arch (builds exist for amd64 and arm64)" ;;
    esac
}

# --- fetching ---------------------------------------------------------------

# One downloader, chosen once. --fail matters: without it curl writes an error
# page to the output file and exits 0, and the failure surfaces later as a
# corrupt archive.
if command -v curl >/dev/null 2>&1; then
    fetch()    { curl -fsSL "$1"; }
    download() { curl -fsSL -o "$2" "$1"; }
elif command -v wget >/dev/null 2>&1; then
    fetch()    { wget -qO- "$1"; }
    download() { wget -qO "$2" "$1"; }
else
    die "neither curl nor wget is available"
fi

latest_version() {
    api="https://api.github.com/repos/$OWNER/$REPO/releases/latest"
    body=$(fetch "$api" 2>/dev/null) || die "cannot reach the GitHub release API.
       If $OWNER/$REPO is still private, GitHub reports it as not found and this
       script cannot install from it yet."
    # Deliberately not jq: this script must run on a machine with nothing on it.
    version=$(printf '%s' "$body" | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"v\{0,1\}\([^"]*\)".*/\1/p' | head -n1)
    [ -n "$version" ] || die "could not read a version from the release API"
    printf '%s' "$version"
}

# --- checksums --------------------------------------------------------------

# sha256sum on Linux, shasum on macOS. Refusing to continue without one is
# deliberate: installing an unverified binary is the thing this script exists
# to avoid.
if command -v sha256sum >/dev/null 2>&1; then
    sha256() { sha256sum "$1" | cut -d' ' -f1; }
elif command -v shasum >/dev/null 2>&1; then
    sha256() { shasum -a 256 "$1" | cut -d' ' -f1; }
else
    die "no sha256sum or shasum available, so the download cannot be verified"
fi

# --- install location -------------------------------------------------------

choose_dir() {
    if [ -n "${BASALTIC_INSTALL_DIR:-}" ]; then
        printf '%s' "$BASALTIC_INSTALL_DIR"
        return
    fi
    if [ -w /usr/local/bin ] 2>/dev/null; then
        printf '%s' /usr/local/bin
        return
    fi
    printf '%s' "$HOME/.local/bin"
}

on_path() {
    case ":${PATH}:" in
        *":$1:"*) return 0 ;;
        *) return 1 ;;
    esac
}

# --- main -------------------------------------------------------------------

OS=$(detect_os)
ARCH=$(detect_arch)

VERSION="${BASALTIC_VERSION:-}"
if [ -z "$VERSION" ]; then
    say "Finding the latest release..."
    VERSION=$(latest_version)
fi
VERSION="${VERSION#v}"

ARCHIVE="${BINARY}_${VERSION}_${OS}_${ARCH}.tar.gz"
BASE="https://github.com/$OWNER/$REPO/releases/download/v$VERSION"

TMP=$(mktemp -d 2>/dev/null || mktemp -d -t basaltic)
# Cleared on any exit, including a failure part-way through: a half-downloaded
# archive left in /tmp is confusing the next time someone looks.
trap 'rm -rf "$TMP"' EXIT INT TERM

say "Downloading basaltic $VERSION for $OS/$ARCH..."
download "$BASE/$ARCHIVE" "$TMP/$ARCHIVE" || die "could not download $BASE/$ARCHIVE"
download "$BASE/checksums.txt" "$TMP/checksums.txt" || die "could not download the checksums"

EXPECTED=$(grep " $ARCHIVE\$" "$TMP/checksums.txt" | cut -d' ' -f1 || true)
[ -n "$EXPECTED" ] || die "checksums.txt does not list $ARCHIVE"

ACTUAL=$(sha256 "$TMP/$ARCHIVE")
if [ "$EXPECTED" != "$ACTUAL" ]; then
    die "checksum mismatch for $ARCHIVE
       expected $EXPECTED
       got      $ACTUAL
       Nothing was installed. This is worth reporting rather than retrying."
fi
say "Checksum verified."

tar -xzf "$TMP/$ARCHIVE" -C "$TMP" || die "could not extract $ARCHIVE"
[ -f "$TMP/$BINARY" ] || die "the archive does not contain $BINARY"

DIR=$(choose_dir)
mkdir -p "$DIR" || die "cannot create $DIR"
if [ ! -w "$DIR" ]; then
    die "$DIR is not writable.
       Either re-run with elevated permissions, or choose somewhere you own:
       curl -fsSL https://get.basaltic.sh/cli | BASALTIC_INSTALL_DIR=\$HOME/.local/bin sh"
fi

# Move into place through a temporary name in the same directory, so an
# interrupted install cannot leave a partial binary where a working one was.
chmod 0755 "$TMP/$BINARY"
mv "$TMP/$BINARY" "$DIR/$BINARY.new"
mv "$DIR/$BINARY.new" "$DIR/$BINARY"

say ""
say "Installed basaltic $VERSION to $DIR/$BINARY"

if ! on_path "$DIR"; then
    say ""
    say "$DIR is not on your PATH. Add it:"
    say ""
    say "    export PATH=\"$DIR:\$PATH\""
    say ""
    say "and add that line to your shell profile to make it permanent."
else
    say ""
    say "Next: basaltic auth login --api-key ACCESS_KEY_ID:SECRET"
fi
