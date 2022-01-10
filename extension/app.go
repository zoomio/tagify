package extension

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/zoomio/inout"
)

const (
	install_path_prefix = "$HOME/bin/tagify_"

	flagVersion = "-version"
)

// App represents the executable for the extension.
type App struct {
	name        string
	source      string
	file        string
	version     string
	execTimeout time.Duration
}

// New creates new instance of App based on the provided parameters.
func New(name, source string) *App {
	return &App{
		name:        name,
		source:      source,
		file:        "",
		version:     "",
		execTimeout: time.Second * 30, // wait for 30 seconds before interrupting the execution
	}
}

// Install installs the App.
func (a *App) Install(ctx context.Context) error {
	if a.source == "" {
		return fmt.Errorf("source is empty for %q", a.name)
	}
	u, err := url.Parse(a.source)
	if err != nil {
		return fmt.Errorf("wrong source format %q: %w", a.source, err)
	}
	r, err := inout.New(ctx, a.source)
	if err != nil {
		return err
	}
	defer r.Close()
	fmt.Printf("downloading %q from %q.\n", a.name, a.source)
	bs, err := io.ReadAll(&r)
	if err != nil {
		return fmt.Errorf("failed to download source %q: %w", a.source, err)
	}
	fName := toFileName(u.Hostname() + "/" + u.Path)
	err = os.WriteFile(fName, bs, os.FileMode(0755))
	if err != nil {
		return fmt.Errorf("failed to save file %q: %w", fName, err)
	}
	a.file = fName
	bs, err = a.Run(ctx, flagVersion)
	if err != nil {
		return fmt.Errorf("failed to install %q: %w", a.name, err)
	}
	a.version = string(bs)
	fmt.Printf("extension %q [%s] has been successfully installed to %q.\n", a.name, a.version, a.file)
	return nil
}

// Run executes the App.
func (a *App) Run(ctx context.Context, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, a.execTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, a.file, args...)
	return cmd.Output()
}

func toFileName(s string) string {
	fName := strings.Replace(s, ".", "-", -1)
	fName = strings.Replace(fName, "/", "_", -1)
	return install_path_prefix + fName
}
