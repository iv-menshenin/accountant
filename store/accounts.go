package store

import (
	"context"
	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/model/uuid"
	"github.com/iv-menshenin/accountant/store/internal/memory"
)

type (
	AccountCollection struct {
		Create  func(context.Context, model.Account) error
		Lookup  func(context.Context, uuid.UUID) (*model.Account, error)
		Replace func(context.Context, uuid.UUID, model.Account) error
		Delete  func(context.Context, uuid.UUID) error
		Find    func(context.Context, FindAccountOption) ([]model.Account, error)
	}
)

func NewAccountMemoryCollection() *AccountCollection {
	var accountMem = memory.New()
	return &AccountCollection{
		Create: func(ctx context.Context, account model.Account) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return mapError(accountMem.Create(account.AccountID, &account))
			}
		},

		Lookup: func(ctx context.Context, id uuid.UUID) (*model.Account, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				acc, err := accountMem.Lookup(id)
				if err != nil {
					return nil, mapError(err)
				}
				return acc.(*model.Account), nil
			}
		},

		Replace: func(ctx context.Context, id uuid.UUID, account model.Account) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return mapError(accountMem.Replace(id, &account))
			}
		},

		Delete: func(ctx context.Context, id uuid.UUID) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return mapError(accountMem.Delete(id))
			}
		},

		Find: func(ctx context.Context, option FindAccountOption) ([]model.Account, error) {
			collection := accountMem.Find(func(i interface{}) bool {
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
		},
	}
}
