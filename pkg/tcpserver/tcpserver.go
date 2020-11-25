package tcpserver

import (
	"fmt"
	"log"
	"net"
)

type Handler interface {
	Handle(conn net.Conn) error
}

type Server struct {
	addr    string
	handler Handler
}

func New(addr string, handler Handler) *Server {
	s := &Server{
		addr:    addr,
		handler: handler,
	}
	return s
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("create listener: %w", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("unable to accept connection: %w", err)
		}
		go func() {
			defer conn.Close()
			err := s.handler.Handle(conn)
			if err != nil {
				log.Printf("error handling a connection: %v", err)
			}
		}()
	}
}
