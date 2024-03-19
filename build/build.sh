#!/bin/bash

# Function to check file existence and return correct path
get_file_path() {
    if [ -f "$1" ]; then
        echo "$1"
    elif [ -f "../$1" ]; then
        echo "../$1"
    else
        echo ""
    fi
}

# Determine the correct paths for main.go and helpers.go
main_path=$(get_file_path "main.go")
helpers_path=$(get_file_path "helpers.go")

# Check if both files are found
if [ -z "$main_path" ] || [ -z "$helpers_path" ]; then
    echo "Error: Required files (main.go, helpers.go) not found."
    exit 1
else
    # If both files are found, proceed with build for each target
    GOOS="darwin" GOARCH="arm64" CGO_ENABLED="0" go build -ldflags "-s -w" -o target/darwin/nsh $main_path $helpers_path
    GOOS="linux" GOARCH="amd64" CGO_ENABLED="0" go build -ldflags "-s -w" -o target/linux/nsh $main_path $helpers_path
    GOOS="windows" GOARCH="amd64" CGO_ENABLED="0" go build -ldflags "-s -w" -o target/win32/nsh.exe $main_path $helpers_path
fi
