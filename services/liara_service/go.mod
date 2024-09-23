module github.com/cardboardrobots/liara_service

go 1.23.1

replace github.com/cardboardrobots/eventsource => ../../modules/eventsource

replace github.com/cardboardrobots/esgrpc => ../../modules/esgrpc

replace github.com/cardboardrobots/eventsource_go => ../../modules/eventsource_go

require (
	github.com/cardboardrobots/baseerror v0.0.2
	github.com/cardboardrobots/config v0.0.0
	github.com/cardboardrobots/esgrpc v0.0.0
	github.com/cardboardrobots/eventsource v0.0.0
	github.com/cardboardrobots/eventsource_go v0.0.0
	github.com/cardboardrobots/listener v0.0.1
	google.golang.org/grpc v1.67.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.1.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
