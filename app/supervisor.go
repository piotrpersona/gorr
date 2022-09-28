package app

import (
	"context"
	"sync"

	"github.com/piotrpersona/gorr/log"
)

type runCfg struct {
	logLevel log.Level
	apps     []Application
}

type Option func(*runCfg)

func WithLogLevel(level log.Level) Option {
	return func(o *runCfg) {
		o.logLevel = level
	}
}

func WithPprof(port int) Option {
	return WithApps(NewPprofHttpServer(port))
}

func WithPrometheus(port int) Option {
	return WithApps(NewPrometheusMetricsHttpServer(port))
}

func WithApps(apps ...Application) Option {
	return func(o *runCfg) {
		o.apps = append(o.apps, apps...)
	}
}

type supervisor struct {
	apps     []Application
	logLevel log.Level
}

func NewSupervisor(opts ...Option) Application {
	cfg := &runCfg{apps: make([]Application, 0)}
	for _, opt := range opts {
		opt(cfg)
	}
	super := &supervisor{apps: cfg.apps, logLevel: cfg.logLevel}
	return super
}

func (s *supervisor) Run(parent context.Context) (done <-chan struct{}, err error) {
	log.Init(parent, s.logLevel)

	doneCh := make(chan struct{})

	go func() {
		name := s.Name()
		s.logStarted(name)
		var wg sync.WaitGroup
		for _, app := range s.apps {
			wg.Add(1)
			go func(app Application) {
				defer wg.Done()
				appName := app.Name()
				s.logStarted(appName)
				appDone, err := app.Run(parent)
				if err != nil {
					return
				}
				s.logStopped(appName)
				<-appDone
			}(app)
		}
		wg.Wait()
		<-log.Sync()
		s.logStopped(name)
		doneCh <- struct{}{}
	}()

	done = doneCh
	return
}

func (s *supervisor) logStarted(name string) {
	log.Infof("app '%s' started", name)
}

func (s *supervisor) logStopped(name string) {
	log.Infof("app '%s' terminated gracefully", name)
}

func (s *supervisor) Name() string {
	return "supervisor"
}
