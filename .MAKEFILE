build:
    $env:GOOS="windows"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"; go build -ldflags "-s -w" -o target/win/nameShift.exe nameShift.go
    $env:GOOS="linux"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"; go build -ldflags "-s -w" -o target/linux/nameShift nameShift.go
    $env:GOOS="darwin"; $env:GOARCH="arm64"; $env:CGO_ENABLED="0"; go build -ldflags "-s -w" -o target/mac/nameShift nameShift.go
    rsrc -ico nameShift.ico -o nameShift.syso
    go build -o target/win/nameShift.exe
