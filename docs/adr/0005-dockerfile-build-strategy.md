# 0001 Dockerfile Build Strategy

## Status
Accepted

## Context
Go monolith with embedded frontend deployed to Raspberry Pi 4
(ARM64) from an x86 dev machine. We need a small, secure,
cross-compiled image.

## Decision
Use a multi-stage build with `golang:1.24-alpine` as the
builder and `gcr.io/distroless/static-debian12` as the final
image. Build with `--platform=$BUILDPLATFORM` so the builder
runs natively on x86. Cross-compile to ARM64 via
`TARGETOS`/`TARGETARCH` Go env vars. Build flags
`CGO_ENABLED=0` and `-ldflags="-s -w"` produce a fully static,
stripped binary. Final image runs as `nonroot:nonroot`.

## Consequences
Final image is ~2–5MB and has no shell or package manager,
minimizing attack surface. Cross-compilation is fast with no
QEMU emulation. Debugging inside the container is not possible,
use structured logs instead. If CGO is ever required, switch
final image to `distroless/base`.
