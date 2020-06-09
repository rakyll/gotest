// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gotest is a tiny program that shells out to `go test`
// and prints the output in color.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/fatih/color"
)

var (
	pass        = color.FgGreen
	skip        = color.FgYellow
	fail        = color.FgHiRed
	gotestFlags = []string{
		"-skipnotest",
	}
	parseFlags = make([]string,0)
	args = make([]string,0)
	ignore *bool
)

const paletteEnv = "GOTEST_PALETTE"

func main() {
	setPalette()
	enableOnCI()

	flagSet := flag.NewFlagSet("gotestFlags",flag.ContinueOnError)
	ignore=flagSet.Bool("skipnotest", false, "skip packages with no test files")

	// separate program specific flags from go test flags
	for _,arg:=range os.Args[1:] {
		for _,flg := range gotestFlags {
			// if the argument is in gotest flags the add it to parseFlags
			if strings.Contains(arg,flg){
				parseFlags=append(parseFlags,arg)
			} else {
				args = append(args,arg)
			}
		}
	}

	flagSet.Parse(parseFlags)

	os.Exit(gotest(args))

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

	sigc := make(chan os.Signal)
	done := make(chan struct{})
	defer func() {
		done <- struct{}{}
	}()
	signal.Notify(sigc)

	go func() {
		for {
			select {
			case sig := <-sigc:
				cmd.Process.Signal(sig)
			case <-done:
				return
			}
		}
	}()

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

func parse(line string) {
	trimmed := strings.TrimSpace(line)
	defer color.Unset()

	var c color.Attribute
	switch {
	case strings.HasPrefix(trimmed, "=== RUN"):
		fallthrough
	case strings.HasPrefix(trimmed, "?"):
		if *ignore {
			return
		}
		color.Unset()
		// passed
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

	color.Set(c)
	fmt.Printf("%s\n", line)
}

func enableOnCI() {
	ci := strings.ToLower(os.Getenv("CI"))
	switch ci {
	case "true":
		fallthrough
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
		fail = c
	}
	if c, ok := colors[vals[1]]; ok {
		pass = c
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
