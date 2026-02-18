#!/usr/bin/env bash

set -euxo pipefail
cp -r ../README.md ../dist/* .

# Use VERSION env var if set (CI), otherwise fetch from GitHub API (local)
if [ -z "${VERSION:-}" ]; then
  VERSION=$(gh api repos/:owner/:repo/tags --jq '.[0].name')
fi

# Strip 'v' prefix and bump version
VERSION="${VERSION#v}"
jq --arg ver "$VERSION" '.version = $ver' package.json > package.json.tmp && mv package.json.tmp package.json
