package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/KirillKhitev/metrics/internal/flags"
	"github.com/KirillKhitev/metrics/internal/logger"
	"github.com/KirillKhitev/metrics/internal/metrics"
)

// DBStorage - хранилище в БД
type DBStorage struct {
	db *pgxpool.Pool // Соединение с БД
}

func (s *DBStorage) UpdateCounter(ctx context.Context, name string, value int64) error {
	for i := 1; i <= AttemptCount; i++ {
		err := s.attemptUpdateCounter(ctx, name, value)
		if err != nil {
			var pgErr *pgconn.PgError
			logger.Log.Error("error by update counter", zap.Error(err))
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != AttemptCount {
				time.Sleep(time.Duration(2*i-1) * time.Second)
				continue
			}
			return err
		}

		break
	}

	return nil
}

func (s *DBStorage) attemptUpdateCounter(ctx context.Context, name string, value int64) error {
	oldValue, _ := s.GetCounter(ctx, name)

	value += oldValue

	_, err := s.db.Exec(ctx, "INSERT INTO counter (name, value) VALUES($1, $2) ON CONFLICT (name) DO UPDATE SET name = $1, value = $2", name, value)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil
}

func (s *DBStorage) UpdateCounters(ctx context.Context, data []metrics.Metrics) error {
	for i := 1; i <= AttemptCount; i++ {
		err := s.attemptUpdateCounters(ctx, data)
		if err != nil {
			var pgErr *pgconn.PgError
			logger.Log.Error("error by update counters", zap.Error(err))
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != AttemptCount {
				time.Sleep(time.Duration(2*i-1) * time.Second)
				continue
			}
			return err
		}

		break
	}

	return nil
}

func (s *DBStorage) attemptUpdateCounters(ctx context.Context, data []metrics.Metrics) error {
	if len(data) == 0 {
		return nil
	}

	batch := &pgx.Batch{}

	metricsForUpdate := map[string]int64{}

	for _, metrica := range data {
		if _, ok := metricsForUpdate[metrica.ID]; !ok {
			oldValue, _ := s.GetCounter(ctx, metrica.ID)

			metricsForUpdate[metrica.ID] = oldValue
		}

		metricsForUpdate[metrica.ID] += *metrica.Delta
	}

	for name, value := range metricsForUpdate {
		batch.Queue("INSERT INTO counter (name, value) VALUES($1, $2) ON CONFLICT (name) DO UPDATE SET name = $1, value = $2", name, value)
	}

	results := s.db.SendBatch(ctx, batch)
	defer results.Close()

	for range metricsForUpdate {
		if _, err := results.Exec(); err != nil {
			return fmt.Errorf("unable query: %w", err)
		}
	}

	return nil
}

func (s *DBStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	var value int64
	var err error

	for i := 1; i <= AttemptCount; i++ {
		value, err = s.attemptGetCounter(ctx, name)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			var pgErr *pgconn.PgError
			logger.Log.Error("error by get counter", zap.Error(err))
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != AttemptCount {
				time.Sleep(time.Duration(2*i-1) * time.Second)
				continue
			}
			return value, err
		}

		break
	}

	return value, err
}

func (s *DBStorage) attemptGetCounter(ctx context.Context, name string) (int64, error) {
	row := s.db.QueryRow(ctx, "SELECT value FROM counter WHERE name = $1", name)

	var value int64

	err := row.Scan(&value)
	if err != nil {
		return value, fmt.Errorf("unable query: %w", err)
	}

	return value, nil
}

func (s *DBStorage) GetCounters(ctx context.Context, data []string) (map[string]int64, error) {
	var result map[string]int64
	var err error

	for i := 1; i <= AttemptCount; i++ {
		result, err = s.attemptGetCounters(ctx, data)
		if err != nil {
			var pgErr *pgconn.PgError
			logger.Log.Error("error by get counters", zap.Error(err))
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != AttemptCount {
				time.Sleep(time.Duration(2*i-1) * time.Second)
				continue
			}
			return result, err
		}

		break
	}

	return result, nil
}

