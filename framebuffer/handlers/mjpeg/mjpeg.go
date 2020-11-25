package mjpeg

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/nogoegst/pixelflut/framebuffer"
)

type Handler struct {
	fb  *framebuffer.FrameBuffer
	mux *http.ServeMux
}

func New(fb *framebuffer.FrameBuffer) *Handler {
	mux := http.NewServeMux()
	h := &Handler{
		fb:  fb,
		mux: mux,
	}
	mux.HandleFunc("/", h.staticHandler)
	mux.HandleFunc("/framebuffer.jpeg", h.mjpegHandler)
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

var staticPage = `<html>
	<head>
		<title>PixelFlut</title>
	</head>
	<body>
		<img src="framebuffer.jpeg"  />
	</body>
</html>
`

func (h *Handler) staticHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(staticPage))
}

func (h *Handler) mjpegHandler(w http.ResponseWriter, r *http.Request) {
	flusher := w.(http.Flusher)
	boundary := "pixelflut"
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary="+boundary)
	fmt.Fprintf(w, "\r\n--%s\r\n", boundary)

	var buf bytes.Buffer
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		buf.Reset()
		h.fb.DumpJPEG(&buf)

		fmt.Fprintf(w, "Content-Type: image/jpeg\r\n")
		fmt.Fprintf(w, "Content-Length: %v\r\n\r\n", buf.Len())
		_, err := io.Copy(w, &buf)
		if err != nil {
			log.Printf("write image: %v", err)
			return
		}
		fmt.Fprintf(w, "\r\n--%s\r\n", boundary)
		flusher.Flush()
		<-ticker.C
	}
}
