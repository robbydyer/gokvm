package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/robbydyer/gokvm/internal/client"
	clientpb "github.com/robbydyer/gokvm/internal/proto/client"
	serverpb "github.com/robbydyer/gokvm/internal/proto/server"
)

type clientCmd struct {
	server     string
	scrollOnly bool
	port       int
	udp        bool
	location   string
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
	f.StringVar(&c.location, "relative-location", "right", "Relative location of this client to the server screen")
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

	s, err := client.New(c.port)
	if err != nil {
		return err
	}
	s.Log.Level = log.DebugLevel
	s.InternalIPSubnet = "192.168.1"
	s.ScrollOnly = c.scrollOnly

	grpcServer := grpc.NewServer()
	clientpb.RegisterClientServer(grpcServer, s)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		s.Log.Info("Shutting down client")
		grpcServer.GracefulStop()
	}()

	ctx := context.Background()

	loc, err := locationToLocation(c.location)
	if err != nil {
		return err
	}

	if err := s.ConnectServer(ctx, c.server, loc); err != nil {
		return err
	}

	if err := grpcServer.Serve(l); err != nil {
		panic(err)
	}

	return nil
}

func locationToLocation(location string) (serverpb.Location, error) {
	switch strings.ToLower(location) {
	case "left":
		return serverpb.Location_LEFT, nil
	case "right":
		return serverpb.Location_RIGHT, nil
	case "up":
		return serverpb.Location_UP, nil
	case "down":
		return serverpb.Location_DOWN, nil
	}

	return 0, fmt.Errorf("invalid location %s", location)
}
