# Ensure target directories exist
$directories = @("target/win", "target/linux", "target/mac")
foreach ($dir in $directories) {
    if (-not (Test-Path $dir)) {
        New-Item -ItemType Directory -Force -Path $dir
    }
}

# Windows build
$env:GOOS="windows"
$env:GOARCH="amd64"
$env:CGO_ENABLED="0"
go build -ldflags "-s -w" -o .\target\win\nameShift.exe .\nameShift.go .\helpers.go

# Linux build
$env:GOOS="linux"
go build -ldflags "-s -w" -o .\target\linux\nameShift .\nameShift.go .\helpers.go

# macOS build
$env:GOOS="darwin"
$env:GOARCH="arm64"
go build -ldflags "-s -w" -o .\target\mac\nameShift .\nameShift.go .\helpers.go

# Reset environment variables for resource generation and final Windows build
Remove-Item Env:\GOOS
Remove-Item Env:\GOARCH
$env:CGO_ENABLED="0"
rsrc -ico logo.ico -o nameShift.syso
go build -ldflags "-s -w -r" -o .\target\win\nameShift.exe .\nameShift.go .\helpers.go

echo "Build Completed!"
