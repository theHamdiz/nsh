// Author: Ahmad Hamdi
//  .\nsh.exe "path/to/directory" "OldText" "NewText" -ignore-config-dirs=true -work-globally=false -concurrent-run=false -case-matching=true -file-extensions=".go,.md"

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ignoreConfig      bool
	workGlobally      bool
	concurrentRun     bool
	caseMatching      bool
	Version           = "0.1.3"
	versionFlag       bool
	fileExtensions    string
	errorsCount       int32
	replacementsCount int32
	errorReport       table.Writer
)

func customFlagParsing() {
	//log.Println("> Inside customFlagParsing")
	for i, arg := range os.Args {
		if strings.HasPrefix(arg, "--") {
			//log.Println("> Inside customFlagParsing for loop")
			os.Args[i] = strings.Replace(arg, "--", "-", -1)
			//log.Printf("> arg before : %s and after: %s\n", arg, os.Args[1])
		}
	}
}

func init() {
	//log.Println("> Entering init function for initializing the flags")
	customFlagParsing()
	flag.BoolVar(&ignoreConfig, "ignore-config-dirs", true, "Ignore .config directories 🚫🐙")
	flag.BoolVar(&ignoreConfig, "i", true, "Ignore .config directories 🚫🐙")

	flag.BoolVar(&workGlobally, "work-globally", false, "Work on folder names, file names, and file contents (default false)🌍✨")
	flag.BoolVar(&workGlobally, "g", false, "Work on folder names, file names, and file contents (default false)🌍✨")

	flag.BoolVar(&concurrentRun, "concurrent-run", false, "Run each folder inside the root directory in a separate goroutine (default false)🏃💨")
	flag.BoolVar(&concurrentRun, "cr", false, "Run each folder inside the root directory in a separate goroutine (default false)🏃💨")

	flag.BoolVar(&caseMatching, "case-matching", true, "Match case when replacing strings (default true) 👔🔍")
	flag.BoolVar(&caseMatching, "cm", true, "Match case when replacing strings (default true) 👔🔍")

	flag.StringVar(&fileExtensions, "file-extensions", ".go", "Comma-separated list of file extensions to process, e.g., '.go,.txt' 📄✂️")
	flag.StringVar(&fileExtensions, "ext", ".go", "Comma-separated list of file extensions to process, e.g., '.go,.txt' 📄✂️")
	flag.StringVar(&fileExtensions, "exts", ".go", "Comma-separated list of file extensions to process, e.g., '.go,.txt' 📄✂️")

	flag.BoolVar(&versionFlag, "version", false, "Get the current program version 🚀")
	flag.BoolVar(&versionFlag, "v", false, "Get the current program version 🚀")
	flag.BoolVar(&versionFlag, "V", false, "Get the current program version 🚀")
}

func ignoreConfigDirs(path string, err error) error {
	dirName := filepath.Base(path)
	if strings.HasPrefix(dirName, ".") && ignoreConfig {
		// If the directory name starts with '.', skip it
		return filepath.SkipDir
	}

	if err != nil {
		if os.IsPermission(err) {
			return filepath.SkipDir // Skip this file or directory but continue walking
		}
		atomic.AddInt32(&errorsCount, 1)
		return fmt.Errorf("> Error while attempting to ignore .config dirs: %w", err)
	}
	return nil
}

func processPath(ctx *AppContext, path string, info os.FileInfo, theStringToBeReplaced, theReplacementString string) error {
	if err := ignoreConfigDirs(path, nil); err != nil {
		row := []table.Row{{"Path", path, "Error", err}}
		ctx.AddError()
		ctx.AddErrorReportRow(row)
		return err
	}

	// Directly use the newly abstracted renameEntity function for files and directories.
	if workGlobally && (info.IsDir() || strings.Contains(info.Name(), theStringToBeReplaced)) {
		if err := renameEntity(ctx, path, theStringToBeReplaced, theReplacementString); err != nil {
			row := []table.Row{{"Path", path, "Error", fmt.Sprintf("Could not rename: %v", err)}}
			ctx.AddError()
			ctx.AddErrorReportRow(row)
			return err
		}
	}

	// For files, check if they should be processed and then process.
	if shouldProcessFile(path, info) {
		return processFile(path, theStringToBeReplaced, theReplacementString)
	}

	return nil
}

