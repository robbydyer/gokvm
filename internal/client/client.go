package client

import (
	"io"
	"context"

	gokvmpb "github.com/robbydyer/gokvm/internal/proto/gokvm"
)

// Client ...
type Client struct {
	Log io.Writer
}

func (s *Client) MouseClick(ctx context.Context, req *gokvmpb.MouseClickRequest) (*gokvmpb.MouseClickResponse, error) {
	_, _ = s.Log.Write([]byte("Got MouseClick\n"))

	return nil, nil
}