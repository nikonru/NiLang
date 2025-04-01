if [[ "$OSTYPE" = "linux-gnu" ]]; then
    platform="linux"
    ext=""
elif [[ "$OSTYPE" = "msys" ]]; then
    platform="windows"
    ext=".exe"
else
    echo "ERROR: unknown platform - $OSTYPE"
    exit 1
fi

if [[ "$1" = "-wasm" ]]; then
    GOOS=js 
    GOARCH=wasm
    ext=".wasm"
fi

go build -o build/nilang$ext  src/main.go

tar -czvf build/nilang-$platform.tar.gz --directory=build nilang$ext
