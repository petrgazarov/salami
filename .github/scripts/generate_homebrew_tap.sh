#!/bin/env bash

set -euo pipefail

export SALAMI_VERSION="${1}"
TAP_FILE="./.github/scripts/salami_tap_template.rb"
TAP_FILE="$(realpath "${TAP_FILE}")"

cd "$(mktemp -d)"

>&2 echo "::info Generating Homebrew Tap..."
>&2 echo "::group::Download release assets"
>&2 gh release download --repo petrgazarov/salami "v${SALAMI_VERSION}" -p 'salami*darwin*.tar.gz' -p 'salami*linux*.tar.gz' --skip-existing
>&2 echo "::endgroup::"

for i in \
  "darwin x64   SALAMI_DARWIN_X64" \
  "darwin arm64 SALAMI_DARWIN_ARM64" \
  "linux x64    SALAMI_LINUX_X64" \
  "linux arm64  SALAMI_LINUX_ARM64" \
  ; do
  # shellcheck disable=SC2086 # intentional, we want to split the strings
  set -- $i # read loop strings as args
  OS="$1"
  ARCH="$2"
  ENV_VAR="$3"
  SHA256="$(sha256sum "salami-${SALAMI_VERSION}-${OS}-${ARCH}.tar.gz" | cut -f1 -d' ')"

  SHA256_VAR="${ENV_VAR}_SHA256"
  URL_VAR="${ENV_VAR}_URL"
  printf -v "${SHA256_VAR}" "%s" "${SHA256}"
  printf -v "${URL_VAR}" "%s" "https://github.com/petrgazarov/salami/releases/download/${SALAMI_VERSION}/salami-${SALAMI_VERSION}-${OS}-${ARCH}.tar.gz"

  export "${SHA256_VAR?}"
  export "${URL_VAR?}"
  >&2 echo "${OS}-${ARCH} SHA256: " "${!SHA256_VAR}"
  >&2 echo "${OS}-${ARCH} URL:    " "${!URL_VAR}"
done

# shellcheck disable=SC2016 # intentional, envsubst requires us to pass variable names with $ prefixes.
envsubst '$SALAMI_VERSION,$SALAMI_DARWIN_X64_URL,$SALAMI_DARWIN_X64_SHA256,$SALAMI_DARWIN_ARM64_URL,$SALAMI_DARWIN_ARM64_SHA256,$SALAMI_LINUX_X64_URL,$SALAMI_LINUX_X64_SHA256,$SALAMI_LINUX_ARM64_URL,$SALAMI_LINUX_ARM64_SHA256' < "${TAP_FILE}"