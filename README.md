# secreplace

[![GoDoc Reference](https://godoc.org/github.com/jakebailey/secreplace?status.svg)](http://godoc.org/github.com/jakebailey/secreplace) [![Go Report Card](https://goreportcard.com/badge/github.com/jakebailey/secreplace)](https://goreportcard.com/report/github.com/jakebailey/secreplace) [![Build Status](https://travis-ci.com/jakebailey/secreplace.svg?branch=master)](https://travis-ci.com/jakebailey/secreplace) [![Coverage Status](https://coveralls.io/repos/github/jakebailey/secreplace/badge.svg?branch=master)](https://coveralls.io/github/jakebailey/secreplace?branch=master)


`secreplace` is a library which aids in the string replacement of "sections".
A section is defined as text enclosed by open and close terminators.

For example, take the following string:

```
Hi, my name is (_NAME_)!
```

`(_` and `_)` are the open and close terminators, respectively. We'd like
to replace `(_NAME_)` with something else. Let's pick a function and run
`ReplaceOne`:

```go
replacer := func(s string) (string, error) {
    return "COOL-"+s, nil
}

out, changed, err := ReplaceOne("Hi, my name is (_NAME_)!", "(_", "_)", replacer)
// out     == "Hi, my name is COOL-NAME!"
// changed == true
// err     == nil
```

The `ReplaceAll` function will repeatedly call `ReplaceOne` until no more
replacements can occur. This allows for nesting, for example:


```go
replacer := func(s string) (string, error) {
    return "COOL"+s, nil
}

out, changed, err := ReplaceAll("Hi, my name is (_(_A_)-(_B_)_)!", "(_", "_)", replacer)
// out     == "Hi, my name is COOL-COOL-A-COOL-B!"
// changed == true
// err     == nil
```

Any errors returned by the replacement function will be returned by
`ReplaceOne` and `ReplaceAll`. Any unmatched terminators will be caught and
also cause errors to be returned.
