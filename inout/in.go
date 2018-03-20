package inout

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// In - Input. This struct provides methods for reading strings
// and numbers from standard input, file input, URLs, and sockets.
type In struct {
	lines []string
}

// NewIn initializes an input stream from STDIN, file or web page.
//
// name - the filename or web page name, reads from STDIN if name is empty.
// Panics on errors.
func NewIn(name string) In {
	var lines []string

	// STDIN
	if name == "" {
		stat, err := os.Stdin.Stat()
		if err != nil {
			panic(fmt.Sprintf("error in reading from STDIN: %v", err))
		}
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return In{}
		}
		lines, err = linesFromReader(bufio.NewReader(os.Stdin))
		if err != nil {
			panic(fmt.Sprintf("error in reading from STDIN: %v", err))
		}	
	// File system	
	} else if _, err := os.Stat(name); err == nil {
		f, err := os.Open(name)
		if err != nil {
			panic(fmt.Sprintf("error in opening file %s for reading: %v", name, err))
		}
		defer f.Close()
		lines, err = linesFromReader(bufio.NewReader(f))
		if err != nil {
			panic(fmt.Sprintf("error in reading from file %s: %v", name, err))
		}
	// HTTP
	} else {
		resp, err := http.Get(name)
		if err != nil {
			panic(fmt.Sprintf("provided name=%s is not a file and not a URL: %v", name, err))
		}

		defer resp.Body.Close()
		lines, err = linesFromReader(resp.Body)
		if err != nil {
			panic(fmt.Sprintf("error in reading from %s: %v", name, err))
		}
	}

	return In{
		lines: lines,
	}
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

// ReadAllInts reads all remaining tokens from this input stream, parses them as integers,
// and returns them as an array of integers.
//
// Returns all remaining lines in this input stream, as an array of integers
func (in *In) ReadAllInts() []int {
	fields := in.ReadAllStrings()
	vals := make([]int, len(fields))
	for i, f := range fields {
		n, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("error in parsing %s: %v", f, err))
		}
		vals[i] = int(n)
	}
	return vals
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
