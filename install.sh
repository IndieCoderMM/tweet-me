#!/usr/bin/env bash
set -euo pipefail

APP="tweet-me"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

usage() {
  cat <<EOF
Install $APP locally (~/.local/bin) or system-wide.

USAGE:
  ./install.sh [--system] [--force]

OPTIONS:
  --system   Install to /usr/local/bin (requires write permission / sudo)
  --force    Overwrite existing binary
  -h, --help Show this help
EOF
}

TARGET_DIR="$HOME/.local/bin"
FORCE=0
SYSTEM=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --system) SYSTEM=1; shift ;;
    --force) FORCE=1; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "Unknown arg: $1"; usage; exit 1 ;;
  esac
done

if [[ $SYSTEM -eq 1 ]]; then
  TARGET_DIR="/usr/local/bin"
fi

echo "Building $APP ..."
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo dev)
go build -trimpath -ldflags "-s -w -X main.version=$VERSION" -o "$REPO_ROOT/$APP" "$REPO_ROOT"

mkdir -p "$TARGET_DIR"
DEST="$TARGET_DIR/$APP"
if [[ -e "$DEST" && $FORCE -eq 0 ]]; then
  echo "Error: $DEST exists (use --force to overwrite)" >&2
  exit 1
fi

if [[ $SYSTEM -eq 1 && ! -w "$TARGET_DIR" ]]; then
  echo "Using sudo to copy to $TARGET_DIR"; sudo cp "$REPO_ROOT/$APP" "$DEST"
else
  cp "$REPO_ROOT/$APP" "$DEST"
fi

echo "Installed $APP $VERSION to $DEST"
case :$PATH: in
  *:"$TARGET_DIR":*) ;; # already on PATH
  *) echo "NOTE: Add $TARGET_DIR to your PATH" ;;
esac
