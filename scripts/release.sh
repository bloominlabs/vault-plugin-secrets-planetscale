VERSION=0.6.0

earthly +build
cp ./bin/vault-plugin-secrets-planetscale ./bin/vault-plugin-secrets-planetscale-linux-amd64
sha256sum ./bin/vault-plugin-secrets-planetscale-linux-amd64 > ./bin/SHA256SUMS

gh release -R bloominlabs/vault-plugin-secrets-planetscale create --title v$VERSION v$VERSION ./bin/*
