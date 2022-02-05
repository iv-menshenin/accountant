package static

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
)

type (
	Resources struct {
		staticFiles string
	}
)

var (
	staticPath = flag.String("www-path", os.Getenv("WWW_PATH"), "http-server html files path")
)

const (
	contentType = "Content-Type"
)

func New() *Resources {
	if *staticPath == "" {
		*staticPath = "./www"
	}
	return &Resources{
		staticFiles: *staticPath,
	}
}

func (r *Resources) SetupRouting(router *mux.Router) {
	router.Path("/js/{filename:[a-z0-9\\-_/]+}.js").Methods(http.MethodGet).Handler(http.HandlerFunc(r.Script))
	router.Path("/html/{filename:[a-z0-9\\-_/]+}.html").Methods(http.MethodGet).Handler(http.HandlerFunc(r.Html))
	router.Path("/css/{filename:[a-z0-9\\-_/]+}.css").Methods(http.MethodGet).Handler(http.HandlerFunc(r.Css))
	router.Path("/png/{filename:[a-z0-9\\-_/]+}.png").Methods(http.MethodGet).Handler(http.HandlerFunc(r.Png))
}

func (r *Resources) Script(w http.ResponseWriter, q *http.Request) {
	fileName := mux.Vars(q)["filename"]
	f, err := os.Open(fmt.Sprintf("%s/js/%s.js", r.staticFiles, fileName))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	w.Header().Set(contentType, "application/javascript; charset=utf-8")
	_, err = io.Copy(w, f)
	if err != nil {
		log.Println(err)
	}
}

func (r *Resources) Html(w http.ResponseWriter, q *http.Request) {
	fileName := mux.Vars(q)["filename"]
	f, err := os.Open(fmt.Sprintf("%s/html/%s.html", r.staticFiles, fileName))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	w.Header().Set(contentType, "text/html; charset=utf-8")
	_, err = io.Copy(w, f)
	if err != nil {
		log.Println(err)
	}
}

func (r *Resources) Css(w http.ResponseWriter, q *http.Request) {
	fileName := mux.Vars(q)["filename"]
	f, err := os.Open(fmt.Sprintf("%s/css/%s.css", r.staticFiles, fileName))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	w.Header().Set(contentType, "text/css; charset=utf-8")
	_, err = io.Copy(w, f)
	if err != nil {
		log.Println(err)
	}
}

func (r *Resources) Png(w http.ResponseWriter, q *http.Request) {
	fileName := mux.Vars(q)["filename"]
	f, err := os.Open(fmt.Sprintf("%s/png/%s.png", r.staticFiles, fileName))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	w.Header().Set(contentType, "image/png")
	_, err = io.Copy(w, f)
	if err != nil {
		log.Println(err)
	}
}
