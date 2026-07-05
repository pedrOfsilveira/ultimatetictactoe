package cli

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pedrofsilveira/ultimatetictactoe/internal/game"
)

type networkMoveModel struct {
	game               *game.Game
	mode               selectionMode
	boardRow, boardCol int
	cellRow, cellCol   int
	move               NetMove
	submitted          bool
	message            string
	cursorVisible      bool
}

func newNetworkMoveModel(g *game.Game) networkMoveModel {
	m := networkMoveModel{game: g, mode: selectBoard, cursorVisible: true}
	if !g.FreeMove {
		m.mode = selectCell
		m.boardRow, m.boardCol = g.NextBoardRow, g.NextBoardCol
	}
	return m
}

func (m networkMoveModel) Init() tea.Cmd { return blinkCursor() }

func (m networkMoveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(cursorBlinkMsg); ok {
		m.cursorVisible = !m.cursorVisible
		return m, blinkCursor()
	}

	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch key.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc":
		if m.mode == selectCell && m.game.FreeMove {
			m.mode = selectBoard
			m.message = ""
		}
	case "up", "k":
		m.moveCursor(-1, 0)
	case "down", "j":
		m.moveCursor(1, 0)
	case "left", "h":
		m.moveCursor(0, -1)
	case "right", "l":
		m.moveCursor(0, 1)
	case "enter", " ":
		if m.mode == selectBoard {
			if m.game.Board.Boards[m.boardRow][m.boardCol].Status != game.Playing {
				m.message = "That board is already finished."
				return m, nil
			}
			m.mode = selectCell
			m.cellRow, m.cellCol = 0, 0
			m.cursorVisible = true
			return m, nil
		}
		board := m.game.Board.Boards[m.boardRow][m.boardCol]
		if board.Cells[m.cellRow][m.cellCol] != game.Empty {
			m.message = "That cell is already occupied."
			return m, nil
		}
		m.move = NetMove{m.boardRow, m.boardCol, m.cellRow, m.cellCol}
		m.submitted = true
		return m, tea.Quit
	}
	return m, nil
}

func (m *networkMoveModel) moveCursor(rowDelta, colDelta int) {
	m.message = ""
	m.cursorVisible = true
	if m.mode == selectBoard {
		m.boardRow = wrap3(m.boardRow + rowDelta)
		m.boardCol = wrap3(m.boardCol + colDelta)
	} else {
		m.cellRow = wrap3(m.cellRow + rowDelta)
		m.cellCol = wrap3(m.cellCol + colDelta)
	}
}

func (m networkMoveModel) View() string {
	status := fmt.Sprintf("Your turn (%s) — select a cell", m.game.CurrentTurn)
	if m.mode == selectBoard {
		status = fmt.Sprintf("Your turn (%s) — select a small board", m.game.CurrentTurn)
	}
	view := titleStyle.Render("ULTIMATE TIC-TAC-TOE — MULTIPLAYER") + "\n\n"
	view += renderInteractiveBoard(m.game, m.boardRow, m.boardCol, m.cellRow, m.cellCol, m.mode == selectCell, m.cursorVisible) + "\n"
	view += status + "\n"
	if m.message != "" {
		view += errorStyle.Render(m.message) + "\n"
	}
	return view + helpStyle.Render("arrows: move • enter: select • esc: back • q: leave") + "\n"
}

func selectNetworkMove(g *game.Game) (NetMove, error) {
	result, err := tea.NewProgram(newNetworkMoveModel(g), tea.WithAltScreen()).Run()
	if err != nil {
		return NetMove{}, err
	}
	model, ok := result.(networkMoveModel)
	if !ok || !model.submitted {
		return NetMove{}, errors.New("left multiplayer game")
	}
	return model.move, nil
}
