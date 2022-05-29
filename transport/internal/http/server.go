package http

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/iv-menshenin/accountant/transport/internal/http/ep"
	"github.com/iv-menshenin/accountant/utils"
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
	fmt.Printf("starting listening on '%s'", t.server.Addr)
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

func makeServer(config utils.Config, handler http.Handler, logger *log.Logger) http.Server {
	var (
		httpPort = config.IntegerConfig("http-port", "port", "PORT", 8080, "http-server port")
		httpHost = config.StringConfig("http-host", "host", "HOST", "", "http-server host")

		httpReadTimeout  = config.DurationConfig("http-read-timeout", "read-timeout", "READ_TIMEOUT", time.Second, "http-read timeout duration")
		httpWriteTimeout = config.DurationConfig("http-write-timeout", "write-timeout", "WRITE_TIMEOUT", time.Second, "http-write timeout duration")
		httpIdleTimeout  = config.DurationConfig("http-idle-timeout", "idle-timeout", "IDLE_TIMEOUT", time.Second, "http-idle timeout duration")

		httpMaxHeaderBytes = config.IntegerConfig("http-max-header-bytes", "max-header-bytes", "MAX_HEADER_BYTES", 4098, "maximum of http-header bytes")
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

func New(config utils.Config, logger *log.Logger, rp RequestProcessor, auth AuthCore) *Server {
	var httpServer = makeServer(config, makeRouter(rp, auth, logger), logger)
	var server = Server{
		server: &httpServer,
	}
	httpServer.ConnState = server.ConnState
	return &server
}
