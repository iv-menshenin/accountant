package business

import (
	"context"
	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/model/uuid"
	"github.com/iv-menshenin/accountant/store"
)

func (a *App) AccountCreate(ctx context.Context, acc model.PostAccountQuery) (*model.Account, error) {
	var account = model.Account{
		AccountID:   uuid.NewUUID(),
		AccountData: acc.AccountData,
	}
	err := a.accounts.Create(ctx, account)
	if err != nil {
		return nil, err
	}
	return a.accounts.Lookup(ctx, account.AccountID)
}

func (a *App) AccountGet(ctx context.Context, q model.GetAccountQuery) (*model.Account, error) {
	account, err := a.accounts.Lookup(ctx, q.ID)
	if err == store.ErrNotFound {
		return nil, model.NotFound{}
	}
	return account, nil
}

func (a *App) AccountSave(ctx context.Context, acc model.PutAccountQuery) (*model.Account, error) {
	account, err := a.accounts.Lookup(ctx, acc.ID)
	if err == store.ErrNotFound {
		return nil, model.NotFound{}
	}
	account.AccountData = acc.Account
	if err = a.accounts.Replace(ctx, acc.ID, *account); err != nil {
		return nil, err
	}
	return account, nil
}

func (a *App) AccountsFind(ctx context.Context) ([]model.Account, error) {
	return nil, nil
}

func (a *App) AccountDelete(ctx context.Context) error {
	return nil
}
