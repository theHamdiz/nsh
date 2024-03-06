// Author: Ahmad Hamdi

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	ignoreConfig  bool
	workGlobally  bool
	concurrentRun bool
	caseMatching  bool
)

func printLogo() {
	teal := color.New(color.FgCyan).SprintFunc()
	multiLineString := `
░▒▓███████▓▒░ ░▒▓██████▓▒░░▒▓██████████████▓▒░░▒▓████████▓▒░░▒▓███████▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓████████▓▒░▒▓████████▓▒░ 
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░         ░▒▓█▓▒░     
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░         ░▒▓█▓▒░     
░▒▓█▓▒░░▒▓█▓▒░▒▓████████▓▒░▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓██████▓▒░  ░▒▓██████▓▒░░▒▓████████▓▒░▒▓█▓▒░▒▓██████▓▒░    ░▒▓█▓▒░     
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░             ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░         ░▒▓█▓▒░     
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░             ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░         ░▒▓█▓▒░     
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓████████▓▒░▒▓███████▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░         ░▒▓█▓▒░
`

	fmt.Println(teal(multiLineString))
}

func customFlagParsing() {
	for i, arg := range os.Args {
		if strings.HasPrefix(arg, "--") {
			os.Args[i] = "-" + arg
		}
	}
}

func init() {
	customFlagParsing()
	flag.BoolVar(&ignoreConfig, "ignore-config-dirs", true, "Ignore .config directories 🚫🐙")
	flag.BoolVar(&workGlobally, "work-globally", true, "Work on folder names, file names, and file contents 🌍✨")
	flag.BoolVar(&concurrentRun, "concurrent-run", true, "Run each folder inside the root directory in a separate goroutine 🏃💨")
	flag.BoolVar(&caseMatching, "case-matching", true, "Match case when replacing strings (default true) 👔🔍")
}
func ignoreConfigDirs(path string, info os.FileInfo, err error) error {
	if err != nil {
		if os.IsPermission(err) {
			return nil // Skip this file or directory but continue walking
		}
		return err
	}
	dirName := filepath.Base(path)
	if info.IsDir() && strings.HasPrefix(dirName, ".") {
		// If the directory name starts with '.', skip it
		return filepath.SkipDir
	}
	return nil
}

func processPath(path string, info os.FileInfo, theStringToBeReplaced, theReplacementString string) error {
	// If global, rename entities first.
	if workGlobally && info.IsDir() || strings.Contains(info.Name(), theStringToBeReplaced) {
		err := renameEntities(path, theStringToBeReplaced, theReplacementString)
		if err != nil {
			return err
		}
	}

	// Now rename file contents
	if !info.IsDir() {
		readFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func(readFile *os.File) {
			err := readFile.Close()
			if err != nil {
				color.Red("Could not save this file: %s 📕🚫", err)
			}
		}(readFile)

		scanner := bufio.NewScanner(readFile)
		var text []string
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, theStringToBeReplaced) {
				line = strings.Replace(line, theStringToBeReplaced, theReplacementString, -1)
			}
			text = append(text, line)
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		writeFile, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, info.Mode())
		if err != nil {
			return err
		}
		defer func(writeFile *os.File) {
			err := writeFile.Close()
			if err != nil {
				color.Red("Could not save this file: %s 🚷📝", err)
			}
		}(writeFile)

		for _, line := range text {
			_, err := writeFile.WriteString(line + "\n")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func renameEntities(startingDirectory, theStringToBeReplaced, theReplacementString string) error {
	var dirs []string
	var files []string

	// First, accumulate directories and files
	err := filepath.Walk(startingDirectory, func(path string, info os.FileInfo, err error) error {
		if err := ignoreConfigDirs(path, info, nil); err != nil {
			return err
		}
		if info.IsDir() {
			dirs = append(dirs, path)
		} else {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	var (
		newPath string
		newName string
	)

	for _, dir := range dirs {
		err := os.Chmod(dir, 0666)
		if err != nil {
			log.Fatalf("Failed to set directory permissions: %v", err)
		}

		// Isolate the directory name from its path
		dirName := filepath.Base(dir)
		currentDir := filepath.Dir(dir)
		parentDir := filepath.Dir(currentDir)

		color.Green("Dir => %s\n", dirName)
		color.Green("Parent Dir => %s\n", parentDir)

		// Replace only in the directory name, not the entire path
		newName = strings.Replace(dirName, theStringToBeReplaced, theReplacementString, -1)
		newPath = filepath.Join(parentDir, newName)

		if err := os.Rename(dir, newPath); err != nil {
			return fmt.Errorf("error renaming %s to %s: %w 😵💔", dir, newPath, err)
		}
	}

	// Rename files
	for _, file := range files {
		err := os.Chmod(file, 0666)
		if err != nil {
			log.Fatalf("Failed to set file permissions: %v", err)
		}

		newPath = filepath.Join(filepath.Dir(file), newName)

		if err := os.Rename(file, newPath); err != nil {
			return fmt.Errorf("error renaming %s to %s: %w 😵💔", file, newPath, err)
		}
	}

	return nil
}

func main() {
	printLogo()
	flag.Parse()

	if len(flag.Args()) != 3 {
		color.Red("Usage: go run script.go -flags <startingDirectory> <theStringToBeReplaced> <theReplacementString> ❗📚👀")
		os.Exit(1)
	}

	args := flag.Args()
	startingDirectory, theStringToBeReplaced, theReplacementString := args[0], args[1], args[2]

	var wg sync.WaitGroup

	err := filepath.Walk(startingDirectory, func(path string, info os.FileInfo, err error) error {
		if concurrentRun && info.IsDir() {
			wg.Add(1)
			go func(path string, info os.FileInfo) {
				defer wg.Done()
				if err := processPath(path, info, theStringToBeReplaced, theReplacementString); err != nil {
					color.Red("%s\n=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-= 🐛💥", err)
				}
			}(path, info)
			return filepath.SkipDir
		} else {
			return processPath(path, info, theStringToBeReplaced, theReplacementString)
		}
	})

	wg.Wait()

	if err != nil {
		color.Red("Error walking through %s: %v 😢👣", startingDirectory, err)
		os.Exit(2)
	} else {
		color.Green("Names Shifted Successfully! 🎉📝✅")
	}
}
