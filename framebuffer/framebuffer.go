package framebuffer

import (
	"bytes"
	"fmt"
	"image"
	"io"

	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
)

type FrameBuffer struct {
	buffer *LockableImage
}

func New(w, h int) *FrameBuffer {
	bounds := image.Rect(0, 0, w, h)
	image := image.NewRGBA(bounds)
	lockableImage := &LockableImage{
		Image: image,
	}
	fb := &FrameBuffer{
		buffer: lockableImage,
	}
	return fb
}

func (fb *FrameBuffer) Bounds() (int, int) {
	return fb.buffer.Image.Bounds().Dx(), fb.buffer.Image.Bounds().Dy()
}

func (fb *FrameBuffer) SetColor(x, y int, c color.Color) {
	fb.buffer.SetColor(x, y, c)
}

func (fb *FrameBuffer) Dump(w io.Writer) error {
	fb.buffer.RLock()
	defer fb.buffer.RUnlock()
	buf := bytes.NewReader(fb.buffer.Image.Pix)

	_, err := io.Copy(w, buf)
	if err != nil {
		return fmt.Errorf("dump raw image: %w", err)
	}
	return nil
}

func (fb *FrameBuffer) DumpJPEG(w io.Writer) error {
	fb.buffer.RLock()
	defer fb.buffer.RUnlock()
	options := &jpeg.Options{
		Quality: 95,
	}
	err := jpeg.Encode(w, fb.buffer.Image, options)
	if err != nil {
		return fmt.Errorf("dump JPEG %w", err)
	}
	return nil
}
