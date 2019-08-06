package util

import (
	"math/rand"
	"os"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// @source https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func RandStringBytesRmndr(n int) string {
	rand.Seed(int64(os.Getpid()+os.Geteuid()+os.Getppid()+os.Getegid()) + time.Now().Unix())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
