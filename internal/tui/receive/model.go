package receive

import (
	"AirBridge/internal/crypto"
	"AirBridge/internal/tui"
	"crypto/rsa"
	"fmt"

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

	payload     string
	payloadPath string
	deleteFile  bool

	statusText string
	err        error
}

// InitialModel initializes the receive model with default values.
func InitialModel(initialPrivKeyPEM []byte, initialPayload string, payloadPath string, deleteFile bool) *Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	ta := textarea.New()
	ta.Placeholder = "Paste Base64 encoded payload here ..."
	ta.ShowLineNumbers = false
	ta.Focus()

	window := tui.Window{}

	var privateKey *rsa.PrivateKey
	var publicKey *rsa.PublicKey
	var encodedKey string
	var step = StepGeneratingKey
	var statusText string
	var err error

	if len(initialPrivKeyPEM) > 0 {
		privateKey, err = crypto.DecodeRSAPrivateKey(initialPrivKeyPEM)
		if err == nil {
			publicKey = &privateKey.PublicKey
			encodedKey, err = crypto.EncodeRSAPublicKey(publicKey)
			if err == nil {
				step = StepAwaitingPayload
			}
		} else {
			// If key decoding fails, fallback to generation and warn user
			statusText = fmt.Sprintf("⚠️ Invalid private key provided: %v. Generating new key...", err)
			privateKey = nil
			publicKey = nil
			err = nil // Clear error to avoid fatal error display
		}
	}

	if initialPayload != "" {
		ta.SetValue(initialPayload)
		// If we already have the key, we can proceed
		if step == StepAwaitingPayload {
			statusText = "Decrypting..."
			// We will rely on Update() init logic or nextStep logic to pick this up,
			// or we can set it here directly if we trust the flow.
			// Let's keep it clean: if we have payload, we update the model state.
		}
	}

	return &Model{
		Window:      window,
		step:        step,
		spinner:     s,
		textarea:    ta,
		privateKey:  privateKey,
		publicKey:   publicKey,
		encodedKey:  encodedKey,
		payload:     initialPayload,
		payloadPath: payloadPath,
		deleteFile:  deleteFile,
		statusText:  statusText,
		err:         err,
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
