package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type model struct {
	spinner  spinner.Model
	quitting bool
	loading  bool
	err      error
	// in the future, we will persist urls
	shouldCreateEndpoint bool
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Meter
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{spinner: s, shouldCreateEndpoint: true}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	default:
		if m.shouldCreateEndpoint {
			m.loading = true

			time.Sleep(time.Second * 8)

			m.loading = false
			return m, nil
		}
		return m, nil
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	if m.loading {
		str := fmt.Sprintf("\n\n   %s Loading your sdump URL...press q to quit\n\n", m.spinner.View())
		return str
	}

	if m.quitting {
		return "\n"
	}
	return "dEfault page here"
}
