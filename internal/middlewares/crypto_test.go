package middlewares

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stepanov-ds/ya-metrics/internal/config/server"
	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
	"github.com/stretchr/testify/assert"
)

func encrypt(plainText []byte, publicKey *rsa.PublicKey) (*utils.EncryptedPayload, error) {
	aesKey := make([]byte, 32)
	if _, err := rand.Read(aesKey); err != nil {
		return nil, fmt.Errorf("ошибка генерации AES ключа: %v", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания AES шифра: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания GCM: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("ошибка генерации nonce: %v", err)
	}

	cipherText := gcm.Seal(nil, nonce, plainText, nil)

	encryptedAESKey, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, aesKey)
	if err != nil {
		return nil, fmt.Errorf("ошибка шифрования AES ключа RSA: %v", err)
	}

	payload := &utils.EncryptedPayload{
		EncryptedAESKey: base64.StdEncoding.EncodeToString(encryptedAESKey),
		CipherText:      base64.StdEncoding.EncodeToString(cipherText),
		Nonce:           base64.StdEncoding.EncodeToString(nonce),
	}

	return payload, nil
}

func readPublicKey(file string) *x509.Certificate {
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

func setupTestRouterWithCrypto(privKey *rsa.PrivateKey) *gin.Engine {
	logger.Initialize("fatal")
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(Crypto(privKey))

	r.POST("/test", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		c.String(http.StatusOK, string(body))
	})

	return r
}

func TestCryptoMiddleware_DecryptsSuccessfully(t *testing.T) {
	privKey := server.ReadPrivateKey("../../private_key.pem")
	r := setupTestRouterWithCrypto(privKey)

	originalData := []byte(`{"test": "value"}`)
	encrypted, err := encrypt(originalData, readPublicKey("../../cert.pem").PublicKey.(*rsa.PublicKey))
	assert.NoError(t, err)

	bodyBytes, _ := json.Marshal(encrypted)

	req, _ := http.NewRequest("POST", "/test", bytes.NewReader(bodyBytes))
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, string(originalData), resp.Body.String())
}

func TestCryptoMiddleware_InvalidJSON(t *testing.T) {
	privKey := server.ReadPrivateKey("../../private_key.pem")
	r := setupTestRouterWithCrypto(privKey)

	req, _ := http.NewRequest("POST", "/test", strings.NewReader("not a JSON"))
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestCryptoMiddleware_NoCryptoKey(t *testing.T) {
	privKey := server.ReadPrivateKey("")
	r := setupTestRouterWithCrypto(privKey)

	originalData := []byte(`{"test": "value"}`)

	req, _ := http.NewRequest("POST", "/test", bytes.NewReader(originalData))
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, string(originalData), resp.Body.String())
}
