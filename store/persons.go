package store

import (
	"context"

	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/model/uuid"
)

type (
	PersonCollection struct {
		accounts *AccountCollection
	}
)

func (p *PersonCollection) Create(ctx context.Context, accountID uuid.UUID, person model.Person) error {
	account, err := p.accounts.Lookup(ctx, accountID)
	if err != nil {
		return mapError(err)
	}
	account.Person = append(account.Person, person)
	return mapError(p.accounts.Replace(ctx, accountID, *account))
}

func (p *PersonCollection) Lookup(ctx context.Context, accountID uuid.UUID, personID uuid.UUID) (*model.Person, error) {
	account, err := p.accounts.Lookup(ctx, accountID)
	if err != nil {
		return nil, mapError(err)
	}
	for i := range account.Person {
		person := account.Person[i]
		if person.PersonID.Equal(personID) {
			return &person, nil
		}
	}
	return nil, ErrNotFound
}

func (p *PersonCollection) Replace(ctx context.Context, accountID uuid.UUID, personID uuid.UUID, person model.Person) error {
	account, err := p.accounts.Lookup(ctx, accountID)
	if err != nil {
		return mapError(err)
	}
	for i := range account.Person {
		current := account.Person[i]
		if current.PersonID.Equal(personID) {
			account.Person[i] = person
			return mapError(p.accounts.Replace(ctx, accountID, *account))
		}
	}
	return ErrNotFound
}

func (p *PersonCollection) Delete(ctx context.Context, accountID uuid.UUID, personID uuid.UUID) error {
	account, err := p.accounts.Lookup(ctx, accountID)
	if err != nil {
		return mapError(err)
	}
	for i := range account.Person {
		current := account.Person[i]
		if current.PersonID.Equal(personID) {
			var tail []model.Person
			if i+1 < len(account.Person) {
				tail = account.Person[i+1:]
			}
			account.Person = append(account.Person[:i], tail...)
			return mapError(p.accounts.Replace(ctx, accountID, *account))
		}
	}
	return ErrNotFound
}

func (p *PersonCollection) Find(ctx context.Context, option FindPersonOption) ([]model.Person, error) {
	var err error
	var accounts = make([]model.Account, 0, 10)
	if option.AccountID == nil {
		accounts, err = p.accounts.Find(ctx, FindAccountOption{
			PersonFullName: option.PersonFullName,
		})
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
	var persons = make([]model.Person, 0, len(accounts))
	for _, account := range accounts {
		for _, person := range account.Person {
			if option.PersonFullName == nil || checkPersonFullName(person, *option.PersonFullName) {
				persons = append(persons, person)
			}
		}
	}
	return persons, nil
}

func NewPersonMemoryCollection(accounts *AccountCollection) *PersonCollection {
	return &PersonCollection{
		accounts: accounts,
	}
}
