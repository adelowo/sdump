package tui

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/adelowo/sdump/config"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/r3labs/sse/v2"
	"golang.org/x/term"
)

type model struct {
	title   string
	spinner spinner.Model

	cfg        *config.Config
	dumpURL    *url.URL
	pubChannel string
	err        error

	requestList list.Model
	httpClient  *http.Client

	sseClient           *sse.Client
	receiveChan         chan item
	detailedRequestView viewport.Model

	detailedRequestViewBuffer *bytes.Buffer
}

func initialModel(cfg *config.Config) model {
	m := model{
		title: "Sdump",
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Line),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205"))),
		),
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: time.Minute,
		},
		requestList:               list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		detailedRequestView:       viewport.New(100, 50),
		detailedRequestViewBuffer: bytes.NewBuffer(nil),
		sseClient:                 sse.NewClient(fmt.Sprintf("%s/events", cfg.HTTP.Domain)),
		receiveChan:               make(chan item),
	}

	m.requestList.Title = "Incoming requests"
	m.requestList.SetShowTitle(true)
	m.requestList.SetFilteringEnabled(false)

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

	return tea.Batch(m.spinner.Tick,
		m.createEndpoint)
}

func (m model) listenForNextItem() tea.Msg {
	err := m.sseClient.Subscribe(m.pubChannel, func(msg *sse.Event) {
		var i item

		if err := json.NewDecoder(bytes.NewBuffer(msg.Data)).Decode(&i); err != nil {
			panic(err)
		}

		m.receiveChan <- i
	})
	if err != nil {
		panic(err)
	}

	return nil
}

func (m model) waitForNextItem() tea.Msg {
	return ItemMsg{item: <-m.receiveChan}
}

func (m model) createEndpoint() tea.Msg {
	// err can be safely ignored
	req, _ := http.NewRequest(http.MethodPost,
		m.cfg.HTTP.Domain,
		strings.NewReader("{}"))

	req.Header.Add("Content-Type", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	var response struct {
		URL struct {
			HumanReadableEndpoint string `json:"human_readable_endpoint,omitempty"`
		} `json:"url,omitempty"`
		SSE struct {
			Channel string `json:"channel,omitempty"`
		} `json:"sse,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		panic(err)
	}

	return DumpURLMsg{
		URL:        response.URL.HumanReadableEndpoint,
		SSEChannel: response.SSE.Channel,
	}
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

		m.pubChannel = msg.SSEChannel
		go m.listenForNextItem()
		return m, m.waitForNextItem

	case ItemMsg:

		m.requestList.InsertItem(0, msg.item)

		if err := highlightCode(m.detailedRequestViewBuffer, msg.item.Request.Body); err != nil {
			panic(err)
		}

		m.detailedRequestView.SetContent(m.detailedRequestViewBuffer.String())
		m.detailedRequestViewBuffer.Reset()

		return m, m.waitForNextItem

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
Waiting for requests on %s .. Ctrl-j/k or arrow up and down to navigate requests`, m.dumpURL), true),
		))

	return m.spinner.View() + browserHeader + strings.Repeat("\n", 5) + m.makeTable()
}

func (m model) makeTable() string {
	return lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Margin(1, 4).Render(m.requestList.View()),
		lipgloss.NewStyle().Padding(1, 4).Render(m.detailedRequestView.View()))
}
