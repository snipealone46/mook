package pkg

import (
	"github.com/acarl005/stripansi"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"os"
	"os/exec"
	"unicode/utf8"
)

type Color = color.Color

func ColorPrintLines(lines [][]string) {
	ClearScreen()
	columnHeaderList := []interface{}{"PodName", "PodState", "RestartCount", "Age", "Ready?", "ExtraInfo"}
	tbl := InitializeTable(columnHeaderList)

	colors := [...]*Color{
		color.New(color.FgCyan),
		color.New(color.FgYellow),
		color.New(color.FgMagenta),
		color.New(color.FgWhite),
		color.New(color.FgGreen),
		color.New(color.FgBlue)}

	for index := 0; index < len(lines); index++ {
		row := lines[index]

		coloredRow := make([]interface{}, len(row))
		colorIndex := index % len(colors)
		for rowIndex := 0; rowIndex < len(row); rowIndex++ {
			coloredRow[rowIndex] = colors[colorIndex].Sprint(row[rowIndex])
		}
		tbl.AddRow(coloredRow...)
	}
	tbl.Print()
}

func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func InitializeTable(columnHeaderList []interface{}) table.Table {
	headerFmt := color.New(color.Underline).SprintfFunc()

	tbl := table.New(columnHeaderList...)
	tbl.WithHeaderFormatter(headerFmt).WithWidthFunc(func(s string) int {
		return utf8.RuneCountInString(stripansi.Strip(s))
	})
	return tbl
}
