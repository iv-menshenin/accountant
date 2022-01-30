package store

import (
	"context"

	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/model/uuid"
)

type (
	ObjectCollection struct {
		accounts *AccountCollection
	}
)

func (p *ObjectCollection) Create(ctx context.Context, accountID uuid.UUID, object model.Object) error {
	account, err := p.accounts.Lookup(ctx, accountID)
	if err != nil {
		return mapError(err)
	}
	account.Object = append(account.Object, object)
	return mapError(p.accounts.Replace(ctx, accountID, *account))
}

func (p *ObjectCollection) Lookup(ctx context.Context, accountID uuid.UUID, objectID uuid.UUID) (*model.Object, error) {
	account, err := p.accounts.Lookup(ctx, accountID)
	if err != nil {
		return nil, mapError(err)
	}
	for i := range account.Object {
		object := account.Object[i]
		if object.ObjectID.Equal(objectID) {
			return &object, nil
		}
	}
	return nil, ErrNotFound
}

func (p *ObjectCollection) Replace(ctx context.Context, accountID uuid.UUID, objectID uuid.UUID, object model.Object) error {
	account, err := p.accounts.Lookup(ctx, accountID)
	if err != nil {
		return mapError(err)
	}
	for i := range account.Object {
		current := account.Object[i]
		if current.ObjectID.Equal(objectID) {
			account.Object[i] = object
			return mapError(p.accounts.Replace(ctx, accountID, *account))
		}
	}
	return ErrNotFound
}

func (p *ObjectCollection) Delete(ctx context.Context, accountID uuid.UUID, objectID uuid.UUID) error {
	account, err := p.accounts.Lookup(ctx, accountID)
	if err != nil {
		return mapError(err)
	}
	for i := range account.Object {
		current := account.Object[i]
		if current.ObjectID.Equal(objectID) {
			var tail []model.Object
			if i+1 < len(account.Object) {
				tail = account.Object[i+1:]
			}
			account.Object = append(account.Object[:i], tail...)
			return mapError(p.accounts.Replace(ctx, accountID, *account))
		}
	}
	return ErrNotFound
}

func (p *ObjectCollection) Find(ctx context.Context, option FindObjectOption) ([]model.Object, error) {
	var err error
	var accounts = make([]model.Account, 0, 10)
	if option.AccountID == nil {
		accounts, err = p.accounts.Find(ctx, FindAccountOption{}) // find all ?
	} else {
		var account *model.Account
		account, err = p.accounts.Lookup(ctx, *option.AccountID)
		if account != nil {
			accounts = append(accounts, *account)
		}
	}
	if err != nil {
		return nil, mapError(err)
	}
	var objects = make([]model.Object, 0, len(accounts))
	for _, account := range accounts {
		for _, object := range account.Object {
			if checkObjectFilter(object, option) {
				objects = append(objects, object)
			}
		}
	}
	return objects, nil
}

func NewObjectMemoryCollection(accounts *AccountCollection) *ObjectCollection {
	return &ObjectCollection{
		accounts: accounts,
	}
}
