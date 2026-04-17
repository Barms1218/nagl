package contracts

import (
	"context"
	"fmt"
	"github.com/Barms1218/nagl/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"math"
)

type ContractService struct {
	store *database.Store
}

func NewContractService(s *database.Store) *ContractService {
	return &ContractService{store: s}
}

func (s *ContractService) ClaimContract(ctx context.Context, c ContractClaimRequest) error {
	params := database.AssignToGuildParams{
		ID:      c.ContractID,
		GuildID: c.GuildID,
	}
	return s.store.AssignToGuild(ctx, params)
}

func (s *ContractService) StartContract(ctx context.Context, c SetContractStatusRequest) error {
	return s.store.ExecTX(ctx, func(q *database.Queries) error {
		contractParams := database.SetContractStatusParams{
			ID:             c.ID,
			GuildID:        c.GuildID,
			ContractStatus: database.ContractStatusEnumInProgress,
		}

		if err := q.SetContractStatus(ctx, contractParams); err != nil {
			return fmt.Errorf("Error occurred starting contract: %w", err)
		}

		memberParams := database.SetMemberStatusParams{
			CurrentActivity: database.ActivityEnumOnContract,
			ContractID:      c.ID,
		}

		if err := q.SetMemberStatus(ctx, memberParams); err != nil {
			return fmt.Errorf("Error occurred when adventurers started contract: %w", err)
		}
		return nil
	})
}

func (s *ContractService) ListContracts(ctx context.Context, filter SearchFilters) ([]database.ListContractsRow, error) {
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

	return s.store.ListContracts(ctx, params)

}

func (s *ContractService) GetPastContractsWithStatus(ctx context.Context, p PastContractsParams) ([]database.GetPastContractsWithStatusRow, error) {
	params := database.GetPastContractsWithStatusParams{
		Status: database.ContractStatusEnum(p.Status),
		SortBy: p.SortBy,
	}

	if params.SortBy == "" {
		params.SortBy = "title"
	}

	return s.store.GetPastContractsWithStatus(ctx, params)
}

func (s *ContractService) GetPartyOnContract(ctx context.Context, contractID uuid.UUID) (database.GetPartyOnContractRow, error) {
	return s.store.GetPartyOnContract(ctx, contractID)
}

func (s *ContractService) SetContractStatus(ctx context.Context, cs SetContractStatusRequest) error {
	return s.store.ExecTX(ctx, func(q *database.Queries) error {
		if err := s.RecordContractStatus(ctx, q, cs); err != nil {
			return err
		}

		if err := s.HandlePartyProgression(ctx, q, cs); err != nil {
			return err
		}

		if database.ContractStatusEnum(cs.NewStatus) == database.ContractStatusEnumComplete ||
			database.ContractStatusEnum(cs.NewStatus) == database.ContractStatusEnumFailed {
			partyStatusParams := database.SetMemberStatusParams{
				CurrentActivity: database.ActivityEnumAvailable,
				ContractID:      cs.ID,
			}
			err := q.SetMemberStatus(ctx, partyStatusParams)
			if err != nil {
				return fmt.Errorf("Error occurred updating party  activity: %w", err)
			}
		}

		return nil

	})

}

func (s *ContractService) RecordContractStatus(
	ctx context.Context,
	q *database.Queries,
	cs SetContractStatusRequest) error {
	contractParams := database.SetContractStatusParams{
		GuildID:        cs.GuildID,
		ID:             cs.ID,
		ContractStatus: database.ContractStatusEnum(cs.NewStatus),
	}

	if err := q.SetContractStatus(ctx, contractParams); err != nil {
		return fmt.Errorf("Error occurred during contract status change: %w", err)
	}

	if err := q.InsertContractHistory(ctx, database.InsertContractHistoryParams{
		GuildID:    cs.GuildID,
		ContractID: cs.ID,
		PartyID:    cs.PartyID,
		Status:     database.ContractStatusEnum(cs.NewStatus),
	}); err != nil {
		return fmt.Errorf("Error occurred when updating contract history: %w", err)
	}

	if err := q.RemovePartyFromContract(ctx, cs.PartyID); err != nil {
		return fmt.Errorf("Error occurred removing party from ended contract: %w", err)
	}

	return nil
}

func (s *ContractService) HandlePartyProgression(
	ctx context.Context,
	q *database.Queries,
	cs SetContractStatusRequest) error {

	partyHistoryParams := database.InsertPartyHistoryParams{
		PartyID:        cs.PartyID,
		ContractStatus: database.ContractStatusEnum(cs.NewStatus),
	}

	if err := q.InsertPartyHistory(ctx, partyHistoryParams); err != nil {
		return fmt.Errorf("Error occurred updating party history: %w", err)
	}

	completedContracts, err := q.CountPartyCompleteContracts(ctx, cs.PartyID)
	if err != nil {
		return fmt.Errorf("Error ocurred during party progression: %w", err)
	}

	if completedContracts > 0 && completedContracts%5 == 0 {
		if err := q.SetPartyRank(ctx, database.SetPartyRankParams{
			PartyRank: int32(math.Round(float64(completedContracts) / 5.0)),
			ID:        cs.PartyID,
		}); err != nil {
			return err
		}
	}

	partyContractParams := database.InsertMemberContractHistoryParams{
		ContractID: cs.ID,
		Status:     database.ContractStatusEnum(cs.NewStatus),
	}

	err = q.InsertMemberContractHistory(ctx, partyContractParams)
	if err != nil {
		return fmt.Errorf("Error occurred updating party history: %w", err)
	}

	memberMetrics, err := q.CountMemberCompleteContracts(ctx, cs.ID)
	if err != nil {
		return fmt.Errorf("Error occurred during party member progression: %w", err)
	}

	for _, member := range memberMetrics {
		if member.CompletedCount > 0 && member.CompletedCount%5 == 0 && member.CurrentRank < 5 {
			if err := q.SetAdventurerRank(ctx, database.SetAdventurerRankParams{
				CurrentRank: int32(math.Round(float64(member.CompletedCount) / 5.0)),
				ID:          member.AdventurerID,
			}); err != nil {
				return fmt.Errorf("Error occurred during party member ranking: %w", err)
			}
		}
	}

	return nil
}
