package persistence

import (
	repositoryerrors "abysscore/internal/common/errors/repository"
	"abysscore/internal/infrastructure/ent"
	"context"
	"github.com/intezya/pkglib/logger"
)

func withTx[T any](ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) (T, error)) (res T, err error) {
	tx, err := client.Tx(ctx)
	if err != nil {
		logger.Log.Errorf("ent.WithTx unexpected error: %v", err)
		return res, repositoryerrors.WrapUnexpectedError(err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	return fn(tx)
}
