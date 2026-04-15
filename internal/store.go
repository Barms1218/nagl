package store

import (
	"context"
	"fmt"
	"github.com/Barms1218/nagl/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
	*database.Queries
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db:      db,
		Queries: database.New(db),
	}
}

func (s *Store) ExecTX(ctx context.Context, fn func(*database.Queries) error) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	q := s.WithTx(tx)

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, erb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
