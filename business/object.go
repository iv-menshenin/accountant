package business

import (
	"context"

	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/model/uuid"
	"github.com/iv-menshenin/accountant/storage"
)

func (a *App) ObjectCreate(ctx context.Context, q model.PostObjectQuery) (*model.Object, error) {
	var object = model.Object{
		ObjectID:   uuid.NewUUID(),
		ObjectData: q.ObjectData,
	}
	err := a.objects.Create(ctx, q.AccountID, object)
	if err != nil {
		return nil, err
	}
	return a.objects.Lookup(ctx, q.AccountID, object.ObjectID)
}

func (a *App) ObjectGet(ctx context.Context, q model.GetObjectQuery) (*model.Object, error) {
	object, err := a.objects.Lookup(ctx, q.AccountID, q.ObjectID)
	if err == storage.ErrNotFound {
		return nil, model.NotFound{}
	}
	return object, nil
}

func (a *App) ObjectSave(ctx context.Context, q model.PutObjectQuery) (*model.Object, error) {
	object, err := a.objects.Lookup(ctx, q.AccountID, q.ObjectID)
	if err == storage.ErrNotFound {
		return nil, model.NotFound{}
	}
	object.ObjectData = q.ObjectData
	if err = a.objects.Replace(ctx, q.AccountID, q.ObjectID, *object); err != nil {
		return nil, err
	}
	return object, nil
}

func (a *App) ObjectDelete(ctx context.Context, q model.DeleteObjectQuery) error {
	err := a.objects.Delete(ctx, q.AccountID, q.ObjectID)
	if err == storage.ErrNotFound {
		return model.NotFound{}
	}
	return err
}

func (a *App) ObjectsFind(ctx context.Context, q model.FindObjectsQuery) ([]model.Object, error) {
	var findOption model.FindObjectOption
	findOption.FillFromQuery(q)
	objects, err := a.objects.Find(ctx, findOption)
	if err != nil {
		if err == storage.ErrNotFound {
			return nil, model.NotFound{}
		}
		return nil, err
	}
	if len(objects) == 0 {
		return nil, model.NotFound{}
	}
	return objects, nil
}
