package business

import (
	"context"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/model/generic"
	"github.com/iv-menshenin/accountant/model/request"
	"github.com/iv-menshenin/accountant/storage"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

func (u *Usr) UserCreate(ctx context.Context, q request.PostUserQuery) (*domain.UserInfo, error) {
	if err := u.checkLoginNotExists(ctx, q.Login); err != nil {
		return nil, err
	}

	perm, err := domain.StringsToPermissions(q.Permissions)
	if err != nil {
		return nil, err
	}
	var i = domain.UserInfo{
		ID:          uuid.NewUUID(),
		Name:        q.Name,
		Surname:     q.Surname,
		EMail:       q.EMail,
		Permissions: perm,
	}
	var idt = domain.UserIdentity{
		Login:    q.Login,
		Password: "",
	}
	return u.users.Create(ctx, i, idt)
}

func (u *Usr) checkLoginNotExists(ctx context.Context, login string) error {
	_, err := u.users.FindByLogin(ctx, login)
	if err == nil {
		err = generic.AlreadyExists{}
	}
	if err == storage.ErrNotFound {
		return nil
	}
	return err
}

func (u *Usr) UserGet(ctx context.Context, q request.GetUserQuery) (*domain.UserInfo, error) {
	user, err := u.users.Lookup(ctx, q.ID)
	if err == storage.ErrNotFound {
		u.getLogger().Warning("user not found %s", q.ID)
		return nil, generic.NotFound{}
	}
	if err != nil {
		u.getLogger().Error("unable to lookup user %s: %s", q.ID, err)
		return nil, err
	}
	return user, nil
}

func (u *Usr) UserSave(ctx context.Context, q request.PutUserQuery) (*domain.UserInfo, error) {
	perm, err := domain.StringsToPermissions(q.Permissions)
	if err != nil {
		return nil, err
	}
	i := domain.UserInfo{
		ID:          uuid.NewUUID(),
		Name:        q.Name,
		Surname:     q.Surname,
		EMail:       q.EMail,
		Permissions: perm,
	}
	info, err := u.users.Update(ctx, i)
	if err == storage.ErrNotFound {
		u.getLogger().Warning("user not fond %s", q.ID)
		return nil, generic.NotFound{}
	}
	if err != nil {
		u.getLogger().Error("unable to update user %s: %s", q.ID, err)
		return nil, err
	}
	return info, nil
}

func (u *Usr) UserDelete(ctx context.Context, q request.DeleteUserQuery) error {
	err := u.users.Delete(ctx, q.ID)
	if err == storage.ErrNotFound {
		u.getLogger().Warning("user not found %s", q.ID)
		return generic.NotFound{}
	}
	if err != nil {
		u.getLogger().Error("unable to delete user %s: %s", q.ID, err)
		return err
	}
	return nil
}

func (u *Usr) UsersFind(ctx context.Context, q request.GetUsersQuery) ([]domain.UserInfo, error) {
	u.users.Find(ctx, q.Pattern)
	found, err := u.users.Find(ctx, q.Pattern)
	if err != nil {
		u.getLogger().Error("unable to find user: %s", err)
		return nil, err
	}
	return found, nil
}
