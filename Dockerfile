# syntax=docker/dockerfile:1

# ---- build stage ----
FROM golang:1.26-alpine AS build
WORKDIR /src

# Cache dependencies first.
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Build the static binary.
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -tags timetzdata -ldflags="-s -w" -o /out/arista-eos-mcp ./cmd/arista-eos-mcp

# ---- runtime stage ----
# scratch: no OS package database, so no Debian tzdata/ca-certificates to go
# stale in scans. Timezone data is embedded in the binary (-tags timetzdata) and
# the CA bundle is copied from the (current) build stage.
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /out/arista-eos-mcp /usr/local/bin/arista-eos-mcp

# HTTP transport listens here by default (override with MCP_HTTP_ADDR).
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/usr/local/bin/arista-eos-mcp"]
