package socketiocli

import (
	"math/rand"
)

var (
	randChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func getRandName(cnt int) string {
	buf := make([]byte, cnt)
	maxcnt := len(randChars)
	for i := 0; i < cnt; i++ {
		rndnum := rand.Int31n(int32(maxcnt - 1))
		buf[i] = byte(randChars[rndnum])
	}
	return string(buf)
}
