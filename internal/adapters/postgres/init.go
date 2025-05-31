package postgres

import (
	storage_config "LostAndFound/internal/config/storage_config"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func NewStorage(dbConfig storage_config.PostgresConfig) (*sql.DB, error) {

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbConfig.Username, dbConfig.Password,
		dbConfig.Host, dbConfig.Port,
		dbConfig.DBName, dbConfig.SSLMode,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(12)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(20 * time.Minute)
	db.SetConnMaxLifetime(10 * time.Minute)

	return db, nil
}
