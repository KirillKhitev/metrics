package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/KirillKhitev/metrics/internal/flags"
	"github.com/KirillKhitev/metrics/internal/logger"
	"github.com/KirillKhitev/metrics/internal/metrics"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"time"
)

type DBStorage struct {
	db *sql.DB
}

func (s *DBStorage) UpdateCounter(ctx context.Context, name string, value int64) error {
	for i := 1; i <= ATTEMPT_COUNT; i++ {
		err := s.attemptUpdateCounter(ctx, name, value)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != ATTEMPT_COUNT {
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

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO counter (name, value) VALUES($1, $2) ON CONFLICT (name) DO UPDATE SET name = $1, value = $2")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, name, value)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *DBStorage) UpdateCounters(ctx context.Context, data []metrics.Metrics) error {
	for i := 1; i <= ATTEMPT_COUNT; i++ {
		err := s.attemptUpdateCounters(ctx, data)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != ATTEMPT_COUNT {
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

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO counter (name, value) VALUES($1, $2) ON CONFLICT (name) DO UPDATE SET name = $1, value = $2")
	if err != nil {
		logger.Log.Error("Error prepare query", zap.Error(err))
		return err
	}

	defer stmt.Close()

	metricsForUpdate := map[string]int64{}

	for _, metrica := range data {
		if _, ok := metricsForUpdate[metrica.ID]; !ok {
			oldValue, _ := s.GetCounter(ctx, metrica.ID)

			metricsForUpdate[metrica.ID] = oldValue
		}

		metricsForUpdate[metrica.ID] += *metrica.Delta
	}

	for name, value := range metricsForUpdate {
		_, err = stmt.ExecContext(ctx, name, value)
		if err != nil {
			logger.Log.Error("Error exec query to DB", zap.Error(err))
			return err
		}
	}

	return tx.Commit()
}

func (s *DBStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	var value int64
	var err error

	for i := 1; i <= ATTEMPT_COUNT; i++ {
		value, err = s.attemptGetCounter(ctx, name)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != ATTEMPT_COUNT {
				time.Sleep(time.Duration(2*i-1) * time.Second)
				continue
			}
			return value, err
		}

		break
	}

	return value, nil
}

func (s *DBStorage) attemptGetCounter(ctx context.Context, name string) (int64, error) {
	row := s.db.QueryRowContext(ctx, "SELECT value FROM counter WHERE name = $1", name)

	var value int64

	err := row.Scan(&value)
	if err != nil {
		return value, err
	}

	return value, nil
}

func (s *DBStorage) GetCounters(ctx context.Context, data []string) (map[string]int64, error) {
	var result map[string]int64
	var err error

	for i := 1; i <= ATTEMPT_COUNT; i++ {
		result, err = s.attemptGetCounters(ctx, data)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != ATTEMPT_COUNT {
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

	stmt, err := s.db.PrepareContext(ctx, "SELECT name, value FROM counter WHERE name = ANY ($1)")
	if err != nil {
		logger.Log.Error("Error prepare query", zap.Error(err))
		return result, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, data)
	if err != nil {
		logger.Log.Error("Error exec query", zap.Error(err))
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		var value int64

		err = rows.Scan(&name, &value)
		if err != nil {
			logger.Log.Error("Error parse values from DB", zap.Error(err))
			return result, err
		}

		result[name] = value
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error("Error query to DB", zap.Error(err))
		return result, err
	}

	return result, nil
}

func (s *DBStorage) UpdateGauge(ctx context.Context, name string, value float64) error {
	for i := 1; i <= ATTEMPT_COUNT; i++ {
		err := s.attemptUpdateGauge(ctx, name, value)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != ATTEMPT_COUNT {
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
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO gauge (name, value) VALUES($1, $2) ON CONFLICT (name) DO UPDATE SET name = $1, value = $2")
	if err != nil {
		logger.Log.Error("Error prepare query", zap.Error(err))
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, name, value)
	if err != nil {
		logger.Log.Error("Error exec query", zap.Error(err))
		return err
	}

	return tx.Commit()
}

func (s *DBStorage) UpdateGauges(ctx context.Context, data []metrics.Metrics) error {
	for i := 1; i <= ATTEMPT_COUNT; i++ {
		err := s.attemptUpdateGauges(ctx, data)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != ATTEMPT_COUNT {
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

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO gauge (name, value) VALUES($1, $2) ON CONFLICT (name) DO UPDATE SET name = $1, value = $2")
	if err != nil {
		logger.Log.Error("Error prepare query", zap.Error(err))
		return err
	}

	defer stmt.Close()

	metricsForUpdate := map[string]float64{}

	for _, metrica := range data {
		metricsForUpdate[metrica.ID] = *metrica.Value
	}

	for name, value := range metricsForUpdate {
		_, err = stmt.ExecContext(ctx, name, value)
		if err != nil {
			logger.Log.Error("Error exec query", zap.Error(err))
			return err
		}
	}

	return tx.Commit()
}

func (s *DBStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	var value float64
	var err error

	for i := 1; i <= ATTEMPT_COUNT; i++ {
		value, err = s.attemptGetGauge(ctx, name)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != ATTEMPT_COUNT {
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
	row := s.db.QueryRowContext(ctx, "SELECT value FROM gauge WHERE name = $1", name)

	var value float64

	err := row.Scan(&value)
	if err != nil {
		return value, err
	}

	return value, nil
}

func (s *DBStorage) GetGauges(ctx context.Context, data []string) (map[string]float64, error) {
	var result map[string]float64
	var err error

	for i := 1; i <= ATTEMPT_COUNT; i++ {
		result, err = s.attemptGetGauges(ctx, data)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != ATTEMPT_COUNT {
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

	stmt, err := s.db.PrepareContext(ctx, "SELECT name, value FROM gauge WHERE name = ANY ($1)")
	if err != nil {
		logger.Log.Error("Error prepare query", zap.Error(err))
		return result, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, data)

	if err != nil {
		logger.Log.Error("Error exec query", zap.Error(err))
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		var value float64

		err = rows.Scan(&name, &value)
		if err != nil {
			logger.Log.Error("Error parse values from DB", zap.Error(err))
			return result, err
		}

		result[name] = value
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error("Error query to DB", zap.Error(err))
		return result, err
	}

	return result, nil
}

func (s *DBStorage) Init() error {
	db, err := sql.Open("pgx", flags.Args.DBConnectionString)
	if err != nil {
		return err
	}

	s.db = db

	for i := 1; i <= ATTEMPT_COUNT; i++ {
		err := s.prepareTables()
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != ATTEMPT_COUNT {
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

	for i := 1; i <= ATTEMPT_COUNT; i++ {
		result, err = s.attemptGetCounterList(ctx)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != ATTEMPT_COUNT {
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
	rows, err := s.db.QueryContext(ctx, "SELECT name, value FROM counter")
	if err != nil {
		logger.Log.Error("Error query to DB", zap.Error(err))
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		var value int64

		err = rows.Scan(&name, &value)
		if err != nil {
			logger.Log.Error("Error parse values from DB", zap.Error(err))
			return result, err
		}

		result[name] = value
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error("Error query to DB", zap.Error(err))
		return result, err
	}

	return result, err
}

func (s *DBStorage) GetGaugeList(ctx context.Context) map[string]float64 {
	var result map[string]float64
	var err error

	for i := 1; i <= ATTEMPT_COUNT; i++ {
		result, err = s.attemptGetGaugeList(ctx)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) && i != ATTEMPT_COUNT {
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
	rows, err := s.db.QueryContext(ctx, "SELECT name, value FROM gauge")
	if err != nil {
		logger.Log.Error("Error query to DB", zap.Error(err))
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		var value float64

		err = rows.Scan(&name, &value)
		if err != nil {
			logger.Log.Error("Error parse values from DB", zap.Error(err))
			return result, err
		}

		result[name] = value
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error("Error query to DB", zap.Error(err))
		return result, err
	}

	return result, err
}

func (s *DBStorage) TrySaveToFile() error {
	return nil
}

func (s *DBStorage) Close() error {
	return s.db.Close()
}

func (s *DBStorage) Ping(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}

func (s *DBStorage) prepareTables() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS counter (name varchar(255) NOT NULL PRIMARY KEY, value bigint)")
	if err != nil {
		return fmt.Errorf("не удалось создать таблицу counter: %w", err)
	}

	_, err = s.db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS gauge (name varchar(255) NOT NULL PRIMARY KEY, value double precision)")
	if err != nil {
		return fmt.Errorf("не удалось создать таблицу gauge: %w", err)
	}

	return nil
}