func renameEntity(ctx *AppContext, entityPath, theStringToBeReplaced, theReplacementString string) error {
	newPath := strings.Replace(entityPath, theStringToBeReplaced, theReplacementString, -1)
	if err := moveFileWithRetry(entityPath, newPath, 6); err != nil {
		if os.IsPermission(err) {
			if permErr := os.Chmod(entityPath, 0666); permErr != nil {
				ctx.AddError()
				return permErr // Permission change failed, return the error
			}
			if retryErr := moveFileWithRetry(entityPath, newPath, 6); retryErr != nil {
				ctx.AddError()
				return retryErr // Rename still failed, return the error
			}
		} else {
			ctx.AddError()
			return err // Non-permission error encountered, return it
		}
	}
	ctx.AddReplacement()
	return nil // Successfully renamed the entity
}

func replaceString(original, toReplace, replacement string, caseSensitive bool) string {
	if caseSensitive {
		return strings.Replace(original, toReplace, replacement, -1)
	}
	regex := regexp.MustCompile("(?i)" + regexp.QuoteMeta(toReplace))
	return regex.ReplaceAllString(original, replacement)
}

func processFile(path, theStringToBeReplaced, theReplacementString string) error {
	originalFile, err := os.Open(path)
	if err != nil {
		atomic.AddInt32(&errorsCount, 1)
		return err
	}
	defer func(originalFile *os.File) {
		err := originalFile.Close()
		if err != nil {
			fmt.Println("> Error while attempting to close the file: %w", err)
		}
	}(originalFile) // We'll still defer the close here, as it's simpler and still safe.

	// Create a temp file. Note: We're not deferring the cleanup here because we want to control it precisely.
	tempFile, err := os.CreateTemp("", "nsh_temp_file_")
	if err != nil {
		atomic.AddInt32(&errorsCount, 1)
		return err
	}

	// Ensure we clean up the temp file in every possible exit path after this point.
	// This defer statement is critical for making sure the temp file is always removed.
	defer func() {
		err := tempFile.Close()
		if err != nil {
			return
		} // Attempt to close the temp file. Ignoring errors here as we're going to delete it anyway.
		err = os.Remove(tempFile.Name())
		if err != nil {
			return
		} // Attempt to remove the file. If this fails, there's not much more we can do.
	}()

	scanner := bufio.NewScanner(originalFile)
	writer := bufio.NewWriter(tempFile)

	for scanner.Scan() {
		line := scanner.Text()
		modifiedLine := replaceString(line, theStringToBeReplaced, theReplacementString, caseMatching)
		if _, err := writer.WriteString(modifiedLine + "\n"); err != nil {
			atomic.AddInt32(&errorsCount, 1)
			return err // No need for additional cleanup, defer will handle it.
		}
		if modifiedLine != line {
			atomic.AddInt32(&replacementsCount, 1) // Increment only if a replacement occurred
		}
	}
	if err := scanner.Err(); err != nil {
		atomic.AddInt32(&errorsCount, 1)
		return err // No need for additional cleanup, defer will handle it.
	}

	if err := writer.Flush(); err != nil {
		atomic.AddInt32(&errorsCount, 1)
		return err // No need for additional cleanup, defer will handle it.
	}

	// Closing the temp file before renaming it. This is necessary on some systems like Windows.
	if err := tempFile.Close(); err != nil {
		atomic.AddInt32(&errorsCount, 1)
		return err // The file is still going to be removed due to the defer.
	}

	// Replace the original file with the temp file.
	if err := moveFileWithRetry(tempFile.Name(), path, 6); err != nil {
		atomic.AddInt32(&errorsCount, 1)
		return err // The file is still going to be removed due to the defer.
	}

	return nil
}

