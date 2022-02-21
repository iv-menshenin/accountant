package business

import (
	"context"

	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/model/uuid"
)

type (
	App struct {
		logger   Logger
		accounts AccountCollection
		persons  PersonCollection
		objects  ObjectsCollection
	}
	Logger interface {
		Warning(format string, args ...interface{})
		Debug(format string, args ...interface{})
		Error(format string, args ...interface{})
	}

	PersonCollection interface {
		Create(context.Context, uuid.UUID, model.Person) error
		Lookup(context.Context, uuid.UUID, uuid.UUID) (*model.Person, error)
		Replace(context.Context, uuid.UUID, uuid.UUID, model.Person) error
		Delete(context.Context, uuid.UUID, uuid.UUID) error
		Find(context.Context, model.FindPersonOption) ([]model.Person, error)
	}

	ObjectsCollection interface {
		Create(context.Context, uuid.UUID, model.Object) error
		Lookup(context.Context, uuid.UUID, uuid.UUID) (*model.Object, error)
		Replace(context.Context, uuid.UUID, uuid.UUID, model.Object) error
		Delete(context.Context, uuid.UUID, uuid.UUID) error
		Find(context.Context, model.FindObjectOption) ([]model.Object, error)
	}

	AccountCollection interface {
		Create(context.Context, model.Account) error
		Lookup(context.Context, uuid.UUID) (*model.Account, error)
		Replace(context.Context, uuid.UUID, model.Account) error
		Delete(context.Context, uuid.UUID) error
		Find(context.Context, model.FindAccountOption) ([]model.Account, error)
	}
)

func New(
	logger Logger,
	accounts AccountCollection,
	persons PersonCollection,
	objects ObjectsCollection,
) *App {
	return &App{
		logger:   logger,
		accounts: accounts,
		persons:  persons,
		objects:  objects,
	}
}
