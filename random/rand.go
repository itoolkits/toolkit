// random util

package random

import (
	"math/rand"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenStrByLen generate string by length
func GenStrByLen(length int) string {
	rst := make([]byte, length)
	for i := range rst {
		rst[i] = charset[rand.Intn(len(charset))]
	}
	return string(rst)
}
