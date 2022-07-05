package ep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/model/request"
)

type (
	UserProcessor interface {
		UserCreate(ctx context.Context, q request.PostUserQuery) (*domain.UserInfo, error)
		UserGet(ctx context.Context, q request.GetUserQuery) (*domain.UserInfo, error)
		UserSave(ctx context.Context, q request.PutUserQuery) (*domain.UserInfo, error)
		UserDelete(ctx context.Context, q request.DeleteUserQuery) error
		UsersFind(ctx context.Context, q request.GetUsersQuery) ([]domain.UserInfo, error)
	}
	Users struct {
		processor UserProcessor
	}
)

const (
	pathSegmentUsers    = "/users"
	parameterNameUserID = "user_id"
)

func (u *Users) SetupRouting(router *mux.Router) {
	userWithIDPath := fmt.Sprintf("%s/{%s:[0-9a-f\\-]+}", pathSegmentUsers, parameterNameUserID)

	router.Path(pathSegmentUsers).Methods(http.MethodPost).Handler(u.PostHandler())
	router.Path(userWithIDPath).Methods(http.MethodGet).Handler(u.GetHandler())
}

func (u *Users) PostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := postUserMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		info, err := u.processor.UserCreate(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, info)
	}
}

func postUserMapper(r *http.Request) (q request.PostUserQuery, err error) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&q)
	if strings.Contains(q.EMail, "@") && !strings.Contains(q.EMail, " ") {
		err = errors.New("email format validation error")
		return
	}
	q.Login = strings.TrimSpace(q.Login)
	if q.Login == "" {
		err = errors.New("login is empty")
		return
	}
	return
}

func (u *Users) GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := getUserMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		user, err := u.processor.UserGet(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, user)
	}
}

func getUserMapper(r *http.Request) (q request.GetUserQuery, err error) {
	id := mux.Vars(r)[parameterNameUserID]
	if id == "" {
		err = errors.New(parameterNameTargetID + " must not be empty")
		return
	}
	err = q.ID.FromString(id)
	return
}

func NewUsersEP(proc UserProcessor) *Users {
	return &Users{
		processor: proc,
	}
}
