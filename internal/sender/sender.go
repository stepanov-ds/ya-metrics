package sender

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"time"

	"github.com/stepanov-ds/ya-metrics/internal/collector"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
)

type Sender interface {
	SendMetric(name string, metric utils.Metrics) (*http.Response, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HTTPSender struct {
	BaseURL string
	Headers http.Header
	Client  HTTPClient
}

func NewHTTPSender(timeout time.Duration, headers http.Header, baseURL string) HTTPSender {
	return HTTPSender{
		BaseURL: baseURL,
		Headers: headers,
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (s *HTTPSender) SendMetric(name string, m utils.Metrics) error {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, s.BaseURL+"/update", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header = s.Headers
	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return err
}

func (s *HTTPSender) SendMetricGzip(name string, m utils.Metrics) error {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)

	if _, err := gzWriter.Write(jsonBytes); err != nil {
		return err
	}
	if err := gzWriter.Close(); err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, s.BaseURL+"/update", &buf)
	if err != nil {
		return err
	}
	req.Header = s.Headers
	req.Header.Add("Content-Encoding", "gzip")
	req.Header.Add("Accept-Encoding", "gzip")

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return err
}

func (s *HTTPSender) send(interval time.Duration, collector *collector.Collector, gzip bool) {
	for {
		for k, v := range collector.GetAllMetrics() {
			if gzip {
				err := s.SendMetricGzip(k, v)
				if err != nil {
					println(err.Error())
				}
			} else {
				err := s.SendMetric(k, v)
				if err != nil {
					println(err.Error())
				}
			}
		}
		time.Sleep(interval)
	}
}

func (s *HTTPSender) Send(interval time.Duration, collector *collector.Collector, gzip bool) {
	go s.send(interval, collector, gzip)
}
