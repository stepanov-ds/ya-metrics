// Package sender implements logic for sending metrics to a remote server.
//
// It provides:
// - HTTP-based metric sender with retry mechanism
// - Gzip compression support
// - Rate limiting and backoff strategies
package sender

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/stepanov-ds/ya-metrics/internal/collector"
	"github.com/stepanov-ds/ya-metrics/internal/config/agent"
	pb "github.com/stepanov-ds/ya-metrics/internal/grpcp/grpc_generated"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Sender is an interface for sending individual metrics.
type Sender interface {
	SendAll(ctx context.Context, wg *sync.WaitGroup, interval time.Duration, collector *collector.Collector, gzip bool)
}

// HTTPClient is an interface wrapping the HTTP client's Do method.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// HTTPSender implements metric sending via HTTP requests.
type HTTPSender struct {
	Headers   http.Header
	Client    HTTPClient
	CryptoKey *rsa.PublicKey
	sem       chan struct{}
	BaseURL   string
}

type GRPCSender struct {
	Timeout time.Duration
	Conn *grpc.ClientConn
	Headers   http.Header
	Client    pb.MetricsTunnelClient
	CryptoKey *rsa.PublicKey
	sem       chan struct{}
	BaseURL   string
}

// NewHTTPSender creates and returns a new HTTPSender instance.
//
// Initializes:
// - Base URL for the server
// - Headers to be used in each request
// - HTTP client with timeout
// - Semaphore based on rate limit
func NewHTTPSender(timeout time.Duration, headers http.Header, baseURL string, rateLimit int, cryptoKey *rsa.PublicKey) HTTPSender {
	return HTTPSender{
		sem:       make(chan struct{}, rateLimit),
		BaseURL:   baseURL,
		CryptoKey: cryptoKey,
		Headers:   headers,
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}

func NewGRPCSender(timeout time.Duration, headers http.Header, baseURL string, rateLimit int, cryptoKey *rsa.PublicKey) GRPCSender {
	conn, err := grpc.NewClient(*agent.EndpointAgent, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err.Error())
	}
	client := pb.NewMetricsTunnelClient(conn)

	return GRPCSender{
		sem:       make(chan struct{}, rateLimit),
		BaseURL:   baseURL,
		CryptoKey: cryptoKey,
		Headers:   headers,
		Client:    client,
		Timeout: timeout,
		Conn: conn,
	}
}

func Encrypt(plainText []byte, publicKey *rsa.PublicKey) (*utils.EncryptedPayload, error) {
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

// SendMetric sends a single metric to the server using HTTP POST.
//
// Applies exponential backoff retry strategy if request fails.
func (s *HTTPSender) SendMetric(m interface{}, path string) error {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if s.CryptoKey != nil {
		encryptedPayload, err1 := Encrypt(jsonBytes, s.CryptoKey)
		if err1 != nil {
			return err1
		}
		jsonBytes, err = json.Marshal(encryptedPayload)
		if err != nil {
			return err
		}
	}

	var hashString string
	if *agent.Key != "" {
		hashString = utils.CalculateHashWithKey(jsonBytes, *agent.Key)
	}

	req, err := http.NewRequest(http.MethodPost, s.BaseURL+path, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header = s.Headers
	if *agent.Key != "" {
		req.Header.Add("HashSHA256", hashString)
	}

	operation := func() (string, error) {
		var resp *http.Response
		resp, err = s.Client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		return "", err
	}

	_, err = backoff.RetryWithData(operation, utils.NewOneThreeFiveBackOff())
	return err
}

// SendMetricGzip sends a single metric to the server using gzip-compressed HTTP POST.
//
// Applies compression and optional payload signing.
// Uses exponential backoff retry strategy if request fails.
func (s *HTTPSender) SendMetricGzip(m interface{}, path string) error {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	var hashString string
	if *agent.Key != "" {
		hashString = utils.CalculateHashWithKey(jsonBytes, *agent.Key)
	}

	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)

	if _, err = gzWriter.Write(jsonBytes); err != nil {
		return err
	}
	if err = gzWriter.Close(); err != nil {
		return err
	}

	if s.CryptoKey != nil {
		encryptedPayload, err1 := Encrypt(buf.Bytes(), s.CryptoKey)
		if err1 != nil {
			return err1
		}
		encryptedBytes, err2 := json.Marshal(encryptedPayload)
		if err != nil {
			return err2
		}
		buf = *bytes.NewBuffer(encryptedBytes)
	}

	req, err := http.NewRequest(http.MethodPost, s.BaseURL+path, &buf)
	if err != nil {
		return err
	}
	ip, err := getOutboundIPForDest(s.BaseURL)
	if err != nil {
		return err
	}
	req.Header = s.Headers.Clone()
	req.Header.Add("Content-Encoding", "gzip")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("X-Real-IP", ip.String())
	if *agent.Key != "" {
		req.Header.Add("HashSHA256", hashString)
	}

	operation := func() (string, error) {
		var resp *http.Response
		resp, err = s.Client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		return "", err
	}

	_, err = backoff.RetryWithData(operation, utils.NewOneThreeFiveBackOff())
	return err
}

// SendAll sends all metrics in bulk at the specified interval.
//
// Uses a semaphore to respect configured rate limit.
func (s HTTPSender) SendAll(ctx context.Context, wg *sync.WaitGroup, interval time.Duration, collector *collector.Collector, gzip bool) {
	defer wg.Done()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			s.sem <- struct{}{}

			var metrics []utils.Metrics
			for _, v := range collector.GetAllMetrics() {
				metrics = append(metrics, v)
			}
			if gzip {
				err := s.SendMetricGzip(metrics, "/updates")
				if err != nil {
					println(err.Error())
				}
			} else {
				err := s.SendMetric(metrics, "/updates")
				if err != nil {
					println(err.Error())
				}
			}
			<-s.sem
		}
	}
}

