package request

import (
	"context"

	"github.com/iv-menshenin/accountant/model/domain"
)

type (
	AccountGetter interface {
		AccountGet(context.Context, GetAccountQuery) (*domain.Account, error)
	}
	AccountSaver interface {
		AccountSave(context.Context, PutAccountQuery) (*domain.Account, error)
	}
	AccountCreator interface {
		AccountCreate(context.Context, PostAccountQuery) (*domain.Account, error)
	}
	AccountDeleter interface {
		AccountDelete(context.Context, DeleteAccountQuery) error
	}
	AccountFinder interface {
		AccountsFind(context.Context, FindAccountsQuery) ([]domain.Account, error)
	}

	ObjectGetter interface {
		ObjectGet(context.Context, GetObjectQuery) (*domain.Object, error)
	}
	ObjectCreator interface {
		ObjectCreate(context.Context, PostObjectQuery) (*domain.Object, error)
	}
	ObjectSaver interface {
		ObjectSave(context.Context, PutObjectQuery) (*domain.Object, error)
	}
	ObjectDeleter interface {
		ObjectDelete(context.Context, DeleteObjectQuery) error
	}
	ObjectFinder interface {
		ObjectsFind(context.Context, FindObjectsQuery) ([]domain.Object, error)
	}

	PersonGetter interface {
		PersonGet(context.Context, GetPersonQuery) (*domain.Person, error)
	}
	PersonCreator interface {
		PersonCreate(context.Context, PostPersonQuery) (*domain.Person, error)
	}
	PersonSaver interface {
		PersonSave(context.Context, PutPersonQuery) (*domain.Person, error)
	}
	PersonDeleter interface {
		PersonDelete(context.Context, DeletePersonQuery) error
	}
	PersonFinder interface {
		PersonsFind(context.Context, FindPersonsQuery) ([]domain.Person, error)
	}

	TargetGetter interface {
		TargetGet(context.Context, GetTargetQuery) (*domain.Target, error)
	}
	TargetCreator interface {
		TargetCreate(context.Context, PostTargetQuery) (*domain.Target, error)
	}
	TargetDeleter interface {
		TargetDelete(context.Context, DeleteTargetQuery) error
	}
	TargetFinder interface {
		TargetsFind(context.Context, FindTargetsQuery) ([]domain.Target, error)
	}

	BillCreator interface {
		BillCreate(context.Context, PostBillQuery) (*domain.Bill, error)
	}
	BillSaver interface {
		BillSave(context.Context, PutBillQuery) (*domain.Bill, error)
	}
	BillGetter interface {
		BillGet(context.Context, GetBillQuery) (*domain.Bill, error)
	}
	BillDeleter interface {
		BillDelete(context.Context, DeleteBillQuery) error
	}
	BillFinder interface {
		BillsFind(context.Context, FindBillsQuery) ([]domain.Bill, error)
	}
)
