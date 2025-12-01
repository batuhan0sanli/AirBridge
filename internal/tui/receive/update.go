package receive

import (
	"AirBridge/internal/tui"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) Init() tea.Cmd {
	m.nextStep()
	return tea.Batch(
		generateKeyCmd(),
		m.spinner.Tick,
		textarea.Blink,
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case keyGeneratedMsg:
		m.privateKey = msg.privateKey
		m.publicKey = msg.publicKey
		m.encodedKey = msg.encodedKey
		m.nextStep()
		return m, nil

	case fileDecryptedMsg:
		m.statusText = "File saved successfully!"
		m.nextStep()
		return m, nil

	case errMsg:
		m.err = msg.error
		m.statusText = ""
		return m, nil

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		headerH := lipgloss.Height(tui.Header())
		footerH := lipgloss.Height(tui.Footer(m.err))
		headerW := lipgloss.Width(tui.Header())
		spacer := 0

		availableHeight := m.Window.Height - (headerH + footerH + spacer)
		if availableHeight < 3 {
			availableHeight = 3
		}
		m.AvailableHeight = availableHeight

		availableWidth := headerW
		m.AvailableWidth = availableWidth
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

		switch m.step {
		case StepGeneratingKey:
			// Wait for key generation
			return m, nil
		case StepAwaitingPayload:
			// Handle Copy Key
			if msg.String() == "ctrl+c" {
				// Already handled by global key? No, global is Ctrl+C for quit.
				// Let's use a specific key for copy, e.g., 'c' or just Enter as requested.
				// User request: "add a button to copy to clipboard".
				// In TUI, we can use a key press. Let's use 'enter' to copy key if focus is not on textarea?
				// Or maybe a specific key binding.
				// But wait, we also need to paste the payload.
				// Let's make it simple:
				// Top part: Public Key. "Press 'c' to copy public key".
				// Bottom part: Textarea.
			}

			// If user presses 'c' (and not typing in textarea?), copy key.
			// But textarea captures input.
			// We can check if textarea is focused.

			// Actually, the user flow is:
			// 1. Receive starts -> Generates Key -> Shows Key.
			// 2. User copies key (Press Enter/c).
			// 3. User sends key to sender.
			// 4. Sender sends payload.
			// 5. User pastes payload into textarea.
			// 6. User presses Enter to decrypt.

			// Let's use 'ctrl+y' to copy key to avoid conflict with typing?
			// Or just a button if we had mouse support.
			// Let's stick to: "Press 'ctrl+k' to copy public key".
			if msg.Type == tea.KeyCtrlK {
				err := clipboard.WriteAll(m.encodedKey)
				if err != nil {
					m.err = err
				} else {
					m.statusText = tui.SuccessStyle.Render("Public key copied to clipboard!")
				}
				return m, nil
			}

			// Handle Textarea input
			m.textarea, cmd = m.textarea.Update(msg)

			// Check for submit (Ctrl+Enter or Enter if single line? Payload is large, so likely multiline textarea)
			// Textarea default is Enter for newline.
			// Let's use Ctrl+Enter to submit.
			if msg.Type == tea.KeyCtrlD { // Ctrl+D for EOF/Submit often used
				// Or just check if payload is pasted?
				// Let's use a specific key to "Process Payload".
			}

			// Let's use the same logic as send module: Enter to submit if it was single line, but here it's base64 blob.
			// Let's use Ctrl+S to save/start decryption?
			if msg.Type == tea.KeyCtrlS {
				m.payload = m.textarea.Value()
				m.statusText = "Decrypting..."
				m.nextStep()
				return m, tea.Batch(
					decryptAndSaveCmd(m.payload, m.privateKey),
					m.spinner.Tick,
				)
			}

			return m, cmd
		}
	}

	switch m.step {
	case StepGeneratingKey, StepDecrypting:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, cmd
}
