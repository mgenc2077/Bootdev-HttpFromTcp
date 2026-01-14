package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	str, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	rline, err := parseRequestLine(strings.SplitN(string(str), "\r\n", 2)[0])
	if err != nil {
		return nil, err
	}

	return &Request{RequestLine: rline}, nil
}

func parseRequestLine(line string) (RequestLine, error) {
	parts := strings.SplitN(line, " ", 3)
	if len(parts) != 3 {
		return RequestLine{}, errors.New("invalid request line")
	}
	//parts[0] != "GET" && parts[0] != "POST" && parts[0] != "OPTIONS"
	if !(parts[0] == "GET" || parts[0] == "POST" || parts[0] == "OPTIONS") || !strings.Contains(parts[1], "/") || strings.TrimPrefix(parts[2], "HTTP/") != "1.1" {
		return RequestLine{}, errors.New("invalid request line")
	} else {
		return RequestLine{
			Method:        parts[0],
			RequestTarget: parts[1],
			HttpVersion:   strings.TrimPrefix(parts[2], "HTTP/"),
		}, nil
	}
}
