package view

import (
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

type model struct {
	lutrisClient lutris.LutrisClient
	games        []lutris.Game
	grid         gamesGrid
	statusBar    string
	selectedGame *lutris.Game
}

type gamesGrid struct {
	cells     [][]lutris.Game
	cursor    CursorPosition
	paginator paginator.Model
	rowCount  int
}

type CursorPosition struct {
	x int
	y int
}

const _GAMES_PER_PAGE = 18

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
		},
	}

	model.updateGameGrid()

	return model
}

func (m model) Init() tea.Cmd {
	return nil
}

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

	m.updateGameGrid()
	m.selectedGame = &m.grid.cells[m.grid.cursor.y][m.grid.cursor.x]

	return m, nil
}

type statusMsg int
type errMsg struct{ err error }

func (m *model) updateGameGrid() {
	var gameLayout = [][]lutris.Game{}

	startIdx, endIdx := m.grid.paginator.GetSliceBounds(len(m.games))

	for i := startIdx; i < endIdx; {
		var rowGames []lutris.Game

		for j := 0; j < m.grid.rowCount; j++ {
			if i < endIdx {
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
