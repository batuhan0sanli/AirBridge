package send

import (
	"AirBridge/internal/tui"
	"AirBridge/pkg"
	"crypto/rsa"
	"os"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

type Step int

const (
	StepUndefined Step = iota
	StepAwaitingFile
	StepReadyingFile
	StepAwaitingPublicKey
	StepReadyingPublicKey
	StepReadyToSend
)

type Model struct {
	tui.Window
	step Step

	filepicker filepicker.Model
	spinner    spinner.Model
	textarea   textarea.Model

	selectedFile string
	file         *os.File
	fileMetadata pkg.FileMetadata

	rawPublicKey string
	publicKey    *rsa.PublicKey

	filePayload string

	statusText string
	err        error
}

func InitialModel() *Model {
	fp := filepicker.New()
	// Todo: Add styles
	styles := filepicker.DefaultStyles()
	fp.Styles = styles

	window := tui.Window{}
	err := tui.ErrTest

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	ta := textarea.New()
	ta.Placeholder = "Paste Base64 encoded public key here ..."
	ta.ShowLineNumbers = false
	ta.Focus()
	return &Model{
		Window:     window,
		step:       StepUndefined,
		spinner:    s,
		filepicker: fp,
		textarea:   ta,
		err:        err}
}

func (m *Model) nextStep() {
	if m.selectedFile == "" {
		m.step = StepAwaitingFile
	} else if m.file == nil {
		m.step = StepReadyingFile
	} else if m.rawPublicKey == "" {
		m.step = StepAwaitingPublicKey
	} else if m.filePayload != "" {
		m.step = StepReadyToSend
	} else if m.publicKey == nil || m.file != nil && m.publicKey != nil && m.filePayload == "" {
		m.step = StepReadyingPublicKey
	} else {
		m.step = StepUndefined
	}
}

func (m *Model) resetError() {
	m.err = nil
}
