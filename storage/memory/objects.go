package memory

import (
	"context"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/storage"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

type (
	ObjectCollection struct {
		mapError func(error) error
		accounts *AccountCollection
	}
)

func (o *ObjectCollection) Create(ctx context.Context, accountID uuid.UUID, object domain.Object) error {
	account, err := o.accounts.Lookup(ctx, accountID)
	if err != nil {
		return o.mapError(err)
	}
	account.Objects = append(account.Objects, object)
	return o.mapError(o.accounts.Replace(ctx, accountID, *account))
}

func (o *ObjectCollection) Lookup(ctx context.Context, accountID uuid.UUID, objectID uuid.UUID) (*domain.Object, error) {
	account, err := o.accounts.Lookup(ctx, accountID)
	if err != nil {
		return nil, o.mapError(err)
	}
	for i := range account.Objects {
		object := account.Objects[i]
		if object.ObjectID.Equal(objectID) {
			return &object, nil
		}
	}
	return nil, storage.ErrNotFound
}

func (o *ObjectCollection) Replace(ctx context.Context, accountID uuid.UUID, objectID uuid.UUID, object domain.Object) error {
	account, err := o.accounts.Lookup(ctx, accountID)
	if err != nil {
		return o.mapError(err)
	}
	for i := range account.Objects {
		current := account.Objects[i]
		if current.ObjectID.Equal(objectID) {
			account.Objects[i] = object
			return o.mapError(o.accounts.Replace(ctx, accountID, *account))
		}
	}
	return storage.ErrNotFound
}

func (o *ObjectCollection) Delete(ctx context.Context, accountID uuid.UUID, objectID uuid.UUID) error {
	account, err := o.accounts.Lookup(ctx, accountID)
	if err != nil {
		return o.mapError(err)
	}
	for i := range account.Objects {
		current := account.Objects[i]
		if current.ObjectID.Equal(objectID) {
			var tail []domain.Object
			if i+1 < len(account.Objects) {
				tail = account.Objects[i+1:]
			}
			account.Objects = append(account.Objects[:i], tail...)
			return o.mapError(o.accounts.Replace(ctx, accountID, *account))
		}
	}
	return storage.ErrNotFound
}

func (o *ObjectCollection) Find(ctx context.Context, option storage.FindObjectOption) ([]domain.Object, error) {
	var err error
	var accounts = make([]domain.Account, 0, 10)
	if option.AccountID == nil {
		accounts, err = o.accounts.Find(ctx, storage.FindAccountOption{}) // find all ?
	} else {
		var account *domain.Account
		account, err = o.accounts.Lookup(ctx, *option.AccountID)
		if account != nil {
			accounts = append(accounts, *account)
		}
	}
	if err != nil {
		return nil, o.mapError(err)
	}
	var objects = make([]domain.Object, 0, len(accounts))
	for _, account := range accounts {
		for _, object := range account.Objects {
			if checkObjectFilter(object, option) {
				objects = append(objects, object)
			}
		}
	}
	return objects, nil
}

func NewObjectCollection(accounts *AccountCollection, mapError func(error) error) *ObjectCollection {
	return &ObjectCollection{
		mapError: mapError,
		accounts: accounts,
	}
}
