package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Generate a random string for hexadecimal. Docker Id use this.
func RandStringHex(n int) (string, error) {
	readByes := make([]byte, n/2)
	if _, err := rand.Read(readByes); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", readByes), nil
}
