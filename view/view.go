package view

import (
	component "lutris-tui/view/components"
	S "lutris-tui/view/styles"
	lutris "lutris-tui/wrapper"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
)

func Start(wrapper lutris.Wrapper, games []lutris.Game) error {
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
	lutris       lutris.Wrapper
	games        []lutris.Game
	cursor       CursorPosition
	paginator    paginator.Model
	pageStartIdx int
	pageEndIdx   int
	statusBar    string
	gamesGrid    [][]lutris.Game
	selectedGame *lutris.Game
	rowCount     int
}

func initialModel(wrapper lutris.Wrapper, games []lutris.Game) model {
	p := paginator.New()
	p.Type = paginator.Arabic
	p.PerPage = _GAMES_PER_PAGE
	p.SetTotalPages(len(games))

	model := model{
		lutris:    wrapper,
		games:     games,
		paginator: p,
		rowCount:  3,
	}

	model.updateGameGrid(games, 0, _GAMES_PER_PAGE)

	return model
}

func (m model) Init() tea.Cmd {
	return nil
}

const _GAMES_PER_PAGE = 18

func (m model) View() string {
	s := ""

	s += component.Main(m.gamesGrid, m.cursor.x, m.cursor.y)

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

func (m *model) updateGameGrid(games []lutris.Game, start int, end int) {
	var gameLayout = [][]lutris.Game{}

	for i := start; i < end; {
		var rowGames []lutris.Game

		for j := 0; j < m.rowCount; j++ {
			if i < end {
				rowGames = append(rowGames, games[i])
				i++
			} else {
				break
			}
		}

		gameLayout = append(gameLayout, rowGames)
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

	m.pageStartIdx, m.pageEndIdx = m.paginator.GetSliceBounds(len(m.games))
	m.updateGameGrid(m.games, m.pageStartIdx, m.pageEndIdx)
	m.selectedGame = &m.gamesGrid[m.cursor.y][m.cursor.x]

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
