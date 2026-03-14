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
python3 -c "
import json, sys
with open('$KEYFILE') as f:
    keys = json.load(f)
keys['$API_KEY'] = '$APP'
with open('$TMPFILE', 'w') as f:
    json.dump(keys, f, indent=2)
"
mv "$TMPFILE" "$KEYFILE"

echo "App:      $APP"
echo "API-nyckel: $API_KEY"
echo "Sparad i: $KEYFILE"
