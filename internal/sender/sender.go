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
	"encoding/json"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/stepanov-ds/ya-metrics/internal/collector"
	"github.com/stepanov-ds/ya-metrics/internal/config/agent"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
)

// Sender is an interface for sending individual metrics.
type Sender interface {
	SendMetric(name string, metric utils.Metrics) (*http.Response, error)
}

// HTTPClient is an interface wrapping the HTTP client's Do method.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// HTTPSender implements metric sending via HTTP requests.
type HTTPSender struct {
	BaseURL string
	Headers http.Header
	Client  HTTPClient
	sem     chan struct{}
}

// NewHTTPSender creates and returns a new HTTPSender instance.
//
// Initializes:
// - Base URL for the server
// - Headers to be used in each request
// - HTTP client with timeout
// - Semaphore based on rate limit
func NewHTTPSender(timeout time.Duration, headers http.Header, baseURL string, rateLimit int) HTTPSender {
	return HTTPSender{
		BaseURL: baseURL,
		Headers: headers,
		Client: &http.Client{
			Timeout: timeout,
		},
		sem: make(chan struct{}, rateLimit),
	}
}

// SendMetric sends a single metric to the server using HTTP POST.
//
// Applies exponential backoff retry strategy if request fails.
func (s *HTTPSender) SendMetric(m interface{}, path string) error {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return err
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
		resp, err := s.Client.Do(req)
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

	if _, err := gzWriter.Write(jsonBytes); err != nil {
		return err
	}
	if err := gzWriter.Close(); err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, s.BaseURL+path, &buf)
	if err != nil {
		return err
	}
	req.Header = s.Headers.Clone()
	req.Header.Add("Content-Encoding", "gzip")
	req.Header.Add("Accept-Encoding", "gzip")
	if *agent.Key != "" {
		req.Header.Add("HashSHA256", hashString)
	}

	operation := func() (string, error) {
		resp, err := s.Client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		return "", err
	}

	_, err = backoff.RetryWithData(operation, utils.NewOneThreeFiveBackOff())
	return err
}

// send sends all metrics individually at the specified interval.
//
// Uses a semaphore to respect configured rate limit.
func (s *HTTPSender) send(interval time.Duration, collector *collector.Collector, gzip bool) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		s.sem <- struct{}{}
		for _, v := range collector.GetAllMetrics() {
			if gzip {
				err := s.SendMetricGzip(v, "/update")
				if err != nil {
					println(err.Error())
				}
			} else {
				err := s.SendMetric(v, "/update")
				if err != nil {
					println(err.Error())
				}
			}
		}
		<-s.sem
	}
}

// sendAll sends all metrics in bulk at the specified interval.
//
// Uses a semaphore to respect configured rate limit.
func (s *HTTPSender) sendAll(interval time.Duration, collector *collector.Collector, gzip bool) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
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

// Send starts the background loop that periodically sends metrics to the server.
func (s *HTTPSender) Send(interval time.Duration, collector *collector.Collector, gzip bool) {
	go s.sendAll(interval, collector, gzip)
}
