package agent

import (
	"crypto/rsa"
	"testing"
)

func TestReadPublicKey_Success(t *testing.T) {
	cryptoKey := "../../../cert.pem"
	cert := ReadPublicKey(cryptoKey).PublicKey.(*rsa.PublicKey)
	if cert == nil {
		t.Error("expected certificate, got nil")
	}
}
func TestReadPublicKey_NoFile(t *testing.T) {
	cryptoKey := "../../../asdasdas.pem"
	cert := ReadPublicKey(cryptoKey)
	if cert != nil {
		t.Error("expected no cert, got cert")
	}
}
func TestReadPublicKey_BadCert1(t *testing.T) {
	cryptoKey := "../../../testconfigs/testCert1.pem"
	cert := ReadPublicKey(cryptoKey)
	if cert != nil {
		t.Error("expected no cert, got cert")
	}
}
func TestReadPublicKey_BadCert2(t *testing.T) {
	cryptoKey := "../../../testconfigs/testCert2.pem"
	cert := ReadPublicKey(cryptoKey)
	if cert != nil {
		t.Error("expected no cert, got cert")
	}
}
