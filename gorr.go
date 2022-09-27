package gorr

import (
	"context"
	"fmt"
	"net/http"
	"time"

	_ "net/http/pprof"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Application interface {
	Run(ctx context.Context) (done <-chan struct{}, err error)
	Name() string
}

type httpServer struct {
	name string
	srv  *http.Server
}

func (s *httpServer) Run(parent context.Context) (done <-chan struct{}, err error) {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("http erver listen error", err)
		}
	}()

	doneCh := make(chan struct{})

	go func() {
		select {
		case <-parent.Done():
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := s.srv.Shutdown(ctx); err != nil {
				fmt.Println("http server shutdown error", err)
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
	router.Path("/promethues").Handler(promhttp.Handler())

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	return &httpServer{srv: srv, name: "prometheus"}
}

func NewPprofHttpServer(port int) Application {
	router := mux.NewRouter()
	router.Path("/debug/pprof").Handler(http.DefaultServeMux)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	return &httpServer{srv: srv, name: "pprof"}
}
