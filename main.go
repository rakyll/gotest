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
	success = color.New(color.FgGreen)
	fail    = color.New(color.FgHiRed)
)

const paletteEnv = "GOTEST_PALETTE"

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

	// success
	case strings.HasPrefix(trimmed, "--- PASS"):
		fallthrough
	case strings.HasPrefix(trimmed, "ok"):
		fallthrough
	case strings.HasPrefix(trimmed, "PASS"):
		c = success

	// failure
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
	v := os.Getenv(paletteEnv)
	if v == "" {
		return
	}
	vals := strings.Split(v, ",")
	if len(vals) != 2 {
		return
	}
	if c, ok := colors[vals[0]]; ok {
		fail = color.New(c)
	}
	if c, ok := colors[vals[1]]; ok {
		success = color.New(c)
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
