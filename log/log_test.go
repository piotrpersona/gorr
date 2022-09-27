package log

import (
	"context"
	"testing"
	"time"
)

func Test_Log(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*10)
	Init(ctx, LevelInfo)
	Infof("Hey %s!", "Adam")
	<-Sync()
}
