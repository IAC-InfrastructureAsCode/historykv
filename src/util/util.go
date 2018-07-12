package util

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func CreateHashPassword(user string, password string) string {
	hashString := fmt.Sprintf("historykv-%s-%s", user, password)

	loginMD5 := md5.New()
	loginMD5.Write([]byte(hashString))
	return hex.EncodeToString(loginMD5.Sum(nil))
}

func CreateRandomToken() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func CreateSessionToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
