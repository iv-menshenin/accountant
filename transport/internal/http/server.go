package http

import (
	"context"
	"flag"
	"fmt"
	"github.com/iv-menshenin/accountant/configstorage"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
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
	httpPort = flag.Int("http-port", configstorage.EnvInt("HTTP_PORT", 8080), "http-server port")
	httpHost = flag.String("http-host", configstorage.EnvString("HTTP_HOST", ""), "http-server host")

	httpReadTimeout  = flag.Duration("http-read-timeout", configstorage.EnvDuration("HTTP_READ_TIMEOUT", time.Second), "http-read timeout duration")
	httpWriteTimeout = flag.Duration("http-write-timeout", configstorage.EnvDuration("HTTP_WRITE_TIMEOUT", time.Second), "http-write timeout duration")
	httpIdleTimeout  = flag.Duration("http-idle-timeout", configstorage.EnvDuration("HTTP_IDLE_TIMEOUT", time.Minute), "http-idle timeout duration")

	httpMaxHeaderBytes = flag.Int("http-max-header-bytes", configstorage.EnvInt("HTTP_MAX_HEADER_BYTES", 4098), "maximum of http-header bytes")
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

func (t *Server) ConnState(_ net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		t.connCounter.Add(1)
	case http.StateClosed:
		t.connCounter.Done()
	}
}

func New(logger *log.Logger, rp RequestProcessor) *Server {
	flag.Parse()
	var httpServer = makeServer(makeRouter(rp), logger)
	var server = Server{
		server: &httpServer,
	}
	httpServer.ConnState = server.ConnState
	return &server
}
