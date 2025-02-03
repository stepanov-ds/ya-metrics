package sender

import (
	"fmt"
	"net/http"
	"time"

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

func (s HTTPSender) SendMetric(name string, metric utils.Metric) (*http.Response, error) {
	var path string
	if metric.IsCounter {
		path = fmt.Sprintf("/update/counter/%s/%d", name, metric.Counter)
	} else {
		path = fmt.Sprintf("/update/gauge/%s/%f", name, metric.Gauge)
	}
	req, err := http.NewRequest(http.MethodPost, s.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header = s.Headers
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, err
}
