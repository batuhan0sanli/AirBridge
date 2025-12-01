package send

import (
	"AirBridge/internal/strutil"
	"AirBridge/internal/tui"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	switch m.step {
	case StepUndefined:
		return tui.View(m.err, "")
	case StepAwaitingFile:
		text := "Please select a file to send:"
		m.filepicker.SetHeight(m.AvailableHeight - 3) // -3 for the text and spacing
		input := m.filepicker.View()
		view := lipgloss.JoinVertical(lipgloss.Left, text, "", input)
		view = tui.MainStyle(m.Window).Render(view)
		return tui.View(m.err, view)
	case StepReadyingFile:
		input := m.spinner.View() + m.statusText
		view := tui.MainStyle(m.Window).Render(input)
		return tui.View(m.err, view)
	case StepAwaitingPublicKey:
		text := "Please paste the recipient's public key and press 'Enter':"
		m.textarea.SetWidth(m.AvailableWidth - 2) // -2 for the spacing
		input := m.textarea.View()
		view := lipgloss.JoinVertical(lipgloss.Left, text, "", input)
		view = tui.MainStyle(m.Window).Render(view)
		return tui.View(m.err, view)
	case StepReadyingPublicKey:
		input := m.spinner.View() + m.statusText
		view := tui.MainStyle(m.Window).Render(input)
		return tui.View(m.err, view)
	case StepReadyToSend:
		text := m.statusText
		if text == "" {
			text = "Press 'Ctrl+K' to copy payload to clipboard."
		}
		payloadText := strutil.TruncateMiddle(m.filePayload, 15)
		input := text + "\n\nPayload: " + payloadText
		view := tui.MainStyle(m.Window).Render(input)
		return tui.View(m.err, view)
	default:
		panic("unhandled default case")
	}

}
