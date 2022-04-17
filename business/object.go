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

func (a *Acc) ObjectCreate(ctx context.Context, q request.PostObjectQuery) (*domain.Object, error) {
	var object = domain.Object{
		ObjectID:   uuid.NewUUID(),
		ObjectData: q.ObjectData,
	}
	err := a.objects.Create(ctx, q.AccountID, object)
	if err != nil {
		return nil, err
	}
	return a.objects.Lookup(ctx, q.AccountID, object.ObjectID)
}

func (a *Acc) ObjectGet(ctx context.Context, q request.GetObjectQuery) (*domain.Object, error) {
	object, err := a.objects.Lookup(ctx, q.AccountID, q.ObjectID)
	if err == storage.ErrNotFound {
		return nil, generic.NotFound{}
	}
	return object, nil
}

func (a *Acc) ObjectSave(ctx context.Context, q request.PutObjectQuery) (*domain.Object, error) {
	object, err := a.objects.Lookup(ctx, q.AccountID, q.ObjectID)
	if err == storage.ErrNotFound {
		return nil, generic.NotFound{}
	}
	object.ObjectData = q.ObjectData
	if err = a.objects.Replace(ctx, q.AccountID, q.ObjectID, *object); err != nil {
		return nil, err
	}
	return object, nil
}

func (a *Acc) ObjectDelete(ctx context.Context, q request.DeleteObjectQuery) error {
	err := a.objects.Delete(ctx, q.AccountID, q.ObjectID)
	if err == storage.ErrNotFound {
		return generic.NotFound{}
	}
	return err
}

func (a *Acc) ObjectsFind(ctx context.Context, q request.FindObjectsQuery) ([]domain.Object, error) {
	var findOption storage2.FindObjectOption
	findOption.FillFromQuery(q)
	objects, err := a.objects.Find(ctx, findOption)
	if err != nil {
		if err == storage.ErrNotFound {
			return nil, generic.NotFound{}
		}
		return nil, err
	}
	if len(objects) == 0 {
		return nil, generic.NotFound{}
	}
	return objects, nil
}
