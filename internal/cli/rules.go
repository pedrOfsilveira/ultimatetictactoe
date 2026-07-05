package cli

import "fmt"

func printRules() {
	fmt.Println(`Ultimate Tic-Tac-Toe Rules

- Win small boards by getting 3 in a row.
- Win the game by winning 3 small boards in a row.
- Your cell move sends the opponent to the matching board.
- If that board is finished, the opponent gets a free move.`)
}
