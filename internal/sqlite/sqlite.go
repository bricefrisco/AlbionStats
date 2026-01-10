package sqlite

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	db *gorm.DB
}

func NewSQLiteDatabase(dsn string) (*SQLite, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.AutoMigrate(
		&PlayerPoll{},
		&PlayerStats{},
	); err != nil {
		return nil, fmt.Errorf("db migrate: %w", err)
	}

	if err := createIndexes(db); err != nil {
		return nil, fmt.Errorf("create indexes: %w", err)
	}

	return &SQLite{
		db: db,
	}, nil
}
