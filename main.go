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
	"time"
)

// Config encapsulates application-wide configurations.
type Config struct {
	IgnoreConfig   bool
	WorkGlobally   bool
	ConcurrentRun  bool
	CaseMatching   bool
	FileExtensions []string
	VersionFlag    bool
	Version        string
}

func NewConfig() *Config {
	cfg := &Config{
		Version: "0.2.1", // Assuming this is a constant for now
	}
	flag.BoolVar(&cfg.IgnoreConfig, "ignore-config-dirs", true, "Ignore .config directories 🚫🐙")
	flag.BoolVar(&cfg.IgnoreConfig, "i", true, "Ignore .config directories 🚫🐙")
	flag.BoolVar(&cfg.WorkGlobally, "work-globally", false, "Work on folder names, file names, and file contents 🌍✨")
	flag.BoolVar(&cfg.WorkGlobally, "g", false, "Work on folder names, file names, and file contents 🌍✨")
	flag.BoolVar(&cfg.ConcurrentRun, "concurrent-run", false, "Run each folder inside the root directory in a separate goroutine 🏃💨")
	flag.BoolVar(&cfg.ConcurrentRun, "cr", false, "Run each folder inside the root directory in a separate goroutine 🏃💨")
	flag.BoolVar(&cfg.CaseMatching, "case-matching", true, "Match case when replacing strings 👔🔍")
	flag.BoolVar(&cfg.CaseMatching, "cm", true, "Match case when replacing strings 👔🔍")
	var fileExtensions string
	flag.StringVar(&fileExtensions, "file-extensions", ".go,.md", "Comma-separated list of file extensions to process, e.g., '.go,.md' 📄✂️")
	flag.StringVar(&fileExtensions, "ext", ".go,.md", "Comma-separated list of file extensions to process, e.g., '.go,.md' 📄✂️")
	flag.StringVar(&fileExtensions, "exts", ".go,.md", "Comma-separated list of file extensions to process, e.g., '.go,.md' 📄✂️")

	flag.Parse()

	cfg.FileExtensions = strings.Split(fileExtensions, ",")
	return cfg
}

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
}

// NameShifter encapsulates all functionalities related to the name shifting process.
type NameShifter struct {
	Config  *Config
	Context *AppContext
}

// NewNameShifter creates a new instance of NameShifter with given configuration and context.
func NewNameShifter(cfg *Config, ctx *AppContext) *NameShifter {
	return &NameShifter{
		Config:  cfg,
		Context: ctx,
	}
}

// collectPaths walks the starting directory and collects all paths.
func (ns *NameShifter) collectPaths(startingDir string) ([]string, error) {
	var paths []string
	err := filepath.Walk(startingDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .config directories if ignoreConfig is true
		if ns.Config.IgnoreConfig && info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		paths = append(paths, path)
		return nil
	})

	return paths, err
}

// ProcessAllPaths decides whether to process paths concurrently or sequentially based on the configuration.
func (ns *NameShifter) ProcessAllPaths(paths []string, theStringToBeReplaced, theReplacementString string) {
	if ns.Config.ConcurrentRun {
		ns.processPathsConcurrently(paths, theStringToBeReplaced, theReplacementString)
	} else {
		ns.processPathsSequentially(paths, theStringToBeReplaced, theReplacementString)
	}
}

// processPathsConcurrently processes paths in parallel using goroutines.
func (ns *NameShifter) processPathsConcurrently(paths []string, theStringToBeReplaced, theReplacementString string) {
	var wg sync.WaitGroup
	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			ns.processSinglePath(p, theStringToBeReplaced, theReplacementString)
		}(path)
	}
	wg.Wait()
}

// processPathsSequentially processes paths one after another.
func (ns *NameShifter) processPathsSequentially(paths []string, theStringToBeReplaced, theReplacementString string) {
	for _, path := range paths {
		ns.processSinglePath(path, theStringToBeReplaced, theReplacementString)
	}
}

