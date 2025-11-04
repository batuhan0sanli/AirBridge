package send

import (
	"AirBridge/internal/tui"
	"AirBridge/pkg"
	"crypto/rsa"
	"os"

	"github.com/charmbracelet/bubbles/filepicker"
)

type Step int

const (
	StepUndefined Step = iota
	StepAwaitingFile
	StepReadyingFile
	StepAwaitingPublicKey
	StepReadyingPublicKey
	StepEncryptingFile
	StepReadyToSend
)

type Model struct {
	tui.Window
	step         Step
	filepicker   filepicker.Model
	selectedFile string
	file         *os.File
	fileMetadata pkg.FileMetadata
	rawPublicKey string
	publicKey    *rsa.PublicKey
	userText     string
	statusText   string
	err          error
}

func InitialModel() *Model {
	fp := filepicker.New()
	// current_path := pkg.GetCurrentDirectory()
	//_, err := os.Getwd()
	//if err != nil {
	//	log.Fatalf("Çalışma dizini alınamadı: %v", err)
	//}
	// filepicker.CurrentDirectory = cwd

	// Todo: Add styles
	styles := filepicker.DefaultStyles()
	fp.Styles = styles

	window := tui.Window{}
	err := tui.ErrTest

	return &Model{Window: window, step: StepUndefined, filepicker: fp, err: err}
}

func (m *Model) nextStep() {
	if m.selectedFile == "" {
		m.step = StepAwaitingFile
	} else if m.file == nil {
		m.step = StepReadyingFile
	} else if m.rawPublicKey == "" {
		m.step = StepAwaitingPublicKey
	} else if m.publicKey == nil {
		m.step = StepReadyingPublicKey
	} else if m.file != nil && m.publicKey != nil && m.step != StepReadyToSend {
		m.step = StepEncryptingFile
	} else if m.file != nil && m.publicKey != nil && m.step == StepEncryptingFile {
		m.step = StepReadyToSend
	} else {
		m.step = StepUndefined
	}
}
