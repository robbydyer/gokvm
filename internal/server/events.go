package server

import (
	"context"

	clientpb "github.com/robbydyer/gokvm/internal/proto/client"
	hook "github.com/robotn/gohook"
)

func (s *Server) processEvent(ctx context.Context, e hook.Event) {
	var activeClient *client
	for _, c := range s.clients {
		if c.isActive {
			activeClient = c
		}
	}

	if activeClient == nil {
		s.Log.Debug("No active client, ignoring event")
		s.Log.Debugf("Ignored event: %v", e)

		return
	}

	if activeClient.clientConn == nil {
		s.Log.Warnf("active client '%s' is missing connection", activeClient.address)

		return
	}

	if e.Kind == hook.MouseMove {
	}

	if e.Kind == hook.MouseWheel {
		s.Log.Warnf("HOOK", e)
		req := &clientpb.MouseScrollRequest{
			X:         e.Rotation,
			Direction: "up",
		}
		if e.Rotation > 0 {
			req.Direction = "down"
		}
		_, err := activeClient.clientConn.MouseScroll(ctx, req)
		if err != nil {
			s.Log.Errorf("MouseScroll failed: %s", err.Error())
		}
		return
	}

	if e.Kind == hook.MouseDown {
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

		_, err := activeClient.clientConn.MouseClick(ctx, req)
		if err != nil {
			s.Log.Errorf("MouseClick failed: %s", err.Error())
		}
		return
	}

	s.Log.Debugf("UNKNOWN HOOK: %v", e)
}
