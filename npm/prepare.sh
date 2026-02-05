#!/usr/bin/env bash

set -euxo pipefail
cp -r ../README.md ../dist/* .

# bump version
LATEST_TAG=$(gh api repos/:owner/:repo/tags --jq '.[0].name' | sed 's/^v//')
jq --arg ver "$LATEST_TAG" '.version = $ver' package.json > package.json.tmp && mv package.json.tmp package.json
