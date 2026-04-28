module github.com/cardboardrobots/liara

go 1.26.2

replace github.com/cardboardrobots/eventsource_go => ../eventsource_go

require (
	github.com/cardboardrobots/baseerror v0.0.2
	github.com/cardboardrobots/eventsource_go v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.80.0
	google.golang.org/protobuf v1.36.11
)

require (
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260427160629-7cedc36a6bc4 // indirect
)
