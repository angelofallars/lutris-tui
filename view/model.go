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
