package main

import (
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/robbydyer/gokvm/internal/client"
	gokvmpb "github.com/robbydyer/gokvm/internal/proto/gokvm"
)

type clientCmd struct {
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
	gokvmpb.RegisterGoKvmServer(grpcServer, &s)

	if err := grpcServer.Serve(l); err != nil {
		panic(err)
	}

	return nil
}
