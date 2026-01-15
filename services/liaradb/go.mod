module github.com/liaradb/liaradb

go 1.25.5

replace github.com/liaradb/eventsource_go => ../../modules/eventsource_go

require (
	github.com/cardboardrobots/baseerror v0.0.2
	github.com/cardboardrobots/config v0.0.0
	github.com/cardboardrobots/errormap v0.0.0
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.3
	github.com/liaradb/eventsource_go v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.78.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260114163908-3f89685c29c3 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