// processSinglePath processes a single path, deciding whether to rename the entity and/or process the file.
func (ns *NameShifter) processSinglePath(path, theStringToBeReplaced, theReplacementString string) {
	info, err := os.Stat(path)
	if err != nil {
		ns.Context.AddError()
		return
	}

	if err := ns.ignoreConfigDirs(path, nil); err != nil {
		ns.Context.AddError()
		return // Skip this path due to error or it being a directory we're ignoring
	}

	// Handling directories and files
	if ns.Config.WorkGlobally && (info.IsDir() || strings.Contains(info.Name(), theStringToBeReplaced)) {
		if err := ns.renameEntity(path, theStringToBeReplaced, theReplacementString); err != nil {
			ns.Context.AddError()
			return
		}
	} else if !info.IsDir() && ns.shouldProcessFile(path, info) {
		if err := ns.processFile(path, theStringToBeReplaced, theReplacementString); err != nil {
			ns.Context.AddError()
			return
		}
	}
}

func (ns *NameShifter) ignoreConfigDirs(path string, err error) error {
	dirName := filepath.Base(path)
	if (strings.HasPrefix(dirName, ".") || strings.HasPrefix(dirName, "venv")) && ns.Config.IgnoreConfig {
		// If the directory name starts with '.', 'venv' skip it
		return filepath.SkipDir
	}

	if err != nil {
		if os.IsPermission(err) {
			return filepath.SkipDir // Skip this file or directory but continue walking
		}
		ns.Context.AddError()
		return fmt.Errorf("> Error while attempting to ignore .config dirs: %w", err)
	}
	return nil
}

// replaceString replaces all occurrences of toReplace with replacement in the original string.
func (ns *NameShifter) replaceString(original, toReplace, replacement string) string {
	if ns.Config.CaseMatching {
		return strings.Replace(original, toReplace, replacement, -1)
	}
	regex := regexp.MustCompile("(?i)" + regexp.QuoteMeta(toReplace))
	return regex.ReplaceAllString(original, replacement)
}

func (ns *NameShifter) processFile(path, theStringToBeReplaced, theReplacementString string) error {
	originalFile, err := os.Open(path)
	if err != nil {
		ns.Context.AddError()
		return err
	}
	defer originalFile.Close()

	// Create a temp file
	tempFile, err := os.CreateTemp("", "nsh_temp_file_")
	if err != nil {
		ns.Context.AddError()
		return err
	}
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name()) // Cleanup temp file regardless of success
	}()

	scanner := bufio.NewScanner(originalFile)
	writer := bufio.NewWriter(tempFile)

	for scanner.Scan() {
		line := scanner.Text()
		modifiedLine := ns.replaceString(line, theStringToBeReplaced, theReplacementString)
		if _, err := writer.WriteString(modifiedLine + "\n"); err != nil {
			ns.Context.AddError()
			return err
		}
		if modifiedLine != line {
			ns.Context.AddReplacement() // Corrected to use the context's method for incrementing
		}
	}
	if err := scanner.Err(); err != nil {
		ns.Context.AddError()
		return err
	}

	if err := writer.Flush(); err != nil {
		ns.Context.AddError()
		return err
	}

	// Ensure the temp file is closed before attempting to rename
	if err := tempFile.Close(); err != nil {
		ns.Context.AddError()
		return err
	}

	// Replace the original file with the temp file
	if err := ns.moveFileWithRetry(tempFile.Name(), path, 6); err != nil {
		ns.Context.AddError()
		return err
	}

	return nil
}

