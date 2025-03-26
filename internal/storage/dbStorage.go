package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stepanov-ds/ya-metrics/internal/logger"
	"github.com/stepanov-ds/ya-metrics/internal/utils"
	"go.uber.org/zap"
)

type DBStorage struct {
	Pool *pgxpool.Pool
}

func NewDBPool(ctx context.Context, dsn string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.Log.Error("NewDBPool", zap.String("error while creating new DB pool", err.Error()))
		return nil
	}
	return pool
}

func NewDBStorage(ctx context.Context, p *pgxpool.Pool) *DBStorage {
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
		_, err := p.Exec(ctx, query)
		return "", retriableHelper(err)
	}

	_, err := backoff.RetryWithData(operation, utils.NewOneThreeFiveBackOff())
	if err != nil {
		logger.Log.Error("NewDBStorage", zap.String("error while creating table in DB", err.Error()))
	}

	return &DBStorage{
		Pool: p,
	}
}

func (st *DBStorage) GetMetric(key string) (utils.Metrics, bool) {
	query := `SELECT "ID", "MType", "Delta", "Value" FROM public.metrics WHERE "ID" = $1;`

	operation := func() (utils.Metrics, error) {
		row := st.Pool.QueryRow(context.Background(), query, key)

		var m utils.Metrics
		err := row.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
		return m, retriableHelper(err)
	}

	metric, err := backoff.RetryWithData(operation, utils.NewOneThreeFiveBackOff())
	if err != nil {
		logger.Log.Error("GetMetric", zap.String("error while select from DB", err.Error()))
		return metric, false
	}
	return metric, true
}

func (st *DBStorage) GetAllMetrics() map[string]utils.Metrics {
	query := `SELECT "ID", "MType", "Delta", "Value" FROM public.metrics;`

	operation := func() (pgx.Rows, error) {
		rows, err := st.Pool.Query(context.Background(), query)
		if err == nil {
			defer rows.Close()
		}
		return rows, retriableHelper(err)
	}

	rows, err := backoff.RetryWithData(operation, utils.NewOneThreeFiveBackOff())
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

func (st *DBStorage) SetMetric(ctx context.Context, key string, value interface{}, counter bool) {
	var query string
	var metricType string
	if counter {
		query = `
		INSERT INTO public.metrics ("ID", "MType", "Delta")
		VALUES ($1, $2, $3)
		ON CONFLICT ("ID") DO UPDATE SET
			"MType" = EXCLUDED."MType",
			"Delta" = COALESCE(public.metrics."Delta", 0) + EXCLUDED."Delta",
			"Value" = NULL;
		`
		metricType = "counter"
	} else {
		query = `
		INSERT INTO public.metrics ("ID", "MType", "Value")
		VALUES ($1, $2, $3)
		ON CONFLICT ("ID") DO UPDATE SET
			"MType" = EXCLUDED."MType",
			"Value" = EXCLUDED."Value",
			"Delta" = Null;
		`
		metricType = "gauge"
	}

	operation := func() (string, error) {
		tx, ok := ctx.Value(utils.Transaction).(pgx.Tx)
		if ok {
			_, err := tx.Exec(ctx, query, key, metricType, value)
			return "", retriableHelper(err)
		}

		_, err := st.Pool.Exec(ctx, query, key, metricType, value)
		return "", retriableHelper(err)
	}

	_, err := backoff.RetryWithData(operation, utils.NewOneThreeFiveBackOff())
	if err != nil {
		logger.Log.Error("SetMetric", zap.String("error while insert in DB", err.Error()))
		return
	}
}

func (st *DBStorage) BeginTransaction(ctx context.Context) (context.Context, error) {
	tx, err := st.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, utils.Transaction, tx)
	return ctx, nil
}

func (st *DBStorage) CommitTransaction(ctx context.Context) error {
	tx, ok := ctx.Value(utils.Transaction).(pgx.Tx)
	if !ok {
		return fmt.Errorf("no transaction found in context")
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (st *DBStorage) RollbackTransaction(ctx context.Context) error {
	tx, ok := ctx.Value(utils.Transaction).(pgx.Tx)
	if !ok {
		return fmt.Errorf("no transaction found in context")
	}
	if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
		return err
	}
	return nil
}

func retriableHelper(err error) error {
	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) {
			if pgerrcode.IsConnectionException(pgErr.Code) ||
				pgerrcode.IsTransactionRollback(pgErr.Code) ||
				pgerrcode.IsInsufficientResources(pgErr.Code) {
				return err
			}
			return backoff.Permanent(err)
		}
	}
	return err
}
