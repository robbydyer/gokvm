package main

import (
	"context"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/robbydyer/gokvm/internal/server"
)

type serverCmd struct{}

func newServerCmd() *cobra.Command {
	c := serverCmd{}

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Server listening service",
		RunE:  c.server,
	}

	return cmd
}

func (c *serverCmd) server(cmd *cobra.Command, args []string) error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		os.Exit(1)
	}()

	s, err := server.New(context.Background(),
		server.WithLogLevel(log.DebugLevel),
	)
	if err != nil {
		return err
	}

	if err := s.ConnectClient("192.168.1.34:10000"); err != nil {
		return err
	}

	return nil
}
