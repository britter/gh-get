#!/usr/bin/env bash
# release.sh — cut a release: finalize CHANGELOG, commit, tag, and push.
# The tag push triggers .github/workflows/release.yml, which builds and
# publishes the precompiled extension binaries.
#
# Usage: ./release.sh <version>    e.g. ./release.sh 2.6.0
set -euo pipefail

version="${1:-}"
[ -n "$version" ] || { echo "usage: $0 <version>  (e.g. 2.6.0)" >&2; exit 1; }
version="${version#v}"                       # accept 2.6.0 or v2.6.0
[[ "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || { echo "invalid version: $version (want X.Y.Z)" >&2; exit 1; }
tag="v$version"

cd "$(dirname "$0")"

# --- preflight guards ---
[ "$(git rev-parse --abbrev-ref HEAD)" = "main" ] || { echo "not on main" >&2; exit 1; }
git diff --quiet && git diff --cached --quiet || { echo "working tree not clean" >&2; exit 1; }
if git rev-parse -q --verify "refs/tags/$tag" >/dev/null; then echo "tag $tag already exists" >&2; exit 1; fi
grep -q '^## \[Unreleased\]$' CHANGELOG.md || { echo "no [Unreleased] section in CHANGELOG.md" >&2; exit 1; }

# Derive the repo URL and previous version from the existing [Unreleased] link,
# e.g. "[Unreleased]: https://github.com/britter/gh-get/compare/v2.5.0...HEAD".
base=$(sed -n 's|^\[Unreleased\]: \(.*\)/compare/.*|\1|p' CHANGELOG.md)
prev=$(sed -n 's|^\[Unreleased\]: .*/compare/\(v[0-9][0-9.]*\)\.\.\.HEAD$|\1|p' CHANGELOG.md)
[ -n "$base" ] && [ -n "$prev" ] || { echo "could not parse [Unreleased] link in CHANGELOG.md" >&2; exit 1; }

# --- finalize CHANGELOG: name the section, repoint [Unreleased], add the version link ---
today=$(date -u +%F)
sed -i \
  -e "s|^## \[Unreleased\]$|## [$version] - $today|" \
  -e "s|^\[Unreleased\]: .*|[Unreleased]: $base/compare/$tag...HEAD\n[$version]: $base/compare/$prev...$tag|" \
  CHANGELOG.md

# --- show the diff and confirm before anything is committed or pushed ---
git --no-pager diff CHANGELOG.md
echo
printf 'Release %s (previous %s). Does the diff look right? This commits, tags, and pushes. [y/N] ' "$tag" "$prev"
read -r reply
if [ "$reply" != "y" ] && [ "$reply" != "Y" ]; then
  git checkout -- CHANGELOG.md
  echo "aborted, CHANGELOG.md reverted"
  exit 1
fi

# --- commit, tag, push ---
git add CHANGELOG.md
git commit -m "chore: release $tag"
git tag -m "Release $tag" "$tag"
git push origin main
git push origin "$tag"

echo "Pushed $tag — the release workflow will build and publish the binaries."
