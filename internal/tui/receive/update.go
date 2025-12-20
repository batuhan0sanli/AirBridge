package receive

import (
	"AirBridge/internal/tui"
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) Init() tea.Cmd {
	m.nextStep()
	if m.privateKey != nil {
		// If we have both key and payload at init, start decryption immediately
		if m.payload != "" {
			m.statusText = "Decrypting..."
			m.step = StepDecrypting // Force step update for consistency
			return tea.Batch(
				decryptAndSaveCmd(m.payload, m.privateKey),
				m.spinner.Tick,
			)
		}

		return tea.Batch(
			m.spinner.Tick,
			textarea.Blink,
		)
	}
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
		// Handle file deletion if requested
		if m.deleteFile && m.payloadPath != "" {
			err := os.Remove(m.payloadPath)
			if err != nil {
				m.statusText += fmt.Sprintf(" (Failed to delete payload: %v)", err)
			} else {
				m.statusText += " (Payload deleted)"
			}
		}
		m.nextStep()
		return m, nil

	case errMsg:
		m.err = msg.error
		m.statusText = ""
		switch m.step {
		case StepDecrypting:
			m.payload = ""
			m.textarea.Reset()
		case StepGeneratingKey:
			m.privateKey = nil
			m.publicKey = nil
			m.encodedKey = ""
		}
		m.nextStep()
		return m, nil

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		headerH := lipgloss.Height(tui.Header())
		footerH := lipgloss.Height(tui.Footer(m.err))
		headerW := lipgloss.Width(tui.Header())
		spacer := 0

		availableHeight := m.Height - (headerH + footerH + spacer)
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

			if msg.Type == tea.KeyEnter {
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
