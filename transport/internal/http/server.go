package http

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type (
	Server struct {
		err         error
		connCounter sync.WaitGroup
		server      *http.Server
		stopOnce    sync.Once
	}
)

func (t *Server) ListenAndServe(errCh chan<- error) {
	if err := t.server.ListenAndServe(); err != http.ErrServerClosed {
		errCh <- err
	}
}

func (t *Server) Shutdown(ctx context.Context) error {
	t.stopOnce.Do(func() {
		t.err = t.server.Shutdown(ctx)
	})
	return t.err
}

var (
	httpPort = flag.Int("http-port", envInt("HTTP_PORT", 8080), "http-server port")
	httpHost = flag.String("http-host", envStr("HTTP_HOST", ""), "http-server host")

	httpReadTimeout  = flag.Duration("http-read-timeout", envDuration("HTTP_READ_TIMEOUT", time.Second), "http-read timeout duration")
	httpWriteTimeout = flag.Duration("http-write-timeout", envDuration("HTTP_WRITE_TIMEOUT", time.Second), "http-write timeout duration")
	httpIdleTimeout  = flag.Duration("http-idle-timeout", envDuration("HTTP_IDLE_TIMEOUT", time.Minute), "http-idle timeout duration")

	httpMaxHeaderBytes = flag.Int("http-max-header-bytes", envInt("HTTP_MAX_HEADER_BYTES", 4098), "maximum of http-header bytes")
)

func makeServer(handler http.Handler, logger *log.Logger) http.Server {
	return http.Server{
		Addr:           fmt.Sprintf("%s:%d", *httpHost, *httpPort),
		Handler:        handler,
		ReadTimeout:    *httpReadTimeout,
		WriteTimeout:   *httpWriteTimeout,
		IdleTimeout:    *httpIdleTimeout,
		MaxHeaderBytes: *httpMaxHeaderBytes,
		ErrorLog:       logger,
	}
}

func makeRouter() http.Handler {
	router := mux.NewRouter()
	router.Path("/account").Methods(http.MethodGet).Handler(nil)
	return router
}

func (t *Server) ConnState(_ net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		t.connCounter.Add(1)
	case http.StateClosed:
		t.connCounter.Done()
	}
}

func New(logger *log.Logger) *Server {
	flag.Parse()
	var httpServer = makeServer(makeRouter(), logger)
	var server = Server{
		server: &httpServer,
	}
	httpServer.ConnState = server.ConnState
	return &server
}
