#!/usr/bin/env bash

set -e

if [ ! -f "test.png" ]; then
  echo "Place a PNG file named test.png in the repository root before running this script."
  exit 1
fi

curl -X POST "http://localhost:8085/upload" \
  -F "image=@test.png" \
  -F "s3=false"

echo
