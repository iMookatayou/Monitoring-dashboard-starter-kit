package services

import (
	"context"
	"log"
	"time"
)

type Aggregator struct {
	Tick time.Duration
	Work func(context.Context) error
}

func (a Aggregator) Start(ctx context.Context) {
	t := time.NewTicker(a.Tick)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("aggregator stop")
			return
		case <-t.C:
			if a.Work != nil {
				_ = a.Work(ctx)
			}
		}
	}
}
