package browser

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	CantOpenBrowserError = errors.New("cannot open browser")
	OsNotSupportedError  = errors.New("your os is not supported")
)

type Browser interface {
	command(string) (*exec.Cmd, error)
	openTab(string) error
}

var OSs []Browser

func OpenBrowser(s string) error {
	if len(OSs) == 0 {
		return OsNotSupportedError
	}
	switch runtime.GOOS {
	case "linux":
		if os.Getenv("DISPLAY") == "" {
			return fmt.Errorf("tried to open %q, no screen found", s)
		}
		fallthrough
	case "darwin":
		if os.Getenv("SSH_CLIENT") != "" || os.Getenv("SSH_TTY") != "" {
			return fmt.Errorf("tried to open %q, but you are running a shell session", s)
		}
	}
	for _, os := range OSs {
		err := os.openTab(s)
		time.Sleep(time.Second * 1)
		if err == nil {
			return nil
		}
	}
	return CantOpenBrowserError
}

type browserCommand struct {
	cmd  string
	args []string
}

var (
	osCommand = map[string]*browserCommand{
		"darwin":  &browserCommand{"open", nil},
		"linux":   &browserCommand{"xdg-open", nil},
		"windows": &browserCommand{"cmd", []string{"/c", "start"}},
	}
	winSchemes = [3]string{"https", "http", "file"}
)

func init() {
	if os, ok := osCommand[runtime.GOOS]; ok {
		OSs = append(OSs, browserCommand{os.cmd, os.args})
	}
}
func (b browserCommand) command(s string) (*exec.Cmd, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	validUrl := ensureValidURL(u)

	b.args = append(b.args, validUrl)

	return exec.Command(b.cmd, b.args...), nil
}

func (b browserCommand) openTab(s string) error {
	cmd, err := b.command(s)
	if err != nil {
		return err
	}

	return cmd.Run()
}
func ensureScheme(u *url.URL) {
	for _, s := range winSchemes {
		if u.Scheme == s {
			return
		}
	}
	u.Scheme = "http"
}
func ensureValidURL(u *url.URL) string {
	ensureScheme(u)
	s := u.String()
	switch runtime.GOOS {
	case "windows":
		s = strings.Replace(s, "&", `^&`, -1)
	}
	return s
}
