package ep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/iv-menshenin/accountant/model"
	"net/http"
)

type (
	AccountGetter interface {
		AccountGet(context.Context, model.GetAccountQuery) (*model.Account, error)
	}
	AccountSaver interface {
		AccountSave(context.Context, model.PutAccountQuery) (*model.Account, error)
	}
	AccountCreator interface {
		AccountCreate(context.Context, model.PostAccountQuery) (*model.Account, error)
	}
	AccountDeleter interface {
		AccountDelete(context.Context, model.DeleteAccountQuery) error
	}
	AccountFinder interface {
		AccountsFind(context.Context, model.FindAccountsQuery) error
	}
	AccountProcessor interface {
		AccountCreator
		AccountGetter
		AccountSaver
		AccountDeleter
	}
	Accounts struct {
		processor AccountProcessor
	}
)

func NewAccountsEP(ap AccountProcessor) *Accounts {
	return &Accounts{
		processor: ap,
	}
}

const (
	accountID = "account_id"
)

func (a *Accounts) SetupRouting(router *mux.Router) {
	const accountsPath = "/accounts"
	accountsWithIDPath := fmt.Sprintf("%s/{%s:[0-9a-f\\-]+}", accountsPath, accountID)

	router.Path(accountsWithIDPath).Methods(http.MethodGet).Handler(a.LookupHandler())
	router.Path(accountsPath).Methods(http.MethodPost).Handler(a.PostHandler())
	router.Path(accountsWithIDPath).Methods(http.MethodPut).Handler(a.PutHandler())
	router.Path(accountsWithIDPath).Methods(http.MethodDelete).Handler(a.DeleteHandler())

}

func (a *Accounts) LookupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := getAccountMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		account, err := a.processor.AccountGet(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, account)
	}
}

func getAccountMapper(r *http.Request) (q model.GetAccountQuery, err error) {
	id := mux.Vars(r)[accountID]
	if id == "" {
		err = errors.New(accountID + " must not be empty")
		return
	}
	err = q.ID.FromString(id)
	return
}

func (a *Accounts) PostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := postAccountMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		account, err := a.processor.AccountCreate(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, account)
	}
}

func postAccountMapper(r *http.Request) (q model.PostAccountQuery, err error) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&q)
	return
}

func (a *Accounts) PutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := putAccountMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		account, err := a.processor.AccountSave(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, account)
	}
}

func putAccountMapper(r *http.Request) (q model.PutAccountQuery, err error) {
	id := mux.Vars(r)[accountID]
	if id == "" {
		err = errors.New(accountID + " must not be empty")
		return
	}
	if err = q.ID.FromString(id); err != nil {
		return q, err
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&q.Account)
	return
}

func (a *Accounts) DeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := deleteAccountMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		err = a.processor.AccountDelete(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeNoContent(w)
	}
}

func deleteAccountMapper(r *http.Request) (q model.DeleteAccountQuery, err error) {
	id := mux.Vars(r)[accountID]
	if id == "" {
		err = errors.New(accountID + " must not be empty")
		return
	}
	err = q.ID.FromString(id)
	return
}
