package store

import (
	"context"
	"fmt"

	"github.com/Barms1218/nagl/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

// type ContractStore interface {
// 	GetContracts(ctx context.Context, sortBy string) ([]database.GetContractsRow, error)
// 	GetContractsWithDifficulty(ctx context.Context, arg database.GetContractsWithDifficultyParams) ([]database.GetContractsWithDifficultyRow, error)
// 	GetContractsWithMinPartySize(ctx context.Context, arg database.GetContractsWithMinPartySizeParams) ([]database.GetContractsWithMinPartySizeRow, error)
// 	GetContractsWithStatus(ctx context.Context, arg database.GetContractsWithStatusParams) ([]database.GetContractsWithStatusRow, error)
// 	GetPartyOnContract(ctx context.Context, id uuid.UUID) ([]database.GetPartyOnContractRow, error)
// 	InsertContract(ctx context.Context, arg database.InsertContractParams) (database.Contract, error)
// 	InsertContractHistory(ctx context.Context, arg database.InsertContractHistoryParams) error
// 	SetContractStatus(ctx context.Context, arg database.SetContractStatusParams) error
// }
//
// func NewContractQueries(s ContractStore) *ContractQueries {
// 	return &ContractQueries{store: s}
// }
//
// type ContractQueries struct {
// 	store ContractStore
// }

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
