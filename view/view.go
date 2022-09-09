package view

import (
	wrapper "lutris-tui/wrapper"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func Start(wrapper wrapper.Wrapper, games []wrapper.Game) error {
	model := initialModel(wrapper, games)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		return err
	}

	return nil
}

type model struct {
	lutris    wrapper.Wrapper
	games     []wrapper.Game
	cursor    int
	paginator paginator.Model
	start     int
	end       int
	statusBar string
}

func initialModel(wrapper wrapper.Wrapper, games []wrapper.Game) model {
	p := paginator.New()
	p.Type = paginator.Arabic
	p.PerPage = GAMES_PER_PAGE
	p.SetTotalPages(len(games))

	return model{
		lutris:    wrapper,
		games:     games,
		paginator: p,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

const GAMES_PER_PAGE = 12

var styleGame = lipgloss.NewStyle().
	Padding(0, 2, 1).
	Margin(1, 0).
	Width(24).
	MaxHeight(3).
	Align(lipgloss.Center).
	Background(lipgloss.Color("#217A9A")).
	Foreground(lipgloss.Color("#FFFFFF"))

var styleGameRunning = styleGame.Copy().
	Background(lipgloss.Color("#603B6F"))

func (m model) View() string {
	s := ""

	for i, game := range m.games[m.start:m.end] {
		var game_cell string

		if game.IsRunning {
			game_cell = styleGameRunning.Render(game.Name)
		} else {
			game_cell = styleGame.Render(game.Name)
		}

		var cursor string

		if i == m.cursor {
			cursor = "▋\n▋"
		} else {
			cursor = " "
		}

		game_field := lipgloss.JoinHorizontal(lipgloss.Center, cursor, " ", game_cell)

		s += game_field
		s += "\n"
	}

	if len(m.statusBar) > 0 {
		s += m.statusBar
	}
	s += "\n"

	return s
}

type statusMsg int
type errMsg struct{ err error }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cursorRealIdx := m.start + m.cursor

	switch msg := msg.(type) {

	case errMsg:
		m.statusBar = msg.err.Error()

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "down", "j":
			if m.cursor < m.end-m.start-1 {
				m.cursor++
			} else if !m.paginator.OnLastPage() {
				m.cursor = 0
				m.paginator.NextPage()
			}

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			} else if m.paginator.Page != 0 {
				m.cursor = GAMES_PER_PAGE - 1
				m.paginator.PrevPage()
			}

		case "enter":
			// Run the game
			if m.games[cursorRealIdx].IsRunning {
				m.games[cursorRealIdx].Stop()
			} else {
				m.games[cursorRealIdx].IsRunning = true
				return m, runGame(&m.games[cursorRealIdx])
			}
		}
	}

	m.start, m.end = m.paginator.GetSliceBounds(len(m.games))

	return m, nil
}

func runGame(game *wrapper.Game) tea.Cmd {
	return func() tea.Msg {
		command, err := game.Start()

		if err != nil {
			return errMsg{err}
		}

		command.Process.Wait()

		game.IsRunning = false

		return statusMsg(0)
	}
}
