#!/bin/bash

# Push Notify Build Script for Linux
set -e

# Define paths
ROOT_PATH=$(pwd)
FRONTEND_PATH="$ROOT_PATH/frontend"
BACKEND_PATH="$ROOT_PATH/backend"
STATIC_DIST_PATH="$BACKEND_PATH/static/dist"

echo -e "\033[0;36m>>> Starting Build Process...\033[0m"

# 1. Build Frontend
echo -e "\033[0;33m>>> Building Frontend...\033[0m"
cd "$FRONTEND_PATH"
npm install
npm run build

# 2. Prepare Backend Static Directory
echo -e "\033[0;33m>>> Preparing Backend Static Assets...\033[0m"
if [ -d "$STATIC_DIST_PATH" ]; then
    rm -rf "$STATIC_DIST_PATH"
fi
mkdir -p "$STATIC_DIST_PATH"

# Copy dist content to backend/static/dist
cp -r "$FRONTEND_PATH/dist/"* "$STATIC_DIST_PATH/"

# 3. Build Backend
echo -e "\033[0;33m>>> Building Backend...\033[0m"
cd "$BACKEND_PATH"
go build -o "$ROOT_PATH/push-notify" main.go

echo -e "\033[0;32m>>> Build Finished! Binary is at: $ROOT_PATH/push-notify\033[0m"
cd "$ROOT_PATH"
