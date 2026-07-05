package cli

import (
	"flag"
	"fmt"
	"os"
)

func Execute() {
	args := os.Args[1:]
	if len(args) == 0 {
		runInteractiveMenu()
		return
	}

	var err error
	switch args[0] {
	case "play", "local":
		err = playGame()
	case "host":
		flags := flag.NewFlagSet("host", flag.ContinueOnError)
		port := flags.String("port", "8080", "TCP port to listen on")
		if err = flags.Parse(args[1:]); err == nil {
			err = HostGame(*port)
		}
	case "join":
		if len(args) != 2 {
			err = fmt.Errorf("usage: uttt join ADDRESS")
		} else {
			err = JoinGame(args[1])
		}
	case "rules":
		printRules()
	case "help", "-h", "--help":
		printHelp()
	default:
		err = fmt.Errorf("unknown command %q\n\nRun 'uttt help' for usage", args[0])
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func runInteractiveMenu() {
	for {
		menu, err := runMenu()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			return
		}

		switch menu.action {
		case menuLocal:
			err = playGame()
		case menuHost:
			err = HostGame("8080")
		case menuJoin:
			err = JoinGame(menu.address)
		case menuQuit, menuNone:
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
		}
	}
}

func printHelp() {
	fmt.Println(`UTTT - Ultimate Tic-Tac-Toe

Usage:
  uttt                 Open the interactive menu
  uttt play|local      Start a local game
  uttt host [--port]   Host a LAN game
  uttt join ADDRESS    Join a LAN game
  uttt rules           Show the rules

Controls: arrow keys move, Enter selects/plays, Esc goes back, q quits.`)
}
