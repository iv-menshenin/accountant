package http

import (
	"github.com/gorilla/mux"
	"github.com/iv-menshenin/accountant/transport/internal/http/ep"
	"github.com/iv-menshenin/accountant/transport/internal/http/static"
	"net/http"
)

type (
	RequestProcessor interface {
		ep.AccountProcessor
	}
)

func makeRouter(rp RequestProcessor) http.Handler {
	router := mux.NewRouter()

	stat := static.New()
	stat.SetupRouting(router.PathPrefix("/www").Subrouter())

	accounts := ep.NewAccountsEP(rp)
	accounts.SetupRouting(router.PathPrefix("/api").Subrouter())

	return router
}
