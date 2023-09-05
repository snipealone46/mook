package cmd

import (
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
)

type Color = color.Color

func LivePrint(lines []string) {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	ColorPrintLines(lines)

	time.Sleep(2 * time.Second)
}

func ColorPrintLines(lines []string) {
	colors := [...]*Color{
		color.New(color.FgCyan),
		color.New(color.FgYellow),
		color.New(color.FgMagenta),
		color.New(color.FgWhite),
		color.New(color.FgGreen),
		color.New(color.FgBlue)}

	for index := 0; index < len(lines); index++ {
		colorIndex := index % len(colors)
		colors[colorIndex].Println(lines[index] + " ")
	}
}
