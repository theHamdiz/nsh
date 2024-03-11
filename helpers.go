package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"os"
	"strings"
)

func resetColors() {
	reset := color.New(color.Reset).SprintFunc()
	fmt.Printf(reset(""))
}

func printSettings() {
	//log.Println("> Entering printSettings")
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
	//log.Println("> Exiting printSettings")
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

func replacementsAndErrorsReport() {
	//log.Println("> Entering the replacementsAndErrorsReport")
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
	//log.Println("> Exiting the replacementsAndErrorsReport")
}

func addRowTo(table table.Writer, r []table.Row) table.Writer {
	table.AppendRows(r)
	table.AppendSeparator()
	return table
}
