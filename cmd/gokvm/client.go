package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/robbydyer/gokvm/internal/client"
	clientpb "github.com/robbydyer/gokvm/internal/proto/client"
)

type clientCmd struct {
	server     string
	scrollOnly bool
	port       int
	udp        bool
}

func newClientCmd() *cobra.Command {
	c := clientCmd{}

	cmd := &cobra.Command{
		Use:   "client",
		Short: "Client listening service",
		RunE:  c.client,
	}

	f := cmd.Flags()
	f.StringVar(&c.server, "server-address", "", "[IP]:[PORT] of server to connect to")
	f.BoolVar(&c.scrollOnly, "scroll-only", false, "Tells the client to ignore all commands except for Mouse scroll")
	f.IntVar(&c.port, "port", 10000, "Listen port")
	f.BoolVar(&c.udp, "udp", false, "Use UDP")
	return cmd
}

func (c *clientCmd) client(cmd *cobra.Command, args []string) error {
	proto := "tcp"
	if c.udp {
		proto = "udp"
	}
	l, err := net.Listen(proto, fmt.Sprintf(":%d", c.port))
	if err != nil {
		return fmt.Errorf("failed to start net listener: %w", err)
	}

	s := client.Client{
		Log: &log.Logger{
			Out:   os.Stderr,
			Level: log.DebugLevel,
		},
		ScrollOnly: c.scrollOnly,
	}
	grpcServer := grpc.NewServer()
	clientpb.RegisterClientServer(grpcServer, &s)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		s.Log.Info("Shutting down client")
		grpcServer.GracefulStop()
	}()

	if err := grpcServer.Serve(l); err != nil {
		panic(err)
	}

	return nil
}
