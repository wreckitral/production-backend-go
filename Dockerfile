# Build stage
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

WORKDIR /app

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build for target platform (ARM64 for Pi)
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w" -o blog .

# Final stage
FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/blog .

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/blog"]
