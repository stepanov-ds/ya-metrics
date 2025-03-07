package storage

import (
	"context"

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

		TABLESPACE pg_default;

		ALTER TABLE IF EXISTS public.metrics
    		OWNER to usr;
	`
	_, err := p.Exec(context.Background(), query)
	if err != nil {
		logger.Log.Error("NewDbStorage", zap.String("error while creating table in DB", err.Error()))
	}

	return &DbStorage{
		Pool: p,
	}
}

func (st *DbStorage) GetMetric(key string) (utils.Metrics, bool) {
	query := "SELECT ID, MType, Delta, Value FROM public.metrics WHERE ID = $1;"

	row := st.Pool.QueryRow(context.Background(), query, key)

	var m utils.Metrics
	if err := row.Scan(&m.ID, &m.MType, &m.Delta, &m.Value); err != nil {
		if err != pgx.ErrNoRows {
			logger.Log.Error("GetMetric", zap.String("error while parsing row result", err.Error()))
		}
		return m, false
	}
	return m, true
}

func (st *DbStorage) GetAllMetrics() map[string]utils.Metrics {
	query := "SELECT ID, MType, Delta, Value FROM public.metrics;"

	rows, err := st.Pool.Query(context.Background(), query)
	if err != nil {
		logger.Log.Error("GetAllMetric", zap.String("error while select from DB", err.Error()))
	}
	defer rows.Close()

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
			"Delta" = public.metrics."Delta" + EXCLUDED."Delta",
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
	if counter {
		_, err := st.Pool.Exec(context.Background(), query1, key, "counter", value)
		if err != nil {
			logger.Log.Error("SetMetric", zap.String("error while insert in DB", err.Error()))
			return
		}
	} else {
		_, err := st.Pool.Exec(context.Background(), query2, key, "gauge", value)
		if err != nil {
			logger.Log.Error("SetMetric", zap.String("error while insert in DB", err.Error()))
			return
		}
	}

}
