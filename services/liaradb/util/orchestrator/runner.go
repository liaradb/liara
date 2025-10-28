package orchestrator

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/liaradb/liaradb/util/mut"
)

type runner struct {
	name  runnerID
	in    chan message
	seed  runnerID
	peers mut.Set[runnerID]
}

func newRunner(name runnerID, seed runnerID) runner {
	var peers mut.Set[runnerID]
	if seed == "" {
		peers = mut.NewSet[runnerID]()
	} else {
		peers = mut.NewSet(seed)
	}

	return runner{
		name:  name,
		in:    make(chan message),
		seed:  seed,
		peers: peers,
	}
}

func (r *runner) send(m message) {
	r.in <- m
}

func (r *runner) run(ctx context.Context, parent reciever) {
	go r.heartbeat(ctx, parent)
	r.recieve(ctx)
}

func (r *runner) recieve(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
		case m := <-r.in:
			r.handle(m)
		}
	}
}

func (r *runner) handle(m message) {
	data, _ := json.MarshalIndent(m, "", "  ")
	log.Println(string(data))
	switch m.Name {
	case "runners":
		switch v := m.Value.(type) {
		case []runnerID:
			for _, n := range v {
				if n != r.name {
					r.peers.Add(n)
				}
			}
		}
	}
}

func (r *runner) heartbeat(ctx context.Context, parent reciever) {
	timer := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
		case <-timer.C:
			r.sendPeers(parent)
		}
	}
}

func (r *runner) sendPeers(parent reciever) {
	id, ok := r.getPeerID()
	if !ok {
		return
	}

	parent.send(message{
		Destination: id,
		Source:      r.name,
		Name:        "runners",
		Value:       r.getRunners(),
	})
}

func (r *runner) getRunners() []runnerID {
	return append(r.peers.Slice(), r.name)
}

func (r *runner) getPeerID() (runnerID, bool) {
	length := len(r.peers)
	if length == 0 {
		return "", false
	}

	i := rand.Int() % len(r.peers)
	index := 0
	for key := range r.peers {
		if i == index {
			return key, true
		}
		index++
	}

	return "", false
}
