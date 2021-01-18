package client

import (
	"context"
	"os"

	"github.com/go-vgo/robotgo"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	clientpb "github.com/robbydyer/gokvm/internal/proto/client"
	serverpb "github.com/robbydyer/gokvm/internal/proto/server"
	"github.com/robbydyer/gokvm/internal/util"
)

// Client ...
type Client struct {
	Log              *log.Logger
	ScrollOnly       bool
	Server           serverpb.ServerClient
	InternalIPSubnet string
	listenPort       int
}

func New(listenPort int) (*Client, error) {
	c := &Client{
		listenPort: listenPort,
		Log: &log.Logger{
			Out:   os.Stderr,
			Level: log.InfoLevel,
		},
		InternalIPSubnet: "192.168.1",
	}

	return c, nil
}

func (c *Client) ConnectServer(ctx context.Context, address string, relativeLocation serverpb.Location) error {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}

	c.Server = serverpb.NewServerClient(conn)

	ip, err := util.GetInternalIP(c.InternalIPSubnet)
	if err != nil {
		return err
	}

	c.Server.RegisterClient(ctx, &serverpb.RegisterClientRequest{
		Ip:       ip,
		Port:     int32(c.listenPort),
		Location: relativeLocation,
	})

	return nil
}

func Shutdown() {

}

func (c *Client) Hello(ctx context.Context, req *clientpb.HelloRequest) (*clientpb.HelloResponse, error) {
	c.Log.Info("Got Hello Request", req)
	return &clientpb.HelloResponse{
		Message: "Hello there",
	}, nil
}

func (c *Client) MouseClick(ctx context.Context, req *clientpb.MouseClickRequest) (*clientpb.MouseClickResponse, error) {
	if c.ScrollOnly {
		c.Log.Info("Ignoring MouseClick")

		return &clientpb.MouseClickResponse{}, nil
	}

	c.Log.Debug("Got MouseClick", req)

	robotgo.MouseClick(req.Button, req.Double)

	return &clientpb.MouseClickResponse{}, nil
}

func (c *Client) MouseMove(ctx context.Context, req *clientpb.MouseMoveRequest) (*clientpb.MouseMoveResponse, error) {
	if c.ScrollOnly {
		c.Log.Info("Ignoring MouseMove")

		return &clientpb.MouseMoveResponse{}, nil
	}
	c.Log.Debug("Got MouseMove", req)
	robotgo.MoveSmoothRelative(int(req.Xrel), int(req.Yrel))

	return &clientpb.MouseMoveResponse{}, nil
}

func (c *Client) MouseScroll(ctx context.Context, req *clientpb.MouseScrollRequest) (*clientpb.MouseScrollResponse, error) {
	c.Log.Debug("Got MouseScroll", req)
	robotgo.ScrollMouse(int(req.X), req.Direction)

	return &clientpb.MouseScrollResponse{}, nil
}
