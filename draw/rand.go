package draw

import (
	"math/rand"
	"time"
)

var (
	runeDigit = []rune("1234567890")
	lenDigit  = len(runeDigit)

	runeAlpha = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	lenAlpha  = len(runeAlpha)

	runeAlphaDigit = []rune("abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	lenAlphaDigit  = len(runeAlphaDigit)
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 随机数字[0-9]
func RandDigit(n int) string {
	buffer := make([]rune, n)
	for i := range buffer {
		buffer[i] = runeDigit[rand.Intn(lenDigit)]
	}
	return string(buffer)
}

// 随机字母[a-zA-Z]
func RandAlpha(n int) string {
	buffer := make([]rune, n)
	for i := range buffer {
		buffer[i] = runeAlpha[rand.Intn(lenAlpha)]
	}
	return string(buffer)
}

// 随机字母+数字[0-9a-zA-Z]
func RandAlphaDigit(n int) string {
	buffer := make([]rune, n)
	for i := range buffer {
		buffer[i] = runeAlphaDigit[rand.Intn(lenAlphaDigit)]
	}
	return string(buffer)
}
