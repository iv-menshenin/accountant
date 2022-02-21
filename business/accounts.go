package business

import (
	"context"

	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/model/uuid"
	"github.com/iv-menshenin/accountant/storage"
)

func (a *App) AccountCreate(ctx context.Context, q model.PostAccountQuery) (*model.Account, error) {
	var account = model.Account{
		AccountID:   uuid.NewUUID(),
		AccountData: q.AccountData,
	}
	err := a.accounts.Create(ctx, account)
	if err != nil {
		a.logger.Error("unable to create account: %s", err)
		return nil, err
	}
	return a.accounts.Lookup(ctx, account.AccountID)
}

func (a *App) AccountGet(ctx context.Context, q model.GetAccountQuery) (*model.Account, error) {
	account, err := a.accounts.Lookup(ctx, q.ID)
	if err == storage.ErrNotFound {
		a.logger.Warning("account not found %s", q.ID)
		return nil, model.NotFound{}
	}
	if err != nil {
		a.logger.Error("unable to lookup account %s: %s", q.ID, err)
		return nil, err
	}
	return account, nil
}

func (a *App) AccountSave(ctx context.Context, q model.PutAccountQuery) (*model.Account, error) {
	account, err := a.accounts.Lookup(ctx, q.ID)
	if err == storage.ErrNotFound {
		a.logger.Warning("account not found: %s", q.ID)
		return nil, model.NotFound{}
	}
	if err != nil {
		a.logger.Error("unable to get account data %s: %s", q.ID, err)
		return nil, err
	}
	account.AccountData = q.Account
	if err = a.accounts.Replace(ctx, q.ID, *account); err != nil {
		a.logger.Error("unable to replace account %s: %s", q.ID, err)
		return nil, err
	}
	return account, nil
}

func (a *App) AccountDelete(ctx context.Context, q model.DeleteAccountQuery) error {
	err := a.accounts.Delete(ctx, q.ID)
	if err == storage.ErrNotFound {
		a.logger.Error("unable to delete account %s: not found", q.ID)
		return model.NotFound{}
	}
	if err != nil {
		a.logger.Error("unable to delete account %s: %s", q.ID, err)
	}
	return err
}

func (a *App) AccountsFind(ctx context.Context, q model.FindAccountsQuery) ([]model.Account, error) {
	var findOption model.FindAccountOption
	findOption.FillFromQuery(q)
	accounts, err := a.accounts.Find(ctx, findOption)
	if err != nil {
		a.logger.Error("unable to find account: %s", err)
		return nil, err
	}
	if len(accounts) == 0 {
		return nil, model.NotFound{}
	}
	return accounts, nil
}
