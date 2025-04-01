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

additional_files=""
if [[ "$1" = "-wasm" ]]; then
    ext=".wasm"
    platform="js"

    if [[ "$OSTYPE" = "linux-gnu" ]]; then
        additional_files="$(go env GOROOT)/misc/wasm/wasm_exec.js $(go env GOROOT)/misc/wasm/wasm_exec.html"
    fi

    GOOS=js GOARCH=wasm go build -o build/nilang$ext src/main.go
else
    go build -o build/nilang$ext src/main.go
fi

tar -czvf build/nilang-$platform.tar.gz --directory=build nilang$ext $additional_files
