module github.com/cardboardrobots/liaractl

go 1.25.0

replace github.com/cardboardrobots/eventsource_go => ../eventsource_go

replace github.com/cardboardrobots/liara => ../liara

require (
	github.com/cardboardrobots/liara v0.0.0
	google.golang.org/grpc v1.76.0
)

require (
	github.com/cardboardrobots/baseerror v0.0.2 // indirect
	github.com/cardboardrobots/eventsource_go v0.0.0-00010101000000-000000000000 // indirect
	github.com/google/uuid v1.6.0 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
