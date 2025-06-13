package router

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stepanov-ds/ya-metrics/internal/storage"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestRoute_RegisterAllRoutes(t *testing.T) {
	// Создаем фиктивное хранилище и пул
	var st storage.Storage
	var p *pgxpool.Pool

	r := setupRouter()
	Route(r, st, p)

	// Получаем список всех зарегистрированных маршрутов
	routes := r.Routes()

	// Проверяем конкретные маршруты
	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"POST", "/update"},
		{"POST", "/value"},
		{"GET", "/"},
		{"GET", "/ping"},
		{"POST", "/updates"},
	}

	for _, expected := range expectedRoutes {
		found := false
		for _, route := range routes {
			if route.Method == expected.method && route.Path == expected.path {
				found = true
				break
			}
		}
		assert.True(t, found, "Route not found: %s %s", expected.method, expected.path)
	}
}
