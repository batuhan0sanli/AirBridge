package receive

import (
	"AirBridge/internal/tui"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	switch m.step {
	case StepUndefined:
		return tui.View(m.err, "")
	case StepGeneratingKey:
		input := m.spinner.View() + " Generating RSA Key Pair..."
		view := tui.MainStyle(m.Window).Render(input)
		return tui.View(m.err, view)
	case StepAwaitingPayload:
		encodedText := m.encodedKey[0:10] + " ... " + m.encodedKey[len(m.encodedKey)-10:]
		// Public Key Section
		keyView := lipgloss.NewStyle().
			Padding(1).
			Border(lipgloss.RoundedBorder()).
			Render(encodedText)

		keyHelp := tui.SubtleStyle.Render("Press 'Ctrl+K' to copy public key")

		// Payload Input Section
		m.textarea.SetWidth(m.AvailableWidth - 4)
		m.textarea.SetHeight(10)
		input := m.textarea.View()

		inputHelp := tui.SubtleStyle.Render("Paste payload above and press 'Ctrl+S' to decrypt and save")

		view := lipgloss.JoinVertical(lipgloss.Left,
			"Your Public Key:",
			keyView,
			keyHelp,
			"",
			"Incoming Payload:",
			input,
			inputHelp,
			"",
			m.statusText,
		)

		view = tui.MainStyle(m.Window).Render(view)
		return tui.View(m.err, view)
	case StepDecrypting:
		input := m.spinner.View() + " Decrypting and Saving..."
		view := tui.MainStyle(m.Window).Render(input)
		return tui.View(m.err, view)
	case StepSuccess:
		text := tui.SuccessStyle.Render("File received and saved successfully!")
		view := tui.MainStyle(m.Window).Render(text)
		return tui.View(m.err, view)
	default:
		return tui.View(m.err, "Unknown Step")
	}
}
