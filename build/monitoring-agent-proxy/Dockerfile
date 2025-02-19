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

ARG GOLDFLAGS

# Build
RUN CGO_ENABLED=0 go build -a -ldflags="${GOLDFLAGS}" -o monitoring-gateway cmd/monitoring-gateway/*

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/monitoring-gateway .
USER 65532:65532

ENTRYPOINT ["/monitoring-gateway","agent-proxy"]
