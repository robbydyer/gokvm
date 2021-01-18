package server

import (
	"context"
	"fmt"
	"os"

	hook "github.com/robotn/gohook"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	clientpb "github.com/robbydyer/gokvm/internal/proto/client"
	serverpb "github.com/robbydyer/gokvm/internal/proto/server"
)

// Server ...
type Server struct {
	Log          *log.Logger
	ctx          context.Context
	clients      []*client
	activeClient string
	mouseVisible bool
}

type client struct {
	clientConn clientpb.ClientClient
	grpcConn   *grpc.ClientConn
	address    string
	location   serverpb.Location
	isActive   bool
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

func (s *Server) ConnectClient(ctx context.Context, address string, location serverpb.Location) (*client, error) {
	var thisClient *client

	for _, c := range s.clients {
		if c.address == address {
			thisClient = c
			break
		}
	}

	if thisClient != nil {
		return thisClient, nil
	}

	thisClient = &client{
		address:  address,
		location: location,
	}
	s.clients = append(s.clients, thisClient)

	var err error
	thisClient.grpcConn, err = grpc.DialContext(ctx, address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	thisClient.clientConn = clientpb.NewClientClient(thisClient.grpcConn)

	resp, err := thisClient.clientConn.Hello(context.Background(), &clientpb.HelloRequest{})
	if err != nil {
		return nil, err
	}

	go s.Log.Infof("Client %s said '%s'", address, resp.Message)

	return thisClient, nil
}

func (s *Server) Shutdown() {
	hook.End()
	for _, c := range s.clients {
		c.grpcConn.Close()
	}
}

func (s *Server) RegisterClient(ctx context.Context, req *serverpb.RegisterClientRequest) (*serverpb.RegisterClientResponse, error) {
	c, err := s.ConnectClient(ctx, fmt.Sprintf("%s:%d", req.Ip, req.Port), req.Location)
	if err != nil {
		return &serverpb.RegisterClientResponse{}, err
	}

	// TODO: Remove after auto-activation is figured out
	c.isActive = true

	return &serverpb.RegisterClientResponse{}, nil
}

func (s *Server) SetClientActive(ctx context.Context, req *serverpb.SetClientActiveRequest) (*serverpb.SetClientActiveResponse, error) {
	found := false
	for _, c := range s.clients {
		if c.address == fmt.Sprintf("%s:%d", req.Ip, req.Port) {
			found = true
		}
	}

	if !found {
		return &serverpb.SetClientActiveResponse{}, fmt.Errorf("client %s:%d does not exist", req.Ip, req.Port)
	}

	for _, c := range s.clients {
		if c.address == fmt.Sprintf("%s:%d", req.Ip, req.Port) {
			c.isActive = true
		}
		c.isActive = false
	}

	return &serverpb.SetClientActiveResponse{}, nil
}
