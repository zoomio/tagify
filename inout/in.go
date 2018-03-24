package inout

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Source types
const (
	STDIn = iota
	FS
	Web
)

// In - Input. This struct provides methods for reading strings
// and numbers from standard input, file input, URLs, and sockets.
type In struct {
	lines []string
	SourceType int
}

// NewIn initializes an input stream from STDIN, file or web page.
//
// name - the filename or web page name, reads from STDIN if name is empty.
// Panics on errors.
func NewIn(name string) In {
	in := In{}

	// STDIN
	if name == "" {
		in.SourceType = STDIn
		stat, err := os.Stdin.Stat()
		if err != nil {
			panic(fmt.Sprintf("error in reading from STDIN: %v", err))
		}
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return in
		}
		in.lines, err = linesFromReader(bufio.NewReader(os.Stdin))
		if err != nil {
			panic(fmt.Sprintf("error in reading from STDIN: %v", err))
		}
	// File system
	} else if _, err := os.Stat(name); err == nil {
		in.SourceType = FS
		f, err := os.Open(name)
		if err != nil {
			panic(fmt.Sprintf("error in opening file %s for reading: %v", name, err))
		}
		defer f.Close()
		in.lines, err = linesFromReader(bufio.NewReader(f))
		if err != nil {
			panic(fmt.Sprintf("error in reading from file %s: %v", name, err))
		}
	// HTTP
	} else {
		in.SourceType = Web
		resp, err := http.Get(name)
		if err != nil {
			panic(fmt.Sprintf("provided name=%s is not a file and not a URL: %v", name, err))
		}

		defer resp.Body.Close()
		in.lines, err = linesFromReader(resp.Body)
		if err != nil {
			panic(fmt.Sprintf("error in reading from %s: %v", name, err))
		}
	}

	return in
}

// NewInFromString ...
func NewInFromString(input string) In {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(input))
	// Set the split function for the scanning operation.
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("error in reading input: %v", err))
	}
	return In{
		lines: lines,
	}
}

// GetLines returns lines.
func (in *In) GetLines() []string {
	return in.lines
}

// ReadAllStrings provides slice of strings from input split by white space.
func (in *In) ReadAllStrings() []string {
	tokens := make([]string, 0)
	if in.lines != nil {
		for _, line := range in.lines {
			tokens = append(tokens, strings.Fields(line)...)
		}
	}
	return tokens
}

func linesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
