package rgb

import (
	"fmt"
	"net"
	"time"

	"github.com/nogoegst/pixelflut/framebuffer"
)

type Handler struct {
	fb *framebuffer.FrameBuffer
}

func New(fb *framebuffer.FrameBuffer) *Handler {
	h := &Handler{
		fb: fb,
	}
	return h
}

func (h *Handler) Handle(conn net.Conn) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	for range ticker.C {
		err := h.fb.Dump(conn)
		if err != nil {
			return fmt.Errorf("writing framebuffer: %w", err)
		}
	}
	return nil
}
