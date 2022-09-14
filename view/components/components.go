package components

import (
	S "lutris-tui/view/styles"
	wrapper "lutris-tui/wrapper"

	"github.com/charmbracelet/lipgloss"
)

func Main(gamesGrid [][]wrapper.Game, cursorX int, cursorY int) string {
	gamesGridView := GamesGrid(gamesGrid, cursorX, cursorY)
	gameStatsView := GameStats(gamesGrid[cursorY][cursorX])
	return lipgloss.JoinHorizontal(lipgloss.Top, gamesGridView, "  ", gameStatsView) + "\n"
}

type GameState int

const (
	GS_NORMAL GameState = iota
	GS_SELECTED
	GS_RUNNING
)

func Game(name string, state GameState) string {
	switch state {
	case GS_SELECTED:
		return S.StyleGameSelected.Render(name)
	case GS_RUNNING:
		return S.StyleGameRunning.Render(name)
	case GS_NORMAL:
	default:
	}
	return S.StyleGame.Render(name)
}

func GamesGrid(grid [][]wrapper.Game, cursorX int, cursorY int) string {
	var gridView string

	for i, row := range grid {
		var columnView string

		for j, game := range row {
			var gameView string

			var gameState GameState

			if game.IsRunning {
				gameState = GS_RUNNING
			} else if cursorX == j && cursorY == i {
				gameState = GS_SELECTED
			} else {
				gameState = GS_NORMAL
			}

			gameView = Game(game.Name, gameState)

			columnView = lipgloss.JoinHorizontal(lipgloss.Center, columnView, " ", gameView)
		}

		gridView += columnView + "\n"
	}

	gridView = S.StyleGamesGrid.Render(gridView)

	return gridView
}

func GameStats(game wrapper.Game) string {
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
