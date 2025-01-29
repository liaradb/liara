module github.com/cardboardrobots/liaractl

go 1.23.1

replace github.com/cardboardrobots/eventsource_go => ../eventsource_go

replace github.com/cardboardrobots/liara => ../liara

require (
	github.com/cardboardrobots/liara v0.0.0
	google.golang.org/grpc v1.67.1
)

require (
	github.com/cardboardrobots/baseerror v0.0.2 // indirect
	github.com/cardboardrobots/eventsource_go v0.0.0-00010101000000-000000000000 // indirect
	github.com/google/uuid v1.6.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
)
