VERSION 0.6
FROM golang:1.19
WORKDIR /vault-plugin-secrets-planetscale

deps:
    COPY go.mod go.sum ./
    RUN go mod download
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

build:
    FROM +deps
    COPY *.go .
    COPY --dir ./cmd .
    RUN CGO_ENABLED=0 go build -o bin/vault-plugin-secrets-planetscale cmd/planetscale/main.go
    SAVE ARTIFACT bin/vault-plugin-secrets-planetscale /planetscale AS LOCAL bin/vault-plugin-secrets-planetscale

test:
    FROM +deps
    COPY *.go .
    ARG TEST_planetscale_DOMAIN=https://test-stratos-host.us.planetscale.com
    RUN --secret PLANETSCALE_SERVICE_TOKEN --secret $PLANETSCALE_SERVICE_TOKEN_ID CGO_ENABLED=0 go test github.com/bloominlabs/vault-plugin-secrets-planetscale

dev:
  BUILD +build
  LOCALLY
  RUN bash ./scripts/dev.sh

all:
  BUILD +build
  BUILD +test
