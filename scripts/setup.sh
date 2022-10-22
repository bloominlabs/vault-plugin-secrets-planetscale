#!/usr/bin/env bash
# envconsul -config envconsul.hcl -once ./scripts/setup.sh
export VAULT_ADDR="http://127.0.0.1:8200" 
export VAULT_TOKEN="root" 
# export TIP_ENDPOINT="https://planetscale.service.consul/query"
# for testing locally with planetscale service
export TIP_ENDPOINT="http://127.0.0.1:8080/query"

 # TODO: get planetscale_identifier from consul 
vault write planetscale/config/root \
  service_token=$SERVICE_TOKEN \
  service_token_name=$SERVICE_TOKEN_NAME

vault write planetscale/roles/fjord organization=bloominlabs database=bloominlabs role=admin branch=main
vault read planetscale/roles/fjord

vault read planetscale/creds/fjord
