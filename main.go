package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

var green = color.New(color.FgGreen)
var red = color.New(color.FgHiRed)

func main() {
	gotest(os.Args[1:])
}

func gotest(args []string) {
	r, w := io.Pipe()

	args = append([]string{"test"}, args...)
	cmd := exec.Command("go", args...)
	cmd.Stderr = w
	cmd.Stdout = w
	cmd.Env = os.Environ()

	go consume(r)

	cmd.Run()
}

func consume(r io.Reader) {
	reader := bufio.NewReader(r)
	for {
		l, _, err := reader.ReadLine()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		parse(r, string(l))
	}
}

var c *color.Color

func parse(r io.Reader, line string) {
	trimmed := strings.TrimSpace(line)

	switch {
	case strings.HasPrefix(trimmed, "=== RUN"):
		c = nil
	case strings.HasPrefix(trimmed, "--- PASS"):
		fallthrough
	case strings.HasPrefix(trimmed, "ok"):
		fallthrough
	case strings.HasPrefix(trimmed, "PASS"):
		c = green
	case strings.HasPrefix(trimmed, "--- FAIL"):
		fallthrough
	case strings.HasPrefix(trimmed, "FAIL"):
		c = red
	}

	if c == nil {
		fmt.Printf("%s\n", line)
	} else {
		c.Printf("%s\n", line)
	}
}
