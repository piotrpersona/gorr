package app

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_httpServers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	done, err := NewPprofHttpServer(8000).Run(ctx)
	assert.Nil(t, err)

	doneProm, err := NewPrometheusMetricsHttpServer(8001).Run(ctx)
	assert.Nil(t, err)

	<-done
	<-doneProm
	assert.True(t, true)
}
