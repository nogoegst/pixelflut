package handler

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"image/color"
	"net"
	"strconv"
	"strings"

	"github.com/nogoegst/pixelflut/framebuffer"
)

type SetRequest struct {
	x, y int
	c    *color.RGBA
}

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
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		sp := strings.Split(line, " ")
		switch sp[0] {
		case "PX":
			r, err := parseSetRequest(line)
			if err != nil {
				return err
			}
			h.fb.SetColor(r.x, r.y, r.c)
		case "SIZE":
			x, y := h.fb.Bounds()
			fmt.Fprintf(conn, "SIZE %d %d\n", x, y)
		default:
			return fmt.Errorf("unsupported command: %v", sp[0])
		}
	}
	return nil
}

func parseSetRequest(s string) (*SetRequest, error) {
	sp := strings.Split(s, " ")

	if len(sp) != 4 {
		return nil, fmt.Errorf("invalid command")
	}

	x, err := strconv.Atoi(sp[1])
	if err != nil {
		return nil, err
	}
	y, err := strconv.Atoi(sp[2])
	if err != nil {
		return nil, err
	}
	rgba, err := hex.DecodeString(sp[3])
	if err != nil {
		return nil, err
	}
	a := uint8(0)
	switch len(rgba) {
	case 3:
		a = 255
	case 4:
		a = uint8(rgba[3])
	default:
		return nil, fmt.Errorf("invalid command")
	}
	req := &SetRequest{
		x: x,
		y: y,
		c: &color.RGBA{
			R: uint8(rgba[0]),
			G: uint8(rgba[1]),
			B: uint8(rgba[2]),
			A: a,
		},
	}

	return req, nil
}
