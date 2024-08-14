set -ex

PATH=/tmp:$PATH
rm -rf completion
mkdir -p completion/bash completion/zsh
go build -o /tmp/platform cmd/platform/main.go
platform completion bash > completion/bash/platform.bash
platform completion zsh > completion/zsh/_platform

go build --tags=vendor,upsun -o /tmp/upsun cmd/platform/main.go
upsun completion bash > completion/bash/upsun.bash
upsun completion zsh > completion/zsh/_upsun

if [ -n "$VENDOR_BINARY" ]; then
    go build --tags=vendor -o /tmp/$VENDOR_BINARY cmd/platform/main.go
    $VENDOR_BINARY completion bash > completion/bash/$VENDOR_BINARY.bash
    $VENDOR_BINARY completion zsh > completion/zsh/_$VENDOR_BINARY
fi
