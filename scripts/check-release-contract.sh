#!/usr/bin/env bash
#
# Three places construct the release asset names: .goreleaser.yaml's
# name_template, AssetName() in internal/selfupdate, and install.sh. They have
# to agree, and nothing in an ordinary build compares them — a rename in one
# would break `basaltic upgrade` and `curl | sh` for everyone already
# installed, while every test still passed.
#
# This builds a real snapshot and checks the other two against what came out.
set -euo pipefail

cd "$(dirname "$0")/.."

command -v goreleaser >/dev/null 2>&1 || {
    echo "goreleaser is not installed: https://goreleaser.com/install/" >&2
    exit 1
}

echo "Building a snapshot..."
goreleaser release --snapshot --clean --skip=publish >/dev/null

VERSION=$(sed -n 's/.*"version"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' dist/metadata.json | head -n1)
[ -n "$VERSION" ] || { echo "could not read the snapshot version from dist/metadata.json" >&2; exit 1; }
echo "Snapshot version: $VERSION"

fail=0

# 1. Go's AssetName names a file goreleaser actually built, for every platform.
while read -r name; do
    if [ ! -f "dist/$name" ]; then
        echo "MISMATCH: selfupdate.AssetName produced $name, which goreleaser did not build" >&2
        fail=1
    fi
done < <(go run ./scripts/assetname "$VERSION")

# 2. install.sh builds the same name. Its construction is
#    ${BINARY}_${VERSION}_${OS}_${ARCH}.tar.gz, reproduced here.
for goos in linux darwin; do
    for goarch in amd64 arm64; do
        want="basaltic_${VERSION}_${goos}_${goarch}.tar.gz"
        if [ ! -f "dist/$want" ]; then
            echo "MISSING: dist/$want — install.sh would look for this and get a 404" >&2
            fail=1
        fi
    done
done

# 3. checksums.txt is in the format install.sh greps for: "<hex>  <filename>".
archive="basaltic_${VERSION}_linux_amd64.tar.gz"
if ! grep -q " ${archive}\$" dist/checksums.txt; then
    echo "MISMATCH: install.sh's checksum lookup finds nothing for $archive" >&2
    echo "          checksums.txt looks like:" >&2
    head -2 dist/checksums.txt >&2
    fail=1
fi

# 4. The binary sits at the archive root under the name both extractors expect.
if ! tar -tzf "dist/$archive" | grep -qx basaltic; then
    echo "MISMATCH: $archive does not contain a top-level 'basaltic'" >&2
    tar -tzf "dist/$archive" >&2
    fail=1
fi

if [ "$fail" -ne 0 ]; then
    echo >&2
    echo "The release artifacts no longer match what install.sh and" >&2
    echo "internal/selfupdate expect. Fix all three together." >&2
    exit 1
fi

echo "Release contract holds: goreleaser, install.sh and selfupdate agree."
