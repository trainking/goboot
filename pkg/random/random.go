package random

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
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

// RandStringNumber 创建指定位数的数字字符串; n < 18 , 常用于短信验证码等情况
func RandStringNumber(n int) (string, error) {
	if n > 18 {
		return "", errors.New("n is over 18")
	}

	return fmt.Sprintf("%0"+strconv.FormatInt(int64(n), 10)+"d", rand.Int63n(int64(math.Pow10(n)))), nil
}
