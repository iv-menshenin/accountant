package mongodb

import (
	"context"
	"fmt"
	"strings"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/storage"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

type (
	ObjectsCollection struct {
		mapError func(error) error
		accounts *AccountsCollection
	}
)

func (o *ObjectsCollection) Create(ctx context.Context, accountID uuid.UUID, object domain.Object) error {
	account, err := o.accounts.Lookup(ctx, accountID)
	if err != nil {
		return o.mapError(err)
	}
	account.Objects = append(account.Objects, object)
	return o.mapError(o.accounts.Replace(ctx, accountID, *account))
}

func (o *ObjectsCollection) Lookup(ctx context.Context, accountID uuid.UUID, objectID uuid.UUID) (*domain.Object, error) {
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

func (o *ObjectsCollection) Replace(ctx context.Context, accountID uuid.UUID, objectID uuid.UUID, object domain.Object) error {
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

func (o *ObjectsCollection) Delete(ctx context.Context, accountID uuid.UUID, objectID uuid.UUID) error {
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

func (o *ObjectsCollection) Find(ctx context.Context, option storage.FindObjectOption) ([]domain.NestedObject, error) {
	var err error
	var accounts = make([]domain.Account, 0, 10)
	if option.AccountID == nil {
		accounts, err = o.accounts.Find(ctx, findAccountByObject(option))
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
	var objects = make([]domain.NestedObject, 0, len(accounts))
	for _, account := range accounts {
		for _, object := range account.Objects {
			if checkObjectFilter(object, option) {
				objects = append(objects, domain.NestedObject{
					Object:    object,
					AccountID: account.AccountID,
				})
			}
		}
	}
	return objects, nil
}

func findAccountByObject(filter storage.FindObjectOption) storage.FindAccountOption {
	return storage.FindAccountOption{
		Address: filter.Address,
		Number:  filter.Number,
	}
}

func checkObjectFilter(object domain.Object, filter storage.FindObjectOption) bool {
	if filter.Address != nil {
		if !checkObjectAddress(object, *filter.Address) {
			return false
		}
	}
	if filter.Number != nil {
		if !checkObjectNumber(object, *filter.Number) {
			return false
		}
	}
	return true
}

func checkObjectAddress(object domain.Object, address string) bool {
	var addressA = strings.ToUpper(address)
	var addressB = strings.ToUpper(fmt.Sprintf("%s %s %s", object.City, object.Village, object.Street))
	return strings.Contains(addressB, addressA)
}

func checkObjectNumber(object domain.Object, number int) bool {
	return object.Number == number
}

func (s *Storage) NewObjectsCollection(accounts *AccountsCollection, mapError func(error) error) *ObjectsCollection {
	return &ObjectsCollection{
		mapError: mapError,
		accounts: accounts,
	}
}