// moveFile handles moving a file from src to dst, working across different file systems/devices.
func (ns *NameShifter) moveFile(src, dst string) error {
	// Open the source file for reading.
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(sourceFile *os.File) {
		err := sourceFile.Close()
		if err != nil {
			fmt.Println("> Error while attempting to close the file: %w", err)
		}
	}(sourceFile)

	// Create the destination file for writing.
	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(destinationFile *os.File) {
		err := destinationFile.Close()
		if err != nil {
			fmt.Println("> Error while attempting to close the file: %w", err)
		}
	}(destinationFile)

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

// moveFileWithRetry handles moving a file from src to dst, working across different file systems/devices, with retries.
func (ns *NameShifter) moveFileWithRetry(src, dst string, maxRetries int) error {
	var lastErr error
	retryDelay := 2 // Initial delay in seconds

	for attempt := 0; attempt < maxRetries; attempt++ {
		err := ns.moveFile(src, dst)
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

func (ns *NameShifter) shouldProcessFile(path string, info os.FileInfo) bool {
	// Immediately return false if it's a directory, no need to check extensions
	if info.IsDir() {
		return false
	}

	// Extract file extension from path for comparison
	fileExt := filepath.Ext(path)

	// Check if the file's extension is in the list of extensions to process
	for _, ext := range ns.Config.FileExtensions {
		if fileExt == ext || (len(ext) > 0 && ext[0] == '.' && fileExt == ext) {
			return true
		}
	}

	return false
}


func (ns *NameShifter) processPath(path string, info os.FileInfo, theStringToBeReplaced, theReplacementString string, cfg *Config) error {
	if err := ns.ignoreConfigDirs(path, nil); err != nil {
		// Uncomment the below if you want the reporter to report failure for skipping config files.
		//row := []table.Row{{"Path", path, "Error", err}}
		//ns.Context.AddError()
		//ns.Context.AddErrorReportRow(row)
		return err
	}

	// Directly use the newly abstracted renameEntity function for files and directories.
	if cfg.WorkGlobally && (info.IsDir() || strings.Contains(info.Name(), theStringToBeReplaced)) {
		if err := ns.renameEntity(path, theStringToBeReplaced, theReplacementString); err != nil {
			row := []table.Row{{"Path", path, "Error", fmt.Sprintf("Could not rename: %v", err)}}
			ns.Context.AddError()
			ns.Context.AddErrorReportRow(row)
			return err
		}
	}

	// For files, check if they should be processed and then process.
	if ns.shouldProcessFile(path, info) {
		return ns.processFile(path, theStringToBeReplaced, theReplacementString)
	}

	return nil
}

func (ns *NameShifter) renameEntity(entityPath, theStringToBeReplaced, theReplacementString string) error {
	// Prepare the new path by replacing the specified string.
	newPath := strings.Replace(entityPath, theStringToBeReplaced, theReplacementString, -1)

	// Attempt to rename (move) the entity.
	if err := ns.moveFileWithRetry(entityPath, newPath, 6); err != nil {
		// Specific handling for permission errors.
		if os.IsPermission(err) {
			// Try changing permissions and retry the move.
			if permErr := os.Chmod(entityPath, 0666); permErr != nil {
				ns.Context.AddError()
				// Wrap the error to provide more context.
				return fmt.Errorf("failed to change permissions for %s: %w", entityPath, permErr)
			}
			// Retry the move operation.
			if retryErr := ns.moveFileWithRetry(entityPath, newPath, 6); retryErr != nil {
				ns.Context.AddError()
				// Wrap the error to provide more context.
				return fmt.Errorf("failed to move %s after changing permissions: %w", entityPath, retryErr)
			}
		} else {
			ns.Context.AddError()
			// Wrap the error to provide more context.
			return fmt.Errorf("failed to move %s: %w", entityPath, err)
		}
	}

	// Log the successful replacement.
	ns.Context.AddReplacement()
	return nil
}

func main() {
	resetColors()
	printLogo()
	flag.Parse()

	customFlagParsing() // Ensure custom flag parsing is called if not already handled by flag.Parse()

	cfg := NewConfig()
	ctx := NewAppContext()
	ns := NewNameShifter(cfg, ctx)

	if cfg.VersionFlag {
		color.Cyan(fmt.Sprintf("\n> NameShifter Version: %s 🚀📚\n", cfg.Version))
		os.Exit(0)
	}

	if len(flag.Args()) < 3 {
		color.Red(fmt.Sprintf("\n> Usage: go run nsh.go <startingDirectory> <theStringToBeReplaced> <theReplacementString> -flags❗📚👀"))
		os.Exit(1)
	}

	args := flag.Args()
	startingDirectory, theStringToBeReplaced, theReplacementString := args[0], args[1], args[2]
	//fmt.Println("> Starting directory:", startingDirectory)
	paths, err := ns.collectPaths(startingDirectory)
	fmt.Println("> Paths:", paths)
	if err != nil {
		fmt.Println("> Error collecting paths:", err)
		os.Exit(1)
	}

	ns.ProcessAllPaths(paths, theStringToBeReplaced, theReplacementString)

	if ctx.errorsCount > 0 {
		ctx.DisplayErrorReport()
	}

	ctx.ReplacementsAndErrorsReport()
	os.Exit(0)
}
