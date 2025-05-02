package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent"
	"github.com/intezya/abyssleague/services/abysscore/internal/pkg/apperrors"
)

func withTxResult[T any](
	ctx context.Context,
	client *ent.Client,
	fn func(tx *ent.Tx) (*T, error),
) (_ *T, err error) {
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
