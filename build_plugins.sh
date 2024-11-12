#!/bin/bash

if [ "$(uname)" == "Darwin" ]; then
    echo "Platform: Darwin"
    EXT="dylib"
else
    echo "Platform: Linux"
    EXT="so"
fi

build_plugin() {
    local path=$1
    local output=$2
    echo "Building $path"
    go build -buildmode=plugin -o $output $path
}

if [ $# -eq 1 ]; then
    if [ -f "plugins/$1" ]; then
        build_plugin "plugins/$1" "${1%.*}.$EXT"
    elif [ -d "plugins/$1" ] && [ -f "plugins/$1/main.go" ]; then
        build_plugin "plugins/$1/main.go" "$1.$EXT"
    else
        echo "Invalid input: $1"
        exit 1
    fi
    exit
fi

echo "Building all plugins"
for file in plugins/*.go; do
    build_plugin "$file" "${file%.*}.$EXT"
done

for dir in plugins/*/; do
    if [ -f "$dir/main.go" ]; then
        build_plugin "$dir/main.go" "${dir%/}.$EXT"
    fi
done