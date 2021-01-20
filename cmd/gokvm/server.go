package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	serverpb "github.com/robbydyer/gokvm/internal/proto/server"
	"github.com/robbydyer/gokvm/internal/server"
)

type serverCmd struct {
	rArgs *rootArgs
	port  int
	udp   bool
}

func newServerCmd(args *rootArgs) *cobra.Command {
	c := serverCmd{
		rArgs: args,
	}

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Server listening service",
		RunE:  c.server,
	}

	f := cmd.Flags()

	f.BoolVar(&c.udp, "udp", false, "Use udp")
	f.IntVar(&c.port, "port", 10000, "Listen port")
	return cmd
}

func (c *serverCmd) server(cmd *cobra.Command, args []string) error {
	proto := "tcp"
	if c.udp {
		proto = "udp"
	}
	l, err := net.Listen(proto, fmt.Sprintf(":%d", c.port))
	if err != nil {
		return fmt.Errorf("failed to start net listener: %w", err)
	}

	s, err := server.New()
	if err != nil {
		return err
	}

	lvl, err := log.ParseLevel(c.rArgs.logLevel)
	if err != nil {
		return err
	}
	s.Log.Level = lvl

	grpcServer := grpc.NewServer()
	serverpb.RegisterServerServer(grpcServer, s)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		s.Shutdown()
		grpcServer.GracefulStop()
	}()

	s.Log.Info("Starting server")
	if err := grpcServer.Serve(l); err != nil {
		panic(err)
	}

	return nil
}
