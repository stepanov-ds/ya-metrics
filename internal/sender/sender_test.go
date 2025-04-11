package sender

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stepanov-ds/ya-metrics/internal/utils"
	"github.com/stretchr/testify/assert"
)

type MockClient struct {
	http.Client
	BaseURL string
	Headers http.Header
	DoFunc  func(*http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestHttpSender_SendMetric(t *testing.T) {
	type args struct {
		metric utils.Metrics
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Request
		wantErr error
	}{
		{},
		// TODO: Add test cases.
		// {
		// 	name: "Positive #1 send gauge",
		// 	args: args{
		// 		name:   "testGauge",
		// 		metric: utils.NewMetricGauge(123.1),
		// 	},
		// 	want: func() *http.Request {
		// 		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:8080/update/gauge/testGauge/%f", 123.1), nil)
		// 		req.Header = map[string][]string{
		// 			"Content-Type": {"text/plain"},
		// 		}
		// 		return req
		// 	}(),
		// 	wantErr: nil,
		// },
		// {
		// 	name: "Positive #2 send counter",
		// 	args: args{
		// 		name:   "testCounter",
		// 		metric: utils.NewMetricCounter(123),
		// 	},
		// 	want: func() *http.Request {
		// 		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:8080/update/counter/testCounter/%d", 123), nil)
		// 		req.Header = map[string][]string{
		// 			"Content-Type": {"text/plain"},
		// 		}
		// 		return req
		// 	}(),
		// 	wantErr: nil,
		// },
		// {
		// 	name: "Negative #1 receive error while send metric",
		// 	args: args{
		// 		name:   "testCounter",
		// 		metric: utils.NewMetricCounter(123),
		// 	},
		// },
		// {
		// 	name: "Negative #2 receive error while creating request",
		// 	args: args{
		// 		name:   "testCounter",
		// 		metric: utils.NewMetricCounter(123),
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.name) >= 8 && tt.name[:8] == "Positive" {
				MockClient := &MockClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						assert.Equal(t, tt.want.URL, req.URL)
						assert.Equal(t, tt.want.Method, req.Method)
						assert.Equal(t, tt.want.Body, req.Body)
						assert.Equal(t, tt.want.Header, req.Header)
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       http.NoBody,
						}, nil
					},
				}
				sender := &HTTPSender{
					BaseURL: "http://localhost:8080",
					Headers: map[string][]string{
						"Content-Type": {"text/plain"},
					},
					Client: MockClient,
				}
				err := sender.SendMetric(tt.args.metric, "/update")
				if err != nil {
					assert.Fail(t, err.Error())
				}
			}
			if len(tt.name) >= 11 && tt.name[:11] == "Negative #1" {
				MockClient := &MockClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						return nil, errors.New("mock error")
					},
				}
				sender := &HTTPSender{
					BaseURL: "http://localhost:8080",
					Headers: map[string][]string{
						"Content-Type": {"text/plain"},
					},
					Client: MockClient,
				}
				err := sender.SendMetric(tt.args.metric, "/update")
				assert.Error(t, err)
			}
			if len(tt.name) >= 11 && tt.name[:11] == "Negative #2" {
				MockClient := &MockClient{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						return nil, nil
					},
				}
				sender := &HTTPSender{
					BaseURL: "\t",
					Headers: map[string][]string{
						"Content-Type": {"text/plain"},
					},
					Client: MockClient,
				}
				err := sender.SendMetric(tt.args.metric, "/update")
				assert.Error(t, err)
			}
		})
	}
}

func TestNewHttpSender(t *testing.T) {
	type args struct {
		timeout time.Duration
		headers http.Header
		baseURL string
	}
	tests := []struct {
		name string
		args args
		want HTTPSender
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewHTTPSender(tt.args.timeout, tt.args.headers, tt.args.baseURL, 1))
		})
	}
}
