package components

import (
	"lutris-tui/lutris"
	S "lutris-tui/view/styles"

	"github.com/charmbracelet/lipgloss"
)

func Main(gamesGrid [][]lutris.Game, cursorX int, cursorY int) string {
	gamesGridView := GamesGrid(gamesGrid, cursorX, cursorY)
	gameStatsView := GameStats(gamesGrid[cursorY][cursorX])
	return lipgloss.JoinHorizontal(lipgloss.Top, gamesGridView, "  ", gameStatsView) + "\n"
}

type gameState int

const (
	_GS_NORMAL gameState = iota
	_GS_SELECTED
	_GS_RUNNING
)

func Game(name string, state gameState) string {
	switch state {
	case _GS_SELECTED:
		return S.StyleGameSelected.Render(name)
	case _GS_RUNNING:
		return S.StyleGameRunning.Render(name)
	case _GS_NORMAL:
	default:
	}
	return S.StyleGame.Render(name)
}

func GamesGrid(grid [][]lutris.Game, cursorX int, cursorY int) string {
	var gridView string

	for i, row := range grid {
		var columnView string

		for j, game := range row {
			var gameView string

			var gameState gameState

			if game.IsRunning {
				gameState = _GS_RUNNING
			} else if cursorX == j && cursorY == i {
				gameState = _GS_SELECTED
			} else {
				gameState = _GS_NORMAL
			}

			gameView = Game(game.Name, gameState)

			columnView = lipgloss.JoinHorizontal(lipgloss.Center, columnView, " ", gameView)
		}

		gridView += columnView + "\n"
	}

	gridView = S.StyleGamesGrid.Render(gridView)

	return gridView
}

func GameStats(game lutris.Game) string {
	s := ""
	s += S.StyleColoredText.Render("Game Stats") + "\n"
	s += KeyValueLine("name", game.Name)
	s += KeyValueLine("platform", game.Platform)
	s += KeyValueLine("runner", game.Runner)

	if len(game.LastPlayed) != 0 {
		s += KeyValueLine("last played", game.LastPlayed)
	}
	if len(game.PlayTime) != 0 {
		s += KeyValueLine("playtime", game.PlayTime)
	}

	if game.IsRunning {
		s += S.StyleDarkerText.Render("Running") + "\n"
	}

	s = S.StyleGameStats.Render(s)

	return s
}

func KeyValueLine(key string, value string) string {
	return S.StyleNormal.Render(key+":") + " " + S.StyleDarkerText.Render(value) + "\n"
}
