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
- `green` for succeeding test cases
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

For a graphical presentation of the colors please see the documentation for the Go package [color](https://pkg.go.dev/mod/github.com/fatih/color), which is the implementation used for by `gotest`.

You can configure the color of the output by setting the environment variable `GOTEST_PALETTE`:

```bash
$ GOTEST_PALETTE="magenta,white" gotest -v github.com/rakyll/hey
```

The output will have `magenta` for failed test cases and `white` for succeeding test cases.

You can specify the color for skipped tests also:

```bash
$ GOTEST_PALETTE="magenta,white,hiyellow" gotest -v github.com/rakyll/hey
```

The order of the colors for are:

1 failed test cases
1 succeeding test cases
1 skipped test cases

If you only want to overwrite the colors for indicating succeeding or skipped test cases, you have to specify the defaults (`hired` being the default for failing test cases):

```bash
$ GOTEST_PALETTE="hired,white,hiyellow" gotest -v github.com/rakyll/hey
```
