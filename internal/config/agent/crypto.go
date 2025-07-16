package agent

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func ReadPublicKey(file string) *x509.Certificate {
	pemData, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "CERTIFICATE" {
		fmt.Println("не удалось декодировать PEM блок сертификата")
		return nil
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return cert
}
