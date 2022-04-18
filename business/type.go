package business

import (
	"context"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/storage"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

type (
	Acc struct {
		accounts  AccountsCollection
		persons   PersonsCollection
		objects   ObjectsCollection
		getLogger func() Logger
	}
	Tar struct {
		targets   TargetsCollection
		getLogger func() Logger
	}
	Bil struct {
		bills     BillsCollection
		getLogger func() Logger
	}
	App struct {
		Acc
		Tar
		Bil
	}
	Logger interface {
		Warning(format string, args ...interface{})
		Debug(format string, args ...interface{})
		Error(format string, args ...interface{})
	}

	PersonsCollection interface {
		Create(context.Context, uuid.UUID, domain.Person) error
		Lookup(context.Context, uuid.UUID, uuid.UUID) (*domain.Person, error)
		Replace(context.Context, uuid.UUID, uuid.UUID, domain.Person) error
		Delete(context.Context, uuid.UUID, uuid.UUID) error
		Find(context.Context, storage.FindPersonOption) ([]domain.Person, error)
	}

	ObjectsCollection interface {
		Create(context.Context, uuid.UUID, domain.Object) error
		Lookup(context.Context, uuid.UUID, uuid.UUID) (*domain.Object, error)
		Replace(context.Context, uuid.UUID, uuid.UUID, domain.Object) error
		Delete(context.Context, uuid.UUID, uuid.UUID) error
		Find(context.Context, storage.FindObjectOption) ([]domain.Object, error)
	}

	AccountsCollection interface {
		Create(context.Context, domain.Account) error
		Lookup(context.Context, uuid.UUID) (*domain.Account, error)
		Replace(context.Context, uuid.UUID, domain.Account) error
		Delete(context.Context, uuid.UUID) error
		Find(context.Context, storage.FindAccountOption) ([]domain.Account, error)
	}

	PaymentsCollection interface {
		Create(context.Context, domain.Payment) error
		Delete(context.Context, uuid.UUID) error
		FindByAccount(context.Context, uuid.UUID) ([]domain.Payment, error)
		FindByIDs(context.Context, []uuid.UUID) ([]domain.Payment, error)
	}

	TargetsCollection interface {
		Create(context.Context, domain.Target) error
		Update(context.Context, uuid.UUID, domain.TargetData) error
		Lookup(context.Context, uuid.UUID) (*domain.Target, error)
		Delete(context.Context, uuid.UUID) error
		FindByPeriod(context.Context, storage.FindTargetOption) ([]domain.Target, error)
	}

	BillsCollection interface {
		Create(context.Context, domain.Bill) error
		Lookup(context.Context, uuid.UUID) (*domain.Bill, error)
		Replace(context.Context, domain.Bill) error
		FindBy(context.Context, *uuid.UUID, *uuid.UUID, *domain.Period) ([]domain.Bill, error)
		Delete(context.Context, uuid.UUID) error
	}
)

func New(
	logger Logger,
	accounts AccountsCollection,
	persons PersonsCollection,
	objects ObjectsCollection,
	targets TargetsCollection,
	bills BillsCollection,
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
		Bil: Bil{
			bills: bills,
			getLogger: func() Logger {
				return logger
			},
		},
	}
}
