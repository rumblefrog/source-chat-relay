package exrouter

import (
	"bytes"
	"encoding/csv"
	"strings"
)

// separator is the separator character for splitting arguments
const separator = ' '

// Args is a helper type for dealing with command arguments
type Args []string

// Get returns the argument at index n
func (a Args) Get(n int) string {
	if n >= 0 && n < len(a) {
		return a[n]
	}
	return ""
}

// After returns all arguments after index n
func (a Args) After(n int) string {
	if n >= 0 && n < len(a) {
		return strings.Join(a[n:], string(separator))
	}
	return ""
}

// ParseArgs parses command arguments
func ParseArgs(content string) Args {
	cv := csv.NewReader(bytes.NewBufferString(content))
	cv.Comma = separator
	fields, err := cv.Read()
	if err != nil {
		return strings.Split(content, string(separator))
	}
	return fields
}
