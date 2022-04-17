package business

import (
	"context"

	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/storage"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

func (a *Acc) PersonCreate(ctx context.Context, q model.PostPersonQuery) (*model.Person, error) {
	var person = model.Person{
		PersonID:   uuid.NewUUID(),
		PersonData: q.PersonData,
	}
	err := a.persons.Create(ctx, q.AccountID, person)
	if err != nil {
		return nil, err
	}
	return a.persons.Lookup(ctx, q.AccountID, person.PersonID)
}

func (a *Acc) PersonGet(ctx context.Context, q model.GetPersonQuery) (*model.Person, error) {
	person, err := a.persons.Lookup(ctx, q.AccountID, q.PersonID)
	if err == storage.ErrNotFound {
		return nil, model.NotFound{}
	}
	return person, nil
}

func (a *Acc) PersonSave(ctx context.Context, q model.PutPersonQuery) (*model.Person, error) {
	person, err := a.persons.Lookup(ctx, q.AccountID, q.PersonID)
	if err == storage.ErrNotFound {
		return nil, model.NotFound{}
	}
	person.PersonData = q.PersonData
	if err = a.persons.Replace(ctx, q.AccountID, q.PersonID, *person); err != nil {
		return nil, err
	}
	return person, nil
}

func (a *Acc) PersonDelete(ctx context.Context, q model.DeletePersonQuery) error {
	err := a.persons.Delete(ctx, q.AccountID, q.PersonID)
	if err == storage.ErrNotFound {
		return model.NotFound{}
	}
	return err
}

func (a *Acc) PersonsFind(ctx context.Context, q model.FindPersonsQuery) ([]model.Person, error) {
	var findOption model.FindPersonOption
	findOption.FillFromQuery(q)
	persons, err := a.persons.Find(ctx, findOption)
	if err != nil {
		if err == storage.ErrNotFound {
			return nil, model.NotFound{}
		}
		return nil, err
	}
	if len(persons) == 0 {
		return nil, model.NotFound{}
	}
	return persons, nil
}
