// Package utils contains utility functions and shared types used across the application.
//
// This file provides a function for calculating HMAC-SHA256 hashes with a given key.
package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// CalculateHashWithKey computes HMAC-SHA256 hash of the input data using the provided key.
//
// Returns the hexadecimal string representation of the hash.
func CalculateHashWithKey(data []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)

	hash := h.Sum(nil)

	hashString := hex.EncodeToString(hash)

	return hashString
}
