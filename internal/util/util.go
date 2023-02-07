package util

import (
	"math/rand"
	"net/url"
)

func IsURLInvalid(s string) bool {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return true
	}
	u, err := url.Parse(s)
	if err != nil || u.Host == "" {
		return true
	}
	return false
}

func GenString() string {
	symbolsForGen := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 5)
	for i := range b {
		b[i] = symbolsForGen[rand.Intn(len(symbolsForGen))]
	}
	return string(b)
}

func CreateURL(s string) string {
	return "http://localhost:8080/" + s
}
