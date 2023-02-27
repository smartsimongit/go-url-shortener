package services

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

var key = generateRandom(16)

func SignName(userName string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(userName))
	dst := h.Sum(nil)
	fmt.Printf("%x", dst)
	return dst
}
func generateRandom(size int) []byte {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil
	}
	return b
}
func VerifyToken(token string) bool {
	tokenBytes, err := hex.DecodeString(token)
	if err != nil {
		return false
	}
	name := tokenBytes[:8]
	sign := tokenBytes[8:]
	fmt.Println("name:", name)                                                    //TODO: убрать
	fmt.Println("pass:", sign)                                                    //TODO: убрать
	fmt.Println("Подпись подлинная - ", hmac.Equal(SignName(string(name)), sign)) //TODO: убрать
	return hmac.Equal(SignName(string(name)), sign)
}

func GetUserFromToken(token string) string {
	tokenBytes, err := hex.DecodeString(token)
	if err != nil {
		return ""
	}
	name := tokenBytes[:8]
	return string(name)
}
