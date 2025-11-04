package send

import (
	"AirBridge/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) Init() tea.Cmd {
	m.nextStep()
	return m.filePath.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		// FilePicker Limit
		bannerHeight := lipgloss.Height(tui.AirBridgeBanner()) + 3
		spacer := 1 // banner ile picker arasında eklediğimiz boş satır
		footer := 0 // hata satırı dinamik; minimumda 0 bırakıyoruz
		availableHeight := msg.Height - (bannerHeight + spacer + footer)
		if availableHeight < 3 {
			availableHeight = 3
		}
		m.filePath.SetHeight(availableHeight)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		//case "enter":
		//	//switch m.step {
		//	//case StepAwaitingFile:
		//	//	m.filePath, cmd = m.filePath.Update(msg)
		//	//}
		//
		//	// Todo: Buraya ayrıca Main Logic Eklenecek. Readying kısımları vs.
		//	// Todo: Burada boş olması drumunda error msg gösterelim.
		//	m.nextStep()
		//	return m, nil

		case "backspace":
			if len(m.userText) > 0 {
				m.userText = m.userText[:len(m.userText)-1]
			}

		default:
			if len(msg.String()) == 1 {
				m.userText += msg.String()
			}
		}
	}

	switch m.step {
	case StepAwaitingFile:
		m.filePath, cmd = m.filePath.Update(msg)
		if didSelect, path := m.filePath.DidSelectFile(msg); didSelect {
			m.selectedFile = path
			m.nextStep()
		}

	default:
		// No default action
	}

	return m, cmd
}
