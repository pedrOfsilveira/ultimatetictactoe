package cli

import (
	"fmt"
	"strings"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pedrofsilveira/ultimatetictactoe/internal/game"
)

type menuAction int

const (
	menuNone menuAction = iota
	menuLocal
	menuHost
	menuJoin
	menuQuit
)

var menuItems = []string{
	"Local multiplayer",
	"Host LAN game",
	"Join LAN game",
	"Rules",
	"Quit",
}

type menuModel struct {
	cursor  int
	joining bool
	address string
	action  menuAction
}

func (m menuModel) Init() tea.Cmd { return nil }

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	if m.joining {
		switch key.Type {
		case tea.KeyEsc:
			m.joining = false
		case tea.KeyBackspace, tea.KeyDelete:
			if len(m.address) > 0 {
				_, size := utf8.DecodeLastRuneInString(m.address)
				m.address = m.address[:len(m.address)-size]
			}
		case tea.KeyEnter:
			if strings.TrimSpace(m.address) != "" {
				m.address = strings.TrimSpace(m.address)
				m.action = menuJoin
				return m, tea.Quit
			}
		case tea.KeyRunes:
			m.address += string(key.Runes)
		}
		return m, nil
	}

	switch key.String() {
	case "ctrl+c", "q":
		m.action = menuQuit
		return m, tea.Quit
	case "up", "k":
		m.cursor = (m.cursor + len(menuItems) - 1) % len(menuItems)
	case "down", "j":
		m.cursor = (m.cursor + 1) % len(menuItems)
	case "enter", " ":
		switch m.cursor {
		case 0:
			m.action = menuLocal
			return m, tea.Quit
		case 1:
			m.action = menuHost
			return m, tea.Quit
		case 2:
			m.joining = true
		case 3:
			m.cursor = 0
			return rulesModel{}, nil
		case 4:
			m.action = menuQuit
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m menuModel) View() string {
	view := titleStyle.Render("ULTIMATE TIC-TAC-TOE") + "\n\n"
	if m.joining {
		view += "Join a LAN game\n\n"
		view += "Host address: " + m.address + "█\n\n"
		return view + helpStyle.Render("type IP:PORT • enter: connect • esc: back") + "\n"
	}

	for i, item := range menuItems {
		marker := "  "
		if i == m.cursor {
			marker = teamStyle(gameTeamForMenu(i)).Render("▶ ")
		}
		view += fmt.Sprintf("%s%s\n", marker, item)
	}
	return view + "\n" + helpStyle.Render("arrows: move • enter: select • q: quit") + "\n"
}

func gameTeamForMenu(index int) game.Team {
	if index%2 == 0 {
		return game.X
	}
	return game.O
}

type rulesModel struct{}

func (rulesModel) Init() tea.Cmd { return nil }

func (rulesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "esc", "enter", " ", "q":
			return menuModel{}, nil
		}
	}
	return rulesModel{}, nil
}

func (rulesModel) View() string {
	return titleStyle.Render("ULTIMATE TIC-TAC-TOE — RULES") + "\n\n" +
		"Win small boards by getting three in a row.\n" +
		"Win the game by winning three small boards in a row.\n" +
		"Your move sends the opponent to the matching small board.\n" +
		"If that board is finished, the opponent may choose any board.\n\n" +
		helpStyle.Render("enter/esc: back") + "\n"
}

func runMenu() (menuModel, error) {
	result, err := tea.NewProgram(menuModel{}, tea.WithAltScreen()).Run()
	if err != nil {
		return menuModel{}, err
	}
	menu, ok := result.(menuModel)
	if !ok {
		return menuModel{}, nil
	}
	return menu, nil
}
