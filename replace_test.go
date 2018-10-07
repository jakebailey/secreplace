package secreplace

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	const (
		open  = "_(_"
		close = "_)_"
	)

	tests := []struct {
		s     string
		start int
		end   int
		ok    bool
		err   error
	}{
		{
			s:     "_(_foo_)_",
			start: 0,
			end:   9,
			ok:    true,
		},
		{
			s:     "_(_bar _(_foo_)__)_",
			start: 7,
			end:   16,
			ok:    true,
		},
		{
			s:     "_(_foo",
			start: -1,
			end:   -1,
			ok:    false,
			err:   ErrNoMatchingClose,
		},
		{
			s:     "_(__(_foo",
			start: -1,
			end:   -1,
			ok:    false,
			err:   ErrNoMatchingClose,
		},
		{
			s:     "foo_)_",
			start: -1,
			end:   -1,
			ok:    false,
			err:   ErrNoMatchingOpen,
		},
		{
			s:     "foo_)__)_",
			start: -1,
			end:   -1,
			ok:    false,
			err:   ErrNoMatchingOpen,
		},
		{
			s:     "foo_)_ _(_bar_)_",
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
		open  = "_(_"
		close = "_)_"
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
			s:        "_(_foo_)_",
			expected: "foo",
		},
		{
			s:         "foo",
			expected:  "foo",
			unchanged: true,
		},
		{
			s:         "_(_foo",
			unchanged: true,
			err:       ErrNoMatchingClose,
		},
		{
			s:         "foo_)_",
			unchanged: true,
			err:       ErrNoMatchingOpen,
		},
		{
			s:        "_(_foo _(_bar_)__)_",
			expected: "_(_foo bar_)_",
		},
		{
			s:        "_(_foo _(_bar_)_ _(_baz_)__)_",
			expected: "_(_foo bar _(_baz_)__)_",
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
		open  = "_(_"
		close = "_)_"
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
			s:        "_(_foo_)_",
			expected: "foo",
		},
		{
			s:         "foo",
			expected:  "foo",
			unchanged: true,
		},
		{
			s:         "_(_foo",
			unchanged: true,
			err:       ErrNoMatchingClose,
		},
		{
			s:         "foo_)_",
			unchanged: true,
			err:       ErrNoMatchingOpen,
		},
		{
			s:        "_(_foo _(_bar_)__)_",
			expected: "foo bar",
		},
		{
			s:        "_(_foo _(_bar_)_ _(_baz_)__)_",
			expected: "foo bar baz",
		},
		{
			s: "Hi, my name is _(__(_A_)_-_(_B_)__)_!",
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
		open  = "_(_"
		close = "_)_"
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
			assert.Equal(t, "", s)
			assert.Equal(t, false, changed)
			assert.Equal(t, errTest, err)
		})
	}
}
