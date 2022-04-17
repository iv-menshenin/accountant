package memory

import (
	"context"

	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/storage/internal/memory"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

type (
	AccountCollection struct {
		mapError func(error) error
		mem      *memory.Memory
	}
)

func (a *AccountCollection) Create(ctx context.Context, account model.Account) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return a.mapError(a.mem.Create(account.AccountID, &account))
	}
}

func (a *AccountCollection) Lookup(ctx context.Context, id uuid.UUID) (*model.Account, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		acc, err := a.mem.Lookup(id)
		if err != nil {
			return nil, a.mapError(err)
		}
		return acc.(*model.Account), nil
	}
}

func (a *AccountCollection) Replace(ctx context.Context, id uuid.UUID, account model.Account) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return a.mapError(a.mem.Replace(id, &account))
	}
}

func (a *AccountCollection) Delete(ctx context.Context, id uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return a.mapError(a.mem.Delete(id))
	}
}

func (a *AccountCollection) Find(ctx context.Context, option model.FindAccountOption) ([]model.Account, error) {
	collection := a.mem.Find(func(i interface{}) bool {
		account := i.(*model.Account)
		return checkAccountFilter(*account, option)
	})
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var results = make([]model.Account, 0, len(collection))
		for _, i := range collection {
			results = append(results, *i.(*model.Account))
		}
		return results, nil
	}
}

func NewAccountCollection(mapError func(error) error) *AccountCollection {
	return &AccountCollection{
		mapError: mapError,
		mem:      memory.New(),
	}
}
