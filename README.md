# gotest

Like `go test` but with colors.

## Installation

```bash
$ go get -u github.com/rakyll/gotest
```

# Usage

Accepts all the arguments and flags `go test` works with.

Example:

```bash
$ gotest -v github.com/rakyll/hey
```

![go test output](https://i.imgur.com/udjWuZx.gif)

The default output colors for `gotest` are:

- `hired` for failed test cases
- `green` for passing test cases
- `yellow` for skipped test cases

`gotest` comes with many colors. All available colors are:

- `black`
- `hiblack`
- `red`
- `hired`
- `green`
- `higreen`
- `yellow`
- `hiyellow`
- `blue`
- `hiblue`
- `magenta`
- `himagenta`
- `cyan`
- `hicyan`
- `white`
- `hiwhite`

For a graphical presentation of the colors please see the documentation for the Go package [fatih/color](https://pkg.go.dev/mod/github.com/fatih/color).

You can configure the color of the output by setting the environment variable `GOTEST_PALETTE`:

```bash
$ GOTEST_PALETTE="magenta,white,hiyellow" gotest -v github.com/rakyll/hey
```

The output will have `magenta` for failed test cases, `white` for passing test cases, and `hiyellow` for skipped test cases.

The order of the colors for are:

1 failed test cases
1 passing test cases
1 skipped test cases

If you only want to overwrite the colors for indicating passing or skipped test cases, you can leave the spot empty.

This example demonstrates:

- `hired` for failed test cases, using the default
- `white` for passing test cases, overwriting the default
- `hiyellow` for skipped test cases, overwriting the default

```bash
$ GOTEST_PALETTE=",white,hiyellow" gotest -v github.com/rakyll/hey
```
