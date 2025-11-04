package send

import (
	"AirBridge/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) Init() tea.Cmd {
	m.nextStep()
	return tea.Batch(
		m.filepicker.Init(),
		m.spinner.Tick,
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

	case errMsg:
		m.err = msg.error
		return m, nil

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		headerH := lipgloss.Height(tui.Header())
		footerH := lipgloss.Height(tui.Footer(m.err))
		spacer := 0

		availableHeight := m.Window.Height - (headerH + footerH + spacer)
		if availableHeight < 3 {
			availableHeight = 3
		}
		m.AvailableHeight = availableHeight
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

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
			m.step = StepReadyingFile
			return m, tea.Batch(
				openFileCmd(path),
				m.spinner.Tick,
			)
		}
	case StepReadyingFile:
		var sCmd tea.Cmd
		m.spinner, sCmd = m.spinner.Update(msg)
		return m, sCmd

	default:
		// No default action
	}

	m.nextStep()
	return m, cmd
}
