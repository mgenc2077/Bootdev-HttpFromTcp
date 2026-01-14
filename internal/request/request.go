package request

import (
	"errors"
	"io"
	"strings"
)

const (
	initialized = iota
	done        = iota
)

type Request struct {
	RequestLine RequestLine
	State       int
}

func (r *Request) parse(data []byte) (int, error) {
	if r.State == initialized {
		parsedLine, size, err := parseRequestLine(string(data))
		if size == 0 && err == nil {
			return 0, nil
		} else if err != nil {
			return 0, err
		} else {
			r.RequestLine = parsedLine
			r.State = done
			return size, nil
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
	req := &Request{State: initialized}
	for {
		if req.State == done {
			break
		}
		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		n, err := reader.Read(buf[readToIndex:])
		if err == io.EOF {
			req.State = done
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
