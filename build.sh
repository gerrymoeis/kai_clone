#!/bin/sh
set -e  # Exit on any error

echo "Gothic Forge Build Script for Leapcell"
echo "======================================"

# 1. Install Git (required for go install)
echo "Installing Git..."
apk add --no-cache git

# 2. Download and install templ binary (has releases)
echo "Downloading templ..."
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    armv7l) ARCH="arm" ;;
esac
TEMPL_VERSION="v0.3.960"
TEMPL_URL="https://github.com/a-h/templ/releases/download/${TEMPL_VERSION}/templ_${OS}_${ARCH}.tar.gz"
wget -q -O templ.tar.gz "$TEMPL_URL"
tar -xzf templ.tar.gz templ
chmod +x templ
rm templ.tar.gz
echo "templ installed"

# 3. Install gotailwindcss using go install (no releases available)
echo "Installing gotailwindcss..."
go install github.com/gotailwindcss/tailwind/cmd/gotailwindcss@latest
echo "gotailwindcss installed"

# Add Go bin to PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# 3. Generate templ templates
echo "Generating templ templates..."
./templ generate

# 4. Build CSS with Tailwind
echo "Building CSS with Tailwind..."
gotailwindcss build -o app/styles/output.css app/styles/input.css

# 5. Build Go server binary
echo "Building Go server..."
go build -o server ./cmd/server

# 6. Cleanup build tools
rm -f templ gotailwindcss

echo "Build complete! Binary: ./server"
