#!/bin/bash

set -e

# Function to detect the operating system
function get_os() {
    uname_out="$(uname -s)"

    case "${uname_out}" in
        Linux*)     os=linux;;
        Darwin*)    os=darwin;;
        *)          os=unknown
    esac

    echo "${os}"
}

# Function to detect the architecture
function get_arch() {
    arch=$(uname -m)

    case "${arch}" in
        x86_64)     arch=amd64;;
        i*86)       arch=386;;
        armv7l)     arch=arm;;
        aarch64)    arch=arm64;;
        arm64)      arch=arm64;;
        *)          arch=unknown
    esac

    echo "${arch}"
}

function getLatestReleaseVersion() {
  if [ -n "${GITHUB_TOKEN}" ]; then
    AUTH_HEADER="Authorization: Bearer ${GITHUB_TOKEN}"
  fi

  # like "v1.2.3"
  latestVersion=$(curl -H "${AUTH_HEADER}" -s https://api.github.com/repos/zeabur/cli/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

  if [ -z "$latestVersion" ]; then
    echo "unknown"
  fi

  echo "${latestVersion}"
}

# Function to decide which binary to use based on the detected OS and architecture
# Input parameters:
# $1: Operating system (linux, darwin, windows)
# $2: Architecture (amd64, 386, arm, arm64)
# $3: Latest release version
# Output:
# The function echoes the binary filename or outputs an error message if unsupported OS/arch.
function get_binary_name() {
    local os="$1"
    local arch="$2"
    local latestVersion="$3"

    # Remove the "v" prefix from the version
    latestVersion="${latestVersion#v}"

    case "${os}_${arch}" in
        linux_amd64)    echo "zeabur_${latestVersion}_linux_amd64";;
        linux_arm64)    echo "zeabur_${latestVersion}_linux_arm64";;
        linux_386)      echo "zeabur_${latestVersion}_linux_386";;
        darwin_amd64)   echo "zeabur_${latestVersion}_darwin_amd64";;
        darwin_arm64)   echo "zeabur_${latestVersion}_darwin_arm64";;
        windows_amd64)  echo "zeabur_${latestVersion}_windows_amd64.exe";;
        windows_arm64)  echo "zeabur_${latestVersion}_windows_arm64.exe";;
        windows_386)    echo "zeabur_${latestVersion}_windows_386.exe";;
        *)              echo "unknown"
    esac
}

function downloadZeabur() {
  local latestVersion="$1"
  local binary="$2"

  # 1. download the release and rename it to "zeabur"
  # 2. count the download count of the release
  fullReleaseUrl="https://github.com/zeabur/cli/releases/download/${latestVersion}/${binary}"

  echo "Downloading Zeabur CLI from: $fullReleaseUrl"
  # use -L to follow redirects

  curl -L -o zeabur $fullReleaseUrl

  echo "Zeabur CLI downloaded completed"

  # grant execution rights
  chmod +x zeabur
}

function showZeaburHelp() {
  echo ""
  # show zeabur help and double check the download is success
  ./zeabur help
}

# Determine the operating system and architecture
OS=$(get_os)

if [ "${OS}" == "unknown" ]; then
    echo "This operating system is not supported by the install script."
    exit 1
fi

echo "Operating system: ${OS}"

ARCH=$(get_arch)

if [ "${ARCH}" == "unknown" ]; then
    echo "This architecture is not supported by the install script."
    exit 1
fi

echo "Architecture: ${ARCH}"

# Latest version of the binary with "v" prefix (change this to the desired version)
latestVersion=$(getLatestReleaseVersion)

if [ "${latestVersion}" == "unknown" ]; then
    echo "Could not find latest version"
    exit 1
fi

echo "Latest version: ${latestVersion}"

binary=$(get_binary_name "${OS}" "${ARCH}" "${latestVersion}")

if [ "${binary}" == "unknown" ]; then
    echo "No pre-built binary available for ${OS}_${ARCH}."
    exit 1
fi

downloadZeabur "${latestVersion}" "${binary}"

showZeaburHelp

sudo mv zeabur /usr/local/bin
