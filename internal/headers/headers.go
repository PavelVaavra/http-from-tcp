package headers

import (
	"strings"
	"unicode"
	"errors"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	if strings.HasPrefix(string(data), "\r\n") {
		return len("\r\n"), true, nil
	}
	headerLine := strings.Split(string(data), "\r\n")
	// \r\n wasn't found
	if len(headerLine) == 1 {
		return 0, false, nil
	}
	parts := strings.SplitN(headerLine[0], ":", 2)
	if len(parts) != 2 {
		return 0, false, errors.New("No colon found.")
	}
	key, value := parts[0], parts[1]
	if len(key) == 0 {
		return 0, false, errors.New("No key value.")
	}
	// Take the last rune (not just byte, since Unicode may be multi-byte)
	r := []rune(key)[len([]rune(key))-1]
	if unicode.IsSpace(r) {
		return 0, false, errors.New("There is a whitespace between key and colon.")
	}
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	h[key] = value
	return len(headerLine[0]) + len("\r\n"), false, nil
}