package main

import (
	"io"
	"os"
	"os/exec"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	first := os.Args[1]
	rest := os.Args[2:]
	if first[0] != 't' {
		os.Exit(1)
	}
	first = first[1:]
	args := []string{}
	if len(first) > 0 {
		args = append(args, first)
		args = append(args, rest...)
	}
	cmd := exec.Command("git", args...)

	wg := sync.WaitGroup{}
	wg.Add(2)

	sOut, err := cmd.StdoutPipe()
	if err != nil {
		os.Exit(4)
	}
	go pumpOut(sOut, os.Stdout, &wg)

	sErr, err := cmd.StderrPipe()
	if err != nil {
		os.Exit(5)
	}
	go pumpOut(sErr, os.Stderr, &wg)

	sIn, err := cmd.StdinPipe()
	if err != nil {
		os.Exit(6)
	}
	go pumpIn(sIn, os.Stdin)

	if err := cmd.Start(); err != nil {
		os.Exit(2)
	}

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		os.Exit(3)
	}

}

func pumpOut(r io.ReadCloser, w io.Writer, g *sync.WaitGroup) {
	if _, err := io.Copy(w, r); err != nil {
		os.Exit(7)
	}
	g.Done()
}

func pumpIn(w io.WriteCloser, r io.Reader) {
	if _, err := io.Copy(w, r); err != nil {
		os.Exit(8)
	}
}
