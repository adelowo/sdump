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
	"github.com/adelowo/sdump/internal/forward"
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
	forwardHandler := forward.New()

	cmd := &cobra.Command{
		Use:   "ssh",
		Short: "Start/run the TUI app",
		RunE: func(_ *cobra.Command, _ []string) error {
			s, err := wish.NewServer(
				wish.WithAddress(fmt.Sprintf("%s:%d", cfg.SSH.Host, cfg.SSH.Port)),
				func(s *ssh.Server) error {
					s.ReversePortForwardingCallback = func(_ ssh.Context, _ string, _ uint32) bool {
						return true
					}

					s.RequestHandlers = map[string]ssh.RequestHandler{
						"tcpip-forward":        forwardHandler.HandleSSHRequest,
						"cancel-tcpip-forward": forwardHandler.HandleSSHRequest,
					}
					return nil
				},
				validateSSHPublicKey(cfg),
				wish.WithMiddleware(
					bm.Middleware(teaHandler(cfg, forwardHandler)),
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

func validateSSHPublicKey(cfg *config.Config) ssh.Option {
	sshKeys := make(map[string]gossh.PublicKey, len(cfg.SSH.AllowList))

	for _, v := range cfg.SSH.AllowList {

		pemBytes, err := os.ReadFile(v)
		if err != nil {
			log.Fatalf("could not fetch ssh key ( %s ).. %v", v, err)
		}

		publicKey, _, _, _, err := gossh.ParseAuthorizedKey(pemBytes)
		if err != nil {
			log.Fatalf("could not parse ssh key ( %s ).. %v", v, err)
		}

		sshKeys[gossh.FingerprintSHA256(publicKey)] = publicKey
	}

	return wish.WithPublicKeyAuth(func(_ ssh.Context, key ssh.PublicKey) bool {
		if len(sshKeys) == 0 {
			return true
		}

		publicKey, ok := sshKeys[gossh.FingerprintSHA256(key)]
		if !ok {
			return false
		}

		return ssh.KeysEqual(publicKey, key)
	})
}

func teaHandler(cfg *config.Config,
	fowardHandler *forward.Forwarder,
) func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		pty, _, active := s.Pty()
		if !active {
			wish.Fatalln(s, "no active terminal, skipping")
			return nil, nil
		}

		tuiModel, err := tui.New(cfg, fowardHandler,
			tui.WithWidth(pty.Window.Width),
			tui.WithHeight(pty.Window.Height),
			tui.WithRemoteAddr(s.RemoteAddr()),
			tui.WithSSHFingerPrint(gossh.FingerprintSHA256(s.PublicKey())),
		)
		if err != nil {
			wish.Fatalln(s, fmt.Errorf("%v...Could not set up TUI session", err))
			return nil, nil
		}

		return tuiModel,
			[]tea.ProgramOption{tea.WithAltScreen()}
	}
}
