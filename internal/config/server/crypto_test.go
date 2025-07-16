package server

import "testing"

func TestReadPublicKey_Success(t *testing.T) {
	cryptoKey := "../../../private_key.pem"
	cert := ReadPrivateKey(cryptoKey)
	if cert == nil {
		t.Error("expected certificate, got nil")
	}
}
