// cmd util

package cmdt

import (
	"fmt"
	"time"

	"github.com/go-cmd/cmd"
)

type Command struct {
	timeout time.Duration
}

// NewCommand create command
func NewCommand(timeout time.Duration) *Command {
	return &Command{
		timeout: timeout,
	}
}

// Do execute
func (c *Command) Do(name string, args ...string) ([]string, error) {
	cmdMan := cmd.NewCmd(name, args...)

	statusChan := cmdMan.Start()

	timeC := time.NewTimer(c.timeout)
	defer timeC.Stop()

	select {
	case status := <-statusChan:
		if status.Error != nil {
			return nil, status.Error
		}
		return status.Stdout, status.Error
	case <-timeC.C:
		err := cmdMan.Stop()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(" %s execute command timeout ", name)
	}
}
