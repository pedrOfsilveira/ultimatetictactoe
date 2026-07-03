package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pedrofsilveira/ultimatetictactoe/internal/game"
	"github.com/spf13/cobra"
)

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Start a new Ultimate Tic-Tac-Toe game",
	Run: func(cmd *cobra.Command, args []string) {
		playGame()
	},
}

func init() {
	rootCmd.AddCommand(playCmd)
}

func playGame() {
	g := game.NewGame("Player 1", "Player 2")

	for g.Status == game.Playing {
		fmt.Println()
		fmt.Println(renderBoard(g))
		fmt.Println("Current turn:", g.CurrentTurn)

		if g.FreeMove {
			fmt.Println("Free move: choose any board")
		} else {
			fmt.Printf("Required board: (%d, %d)\n", g.NextBoardRow, g.NextBoardCol)
		}

		var boardRow int
		var boardCol int
		var cellRow int
		var cellCol int

		if g.FreeMove {
			fmt.Print("Choose board row (0-2): ")
			fmt.Scan(&boardRow)

			fmt.Print("Choose board col (0-2): ")
			fmt.Scan(&boardCol)
		} else {
			boardRow = g.NextBoardRow
			boardCol = g.NextBoardCol
		}

		fmt.Print("Choose cell row (0-2): ")
		fmt.Scan(&cellRow)

		fmt.Print("Choose cell col (0-2): ")
		fmt.Scan(&cellCol)

		clearScreen()

		err := g.MakeMove(boardRow, boardCol, cellRow, cellCol)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
	}

	fmt.Println()
	fmt.Println(g.Board)

	if g.Winner != game.Empty {
		fmt.Println("Winner:", g.Winner)
	} else {
		fmt.Println("Draw!")
	}
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}
