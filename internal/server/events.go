package server

import (
	"context"

	clientpb "github.com/robbydyer/gokvm/internal/proto/client"
	hook "github.com/robotn/gohook"
)

func (s *Server) processEvent(ctx context.Context, e hook.Event) {
	if s.activeClient == "" {
		go s.Log.Debugf("No active client, igoring event: %v", e)
		return
	}

	activeClient, ok := s.clients[s.activeClient]
	if !ok {
		go s.Log.Errorf("active client is not connected, trying reconnect")
		if err := s.ConnectClient(s.activeClient); err != nil {
			go s.Log.Errorf("failed to connect to active client: %s", err.Error())
			s.activeClient = ""
			return
		}
	}

	if e.Kind == hook.MouseDown {
		s.Log.Debug("HOOK", e)
		req := &clientpb.MouseClickRequest{}
		for k, v := range hook.MouseMap {
			if e.Button == v {
				req.Button = k
				break
			}
		}
		if e.Clicks > 1 {
			req.Double = true
		}

		_, err := activeClient.MouseClick(ctx, req)
		if err != nil {
			s.Log.Errorf("MouseClick failed: %s", err.Error())
		}
		return
	}

	if e.Kind == hook.MouseWheel {
		s.Log.Debug("HOOK", e)
		req := &clientpb.MouseScrollRequest{
			X:         e.Rotation,
			Direction: "up",
		}
		if e.Rotation > 0 {
			req.Direction = "down"
		}
		_, err := activeClient.MouseScroll(ctx, req)
		if err != nil {
			s.Log.Errorf("MouseScroll failed: %s", err.Error())
		}
		return
	}

	s.Log.Debugf("UNKNOWN HOOK: %v", e)
}
