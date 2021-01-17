package client

import (
	"context"

	"github.com/go-vgo/robotgo"
	log "github.com/sirupsen/logrus"

	gokvmpb "github.com/robbydyer/gokvm/internal/proto/gokvm"
)

// Client ...
type Client struct {
	Log        *log.Logger
	ScrollOnly bool
}

func (c *Client) Hello(ctx context.Context, req *gokvmpb.HelloRequest) (*gokvmpb.HelloResponse, error) {
	c.Log.Info("Got Hello Request", req)
	return &gokvmpb.HelloResponse{
		Message: "Hello there",
	}, nil
}

func (c *Client) MouseClick(ctx context.Context, req *gokvmpb.MouseClickRequest) (*gokvmpb.MouseClickResponse, error) {
	if c.ScrollOnly {
		c.Log.Info("Ignoring MouseClick")

		return &gokvmpb.MouseClickResponse{}, nil
	}

	c.Log.Debug("Got MouseClick", req)

	robotgo.MouseClick(req.Button, req.Double)

	return &gokvmpb.MouseClickResponse{}, nil
}

func (c *Client) MouseMove(ctx context.Context, req *gokvmpb.MouseMoveRequest) (*gokvmpb.MouseMoveResponse, error) {
	if c.ScrollOnly {
		c.Log.Info("Ignoring MouseMove")

		return &gokvmpb.MouseMoveResponse{}, nil
	}
	c.Log.Debug("Got MouseMove", req)
	robotgo.MoveSmoothRelative(int(req.Xrel), int(req.Yrel))

	return &gokvmpb.MouseMoveResponse{}, nil
}

func (c *Client) MouseScroll(ctx context.Context, req *gokvmpb.MouseScrollRequest) (*gokvmpb.MouseScrollResponse, error) {
	c.Log.Debug("Got MouseScroll", req)
	robotgo.ScrollMouse(int(req.X), req.Direction)

	return &gokvmpb.MouseScrollResponse{}, nil
}
