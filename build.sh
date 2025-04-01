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

    GOOS=js GOARCH=wasm go build -o build/nilang$ext src/wasm.go

    additional_files="wasm_exec.js index.html"
    cp "wasm/wasm_exec.js" "build/wasm_exec.js"
    cp "wasm/index.html" "build/index.html"
else
    go build -o build/nilang$ext src/main.go
fi

tar -czvf build/nilang-$platform.tar.gz --directory=build nilang$ext $additional_files
