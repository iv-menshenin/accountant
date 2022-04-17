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

func (a *Acc) PersonsFind(ctx context.Context, q request.FindPersonsQuery) ([]domain.Person, error) {
	var findOption storage.FindPersonOption
	findOption.FillFromQuery(q)
	persons, err := a.persons.Find(ctx, findOption)
	if err != nil {
		if err == storage.ErrNotFound {
			return nil, generic.NotFound{}
		}
		return nil, err
	}
	if len(persons) == 0 {
		return nil, generic.NotFound{}
	}
	return persons, nil
}
