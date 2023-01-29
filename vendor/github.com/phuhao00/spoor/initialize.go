package spoor

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var (
	pid      = 0
	program  = ""
	host     = ""
	userName = ""
	pidStr   = ""
)

func init() {
	pid = os.Getpid()
	program = filepath.Base(os.Args[0])
	pidStr = fmt.Sprintf(" pid:%05d ", pid)

	h, err := os.Hostname()
	if err == nil {
		host = shortHostname(h)
	}

	current, err := user.Current()
	if err == nil {
		userName = current.Username
	}

	userName = strings.Replace(userName, `\`, "_", -1)
}

func shortHostname(hostname string) string {
	if i := strings.Index(hostname, "."); i >= 0 {
		return hostname[:i]
	}
	return hostname
}
