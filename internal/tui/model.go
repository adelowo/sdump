package tui

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type model struct {
	title   string
	spinner spinner.Model

	dumpURL *url.URL
	err     error

	httpRequests incomingHTTPRequestMsg
	// requestList is shown on the LHS side of the TUI
	requestList list.Model
	// viewport because it shouldn't be editable.
	// A textarea is editable and does not work for us here
	detailedRequestView viewport.Model
}

func initialModel() model {
	m := model{
		title: "Sdump",
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Line),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205"))),
		),
		requestList: list.New([]list.Item{
			item{
				title: "oops",
				desc:  "oops oops oops oops",
				ip:    "0.0.0.0",
			},
			item{
				title: "omo",
				desc:  "omo oops oops oops",
				ip:    "0.0.0.0",
			},
		}, list.NewDefaultDelegate(), 0, 0),
		detailedRequestView: viewport.New(100, 50),
	}

	m.requestList.Title = "Incoming requests"
	m.requestList.SetShowTitle(true)

	b := new(bytes.Buffer)

	err := highlightCode(b, `
		{"name": "lanre"}
	`)
	// TODO: handle this probably. TUI design is the most important
	// bit right now
	if err != nil {
		panic(err)
	}

	m.detailedRequestView.SetContent(b.String())

	return m
}

func (m model) isInitialized() bool { return m.dumpURL != nil }

func (m model) Init() tea.Cmd {
	tea.SetWindowTitle(m.title)

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return nil
	}

	m.detailedRequestView.Width = width
	m.detailedRequestView.Height = height

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

	case tea.WindowSizeMsg:

		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		m.requestList.SetSize(msg.Width-20-h, msg.Height-20-v)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, nil

		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	m.requestList, cmd = m.requestList.Update(msg)
	cmds = append(cmds, cmd)

	m.detailedRequestView, cmd = m.detailedRequestView.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
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

	return m.spinner.View() + browserHeader + strings.Repeat("\n", 5) + m.makeTable()
}

func (m model) makeTable() string {
	return lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Margin(1, 4).Render(m.requestList.View()),
		lipgloss.NewStyle().Padding(1, 4).Render(m.detailedRequestView.View()))
}
