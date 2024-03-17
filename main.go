// Author: Ahmad Hamdi
//  .\nsh.exe "path/to/directory" "OldText" "NewText" -ignore-config-dirs=true -work-globally=false -concurrent-run=false -case-matching=true -file-extensions=".go,.md"

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	ignoreConfig      bool
	workGlobally      bool
	concurrentRun     bool
	caseMatching      bool
	Version           = "0.1.2"
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
		errorsCount++
		return err
	}
	return nil
}

func processPath(path string, info os.FileInfo, theStringToBeReplaced, theReplacementString string) error {
	if err := ignoreConfigDirs(path, nil); err != nil {
		row := []table.Row{{1, fmt.Sprintf("%v 🚨", path), err}}
		addRowTo(errorReport, row)
		return nil
	}

	// If global, rename entities first.
	if workGlobally && info.IsDir() || strings.Contains(info.Name(), theStringToBeReplaced) && workGlobally {
		err := renameEntities(path, theStringToBeReplaced, theReplacementString)
		if err != nil {
			row := []table.Row{{1, fmt.Sprintf("Could not rename this directory: %s 📕🚫", path), err}}
			addRowTo(errorReport, row)
			return nil
		}
	}

	// If fileExtensions are provided
	if fileExtensions != "" {
		exts := strings.Split(fileExtensions, ",")
		match := false
		for _, ext := range exts {
			if strings.HasSuffix(info.Name(), ext) {
				match = true
				break
			}
		}
		if !match {
			return nil // Skip this file if not a single match was found!
		}
	}

	// Now rename file contents
	if !info.IsDir() {
		readFile, err := os.Open(path)
		if err != nil {
			errorsCount++
			row := []table.Row{{3, fmt.Sprintf("Error opening %s 🚨", path), err}}
			addRowTo(errorReport, row)
			return nil
		}
		defer func(readFile *os.File) {
			err := readFile.Close()
			if err != nil {
				errorsCount++
				row := []table.Row{{3, fmt.Sprintf("Could not save the file 🚨"), err}}
				addRowTo(errorReport, row)
			}
		}(readFile)

		scanner := bufio.NewScanner(readFile)
		var txt []string
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, theStringToBeReplaced) {
				line = strings.Replace(line, theStringToBeReplaced, theReplacementString, -1)
				replacementsCount++
			}
			txt = append(txt, line)
		}

		if err := scanner.Err(); err != nil {
			errorsCount++
			row := []table.Row{{3, fmt.Sprintf("Error processing %s 🚨", path), err}}
			addRowTo(errorReport, row)
			return nil
		}

		writeFile, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, info.Mode())
		if err != nil {
			errorsCount++
			row := []table.Row{{3, fmt.Sprintf("Error opening %s 🚨", path), err}}
			addRowTo(errorReport, row)
			return nil
		}
		defer func(writeFile *os.File) {
			err := writeFile.Close()
			if err != nil {
				errorsCount++
				row := []table.Row{{3, fmt.Sprintf("Error saving file 🚨"), err}}
				addRowTo(errorReport, row)
			}
		}(writeFile)

		for _, line := range txt {
			_, err := writeFile.WriteString(line + "\n")
			if err != nil {
				errorsCount++
				row := []table.Row{{3, fmt.Sprintf("Error processing %s 🚨", path), err}}
				addRowTo(errorReport, row)
				return nil
			}
		}
	}

	resetColors()
	return nil
}

func renameEntities(startingDirectory, theStringToBeReplaced, theReplacementString string) error {
	var dirs []string
	var files []string

	// First, accumulate directories and files
	err := filepath.Walk(startingDirectory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			//log.Printf("> renameEntities > appending dir: %s\n", path)
			dirs = append(dirs, path)
		} else {
			//log.Printf("> renameEntities > appending file: %s\n", path)
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		row := []table.Row{{3, err}}
		addRowTo(errorReport, row)
		errorsCount++
		return nil
	}

	var (
		newPath string
		newName string
	)

	for _, dir := range dirs {

		err := os.Chmod(dir, 0666)
		if err != nil {
			errorsCount++
			row := []table.Row{{3, "Failed to set directory permissions", err}}
			addRowTo(errorReport, row)
		}

		// Isolate the directory name from its path
		dirName := filepath.Base(dir)
		parentDir := filepath.Dir(dir)
		//grandDir := filepath.Dir(parentDir)

		color.Green("\n> Dir => %s\n", dirName)
		color.Green("\n> Parent Dir => %s\n", parentDir)

		// Replace only in the directory name, not the entire path
		newName = strings.Replace(dirName, theStringToBeReplaced, theReplacementString, -1)
		newPath = filepath.Join(parentDir, newName)

		if err := os.Rename(dir, newPath); err != nil {
			errorsCount++
			row := []table.Row{{3, fmt.Sprintf("Error renaming %s to %s 😵💔", dir, newPath), err}}
			addRowTo(errorReport, row)
			return err
		}
		replacementsCount++
	}

	// Rename files
	for _, file := range files {
		err := os.Chmod(file, 0666)
		if err != nil {
			errorsCount++
			row := []table.Row{{3, "Failed to set file permissions", err}}
			addRowTo(errorReport, row)
		}

		newPath = filepath.Join(filepath.Dir(file), newName)

		if err := os.Rename(file, newPath); err != nil {
			errorsCount++
			row := []table.Row{{3, fmt.Sprintf("Error renaming %s to %s 😵💔", file, newPath), err}}
			addRowTo(errorReport, row)
			return err
		}
		replacementsCount++
	}

	resetColors()
	return nil
}

func main() {
	resetColors()
	printLogo()
	flag.Parse()

	if versionFlag {
		color.Cyan(fmt.Sprintf("\n> NameShifter Version: %s 🚀📚\n", Version))
		os.Exit(0)
	}

	errorReport = table.NewWriter()
	if len(flag.Args()) < 3 {
		color.Red(fmt.Sprintf("\n> Usage: go run nsh.go <startingDirectory> <theStringToBeReplaced> <theReplacementString> -flags❗📚👀"))
		os.Exit(1)
	}

	args := flag.Args()
	startingDirectory, theStringToBeReplaced, theReplacementString := args[0], args[1], args[2]

	printSettings()
	var wg sync.WaitGroup

	err := filepath.Walk(startingDirectory, func(path string, info os.FileInfo, err error) error {
		if concurrentRun && info.IsDir() {
			wg.Add(1)
			go func(path string, info os.FileInfo) {
				defer wg.Done()
				if err := processPath(path, info, theStringToBeReplaced, theReplacementString); err != nil {
					row := []table.Row{{4, fmt.Sprintf("🚨"), err}}
					addRowTo(errorReport, row)
				}
			}(path, info)
			return filepath.SkipDir
		} else {
			return processPath(path, info, theStringToBeReplaced, theReplacementString)
		}
	})

	wg.Wait()

	if err != nil {
		row := []table.Row{{3, fmt.Sprintf("Error walking through %s 😢👣", startingDirectory), err}}
		addRowTo(errorReport, row)
		//os.Exit(2)
	} else {
		color.Green(fmt.Sprintf("\n> Names Shifted Successfully! 🎉📝✅\n"))
	}

	displayErrorReport()
	replacementsAndErrorsReport()
	os.Exit(0)
}
