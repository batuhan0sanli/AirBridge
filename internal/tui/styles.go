package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	BaseColor    = lipgloss.Color("#E0E0E0") // açık gri
	AccentColor  = lipgloss.Color("#00C6FF") // mavi-turkuaz ton
	WarningColor = lipgloss.Color("#FFAD33") // turuncu
	ErrorColor   = lipgloss.Color("#FF4D4D") // kırmızı
	SuccessColor = lipgloss.Color("#4DFF88") // yeşil
)

var FooterStyle = lipgloss.NewStyle().
	Foreground(BaseColor).
	Padding(0, 1).
	Margin(1, 0, 1, 0).
	Align(lipgloss.Left)

//Border(lipgloss.RoundedBorder()).
//BorderForeground(AccentColor)

var TitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(AccentColor).
	// Background(lipgloss.Color("#1A1A1A")).
	Padding(0, 2).
	Margin(1, 0, 1, 0).
	Align(lipgloss.Center).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(AccentColor)

var (
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(BaseColor).
			Italic(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#B0B0B0"))

	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor)

	WarningStyle = lipgloss.NewStyle().
			Foreground(WarningColor)
)

func asciiTitle() string {
	return `
           _        ____       _     _            
     /\   (_)      |  _ \     (_)   | |           
    /  \   _ _ __  | |_) |_ __ _  __| | __ _  ___ 
   / /\ \ | | '__| |  _ <| '__| |/ _' |/ _' |/ _ \
  / ____ \| | |    | |_) | |  | | (_| | (_| |  __/
 /_/    \_\_|_|    |____/|_|  |_|\__,_|\__, |\___|
                                        __/ |     
                                       |___/
`
}

func AirBridgeBanner() string {
	title := lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true).
		Render(asciiTitle())

	sub := SubtitleStyle.Render("Secure, Simple, and Fast File Transfer")

	box := TitleStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center, title, sub),
	)

	return box
}

func errorView(err error) string {
	if err != nil {
		errMsg := "⚠️  " + err.Error()
		return ErrorStyle.Render(errMsg)
	}
	return ""
}

func Header() string {
	return AirBridgeBanner()
}

func MainStyle(window Window) lipgloss.Style {
	mainStyle := lipgloss.NewStyle().
		Height(window.AvailableHeight).    // Kalan tüm yüksekliği buna ver
		MaxHeight(window.AvailableHeight). // Yüksekliği sığdığından emin ol
		Width(window.Width).               // Tam genişliği kullan
		Padding(0, 2)
	return mainStyle
}

func Footer(err error) string {
	// Todo: Change input to Model
	quitText := "(esc to quit)"
	view := lipgloss.JoinVertical(lipgloss.Left, errorView(err), quitText)
	view = FooterStyle.Render(view)
	return view
}

func View(window Window, err error, view string) string {
	banner := Header()
	footer := Footer(err)
	return lipgloss.JoinVertical(lipgloss.Left, banner, view, footer)
}
