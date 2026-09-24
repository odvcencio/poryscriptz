#!/usr/bin/env bash
set -euo pipefail

project_root=$(cd "$(dirname "$0")/.." && pwd)
package_version=$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1], encoding="utf-8"))["version"])' "$project_root/editors/vscode/package.json")
version=${1:-$package_version}
if [[ "$version" != "$package_version" ]]; then
    echo "release version $version does not match the VS Code package version $package_version" >&2
    exit 2
fi
release_dir="$project_root/dist"
bin_dir=$(mktemp -d)
trap 'rm -rf "$bin_dir"' EXIT
mkdir -p "$release_dir"

cd "$project_root"
GOWORK=off go test ./...

for platform in windows/amd64 windows/arm64 darwin/arm64 darwin/amd64 linux/amd64 linux/arm64; do
    os=${platform%/*}
    arch=${platform#*/}
    mkdir -p "$bin_dir/$os-$arch"
    exe=poryz
    if [[ "$os" == windows ]]; then exe=poryz.exe; fi
    GOWORK=off CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" \
        go build -trimpath -ldflags '-s -w' -o "$bin_dir/$os-$arch/$exe" ./cmd/poryscriptz
    if [[ "$os" == windows ]]; then
        (cd "$bin_dir/$os-$arch" && python3 -m zipfile -c "$release_dir/poryz_${version}_${os}_${arch}.zip" "$exe")
    else
        tar -C "$bin_dir/$os-$arch" -czf "$release_dir/poryz_${version}_${os}_${arch}.tar.gz" "$exe"
    fi
done

(cd editors/vscode && npm ci --no-audit --no-fund && npm run package)
cp "editors/vscode/dist/poryz-${version}.vsix" "$release_dir/"
(cd "$release_dir" && sha256sum poryz_${version}_*.zip poryz_${version}_*.tar.gz poryz-${version}.vsix > checksums.txt)
echo "Release artifacts are in $release_dir"
