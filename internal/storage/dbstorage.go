package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/KirillKhitev/metrics/internal/flags"
	"github.com/KirillKhitev/metrics/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"time"
)

type DBStorage struct {
	db *sql.DB
}

func (s *DBStorage) UpdateCounter(ctx context.Context, name string, value int64) error {
	oldValue, _ := s.GetCounter(ctx, name)

	value += oldValue

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO counter (name, value) VALUES($1, $2) ON CONFLICT (name) DO UPDATE SET name = $1, value = $2")
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

func (s *DBStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	row := s.db.QueryRowContext(ctx, "SELECT value FROM counter WHERE name = $1", name)

	var value int64

	err := row.Scan(&value)
	if err != nil {
		return value, err
	}

	return value, nil
}

func (s *DBStorage) UpdateGauge(ctx context.Context, name string, value float64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO gauge (name, value) VALUES($1, $2) ON CONFLICT (name) DO UPDATE SET name = $1, value = $2")
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

func (s *DBStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	row := s.db.QueryRowContext(ctx, "SELECT value FROM gauge WHERE name = $1", name)

	var value float64

	err := row.Scan(&value)
	if err != nil {
		return value, err
	}

	return value, nil
}

func (s *DBStorage) Init() error {
	db, err := sql.Open("pgx", flags.Args.DBConnectionString)
	if err != nil {
		return err
	}

	s.db = db

	s.prepareTables()

	return nil
}

func (s *DBStorage) GetCounterList(ctx context.Context) map[string]int64 {
	result := make(map[string]int64)
	rows, err := s.db.QueryContext(ctx, "SELECT name, value FROM counter")
	if err != nil {
		logger.Log.Error("Error query to DB", zap.Error(err))
		return result
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		var value int64

		err = rows.Scan(&name, &value)
		if err != nil {
			logger.Log.Error("Error parse values from DB", zap.Error(err))
			return result
		}

		result[name] = value
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error("Error query to DB", zap.Error(err))
		return result
	}

	return result
}

func (s *DBStorage) GetGaugeList(ctx context.Context) map[string]float64 {
	result := make(map[string]float64)
	rows, err := s.db.QueryContext(ctx, "SELECT name, value FROM gauge")
	if err != nil {
		logger.Log.Error("Error query to DB", zap.Error(err))
		return result
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		var value float64

		err = rows.Scan(&name, &value)
		if err != nil {
			logger.Log.Error("Error parse values from DB", zap.Error(err))
			return result
		}

		result[name] = value
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error("Error query to DB", zap.Error(err))
		return result
	}

	return result
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
