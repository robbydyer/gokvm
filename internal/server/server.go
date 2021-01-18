package server

import (
	"context"
	"fmt"
	"os"
	"time"

	hook "github.com/robotn/gohook"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	clientpb "github.com/robbydyer/gokvm/internal/proto/client"
	serverpb "github.com/robbydyer/gokvm/internal/proto/server"
)

// Server ...
type Server struct {
	Log          *log.Logger
	ctx          context.Context
	clientAddrs  []string
	conns        []*grpc.ClientConn
	clients      []clientpb.ClientClient
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
		for _, addr := range s.clientAddrs {
			_ = s.ConnectClient(addr)
		}
		time.Sleep(10 * time.Second)
	}()

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

	s.clientAddrs = append(s.clientAddrs, address)

	var conn *grpc.ClientConn

	conn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}

	s.conns = append(s.conns, conn)

	client := clientpb.NewClientClient(conn)
	s.clients = append(s.clients, client)

	resp, err := client.Hello(context.Background(), &clientpb.HelloRequest{})
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

func (s *Server) RegisterClient(ctx context.Context, req *serverpb.RegisterClientRequest) (*serverpb.RegisterClientResponse, error) {
	if err := s.ConnectClient(fmt.Sprintf("%s:%d", req.Ip, req.Port)); err != nil {
		return &serverpb.RegisterClientResponse{}, err
	}

	return &serverpb.RegisterClientResponse{}, nil
}
