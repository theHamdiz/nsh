@echo off
IF EXIST build.ps1 (
    powershell -ExecutionPolicy Bypass -File "build.ps1"
    pause
) ELSE (
    IF EXIST build\build.ps1 (
        powershell -ExecutionPolicy Bypass -File "build\\build.ps1"
        pause
    ) ELSE (
        echo build.ps1 not found in the current directory or build\ directory.
        exit /b 1
    )
)

