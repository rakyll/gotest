// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gotest is a tiny program that shells out to `go test`
// and prints the output in color.
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/fatih/color"
)

var (
	pass = color.New(color.FgGreen)
	skip = color.New(color.FgYellow)
	fail = color.New(color.FgHiRed)
)

func main() {
	setPalette()
	enableOnCI()
	os.Exit(gotest(os.Args[1:]))
}

func gotest(args []string) int {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	r, w := io.Pipe()
	defer w.Close()

	args = append([]string{"test"}, args...)
	cmd := exec.Command("go", args...)
	cmd.Stderr = w
	cmd.Stdout = w
	cmd.Env = os.Environ()

	go consume(&wg, r)

	if err := cmd.Run(); err != nil {
		if ws, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
			return ws.ExitStatus()
		}
		return 1
	}
	return 0
}

func consume(wg *sync.WaitGroup, r io.Reader) {
	defer wg.Done()
	reader := bufio.NewReader(r)
	for {
		l, _, err := reader.ReadLine()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Print(err)
			return
		}
		parse(string(l))
	}
}

var c *color.Color

func parse(line string) {
	trimmed := strings.TrimSpace(line)

	switch {
	case strings.HasPrefix(trimmed, "=== RUN"):
		fallthrough
	case strings.HasPrefix(trimmed, "?"):
		c = nil

	// succeeded
	case strings.HasPrefix(trimmed, "--- PASS"):
		fallthrough
	case strings.HasPrefix(trimmed, "ok"):
		fallthrough
	case strings.HasPrefix(trimmed, "PASS"):
		c = pass

	// skipped
	case strings.HasPrefix(trimmed, "--- SKIP"):
		c = skip

	// failed
	case strings.HasPrefix(trimmed, "--- FAIL"):
		fallthrough
	case strings.HasPrefix(trimmed, "FAIL"):
		c = fail
	}

	if c == nil {
		fmt.Printf("%s\n", line)
		return
	}
	c.Printf("%s\n", line)
}

func enableOnCI() {
	ci := strings.ToLower(os.Getenv("CI"))
	switch ci {
	case "travis":
		fallthrough
	case "appveyor":
		fallthrough
	case "gitlab_ci":
		fallthrough
	case "circleci":
		color.NoColor = false
	}
}

func setPalette() {
	envArray := [3]string{"GOTEST_SKIP_COLOR", "GOTEST_FAIL_COLOR", "GOTEST_PASS_COLOR"}

	for _, e := range envArray {
		v := os.Getenv(e)
		if v == "" {
			continue
		}
		if c, ok := colors[v]; ok {
			switch e {
			case "GOTEST_FAIL_COLOR":
				fail = color.New(c)
			case "GOTEST_PASS_COLOR":
				pass = color.New(c)
			case "GOTEST_SKIP_COLOR":
				skip = color.New(c)
			}
		}
	}
}

var colors = map[string]color.Attribute{
	"black":     color.FgBlack,
	"hiblack":   color.FgHiBlack,
	"red":       color.FgRed,
	"hired":     color.FgHiRed,
	"green":     color.FgGreen,
	"higreen":   color.FgHiGreen,
	"yellow":    color.FgYellow,
	"hiyellow":  color.FgHiYellow,
	"blue":      color.FgBlue,
	"hiblue":    color.FgHiBlue,
	"magenta":   color.FgMagenta,
	"himagenta": color.FgHiMagenta,
	"cyan":      color.FgCyan,
	"hicyan":    color.FgHiCyan,
	"white":     color.FgWhite,
	"hiwhite":   color.FgHiWhite,
}
