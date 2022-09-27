package app

import (
	"context"
	"sync"

	"github.com/piotrpersona/gorr/log"
)

type supervisor struct {
	apps []Application
}

func NewSupervisor(apps ...Application) Application {
	return &supervisor{apps: apps}
}

func (s *supervisor) Run(ctx context.Context) (done <-chan struct{}, err error) {
	doneCh := make(chan struct{})

	go func() {
		var wg sync.WaitGroup
		for _, app := range s.apps {
			wg.Add(1)
			go func(app Application) {
				defer wg.Done()
				name := app.Name()
				log.Infof("app '%s' started", name)
				appDone, err := app.Run(ctx)
				if err != nil {
					return
				}
				<-appDone
				log.Infof("app '%s' terminated gracefully", name)
			}(app)
			wg.Wait()
			doneCh <- struct{}{}
		}
	}()

	done = doneCh
	return
}

func (s *supervisor) Name() string {
	return "supervisor"
}
