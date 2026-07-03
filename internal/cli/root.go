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
  play    Start a local game
  local   Start a local game
  host    Host a LAN game on port 8080
  join    Join a LAN game (join IP:PORT)
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

		case "local":
			playGame()

		case "host":
			if err := HostGame("8080"); err != nil {
				fmt.Println("Host error:", err)
			}

		default:
			if strings.HasPrefix(input, "join ") {
				address := strings.TrimSpace(strings.TrimPrefix(input, "join "))
				if err := JoinGame(address); err != nil {
					fmt.Println("Join error:", err)
				}
				continue
			}
			fmt.Println("Unknown command:", input)
			fmt.Println("Type 'help' to see available commands.")
			continue

		case "rules":
			printRules()

		case "help":
			printInteractiveHelp()

		case "exit", "quit":
			fmt.Println("Goodbye!")
			return

		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}

func printInteractiveHelp() {
	fmt.Println(`
Available commands:
  play    Start a local game
  local   Start a local game
  host    Host a LAN game on port 8080
  join    Join a LAN game (join IP:PORT)
  rules   Show the rules
  help    Show this message
  exit    Close the program`)
}
