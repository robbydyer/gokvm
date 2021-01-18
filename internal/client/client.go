package client

import (
	"context"

	"github.com/go-vgo/robotgo"
	log "github.com/sirupsen/logrus"

	clientpb "github.com/robbydyer/gokvm/internal/proto/client"
)

// Client ...
type Client struct {
	Log        *log.Logger
	ScrollOnly bool
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
