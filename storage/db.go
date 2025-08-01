package storage

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	DB *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *DB {
	return &DB{
		DB: db,
	}
}

type Queries interface {
	SearchCards()
}
