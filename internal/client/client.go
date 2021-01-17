package client

import (
	"context"
	"fmt"
	"io"

	"github.com/go-vgo/robotgo"

	gokvmpb "github.com/robbydyer/gokvm/internal/proto/gokvm"
)

// Client ...
type Client struct {
	Log io.Writer
}

func (c *Client) Hello(ctx context.Context, req *gokvmpb.HelloRequest) (*gokvmpb.HelloResponse, error) {
	c.Log.Write([]byte("Got Hello request\n"))
	return &gokvmpb.HelloResponse{
		Message: "Hello there",
	}, nil
}

func (c *Client) MouseClick(ctx context.Context, req *gokvmpb.MouseClickRequest) (*gokvmpb.MouseClickResponse, error) {
	go c.Log.Write([]byte(fmt.Sprintf("Got MouseClick: %v\n", req)))

	robotgo.MouseClick(req.Button, req.Double)

	return &gokvmpb.MouseClickResponse{}, nil
}

func (c *Client) MouseMove(ctx context.Context, req *gokvmpb.MouseMoveRequest) (*gokvmpb.MouseMoveResponse, error) {
	_, _ = c.Log.Write([]byte("Got MouseMove\n"))
	robotgo.MoveSmoothRelative(int(req.Xrel), int(req.Yrel))

	return &gokvmpb.MouseMoveResponse{}, nil
}

func (c *Client) MouseScroll(ctx context.Context, req *gokvmpb.MouseScrollRequest) (*gokvmpb.MouseScrollResponse, error) {
	go c.Log.Write([]byte(fmt.Sprintf("Got MouseScroll: %v\n", req)))
	robotgo.ScrollMouse(int(req.X), req.Direction)

	return &gokvmpb.MouseScrollResponse{}, nil
}
