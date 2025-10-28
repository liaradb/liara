package orchestrator

import (
	"context"
)

type orchestrator struct {
	in      chan message
	runners map[runnerID]runner
}

func NewOrchestrator() orchestrator {
	return orchestrator{
		in: make(chan message),
		runners: map[runnerID]runner{
			"seed":    newRunner("seed", ""),
			"child00": newRunner("child00", "seed"),
			"child01": newRunner("child01", "seed"),
		},
	}
}

func (o *orchestrator) Run(ctx context.Context) {
	for _, runner := range o.runners {
		go runner.run(ctx, o)
	}
	go o.recieve(ctx)
	<-ctx.Done()
}

func (o *orchestrator) send(m message) {
	o.in <- m
}

func (o *orchestrator) recieve(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
		case m := <-o.in:
			runner, ok := o.runners[m.Destination]
			if !ok {
				continue
			}

			runner.send(m)
		}
	}
}
