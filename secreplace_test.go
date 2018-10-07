package secreplace

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	const (
		open  = "(_"
		close = "_)"
	)

	tests := []struct {
		s     string
		start int
		end   int
		ok    bool
		err   error
	}{
		{
			s:     "(_foo_)",
			start: 0,
			end:   7,
			ok:    true,
		},
		{
			s:     "(_bar (_foo_)_)",
			start: 6,
			end:   13,
			ok:    true,
		},
		{
			s:     "(_foo",
			start: -1,
			end:   -1,
			ok:    false,
			err:   ErrNoMatchingClose,
		},
		{
			s:     "(_(_foo",
			start: -1,
			end:   -1,
			ok:    false,
			err:   ErrNoMatchingClose,
		},
		{
			s:     "foo_)",
			start: -1,
			end:   -1,
			ok:    false,
			err:   ErrNoMatchingOpen,
		},
		{
			s:     "foo_)_)",
			start: -1,
			end:   -1,
			ok:    false,
			err:   ErrNoMatchingOpen,
		},
		{
			s:     "foo_) (_bar_)",
			start: -1,
			end:   -1,
			ok:    false,
			err:   ErrNoMatchingOpen,
		},
	}

	for _, test := range tests {
		t.Run(test.s, func(t *testing.T) {
			start, end, ok, err := Find(test.s, open, close)
			assert.Equal(t, test.start, start)
			assert.Equal(t, test.end, end)
			assert.Equal(t, test.ok, ok)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestReplaceOne(t *testing.T) {
	const (
		open  = "(_"
		close = "_)"
	)

	identity := func(s string) (string, error) {
		return s, nil
	}

	tests := []struct {
		s         string
		f         func(string) (string, error)
		expected  string
		unchanged bool
		err       error
	}{
		{
			s:        "(_foo_)",
			expected: "foo",
		},
		{
			s:         "foo",
			expected:  "foo",
			unchanged: true,
		},
		{
			s:         "(_foo",
			expected:  "(_foo",
			unchanged: true,
			err:       ErrNoMatchingClose,
		},
		{
			s:         "foo_)",
			expected:  "foo_)",
			unchanged: true,
			err:       ErrNoMatchingOpen,
		},
		{
			s:        "(_foo (_bar_)_)",
			expected: "(_foo bar_)",
		},
		{
			s:        "(_foo (_bar_) (_baz_)_)",
			expected: "(_foo bar (_baz_)_)",
		},
	}

	for _, test := range tests {
		f := identity
		if test.f != nil {
			f = test.f
		}

		t.Run(test.s, func(t *testing.T) {
			got, changed, err := ReplaceOne(test.s, open, close, f)
			assert.Equal(t, test.expected, got)
			assert.Equal(t, !test.unchanged, changed)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestReplaceAll(t *testing.T) {
	const (
		open  = "(_"
		close = "_)"
	)

	identity := func(s string) (string, error) {
		return s, nil
	}

	tests := []struct {
		s         string
		f         func(string) (string, error)
		expected  string
		unchanged bool
		err       error
	}{
		{
			s:        "(_foo_)",
			expected: "foo",
		},
		{
			s:         "foo",
			expected:  "foo",
			unchanged: true,
		},
		{
			s:         "(_foo",
			expected:  "(_foo",
			unchanged: true,
			err:       ErrNoMatchingClose,
		},
		{
			s:         "foo_)",
			expected:  "foo_)",
			unchanged: true,
			err:       ErrNoMatchingOpen,
		},
		{
			s:        "(_foo (_bar_)_)",
			expected: "foo bar",
		},
		{
			s:        "(_foo (_bar_) (_baz_)_)",
			expected: "foo bar baz",
		},
		{
			s: "Hi, my name is (_(_A_)-(_B_)_)!",
			f: func(s string) (string, error) {
				return "COOL-" + s, nil
			},
			expected: "Hi, my name is COOL-COOL-A-COOL-B!",
		},
	}

	for _, test := range tests {
		t.Run(test.s, func(t *testing.T) {
			f := identity
			if test.f != nil {
				f = test.f
			}

			got, changed, err := ReplaceAll(test.s, open, close, f)
			assert.Equal(t, test.expected, got)
			assert.Equal(t, !test.unchanged, changed)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestReplaceFuncErr(t *testing.T) {
	const (
		open  = "(_"
		close = "_)"
		test  = open + "foo" + close
	)

	errTest := errors.New("test error")

	f := func(s string) (string, error) {
		return "", errTest
	}

	funcs := []struct {
		name string
		f    func(s string, open, close string, f func(string) (string, error)) (string, bool, error)
	}{
		{"ReplaceOne", ReplaceOne},
		{"ReplaceAll", ReplaceAll},
	}

	for _, tf := range funcs {
		t.Run(tf.name, func(t *testing.T) {
			s, changed, err := tf.f(test, open, close, f)
			assert.Equal(t, test, s)
			assert.Equal(t, false, changed)
			assert.Equal(t, errTest, err)
		})
	}
}
