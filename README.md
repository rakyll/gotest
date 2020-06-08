# gotest

[![CircleCI](https://circleci.com/gh/rakyll/gotest.svg?style=svg)](https://circleci.com/gh/rakyll/gotest)

Like `go test` but with colors.

## Installation

```
$ go get -u github.com/rakyll/gotest
```

# Usage

Accepts all the arguments and flags `go test` works with.

Example:

```
$ gotest -v github.com/rakyll/hey
```
![go test output](https://i.imgur.com/udjWuZx.gif)

gotest comes with many colors! Configure the color of the output by setting the following environment variables:

- `GOTEST_FAIL`
- `GOTEST_PASS`
- `GOTEST_SKIP`

Alternatively you can use a single environment supporting a list of colors, in the order: fail and pass.

```
$ GOTEST_PALETTE="magenta,white"
```

The output will have magenta for failed cases, white for success.
Available colors: black, hiblack, red, hired, green, higreen, yellow, hiyellow, blue, hiblue, magenta, himagenta, cyan, hicyan, white, hiwhite.

Do note that the individually set environment variables take precedence over the palette variable

For the setting:

```
$ GOTEST_PASS="hiblue" GOTEST_PALETTE="magenta,white"
```

The output will have magenta for failed cases, hiblue for success.
