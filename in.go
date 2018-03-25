package tagify

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Content types
const (
	Unknown ContentType = iota
	Text
	HTML
)

var (
	contentTypes = [...]string{
		"Unknown",
		"Text",
		"HTML",
	}
)

// ContentType ...
type ContentType int

// ContentTypeOf returns ContentType based on string value.
func ContentTypeOf(contentType string) ContentType {
	for i, key := range contentTypes {
		if key == contentType {
			return ContentType(i)
		}
	}
	return Unknown
}

// String ...
func (contentType ContentType) String() string {
	if contentType < Text || contentType > HTML {
		return "Unknown"
	}
	return contentTypes[contentType]
}

// In - Input. This struct provides methods for reading strings
// and numbers from standard input, file input, URLs, and sockets.
type In struct {
	lines []string
	ContentType
}

// NewIn initializes an input stream from STDIN, file or web page.
//
// name - the filename or web page name, reads from STDIN if name is empty.
// Panics on errors.
func NewIn(name string) (In, error) {
	if name == "" {
		return inFromStdIn()
	} else if _, err := os.Stat(name); err == nil {
		return inFromFile(name)
	} else {
		return inFromHTTP(name)
	}
}

// NewInFromString ...
func NewInFromString(input string, contentType ContentType) (In, error) {
	var err error
	in := In{ContentType: contentType}
	in.lines, err = linesFromReader(strings.NewReader(input))
	return in, err
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

func inFromStdIn() (In, error) {
	in := In{}
	stat, err := os.Stdin.Stat()
	if err != nil {
		return in, fmt.Errorf("error in reading from STDIN: %v", err)
	}
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return in, nil
	}
	in.lines, err = linesFromReader(bufio.NewReader(os.Stdin))
	if err != nil {
		return in, fmt.Errorf("error in reading from STDIN: %v", err)
	}
	return in, nil
}

func inFromFile(file string) (In, error) {
	in := In{}
	f, err := os.Open(file)
	if err != nil {
		return in, fmt.Errorf("error in opening file %s for reading: %v", file, err)
	}
	defer f.Close()
	in.lines, err = linesFromReader(bufio.NewReader(f))
	if err != nil {
		return in, fmt.Errorf("error in reading from file %s: %v", file, err)
	}
	return in, nil
}

func inFromHTTP(address string) (In, error) {
	in := In{ContentType: HTML}
	resp, err := http.Get(address)
	if err != nil {
		return in, fmt.Errorf("provided name=\"%s\" is not a file and not a URL: %v", address, err)
	}

	defer resp.Body.Close()
	in.lines, err = linesFromReader(resp.Body)
	if err != nil {
		return in, fmt.Errorf("error in reading from %s: %v", address, err)
	}
	return in, nil
}
