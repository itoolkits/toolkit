// rsync util

package cmdt

import (
	"fmt"
	"strings"
	"time"
)

type RSync struct {
	Timeout time.Duration
	User    string
	RHost   string
	RPath   string
	LPath   string

	Args []string
}

// Do execute
func (r *RSync) Do() ([]string, error) {
	if r.Args == nil {
		r.Args = make([]string, 0)
	}
	if len(r.Args) < 1 {
		r.Args = append(r.Args, "-avz")
		r.Args = append(r.Args, "-e", "ssh")
	}
	if r.Timeout < 1 {
		r.Timeout = time.Minute * 10
	}

	r.Args = append(r.Args, fmt.Sprintf("%s@%s:%s", r.User, r.RHost, r.RPath))
	r.Args = append(r.Args, r.LPath)

	lines, err := NewCommand(r.Timeout).Do("rsync ", r.Args...)
	if err != nil {
		return nil, err
	}
	fileList := make([]string, 0)
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		if strings.Contains(l, "receiving incremental file list") {
			continue
		}
		if strings.HasPrefix(l, "sent ") {
			continue
		}
		if strings.HasPrefix(l, "total size is ") {
			continue
		}
		fileList = append(fileList, l)
	}
	return fileList, nil
}
