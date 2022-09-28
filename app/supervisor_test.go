package app

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Supervisor(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)
	defer cancel()

	done, err := NewSupervisor(WithPprof(3333), WithPrometheus(4444)).Run(ctx)
	assert.Nil(t, err)
	<-done
}
