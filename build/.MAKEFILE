build:
    $env:GOOS="windows"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"; go build -ldflags "-s -w" -o target/win32/nsh.exe main.go helpers.go
    $env:GOOS="linux"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"; go build -ldflags "-s -w" -o target/linux/nsh main.go helpers.go
    $env:GOOS="darwin"; $env:GOARCH="arm64"; $env:CGO_ENABLED="0"; go build -ldflags "-s -w" -o target/darwin/nsh main.go helpers.go
    rsrc -ico nsh.ico -o nsh.syso
    go build -o target/win32/nsh.exe