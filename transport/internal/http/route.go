package http

import (
	"flag"
	"log"
	"net/http"
	"os"

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
	}
)

const (
	PathAuth = "/auth"
	PathAPI  = "/api"
	PathWWW  = "/www"
)

var (
	startPath = flag.String("www-start", os.Getenv("HTML_START"), "http-server homepage")
)

func makeRouter(rp RequestProcessor, auth AuthCore, logger *log.Logger) http.Handler {
	router := mux.NewRouter()
	router.Methods(http.MethodOptions).Handler(http.HandlerFunc(optionHandler))
	router.Path("/").Methods(http.MethodGet).Handler(http.HandlerFunc(func(w http.ResponseWriter, q *http.Request) {
		if *startPath == "/" || *startPath == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.Redirect(w, q, *startPath, http.StatusFound)
	}))

	wwwSubRouter := router.PathPrefix(PathWWW).Subrouter()

	stat := static.New(logger)
	stat.SetupRouting(wwwSubRouter)

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

	authSubRouter := router.PathPrefix(PathAuth).Subrouter()

	authP := ep.NewAuthEP(auth)
	authP.SetupRouting(authSubRouter)

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
