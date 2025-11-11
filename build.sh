#!/bin/bash
set -e  # Exit on any error

echo "🔧 Gothic Forge Build Script for Leapcell"
echo "=========================================="

# 1. Install build tools
echo "📦 Installing build tools..."
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/bep/gotailwindcss/v2@latest

# 2. Add Go bin to PATH (for installed tools)
export PATH="$PATH:$(go env GOPATH)/bin"

# 3. Generate templ templates
echo "🎨 Generating templ templates..."
templ generate

# 4. Build CSS with Tailwind
echo "💅 Building CSS with Tailwind..."
gotailwindcss -i app/styles/input.css -o app/styles/output.css --minify

# 5. Build Go server binary
echo "🚀 Building Go server..."
go build -o server ./cmd/server

echo "✅ Build complete! Binary: ./server"
