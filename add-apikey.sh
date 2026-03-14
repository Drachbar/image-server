#!/bin/sh
set -e

KEYFILE="${KEYFILE:-/etc/imageserver/apikeys.json}"

if [ -z "$1" ]; then
  echo "Användning: $0 <app-namn> [nyckelfil]"
  echo "  app-namn   namn på applikationen (t.ex. min-app)"
  echo "  nyckelfil  sökväg till JSON-fil (standard: $KEYFILE)"
  exit 1
fi

APP="$1"
if [ -n "$2" ]; then
  KEYFILE="$2"
fi

API_KEY=$(openssl rand -hex 32)

if [ ! -f "$KEYFILE" ]; then
  echo "{}" > "$KEYFILE"
fi

TMPFILE=$(mktemp)
jq --arg key "$API_KEY" --arg app "$APP" '.[$key] = $app' "$KEYFILE" > "$TMPFILE"
mv "$TMPFILE" "$KEYFILE"

echo "App:      $APP"
echo "API-nyckel: $API_KEY"
echo "Sparad i: $KEYFILE"
