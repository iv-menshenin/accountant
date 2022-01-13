package http

import (
	"github.com/gorilla/mux"
	"github.com/iv-menshenin/accountant/model"
	"net/http"
)

type (
	AccountGetter interface {
		AccountGet(q model.GetAccountQuery) (*model.Account, error)
	}
	RequestProcessor interface {
		AccountGetter
	}
)

func makeRouter(rp RequestProcessor) http.Handler {
	router := mux.NewRouter()

	router.Path("/account").Methods(http.MethodGet).Handler(makeGetAccountHandler(rp))

	return router
}

func makeGetAccountHandler(rp RequestProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := getAccountMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		account, err := rp.AccountGet(q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, account)
	}
}
