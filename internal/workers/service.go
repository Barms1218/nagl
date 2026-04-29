package workers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/rand/v2"
	"strconv"

	"github.com/Barms1218/nagl/internal/contracts"
	"github.com/Barms1218/nagl/internal/database"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const ContractExpiredStream = "contracts:expired"
const WorkerGroupName = "workers"

type WorkerService struct {
	redis           *redis.Client
	store           *database.Store
	contractService *contracts.ContractService
	logger          *slog.Logger
}

func NewWorkerService(
	r *redis.Client,
	s *database.Store,
	c *contracts.ContractService,
	l *slog.Logger) *WorkerService {
	return &WorkerService{
		redis:           r,
		store:           s,
		contractService: c,
		logger:          l,
	}
}

func (s *WorkerService) Start(ctx context.Context, numWorkers int) error {
	err := s.redis.XGroupCreateMkStream(ctx, "contracts:expired", "workers", "0").Err()
	if err != nil {
		if err.Error() != "BUSYGROUP Consumer Group name already exists" {
			return fmt.Errorf("Failed to create consumer group: %w", err)
		}
	}
	for i := range numWorkers {
		go s.runWorkers(ctx, fmt.Sprintf("worker-%d", i))
	}
	return nil
}

func (s *WorkerService) runWorkers(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msgs, err := s.redis.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    WorkerGroupName,
			Consumer: name,
			Streams:  []string{ContractExpiredStream, ">"},
			Count:    1,
			Block:    0,
		}).Result()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			continue
		}

		for _, msg := range msgs[0].Messages {
			if err := s.handleMessage(ctx, msg); err != nil {
				s.logger.Error("Failed to handle message", "error", err, "worker", name)
				continue
			}
			s.redis.XAck(ctx, ContractExpiredStream, WorkerGroupName, msg.ID)
		}
	}
}

func (s *WorkerService) handleMessage(ctx context.Context, msg redis.XMessage) error {
	contractID, err := uuid.Parse(msg.Values["contract_id"].(string))
	if err != nil {
		return fmt.Errorf("Could not parse contract ID: %w", err)
	}
	guildID, err := uuid.Parse(msg.Values["guild_id"].(string))
	if err != nil {
		return fmt.Errorf("Could not parse guild ID: %w", err)
	}
	partyID, err := uuid.Parse(msg.Values["party_id"].(string))
	if err != nil {
		return fmt.Errorf("Could not parse party ID: %w", err)
	}
	difficulty, err := strconv.Atoi(msg.Values["difficulty"].(string))
	if err != nil {
		return fmt.Errorf("Could not parse contract difficulty: %w", err)
	}

	var successChance float64

	party, err := s.store.GetParty(ctx, partyID)
	if err != nil {
		return fmt.Errorf("Could not retrieve party details: %w", err)
	}

	adventurers, err := s.store.GetMemberDetails(ctx, database.UUIDToPgtype(partyID))
	if err != nil {
		return fmt.Errorf("Could not retrieve party member details: %w", err)
	}

	if party.PartyRank >= int32(difficulty) {
		successChance += 0.50
	}

	for _, a := range adventurers {
		delta := float64(a.CurrentRank) - float64(difficulty)
		// normalize delta: clamp contribution per member between -1 and +1
		normalized := math.Max(-1.0, math.Min(1.0, delta/float64(len(adventurers))))
		successChance += (1.0 + normalized) / 2.0
	}

	// Clamp final result to [0, 1]
	successChance = math.Max(0.0, math.Min(1.0, successChance))

	outcome := rand.Float64()

	var status database.ContractStatusEnum
	if outcome > successChance {
		status = database.ContractStatusEnumFailed
	} else {
		status = database.ContractStatusEnumComplete
	}

	if err := s.contractService.SetContractStatus(ctx, contracts.SetContractStatusRequest{
		GuildID:    guildID,
		ContractID: contractID,
		PartyID:    partyID,
		NewStatus:  string(status),
		Difficulty: int32(difficulty),
	}); err != nil {
		return fmt.Errorf("Failed to record contract outcome: %w", err)
	}

	return nil
}
