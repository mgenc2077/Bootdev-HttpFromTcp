package request

import (
	"errors"
	"io"
	"strings"

	"mgenc2077/httpfromtcp/internal/headers"
)

const (
	initialized                = iota
	done                       = iota
	requestStateParsingHeaders = iota
	requestStateDone           = iota
)

type Request struct {
	RequestLine RequestLine
	State       int
	Headers     headers.Headers
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.State != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			totalBytesParsed += n
			return totalBytesParsed, err
		}
		if n == 0 {
			// Not enough data to continue parsing, return what we have so far
			return totalBytesParsed, nil
		}
		totalBytesParsed += n
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	if r.State == initialized {
		parsedLine, size, err := parseRequestLine(string(data))
		if size == 0 && err == nil {
			return 0, nil
		} else if err != nil {
			return 0, err
		} else {
			r.RequestLine = parsedLine
			r.State = requestStateParsingHeaders
			return size, nil
		}
	} else if r.State == requestStateParsingHeaders {
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		} else if done {
			r.State = requestStateDone
			return n, nil
		} else {
			return n, nil
		}
	} else if r.State == done {
		return 0, errors.New("error: trying to read data in a done state")
	} else {
		return 0, errors.New("error: unknown state")
	}
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	req := &Request{State: initialized, Headers: headers.NewHeaders()}
	for {
		if req.State == requestStateDone {
			break
		}
		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		n, err := reader.Read(buf[readToIndex:])
		if err == io.EOF {
			break
		}
		readToIndex += n
		parsedBytes, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[parsedBytes:readToIndex])
		readToIndex -= parsedBytes
	}

	if req.State != requestStateDone {
		return nil, io.ErrUnexpectedEOF
	}

	return req, nil
}

func parseRequestLine(line string) (RequestLine, int, error) {
	firstNewline := strings.Index(line, "\r\n")
	if firstNewline == -1 {
		return RequestLine{}, 0, nil
	}
	// Slice the line to include only the request line, excluding \r\n and subsequent data
	requestLineContent := line[:firstNewline]
	parts := strings.SplitN(requestLineContent, " ", 3)
	if len(parts) != 3 {
		return RequestLine{}, 0, errors.New("invalid request line")
	}

	if !(parts[0] == "GET" || parts[0] == "POST" || parts[0] == "OPTIONS") || !strings.Contains(parts[1], "/") || strings.TrimPrefix(parts[2], "HTTP/") != "1.1" {
		return RequestLine{}, 0, errors.New("invalid request line")
	} else {
		return RequestLine{
			Method:        parts[0],
			RequestTarget: parts[1],
			HttpVersion:   strings.TrimPrefix(parts[2], "HTTP/"),
		}, firstNewline + 2, nil
	}
}
