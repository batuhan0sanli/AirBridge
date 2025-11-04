package tui

import tea "github.com/charmbracelet/bubbletea"

type BaseModel interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() string
}
