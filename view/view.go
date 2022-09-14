package view

import (
	component "lutris-tui/view/components"
	S "lutris-tui/view/styles"
)

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
