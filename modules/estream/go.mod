module github.com/cardboardrobots/estream

go 1.25.0

replace github.com/cardboardrobots/eventsource => ../eventsource

replace github.com/cardboardrobots/liara => ../liara

require (
	github.com/cardboardrobots/liara v0.0.0
	github.com/nats-io/nats.go v1.37.0
)

require (
	github.com/cardboardrobots/baseerror v0.0.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.17.10 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
)
