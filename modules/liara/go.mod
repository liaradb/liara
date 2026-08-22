module github.com/cardboardrobots/liara

go 1.27.0

replace github.com/cardboardrobots/eventsource_go => ../eventsource_go

require (
	github.com/cardboardrobots/baseerror v0.0.2
	github.com/cardboardrobots/eventsource_go v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.83.1
	google.golang.org/protobuf v1.36.12
)

require (
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260819154853-08b0e4226688 // indirect
)
