package client

import (
	"context"
	"io"

	"github.com/go-vgo/robotgo"

	gokvmpb "github.com/robbydyer/gokvm/internal/proto/gokvm"
)

// Client ...
type Client struct {
	Log io.Writer
}

func (c *Client) MouseClick(ctx context.Context, req *gokvmpb.MouseClickRequest) (*gokvmpb.MouseClickResponse, error) {
	_, _ = c.Log.Write([]byte("Got MouseClick\n"))

	return nil, nil
}

func (c *Client) MouseMove(ctx context.Context, req *gokvmpb.MouseMoveRequest) (*gokvmpb.MouseMoveResponse, error) {
	robotgo.MoveSmoothRelative(int(req.Xrel), int(req.Yrel))

	return &gokvmpb.MouseMoveResponse{}, nil
}
