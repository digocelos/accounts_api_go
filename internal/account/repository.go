package account

import "context"

type Repository interface {
	InsertIdempotentByDocument(ctx context.Context, a Account) (inserted bool, err error)
	GetById(ctx context.Context, id string) (Account, error)
	GetByDocument(ctx context.Context, document string) (Account, error)
	UpdateWithOptimisticLock(ctx context.Context, id string, expectedVersion int, name *string, email *string) (Account, error)
}
