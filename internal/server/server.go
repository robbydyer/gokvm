package server

import (
	"context"
	"os"

	hook "github.com/robotn/gohook"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	gokvmpb "github.com/robbydyer/gokvm/internal/proto/gokvm"
)

// Server ...
type Server struct {
	Log          *log.Logger
	ctx          context.Context
	clientAddrs  []string
	conns        []*grpc.ClientConn
	clients      []gokvmpb.GoKvmClient
	mouseVisible bool
}

type OptionFunc func(*Server) error

func WithLogLevel(level log.Level) OptionFunc {
	return func(s *Server) error {
		s.Log.Level = level
		return nil
	}
}

func New(ctx context.Context, opts ...OptionFunc) (*Server, error) {
	s := &Server{
		Log: &log.Logger{
			Out:   os.Stderr,
			Level: log.DebugLevel,
		},
		ctx: ctx,
	}

	for _, f := range opts {
		if err := f(s); err != nil {
			return nil, err
		}
	}

	s.Log.Info("Starting server")

	go func() {
		ev := hook.Start()

		for e := range ev {
			s.processEvent(ctx, e)
		}
	}()

	return s, nil
}

func (s *Server) CheckConnection(address string) (bool, error) {
	for _, conn := range s.conns {
		if conn.Target() == address {
			if conn.GetState() == connectivity.Ready || conn.GetState() == connectivity.Connecting {
				return true, nil
			}

			return false, nil
		}
	}

	return false, nil
}

func (s *Server) ConnectClient(address string) error {
	connected, err := s.CheckConnection(address)
	if err != nil {
		return err
	}

	if connected {
		s.Log.Debug("Client already connected", address)
		return nil
	}

	var conn *grpc.ClientConn

	conn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}

	s.conns = append(s.conns, conn)

	client := gokvmpb.NewGoKvmClient(conn)
	s.clients = append(s.clients, client)

	resp, err := client.Hello(context.Background(), &gokvmpb.HelloRequest{})
	if err != nil {
		return err
	}

	s.Log.Printf("Client %s said '%s'", address, resp.Message)

	return nil
}

func (s *Server) Shutdown() {
	hook.End()
	for _, c := range s.conns {
		c.Close()
	}
}

func (s *Server) processEvent(ctx context.Context, e hook.Event) error {
	if e.Kind == hook.MouseDown {
		s.Log.Debug("HOOK", e)
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
		for _, c := range s.clients {
			_, err := c.MouseClick(ctx, req)
			return err
		}
	}

	if e.Kind == hook.MouseWheel {
		s.Log.Debug("HOOK", e)
		req := &gokvmpb.MouseScrollRequest{
			X:         e.Rotation,
			Direction: "up",
		}
		if e.Rotation > 0 {
			req.Direction = "down"
		}
		for _, c := range s.clients {
			_, err := c.MouseScroll(ctx, req)

			return err
		}
	}

	s.Log.Debug("UNKNOWN HOOK", e)

	return nil
}
