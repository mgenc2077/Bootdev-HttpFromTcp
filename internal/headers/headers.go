package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

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

	val := parts[1]
	h[strings.TrimSpace(key)] = strings.TrimSpace(val)

	return firstNewline + 2, false, nil
}
