package mongodb

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/storage"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

func Test_Users(t *testing.T) {
	once.Do(initTestEnv)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err := testStorage.mongo.Users().Collection.DeleteMany(ctx, bson.D{})
	if err != nil {
		t.Fatalf("cannot delete users: %s", err)
	}

	users := testStorage.NewUsersCollection(storage.MapMongodbErrors)
	var (
		info = domain.UserInfo{
			ID:          uuid.NewUUID(),
			Name:        "Igor",
			Surname:     "Menshenin",
			EMail:       "igor_m@devaliada.ru",
			Permissions: []domain.Permission{"read", "write"},
		}
		identity = domain.UserIdentity{
			Login:    "devalio",
			Password: "dfslkgjsdhgshvlsungs",
		}
	)

	t.Run("Create User", func(t *testing.T) {
		_, err = users.Create(ctx, info, identity)
		if err != nil {
			t.Fatalf("cannot insert user: %s", err)
		}

		got, err := users.FindByLogin(ctx, identity.Login)
		if err != nil {
			t.Fatalf("cannot find user: %s", err)
		}
		if !reflect.DeepEqual(info, *got) {
			t.Errorf("matching error\nwant: %+v\n got: %+v", info, *got)
		}
	})

	t.Run("Find User", func(t *testing.T) {
		got, err := users.Find(ctx, "igor_m")
		if err != nil {
			t.Fatalf("cannot find user: %s", err)
		}
		if !reflect.DeepEqual([]domain.UserInfo{info}, got) {
			t.Errorf("matching error\nwant: %+v\n got: %+v", info, got)
		}
	})

	t.Run("Update User", func(t *testing.T) {
		var newInfo = domain.UserInfo{
			ID:          info.ID,
			Name:        "Igor",
			Surname:     "Kolbasa",
			Permissions: []domain.Permission{"get", "set"},
		}
		_, err = users.Update(ctx, newInfo)
		if err != nil {
			t.Fatalf("cannot update user: %s", err)
		}

		got, err := users.FindByLogin(ctx, identity.Login)
		if err != nil {
			t.Fatalf("cannot find user: %s", err)
		}
		if !reflect.DeepEqual(newInfo, *got) {
			t.Errorf("matching error\nwant: %+v\n got: %+v", newInfo, *got)
		}
	})

	t.Run("Delete User", func(t *testing.T) {
		if err := users.Delete(ctx, info.ID); err != nil {
			t.Errorf("cannot delete user: %s", err)
		}

		_, err := users.FindByLogin(ctx, identity.Login)
		if !errors.Is(err, storage.ErrNotFound) {
			t.Fatalf("expected NotFound error, got: %s", err)
		}
	})
}
