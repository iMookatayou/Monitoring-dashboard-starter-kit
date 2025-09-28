package services
import (
	"conntext"
	"log"
	"time"
)

type AppRegator struct {
	Tick time.Duration
	Work func(context.Context) error
}

func (a AppRegator) Start(ctx context.Context) {
 	t := time.NewTicker(a.Tick)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("AppRegator stopped"); return
		case <-t.C:
			if a.Work != nil { _ = a.Work(ctx) }
		}
	}
}