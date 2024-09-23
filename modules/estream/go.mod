module github.com/cardboardrobots/estream

go 1.23.1

replace github.com/cardboardrobots/eventsource => ../eventsource

require (
	github.com/cardboardrobots/eventsource v0.0.0
	github.com/nats-io/nats.go v1.33.1
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
