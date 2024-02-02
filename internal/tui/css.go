package tui

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	color          = lipgloss.AdaptiveColor{Light: "#111222", Dark: "#FAFAFA"}
	feintColor     = lipgloss.AdaptiveColor{Light: "#333333", Dark: "#888888"}
	faintBuleColor = lipgloss.Color("#428BCA")

	errorStyle = lipgloss.NewStyle().BorderForeground(lipgloss.Color("9")).
			Border(lipgloss.RoundedBorder()).
			Align(lipgloss.Center).
			Margin(0, 0, 0, 1).
			Padding(0, 2, 0, 2)

	defaultTextStyle = lipgloss.NewStyle().Foreground(color)
)

func getTableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	return s
}

func showError(err error) string {
	return errorStyle.Render(lipgloss.Place(200, 3, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, err.Error(),
			"",
			"Press Ctrl-c to shut down",
		)))
}

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

func highlightCode(w io.Writer, s, colorscheme string) error {
	err := quick.Highlight(w, s, "json", "terminal256", colorscheme)
	return err
}

func prettyPrintJSON(str string) (string, error) {
	var b bytes.Buffer
	if err := json.Indent(&b, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return b.String(), nil
}
