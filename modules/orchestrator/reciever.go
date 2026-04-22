package orchestrator

type message struct {
	Destination runnerID `json:"destination"`
	Source      runnerID `json:"source"`
	Name        string   `json:"name"`
	Value       any      `json:"value"`
}

type reciever interface {
	send(message)
}

type runnerID string
