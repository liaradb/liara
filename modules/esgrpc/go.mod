module github.com/cardboardrobots/esgrpc

go 1.23.1

replace github.com/cardboardrobots/eventsource => ../../modules/eventsource

replace github.com/cardboardrobots/eventsource_go => ../../modules/eventsource_go

require (
	github.com/cardboardrobots/eventsource v0.0.0
	github.com/cardboardrobots/eventsource_go v0.0.0
	google.golang.org/grpc v1.67.0
	google.golang.org/protobuf v1.34.2
)

require (
	github.com/google/uuid v1.6.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
)
