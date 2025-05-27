package utils

import (
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

func NewUUID() string {
	return uuid.NewString()
}

func NewGuidStr() string {
	return strings.ReplaceAll(NewUUID(), "-", "")
}

// GenerateRandomNumber 生成6位随机数
func GenerateRandomNumber() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(900000) + 100000
}

// GenerateRandomStringWithLength 生成指定长度的随机字符串
func GenerateRandomStringWithLength(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(rand.Intn(1000))))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
