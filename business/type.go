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
	Pay struct {
		payments  PaymentsCollection
		getLogger func() Logger
	}
	Usr struct {
		users     UsersCollection
		getLogger func() Logger
	}

	App struct {
		Acc
		Tar
		Bil
		Pay
		Usr
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
	}

	ObjectsCollection interface {
		Create(context.Context, uuid.UUID, domain.Object) error
		Lookup(context.Context, uuid.UUID, uuid.UUID) (*domain.Object, error)
		Replace(context.Context, uuid.UUID, uuid.UUID, domain.Object) error
		Delete(context.Context, uuid.UUID, uuid.UUID) error
		Find(context.Context, storage.FindObjectOption) ([]domain.NestedObject, error)
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
		Lookup(context.Context, uuid.UUID) (*domain.Payment, error)
		Replace(context.Context, uuid.UUID, domain.Payment) error
		Delete(context.Context, uuid.UUID) error
		FindBy(context.Context, *uuid.UUID, *uuid.UUID, *uuid.UUID, *uuid.UUID) ([]domain.Payment, error)
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

	UsersCollection interface {
		Create(ctx context.Context, info domain.UserInfo, identity domain.UserIdentity) (*domain.UserInfo, error)
		Lookup(ctx context.Context, ID uuid.UUID) (*domain.UserInfo, error)
		FindByLogin(ctx context.Context, login string) (*domain.UserInfo, error)
		Update(ctx context.Context, info domain.UserInfo) (*domain.UserInfo, error)
		Delete(ctx context.Context, ID uuid.UUID) error
		Find(ctx context.Context, searchPattern string) ([]domain.UserInfo, error)
	}
)

func New(
	logger Logger,
	accounts AccountsCollection,
	persons PersonsCollection,
	objects ObjectsCollection,
	targets TargetsCollection,
	bills BillsCollection,
	payments PaymentsCollection,
	users UsersCollection,
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
		Pay: Pay{
			payments: payments,
			getLogger: func() Logger {
				return logger
			},
		},
		Usr: Usr{
			users: users,
			getLogger: func() Logger {
				return logger
			},
		},
	}
}
