// Author: Ahmad Hamdi
//  .\nameShift.exe "path/to/directory" "OldText" "NewText" -ignore-config-dirs=true -work-globally=false -concurrent-run=false -case-matching=true -file-extensions=".go,.md"

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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
	fileExtensions    string
	errorsCount       int32
	replacementsCount int32
	errorReport       table.Writer
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

func resetColors() {
	reset := color.New(color.Reset).SprintFunc()
	fmt.Printf(reset(""))
}

func printSettings() {
	magenta := color.New(color.FgHiMagenta).SprintFunc()
	fmt.Println(magenta(fmt.Sprintf("\n> nameShifter called with the following arguments:\n")))
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	header := []string{"#", "Argument Name", "Argument Value"}
	t.AppendHeader(table.Row{"#", "Argument Name", "Argument Value"})
	for i := range os.Args {
		argName := ""
		if i == 0 {
			argName = "binaryPath"
		} else if i == 1 {
			argName = "startingDir"
		} else if i == 2 {
			argName = "stringToBeReplaced"
		} else if i == 3 {
			argName = "replacementString"
		} else {
			argName = strings.Split(os.Args[i], "=")[0]
		}
		t.AppendRows([]table.Row{
			{i, argName, os.Args[i]},
		})
		t.AppendSeparator()
	}

	t.AppendFooter(table.Row{"", "", "Total", len(os.Args)})
	t = formatColumn(t, header)
	t.SetStyle(table.StyleColoredMagentaWhiteOnBlack)
	t.Render()
	fmt.Println("")
	resetColors()
}

func customFlagParsing() {
	for i, arg := range os.Args {
		if strings.HasPrefix(arg, "--") {
			os.Args[i] = "-" + arg
		}
	}
}

func replacementsAndErrorsReport() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	header := []string{"#", "Replacements Made", "Errors Encountered"}
	t.AppendHeader(table.Row{"#", "Replacements Made", "Errors Encountered"})
	t.AppendRows([]table.Row{
		{1, replacementsCount, errorsCount},
	})
	t.AppendSeparator()
	t.AppendFooter(table.Row{">", "Shifted", ""})

	t = formatColumn(t, header)
	t.SetStyle(table.StyleColoredBlackOnYellowWhite)
	t.Render()
	fmt.Println("")
	resetColors()
}

func addRowTo(table table.Writer, r []table.Row) table.Writer {
	table.AppendRows(r)
	table.AppendSeparator()
	return table
}

func init() {
	customFlagParsing()
	flag.BoolVar(&ignoreConfig, "ignore-config-dirs", true, "Ignore .config directories 🚫🐙")
	flag.BoolVar(&workGlobally, "work-globally", false, "Work on folder names, file names, and file contents (default false)🌍✨")
	flag.BoolVar(&concurrentRun, "concurrent-run", false, "Run each folder inside the root directory in a separate goroutine (default false)🏃💨")
	flag.BoolVar(&caseMatching, "case-matching", true, "Match case when replacing strings (default true) 👔🔍")
	flag.StringVar(&fileExtensions, "file-extensions", ".go", "Comma-separated list of file extensions to process, e.g., '.go,.txt' 📄✂️")
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
	if workGlobally && info.IsDir() || strings.Contains(info.Name(), theStringToBeReplaced) {
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
			dirs = append(dirs, path)
		} else {
			files = append(files, path)
		}
		replacementsCount++
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
		currentDir := filepath.Dir(dir)
		parentDir := filepath.Dir(currentDir)

		color.Green("\n> Dir => %s\n", dirName)
		color.Green("\n> Parent Dir => %s\n", parentDir)

		// Replace only in the directory name, not the entire path
		newName = strings.Replace(dirName, theStringToBeReplaced, theReplacementString, -1)
		newPath = filepath.Join(parentDir, newName)

		if err := os.Rename(dir, newPath); err != nil {
			errorsCount++
			row := []table.Row{{3, fmt.Errorf("\n> error renaming %s to %s 😵💔", dir, newPath), err}}
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
			row := []table.Row{{3, fmt.Errorf("\n> error renaming %s to %s 😵💔", file, newPath), err}}
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
	errorReport = table.NewWriter()
	if len(flag.Args()) < 3 {
		color.Red(fmt.Sprintf("\n> Usage: go run script.go <startingDirectory> <theStringToBeReplaced> <theReplacementString> -flags❗📚👀"))
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

func displayErrorReport() {
	errorReport.SetOutputMirror(os.Stdout)
	header := []string{"#", "Directory", "Error Details"}
	errorReport.AppendHeader(table.Row{"#", "Directory", "Error Details"})
	errorReport.AppendFooter(table.Row{"Error", "Report", "Done"})
	errorReport = formatColumn(errorReport, header)
	errorReport.SetStyle(table.StyleColoredRedWhiteOnBlack)
	errorReport.Render()
	fmt.Println("")
	resetColors()
}

func formatColumn(t table.Writer, columns []string) table.Writer {
	nameTransformer := text.Transformer(func(val interface{}) string {
		return text.Bold.Sprint(val)
	})

	count := len(columns)
	columnConfigs := make([]table.ColumnConfig, count)

	for i := 0; i < count; i++ {
		width := 32
		if columns[i] == "#" {
			width = 8
		} else {
			width = 32
		}
		columnConfigs[i] = table.ColumnConfig{
			Name:        columns[i],
			Align:       text.AlignLeft,
			AlignFooter: text.AlignLeft,
			AlignHeader: text.AlignLeft,
			//Colors:            text.Colors{text.BgBlack, text.FgRed},
			//ColorsHeader:      text.Colors{text.BgRed, text.FgBlack, text.Bold},
			//ColorsFooter:      text.Colors{text.BgRed, text.FgBlack},
			Hidden:            false,
			Transformer:       nameTransformer,
			TransformerFooter: nameTransformer,
			TransformerHeader: nameTransformer,
			VAlign:            text.VAlignMiddle,
			VAlignFooter:      text.VAlignTop,
			VAlignHeader:      text.VAlignBottom,
			WidthMin:          width,
			WidthMax:          width,
		}
	}

	t.SetColumnConfigs(columnConfigs)

	return t
}
