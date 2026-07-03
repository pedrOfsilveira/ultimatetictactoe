package cli

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/pedrofsilveira/ultimatetictactoe/internal/game"
)

var (
	xColor        = color.New(color.FgHiRed).SprintFunc()
	oColor        = color.New(color.FgHiBlue).SprintFunc()
	requiredColor = color.New(color.FgHiCyan).SprintFunc()
	wonColor      = color.New(color.FgHiGreen).SprintFunc()
)

func renderBoard(g *game.Game) string {
	var result strings.Builder

	for boardRow := range 3 {
		left := renderSmallBoard(g, boardRow, 0)
		middle := renderSmallBoard(g, boardRow, 1)
		right := renderSmallBoard(g, boardRow, 2)

		for line := range len(left) {
			result.WriteString(left[line])
			result.WriteString(" || ")
			result.WriteString(middle[line])
			result.WriteString(" || ")
			result.WriteString(right[line])
			result.WriteString("\n")
		}

		if boardRow < 2 {
			result.WriteString("=========================================\n")
		}
	}

	return result.String()
}

func renderSmallBoard(g *game.Game, boardRow, boardCol int) []string {
	smallBoard := g.Board.Boards[boardRow][boardCol]

	lines := smallBoardLines(smallBoard)

	if !g.FreeMove && boardRow == g.NextBoardRow && boardCol == g.NextBoardCol {
		for i := range lines {
			lines[i] = highlightBoardBorders(lines[i])
		}
	}

	return lines
}

func smallBoardLines(b game.SmallBoard) []string {
	if b.Status == game.Finished && b.Winner != game.Empty {
		return []string{
			"           ",
			wonColor(fmt.Sprintf("     %s     ", b.Winner)),
			"           ",
			"-----------",
			wonColor("    WON    "),
		}
	}

	return []string{
		fmt.Sprintf(" %s | %s | %s ", renderTeam(b.Cells[0][0]), renderTeam(b.Cells[0][1]), renderTeam(b.Cells[0][2])),
		"-----------",
		fmt.Sprintf(" %s | %s | %s ", renderTeam(b.Cells[1][0]), renderTeam(b.Cells[1][1]), renderTeam(b.Cells[1][2])),
		"-----------",
		fmt.Sprintf(" %s | %s | %s ", renderTeam(b.Cells[2][0]), renderTeam(b.Cells[2][1]), renderTeam(b.Cells[2][2])),
	}
}

func renderTeam(team game.Team) string {
	switch team {
	case game.X:
		return xColor("X")
	case game.O:
		return oColor("O")
	default:
		return " "
	}
}

func highlightBoardBorders(line string) string {
	var result strings.Builder

	for _, char := range line {
		switch char {
		case '|', '-':
			result.WriteString(requiredColor(string(char)))
		default:
			result.WriteRune(char)
		}
	}

	return result.String()
}
