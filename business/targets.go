package business

import (
	"context"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/model/generic"
	"github.com/iv-menshenin/accountant/model/request"
	storage2 "github.com/iv-menshenin/accountant/model/storage"
	"github.com/iv-menshenin/accountant/storage"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

func (a *Tar) TargetGet(ctx context.Context, q request.GetTargetQuery) (*domain.Target, error) {
	target, err := a.targets.Lookup(ctx, q.TargetID)
	if err == storage.ErrNotFound {
		a.getLogger().Warning("target not found %s", q.TargetID)
		return nil, generic.NotFound{}
	}
	if err != nil {
		a.getLogger().Error("unable to lookup target %s: %s", q.TargetID, err)
		return nil, err
	}
	return target, nil
}

func (a *Tar) TargetCreate(ctx context.Context, data request.PostTargetQuery) (*domain.Target, error) {
	var target = domain.Target{
		TargetHead: domain.TargetHead{
			TargetID: uuid.NewUUID(),
			Type:     data.Type,
		},
		TargetData: data.Target,
	}
	err := a.targets.Create(ctx, target)
	if err != nil {
		a.getLogger().Error("unable to create target: %s", err)
		return nil, err
	}
	return a.targets.Lookup(ctx, target.TargetID)
}

func (a *Tar) TargetDelete(ctx context.Context, q request.DeleteTargetQuery) error {
	err := a.targets.Delete(ctx, q.TargetID)
	if err == storage.ErrNotFound {
		a.getLogger().Error("unable to delete target %s: not found", q.TargetID)
		return generic.NotFound{}
	}
	if err != nil {
		a.getLogger().Error("unable to delete target %s: %s", q.TargetID, err)
	}
	return err
}

func (a *Tar) TargetsFind(ctx context.Context, q request.FindTargetQuery) ([]domain.Target, error) {
	var findOption = storage2.FindTargetOption{
		ShowClosed: q.ShowClosed,
	}
	if q.Period != nil && q.Period.Year > 0 {
		findOption.Year = q.Period.Year
	}
	if q.Period != nil && q.Period.Month > 0 {
		findOption.Month = q.Period.Month
	}

	targets, err := a.targets.FindByPeriod(ctx, findOption)
	if err != nil {
		a.getLogger().Error("unable to find targets: %s", err)
		return nil, err
	}
	if len(targets) == 0 {
		return nil, generic.NotFound{}
	}
	return targets, nil
}
