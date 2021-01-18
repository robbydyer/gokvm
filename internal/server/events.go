package server

import (
	"context"

	clientpb "github.com/robbydyer/gokvm/internal/proto/client"
	hook "github.com/robotn/gohook"
)

func (s *Server) processEvent(ctx context.Context, e hook.Event) {
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
		for _, c := range s.clients {
			c := c
			go func() {
				_, err := c.MouseClick(ctx, req)
				if err != nil {
					s.Log.Errorf("MouseClick failed: %s", err.Error())
				}
			}()
		}
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
		for _, c := range s.clients {
			c := c
			go func() {
				_, err := c.MouseScroll(ctx, req)
				if err != nil {
					s.Log.Errorf("MouseScroll failed: %s", err.Error())
				}
			}()
		}
	}

	s.Log.Debugf("UNKNOWN HOOK: %v", e)
}
