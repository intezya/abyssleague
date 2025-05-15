package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
	"github.com/intezya/abyssleague/services/abysscore/internal/pkg/apperrors"
)

func withTxResult[T any](
	ctx context.Context,
	client *ent.Client,
	fn func(tx *ent.Tx) (*T, error),
) (_ *T, err error) {
	ctx, span := tracer.StartSpan(ctx, "withTxResult")
	defer span.End()

	tx, err := client.Tx(ctx)
	if err != nil {
		return nil, apperrors.WrapUnexpectedError(fmt.Errorf("start transaction: %w", err))
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()

			panic(p)
		}

		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = apperrors.WrapUnexpectedError(errors.Join(err, rollbackErr))
			}

			return
		}

		if commitErr := tx.Commit(); commitErr != nil {
			err = apperrors.WrapUnexpectedError(commitErr)
		}
	}()

	var res *T
	res, err = fn(tx)

	return res, err
}

func WithTxResultTx[T any](
	ctx context.Context,
	tx *ent.Tx,
	fn func(tx *ent.Tx) (*T, error),
) (_ *T, err error) {
	ctx, span := tracer.StartSpan(ctx, "WithTxResultTx")
	defer span.End()

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()

			panic(p)
		}

		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = apperrors.WrapUnexpectedError(errors.Join(err, rollbackErr))
			}

			return
		}

		if commitErr := tx.Commit(); commitErr != nil {
			err = apperrors.WrapUnexpectedError(commitErr)
		}
	}()

	var res *T
	res, err = fn(tx)

	return res, err
}

func WithTx(
	ctx context.Context,
	tx *ent.Tx,
	fn func(tx *ent.Tx) error,
) (err error) {
	ctx, span := tracer.StartSpan(ctx, "WithTxResult2Tx")
	defer span.End()

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()

			panic(p)
		}

		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = apperrors.WrapUnexpectedError(errors.Join(err, rollbackErr))
			}

			return
		}

		if commitErr := tx.Commit(); commitErr != nil {
			err = apperrors.WrapUnexpectedError(commitErr)
		}
	}()

	err = fn(tx)

	return err
}
