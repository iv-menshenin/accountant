package mongodb

import (
	"context"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/storage"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

type (
	PersonsCollection struct {
		mapError func(error) error
		accounts *AccountsCollection
	}
)

func (p *PersonsCollection) Create(ctx context.Context, accountID uuid.UUID, person domain.Person) error {
	account, err := p.accounts.Lookup(ctx, accountID)
	if err != nil {
		return p.mapError(err)
	}
	account.Persons = append(account.Persons, person)
	return p.mapError(p.accounts.Replace(ctx, accountID, *account))
}

func (p *PersonsCollection) Lookup(ctx context.Context, accountID uuid.UUID, personID uuid.UUID) (*domain.Person, error) {
	account, err := p.accounts.Lookup(ctx, accountID)
	if err != nil {
		return nil, p.mapError(err)
	}
	for i := range account.Persons {
		person := account.Persons[i]
		if person.PersonID.Equal(personID) {
			return &person, nil
		}
	}
	return nil, storage.ErrNotFound
}

func (p *PersonsCollection) Replace(ctx context.Context, accountID uuid.UUID, personID uuid.UUID, person domain.Person) error {
	account, err := p.accounts.Lookup(ctx, accountID)
	if err != nil {
		return p.mapError(err)
	}
	for i := range account.Persons {
		current := account.Persons[i]
		if current.PersonID.Equal(personID) {
			account.Persons[i] = person
			return p.mapError(p.accounts.Replace(ctx, accountID, *account))
		}
	}
	return storage.ErrNotFound
}

func (p *PersonsCollection) Delete(ctx context.Context, accountID uuid.UUID, personID uuid.UUID) error {
	account, err := p.accounts.Lookup(ctx, accountID)
	if err != nil {
		return p.mapError(err)
	}
	for i := range account.Persons {
		current := account.Persons[i]
		if current.PersonID.Equal(personID) {
			var tail []domain.Person
			if i+1 < len(account.Persons) {
				tail = account.Persons[i+1:]
			}
			account.Persons = append(account.Persons[:i], tail...)
			return p.mapError(p.accounts.Replace(ctx, accountID, *account))
		}
	}
	return storage.ErrNotFound
}

func (p *PersonsCollection) Find(ctx context.Context, option storage.FindPersonOption) ([]domain.Person, error) {
	var err error
	var accounts = make([]domain.Account, 0, 10)
	if option.AccountID != nil {
		var account *domain.Account
		account, err = p.accounts.Lookup(ctx, *option.AccountID)
		if account != nil {
			accounts = append(accounts, *account)
		}
	}
	if err != nil {
		return nil, p.mapError(err)
	}
	var persons = make([]domain.Person, 0, len(accounts))
	for _, account := range accounts {
		persons = append(persons, account.Persons...)
	}
	return persons, nil
}

func (s *Storage) NewPersonsCollection(accounts *AccountsCollection, mapError func(error) error) *PersonsCollection {
	return &PersonsCollection{
		mapError: mapError,
		accounts: accounts,
	}
}
