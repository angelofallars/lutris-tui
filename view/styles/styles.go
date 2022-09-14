package styles

import "github.com/charmbracelet/lipgloss"

var ColorFg1 = lipgloss.Color("#FFFFFF")
var ColorFg2 = lipgloss.Color("#A6A6A6")
var ColorCellBg1 = lipgloss.Color("#257693")
var ColorCellBg2 = lipgloss.Color("#0090C4")
var ColorCellBg3 = lipgloss.Color("#603B6F")

var StyleNormal = lipgloss.NewStyle().
	Foreground(ColorFg1)

var StyleDarkerText = lipgloss.NewStyle().
	Foreground(ColorFg2)

var StyleColoredText = lipgloss.NewStyle().
	Foreground(ColorCellBg1)

var StyleGamesGrid = lipgloss.NewStyle().
	Height(3*6 + 1)

var StyleGame = lipgloss.NewStyle().
	PaddingBottom(2).
	PaddingRight(2).
	Margin(1, 0).
	Width(22).
	Height(4).
	MaxHeight(4).
	Align(lipgloss.Left).
	Background(ColorCellBg1).
	Foreground(ColorFg1)

var StyleGameSelected = StyleGame.Copy().
	Background(ColorCellBg2)

var StyleGameRunning = StyleGame.Copy().
	Background(ColorCellBg3)

var StyleGameStats = lipgloss.NewStyle().
	MarginTop(1)
