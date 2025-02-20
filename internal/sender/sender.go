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
	SendMetric(name string, metric utils.Metric) (*http.Response, error)
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

func (s *HTTPSender) SendMetric(name string, m utils.Metric) (*http.Response, error) {
	metric := m.ConstructJsonObj(name)
	jsonBytes, err := json.Marshal(metric)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, s.BaseURL+"/update", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header = s.Headers
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp, err
}

func (s *HTTPSender) SendMetricGzip(name string, m utils.Metric) (*http.Response, error) {
	metric := m.ConstructJsonObj(name)
	jsonBytes, err := json.Marshal(metric)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)

	if _, err := gzWriter.Write(jsonBytes); err != nil {
		return nil, err
	}
	if err := gzWriter.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, s.BaseURL+"/update", &buf)
	if err != nil {
		return nil, err
	}
	req.Header = s.Headers
	req.Header.Add("Content-Encoding", "gzip")
	req.Header.Add("Accept-Encoding", "gzip")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp, err
}

func (s *HTTPSender) send(interval time.Duration, collector *collector.Collector, gzip bool) {
	for {
		for k, v := range collector.GetAllMetrics() {
			if gzip {
				_, err := s.SendMetricGzip(k, v)
				if err != nil {
					println(err.Error())
				}
			} else {
				_, err := s.SendMetric(k, v)
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
