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
		ep.PersonProcessor
		ep.ObjectProcessor
	}
)

func makeRouter(rp RequestProcessor) http.Handler {
	router := mux.NewRouter()
	wwwSubRouter := router.PathPrefix("/www").Subrouter()
	apiSubRouter := router.PathPrefix("/api").Subrouter()

	stat := static.New()
	stat.SetupRouting(wwwSubRouter)

	accounts := ep.NewAccountsEP(rp)
	accounts.SetupRouting(apiSubRouter)

	persons := ep.NewPersonsEP(rp)
	persons.SetupRouting(apiSubRouter)

	objects := ep.NewObjectsEP(rp)
	objects.SetupRouting(apiSubRouter)

	return router
}
