#!/usr/bin/env bash

# Script to generate the pb code for the api

# Stop at the first error
set -e

# Get the directory of the project
SCRIPT_PATH=${BASH_SOURCE[0]}
# Follow symlinks
while [ -h "$SOURCE" ] ; 
    do SCRIPT_PATH="$(readlink "$SCRIPT_PATH")"; 
done
PROJECT_DIR="$( cd -P "$( dirname "${SCRIPT_PATH}" )/.." && pwd )"

PROTO_DIR="${PROJECT_DIR}/api/pb"
cd "$PROTO_DIR"

echo "==> Removing old .pb.go files"
# Ignore errors if no files are found
set +e
rm ./*.pb.go 2>/dev/null
set -e

echo "==> Collecting .proto files"
PROTO_FILES=$( ls ./*.proto )

echo "==> Generating new .pb.go files"
for FILE in "${PROTO_FILES[@]}"; do
    protoc \
        --go_out=. \
        --go_opt=paths=source_relative \
        --go-grpc_out=. \
        --go-grpc_opt=paths=source_relative \
        "${FILE}"
done
