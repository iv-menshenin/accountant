package business

import (
	"context"

	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

type (
	Tar struct {
		targets   TargetsCollection
		getLogger func() Logger
	}
	Acc struct {
		accounts  AccountsCollection
		persons   PersonsCollection
		objects   ObjectsCollection
		getLogger func() Logger
	}
	App struct {
		Acc
		Tar
	}
	Logger interface {
		Warning(format string, args ...interface{})
		Debug(format string, args ...interface{})
		Error(format string, args ...interface{})
	}

	PersonsCollection interface {
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

	AccountsCollection interface {
		Create(context.Context, model.Account) error
		Lookup(context.Context, uuid.UUID) (*model.Account, error)
		Replace(context.Context, uuid.UUID, model.Account) error
		Delete(context.Context, uuid.UUID) error
		Find(context.Context, model.FindAccountOption) ([]model.Account, error)
	}

	PaymentsCollection interface {
		Create(context.Context, model.Payment) error
		Delete(context.Context, uuid.UUID) error
		FindByAccount(context.Context, uuid.UUID) ([]model.Payment, error)
		FindByIDs(context.Context, []uuid.UUID) ([]model.Payment, error)
	}

	TargetsCollection interface {
		Create(context.Context, model.Target) error
		Lookup(context.Context, uuid.UUID) (*model.Target, error)
		Delete(context.Context, uuid.UUID) error
		FindByPeriod(context.Context, model.FindTargetOption) ([]model.Target, error)
	}

	BillsCollection interface {
		Create(context.Context, model.Bill) error
		Lookup(context.Context, uuid.UUID) (*model.Bill, error)
		FindByAccount(context.Context, uuid.UUID) ([]model.Bill, error)
		FindByPeriod(context.Context, model.Period) ([]model.Bill, error)
		Delete(context.Context, uuid.UUID) error
	}
)

func New(
	logger Logger,
	accounts AccountsCollection,
	persons PersonsCollection,
	objects ObjectsCollection,
	targets TargetsCollection,
) *App {
	return &App{
		Acc: Acc{
			accounts: accounts,
			persons:  persons,
			objects:  objects,
			getLogger: func() Logger {
				return logger
			},
		},
		Tar: Tar{
			targets: targets,
			getLogger: func() Logger {
				return logger
			},
		},
	}
}
