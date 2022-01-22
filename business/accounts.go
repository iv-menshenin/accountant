package business

import (
	"context"

	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/model/uuid"
	"github.com/iv-menshenin/accountant/store"
)

func (a *App) AccountCreate(ctx context.Context, q model.PostAccountQuery) (*model.Account, error) {
	var account = model.Account{
		AccountID:   uuid.NewUUID(),
		AccountData: q.AccountData,
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

func (a *App) AccountSave(ctx context.Context, q model.PutAccountQuery) (*model.Account, error) {
	account, err := a.accounts.Lookup(ctx, q.ID)
	if err == store.ErrNotFound {
		return nil, model.NotFound{}
	}
	account.AccountData = q.Account
	if err = a.accounts.Replace(ctx, q.ID, *account); err != nil {
		return nil, err
	}
	return account, nil
}

func (a *App) AccountDelete(ctx context.Context, q model.DeleteAccountQuery) error {
	err := a.accounts.Delete(ctx, q.ID)
	if err == store.ErrNotFound {
		return model.NotFound{}
	}
	return err
}

func (a *App) AccountsFind(ctx context.Context, q model.FindAccountsQuery) ([]model.Account, error) {
	var findOption store.FindAccountOption
	findOption.FillFromQuery(q)
	accounts, err := a.accounts.Find(ctx, findOption)
	if err != nil {
		return nil, err
	}
	if len(accounts) == 0 {
		return nil, model.NotFound{}
	}
	return accounts, nil
}
