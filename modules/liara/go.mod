module github.com/cardboardrobots/liara

go 1.23.1

replace github.com/cardboardrobots/eventsource => ../../modules/eventsource

replace github.com/cardboardrobots/esgrpc => ../../modules/esgrpc

replace github.com/cardboardrobots/eventsource_go => ../../modules/eventsource_go

require github.com/cardboardrobots/eventsource v0.0.0

require github.com/google/uuid v1.6.0 // indirect
