@echo off
IF EXIST install.py (
    python install.py %*
) ELSE (
    IF EXIST build\install.py (
        python build\install.py %*
    ) ELSE (
        echo install.py not found in the current directory or build\ directory.
        exit /b 1
    )
)
