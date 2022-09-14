package view

import "github.com/charmbracelet/lipgloss"

var colorFg1 = lipgloss.Color("#FFFFFF")
var colorFg2 = lipgloss.Color("#A6A6A6")
var colorCellBg1 = lipgloss.Color("#257693")
var colorCellBg2 = lipgloss.Color("#0090C4")
var colorCellBg3 = lipgloss.Color("#603B6F")

var styleNormal = lipgloss.NewStyle().
	Foreground(colorFg1)

var styleDarkerText = lipgloss.NewStyle().
	Foreground(colorFg2)

var styleColoredText = lipgloss.NewStyle().
	Foreground(colorCellBg1)

var styleGamesView = lipgloss.NewStyle().
	Height(3*_GAMES_PER_PAGE + 1)

var styleGame = lipgloss.NewStyle().
	PaddingBottom(2).
	PaddingRight(2).
	Margin(1, 0).
	Width(30).
	MaxHeight(3).
	Align(lipgloss.Left).
	Background(colorCellBg1).
	Foreground(colorFg1)

var styleGameSelected = styleGame.Copy().
	Background(colorCellBg2)

var styleGameRunning = styleGame.Copy().
	Background(colorCellBg3)

var styleGameStats = lipgloss.NewStyle().
	MarginTop(1)
