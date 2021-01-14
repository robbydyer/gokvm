package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"github.com/spf13/cobra"

	gokvmpb "github.com/robbydyer/gokvm/internal/proto/gokvm"
	"github.com/robbydyer/gokvm/internal/client"
)

type clientCmd struct {}

func newClientCmd() *cobra.Command {
	c := clientCmd{}

	cmd := &cobra.Command{
		Use: "client",
		Short: "Client listening service",
		RunE: c.client,
	}

	return cmd
}

func (c *clientCmd) client(cmd *cobra.Command, args []string) error {
	l, err := net.Listen("tcp", ":10000")
	if err != nil {
		return fmt.Errorf("failed to start net listener: %w", err)
	}

	s := client.Client{}
	grpcServer := grpc.NewServer()
	gokvmpb.RegisterGoKvmServer(grpcServer, &s)

	if err := grpcServer.Serve(l); err != nil {
		panic(err)
	}

	return nil
}