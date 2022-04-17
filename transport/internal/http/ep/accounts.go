package ep

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/iv-menshenin/accountant/model/request"
)

type (
	AccountProcessor interface {
		request.AccountCreator
		request.AccountGetter
		request.AccountSaver
		request.AccountDeleter
		request.AccountFinder
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
	pathSegmentAccounts    = "/accounts"
	parameterNameAccountID = "account_id"
	parameterNameAccount   = "account"
)

func (a *Accounts) SetupRouting(router *mux.Router) {
	accountsWithIDPath := fmt.Sprintf("%s/{%s:[0-9a-f\\-]+}", pathSegmentAccounts, parameterNameAccountID)

	router.Path(accountsWithIDPath).Methods(http.MethodGet).Handler(a.LookupHandler())
	router.Path(pathSegmentAccounts).Methods(http.MethodPost).Handler(a.PostHandler())
	router.Path(accountsWithIDPath).Methods(http.MethodPut).Handler(a.PutHandler())
	router.Path(accountsWithIDPath).Methods(http.MethodDelete).Handler(a.DeleteHandler())
	router.Path(pathSegmentAccounts).Methods(http.MethodGet).Handler(a.FindHandler())

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

func getAccountMapper(r *http.Request) (q request.GetAccountQuery, err error) {
	id := mux.Vars(r)[parameterNameAccountID]
	if id == "" {
		err = errors.New(parameterNameAccountID + " must not be empty")
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

func postAccountMapper(r *http.Request) (q request.PostAccountQuery, err error) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&q.AccountData)
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

func putAccountMapper(r *http.Request) (q request.PutAccountQuery, err error) {
	id := mux.Vars(r)[parameterNameAccountID]
	if id == "" {
		err = errors.New(parameterNameAccountID + " must not be empty")
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

func deleteAccountMapper(r *http.Request) (q request.DeleteAccountQuery, err error) {
	id := mux.Vars(r)[parameterNameAccountID]
	if id == "" {
		err = errors.New(parameterNameAccountID + " must not be empty")
		return
	}
	err = q.ID.FromString(id)
	return
}

func (a *Accounts) FindHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := findAccountMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		accounts, err := a.processor.AccountsFind(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, accounts)
	}
}

func findAccountMapper(r *http.Request) (q request.FindAccountsQuery, err error) {
	params := queryParams{r: r}
	if account, ok := params.vars(parameterNameAccount); ok {
		q.Account = &account
	}
	return
}
