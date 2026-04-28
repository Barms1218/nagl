package workers

import (
	"context"
	"errors"
	"fmt"
	"github.com/Barms1218/nagl/internal/contracts"
	"github.com/Barms1218/nagl/internal/database"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"math"
	"math/rand/v2"
	"strconv"
)

const ContractExpiredStream = "contracts:expired"
const WorkerGroupName = "workers"

type WorkerService struct {
	redis           *redis.Client
	store           *database.Store
	contractService *contracts.ContractService
}

func NewWorkerService(r *redis.Client, s *database.Store, c *contracts.ContractService) *WorkerService {
	return &WorkerService{
		redis:           r,
		store:           s,
		contractService: c,
	}
}

func (s *WorkerService) Start(ctx context.Context, numWorkers int) {
	for i := range numWorkers {
		go s.runWorkers(ctx, fmt.Sprintf("worker-%d", i))
	}
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
			s.handleMessage(ctx, msg)
			s.redis.XAck(ctx, ContractExpiredStream, WorkerGroupName, msg.ID)
		}
	}
}

func (s *WorkerService) handleMessage(ctx context.Context, msg redis.XMessage) (float64, error) {
	partyID, err := uuid.Parse(msg.Values["party_id"].(string))
	if err != nil {
		return 0, fmt.Errorf("Could not parse party ID: %w", err)
	}
	difficulty, err := strconv.Atoi(msg.Values["difficulty"].(string))
	if err != nil {
		return 0, fmt.Errorf("Could not parse contract difficulty: %w", err)
	}

	var successChance float64

	party, err := s.store.GetParty(ctx, partyID)
	if err != nil {
		return 0, fmt.Errorf("Could not retrieve party details: %w", err)
	}

	adventurers, err := s.store.GetMemberDetails(ctx, database.UUIDToPgtype(partyID))
	if err != nil {
		return 0, fmt.Errorf("Could not retrieve party member details: %w", err)
	}

	if party.PartyRank >= int32(difficulty) {

	}

	for _, a := range adventurers {
		delta := float64(a.CurrentRank) - float64(difficulty)
		// normalize delta: clamp contribution per member between -1 and +1
		normalized := math.Max(-1.0, math.Min(1.0, delta/5.0))
		successChance += (1.0 + normalized) / 2.0
	}

	// Clamp final result to [0, 1]
	successChance = math.Max(0.0, math.Min(1.0, successChance))

	outcome := rand.Float64()

	return outcome, nil
}
