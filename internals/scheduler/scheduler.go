package scheduler

import (
	"context"

	scheduling "github.com/codeshelldev/gotl/pkg/scheduler"
)

var scheduler = scheduling.New()
var cancel context.CancelFunc

func Start() {
	var ctx context.Context
	ctx, cancel = context.WithCancel(context.Background())

	go scheduler.Run(ctx)

	StartRequestScheduler()
}

func Stop() {
	cancel()
}