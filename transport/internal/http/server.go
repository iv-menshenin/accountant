package http

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/transport/internal/http/ep"
)

type (
	Server struct {
		err         error
		connCounter sync.WaitGroup
		server      *http.Server
		stopOnce    sync.Once
	}
	AuthCore interface {
		ep.AuthProcessor
		Middleware() func(h http.Handler) http.Handler
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

func makeServer(config model.Config, handler http.Handler, logger *log.Logger) http.Server {
	var (
		httpPort = config.IntegerConfig("http-port", "http-port", "HTTP_PORT", 8080, "http-server port")
		httpHost = config.StringConfig("http-host", "http-host", "HTTP_HOST", "", "http-server host")

		httpReadTimeout  = config.DurationConfig("http-read-timeout", "http-read-timeout", "HTTP_READ_TIMEOUT", time.Second, "http-read timeout duration")
		httpWriteTimeout = config.DurationConfig("http-write-timeout", "http-write-timeout", "HTTP_WRITE_TIMEOUT", time.Second, "http-write timeout duration")
		httpIdleTimeout  = config.DurationConfig("http-idle-timeout", "http-idle-timeout", "HTTP_IDLE_TIMEOUT", time.Second, "http-idle timeout duration")

		httpMaxHeaderBytes = config.IntegerConfig("http-max-header-bytes", "http-max-header-bytes", "HTTP_MAX_HEADER_BYTES", 4098, "maximum of http-header bytes")
	)
	if err := config.Init(); err != nil {
		panic(err)
	}
	return http.Server{
		Addr:           fmt.Sprintf("%s:%d", *httpHost, *httpPort),
		Handler:        handler,
		ReadTimeout:    *httpReadTimeout,
		WriteTimeout:   *httpWriteTimeout,
		IdleTimeout:    *httpIdleTimeout,
		MaxHeaderBytes: int(*httpMaxHeaderBytes),
		ErrorLog:       logger,
	}
}

func (t *Server) ConnState(_ net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		t.connCounter.Add(1)
	case http.StateClosed:
		t.connCounter.Done()
	}
}

func New(config model.Config, logger *log.Logger, rp RequestProcessor, auth AuthCore) *Server {
	var httpServer = makeServer(config, makeRouter(rp, auth), logger)
	var server = Server{
		server: &httpServer,
	}
	httpServer.ConnState = server.ConnState
	return &server
}
