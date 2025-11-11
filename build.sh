#!/bin/sh
set -e  # Exit on any error

echo "Gothic Forge Build Script for Leapcell"
echo "======================================"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
# Keep original arch for gotailwindcss (uses x86_64, not amd64)
ARCH_ORIG=$(uname -m)
case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    armv7l) ARCH="arm" ;;
esac

# 1. Download and install templ
echo "Downloading templ..."
TEMPL_VERSION="v0.3.960"
TEMPL_URL="https://github.com/a-h/templ/releases/download/${TEMPL_VERSION}/templ_${OS}_${ARCH}.tar.gz"
wget -q -O templ.tar.gz "$TEMPL_URL"
tar -xzf templ.tar.gz templ
chmod +x templ
rm templ.tar.gz
echo "templ installed"

# 2. Download and install gotailwindcss
echo "Downloading gotailwindcss..."
# Correct repository: github.com/gotailwindcss/tailwind
TAILWIND_VERSION="v0.5.0"
OS_CAPS="Linux"
TAILWIND_URL="https://github.com/gotailwindcss/tailwind/releases/download/${TAILWIND_VERSION}/tailwindcss-extra_${TAILWIND_VERSION}_${OS_CAPS}_${ARCH_ORIG}.tar.gz"
wget -q -O tailwind.tar.gz "$TAILWIND_URL"
tar -xzf tailwind.tar.gz tailwindcss-extra
mv tailwindcss-extra gotailwindcss
chmod +x gotailwindcss
rm tailwind.tar.gz
echo "gotailwindcss installed"

# 3. Generate templ templates
echo "Generating templ templates..."
./templ generate

# 4. Build CSS with Tailwind
echo "Building CSS with Tailwind..."
./gotailwindcss -i app/styles/input.css -o app/styles/output.css --minify

# 5. Build Go server binary
echo "Building Go server..."
go build -o server ./cmd/server

# 6. Cleanup build tools
rm -f templ gotailwindcss

echo "Build complete! Binary: ./server"
