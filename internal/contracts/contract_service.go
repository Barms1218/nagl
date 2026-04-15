package contracts

import (
	"context"

	"github.com/Barms1218/nagl/internal"
	"github.com/Barms1218/nagl/internal/database"
	"github.com/jackc/pgx/v5/pgtype"
)

type ContractService struct {
	store *store.Store
}

func NewContractService(s *store.Store) *ContractService {
	return &ContractService{store: s}
}

func (s *ContractService) ListContracts(ctx context.Context, filter SearchFilters) error {
	params := database.ListContractsParams{
		SortBy: filter.SortBy,
	}

	if filter.Difficulty != nil {
		params.Difficulty = pgtype.Int4{Int32: int32(*filter.Difficulty), Valid: true}
	}
	if filter.PartySize != nil {
		params.MinPartySize = pgtype.Int4{Int32: int32(*filter.PartySize), Valid: true}
	}
	if filter.Status != nil {
		params.Status = database.NullContractStatusEnum{
			ContractStatusEnum: database.ContractStatusEnum(*filter.Status),
			Valid:              true,
		}
	}

	if params.SortBy == "" {
		params.SortBy = "title"
	}

	return s.ListContracts(ctx, params)

}
