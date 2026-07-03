package game

import (
	"errors"
	"fmt"
)

type Game struct {
	Player1     Player
	Player2     Player
	CurrentTurn Team
	Board       Board
	Winner      Team
	Status      Status

	NextBoardRow int
	NextBoardCol int
	FreeMove     bool
}

type Player struct {
	Name string
	Team Team
}

type Team string

const (
	X     Team = "X"
	O     Team = "O"
	Empty Team = " "
)

type Status string

const (
	Playing  Status = "playing"
	Finished Status = "finished"
)

type SmallBoard struct {
	Cells  [3][3]Team
	Winner Team
	Status Status
}

type Board struct {
	Boards [3][3]SmallBoard
}

func NewGame(player1Name, player2Name string) *Game {
	return &Game{
		Player1:     Player{Name: player1Name},
		Player2:     Player{Name: player2Name},
		CurrentTurn: X,
		Winner:      Empty,
		Status:      Playing,
		Board:       NewBoard(),
		FreeMove:    true,
	}
}

func NewBoard() Board {
	board := Board{}

	for row := range 3 {
		for col := range 3 {
			board.Boards[row][col] = NewSmallBoard()
		}
	}

	return board
}

func NewSmallBoard() SmallBoard {
	return SmallBoard{
		Cells: [3][3]Team{
			{Empty, Empty, Empty},
			{Empty, Empty, Empty},
			{Empty, Empty, Empty},
		},
		Winner: Empty,
		Status: Playing,
	}
}

func (g *Game) PickX(player int) {
	if player == 1 {
		g.Player1.Team = X
		g.Player2.Team = O
	} else {
		g.Player1.Team = O
		g.Player2.Team = X
	}

	g.CurrentTurn = X
}

func (g *Game) MakeMove(boardRow, boardCol, cellRow, cellCol int) error {
	if g.Status != Playing {
		return errors.New("game has already finished")
	}

	if boardRow < 0 || boardRow > 2 || boardCol < 0 || boardCol > 2 {
		return errors.New("invalid board")
	}

	if cellRow < 0 || cellRow > 2 || cellCol < 0 || cellCol > 2 {
		return errors.New("invalid cell")
	}

	if !g.FreeMove && (boardRow != g.NextBoardRow || boardCol != g.NextBoardCol) {
		return errors.New("you must play in the required board")
	}

	smallBoard := &g.Board.Boards[boardRow][boardCol]

	if smallBoard.Status != Playing {
		return errors.New("this small board is already finished")
	}

	if smallBoard.Cells[cellRow][cellCol] != Empty {
		return errors.New("the cell is not empty")
	}

	smallBoard.Cells[cellRow][cellCol] = g.CurrentTurn

	if smallWinner := smallBoard.checkWinner(); smallWinner != Empty {
		smallBoard.Winner = smallWinner
		smallBoard.Status = Finished
	} else if smallBoard.isDraw() {
		smallBoard.Status = Finished
	}

	if winner := g.Board.checkWinner(); winner != Empty {
		g.Winner = winner
		g.Status = Finished
		return nil
	}

	if g.Board.isDraw() {
		g.Status = Finished
		return nil
	}

	nextBoard := &g.Board.Boards[cellRow][cellCol]

	if nextBoard.Status != Playing {
		g.FreeMove = true
	} else {
		g.FreeMove = false
		g.NextBoardRow = cellRow
		g.NextBoardCol = cellCol
	}

	g.switchTurn()

	return nil
}

func (b Board) String() string {
	var result string

	for boardRow := range 3 {
		left := b.Boards[boardRow][0].Lines()
		middle := b.Boards[boardRow][1].Lines()
		right := b.Boards[boardRow][2].Lines()

		for line := range len(left) {
			result += left[line]
			result += " || "
			result += middle[line]
			result += " || "
			result += right[line]
			result += "\n"
		}

		if boardRow < 2 {
			result += "=========================================\n"
		}
	}

	return result
}

func (b SmallBoard) Lines() []string {
	if b.Status == Finished && b.Winner != Empty {
		return []string{
			"           ",
			fmt.Sprintf("     %s     ", b.Winner),
			"           ",
			"-----------",
			"    WON    ",
		}
	}

	return []string{
		fmt.Sprintf(" %s | %s | %s ", b.Cells[0][0], b.Cells[0][1], b.Cells[0][2]),
		"-----------",
		fmt.Sprintf(" %s | %s | %s ", b.Cells[1][0], b.Cells[1][1], b.Cells[1][2]),
		"-----------",
		fmt.Sprintf(" %s | %s | %s ", b.Cells[2][0], b.Cells[2][1], b.Cells[2][2]),
	}
}

func (b *SmallBoard) checkWinner() Team {
	for i := range 3 {
		if b.Cells[i][0] == b.Cells[i][1] &&
			b.Cells[i][1] == b.Cells[i][2] &&
			b.Cells[i][0] != Empty {
			return b.Cells[i][0]
		}

		if b.Cells[0][i] == b.Cells[1][i] &&
			b.Cells[1][i] == b.Cells[2][i] &&
			b.Cells[0][i] != Empty {
			return b.Cells[0][i]
		}
	}

	if b.Cells[0][0] == b.Cells[1][1] &&
		b.Cells[1][1] == b.Cells[2][2] &&
		b.Cells[0][0] != Empty {
		return b.Cells[0][0]
	}

	if b.Cells[0][2] == b.Cells[1][1] &&
		b.Cells[1][1] == b.Cells[2][0] &&
		b.Cells[0][2] != Empty {
		return b.Cells[0][2]
	}

	return Empty
}

func (b *SmallBoard) isDraw() bool {
	for _, row := range b.Cells {
		for _, cell := range row {
			if cell == Empty {
				return false
			}
		}
	}

	return true
}

func (b *Board) checkWinner() Team {
	for i := range 3 {
		if b.Boards[i][0].Winner == b.Boards[i][1].Winner &&
			b.Boards[i][1].Winner == b.Boards[i][2].Winner &&
			b.Boards[i][0].Winner != Empty {
			return b.Boards[i][0].Winner
		}

		if b.Boards[0][i].Winner == b.Boards[1][i].Winner &&
			b.Boards[1][i].Winner == b.Boards[2][i].Winner &&
			b.Boards[0][i].Winner != Empty {
			return b.Boards[0][i].Winner
		}
	}

	if b.Boards[0][0].Winner == b.Boards[1][1].Winner &&
		b.Boards[1][1].Winner == b.Boards[2][2].Winner &&
		b.Boards[0][0].Winner != Empty {
		return b.Boards[0][0].Winner
	}

	if b.Boards[0][2].Winner == b.Boards[1][1].Winner &&
		b.Boards[1][1].Winner == b.Boards[2][0].Winner &&
		b.Boards[0][2].Winner != Empty {
		return b.Boards[0][2].Winner
	}

	return Empty
}

func (b *Board) isDraw() bool {
	for _, row := range b.Boards {
		for _, smallBoard := range row {
			if smallBoard.Status == Playing {
				return false
			}
		}
	}

	return true
}

func (g *Game) switchTurn() {
	if g.CurrentTurn == X {
		g.CurrentTurn = O
	} else {
		g.CurrentTurn = X
	}
}
