package send

import (
	"AirBridge/internal/tui"
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func errorView(err error) string {
	if err != nil {
		return fmt.Sprintf("⚠️  %s", err)
	}
	return ""
}

func (m *Model) View() string {
	banner := tui.AirBridgeBanner()
	switch m.step {
	case StepUndefined:
		return fmt.Sprintf(
			"\n%s\nError: Sender\n\n%s\n[q] Çıkış\n(esc to quit)",
			tui.AirBridgeBanner(),
			errorView(m.err),
		)

	case StepAwaitingFile:
		input := m.filePath.View()
		view := lipgloss.JoinVertical(
			lipgloss.Left,
			banner,
			"",
			input,
			errorView(m.err),
			"(esc to quit)",
		)
		return view
	case StepReadyingFile:
		input := "Readying file"
		view := lipgloss.JoinVertical(
			lipgloss.Left,
			banner,
			"",
			input,
			errorView(m.err),
			"(esc to quit)",
		)
		return view
	default:
		panic("unhandled default case")
	}
	return ""
}
