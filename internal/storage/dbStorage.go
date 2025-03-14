package storage

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
	"go.uber.org/zap"
)

type DbStorage struct {
	Pool *pgxpool.Pool
}

func NewDbPool(ctx context.Context, dsn string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.Log.Error("NewDbPool", zap.String("error while creating new DB pool", err.Error()))
		return nil
	}
	return pool
}

func NewDbStorage(p *pgxpool.Pool) *DbStorage {
	query := `
		CREATE TABLE IF NOT EXISTS public.metrics
		(
    		"ID" character varying(255) COLLATE pg_catalog."default" NOT NULL,
    		"MType" character varying(255) COLLATE pg_catalog."default" NOT NULL,
    		"Delta" bigint,
    		"Value" double precision,
    		CONSTRAINT metrics_pkey PRIMARY KEY ("ID")
		)
	`

	operation := func() (string, error) {
		_, err := p.Exec(context.Background(), query)
		return "", err
	}

	_, err := backoff.RetryWithData( operation, utils.NewConstantIncreaseBackOff(time.Second, time.Second*2, 3))
	if err != nil {
		logger.Log.Error("NewDbStorage", zap.String("error while creating table in DB", err.Error()))
	}

	return &DbStorage{
		Pool: p,
	}
}

func (st *DbStorage) GetMetric(key string) (utils.Metrics, bool) {
	query := `SELECT "ID", "MType", "Delta", "Value" FROM public.metrics WHERE "ID" = $1;`

	operation := func() (utils.Metrics, error) {
		row := st.Pool.QueryRow(context.Background(), query, key)

		var m utils.Metrics
		err := row.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
		return m, err
	}

	metric, err := backoff.RetryWithData( operation, utils.NewConstantIncreaseBackOff(time.Second, time.Second*2, 3))
	if err != nil {
		logger.Log.Error("GetMetric", zap.String("error while select from DB", err.Error()))
		return metric, false
	}
	return metric, true
}

func (st *DbStorage) GetAllMetrics() map[string]utils.Metrics {
	query := `SELECT "ID", "MType", "Delta", "Value" FROM public.metrics;`

	operation := func() (pgx.Rows, error) {
		rows, err := st.Pool.Query(context.Background(), query)
		if err != nil {
			return rows, err
		}
		defer rows.Close()
		return rows, err
	}

	rows, err := backoff.RetryWithData( operation, utils.NewConstantIncreaseBackOff(time.Second, time.Second*2, 3))
	if err != nil {
		logger.Log.Error("GetAllMetric", zap.String("error while select from DB", err.Error()))
	}

	metrics := make(map[string]utils.Metrics)
	for rows.Next() {
		var m utils.Metrics
		if err := rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value); err != nil {
			logger.Log.Error("GetAllMetric", zap.String("error while parsing row result", err.Error()))
		}
		metrics[m.ID] = m
	}
	return metrics
}

func (st *DbStorage) SetMetric(key string, value interface{}, counter bool) {
	query1 := `
		INSERT INTO public.metrics ("ID", "MType", "Delta")
		VALUES ($1, $2, $3)
		ON CONFLICT ("ID") DO UPDATE SET
			"MType" = EXCLUDED."MType",
			"Delta" = COALESCE(public.metrics."Delta", 0) + EXCLUDED."Delta",
			"Value" = NULL;
	`
	query2 := `
		INSERT INTO public.metrics ("ID", "MType", "Value")
		VALUES ($1, $2, $3)
		ON CONFLICT ("ID") DO UPDATE SET
			"MType" = EXCLUDED."MType",
			"Value" = EXCLUDED."Value",
			"Delta" = Null;
	`

	operation := func() (string, error) {
		if counter {
			_, err := st.Pool.Exec(context.Background(), query1, key, "counter", value)
			return "", err
		} else {
			_, err := st.Pool.Exec(context.Background(), query2, key, "gauge", value)
			return "", err
		}
	}

	_, err := backoff.RetryWithData( operation, utils.NewConstantIncreaseBackOff(time.Second, time.Second*2, 3))
	if err != nil {
		logger.Log.Error("SetMetric", zap.String("error while insert in DB", err.Error()))
		return
	}
}
