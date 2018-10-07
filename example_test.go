package secreplace_test

import (
	"fmt"

	"github.com/jakebailey/secreplace"
)

func ExampleReplaceOne() {
	const (
		open  = "(_"
		close = "_)"
	)

	replacer := func(s string) (string, error) {
		return s, nil
	}

	out, changed, err := secreplace.ReplaceOne("Hello, (_NAME_)!", open, close, replacer)
	fmt.Println(out)
	fmt.Println(changed)
	fmt.Println(err)
	// Output:
	// Hello, NAME!
	// true
	// <nil>
}

func ExampleReplaceAll() {
	const (
		open  = "(_"
		close = "_)"
	)

	replacer := func(s string) (string, error) {
		return s, nil
	}

	out, changed, err := secreplace.ReplaceAll("(_foo (_bar_) (_baz (_qux_)_)_)", open, close, replacer)
	fmt.Println(out)
	fmt.Println(changed)
	fmt.Println(err)
	// Output:
	// foo bar baz qux
	// true
	// <nil>
}
