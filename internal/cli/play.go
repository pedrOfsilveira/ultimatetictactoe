package cli

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pedrofsilveira/ultimatetictactoe/internal/game"
)

type selectionMode int

const (
	selectBoard selectionMode = iota
	selectCell
)

type gameModel struct {
	game               *game.Game
	mode               selectionMode
	boardRow, boardCol int
	cellRow, cellCol   int
	message            string
	cursorVisible      bool
}

type cursorBlinkMsg time.Time

const cursorBlinkInterval = 500 * time.Millisecond

func blinkCursor() tea.Cmd {
	return tea.Tick(cursorBlinkInterval, func(t time.Time) tea.Msg {
		return cursorBlinkMsg(t)
	})
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

func newGameModel() gameModel {
	return gameModel{game: game.NewGame("Player 1", "Player 2"), mode: selectBoard, cursorVisible: true}
}

func (m gameModel) Init() tea.Cmd { return blinkCursor() }

func (m gameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		return m.choose()
	}
	return m, nil
}

func (m *gameModel) moveCursor(rowDelta, colDelta int) {
	if m.game.Status == game.Finished {
		return
	}
	m.message = ""
	m.cursorVisible = true
	if m.mode == selectBoard {
		m.boardRow = wrap3(m.boardRow + rowDelta)
		m.boardCol = wrap3(m.boardCol + colDelta)
		return
	}
	m.cellRow = wrap3(m.cellRow + rowDelta)
	m.cellCol = wrap3(m.cellCol + colDelta)
}

func (m gameModel) choose() (tea.Model, tea.Cmd) {
	if m.game.Status == game.Finished {
		return newGameModel(), nil
	}
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

	if err := m.game.MakeMove(m.boardRow, m.boardCol, m.cellRow, m.cellCol); err != nil {
		m.message = err.Error()
		return m, nil
	}
	m.message = ""
	if m.game.Status == game.Playing {
		if m.game.FreeMove {
			m.mode = selectBoard
		} else {
			m.mode = selectCell
			m.boardRow, m.boardCol = m.game.NextBoardRow, m.game.NextBoardCol
		}
		m.cellRow, m.cellCol = 0, 0
		m.cursorVisible = true
	}
	return m, nil
}

func (m gameModel) View() string {
	var status string
	if m.game.Status == game.Finished {
		status = gameResult(m.game) + "  Press Enter to play again."
	} else if m.mode == selectBoard {
		status = fmt.Sprintf("%s's turn — select a small board", m.game.CurrentTurn)
	} else {
		status = fmt.Sprintf("%s's turn — select a cell in board (%d, %d)", m.game.CurrentTurn, m.boardRow, m.boardCol)
	}

	view := titleStyle.Render("ULTIMATE TIC-TAC-TOE") + "\n\n"
	view += renderInteractiveBoard(m.game, m.boardRow, m.boardCol, m.cellRow, m.cellCol, m.mode == selectCell, m.cursorVisible) + "\n"
	view += status + "\n"
	if m.message != "" {
		view += errorStyle.Render(m.message) + "\n"
	}
	if m.game.Status == game.Playing {
		view += helpStyle.Render("arrows: move • enter/space: select • esc: back • q: quit")
	} else {
		view += helpStyle.Render("enter: new game • q: quit")
	}
	return view + "\n"
}

func wrap3(value int) int { return (value + 3) % 3 }

func playGame() error {
	_, err := tea.NewProgram(newGameModel(), tea.WithAltScreen()).Run()
	return err
}