// moveFile handles moving a file from src to dst, working across different file systems/devices.
func moveFile(src, dst string) error {
	// Open the source file for reading.
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file for writing.
	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// Copy the contents of the source file to the destination file.
	if _, err := io.Copy(destinationFile, sourceFile); err != nil {
		return err
	}

	// Ensure the destination file is fully written and closed.
	if err := destinationFile.Sync(); err != nil {
		return err
	}

	// Delete the original (source) file.
	if err := os.Remove(src); err != nil {
		return err
	}

	return nil
}

func moveFileWithRetry(src, dst string, maxRetries int) error {
	var lastErr error
	retryDelay := 2 // Initial delay in seconds

	for attempt := 0; attempt < maxRetries; attempt++ {
		err := moveFile(src, dst)
		if err == nil {
			return nil // Success, file moved
		}

		lastErr = err
		// Log or print the retry attempt and wait time, can be helpful for debugging
		fmt.Printf("> Attempt %d failed to move file. Retrying in %d seconds...\n", attempt+1, retryDelay)

		time.Sleep(time.Duration(retryDelay) * time.Second)
		retryDelay *= 2 // Exponential increase of the wait time for the next retry
	}

	return fmt.Errorf("> moveFileWithRetry failed after %d attempts: %v", maxRetries, lastErr)
}

func shouldProcessFile(path string, info os.FileInfo) bool {
	if !info.IsDir() {
		exts := strings.Split(fileExtensions, ",")
		for _, ext := range exts {
			if strings.HasSuffix(path, ext) {
				return true
			}
		}
	}
	return false
}

func collectPaths(startingDir string, ignoreConfig bool) ([]string, error) {
	var paths []string
	err := filepath.Walk(startingDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .config directories if ignoreConfig is true
		if ignoreConfig && info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		paths = append(paths, path)
		return nil
	})

	return paths, err
}

func processPaths(ctx *AppContext, paths []string, theStringToBeReplaced, theReplacementString string) {
	var wg sync.WaitGroup

	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()

			info, err := os.Stat(p)
			if err != nil {
				// Handle error, possibly update ctx with the error information
				ctx.AddError()
				return
			}

			if err := processPath(ctx, p, info, theStringToBeReplaced, theReplacementString); err != nil {
				// Handle error, possibly update ctx with the error information
				ctx.AddError()
			}
		}(path)
	}

	wg.Wait()
}

func main() {
	resetColors()
	printLogo()
	flag.Parse()
	ctx := NewAppContext()

	if versionFlag {
		color.Cyan(fmt.Sprintf("\n> NameShifter Version: %s 🚀📚\n", Version))
		os.Exit(0)
	}

	if len(flag.Args()) < 3 {
		color.Red(fmt.Sprintf("\n> Usage: go run nsh.go <startingDirectory> <theStringToBeReplaced> <theReplacementString> -flags❗📚👀"))
		os.Exit(1)
	}

	args := flag.Args()
	startingDirectory, theStringToBeReplaced, theReplacementString := args[0], args[1], args[2]
	paths, err := collectPaths(startingDirectory, ignoreConfig)
	if err != nil {
		fmt.Println("Error collecting paths:", err)
		os.Exit(1)
	}
	printSettings()

	var wg sync.WaitGroup

	if concurrentRun {
		processPaths(ctx, paths, theStringToBeReplaced, theReplacementString)
	} else {
		for _, path := range paths {
			info, err := os.Stat(path)
			if err != nil {
				// Handle error, possibly update ctx with the error information
				ctx.AddError()
				continue
			}
			err = processPath(ctx, path, info, theStringToBeReplaced, theReplacementString)
			if err != nil {
				//fmt.Println("> Error processing path:", err)
				continue
			}
		}
	}

	wg.Wait()

	if err != nil {
		row := []table.Row{{3, fmt.Sprintf("Error walking through %s 😢👣", startingDirectory), err}}
		ctx.AddErrorReportRow(row)
		//os.Exit(2)
	} else {
		color.Green(fmt.Sprintf("\n> Names Shifted Successfully! 🎉📝✅\n"))
	}

	ctx.DisplayErrorReport()
	ctx.ReplacementsAndErrorsReport()
	os.Exit(0)
}
