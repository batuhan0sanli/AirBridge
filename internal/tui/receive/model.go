package receive

import (
	"AirBridge/internal/tui"
	"crypto/rsa"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

type Step int

const (
	StepUndefined Step = iota
	StepGeneratingKey
	StepAwaitingPayload
	StepDecrypting
	StepSuccess
)

type Model struct {
	tui.Window
	step Step

	spinner  spinner.Model
	textarea textarea.Model

	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	encodedKey string

	payload string

	statusText string
	err        error
}

// InitialModel initializes the receive model with default values.
func InitialModel() *Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	ta := textarea.New()
	ta.Placeholder = "Paste Base64 encoded payload here ..."
	ta.ShowLineNumbers = false
	ta.Focus()

	window := tui.Window{}

	return &Model{
		Window:   window,
		step:     StepUndefined,
		spinner:  s,
		textarea: ta,
		err:      nil,
	}
}

func (m *Model) nextStep() {
	if m.privateKey == nil {
		m.step = StepGeneratingKey
	} else if m.payload == "" {
		m.step = StepAwaitingPayload
	} else if m.statusText == "Decrypting..." {
		m.step = StepDecrypting
	} else {
		m.step = StepSuccess
	}
}

func (m *Model) resetError() {
	m.err = nil
}
