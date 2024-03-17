# Ensure target directories exist
$directories = @("target/win32", "target/linux", "target/darwin")
foreach ($dir in $directories) {
    if (-not (Test-Path $dir)) {
        New-Item -ItemType Directory -Force -Path $dir
    }
}

# Windows build
$env:GOOS="windows"
$env:GOARCH="amd64"
$env:CGO_ENABLED="0"
go build -ldflags "-s -w" -o target\win\nsh.exe main.go helpers.go

# Linux build
$env:GOOS="linux"
go build -ldflags "-s -w" -o target\linux\nsh main.go helpers.go

# macOS build
$env:GOOS="darwin"
$env:GOARCH="arm64"
go build -ldflags "-s -w" -o target\mac\nsh main.go helpers.go

# Reset environment variables for resource generation and final Windows build
Remove-Item Env:\GOOS
Remove-Item Env:\GOARCH
$env:CGO_ENABLED="0"
rsrc -ico nsh.ico -o nsh.syso
go build -ldflags "-s -w -r" -o target\win\nsh.exe main.go helpers.go

echo "Build Completed!"
