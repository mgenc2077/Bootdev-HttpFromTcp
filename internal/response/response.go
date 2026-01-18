package response

import (
	"fmt"
	"io"
	"mgenc2077/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

type Writer struct {
	io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{Writer: w}
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.Write(p)
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	switch statusCode {
	case StatusOK:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		return err
	case StatusBadRequest:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		return err
	case StatusInternalServerError:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		return err
	default:
		_, err := w.Write([]byte("\r\n"))
		return err
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

// Modifies the default headers with the provided headers
func ModifyDefaultHeaders(h headers.Headers) headers.Headers {
	if h["Connection"] == "" {
		h.Set("Connection", "close")
	}
	if h["content-type"] == "" {
		h.Set("Content-Type", "text/plain")
	}
	return h
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	for key, value := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	size := len(p)
	// Write the size of the chunk in hexadecimal followed by \r\n
	_, err := fmt.Fprintf(w, "%x\r\n", size)
	if err != nil {
		return 0, err
	}
	// Write the actual chunk data
	n, err := w.Write(p)
	if err != nil {
		return n, err
	}
	// Write the trailing \r\n after the chunk data
	_, err = w.Write([]byte("\r\n"))
	if err != nil {
		return n, err
	}
	return n, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	_, err := w.Write([]byte("0\r\n\r\n"))
	if err != nil {
		return 0, err
	}
	return 0, nil
}
