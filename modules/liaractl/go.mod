module github.com/cardboardrobots/liaractl

go 1.25.5

replace github.com/cardboardrobots/eventsource_go => ../eventsource_go

replace github.com/cardboardrobots/liara => ../liara

require (
	github.com/cardboardrobots/liara v0.0.0
	google.golang.org/grpc v1.78.0
)

require (
	github.com/cardboardrobots/baseerror v0.0.2 // indirect
	github.com/cardboardrobots/eventsource_go v0.0.0-00010101000000-000000000000 // indirect
	github.com/google/uuid v1.6.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251222181119-0a764e51fe1b // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
