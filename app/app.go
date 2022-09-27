package app

import "context"

type Application interface {
	Run(ctx context.Context) (done <-chan struct{}, err error)
	Name() string
}
