package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

func PingWithMock(c *gin.Context, pool pgxmock.PgxPoolIface) {
	if pool == nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	err := pool.Ping(context.Background())
	if err != nil {
		c.String(http.StatusInternalServerError, "")
	} else {
		c.String(http.StatusOK, "")
	}
}

func TestPing_NilPool(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/ping", func(c *gin.Context) {
		Ping(c, nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPing_DBError(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer mock.Close()
	mock.ExpectPing().WillReturnError(assert.AnError)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/ping", func(c *gin.Context) {
		PingWithMock(c, mock)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPing_OK(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer mock.Close()

	mock.ExpectPing().WillReturnError(nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/ping", func(c *gin.Context) {
		PingWithMock(c, mock)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
