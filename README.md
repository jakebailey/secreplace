# secreplace

`secreplace` is a library which aids in the string replacement of "sections".
A section is defined as text enclosed by open and close terminators.

For example, take the following string:

```
Hi, my name is _(_NAME_)_!
```

`_(_` and `_)_` are the open and close terminators, respectively. We'd like
to replace `_(_NAME_)_` with something else. Let's pick a function and run
`ReplaceOne`:

```go
identity := func(s string) (string, error) {
    return "COOL-"+s, nil
}

out, changed, err := ReplaceOne("Hi, my name is _(_NAME_)_!", "_(_", "_)_", identity)
// out     == "Hi, my name is COOL-NAME!"
// changed == true
// err     == nil
```

The `ReplaceAll` function will repeatedly call `ReplaceOne` until no more
replacements can occur. This allows for nesting, for example:


```go
identity := func(s string) (string, error) {
    return "COOL"+s, nil
}

out, changed, err := ReplaceAll("Hi, my name is _(__(_A_)_-_(_B_)__)_!", "_(_", "_)_", identity)
// out     == "Hi, my name is COOL-COOL-A-COOL-B!"
// changed == true
// err     == nil
```

Any errors returned by the replacement function will be returned by
`ReplaceOne` and `ReplaceAll`. Any unmatched terminators will be caught and
also cause errors to be returned.
