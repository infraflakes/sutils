#!/usr/bin/env bash

# Exit immediately if a command exits with a non-zero status
set -e

REPO="infraflakes/sutils"
BINARY_NAME="sn"
INSTALL_DIR="/usr/local/bin"

echo "Checking for the latest release of $REPO..."

# 1. Fetch the latest release download URL via GitHub API
DOWNLOAD_URL=$(curl -s https://api.github.com/repos/$REPO/releases/latest | \
    grep -oP '"browser_download_url": "\K[^"\s]+\_linux\_amd64\.tar\.gz(?=")')

if [ -z "$DOWNLOAD_URL" ]; then
    echo "Error: Could not find a linux_amd64.tar.gz asset for the latest release."
    exit 1
fi

echo "Found latest release: $DOWNLOAD_URL"

# 2. Create a temporary directory for downloading and extracting
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT # Clean up the temp dir even if the script fails later

echo "Downloading..."
curl -L "$DOWNLOAD_URL" -o "$TMP_DIR/sutils.tar.gz"

echo "Extracting..."
tar -xzf "$TMP_DIR/sutils.tar.gz" -C "$TMP_DIR"

# 3. Determine whether to use sudo or doas
if [ -f "$TMP_DIR/$BINARY_NAME" ]; then
    echo "Installing '$BINARY_NAME' to $INSTALL_DIR..."
    
    if command -v sudo >/dev/null 2>&1; then
        PRIV_CMD="sudo"
    elif command -v doas >/dev/null 2>&1; then
        PRIV_CMD="doas"
    else
        echo "Error: Neither 'sudo' nor 'doas' was found. Cannot elevate privileges."
        exit 1
    fi

    # Move the binary using the detected privilege elevator
    $PRIV_CMD mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    $PRIV_CMD chmod +x "$INSTALL_DIR/$BINARY_NAME"

    echo "Successfully installed $BINARY_NAME to $INSTALL_DIR!"
else
    echo "Error: Binary '$BINARY_NAME' not found in the extracted archive."
    exit 1
fi
