package reader

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Reader - Input. This struct provides methods for reading strings
// and numbers from standard input, file input, URLs, and sockets.
type Reader struct {
	reader  io.Reader
	scanner *bufio.Scanner
}

// NewIn initializes an input from STDIN, file or web page.
//
// name - the filename or web page name, reads from STDIN if name is empty.
// Panics on errors.
func NewIn(name string) (Reader, error) {
	var reader io.Reader
	in := Reader{}

	// STDIN
	if name == "" {
		stat, err := os.Stdin.Stat()
		if err != nil {
			return in, fmt.Errorf("error in reading from STDIN: %v", err)
		}
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return in, errors.New("unsupported mode")
		}
		reader = bufio.NewReader(os.Stdin)
		// File system
	} else if _, err := os.Stat(name); err == nil {
		f, err := os.Open(name)
		if err != nil {
			return in, fmt.Errorf("error in opening file %s for reading: %v", name, err)
		}
		reader = bufio.NewReader(f)
		// HTTP
	} else {
		resp, err := http.Get(name)
		if err != nil {
			return in, fmt.Errorf("provided name=%s is not a file and not a URL: %v", name, err)
		}
		reader = resp.Body
	}

	return Reader{
		reader:  reader,
		scanner: bufio.NewScanner(reader),
	}, nil
}

// NewInFromString initializes an input from string.
func NewInFromString(input string) Reader {
	r := strings.NewReader(input)
	return Reader{
		reader:  r,
		scanner: bufio.NewScanner(r),
	}
}

// ReadString ...
func (in *Reader) ReadString() (string, error) {
	var text string
	if in.scanner.Scan() {
		text = in.scanner.Text()
	}
	err := in.scanner.Err()
	if err != nil {
		return "", err
	}
	return text, nil
}

// ReadAllStrings provides slice of strings from input split by white space.
func (in *Reader) ReadAllStrings() ([]string, error) {
	tokens := make([]string, 0)
	lines, err := in.linesFromReader()
	if err != nil {
		return tokens, fmt.Errorf("error in reading from scanner: %v", err)
	}
	for _, line := range lines {
		tokens = append(tokens, strings.Fields(line)...)
	}
	return tokens, nil
}

// Close ...
func (in *Reader) Close() {
	if closer, ok := in.reader.(io.Closer); ok {
		err := closer.Close()
		if err != nil {
			panic(err)
		}
	}
}

func (in *Reader) linesFromReader() ([]string, error) {
	defer in.Close()
	var lines []string
	for in.scanner.Scan() {
		lines = append(lines, in.scanner.Text())
	}
	if err := in.scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
