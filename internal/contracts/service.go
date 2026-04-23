package contracts

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"sync"

	"github.com/Barms1218/nagl/internal/database"
	"github.com/google/uuid"
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
		GuildID: database.UUIDToPgtype(c.GuildID),
	}
	return s.store.AssignToGuild(ctx, params)
}

func (s *ContractService) StartContract(ctx context.Context, c SetContractStatusRequest) error {
	return s.store.ExecTX(ctx, func(q *database.Queries) error {
		contractParams := database.SetContractStatusParams{
			ID:             c.ContractID,
			GuildID:        database.UUIDToPgtype(c.GuildID),
			ContractStatus: database.ContractStatusEnumInProgress,
		}

		if err := q.SetContractStatus(ctx, contractParams); err != nil {
			return err
		}

		memberParams := database.SetMemberStatusParams{
			CurrentActivity: database.ActivityEnumOnContract,
			ContractID:      database.UUIDToPgtype(c.ContractID),
		}

		if err := q.SetMemberStatus(ctx, memberParams); err != nil {
			return err
		}
		return nil
	})
}

func (s *ContractService) ListAvailableContracts(ctx context.Context, sf SearchFilters) ([]ListContractsResponse, error) {
	var params database.ListAvailableContractsParams
	if sf.MinDifficulty != nil {
		params.MinDifficulty = database.IntToPgtype(*sf.MinDifficulty)
	}
	if sf.MaxDifficulty != nil {
		params.MaxDifficulty = database.IntToPgtype(*sf.MaxDifficulty)
	}
	if sf.Status != nil {
		params.Status = database.NullContractStatusEnum{ContractStatusEnum: database.ContractStatusEnum(*sf.Status), Valid: true}
	}
	if sf.PartySize != nil {
		params.PartySize = database.IntToPgtype(*sf.PartySize)
	}
	if sf.SortBy != nil {
		params.SortBy = *sf.SortBy
	} else {
		params.SortBy = "title"
	}

	models, err := s.store.ListAvailableContracts(ctx, params)
	if err != nil {
		return nil, err
	}

	contracts := make([]ListContractsResponse, 0, len(models))
	for _, m := range models {
		c := ListContractsResponse{
			ID:             m.ID,
			Title:          m.Title.String,
			Difficulty:     m.Difficulty,
			RecPartySize:   m.RecPartySize,
			ContractStatus: string(m.ContractStatus),
		}
		contracts = append(contracts, c)
	}

	return contracts, nil
}

func (s *ContractService) ListGuildContracts(ctx context.Context, guildID uuid.UUID, sf SearchFilters) ([]ListContractsResponse, error) {
	var params database.ListGuildContractsParams

	params.GuildID = database.UUIDToPgtype(guildID)
	if sf.MinDifficulty != nil {
		params.MinDifficulty = database.IntToPgtype(*sf.MinDifficulty)
	}
	if sf.MaxDifficulty != nil {
		params.MaxDifficulty = database.IntToPgtype(*sf.MaxDifficulty)
	}
	if sf.Status != nil {
		params.Status = database.NullContractStatusEnum{ContractStatusEnum: database.ContractStatusEnum(*sf.Status), Valid: true}
	}
	if sf.PartySize != nil {
		params.PartySize = database.IntToPgtype(*sf.PartySize)
	}
	if sf.SortBy != nil {
		params.SortBy = *sf.SortBy
	} else {
		params.SortBy = "title"
	}

	models, err := s.store.ListGuildContracts(ctx, params)
	if err != nil {
		return nil, err
	}

	contracts := make([]ListContractsResponse, 0, len(models))
	for _, m := range models {
		c := ListContractsResponse{
			ID:             m.ID,
			Title:          m.Title.String,
			Difficulty:     m.Difficulty,
			RecPartySize:   m.RecPartySize,
			ContractStatus: string(m.ContractStatus),
		}
		contracts = append(contracts, c)
	}

	return contracts, nil
}

func (s *ContractService) GetAvailableContractDetails(ctx context.Context, contractID uuid.UUID) (ContractDetailsResponse, error) {
	model, err := s.store.GetAvailableContractDetails(ctx, contractID)
	if err != nil {
		return ContractDetailsResponse{}, err
	}

	contract := ContractDetailsResponse{
		ID:           model.ID,
		Title:        model.Title.String,
		Description:  model.Description.String,
		Difficulty:   GetDifficultyString(model.Difficulty),
		RecPartySize: model.RecPartySize,
		Status:       string(model.ContractStatus),
	}

	return contract, nil
}

