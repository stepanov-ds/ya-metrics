package sender

import (
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stepanov-ds/ya-metrics/internal/collector"
	"github.com/stepanov-ds/ya-metrics/internal/config/agent"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func generateRSAKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return privateKey, &privateKey.PublicKey
}

func createTestServer(handler http.HandlerFunc) string {
	srv := httptest.NewServer(handler)
	return srv.URL
}

func mockCollector() *collector.Collector {
	c := collector.NewCollector(&sync.Map{})
	c.CollectMetrics()
	return c
}

func TestEncryptDecrypt(t *testing.T) {
	_, pub := generateRSAKeys(t)

	plainText := []byte("secret_data")
	payload, err := Encrypt(plainText, pub)
	require.NoError(t, err)

	assert.NotEmpty(t, payload.EncryptedAESKey)
	assert.NotEmpty(t, payload.CipherText)
	assert.NotEmpty(t, payload.Nonce)
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestSendMetric(t *testing.T) {
	metric := utils.Metrics{
		ID:    "test_gauge",
		MType: "gauge",
		Value: new(float64),
	}
	*metric.Value = 1.23

	handler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var received utils.Metrics
		err := json.Unmarshal(body, &received)
		require.NoError(t, err)
		assert.Equal(t, metric.ID, received.ID)
		w.WriteHeader(http.StatusOK)
	}

	serverURL := createTestServer(http.HandlerFunc(handler))

	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")

	sender := NewHTTPSender(5*time.Second, headers, serverURL, 1, nil)

	err := sender.SendMetric(metric, "/update")
	assert.NoError(t, err)
}

func TestSendMetricWithHash(t *testing.T) {
	key := "test_secret_key"
	agent.Key = &key

	metric := utils.Metrics{
		ID:    "test_counter",
		MType: "counter",
		Delta: new(int64),
	}
	*metric.Delta = 42

	handler := func(w http.ResponseWriter, r *http.Request) {
		hash := r.Header.Get("HashSHA256")
		body, _ := io.ReadAll(r.Body)
		expectedHash := utils.CalculateHashWithKey(body, key)

		assert.Equal(t, expectedHash, hash)
		w.WriteHeader(http.StatusOK)
	}

	serverURL := createTestServer(http.HandlerFunc(handler))

	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")

	sender := NewHTTPSender(5*time.Second, headers, serverURL, 1, nil)

	err := sender.SendMetric(metric, "/update")
	assert.NoError(t, err)
}

func TestSendMetricGzip(t *testing.T) {
	metric := utils.Metrics{
		ID:    "test_gauge",
		MType: "gauge",
		Value: new(float64),
	}
	*metric.Value = 1.23

	handler := func(w http.ResponseWriter, r *http.Request) {

		assert.Equal(t, "gzip", r.Header.Get("Content-Encoding"))
		assert.Equal(t, "gzip", r.Header.Get("Accept-Encoding"))

		gzReader, err := gzip.NewReader(r.Body)
		require.NoError(t, err)
		defer gzReader.Close()

		body, err := io.ReadAll(gzReader)
		require.NoError(t, err)

		var received utils.Metrics
		err = json.Unmarshal(body, &received)
		require.NoError(t, err)
		assert.Equal(t, metric.ID, received.ID)
		w.WriteHeader(http.StatusOK)
	}

	serverURL := createTestServer(http.HandlerFunc(handler))

	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")

	sender := NewHTTPSender(5*time.Second, headers, serverURL, 1, nil)

	err := sender.SendMetricGzip(metric, "/update")
	assert.NoError(t, err)
}

func TestSendAll(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := new(sync.WaitGroup)
	wg.Add(1)

	c := mockCollector()

	handler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var received []utils.Metrics
		err := json.Unmarshal(body, &received)
		require.NoError(t, err)
		assert.Len(t, received, 29)
		w.WriteHeader(http.StatusOK)
	}

	serverURL := createTestServer(http.HandlerFunc(handler))

	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")

	s := NewHTTPSender(5*time.Second, headers, serverURL, 1, nil)

	go s.SendAll(ctx, wg, 100*time.Millisecond, c, false)

	time.Sleep(300 * time.Millisecond)
	cancel()
	wg.Wait()
}

func TestNewHttpSender(t *testing.T) {
	type args struct {
		headers http.Header
		baseURL string
		timeout time.Duration
	}
	tests := []struct {
		name string
		want HTTPSender
		args args
	}{
		// TODO: Add test cases.
		{
			name: "Positive #1 create HttpSender",
			args: args{
				timeout: time.Second * 2,
				headers: map[string][]string{
					"Content-Type": {"text/plain"},
				},
				baseURL: "localhost:8080",
			},
			want: HTTPSender{
				BaseURL: "localhost:8080",
				Headers: map[string][]string{
					"Content-Type": {"text/plain"},
				},
				Client: &http.Client{
					Timeout: time.Second * 2,
				},
				sem: make(chan struct{}, 1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cryptoKey := "../../cert.pem"
			sender := NewHTTPSender(tt.args.timeout, tt.args.headers, tt.args.baseURL, 1, agent.ReadPublicKey(cryptoKey).PublicKey.(*rsa.PublicKey))
			assert.Equal(t, tt.want.BaseURL, sender.BaseURL)
			assert.True(t, reflect.DeepEqual(tt.want.Headers, sender.Headers))
		})
	}
}
