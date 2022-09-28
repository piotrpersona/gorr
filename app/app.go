package app

import "context"

// Application will Run until done or error is returned.
type Application interface {
	Run(ctx context.Context) (done <-chan struct{}, err error)
	Name() string
}
