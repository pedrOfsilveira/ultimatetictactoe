package main

import (
	"os"
	"syscall"

	"github.com/fatih/color"
	"github.com/pedrofsilveira/ultimatetictactoe/internal/cli"
)

var allocConsole = syscall.NewLazyDLL("kernel32.dll").NewProc("AllocConsole")

func main() {
	result, _, _ := allocConsole.Call()
	if result == 0 {
		return
	}

	consoleIn, err := os.OpenFile("CONIN$", os.O_RDWR, 0)
	if err != nil {
		return
	}
	defer consoleIn.Close()

	consoleOut, err := os.OpenFile("CONOUT$", os.O_RDWR, 0)
	if err != nil {
		return
	}
	defer consoleOut.Close()

	os.Stdin = consoleIn
	os.Stdout = consoleOut
	os.Stderr = consoleOut
	color.NoColor = false

	cli.Execute()
}
