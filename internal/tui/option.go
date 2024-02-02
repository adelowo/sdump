package tui

import "github.com/adelowo/sdump/config"

type Option func(*model)

func WithConfig(cfg *config.Config) Option {
	return func(m *model) {
		m.cfg = cfg
	}
}

func WithHeight(height int) Option {
	return func(m *model) {
		m.height = height
	}
}

func WithWidth(width int) Option {
	return func(m *model) {
		m.width = width
	}
}

func WithSSHFingerPrint(fingerPrint string) Option {
	return func(m *model) {
		m.sshFingerPrint = fingerPrint
	}
}
