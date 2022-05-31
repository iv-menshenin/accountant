package static

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
)

type (
	Resources struct {
		staticFiles string
		logger      *log.Logger
	}
)

var (
	staticPath = flag.String("www-path", os.Getenv("WWW_PATH"), "http-server html files path")
)

const (
	contentType = "Content-Type"
)

func New(logger *log.Logger) *Resources {
	if *staticPath == "" {
		*staticPath = "./www"
	}
	return &Resources{
		staticFiles: *staticPath,
		logger:      logger,
	}
}

func (r *Resources) SetupRouting(router *mux.Router) {
	router.Path("/js/{filename:[a-z0-9\\-_/]+}.js").Methods(http.MethodGet).Handler(http.HandlerFunc(r.Script))
	router.Path("/html/{filename:[a-z0-9\\-_/]+}.html").Methods(http.MethodGet).Handler(http.HandlerFunc(r.Html))
	router.Path("/css/{filename:[a-z0-9\\-_/]+}.css").Methods(http.MethodGet).Handler(http.HandlerFunc(r.Css))
	router.Path("/png/{filename:[a-z0-9\\-_/]+}.png").Methods(http.MethodGet).Handler(http.HandlerFunc(r.Png))
	router.Methods(http.MethodGet).Handler(http.HandlerFunc(r.Any))
}

func (r *Resources) Script(w http.ResponseWriter, q *http.Request) {
	fileName := mux.Vars(q)["filename"]
	f, err := os.Open(fmt.Sprintf("%s/js/%s.js", r.staticFiles, fileName))
	if err != nil {
		r.logger.Printf("FILE ERROR: %s [%v]\n", fileName, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	w.Header().Set(contentType, "application/javascript; charset=utf-8")
	_, err = io.Copy(w, f)
	if err != nil {
		r.logger.Println(err)
	}
}

func (r *Resources) Html(w http.ResponseWriter, q *http.Request) {
	fileName := mux.Vars(q)["filename"]
	f, err := os.Open(fmt.Sprintf("%s/html/%s.html", r.staticFiles, fileName))
	if err != nil {
		r.logger.Printf("FILE ERROR: %s [%v]\n", fileName, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	w.Header().Set(contentType, "text/html; charset=utf-8")
	_, err = io.Copy(w, f)
	if err != nil {
		r.logger.Println(err)
	}
}

func (r *Resources) Css(w http.ResponseWriter, q *http.Request) {
	fileName := mux.Vars(q)["filename"]
	f, err := os.Open(fmt.Sprintf("%s/css/%s.css", r.staticFiles, fileName))
	if err != nil {
		r.logger.Printf("FILE ERROR: %s [%v]\n", fileName, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	w.Header().Set(contentType, "text/css; charset=utf-8")
	_, err = io.Copy(w, f)
	if err != nil {
		r.logger.Println(err)
	}
}

func (r *Resources) Png(w http.ResponseWriter, q *http.Request) {
	fileName := mux.Vars(q)["filename"]
	f, err := os.Open(fmt.Sprintf("%s/png/%s.png", r.staticFiles, fileName))
	if err != nil {
		r.logger.Printf("FILE ERROR: %s [%v]\n", fileName, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	w.Header().Set(contentType, "image/png")
	_, err = io.Copy(w, f)
	if err != nil {
		r.logger.Println(err)
	}
}

var startPath = flag.String("www-start", os.Getenv("HTML_START"), "http-server homepage")

func (r *Resources) Any(w http.ResponseWriter, q *http.Request) {
	fileName := q.URL.Path
	if fileName == "" {
		if *startPath == "/" || *startPath == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.Redirect(w, q, *startPath, http.StatusFound)
	}
	if strings.Contains(fileName, "../") {
		r.logger.Printf("DETECTED UPLEVEL: %s\n", fileName)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	f, err := os.Open(r.staticFiles + fileName)
	if err != nil {
		r.logger.Printf("FILE ERROR: %s [%v]\n", fileName, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()

	mimeType := mime.TypeByExtension(path.Ext(fileName))
	if mimeType != "" {
		w.Header().Set(contentType, mimeType)
	}
	_, err = io.Copy(w, f)
	if err != nil {
		r.logger.Println(err)
	}
}
