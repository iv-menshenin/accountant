package business

import (
	"context"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/model/generic"
	"github.com/iv-menshenin/accountant/model/request"
	"github.com/iv-menshenin/accountant/storage"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

func (a *Acc) AccountCreate(ctx context.Context, q request.PostAccountQuery) (*domain.Account, error) {
	var account = domain.Account{
		AccountID:   uuid.NewUUID(),
		AccountData: q.AccountData,
	}
	err := a.accounts.Create(ctx, account)
	if err != nil {
		a.getLogger().Error("unable to create account: %s", err)
		return nil, err
	}
	return a.accounts.Lookup(ctx, account.AccountID)
}

func (a *Acc) AccountGet(ctx context.Context, q request.GetAccountQuery) (*domain.Account, error) {
	account, err := a.accounts.Lookup(ctx, q.ID)
	if err == storage.ErrNotFound {
		a.getLogger().Warning("account not found %s", q.ID)
		return nil, generic.NotFound{}
	}
	if err != nil {
		a.getLogger().Error("unable to lookup account %s: %s", q.ID, err)
		return nil, err
	}
	return account, nil
}

func (a *Acc) AccountSave(ctx context.Context, q request.PutAccountQuery) (*domain.Account, error) {
	account, err := a.accounts.Lookup(ctx, q.ID)
	if err == storage.ErrNotFound {
		a.getLogger().Warning("account not found: %s", q.ID)
		return nil, generic.NotFound{}
	}
	if err != nil {
		a.getLogger().Error("unable to get account data %s: %s", q.ID, err)
		return nil, err
	}
	account.AccountData = q.Account
	if err = a.accounts.Replace(ctx, q.ID, *account); err != nil {
		a.getLogger().Error("unable to replace account %s: %s", q.ID, err)
		return nil, err
	}
	return account, nil
}

func (a *Acc) AccountDelete(ctx context.Context, q request.DeleteAccountQuery) error {
	err := a.accounts.Delete(ctx, q.ID)
	if err == storage.ErrNotFound {
		a.getLogger().Error("unable to delete account %s: not found", q.ID)
		return generic.NotFound{}
	}
	if err != nil {
		a.getLogger().Error("unable to delete account %s: %s", q.ID, err)
	}
	return err
}

func (a *Acc) AccountsFind(ctx context.Context, q request.FindAccountsQuery) ([]domain.Account, error) {
	var findOption storage.FindAccountOption
	findOption.FillFromQuery(q)
	accounts, err := a.accounts.Find(ctx, findOption)
	if err != nil {
		a.getLogger().Error("unable to find account: %s", err)
		return nil, err
	}
	if len(accounts) == 0 {
		return nil, generic.NotFound{}
	}
	return accounts, nil
}
