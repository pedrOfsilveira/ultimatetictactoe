package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "Show the game rules",
	Run: func(cmd *cobra.Command, args []string) {
		printRules()
	},
}

func init() {
	rootCmd.AddCommand(rulesCmd)
}

func printRules() {
	fmt.Println(`
Ultimate Tic-Tac-Toe Rules

- Win small boards by getting 3 in a row.
- Win the game by winning 3 small boards in a row.
- Your cell move sends the opponent to the matching board.
- If that board is finished, the opponent gets a free move.`)
}
