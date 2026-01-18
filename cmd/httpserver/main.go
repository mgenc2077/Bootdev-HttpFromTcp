package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"mgenc2077/httpfromtcp/internal/headers"
	"mgenc2077/httpfromtcp/internal/request"
	"mgenc2077/httpfromtcp/internal/response"
	"mgenc2077/httpfromtcp/internal/server"
)

const port = 42069

func handler(w *response.Writer, req *request.Request) {
	body := []byte("Your problem is not my problem\n")
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		w.WriteStatusLine(response.StatusBadRequest)
		body = []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>\n`)
		w.WriteHeaders(response.ModifyDefaultHeaders(headers.Headers{"content-type": "text/html", "content-length": fmt.Sprintf("%d", len(body))}))
		w.Write(body)
	case "/myproblem":
		w.WriteStatusLine(response.StatusInternalServerError)
		body = []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
		w.WriteHeaders(response.ModifyDefaultHeaders(headers.Headers{"content-type": "text/html", "content-length": fmt.Sprintf("%d", len(body))}))
		w.Write(body)
	case "/httpbin/stream/100":
		w.WriteStatusLine(response.StatusOK)
		w.WriteHeaders(response.ModifyDefaultHeaders(headers.Headers{"content-type": "application/json", "Transfer-Encoding": "chunked"}))

		resp, _ := http.Get("https://httpbin.org/stream/100")
		for {
			buf := make([]byte, 1024)
			n, err := resp.Body.Read(buf)
			if n > 0 {
				w.WriteChunkedBody(buf[:n])
			}
			if err != nil {
				break
			}
		}
		w.WriteChunkedBodyDone()

	case "/":
		w.WriteStatusLine(response.StatusOK)
		body = []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
		w.WriteHeaders(response.ModifyDefaultHeaders(headers.Headers{"content-type": "text/html", "content-length": fmt.Sprintf("%d", len(body))}))
		w.Write(body)
	}
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