func (s *DBStorage) attemptGetCounters(ctx context.Context, data []string) (map[string]int64, error) {
	result := make(map[string]int64)

	rows, err := s.db.Query(ctx, "SELECT name, value FROM counter WHERE name = ANY ($1)", data)
	if err != nil {
		return result, fmt.Errorf("unable query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		var value int64

		err = rows.Scan(&name, &value)
		if err != nil {
			return result, fmt.Errorf("unable to scan row: %w", err)
		}

		result[name] = value
	}

	if err := rows.Err(); err != nil {
		return result, fmt.Errorf("cursor error: %w", err)
	}

	return result, nil
}

func (s *DBStorage) UpdateGauge(ctx context.Context, name string, value float64) error {
	for i := 1; i <= AttemptCount; i++ {
		err := s.attemptUpdateGauge(ctx, name, value)
		if err != nil {
			var pgErr *pgconn.PgError
			logger.Log.Error("error by update gauge", zap.Error(err))
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != AttemptCount {
				time.Sleep(time.Duration(2*i-1) * time.Second)
				continue
			}
			return err
		}

		break
	}

	return nil
}

func (s *DBStorage) attemptUpdateGauge(ctx context.Context, name string, value float64) error {
	_, err := s.db.Exec(ctx, "INSERT INTO gauge (name, value) VALUES($1, $2) ON CONFLICT (name) DO UPDATE SET name = $1, value = $2", name, value)
	if err != nil {
		return fmt.Errorf("unable to query: %w", err)
	}

	return nil
}

func (s *DBStorage) UpdateGauges(ctx context.Context, data []metrics.Metrics) error {
	for i := 1; i <= AttemptCount; i++ {
		err := s.attemptUpdateGauges(ctx, data)
		if err != nil {
			var pgErr *pgconn.PgError
			logger.Log.Error("error by update gauges", zap.Error(err))
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != AttemptCount {
				time.Sleep(time.Duration(2*i-1) * time.Second)
				continue
			}
			return err
		}

		break
	}

	return nil
}

func (s *DBStorage) attemptUpdateGauges(ctx context.Context, data []metrics.Metrics) error {
	if len(data) == 0 {
		return nil
	}
	batch := &pgx.Batch{}

	metricsForUpdate := map[string]float64{}

	for _, metrica := range data {
		metricsForUpdate[metrica.ID] = *metrica.Value
	}

	for name, value := range metricsForUpdate {
		batch.Queue("INSERT INTO gauge (name, value) VALUES($1, $2) ON CONFLICT (name) DO UPDATE SET name = $1, value = $2", name, value)
	}

	results := s.db.SendBatch(ctx, batch)
	defer results.Close()

	for range metricsForUpdate {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("unable query: %w", err)
		}
	}

	return nil
}

func (s *DBStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	var value float64
	var err error

	for i := 1; i <= AttemptCount; i++ {
		value, err = s.attemptGetGauge(ctx, name)
		if err != nil {
			var pgErr *pgconn.PgError
			logger.Log.Error("error by get gauge", zap.Error(err))
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != AttemptCount {
				time.Sleep(time.Duration(2*i-1) * time.Second)
				continue
			}
			return value, err
		}

		break
	}

	return value, nil
}

func (s *DBStorage) attemptGetGauge(ctx context.Context, name string) (float64, error) {
	row := s.db.QueryRow(ctx, "SELECT value FROM gauge WHERE name = $1", name)

	var value float64

	err := row.Scan(&value)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return value, fmt.Errorf("unable query: %w", err)
	}

	return value, nil
}

func (s *DBStorage) GetGauges(ctx context.Context, data []string) (map[string]float64, error) {
	var result map[string]float64
	var err error

	for i := 1; i <= AttemptCount; i++ {
		result, err = s.attemptGetGauges(ctx, data)
		if err != nil {
			var pgErr *pgconn.PgError
			logger.Log.Error("error by get gauges", zap.Error(err))
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != AttemptCount {
				time.Sleep(time.Duration(2*i-1) * time.Second)
				continue
			}
			return result, err
		}

		break
	}

	return result, nil
}

