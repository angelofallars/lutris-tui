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

type cursorPosition struct {
	x int
	y int
}

type model struct {
	lutris    wrapper.Wrapper
	games     []wrapper.Game
	cursor    cursorPosition
	paginator paginator.Model
	start     int
	end       int
	statusBar string
	gamesView [][]wrapper.Game
}

func initialModel(wrapper wrapper.Wrapper, games []wrapper.Game) model {
	p := paginator.New()
	p.Type = paginator.Arabic
	p.PerPage = _GAMES_PER_PAGE
	p.SetTotalPages(len(games))

	return model{
		lutris:    wrapper,
		games:     games,
		paginator: p,
		gamesView: paginateTwoColumnGames(games, 0, _GAMES_PER_PAGE),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

const _GAMES_PER_PAGE = 12

func (m model) View() string {
	s := ""

	gamesView := ""

	var selected_game wrapper.Game

	for i, row := range m.gamesView {
		var columnView string

		for j, game := range row {
			var gameCell string

			if game.IsRunning {
				gameCell = styleGameRunning.Render(game.Name)
			} else if j == m.cursor.x && i == m.cursor.y {
				gameCell = styleGameSelected.Render(game.Name)
				selected_game = game
			} else {
				gameCell = styleGame.Render(game.Name)
			}

			columnView = lipgloss.JoinHorizontal(lipgloss.Center, columnView, " ", gameCell)
		}

		gamesView += columnView + "\n"
	}

	gamesView = styleGamesView.Render(gamesView)
	gameStats := showGameStats(selected_game)

	mainView := lipgloss.JoinHorizontal(lipgloss.Top, gamesView, "  ", gameStats)

	s += mainView + "\n"

	s += styleDarkerText.Render("  ──────────────────────────────") + "\n"

	s += "               " + styleDarkerText.Render(m.paginator.View()) + "\n"
	s += "  " + styleNormal.Render("↑/k - up, ↓/j - down, q - quit") + "\n"

	if len(m.statusBar) > 0 {
		s += styleNormal.Render(m.statusBar)
	}
	s += "\n"

	s += styleDarkerText.Render("  LUTRIS TUI WRAPPER (alpha)") + "\n"

	return s
}

func paginateTwoColumnGames(games []wrapper.Game, start int, end int) [][]wrapper.Game {
	var gameLayout = [][]wrapper.Game{}

	for i := start; i < end; i++ {
		if i+1 < end {
			gameLayout = append(gameLayout, []wrapper.Game{games[i], games[i+1]})
			i++
		} else {
			gameLayout = append(gameLayout, []wrapper.Game{games[i]})
		}
	}

	return gameLayout
}

func showGameStats(game wrapper.Game) string {
	s := ""
	s += styleColoredText.Render("Game Stats") + "\n"
	s += makeKeyValueLine("name", game.Name)
	s += makeKeyValueLine("platform", game.Platform)
	s += makeKeyValueLine("runner", game.Runner)

	if len(game.LastPlayed) != 0 {
		s += makeKeyValueLine("last played", game.LastPlayed)
	}
	if len(game.PlayTime) != 0 {
		s += makeKeyValueLine("playtime", game.PlayTime)
	}

	if game.IsRunning {
		s += styleDarkerText.Render("Running") + "\n"
	}

	s = styleGameStats.Render(s)

	return s
}

func makeKeyValueLine(key string, value string) string {
	return styleNormal.Render(key+":") + " " + styleDarkerText.Render(value) + "\n"
}

type statusMsg int
type errMsg struct{ err error }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cursorRealIdx := m.start

	switch msg := msg.(type) {

	case errMsg:
		m.statusBar = msg.err.Error()

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "down", "j":
			if m.cursor.y < len(m.gamesView)-1 {
				gamesBelowCount := len(m.gamesView[m.cursor.y+1])
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
			if m.cursor.x < len(m.gamesView[m.cursor.y])-1 {
				m.cursor.x++
			}

		case "left", "h":
			if m.cursor.x > 0 {
				m.cursor.x--
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
	m.gamesView = paginateTwoColumnGames(m.games, m.start, m.end)

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
