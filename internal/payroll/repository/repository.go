package repository

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
	pool *pgxpool.Pool
	log  *zap.Logger
}

func NewRepository(pool *pgxpool.Pool, log *zap.Logger) *Repository {
	return &Repository{pool: pool, log: log}
}

func repoErr(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("payroll_repo.%s: %w", op, err)
}
