package business

import (
	"context"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/model/generic"
	"github.com/iv-menshenin/accountant/model/request"
	"github.com/iv-menshenin/accountant/storage"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

func (b *Bil) BillCreate(ctx context.Context, q request.PostBillQuery) (*domain.Bill, error) {
	var bill = domain.Bill{
		BillID:    uuid.NewUUID(),
		AccountID: q.AccountID,
		BillData:  q.Data,
	}
	err := b.bills.Create(ctx, bill)
	if err != nil {
		b.getLogger().Error("unable to create bill: %s", err)
		return nil, err
	}
	return b.bills.Lookup(ctx, bill.BillID)
}

func (b *Bil) BillGet(ctx context.Context, q request.GetBillQuery) (*domain.Bill, error) {
	bill, err := b.bills.Lookup(ctx, q.BillID)
	if err == storage.ErrNotFound {
		b.getLogger().Warning("bill not found %s", q.BillID)
		return nil, generic.NotFound{}
	}
	if err != nil {
		b.getLogger().Error("unable to lookup bill %s: %s", q.BillID, err)
		return nil, err
	}
	return bill, nil
}

func (b *Bil) BillSave(ctx context.Context, q request.PutBillQuery) (*domain.Bill, error) {
	bill, err := b.bills.Lookup(ctx, q.BillID)
	if err == storage.ErrNotFound {
		b.getLogger().Warning("bill not found: %s", q.BillID)
		return nil, generic.NotFound{}
	}
	if err != nil {
		b.getLogger().Error("unable to get bill data %s: %s", q.BillID, err)
		return nil, err
	}
	bill.BillData = q.Data
	if err = b.bills.Replace(ctx, *bill); err != nil {
		b.getLogger().Error("unable to replace bill %s: %s", q.BillID, err)
		return nil, err
	}
	return bill, nil
}

func (b *Bil) BillDelete(ctx context.Context, q request.DeleteBillQuery) error {
	err := b.bills.Delete(ctx, q.BillID)
	if err == storage.ErrNotFound {
		b.getLogger().Error("unable to delete bill %s: not found", q.BillID)
		return generic.NotFound{}
	}
	if err != nil {
		b.getLogger().Error("unable to delete bill %s: %s", q.BillID, err)
	}
	return err
}

func (b *Bil) BillsFind(ctx context.Context, q request.FindBillsQuery) (bills []domain.Bill, err error) {
	bills, err = b.bills.FindBy(ctx, q.AccountID, q.TargetID, q.Period)
	if err != nil {
		b.getLogger().Error("unable to find bills: %s", err)
		return nil, err
	}
	if len(bills) == 0 {
		return nil, generic.NotFound{}
	}
	return
}
