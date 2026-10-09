module github.com/cardboardrobots/liara

go 1.27.2

replace github.com/cardboardrobots/eventsource_go => ../eventsource_go

require (
	github.com/cardboardrobots/baseerror v0.0.2
	github.com/cardboardrobots/eventsource_go v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.84.0
	google.golang.org/protobuf v1.36.12
)

require (
	golang.org/x/net v0.60.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20261005182115-fad411399dd8 // indirect
)
