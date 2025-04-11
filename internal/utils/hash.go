package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func CalculateHashWithKey(data []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)

	hash := h.Sum(nil)

	hashString := hex.EncodeToString(hash)

	return hashString
}
