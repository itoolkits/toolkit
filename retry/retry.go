// retry util

package retry

import (
	"math/rand"
	"time"
)

var (
	retryInterval = time.Duration(500+rand.Intn(2000)) * time.Millisecond
)

// Do - dor retry fn
func Do(fn func() error, retries int) error {
	var err error
	for i := 0; i < retries; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		// sleep random
		time.Sleep(retryInterval)
	}
	return err
}
