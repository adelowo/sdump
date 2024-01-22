package tui

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/charmbracelet/lipgloss"
)

var (
	color = lipgloss.AdaptiveColor{Light: "#111222", Dark: "#FAFAFA"}
	// primaryColor = lipgloss.Color("#4636f5")
	// greenColor   = lipgloss.Color("#9dcc3a")
	// redColor     = lipgloss.Color("#ff0000")
	// whiteColor   = lipgloss.Color("#ffffff")
	// blackColor   = lipgloss.Color("#000000")
	// orangeColor  = lipgloss.Color("#D3A347")
	feintColor = lipgloss.AdaptiveColor{Light: "#333333", Dark: "#888888"}
	// fuschiaColor   = lipgloss.Color("#EF5DA8")
	faintBuleColor = lipgloss.Color("#428BCA")

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

func highlightCode(w io.Writer, s string) error {
	// TODO: make monokai configurable
	err := quick.Highlight(w, s, "json", "terminal256", "monokai")
	return err
}

func prettyPrintJSON(str string) (string, error) {
	var b bytes.Buffer
	if err := json.Indent(&b, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return b.String(), nil
}