func (s *DBStorage) attemptGetGauges(ctx context.Context, data []string) (map[string]float64, error) {
	result := make(map[string]float64)

	rows, err := s.db.Query(ctx, "SELECT name, value FROM gauge WHERE name = ANY ($1)", data)
	if err != nil {
		return result, fmt.Errorf("unable query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		var value float64

		err = rows.Scan(&name, &value)
		if err != nil {
			return result, fmt.Errorf("unable to scan row: %w", err)
		}

		result[name] = value
	}

	if err := rows.Err(); err != nil {
		return result, fmt.Errorf("cursor error: %w", err)
	}

	return result, nil
}

func (s *DBStorage) Init(ctx context.Context) error {
	db, err := pgxpool.New(ctx, flags.Args.DBConnectionString)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	s.db = db

	for i := 1; i <= AttemptCount; i++ {
		err := s.prepareTables(ctx)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != AttemptCount {
				time.Sleep(time.Duration(2*i-1) * time.Second)
				continue
			}

			return err
		}

		break
	}

	return nil
}

func (s *DBStorage) GetCounterList(ctx context.Context) map[string]int64 {
	var result map[string]int64
	var err error

	for i := 1; i <= AttemptCount; i++ {
		result, err = s.attemptGetCounterList(ctx)
		if err != nil {
			var pgErr *pgconn.PgError
			logger.Log.Error("error by update counters list", zap.Error(err))
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != AttemptCount {
				time.Sleep(time.Duration(2*i-1) * time.Second)
				continue
			}
			return result
		}

		break
	}

	return result
}

func (s *DBStorage) attemptGetCounterList(ctx context.Context) (map[string]int64, error) {
	result := make(map[string]int64)
	rows, err := s.db.Query(ctx, "SELECT name, value FROM counter")
	if err != nil {
		return result, fmt.Errorf("unable query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		var value int64

		errScan := rows.Scan(&name, &value)
		if errScan != nil {
			return result, fmt.Errorf("unable to scan row: %w", errScan)
		}

		result[name] = value
	}

	if errCursor := rows.Err(); errCursor != nil {
		return result, fmt.Errorf("cursor error: %w", errCursor)
	}

	return result, err
}

func (s *DBStorage) GetGaugeList(ctx context.Context) map[string]float64 {
	var result map[string]float64
	var err error

	for i := 1; i <= AttemptCount; i++ {
		result, err = s.attemptGetGaugeList(ctx)
		if err != nil {
			var pgErr *pgconn.PgError
			logger.Log.Error("error by get gauges list", zap.Error(err))
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != AttemptCount {
				time.Sleep(time.Duration(2*i-1) * time.Second)
				continue
			}
			return result
		}

		break
	}

	return result
}

func (s *DBStorage) attemptGetGaugeList(ctx context.Context) (map[string]float64, error) {
	result := make(map[string]float64)
	rows, err := s.db.Query(ctx, "SELECT name, value FROM gauge")
	if err != nil {
		return result, fmt.Errorf("unable to query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		var value float64

		errScan := rows.Scan(&name, &value)
		if errScan != nil {
			return result, fmt.Errorf("unable to scan row: %w", errScan)
		}

		result[name] = value
	}

	if errCursor := rows.Err(); errCursor != nil {
		return result, fmt.Errorf("cursor error: %w", errCursor)
	}

	return result, err
}

func (s *DBStorage) TrySaveToFile() error {
	return nil
}

func (s *DBStorage) Close() error {
	s.db.Close()

	logger.Log.Info("Close Storage")

	return nil
}

func (s *DBStorage) Ping(ctx context.Context) error {
	if err := s.db.Ping(ctx); err != nil {
		return err
	}

	return nil
}

func (s *DBStorage) prepareTables(ctx context.Context) error {
	_, err := s.db.Exec(ctx, "CREATE TABLE IF NOT EXISTS counter (name varchar(255) NOT NULL PRIMARY KEY, value bigint)")
	if err != nil {
		return fmt.Errorf("не удалось создать таблицу counter: %w", err)
	}

	_, err = s.db.Exec(ctx, "CREATE TABLE IF NOT EXISTS gauge (name varchar(255) NOT NULL PRIMARY KEY, value double precision)")
	if err != nil {
		return fmt.Errorf("не удалось создать таблицу gauge: %w", err)
	}

	return nil
}
