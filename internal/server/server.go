package server

import (
	"context"
	"fmt"
	"io"
	"runtime"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"google.golang.org/grpc"

	gokvmpb "github.com/robbydyer/gokvm/internal/proto/gokvm"
)

// Server ...
type Server struct {
	Log          io.Writer
	conn         *grpc.ClientConn
	client       gokvmpb.GoKvmClient
	mouseVisible bool
}

func (s *Server) ConnectClient(address string) error {
	var conn *grpc.ClientConn
	var err error

	conn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}

	s.conn = conn

	s.client = gokvmpb.NewGoKvmClient(s.conn)

	resp, err := s.client.Hello(context.Background(), &gokvmpb.HelloRequest{})
	if err != nil {
		return err
	}

	_, _ = s.Log.Write([]byte(fmt.Sprintf("Client said '%s'\n", resp.Message)))

	ev := hook.Start()
	defer hook.End()

	ctx := context.Background()
	for e := range ev {
		s.processEvent(ctx, e)
	}

	for {
		runtime.Gosched()
	}

}

func (s *Server) Shutdown() {
	s.conn.Close()
	robotgo.EventEnd()
}

func (s *Server) processEvent(ctx context.Context, e hook.Event) error {
	if e.Kind == hook.MouseDown {
		_, _ = s.Log.Write([]byte(fmt.Sprintf("HOOK: %v\n", e)))
		req := &gokvmpb.MouseClickRequest{}
		for k, v := range hook.MouseMap {
			if e.Button == v {
				req.Button = k
				break
			}
		}
		if e.Clicks > 1 {
			req.Double = true
		}
		_, err := s.client.MouseClick(ctx, req)
		return err
	}

	if e.Kind == hook.MouseWheel {
		_, _ = s.Log.Write([]byte(fmt.Sprintf("HOOK: %v\n", e)))
		req := &gokvmpb.MouseScrollRequest{
			X:         e.Rotation,
			Direction: "up",
		}
		if e.Rotation > 0 {
			req.Direction = "down"
		}
		_, err := s.client.MouseScroll(ctx, req)

		return err
	}

	return nil
}
