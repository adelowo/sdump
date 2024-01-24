package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adelowo/sdump/config"
	"github.com/adelowo/sdump/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/spf13/cobra"
	gossh "golang.org/x/crypto/ssh"
)

func createSSHCommand(rootCmd *cobra.Command, cfg *config.Config) {
	cmd := &cobra.Command{
		Use:   "ssh",
		Short: "Start/run the TUI app",
		RunE: func(_ *cobra.Command, _ []string) error {
			s, err := wish.NewServer(
				wish.WithAddress(fmt.Sprintf("%s:%d", cfg.SSH.Host, cfg.SSH.Port)),
				wish.WithMiddleware(
					bm.Middleware(teaHandler(cfg)),
					lm.Middleware(),
				),
			)
			if err != nil {
				return err
			}

			for _, v := range cfg.SSH.Identities {

				pemBytes, err := os.ReadFile(v)
				if err != nil {
					return err
				}

				signer, err := gossh.ParsePrivateKey(pemBytes)
				if err != nil {
					return err
				}

				s.AddHostKey(signer)
			}

			done := make(chan os.Signal, 1)
			signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			log.Info("Starting SSH server", "host", cfg.SSH.Host, "port", cfg.SSH.Port)

			go func() {
				if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
					log.Error("could not start server", "error", err)
					done <- nil
				}
			}()

			<-done
			log.Info("Stopping SSH server")
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer func() { cancel() }()
			if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
				log.Error("could not stop server", "error", err)
				return err
			}

			return nil
		},
	}

	rootCmd.AddCommand(cmd)
}

func teaHandler(cfg *config.Config) func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		_, _, active := s.Pty()
		if !active {
			wish.Fatalln(s, "no active terminal, skipping")
			return nil, nil
		}

		return tui.InitialModel(cfg), []tea.ProgramOption{tea.WithAltScreen()}
	}
}
