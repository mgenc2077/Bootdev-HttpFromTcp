package server

import (
	"fmt"
	"io"
	"mgenc2077/httpfromtcp/internal/request"
	"mgenc2077/httpfromtcp/internal/response"
	"net"
	"sync/atomic"
)

type Server struct {
	Listener   net.Listener
	Connection net.Conn
	Status     atomic.Bool
	Handler    Handler
}

type Handler func(w *response.Writer, req *request.Request)

type HandlerError struct {
	StatusCode int
	Err        error
}

// Creates a net.Listener and returns a new Server instance
func Serve(port int, handler Handler) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		Listener: l,
		Handler:  handler,
	}
	s.Status.Store(true)
	go s.listen()
	return s, nil
}

// Closes the listener and the server
func (s *Server) Close() error {
	if err := s.Connection.Close(); err != nil {
		return err
	}
	if err := s.Listener.Close(); err != nil {
		return err
	}
	return nil
}

// Uses a loop to .Accept new connections as they come in, and handles each one in a new goroutine
func (s *Server) listen() error {
	for {
		// Wait for a connection.
		conn, err := s.Listener.Accept()
		if err != nil {
			return err
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			if s.Status.Load() == true {
				s.handle(c)
			}
		}(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	s.Connection = conn
	req, err := request.RequestFromReader(conn)
	if err != nil && err != io.EOF {
		writer := response.NewWriter(conn)
		writer.WriteStatusLine(response.StatusBadRequest)
		writer.WriteHeaders(response.GetDefaultHeaders(0))
		conn.Close()
		return
	}

	writer := response.NewWriter(conn)
	s.Handler(writer, req)
	conn.Close()
}
