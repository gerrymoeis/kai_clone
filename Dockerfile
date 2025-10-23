# Gothic Forge v3 - Production Dockerfile
#
# This Dockerfile uses multi-stage builds to:
# 1. Keep the final image small (security + faster deployments)
# 2. Separate build-time dependencies from runtime
# 3. Follow Docker best practices and least-privilege principles
#
# Educational Resources:
# - Multi-stage builds: https://docs.docker.com/build/building/multi-stage/
# - Security best practices: https://docs.docker.com/develop/security-best-practices/

# ═══════════════════════════════════════════════════════════════════════════
# Stage 1: Builder
# ═══════════════════════════════════════════════════════════════════════════
# Why golang:alpine? 
# - Alpine Linux is minimal (~5MB base) vs debian (~124MB)
# - Includes go toolchain for building
# - Smaller attack surface (fewer packages = fewer vulnerabilities)
FROM golang:1.24-alpine AS builder

# Install build dependencies
# - git: Required for go mod download with private repos
# - ca-certificates: Required for HTTPS requests during build
# - tzdata: Timezone data (copied to final image)
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go.mod and go.sum first (Docker layer caching optimization)
# Why? If dependencies haven't changed, Docker reuses this layer
# This makes rebuilds much faster when only code changes
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Copy the entire source code
COPY . .

# Build the application
# Flags explained:
# -ldflags="-s -w": Strip debug info and symbol table (reduces binary size by ~30%)
# -trimpath: Remove file system paths from binary (security + reproducible builds)
# CGO_ENABLED=0: Disable CGO for fully static binary (portable across Linux distros)
# GOOS=linux: Target Linux (even if building on Windows/macOS)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -trimpath \
    -o server \
    ./cmd/server

# ═══════════════════════════════════════════════════════════════════════════
# Stage 2: Runtime
# ═══════════════════════════════════════════════════════════════════════════
# Why scratch/alpine?
# Option A (scratch): Absolute minimal (~binary only, 15-20MB total)
# Option B (alpine): Minimal Linux with shell (~30MB total, easier debugging)
# We choose alpine for production debugging capabilities
FROM alpine:latest

# Install runtime dependencies (minimal)
# - ca-certificates: HTTPS certificate validation
# - tzdata: Timezone support (for logs, timestamps)
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
# Why? Running as root is a security risk. If the container is compromised,
# attacker has root privileges. Non-root limits damage.
# UID 1000 is standard for non-root users
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Create necessary directories with proper permissions
RUN mkdir -p /app/app/static /app/app/styles /app/app/db && \
    chown -R appuser:appuser /app

# Set working directory
WORKDIR /app

# Copy compiled binary from builder stage
COPY --from=builder /build/server /app/server

# Copy application assets
# These are needed for Templ templates, static files, and database migrations
COPY --chown=appuser:appuser app/static ./app/static
COPY --chown=appuser:appuser app/styles ./app/styles
COPY --chown=appuser:appuser app/db ./app/db

# Copy timezone data for consistent timestamps across deployments
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Switch to non-root user (all subsequent commands run as this user)
USER appuser

# Expose port (documentation only, doesn't actually publish the port)
# The application should read HTTP_PORT from environment
EXPOSE 8080

# Health check (Docker/Kubernetes use this to determine container health)
# Why /readyz? It checks if the app is ready to serve traffic (DB connected, etc.)
# Interval: Check every 30s
# Timeout: Wait max 3s for response
# Start period: Wait 10s after container starts before first check
# Retries: Mark unhealthy after 3 consecutive failures
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/readyz || exit 1

# Environment variables (defaults, override in deployment)
ENV APP_ENV=production \
    HTTP_HOST=0.0.0.0 \
    HTTP_PORT=8080 \
    LOG_FORMAT=json

# Run the application
# Why array form ["cmd", "arg"]? 
# - Proper signal handling (SIGTERM for graceful shutdown)
# - No shell intermediary (more secure, faster startup)
CMD ["/app/server"]
