package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/anthonynsimon/bild/adjust"
)

func RGBA(c *color.RGBA) string {
	if c.A == 255 {
		return fmt.Sprintf("%02x %02x %02x", c.R, c.G, c.B)
	}
	return fmt.Sprintf("%02x%02x%02x%02x", c.R, c.G, c.B, c.A)
}

type Conn struct {
	c net.Conn
}

func (c *Conn) Size() (int, int, error) {
	_, err := c.c.Write([]byte("SIZE\n"))
	if err != nil {
		return 0, 0, err
	}
	scanner := bufio.NewScanner(c.c)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}
	sp := strings.Split(scanner.Text(), " ")
	w, _ := strconv.Atoi(sp[1])
	h, _ := strconv.Atoi(sp[2])

	return w, h, nil
}

func (c *Conn) SetPixel(x, y int, color *color.RGBA) error {
	s := fmt.Sprintf("PX %d %d %s\n", x, y, RGBA(color))
	_, err := fmt.Fprintf(c.c, "%s", s)
	if err != nil {
		return err
	}
	return nil
}

func (c *Conn) WriteImage(img *image.RGBA) error {
	buf := &bytes.Buffer{}
	w, h := img.Bounds().Dx(), img.Bounds().Dy()

	for yi := 0; yi < h; yi++ {
		for xi := 0; xi < w; xi++ {
			color := img.RGBAAt(xi, yi)
			fmt.Fprintf(buf, "PX %d %d %s\n", xi, yi, RGBA(&color))
		}
	}
	_, err := io.Copy(c.c, buf)
	return err
}

func main() {
	var addr = flag.String("addr", "localhost:1337", "Address of pixelflut server")
	flag.Parse()

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	c := &Conn{
		c: conn,
	}

	w, h, err := c.Size()
	if err != nil {
		log.Fatal(err)
	}

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for n := 0; n < 100000; n++ {
		clr := color.RGBA{
			R: uint8(100 * (n % 3)),
			G: uint8(100 * ((n + 1) % 3)),
			B: uint8(100 * ((n + 2) % 3)),
			A: 100,
		}

		fn := func(c color.RGBA) color.RGBA {
			return clr
		}

		img = adjust.Apply(img, fn)
		log.Printf("Setting color to %+v", clr)
		if err := c.WriteImage(img); err != nil {
			log.Fatal(err)
		}

		time.Sleep(300 * time.Millisecond)
	}
}
