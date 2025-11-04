package send

import (
	"AirBridge/internal/tui"
)

func (m *Model) View() string {
	switch m.step {
	case StepUndefined:
		return tui.View(m.err, "")
	case StepAwaitingFile:
		m.filepicker.SetHeight(m.AvailableHeight - 1)
		view := m.filepicker.View()
		return tui.View(m.err, view)
	case StepReadyingFile:
		input := m.spinner.View() + m.statusText
		view := tui.MainStyle(m.Window).Render(input)
		return tui.View(m.err, view)
	case StepAwaitingPublicKey:
		input := "Public Key"
		view := tui.MainStyle(m.Window).Render(input)
		return tui.View(m.err, view)
	default:
		panic("unhandled default case")
	}
	return ""
}
