package view

import (
	component "lutris-tui/view/components"
	S "lutris-tui/view/styles"
	lutris "lutris-tui/wrapper"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
)

func Start(lutrisClient lutris.LutrisClient, games []lutris.Game) error {
	model := initialModel(lutrisClient, games)
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

type gamesGrid struct {
	cells     [][]lutris.Game
	cursor    CursorPosition
	paginator paginator.Model
	start     int
	end       int
	rowCount  int
}

type model struct {
	lutrisClient lutris.LutrisClient
	games        []lutris.Game
	grid         gamesGrid
	statusBar    string
	selectedGame *lutris.Game
}

func initialModel(lutrisClient lutris.LutrisClient, games []lutris.Game) model {
	p := paginator.New()
	p.Type = paginator.Arabic
	p.PerPage = _GAMES_PER_PAGE
	p.SetTotalPages(len(games))

	model := model{
		lutrisClient: lutrisClient,
		games:        games,
		grid: gamesGrid{
			paginator: p,
			rowCount:  3,
			start:     0,
			end:       _GAMES_PER_PAGE,
		},
	}

	model.updateGameGrid()

	return model
}

func (m model) Init() tea.Cmd {
	return nil
}

const _GAMES_PER_PAGE = 18

func (m model) View() string {
	s := ""

	s += component.Main(m.grid.cells, m.grid.cursor.x, m.grid.cursor.y)

	s += S.StyleDarkerText.Render("  ──────────────────────────────") + "\n"

	s += "               " + S.StyleDarkerText.Render(m.grid.paginator.View()) + "\n"
	s += "  " + S.StyleNormal.Render("↑/k - up, ↓/j - down, q - quit") + "\n"

	if len(m.statusBar) > 0 {
		s += S.StyleNormal.Render(m.statusBar)
	}
	s += "\n"

	s += S.StyleDarkerText.Render("  LUTRIS TUI WRAPPER (alpha)") + "\n"

	return s
}

func (m *model) updateGameGrid() {
	var gameLayout = [][]lutris.Game{}

	for i := m.grid.start; i < m.grid.end; {
		var rowGames []lutris.Game

		for j := 0; j < m.grid.rowCount; j++ {
			if i < m.grid.end {
				rowGames = append(rowGames, m.games[i])
				i++
			} else {
				break
			}
		}

		gameLayout = append(gameLayout, rowGames)
	}

	m.grid.cells = gameLayout
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
			if m.grid.cursor.y < len(m.grid.cells)-1 {
				gamesBelowCount := len(m.grid.cells[m.grid.cursor.y+1])
				if m.grid.cursor.x < gamesBelowCount {
					m.grid.cursor.y++
					break
				}
			}
			if !m.grid.paginator.OnLastPage() {
				m.grid.cursor.y = 0
				m.grid.paginator.NextPage()
			}

		case "up", "k":
			if m.grid.cursor.y > 0 {
				m.grid.cursor.y--
			} else if m.grid.paginator.Page != 0 {
				m.grid.cursor.y = 0
				m.grid.paginator.PrevPage()
			}

		case "right", "l":
			if m.grid.cursor.x < len(m.grid.cells[m.grid.cursor.y])-1 {
				m.grid.cursor.x++
			}

		case "left", "h":
			if m.grid.cursor.x > 0 {
				m.grid.cursor.x--
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

	m.grid.start, m.grid.end = m.grid.paginator.GetSliceBounds(len(m.games))
	m.updateGameGrid()
	m.selectedGame = &m.grid.cells[m.grid.cursor.y][m.grid.cursor.x]

	return m, nil
}

func runGame(game *lutris.Game) tea.Cmd {
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
