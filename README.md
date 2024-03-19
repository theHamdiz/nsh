# 🚀 `nsh` (nameShift) Documentation

## Overview

`nsh`, formerly known as nameShift, is a versatile tool designed for comprehensive string transformations across files and directories. It offers a range of functionalities tailored for renaming files, modifying file contents, and filtering operations based on file extensions. `nsh` supports both synchronous and concurrent processing, accommodates case-sensitive or case-agnostic operations, and provides detailed reports on modifications.

## Key Features

- **File and Directory Renaming**: Easily rename files and directories.
- **String Replacement**: Perform string replacements within file contents.
- **File Extension Filtering**: Focus operations on files with specific extensions.
- **Processing Modes**: Choose between concurrent or synchronous processing.
- **Case Sensitivity Options**: Operate in either case-sensitive or case-agnostic mode.
- **Configurable Directory Exclusion**: Optionally include or exclude config directories.
- **Detailed Reporting**: Generate tabular reports detailing modifications and errors.
- **Flexible Flag Handling**: Use either short or long-form command-line flags.

## 🚧 **Build Instructions** 🚧

### Windows

1. Open the command prompt.
2. Navigate to the `nsh` root directory containing `build.bat`.
3. Execute the build script:
```bash
✅   .\build.bat
```

### Unix

1. Open the terminal.
2. Navigate to the `nsh` root directory containing `build.sh`.
3. Run the build script:
```zsh
✅   ./build.sh
```

## Installing `nsh` System-Wide

### Unix Systems

- Note: Installation may require elevated privileges (`sudo`).
- To install, run:

```zsh
✅  sudo python3 build/install.py
```

Or:

```zsh
✅  sudo ./build/install.py
```

Or simply use

```zsh
✅  go install
```

> If you're Ok with just installing it inside  `$GOPATH/bin` directory.

### Windows

- Execute the installation script:
```bash
✅  python build\\install.py
```

Or directly run:

```bash
✅  build\\install
```

Or simply use

```zsh
✅  go install
```

> If you're Ok with just installing it inside  `$GOPATH/bin` directory.

## Usage Examples

### Windows

```bash
✅ .\`nsh`.exe "path\\to\\directory" "OldText" "NewText" --ignore-config-dirs=true --work-globally=false --concurrent-run=false --case-matching=true --file-extensions=".go,.md"
```
Or, for an installed tool:
```bash
✅ `nsh` "path\\to\\directory" "OldText" "NewText" -i=true -g=false -cr=false -cm=true --exts=".go,.md"
```

### Unix Systems

```zsh
✅ ./`nsh` "path/to/directory" "OldText" "NewText" --ignore-config-dirs=true --work-globally=false -concurrent-run=false -case-matching=true --file-extensions=".go,.md"
```
Or for an installed tool:
```zsh
✅ `nsh` "path/to/directory" "OldText" "NewText" -i=true -g=false -cr=false -cm=true --ext=".go,.md"
```

## Advanced Options and Flexibility

`nsh` accommodates different user preferences with dual parameter formats (verbose and shorthand) and has a forgiving approach to typos and parameter variations. Its flexibility extends to accepting both `ext` and `exts` for specifying file extensions.

## Future Enhancements

- [ ] **GUI Integration**: Bringing the power of ``nsh`` to a graphical user interface.
- [ ] **Cross-Platform Package Managers**: Aim to distribute ``nsh`` through package managers like Homebrew, apt, and others, making installation a breeze.
- [ ] **Advanced Pattern Matching**: Implement regex support for the adventurers who need to capture or transform more complex string patterns.
- [ ] **Localization Support**: Support multiple languages.
- [ ] **Plugin Ecosystem**: Enabling the community to extend ``nsh`` with their own plugins.
- [ ] **FFI Function Exposure**: Enabling the community to use ``nsh`` outside of the go realm.
