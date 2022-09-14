package view

import (
	component "lutris-tui/view/components"
	S "lutris-tui/view/styles"
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

type CursorPosition struct {
	x int
	y int
}

type model struct {
	lutris       wrapper.Wrapper
	games        []wrapper.Game
	cursor       CursorPosition
	paginator    paginator.Model
	start        int
	end          int
	statusBar    string
	gamesGrid    [][]wrapper.Game
	selectedGame *wrapper.Game
}

func initialModel(wrapper wrapper.Wrapper, games []wrapper.Game) model {
	p := paginator.New()
	p.Type = paginator.Arabic
	p.PerPage = _GAMES_PER_PAGE
	p.SetTotalPages(len(games))

	model := model{
		lutris:    wrapper,
		games:     games,
		paginator: p,
	}

	model.updateGameGrid(games, 0, _GAMES_PER_PAGE)

	return model
}

func (m model) Init() tea.Cmd {
	return nil
}

const _GAMES_PER_PAGE = 12

func (m model) View() string {
	s := ""

	gamesView := component.GamesGrid(m.gamesGrid, m.cursor.x, m.cursor.y)

	var gameStats string
	if m.selectedGame != nil {
		gameStats = component.GameStats(*m.selectedGame)
	}

	mainView := lipgloss.JoinHorizontal(lipgloss.Top, gamesView, "  ", gameStats)

	s += mainView + "\n"

	s += S.StyleDarkerText.Render("  ──────────────────────────────") + "\n"

	s += "               " + S.StyleDarkerText.Render(m.paginator.View()) + "\n"
	s += "  " + S.StyleNormal.Render("↑/k - up, ↓/j - down, q - quit") + "\n"

	if len(m.statusBar) > 0 {
		s += S.StyleNormal.Render(m.statusBar)
	}
	s += "\n"

	s += S.StyleDarkerText.Render("  LUTRIS TUI WRAPPER (alpha)") + "\n"

	return s
}

func (m *model) updateGameGrid(games []wrapper.Game, start int, end int) {
	var gameLayout = [][]wrapper.Game{}

	for i := start; i < end; i++ {
		if i+1 < end {
			gameLayout = append(gameLayout, []wrapper.Game{games[i], games[i+1]})
			i++
		} else {
			gameLayout = append(gameLayout, []wrapper.Game{games[i]})
		}
	}

	m.gamesGrid = gameLayout
}

type statusMsg int
type errMsg struct{ err error }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case errMsg:
		m.statusBar = msg.err.Error()

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "down", "j":
			if m.cursor.y < len(m.gamesGrid)-1 {
				gamesBelowCount := len(m.gamesGrid[m.cursor.y+1])
				if m.cursor.x < gamesBelowCount {
					m.cursor.y++
					break
				}
			}
			if !m.paginator.OnLastPage() {
				m.cursor.y = 0
				m.paginator.NextPage()
			}

		case "up", "k":
			if m.cursor.y > 0 {
				m.cursor.y--
			} else if m.paginator.Page != 0 {
				m.cursor.y = 0
				m.paginator.PrevPage()
			}

		case "right", "l":
			if m.cursor.x < len(m.gamesGrid[m.cursor.y])-1 {
				m.cursor.x++
			}

		case "left", "h":
			if m.cursor.x > 0 {
				m.cursor.x--
			}

		case "enter":
			// Run the game

			if m.selectedGame.IsRunning {
				m.selectedGame.Stop()
			} else {
				m.selectedGame.IsRunning = true
				return m, runGame(m.selectedGame)
			}
		}
	}

	m.start, m.end = m.paginator.GetSliceBounds(len(m.games))
	m.updateGameGrid(m.games, m.start, m.end)
	m.selectedGame = &m.gamesGrid[m.cursor.y][m.cursor.x]

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
