package main

import (
	"errors"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"mgenc2077/httpfromtcp/internal/request"
	"mgenc2077/httpfromtcp/internal/server"
)

const port = 42069

func handler(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{
			StatusCode: 400,
			Err:        errors.New("Your problem is not my problem\n"),
		}
	case "/myproblem":
		return &server.HandlerError{
			StatusCode: 500,
			Err:        errors.New("Woopsie, my bad\n"),
		}
	case "/use-nvim":
		w.Write([]byte("All good, frfr\n"))
		return nil
	}
	return nil
}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
