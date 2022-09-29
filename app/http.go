package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/piotrpersona/gorr/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type httpServer struct {
	name string
	srv  *http.Server
}

func (s *httpServer) Run(parent context.Context) (done <-chan struct{}, err error) {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("http erver listen error, err: %s", err)
		}
	}()

	doneCh := make(chan struct{})

	go func() {
		select {
		case <-parent.Done():
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := s.srv.Shutdown(ctx); err != nil {
				log.Errorf("http server shutdown err: %s", err)
			}
			doneCh <- struct{}{}
		}
	}()

	done = doneCh
	return
}

func (s *httpServer) Name() string {
	return s.name
}

func NewPrometheusMetricsHttpServer(port int) Application {
	router := mux.NewRouter()
	router.Path("/prometheus").Handler(promhttp.Handler())
	return NewHttpServer(router, port, "prometheus")
}

func NewPprofHttpServer(port int) Application {
	router := mux.NewRouter()
	router.PathPrefix("/debug/pprof").Handler(http.DefaultServeMux)
	return NewHttpServer(router, port, "pprof")
}

func NewHttpServer(router *mux.Router, port int, name string) Application {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}
	return &httpServer{srv: srv, name: name}
}
