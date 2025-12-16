module github.com/cardboardrobots/liara

go 1.25.5

replace github.com/cardboardrobots/eventsource_go => ../eventsource_go

require (
	github.com/cardboardrobots/baseerror v0.0.2
	github.com/cardboardrobots/eventsource_go v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.77.0
	google.golang.org/protobuf v1.36.11
)

require (
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251213004720-97cd9d5aeac2 // indirect
)
