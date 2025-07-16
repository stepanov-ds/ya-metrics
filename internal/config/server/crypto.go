package server

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/stepanov-ds/ya-metrics/internal/logger"
)

func ReadPrivateKey(file string) *rsa.PrivateKey {
	pemData, err := os.ReadFile(file)
	if err != nil {
		logger.Log.Info(err.Error())
		return nil
	}

	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		logger.Log.Error("failed to decode PEM block containing private key")
		return nil
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil
	}

	return privKey
}
