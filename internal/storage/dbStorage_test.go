package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
)

func init() {
	logger.Initialize("fatal")
}

type MockPool struct {
	mock.Mock
}

func (m *MockPool) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	argsCall := m.Called(ctx, sql, args)
	row, _ := argsCall.Get(0).(pgx.Row)
	return row
}

func (m *MockPool) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	argsCall := m.Called(ctx, sql, args)
	var tag pgconn.CommandTag
	if res := argsCall.Get(0); res != nil {
		tag = res.(pgconn.CommandTag)
	}
	return tag, argsCall.Error(1)
}

func (m *MockPool) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	argsCall := m.Called(ctx, sql, args)
	rows, _ := argsCall.Get(0).(pgx.Rows)
	return rows, argsCall.Error(1)
}

func (m *MockPool) Begin(ctx context.Context) (pgx.Tx, error) {
	argsCall := m.Called(ctx)
	tx, _ := argsCall.Get(0).(pgx.Tx)
	return tx, argsCall.Error(1)
}

type MockTx struct {
	mock.Mock
}

func (m *MockTx) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTx) Rollback(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTx) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	argsCall := m.Called(ctx, sql, args)
	var tag pgconn.CommandTag
	if res := argsCall.Get(0); res != nil {
		tag = res.(pgconn.CommandTag)
	}
	return tag, argsCall.Error(1)
}

func (m *MockTx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	argsCall := m.Called(ctx, sql, args)
	rows, _ := argsCall.Get(0).(pgx.Rows)
	return rows, argsCall.Error(1)
}

type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}

type MockRows struct {
	mock.Mock
}

func (m *MockRows) Scan(dest ...interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}

func (m *MockRows) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockRows) Close() {
	m.Called()
}

func TestDBStorage_GetMetric_Success(t *testing.T) {
	mockPool := new(MockPool)
	dbStorage := &DBStorage{
		Pool: mockPool,
	}
	key := "test_metric"

	met := utils.NewMetrics(key, 3.14, false)

	mockRow := new(MockRow)
	mockRow.On("Scan", mock.MatchedBy(func(dest []interface{}) bool {
		if len(dest) != 4 {
			return false
		}
		*(dest[0].(*string)) = met.ID
		*(dest[1].(*string)) = met.MType
		*(dest[2].(**int64)) = met.Delta
		*(dest[3].(**float64)) = met.Value
		return true
	})).Return(nil)

	mockPool.On("QueryRow", context.Background(),
		"SELECT \"ID\", \"MType\", \"Delta\", \"Value\" FROM public.metrics WHERE \"ID\" = $1;",
		[]interface{}{key}).Return(mockRow)

	metric, found := dbStorage.GetMetric(key)
	require.True(t, found)
	assert.Equal(t, key, metric.ID)
	assert.Equal(t, "gauge", metric.MType)
	assert.InDelta(t, 3.14, *metric.Value, 0.001)
}

func TestDBStorage_GetMetric_NotFound(t *testing.T) {
	mockPool := new(MockPool)
	dbStorage := &DBStorage{
		Pool: mockPool,
	}

	key := "unknown"

	mockRow := new(MockRow)
	mockRow.On("Scan", mock.Anything).Return(errors.New("no rows"))

	mockPool.On("QueryRow", context.Background(),
		"SELECT \"ID\", \"MType\", \"Delta\", \"Value\" FROM public.metrics WHERE \"ID\" = $1;",
		[]interface{}{key}).Return(mockRow)

	metric, found := dbStorage.GetMetric(key)
	assert.False(t, found)
	assert.Equal(t, utils.Metrics{}, metric)
}

func TestDBStorage_SetMetric_Gauge(t *testing.T) {
	mockPool := new(MockPool)
	dbStorage := &DBStorage{
		Pool: mockPool,
	}

	ctx := context.Background()
	key := "test_gauge"
	value := 3.14

	mockPool.On("Exec",
		ctx,
		mock.Anything,
		mock.MatchedBy(func(args []interface{}) bool {
			return len(args) == 3 && args[0] == key && args[1] == "gauge" && args[2] == value
		}),
	).Return(pgconn.CommandTag{}, nil)
	dbStorage.SetMetric(ctx, key, value, false)

	mockPool.AssertExpectations(t)
}

func TestDBStorage_SetMetric_Counter(t *testing.T) {
	mockPool := new(MockPool)
	dbStorage := &DBStorage{
		Pool: mockPool,
	}

	ctx := context.Background()
	key := "test_counter"
	value := int64(42)

	mockPool.On("Exec",
		ctx,
		mock.Anything,
		mock.MatchedBy(func(args []interface{}) bool {
			return len(args) == 3 && args[0] == key && args[1] == "counter" && args[2] == value
		}),
	).Return(pgconn.CommandTag{}, nil)

	dbStorage.SetMetric(ctx, key, value, true)

	mockPool.AssertExpectations(t)
}
