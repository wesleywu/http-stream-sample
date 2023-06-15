package streamio

import (
	"bufio"
	"bytes"
	"io"
)

// NewScanner returns a new Scanner to read from r.
// The split function is splitByTwoLineEnds.
func NewScanner(r io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	scanner.Split(splitByTwoLineEnds)
	return scanner
}

// ScanLines is a split function for a Scanner that returns each block of
// text, stripped of any trailing of two consecutive end-of-line markers.
// The returned block may be empty.
//
// The end-of-line marker should be a newline character '\n' without carriage return '\r'
// The last non-empty line of input will be returned even if it has no
// newline.
func splitByTwoLineEnds(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte("\n\n")); i >= 0 {
		// We have a full newline-terminated line.
		return i + 2, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
