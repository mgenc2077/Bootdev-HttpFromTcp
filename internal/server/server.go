package server

import (
	"fmt"
	"mgenc2077/httpfromtcp/internal/response"
	"net"
	"sync/atomic"
)

type Server struct {
	Listener   net.Listener
	Connection net.Conn
	Status     atomic.Bool
}

// Creates a net.Listener and returns a new Server instance
func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		Listener: l,
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

// Handles a single connection by writing the response and then closing the connection
func (s *Server) handle(conn net.Conn) {
	s.Connection = conn
	response.WriteStatusLine(conn, response.StatusOK)
	response.GetDefaultHeaders(0)
	response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	conn.Close()
}
