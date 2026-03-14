#!/bin/sh
set -e

PROJECT_DIR="/root/projects/image-server"
BINARY_DEST="/opt/imageserver/image-server"

cd "$PROJECT_DIR"

echo "Git pull..."
git pull

echo "Bygger binär..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o image-server

echo "Byter ut binär..."
mv image-server "$BINARY_DEST"

echo "Startar om tjänsten..."
rc-service image-server restart

echo "Klart!"
