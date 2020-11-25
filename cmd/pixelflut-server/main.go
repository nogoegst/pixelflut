package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/nogoegst/pixelflut/framebuffer"
	"github.com/nogoegst/pixelflut/framebuffer/handlers/mjpeg"
	"github.com/nogoegst/pixelflut/framebuffer/handlers/rgb"
	pfhandler "github.com/nogoegst/pixelflut/handler"
	"github.com/nogoegst/pixelflut/pkg/tcpserver"
	"golang.org/x/sync/errgroup"
)

func main() {
	var width = flag.Int("w", 800, "Framebuffer width")
	var height = flag.Int("h", 600, "Framebuffer height")
	var pixelflutPort = flag.String("pixelflut-port", "1337", "PixelFlut protocol port")
	var framebufferPort = flag.String("framebuffer-port", "1338", "Framebuffer raw port")
	var httpPort = flag.String("http-port", "80", "Framebuffer MJPEG over HTTP port")
	flag.Parse()

	fb := framebuffer.New(*width, *height)

	pixelflutHandler := pfhandler.New(fb)
	pixelflutServer := tcpserver.New(":"+*pixelflutPort, pixelflutHandler)

	framebufferHandler := rgb.New(fb)
	framebufferServer := tcpserver.New(":"+*framebufferPort, framebufferHandler)

	mjpegHandler := mjpeg.New(fb)

	var eg errgroup.Group
	eg.Go(pixelflutServer.ListenAndServe)
	eg.Go(framebufferServer.ListenAndServe)
	eg.Go(func() error {
		return http.ListenAndServe(":"+*httpPort, mjpegHandler)
	})

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}
