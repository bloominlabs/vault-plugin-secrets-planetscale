#!/usr/bin/env bash
set -eEuo pipefail

MNT_PATH="planetscale"
PLUGIN_NAME="vault-plugin-secrets-planetscale"

#
# Helper script for local development. Automatically builds and registers the
# plugin. Requires `vault` is installed and available on $PATH.
#

# Get the right dir
DIR="$(cd "$(dirname "$(readlink "$0")")" && pwd)"

echo "==> Starting dev"

echo "--> Scratch dir"
echo "    Creating"
SCRATCH="${DIR}/tmp"
mkdir -p "${SCRATCH}/plugins"

function cleanup {
  echo ""
  echo "==> Cleaning up"
  kill -INT "${VAULT_PID}"
  rm -rf "${SCRATCH}"
}
trap cleanup EXIT

echo "--> Starting server"

export VAULT_TOKEN="root"
export VAULT_ADDR="http://127.0.0.1:8200"

vault server \
  -dev \
  -dev-plugin-init \
  -dev-plugin-dir "./bin/" \
  -dev-root-token-id "root" \
  -log-level "debug" \
  &
sleep 2
VAULT_PID=$!

echo "    Mouting plugin"
vault secrets enable -path=${MNT_PATH} -plugin-name=${PLUGIN_NAME} plugin

echo "==> Ready!"
wait ${VAULT_PID}
