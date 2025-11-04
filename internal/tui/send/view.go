package send

import (
	"AirBridge/internal/tui"
)

func (m *Model) View() string {
	switch m.step {
	case StepUndefined:
		return tui.View(m.Window, m.err, "")
	case StepAwaitingFile:
		m.filePath.SetHeight(m.AvailableHeight - 1)
		view := m.filePath.View()
		return tui.View(m.Window, m.err, view)
	case StepReadyingFile:
		input := "Readying file"
		view := tui.MainStyle(m.Window).Render(input)
		return tui.View(m.Window, m.err, view)
	default:
		panic("unhandled default case")
	}
	return ""
}
