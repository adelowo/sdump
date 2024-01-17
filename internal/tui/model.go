package tui

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	title   string
	spinner spinner.Model

	dumpURL *url.URL
	err     error
}

func initialModel() model {
	return model{
		title: "Sdump",
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Line),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205"))),
		),
	}
}

func (m model) isInitialized() bool { return m.dumpURL != nil }

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, func() tea.Msg {
		time.Sleep(time.Second * 2)
		return DumpURLMsg{
			URL: "https://lanre.wtf",
		}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case DumpURLMsg:
		var err error
		if strings.Trim(msg.URL, "") == "" {
			m.err = errors.New("an error occurred while setting up URL")
			return m, cmd
		}

		m.dumpURL, err = url.Parse(msg.URL)
		if err != nil {
			m.err = err
			return m, cmd
		}

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, nil

		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf(`%s. Please click CTRL+C to quit...%v`,
			strings.Repeat("❌", 10), m.err)
	}

	if !m.isInitialized() {
		return lipgloss.Place(
			200, 3,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Center,
				boldenString("Generating your URL... press CTRL+C to quit", true),
				strings.Repeat(m.spinner.View(), 20),
			))
	}

	browserHeader := lipgloss.Place(
		200, 3,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center,
			boldenString("Inspecting incoming HTTP requests", true),
			boldenString(fmt.Sprintf(`
Waiting for requests on %s.. Ctrl-j/k or arrow up and down to navigate requests`, m.dumpURL), true),
		))

	return m.spinner.View() + browserHeader
}
