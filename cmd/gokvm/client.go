package main

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/robbydyer/gokvm/internal/client"
	gokvmpb "github.com/robbydyer/gokvm/internal/proto/gokvm"
)

type clientCmd struct{}

func newClientCmd() *cobra.Command {
	c := clientCmd{}

	cmd := &cobra.Command{
		Use:   "client",
		Short: "Client listening service",
		RunE:  c.client,
	}

	return cmd
}

func (c *clientCmd) client(cmd *cobra.Command, args []string) error {
	l, err := net.Listen("tcp", ":10000")
	if err != nil {
		return fmt.Errorf("failed to start net listener: %w", err)
	}

	s := client.Client{
		Log: os.Stdout,
	}
	grpcServer := grpc.NewServer()
	gokvmpb.RegisterGoKvmServer(grpcServer, &s)

	if err := grpcServer.Serve(l); err != nil {
		panic(err)
	}

	return nil
}
