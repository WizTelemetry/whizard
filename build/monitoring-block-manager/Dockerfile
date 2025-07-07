# Build the manager binary
FROM golang:1.24 as builder
ARG GOPROXY

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN GOPROXY=$GOPROXY go mod download

# Copy the go source
COPY cmd/ cmd/
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 go build -a -o block-manager cmd/monitoring-block-manager/block-manager.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM thanosio/thanos:v0.39.1
WORKDIR /
COPY --from=builder /workspace/block-manager .
USER 65532:65532

ENTRYPOINT ["/block-manager"]
