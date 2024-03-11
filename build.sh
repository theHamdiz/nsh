GOOS="darwin" GOARCH="arm64" CGO_ENABLED="0" go build -ldflags "-s -w" -o target/mac/nameShift nameShift.go helpers.go
GOOS="linux" GOARCH="amd64" CGO_ENABLED="0" go build -ldflags "-s -w" -o target/linux/nameShift nameShift.go helpers.go
GOOS="windows" GOARCH="amd64" CGO_ENABLED="0" go build -ldflags "-s -w" -o target/windows/nameShift.exe nameShift.go helpers.go