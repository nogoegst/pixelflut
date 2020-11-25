package framebuffer

import (
	"image"
	"image/color"
	"sync"
)

type LockableImage struct {
	sync.RWMutex
	Image *image.RGBA
}

func (li *LockableImage) SetColor(x, y int, c color.Color) {
	li.Lock()
	li.Image.Set(x, y, c)
	li.Unlock()
}
