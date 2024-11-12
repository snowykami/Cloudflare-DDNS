if [ "$(uname)" == "Darwin" ]; then
    echo "Platform: Darwin"
    EXT="dylib"
else
    echo "Platform: Linux"
    EXT="so"
fi

if [ $# -eq 1 ]; then
    echo "Building $1"
    go build -buildmode=plugin -o ${1%.*}.$EXT plugins/$1
    exit
fi

echo "Building all plugins"
for file in plugins/*.go
do  
    echo "Building $file"
    go build -buildmode=plugin -o ${file%.*}.$EXT plugins/$file
done