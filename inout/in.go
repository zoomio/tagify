package inout

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gpestana/htmlizer"
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

		// will trim out all the tabs from text
		hizer, err := htmlizer.New([]rune{'\t'})
		if err != nil {
			panic(fmt.Sprintf("error in triming content from %s: %v", name, err))
		}

		for i, line := range lines {
			hizer.Load(line)
			lines[i] = hizer.HumanReadable()
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
