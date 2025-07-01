// check zone util
// https://github.com/isc-projects/bind9/blob/main/bin/check/named-compilezone.rst

package dnt

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/itoolkits/toolkit/cmdt"
	"github.com/itoolkits/toolkit/filet"
	"github.com/itoolkits/toolkit/random"
)

type CheckZone struct {
	Bin string

	InputFmt  string
	OutputFmt string

	Output string // dir or path
	Input  string

	Zone string
	View string

	JNL bool

	Timeout time.Duration
}

// CompileRR - compile to rr list
func (c *CheckZone) CompileRR() ([]*RR, error) {
	if c.Output == "" {
		c.Output = "/tmp"
	}

	err := filet.MkDir(c.Output)
	if err != nil {
		return nil, err
	}

	c.Output = filepath.Join(c.Output,
		fmt.Sprintf("%s.zone.txt", random.GenStrByLen(10)))

	defer func() {
		err := filet.RemoveFile(c.Output)
		if err != nil {
			slog.Error("remove file error", "file", c.Output, "error", err.Error())
		}
	}()

	_, err = c.Compile()
	if err != nil {
		return nil, err
	}

	rrs, err := ParseTextZoneFile(c.Output, c.Zone, c.View)
	if err != nil {
		return nil, err
	}
	return rrs, nil
}

// Compile compile zone
func (c *CheckZone) Compile() ([]string, error) {
	if c.Timeout < 1 {
		c.Timeout = time.Minute * 10
	}
	args := make([]string, 0)
	if c.InputFmt != "" {
		args = append(args, "-f", c.InputFmt)
	}
	if c.OutputFmt != "" {
		args = append(args, "-F", c.OutputFmt)
	}
	if c.Output != "" {
		args = append(args, "-o", c.Output)
	}
	if c.JNL {
		args = append(args, "-j")
	}

	args = append(args, c.Zone, c.Input)

	cmd := cmdt.NewCommand(c.Timeout)

	return cmd.Do(c.Bin, args...)
}
