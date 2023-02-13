package util

import (
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"
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
	rand.Seed(time.Now().UnixNano())
	chars := []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String() //
}

func CreateURL(s string) string {
	return getBaseURL() + "/" + s
}

func GetServerAddress() string {
	addr := os.Getenv(serverAddress)
	if addr == "" {
		return ":8080"
	}
	return addr

}
func getBaseURL() string {
	s := os.Getenv(baseURL)
	if s == "" {
		return "http://localhost:8080"
	}
	return s
}

func GetStorageFileName() string {
	s := os.Getenv(fileStorageUrl)
	return s
}
