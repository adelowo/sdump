package tui

import (
	"github.com/adelowo/sdump/config"
	tea "github.com/charmbracelet/bubbletea"
)

type App struct{}

func New(cfg *config.Config) *tea.Program {
	p := tea.NewProgram(InitialModel(cfg), tea.WithAltScreen())

	return p
}
