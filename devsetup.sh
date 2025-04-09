#!/bin/bash
set -e

OS="$(uname -s)"

# --- os-specific helpers ---
function is_mac() {
  [[ "$OS" == "Darwin" ]]
}

function is_linux() {
  [[ "$OS" == "Linux" ]]
}

function is_windows() {
  [[ "$OS" =~ MINGW|MSYS|CYGWIN ]]
}

# --- check go ---
if ! command -v go &> /dev/null; then
  echo "❌ Go is not installed."
  if is_mac; then
    echo "➡️  Install with: brew install go"
  elif is_linux; then
    echo "➡️  Install from https://go.dev/dl/"
  elif is_windows; then
    echo "➡️  Download installer: https://go.dev/dl/"
  fi
  exit 1
fi

echo "✅ go is installed: $(go version)"
echo "🔥 setting up dev environment..."

# --- .env setup ---
if [ ! -f .env.local ]; then
  echo "📄 copying .env.example to .env.local"
  cp .env.example .env.local
else
  echo "✅ .env.local already exists"
fi

# --- docker ---
if ! command -v docker &> /dev/null; then
  echo "🐳 docker not found."
  if is_mac; then
    echo "➡️  Install Docker Desktop from https://www.docker.com/products/docker-desktop/"
  elif is_linux; then
    echo "📦 installing docker..."
    curl -fsSL https://get.docker.com | bash
    sudo usermod -aG docker "$USER"
    echo "⚠️  log out and log back in to use docker without sudo"
  elif is_windows; then
    echo "➡️  Install Docker Desktop from https://www.docker.com/products/docker-desktop/"
  fi
else
  echo "✅ docker installed"
fi

# --- protoc ---
if ! command -v protoc &> /dev/null; then
  echo "📦 installing protoc..."
  if is_mac; then
    brew install protobuf
  elif is_linux; then
    PROTOC_ZIP=protoc-24.4-linux-x86_64.zip
    curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v24.4/$PROTOC_ZIP
    sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
    sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
    rm -f $PROTOC_ZIP
  else
    echo "⚠️  please install protoc manually: https://github.com/protocolbuffers/protobuf/releases"
  fi
else
  echo "✅ protoc already installed"
fi

# --- stringer ---
if ! command -v stringer &> /dev/null; then
  echo "🧵 installing stringer..."
  go install golang.org/x/tools/cmd/stringer@latest

  if [[ ":$PATH:" != *":$(go env GOPATH)/bin:"* ]]; then
    echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
    export PATH=$PATH:$(go env GOPATH)/bin
  fi
else
  echo "✅ stringer already installed"
fi

echo "✅ setup complete!"
