package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"math/rand"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

type AppContext struct {
	errorsCount       int32
	replacementsCount int32
	errorReport       table.Writer
	mutex             sync.Mutex // Protects errorReport updates.
}

func NewAppContext() *AppContext {
	return &AppContext{
		errorReport: table.NewWriter(),
	}
}

func (ctx *AppContext) AddError() {
	atomic.AddInt32(&ctx.errorsCount, 1)
}

func (ctx *AppContext) AddReplacement() {
	atomic.AddInt32(&ctx.replacementsCount, 1)
}

func (ctx *AppContext) AddErrorReportRow(row []table.Row) {
	ctx.mutex.Lock()
	ctx.errorReport.AppendRows(row)
	ctx.mutex.Unlock()
}

func (ctx *AppContext) DisplayErrorReport() {
	ctx.errorReport.SetOutputMirror(os.Stdout)
	header := table.Row{"#", "Directory", "Error Details"}
	ctx.errorReport.AppendHeader(header)
	ctx.errorReport.AppendFooter(table.Row{"Error", "Report", "Done"})
	ctx.errorReport = formatColumn(ctx.errorReport, header)
	ctx.errorReport.SetStyle(table.StyleColoredRedWhiteOnBlack)
	ctx.errorReport.Render()
	fmt.Println("")
	resetColors() // Assuming resetColors is a function that resets terminal color settings.
}

func (ctx *AppContext) ReplacementsAndErrorsReport() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	header := table.Row{"#", "Replacements Made", "Errors Encountered"}
	t.AppendHeader(header)
	t.AppendRows([]table.Row{
		{1, atomic.LoadInt32(&ctx.replacementsCount), atomic.LoadInt32(&ctx.errorsCount)},
	})
	t.AppendSeparator()
	t.AppendFooter(table.Row{">", "Shifted", ""})

	t = formatColumn(t, header)
	t.SetStyle(table.StyleColoredBlackOnYellowWhite)
	t.Render()
	fmt.Println("")
	resetColors() // Assuming resetColors is a function that resets terminal color settings.
}

func resetColors() {
	reset := color.New(color.Reset).SprintFunc()
	fmt.Printf(reset(""))
}

func printLogo() {
	// Construct a list of logos.
	// Pick a random one and display it.

	teal := color.New(color.FgCyan).SprintFunc()
	logo1 := `

'##::: ##::'######::'##::::'##:
 ###:: ##:'##... ##: ##:::: ##:
 ####: ##: ##:::..:: ##:::: ##:
 ## ## ##:. ######:: #########:
 ##. ####::..... ##: ##.... ##:
 ##:. ###:'##::: ##: ##:::: ##:
 ##::. ##:. ######:: ##:::: ##:
..::::..:::......:::..:::::..::

`

	logo2 := `

                      
 #    #  ####  #    # 
 ##   # #      #    # 
 # #  #  ####  ###### 
 #  # #      # #    # 
 #   ## #    # #    # 
 #    #  ####  #    # 

`

	logo3 := `

███╗   ██╗███████╗██╗  ██╗
████╗  ██║██╔════╝██║  ██║
██╔██╗ ██║███████╗███████║
██║╚██╗██║╚════██║██╔══██║
██║ ╚████║███████║██║  ██║
╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝

`

	logo4 := `

  __   __  ______   __   __   
 /_/\ /\_\/ ____/\ /\_\ /_/\  
 ) ) \ ( () ) __\/( ( (_) ) ) 
/_/   \ \_\\ \ \   \ \___/ /  
\ \ \   / /_\ \ \  / / _ \ \  
 )_) \ (_()____) )( (_( )_) ) 
 \_\/ \/_/\____\/  \/_/ \_\/  

`

	logo5 := `

 _        _______          
( (    /|(  ____ \|\     /|
|  \  ( || (    \/| )   ( |
|   \ | || (_____ | (___) |
| (\ \) |(_____  )|  ___  |
| | \   |      ) || (   ) |
| )  \  |/\____) || )   ( |
|/    )_)\_______)|/     \|

`

	logo6 := `

░▒▓███████▓▒░ ░▒▓███████▓▒░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓█▓▒░░▒▓█▓▒░░▒▓██████▓▒░░▒▓████████▓▒░ 
░▒▓█▓▒░░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓█▓▒░░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓█▓▒░░▒▓█▓▒░▒▓███████▓▒░░▒▓█▓▒░░▒▓█▓▒░ 

`
	logos := []string{logo1, logo2, logo3, logo4, logo5, logo6}
	index := rand.Intn(len(logos))
	fmt.Println(teal(logos[index]))
}

func printSettings() {
	//log.Println("> Entering printSettings")
	magenta := color.New(color.FgHiMagenta).SprintFunc()
	fmt.Println(magenta(fmt.Sprintf("\n> nsh called with the following arguments:\n")))
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	header := table.Row{"#", "Argument Name", "Argument Value"}
	t.AppendHeader(header)
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

// formatColumn configures column properties based on provided column names within the table.Row.
func formatColumn(t table.Writer, row table.Row) table.Writer {
	nameTransformer := text.Transformer(func(val interface{}) string {
		// Ensure conversion to string for the transformation, as val is of type interface{}.
		return text.Bold.Sprint(val)
	})

	columnConfigs := make([]table.ColumnConfig, len(row))

	for i, column := range row {
		width := 32
		columnName, ok := column.(string)
		if !ok {
			// If the column isn't a string, skip configuration.
			// Alternatively, you could convert to string or handle differently.
			continue
		}

		if columnName == "#" {
			width = 8
		}

		columnConfigs[i] = table.ColumnConfig{
			Name:              columnName,
			Align:             text.AlignLeft,
			AlignFooter:       text.AlignLeft,
			AlignHeader:       text.AlignLeft,
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

	(t).SetColumnConfigs(columnConfigs)

	return t
}
