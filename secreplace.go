// Package secreplace aids in replacing sections of text which are surrounded
// by known beginning and end terminators. Sections can be nested, and
// secreplace can replace all until no changes can be made.
package secreplace

import (
	"errors"
	"strings"
)

var (
	// ErrNoMatchingOpen is returned when there is no matching open for a close.
	ErrNoMatchingOpen = errors.New("secreplace: no matching open")
	// ErrNoMatchingClose is returned when there is no matching close for an open.
	ErrNoMatchingClose = errors.New("secreplace: no matching close")
)

// Find searches for the first, most interior section of input text surrounded
// by open and close. It returns the range of text containing the section in
// the form [start, end) including the open/close strings, a boolean that is
// true when a match is found, and an error.
func Find(s string, open, close string) (start, end int, ok bool, err error) {
	closeIdx := strings.Index(s, close)
	if closeIdx == -1 {
		openIdx := strings.Index(s, open)
		if openIdx != -1 {
			return -1, -1, false, ErrNoMatchingClose
		}

		return -1, -1, false, nil
	}

	realCloseIdx := closeIdx + len(close)

	s = s[:realCloseIdx]

	openIdx := strings.LastIndex(s, open)
	if openIdx == -1 {
		return -1, -1, false, ErrNoMatchingOpen
	}

	return openIdx, realCloseIdx, true, nil
}

// ReplaceOne replaces a single match found by Find. It calls f on the text
// between the open and close strings, and returns a new string with the
// whole section replaced. Errors produced by Find and f are propogated.
func ReplaceOne(s string, open, close string, f func(string) (string, error)) (out string, changed bool, err error) {
	start, end, ok, err := Find(s, open, close)
	if err != nil {
		return "", false, err
	}
	if !ok {
		return s, false, nil
	}

	prefix := s[:start]
	middle := s[start+len(open) : end-len(close)]
	suffix := s[end:]

	replaced, err := f(middle)
	if err != nil {
		return "", false, err
	}

	return prefix + replaced + suffix, true, nil
}

// ReplaceAll calls ReplaceOne repeatedly until no more replacements can be
// made, or an error occurs.
func ReplaceAll(s string, open, close string, f func(string) (string, error)) (out string, changed bool, err error) {
	for {
		replaced, c, err := ReplaceOne(s, open, close, f)
		if err != nil {
			return "", false, err
		}
		if !c {
			break
		}

		s = replaced
		changed = true
	}

	return s, changed, nil
}
