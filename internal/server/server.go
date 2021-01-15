package server

import (
	"net"

	"google.golang.org/grpc"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"

	gokvmpb "github.com/robbydyer/gokvm/internal/proto/gokvm"
)

// Server ...
type Server struct {
	Log io.Writer
	conn *grpc.ClientConn
	client gokvmpb.GoKvmClient
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

	return nil
}

func (s *Server) Shutdown() {
	s.conn.Close()
	robotgo.EventEnd()
}

func (s *Server) AddHooks() error {
	robotgo.EventHook(hook.MouseDown, []string{"left"}, func(e hook.Event) {
		req := &gokvmpb.MouseClickRequest{
			Button: "left",
			Double: false,
		}
		_, err := s.client.ClickMouse(context.Background(), req)
		if err != nil {
			_, _ = s.Log.Write([]byte(fmt.Sprintf("Left mouse click failed: %s", err.Error()))
		}
	})

	e := robotgo.EventStart()
	<-robotgo.EventProcess(e)

	return nil
}