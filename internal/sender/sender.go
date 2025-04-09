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

func (s *HTTPSender) send(interval time.Duration, collector *collector.Collector, gzip bool) {
	for {
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
		time.Sleep(interval) //поменять на time.ticker или time.after
	}
}

func (s *HTTPSender) sendAll(interval time.Duration, collector *collector.Collector, gzip bool) {
	for {
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
		time.Sleep(interval) //поменять на time.ticker или time.after
	}
}

func (s *HTTPSender) Send(interval time.Duration, collector *collector.Collector, gzip bool) {
	go s.sendAll(interval, collector, gzip)
}
