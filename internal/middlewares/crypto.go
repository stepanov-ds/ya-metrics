package middlewares

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
)

func decrypt(payload *utils.EncryptedPayload, privKey *rsa.PrivateKey) ([]byte, error) {
	encryptedAESKey, err := base64.StdEncoding.DecodeString(payload.EncryptedAESKey)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования AES ключа: %v", err)
	}

	cipherText, err := base64.StdEncoding.DecodeString(payload.CipherText)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования данных: %v", err)
	}

	nonce, err := base64.StdEncoding.DecodeString(payload.Nonce)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования nonce: %v", err)
	}

	aesKey, err := rsa.DecryptPKCS1v15(rand.Reader, privKey, encryptedAESKey)
	if err != nil {
		return nil, fmt.Errorf("ошибка расшифровки AES ключа: %v", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания AES шифра: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания GCM: %v", err)
	}

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка расшифровки данных: %v", err)
	}

	return plainText, nil
}

func Crypto(privKey *rsa.PrivateKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		if privKey != nil {
			var encryptedPayload utils.EncryptedPayload
			if err := c.ShouldBindBodyWithJSON(&encryptedPayload); err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			decrypted, err := decrypt(&encryptedPayload, privKey)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(decrypted))
		}
		c.Next()
	}
}