func (s GRPCSender) SendAll(ctx context.Context, wg *sync.WaitGroup, interval time.Duration, collector *collector.Collector, gzip bool) {
	println("start sendAll")
	defer wg.Done()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer s.Conn.Close()

	for {
		select {
		case <-ctx.Done():
			println("end sendAll")
			return

		case <-ticker.C:
			s.sem <- struct{}{}

			var metrics []utils.Metrics
			for _, v := range collector.GetAllMetrics() {
				metrics = append(metrics, v)
			}
			err := s.SendMetricGzip(ctx, metrics, "/updates")
			if err != nil {
				println(err.Error())
			}
			<-s.sem
		}
	}
}

func (s *GRPCSender) SendMetricGzip(ctx context.Context, m interface{}, path string) error {
	println("start sendMetricGzip")

	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	var hashString string
	if *agent.Key != "" {
		hashString = utils.CalculateHashWithKey(jsonBytes, *agent.Key)
	}

	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)

	if _, err = gzWriter.Write(jsonBytes); err != nil {
		return err
	}
	if err = gzWriter.Close(); err != nil {
		return err
	}

	if s.CryptoKey != nil {
		encryptedPayload, err1 := Encrypt(buf.Bytes(), s.CryptoKey)
		if err1 != nil {
			return err1
		}
		encryptedBytes, err2 := json.Marshal(encryptedPayload)
		if err2 != nil {
			return err2
		}
		buf = *bytes.NewBuffer(encryptedBytes)
	}
	ip, err := getOutboundIPForDest(s.BaseURL)
	if err != nil {
		return err
	}
	reqHeaders := map[string]string{
		"Content-Encoding": "gzip",
		"Accept-Encoding": "gzip",
		"X-Real-IP": ip.String(),
	}
	println("x-real-ip " + ip.String())
	if *agent.Key != "" {
		reqHeaders["HashSHA256"] = hashString
	}
	grpcReq := &pb.HTTPRequestPayload{
		Method: "POST",
		Path:   path, // Путь к вашему HTTP-роуту
		Headers: &pb.Header{
			Values: reqHeaders,
		},
		Body: buf.Bytes(),
	}

	

	operation := func() (string, error) {
		println("start operation")

		ctx1, cancel := context.WithTimeout(ctx, s.Timeout)
		defer cancel()
		response, err1 := s.Client.HandleHTTPRequest(ctx1, grpcReq)
		if err1 != nil {
			return "", err1
		}
		if response.StatusCode != 200 {
			return "", errors.New("bad status code " + strconv.FormatInt(int64(response.StatusCode), 10))
		}
		println("end operation")

		return "", nil
	}

	_, err = backoff.RetryWithData(operation, utils.NewOneThreeFiveBackOff())
	println("end sendMetricGzip")

	return err
}

func getOutboundIPForDest(destinationAddr string) (net.IP, error) {
	addr := strings.TrimPrefix(destinationAddr, "http://")
	if !strings.Contains(addr, ":") {
		addr = addr + ":80"
	}

	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}
