package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/pedrofsilveira/ultimatetictactoe/internal/game"
)

var (
	xColor        = color.New(color.FgHiRed).SprintFunc()
	oColor        = color.New(color.FgHiBlue).SprintFunc()
	requiredColor = color.New(color.FgHiCyan).SprintFunc()
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

// renderInteractiveBoard keeps every cell the same printable width while using
// color to distinguish the selected board and reverse video for its cursor.
func renderInteractiveBoard(g *game.Game, selectedBoardRow, selectedBoardCol, selectedCellRow, selectedCellCol int, selectingCell, cursorVisible bool) string {
	var result strings.Builder

	for boardRow := range 3 {
		boards := make([][]string, 3)
		for boardCol := range 3 {
			board := g.Board.Boards[boardRow][boardCol]
			lines := smallBoardLinesInteractive(board, g.CurrentTurn, boardRow, boardCol, selectedBoardRow, selectedBoardCol, selectedCellRow, selectedCellCol, selectingCell, cursorVisible)
			if boardRow == selectedBoardRow && boardCol == selectedBoardCol && g.Status == game.Playing {
				for i := range lines {
					// Color only the grid. Wrapping the complete line in a style
					// would override the nested player-colored cell cursor.
					lines[i] = highlightBoardBorders(lines[i])
				}
			}
			boards[boardCol] = lines
		}

		for line := range 5 {
			result.WriteString(boards[0][line])
			result.WriteString(" || ")
			result.WriteString(boards[1][line])
			result.WriteString(" || ")
			result.WriteString(boards[2][line])
			result.WriteByte('\n')
		}
		if boardRow < 2 {
			result.WriteString("=========================================\n")
		}
	}
	return result.String()
}

func smallBoardLinesInteractive(b game.SmallBoard, currentTurn game.Team, boardRow, boardCol, selectedBoardRow, selectedBoardCol, selectedCellRow, selectedCellCol int, selectingCell, cursorVisible bool) []string {
	if b.Status == game.Finished && b.Winner != game.Empty {
		return smallBoardLines(b)
	}

	lines := make([]string, 0, 5)
	for row := range 3 {
		cells := make([]string, 3)
		for col := range 3 {
			cell := renderTeam(b.Cells[row][col])
			if selectingCell && cursorVisible && boardRow == selectedBoardRow && boardCol == selectedBoardCol && row == selectedCellRow && col == selectedCellCol {
				if b.Cells[row][col] == game.Empty {
					cell = teamStyle(currentTurn).Render("█")
				} else {
					cell = cursorStyle(currentTurn).Render(cell)
				}
			}
			cells[col] = cell
		}
		lines = append(lines, fmt.Sprintf(" %s | %s | %s ", cells[0], cells[1], cells[2]))
		if row < 2 {
			lines = append(lines, "-----------")
		}
	}
	return lines
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
		winnerStyle := teamStyle(b.Winner)
		return []string{
			"           ",
			winnerStyle.Render(fmt.Sprintf("     %s     ", b.Winner)),
			"           ",
			"-----------",
			winnerStyle.Render("    WON    "),
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

func cursorStyle(team game.Team) lipgloss.Style {
	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("230"))
	if team == game.X {
		return style.Background(lipgloss.Color("196"))
	}
	return style.Background(lipgloss.Color("27"))
}

func teamStyle(team game.Team) lipgloss.Style {
	if team == game.X {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("27")).Bold(true)
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
