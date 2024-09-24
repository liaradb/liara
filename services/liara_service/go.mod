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
	github.com/lib/pq v1.10.9
	google.golang.org/grpc v1.67.0
	modernc.org/sqlite v1.30.1
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.1.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/gc/v3 v3.0.0-20240107210532-573471604cb6 // indirect
	modernc.org/libc v1.52.1 // indirect
	modernc.org/mathutil v1.6.0 // indirect
	modernc.org/memory v1.8.0 // indirect
	modernc.org/strutil v1.2.0 // indirect
	modernc.org/token v1.1.0 // indirect
)
