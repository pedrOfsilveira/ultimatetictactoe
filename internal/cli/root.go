package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "uttt",
	Short: "Ultimate Tic-Tac-Toe CLI",
	Long:  "UTTT - Ultimate Tic-Tac-Toe CLI",
	Run: func(cmd *cobra.Command, args []string) {
		startInteractiveMode()
	},
}

func Execute() {
	// The app supports being opened from Windows Explorer through the launcher.
	// Cobra's default mousetrap warning would incorrectly reject that launch.
	cobra.MousetrapHelpText = ""

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func startInteractiveMode() {
	fmt.Println(`
THIS IS UTTT - ULTIMATE TIC-TAC-TOE

Available commands:
  play    Start a new game
  rules   Show the rules
  help    Show this message
  exit    Close the program`)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("uttt> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		switch input {
		case "":
			continue

		case "play":
			playGame()

		case "rules":
			printRules()

		case "help":
			printInteractiveHelp()

		case "exit", "quit":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Unknown command:", input)
			fmt.Println("Type 'help' to see available commands.")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}

func printInteractiveHelp() {
	fmt.Println(`
Available commands:
  play    Start a new game
  rules   Show the rules
  help    Show this message
  exit    Close the program`)
}
