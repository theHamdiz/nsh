# Function to perform the build process
function Build-GoApplication {
    param (
        [string]$mainGoPath
    )

    echo "Starting the build process, ignore any previous errors!"
    # Ensure target directories exist
    $directories = @("target/win32", "target/linux", "target/darwin")
    foreach ($dir in $directories) {
        if (-not (Test-Path $dir)) {
            New-Item -ItemType Directory -Force -Path $dir
        }
    }

    # Set main.go directory for relative paths in build commands
    $pwd = Get-Location
    Set-Location (Get-Item $mainGoPath).DirectoryName

    # Windows build
    $env:GOOS="windows"
    $env:GOARCH="amd64"
    $env:CGO_ENABLED="0"
    go build -ldflags "-s -w" -o $pwd\target\win32\nsh.exe main.go helpers.go

    # Linux build
    $env:GOOS="linux"
    go build -ldflags "-s -w" -o $pwd\target\linux\nsh main.go helpers.go

    # macOS build
    $env:GOOS="darwin"
    $env:GOARCH="arm64"
    go build -ldflags "-s -w" -o $pwd\target\darwin\nsh main.go helpers.go

    # Reset environment variables for resource generation and final Windows build
    Remove-Item Env:\GOOS
    Remove-Item Env:\GOARCH
    $env:CGO_ENABLED="0"
    rsrc -ico assets\nsh.ico -o assets\nsh.syso
    go build -ldflags "-s -w -r" -o $pwd\target\win32\nsh.exe main.go helpers.go

    # Restore previous working directory
    Set-Location $pwd
    echo "Build Completed!"
}

# Check for main.go in the current directory
if (Test-Path -Path ".\main.go") {
    Build-GoApplication ".\main.go"
} elseif (Test-Path -Path "..\main.go") {
    # If main.go is not found, check the parent directory and modify paths accordingly
    Build-GoApplication "..\main.go"
} else {
    Write-Error "main.go was not found in the current or parent directory."
    exit 1
}
