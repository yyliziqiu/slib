package sstr

import (
	"math/rand"
	"strings"
)

func Truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func TrimSplit(str string, sep string) []string {
	ret := strings.Split(str, sep)
	for i := 0; i < len(ret); i++ {
		ret[i] = strings.TrimSpace(ret[i])
	}
	return ret
}

var (
	_randDigits     = "0123456789"
	_randAlphabets  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	_randAllCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func RandDigits(length int) string {
	return Rand(_randDigits, length)
}

func RandAlphabets(length int) string {
	return Rand(_randAlphabets, length)
}

func RandString(length int) string {
	return Rand(_randAllCharset, length)
}

func Rand(charset string, length int) string {
	var sb strings.Builder
	sb.Grow(length)

	for i := 0; i < length; i++ {
		sb.WriteByte(charset[rand.Intn(62)])
	}

	return sb.String()
}
