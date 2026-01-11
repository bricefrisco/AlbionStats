package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db *gorm.DB
}

func NewPostgresDatabase(dsn string) (*Postgres, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return &Postgres{
		db: db,
	}, nil
}