func (s *ContractService) GetActiveContractDetails(ctx context.Context, contractID uuid.UUID) (ActiveContractDetailsResponse, error) {
	model, err := s.store.GetContractDetailsByID(ctx, contractID)
	if err != nil {
		return ActiveContractDetailsResponse{}, err
	}

	contract := ActiveContractDetailsResponse{
		ID:           model.ID,
		Title:        model.Title.String,
		GuildName:    model.GuildName,
		PartyName:    model.PartyName,
		PartyStatus:  string(model.PartyStatus),
		Difficulty:   model.Difficulty,
		RecPartySize: model.RecPartySize,
		Description:  model.Description.String,
		Status:       string(model.ContractStatus),
	}

	return contract, nil
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
	return s.store.GetPartyOnContract(ctx, database.UUIDToPgtype(contractID))
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
				ContractID:      database.UUIDToPgtype(cs.ContractID),
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
		GuildID:        database.UUIDToPgtype(cs.GuildID),
		ID:             cs.ContractID,
		ContractStatus: database.ContractStatusEnum(cs.NewStatus),
	}

	if err := q.SetContractStatus(ctx, contractParams); err != nil {
		return fmt.Errorf("Error occurred during contract status change: %w", err)
	}

	if err := q.InsertContractHistory(ctx, database.InsertContractHistoryParams{
		GuildID:    cs.GuildID,
		ContractID: cs.ContractID,
		PartyID:    cs.PartyID,
		Difficulty: cs.Difficulty,
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
		ContractID:     cs.ContractID,
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
		ContractID: database.UUIDToPgtype(cs.ContractID),
		Status:     database.ContractStatusEnum(cs.NewStatus),
	}

	adventurerIds, err := q.InsertMemberContractHistory(ctx, partyContractParams)
	if err != nil {
		return fmt.Errorf("Error occurred updating party history: %w", err)
	}

	for i := range adventurerIds {
		memberMetrics, err := q.CountMemberCompleteContracts(ctx, adventurerIds[i])
		if err != nil {
			return fmt.Errorf("Error occurred during party member progression: %w", err)
		}
		if memberMetrics.CompletedCount > 0 && memberMetrics.CompletedCount%5 == 0 && memberMetrics.CurrentRank < 5 {
			if err := q.SetAdventurerRank(ctx, database.SetAdventurerRankParams{

				CurrentRank: int32(math.Round(float64(memberMetrics.CompletedCount) / 5.0)),
				ID:          memberMetrics.AdventurerID,
			}); err != nil {
				return fmt.Errorf("Error occurred during party member ranking: %w", err)
			}
		}
	}

	return nil
}

func (s *ContractService) CheckExpiredContracts(ctx context.Context) error {
	var errs []error
	expired, err := s.store.GetExpiredContracts(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("Error occurred when checking expired contracts: %w", err))
	}

	errCh := make(chan error, len(expired))
	var wg sync.WaitGroup

	for _, c := range expired {
		wg.Add(1)
		go func(database.GetExpiredContractsRow) {
			defer wg.Done()
			status, err := s.EvaluateContract(ctx, c)
			if err != nil {
				errCh <- fmt.Errorf("Could not evaluate :colo")
			}

			update := SetContractStatusRequest{
				GuildID:    c.GuildID,
				ContractID: c.ContractID,
				PartyID:    c.PartyID,
				Difficulty: c.Difficulty,
				NewStatus:  status,
			}

			if err := s.SetContractStatus(ctx, update); err != nil {
				errCh <- fmt.Errorf("Could not update contract status: %w", err)
			}
		}(c)

	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (s *ContractService) EvaluateContract(ctx context.Context, contract database.GetExpiredContractsRow) (string, error) {
	var successChance float64
	party, err := s.store.GetPartyOnContract(ctx, database.UUIDToPgtype(contract.ContractID))
	if err != nil {
		return "", fmt.Errorf("Could not retrieve party details: %w", err)
	}

	adventurers, err := s.store.GetMemberDetails(ctx, database.UUIDToPgtype(party.ID))
	if err != nil {
		return "", fmt.Errorf("Could not retrieve party member details: %w", err)
	}

	if party.PartyRank >= contract.Difficulty {
		successChance += 0.50
	}

	// Individual members — split remaining 50% across party
	memberWeight := 0.50 / float64(len(adventurers))
	for _, a := range adventurers {
		delta := float64(a.CurrentRank) - float64(contract.Difficulty)
		// normalize delta: clamp contribution per member between -1 and +1
		normalized := math.Max(-1.0, math.Min(1.0, delta/5.0))
		successChance += memberWeight * (1.0 + normalized) / 2.0
	}

	// Clamp final result to [0, 1]
	successChance = math.Max(0.0, math.Min(1.0, successChance))

	outcome := rand.Float64()
	if outcome <= successChance {
		return string(database.ContractStatusEnumComplete), nil
	}
	return string(database.ContractStatusEnumFailed), nil
}

func GetDifficultyString(difficulty int32) string {
	var diffString string
	switch difficulty {
	case 1:
		diffString = "Common"
	case 2:
		diffString = "Challening"
	case 3:
		diffString = "Dangerous"
	case 4:
		diffString = "Deadly"
	case 5:
		diffString = "Fatal"
	default:
		diffString = ""
	}
	return diffString
}
