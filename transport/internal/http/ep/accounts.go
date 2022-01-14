package ep

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/iv-menshenin/accountant/model"
	"net/http"
)

type (
	AccountGetter interface {
		AccountGet(q model.GetAccountQuery) (*model.Account, error)
	}
	AccountSaver interface {
		AccountSave(acc model.SaveAccountQuery) (*model.Account, error)
	}
	AccountProcessor interface {
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
		account, err := a.processor.AccountGet(q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, account)
	}
}

func getAccountMapper(r *http.Request) (q model.GetAccountQuery, err error) {
	if id := mux.Vars(r)[accountID]; id != "" {

		return
	}
	return q, errors.New(accountID + " must not be empty")
}

func (a *Accounts) PostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := postAccountMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		account, err := a.processor.AccountSave(q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, account)
	}
}

func postAccountMapper(r *http.Request) (q model.SaveAccountQuery, err error) {
	return q, errors.New("not implemented")
}

func NewAccountsEP(ap AccountProcessor) *Accounts {
	return &Accounts{
		processor: ap,
	}
}
