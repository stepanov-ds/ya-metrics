//Package grpcp is for grpc functions
package grpcp

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	pb "github.com/stepanov-ds/ya-metrics/internal/grpcp/grpc_generated"
	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"go.uber.org/zap"
)

type GrpcHandler struct {
	pb.UnimplementedMetricsTunnelServer
	Router *gin.Engine
}

func (h *GrpcHandler) HandleHTTPRequest(ctx context.Context, in *pb.HTTPRequestPayload) (*pb.HTTPResponsePayload, error) {
	logger.Log.Info("start grpc HandleHTTPRequest")
	reqBody := bytes.NewReader(in.Body)
	req, err := http.NewRequest(in.Method, in.Path, reqBody)
	if err != nil {
		logger.Log.Error("grpc handler new request error", zap.Error(err))
		return nil, err
	}

	for key, value := range in.Headers.Values {
		req.Header.Set(key, value)
	}

	recorder := httptest.NewRecorder()

	h.Router.ServeHTTP(recorder, req)

	headers := make(map[string]string)
	for key, values := range recorder.Header() {
		headers[key] = strings.Join(values, ", ")
	}

	resBody, err := io.ReadAll(recorder.Body)
	if err != nil {
		logger.Log.Error("grpc handler read response body error", zap.Error(err))
		return nil, err
	}

	response := &pb.HTTPResponsePayload{
		StatusCode: int32(recorder.Code),
		Headers: &pb.Header{Values: headers},
		Body: resBody,
	}
	logger.Log.Info("grpc handler response", zap.String("response", string(response.Body)))
	logger.Log.Info("end grpc HandleHTTPRequest")
	return response, nil
}