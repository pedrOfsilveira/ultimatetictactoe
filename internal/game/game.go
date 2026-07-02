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

type Board struct {
	Cells [3][3]Team
}

func (b *Board) String() string {
	return fmt.Sprintf(
		" %s | %s | %s \n-----------\n %s | %s | %s \n-----------\n %s | %s | %s ",
		b.Cells[0][0],
		b.Cells[0][1],
		b.Cells[0][2],
		b.Cells[1][0],
		b.Cells[1][1],
		b.Cells[1][2],
		b.Cells[2][0],
		b.Cells[2][1],
		b.Cells[2][2],
	)
}

func (g *Game) PickX(player int) {
	if player == 1 {
		g.Player1.Team = X
		g.Player2.Team = O
	} else {
		g.Player2.Team = X
		g.Player1.Team = O
	}
	g.CurrentTurn = X
}

func NewGame(player1Name, player2Name string) *Game {
	game := Game{
		Player1: Player{
			Name: player1Name,
		},
		Player2: Player{
			Name: player2Name,
		},

		CurrentTurn: X,

		Winner: Empty,
		Status: Playing,

		Board: NewBoard(),
	}
	return &game
}

func NewBoard() Board {
	return Board{
		Cells: [3][3]Team{
			{Empty, Empty, Empty},
			{Empty, Empty, Empty},
			{Empty, Empty, Empty},
		},
	}
}

func (g *Game) MakeMove(row, col int) error {
	if g.Status != Playing {
		return errors.New("game has already finished")
	}

	if row < 0 || row > 2 {
		return errors.New("invalid cell")
	}
	if col < 0 || col > 2 {
		return errors.New("invalid cell")
	}

	if g.Board.Cells[row][col] != Empty {
		return errors.New("the cell is not empty")
	}

	g.Board.Cells[row][col] = g.CurrentTurn

	winner := g.checkWinner()

	if winner != Empty {
		g.Winner = winner
		g.Status = Finished
		return nil
	}

	if g.isDraw() {
		g.Status = Finished
		return nil
	}

	g.switchTurn()

	return nil
}

func (g *Game) checkWinner() Team {
	for i := range 3 {
		// horizontais
		if g.Board.Cells[i][0] == g.Board.Cells[i][1] &&
			g.Board.Cells[i][1] == g.Board.Cells[i][2] &&
			g.Board.Cells[i][0] != Empty {
			return g.Board.Cells[i][0]
		}

		// verticais
		if g.Board.Cells[0][i] == g.Board.Cells[1][i] &&
			g.Board.Cells[1][i] == g.Board.Cells[2][i] &&
			g.Board.Cells[0][i] != Empty {
			return g.Board.Cells[0][i]
		}
	}

	// diagonal principal
	if g.Board.Cells[0][0] == g.Board.Cells[1][1] &&
		g.Board.Cells[1][1] == g.Board.Cells[2][2] &&
		g.Board.Cells[0][0] != Empty {
		return g.Board.Cells[0][0]
	}

	// diagonal secundária
	if g.Board.Cells[0][2] == g.Board.Cells[1][1] &&
		g.Board.Cells[1][1] == g.Board.Cells[2][0] &&
		g.Board.Cells[0][2] != Empty {
		return g.Board.Cells[0][2]
	}

	return Empty
}

func (g *Game) switchTurn() {
	if g.CurrentTurn == X {
		g.CurrentTurn = O
	} else {
		g.CurrentTurn = X
	}
}

func (g *Game) isDraw() bool {
	for _, row := range g.Board.Cells {
		for _, cell := range row {
			if cell == Empty {
				return false
			}
		}
	}

	return true
}
