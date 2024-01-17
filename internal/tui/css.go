package tui

import "github.com/charmbracelet/lipgloss"

var (
	color        = lipgloss.AdaptiveColor{Light: "#111222", Dark: "#FAFAFA"}
	primaryColor = lipgloss.Color("#4636f5")
	greenColor   = lipgloss.Color("#9dcc3a")
	redColor     = lipgloss.Color("#ff0000")
	whiteColor   = lipgloss.Color("#ffffff")
	blackColor   = lipgloss.Color("#000000")
	orangeColor  = lipgloss.Color("#D3A347")
	feintColor   = lipgloss.AdaptiveColor{Light: "#333333", Dark: "#888888"}
	irisColor    = lipgloss.Color("#5D5FEF")
	fuschiaColor = lipgloss.Color("#EF5DA8")

	defaultTextStyle = lipgloss.NewStyle().Foreground(color)
)

func makeString(s string, withFeint bool) string {
	style := defaultTextStyle.Copy()

	if withFeint {
		style = style.Foreground(feintColor)
	}

	return style.Render(s)
}

func boldenString(s string, withFeint bool) string {
	style := defaultTextStyle.Copy().Bold(true)

	if withFeint {
		style = style.Foreground(feintColor)
	}

	return style.Render(s)
}
