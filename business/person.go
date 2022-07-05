package business

import (
	"context"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/model/generic"
	"github.com/iv-menshenin/accountant/model/request"
	"github.com/iv-menshenin/accountant/storage"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

func (a *Acc) PersonCreate(ctx context.Context, q request.PostPersonQuery) (*domain.Person, error) {
	var person = domain.Person{
		PersonID:   uuid.NewUUID(),
		PersonData: q.PersonData,
	}
	err := a.persons.Create(ctx, q.AccountID, person)
	if err != nil {
		return nil, err
	}
	return a.persons.Lookup(ctx, q.AccountID, person.PersonID)
}

func (a *Acc) PersonGet(ctx context.Context, q request.GetPersonQuery) (*domain.Person, error) {
	person, err := a.persons.Lookup(ctx, q.AccountID, q.PersonID)
	if err == storage.ErrNotFound {
		return nil, generic.NotFound{}
	}
	return person, nil
}

func (a *Acc) PersonSave(ctx context.Context, q request.PutPersonQuery) (*domain.Person, error) {
	person, err := a.persons.Lookup(ctx, q.AccountID, q.PersonID)
	if err == storage.ErrNotFound {
		return nil, generic.NotFound{}
	}
	person.PersonData = q.PersonData
	if err = a.persons.Replace(ctx, q.AccountID, q.PersonID, *person); err != nil {
		return nil, err
	}
	return person, nil
}

func (a *Acc) PersonDelete(ctx context.Context, q request.DeletePersonQuery) error {
	err := a.persons.Delete(ctx, q.AccountID, q.PersonID)
	if err == storage.ErrNotFound {
		return generic.NotFound{}
	}
	return err
}

func (a *Acc) PersonsFind(ctx context.Context, q request.FindPersonsQuery) ([]domain.NestedPerson, error) {
	var accounts []domain.Account
	if q.AccountID != nil {
		acc, err := a.accounts.Lookup(ctx, *q.AccountID)
		if err != nil {
			if err == storage.ErrNotFound {
				err = generic.NotFound{}
			}
			return nil, err
		}
		accounts = append(accounts, *acc)
	} else {
		var findOpts string
		if q.PersonFullName != nil {
			findOpts = *q.PersonFullName
		}
		accs, err := a.accounts.Find(ctx, storage.FindAccountOption{
			ByPerson: findOpts,
		})
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, accs...)
	}
	var persons = make([]domain.NestedPerson, 0, len(accounts))
	for _, acc := range accounts {
		for _, person := range acc.Persons {
			persons = append(persons, domain.NestedPerson{
				Person:    person,
				AccountID: acc.AccountID,
			})
		}
	}
	return persons, nil
}
