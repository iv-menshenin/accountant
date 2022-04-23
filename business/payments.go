package business

import (
	"context"
	"github.com/iv-menshenin/accountant/model/generic"
	"github.com/iv-menshenin/accountant/storage"
	"github.com/iv-menshenin/accountant/utils/uuid"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/model/request"
)

func (p *Pay) PaymentCreate(ctx context.Context, q request.PostPaymentQuery) (*domain.Payment, error) {
	var payment = domain.Payment{
		PaymentID:   uuid.NewUUID(),
		AccountID:   q.AccountID,
		PaymentData: q.Data,
	}
	err := p.payments.Create(ctx, payment)
	if err != nil {
		return nil, err
	}
	return p.payments.Lookup(ctx, payment.PaymentID)
}

func (p *Pay) PaymentGet(ctx context.Context, q request.GetPaymentQuery) (*domain.Payment, error) {
	object, err := p.payments.Lookup(ctx, q.PaymentID)
	if err == storage.ErrNotFound {
		return nil, generic.NotFound{}
	}
	return object, nil
}

func (p *Pay) PaymentDelete(ctx context.Context, q request.DeletePaymentQuery) error {
	err := p.payments.Delete(ctx, q.PaymentID)
	if err == storage.ErrNotFound {
		return generic.NotFound{}
	}
	return err
}

func (p *Pay) PaymentsFind(ctx context.Context, q request.FindPaymentsQuery) ([]domain.Payment, error) {
	payments, err := p.payments.FindBy(ctx, q.AccountID, q.PersonID, q.ObjectID, q.TargetID)
	if err != nil {
		if err == storage.ErrNotFound {
			return nil, generic.NotFound{}
		}
		return nil, err
	}
	if len(payments) == 0 {
		return nil, generic.NotFound{}
	}
	return payments, nil
}
