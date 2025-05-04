package persistence

import (
	"context"

	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/mapper"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent/bannedhardwareid"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
	"github.com/intezya/abyssleague/services/abysscore/pkg/optional"
)

type BannedHardwareIDRepository struct {
	client *ent.Client
}

func NewBannedHardwareIDRepository(client *ent.Client) *BannedHardwareIDRepository {
	return &BannedHardwareIDRepository{client: client}
}

func (r *BannedHardwareIDRepository) Create(
	ctx context.Context,
	hardwareID string,
	reason optional.String,
) (*dto.BannedHardwareID, error) {
	ctx, span := tracer.StartSpan(ctx, "BannedHardwareIDRepository.Create")
	defer span.End()

	// TODO implement me
	panic("implement me")
}

func (r *BannedHardwareIDRepository) FindByHardwareID(
	ctx context.Context,
	hardwareID string,
) (*dto.BannedHardwareID, error) {
	ctx, span := tracer.StartSpan(ctx, "BannedHardwareIDRepository.FindByHardwareID")
	defer span.End()

	// TODO implement me
	panic("implement me")
}

func (r *BannedHardwareIDRepository) DeleteByHardwareID(
	ctx context.Context,
	hardwareID string,
) error {
	ctx, span := tracer.StartSpan(ctx, "BannedHardwareIDRepository.DeleteByHardwareID")
	defer span.End()

	// TODO implement me
	panic("implement me")
}

func (r *BannedHardwareIDRepository) TxFindByHardwareID(
	ctx context.Context,
	tx *ent.Tx,
	hardwareID string,
) (*dto.BannedHardwareID, error) {
	ctx, span := tracer.StartSpan(ctx, "BannedHardwareIDRepository.TxFindByHardwareID")
	defer span.End()

	found, err := tx.BannedHardwareID.
		Query().
		Where(bannedhardwareid.HardwareIDEQ(hardwareID)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return mapper.ToBannedHardwareIDFromEnt(found), nil
}
