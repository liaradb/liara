module github.com/liaradb/liaradb

go 1.25.0

replace github.com/cardboardrobots/eventsource_go => ../../modules/eventsource_go

require (
	github.com/cardboardrobots/assert v0.0.2
	github.com/cardboardrobots/baseerror v0.0.2
	github.com/cardboardrobots/eventsource_go v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	google.golang.org/protobuf v1.36.10
)

require (
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
	google.golang.org/grpc v1.67.0 // indirect
)
