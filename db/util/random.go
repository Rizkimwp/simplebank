package util

import (
	"math/rand"
	"strings"
	"time"
)


const alpabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}


func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alpabet)
	for i := 0; i < n; i++ {
		c := alpabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}


func RandomOwner() string {
	return RandomString(6)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomBalance() int64 {
	return RandomInt(0, 1000)
}