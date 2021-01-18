package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/robbydyer/gokvm/internal/server"
)

type serverCmd struct {
	clients []string
}

func newServerCmd() *cobra.Command {
	c := serverCmd{}

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Server listening service",
		RunE:  c.server,
	}

	f := cmd.Flags()

	f.StringArrayVar(&c.clients, "client", []string{"192.168.1.34:10000"}, "clients to connect to")
	return cmd
}

func (c *serverCmd) server(cmd *cobra.Command, args []string) error {

	s, err := server.New(context.Background(),
		server.WithLogLevel(log.DebugLevel),
	)
	if err != nil {
		return err
	}

	go func() {
		for _, clientAddr := range c.clients {
			_ = s.ConnectClient(clientAddr)
		}
		time.Sleep(30 * time.Second)
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	s.Log.Info("Shutting down server")
	s.Shutdown()

	return nil
}
