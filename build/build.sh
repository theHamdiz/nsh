GOOS="darwin" GOARCH="arm64" CGO_ENABLED="0" go build -ldflags "-s -w" -o target/darwin/nsh main.go helpers.go
GOOS="linux" GOARCH="amd64" CGO_ENABLED="0" go build -ldflags "-s -w" -o target/linux/nsh main.go helpers.go
GOOS="windows" GOARCH="amd64" CGO_ENABLED="0" go build -ldflags "-s -w" -o target/win32/nsh.exe main.go helpers.go