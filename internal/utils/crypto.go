package utils

type EncryptedPayload struct {
	EncryptedAESKey string `json:"aes_key"` // base64
	CipherText      string `json:"data"`    // base64
	Nonce           string `json:"nonce"`   // base64
}
