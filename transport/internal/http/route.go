package http

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/iv-menshenin/accountant/transport/internal/http/ep"
	"github.com/iv-menshenin/accountant/transport/internal/http/static"
)

type (
	RequestProcessor interface {
		ep.AccountProcessor
		ep.PersonProcessor
		ep.ObjectProcessor
		ep.TargetProcessor
		ep.BillProcessor
		ep.PaymentProcessor
		ep.UserProcessor
	}
)

const (
	PathAuth = "/auth"
	PathAPI  = "/api"
	PathWWW  = "/"
)

func makeRouter(rp RequestProcessor, auth AuthCore, logger *log.Logger) http.Handler {
	router := mux.NewRouter()
	router.Methods(http.MethodOptions).Handler(http.HandlerFunc(optionHandler))

	apiSubRouter := router.PathPrefix(PathAPI).Subrouter()
	if auth != nil {
		apiSubRouter.Use(auth.Middleware())
	}

	accounts := ep.NewAccountsEP(rp)
	accounts.SetupRouting(apiSubRouter)

	persons := ep.NewPersonsEP(rp)
	persons.SetupRouting(apiSubRouter)

	objects := ep.NewObjectsEP(rp)
	objects.SetupRouting(apiSubRouter)

	targets := ep.NewTargetsEP(rp)
	targets.SetupRouting(apiSubRouter)

	bills := ep.NewBillsEP(rp)
	bills.SetupRouting(apiSubRouter)

	payments := ep.NewPaymentsEP(rp)
	payments.SetupRouting(apiSubRouter)

	users := ep.NewUsersEP(rp)
	users.SetupRouting(apiSubRouter)

	authSubRouter := router.PathPrefix(PathAuth).Subrouter()
	authP := ep.NewAuthEP(auth)
	authP.SetupRouting(authSubRouter)

	wwwSubRouter := router.PathPrefix(PathWWW).Subrouter()
	stat := static.New(logger)
	stat.SetupRouting(wwwSubRouter)

	router.Use(DontCareAboutCORS)
	return router
}

const (
	corsAllowHeaders     = "X-Auth-Token,X-Requested-With"
	corsAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "true"
)

func DontCareAboutCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		w.Header().Set("Access-Control-Allow-Headers", corsAllowHeaders)
		w.Header().Set("Access-Control-Allow-Methods", corsAllowMethods)
		w.Header().Set("Access-Control-Allow-Origin", corsAllowOrigin)

		next.ServeHTTP(w, r)
	})
}

func optionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", corsAllowCredentials)
	w.Header().Set("Access-Control-Allow-Headers", corsAllowHeaders)
	w.Header().Set("Access-Control-Allow-Methods", corsAllowMethods)
	w.Header().Set("Access-Control-Allow-Origin", corsAllowOrigin)
	w.WriteHeader(http.StatusOK)
}
