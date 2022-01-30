package business

import (
	"context"
	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/model/uuid"
	"github.com/iv-menshenin/accountant/store"
)

type (
	App struct {
		accounts AccountCollection
		persons  PersonCollection
		objects  ObjectsCollection
	}

	PersonCollection interface {
		Create(context.Context, uuid.UUID, model.Person) error
		Lookup(context.Context, uuid.UUID, uuid.UUID) (*model.Person, error)
		Replace(context.Context, uuid.UUID, uuid.UUID, model.Person) error
		Delete(context.Context, uuid.UUID, uuid.UUID) error
		Find(context.Context, store.FindPersonOption) ([]model.Person, error)
	}

	ObjectsCollection interface {
		Create(context.Context, uuid.UUID, model.Object) error
		Lookup(context.Context, uuid.UUID, uuid.UUID) (*model.Object, error)
		Replace(context.Context, uuid.UUID, uuid.UUID, model.Object) error
		Delete(context.Context, uuid.UUID, uuid.UUID) error
		Find(context.Context, store.FindObjectOption) ([]model.Object, error)
	}

	AccountCollection interface {
		Create(context.Context, model.Account) error
		Lookup(context.Context, uuid.UUID) (*model.Account, error)
		Replace(context.Context, uuid.UUID, model.Account) error
		Delete(context.Context, uuid.UUID) error
		Find(context.Context, store.FindAccountOption) ([]model.Account, error)
	}
)

func New(
	accounts AccountCollection,
	persons PersonCollection,
	objects ObjectsCollection,
) *App {
	return &App{
		accounts: accounts,
		persons:  persons,
		objects:  objects,
	}
}
