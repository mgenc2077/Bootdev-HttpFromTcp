package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Set(key, value string) {
	h[key] = value
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	strData := string(data)
	firstNewline := strings.Index(strData, "\r\n")
	if firstNewline == -1 {
		return 0, false, nil
	}
	if firstNewline == 0 {
		return 2, true, nil
	}

	line := strData[:firstNewline]
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return 0, false, errors.New("invalid header format")
	}

	key := parts[0]
	// Check for space before colon
	if strings.HasSuffix(key, " ") {
		return 0, false, errors.New("invalid header format: space before colon")
	}
	if !isFieldName(strings.TrimSpace(key)) {
		return 0, false, errors.New("invalid header field name")
	}

	val := parts[1]
	h[strings.ToLower(strings.TrimSpace(key))] = strings.TrimSpace(val)

	return firstNewline + 2, false, nil
}

func isFieldName(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ { // bytes; RFC token is ASCII
		c := s[i]

		// ALPHA
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			continue
		}
		// DIGIT
		if c >= '0' && c <= '9' {
			continue
		}
		// tchar specials
		switch c {
		case '!', '#', '$', '%', '&', '\'', '*',
			'+', '-', '.', '^', '_', '`', '|', '~':
			continue
		default:
			return false
		}
	}
	return true
}
