package main

import (
	"fmt"

	"github.com/pedrofsilveira/ultimatetictactoe/internal/game"
)

func main() {
	g := game.NewGame("Player 1", "Player 2")

	for g.Status == game.Playing {
		fmt.Println()
		fmt.Println(&g.Board)
		fmt.Println("Current turn:", g.CurrentTurn)

		var row int
		var col int

		fmt.Print("Choose row (0-2): ")
		fmt.Scan(&row)

		fmt.Print("Choose col (0-2): ")
		fmt.Scan(&col)

		err := g.MakeMove(row, col)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
	}

	fmt.Println()
	fmt.Println(&g.Board)

	if g.Winner != game.Empty {
		fmt.Println("Winner:", g.Winner)
	} else {
		fmt.Println("Draw!")
	}
}
