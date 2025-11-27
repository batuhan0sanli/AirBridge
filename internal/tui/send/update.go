package send

import (
	"AirBridge/internal/tui"
	"AirBridge/pkg"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) Init() tea.Cmd {
	m.nextStep()
	return tea.Batch(
		m.filepicker.Init(),
		m.spinner.Tick,
		textarea.Blink,
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case fileOpenedMsg:
		m.file = msg.file
		m.statusText = "Extracting metadata"
		return m, tea.Batch(
			extractMetadataCmd(m.file),
			m.spinner.Tick,
		)

	case metadataExtractedMsg:
		m.fileMetadata = msg.metadata
		m.statusText = ""
		m.err = nil
		m.nextStep()
		return m, nil

	case smallFilePayloadMsg:
		m.filePayload = msg.payload
		m.nextStep()
		return m, nil

	case errMsg:
		m.err = msg.error
		switch m.step {
		case StepReadyingPublicKey:
			m.rawPublicKey = ""
			m.publicKey = nil
		case StepReadyingFile:
			m.selectedFile = ""
			m.file = nil
			m.fileMetadata = pkg.FileMetadata{}
		default:
			// No default action
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
		//case tea.KeyEnter:
		//	switch m.step {
		//	case StepAwaitingPublicKey:
		//		m.rawPublicKey = m.textarea.Value()
		//		m.textarea.Reset()
		//		m.nextStep()
		//	default:
		//	}

		//case "enter":
		//	//switch m.step {
		//	//case StepAwaitingFile:
		//	//	m.filepicker, cmd = m.filepicker.Update(msg)
		//	//}
		//
		//	// Todo: Buraya ayrıca Main Logic Eklenecek. Readying kısımları vs.
		//	// Todo: Burada boş olması drumunda error msg gösterelim.
		//	m.nextStep()
		//	return m, nil

		//case "backspace":
		//	if len(m.userText) > 0 {
		//		m.userText = m.userText[:len(m.userText)-1]
		//	}

		default:
			//if len(msg.String()) == 1 {
			//	m.userText += msg.String()
			//}
		}
	}

	switch m.step {
	case StepAwaitingFile:
		m.filepicker, cmd = m.filepicker.Update(msg)
		if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
			m.selectedFile = path
			m.statusText = "Opening file"
			m.nextStep()
			return m, tea.Batch(
				openFileCmd(path),
				m.spinner.Tick,
			)
		}
	case StepReadyingFile:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case StepAwaitingPublicKey:
		m.textarea, cmd = m.textarea.Update(msg)
		if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEnter {
			rawPublicKey := m.textarea.Value()
			m.rawPublicKey = rawPublicKey
			m.textarea.Reset()
			m.statusText = "Processing public key"
			m.nextStep()
			return m, tea.Batch(
				processPublicKeyCmd(m.rawPublicKey, m.file, m.fileMetadata),
				m.spinner.Tick,
			)
		}
		return m, cmd
	case StepReadyingPublicKey:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case StepReadyToSend:
		if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEnter {
			err := clipboard.WriteAll(m.filePayload)
			if err != nil {
				m.err = err
				return m, nil
			}
			m.statusText = "File payload copied to clipboard. Exiting..."
			return m, nil
		}
	default:
		// No default action
	}
	//m.nextStep()
	return m, cmd
}
