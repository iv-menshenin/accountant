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
		AccountSave(context.Context, model.PostAccountQuery) (*model.Account, error)
	}
	AccountCreator interface {
		AccountCreate(context.Context, model.PostAccountQuery) (*model.Account, error)
	}
	AccountProcessor interface {
		AccountCreator
		AccountGetter
		AccountSaver
	}
	Accounts struct {
		processor AccountProcessor
	}
)

const (
	accountID = "account_id"
)

func (a *Accounts) LookupPathPattern() string {
	return fmt.Sprintf("/account/{%s:[0-9a-f]+}", accountID)
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

func NewAccountsEP(ap AccountProcessor) *Accounts {
	return &Accounts{
		processor: ap,
	}
}
