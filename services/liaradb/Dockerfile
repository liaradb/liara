FROM golang:1.25.4-alpine AS build

WORKDIR /build/modules/eventsource_go
COPY --from=eventsource_go go.mod go.sum ./

WORKDIR /build/services/liaradb
COPY go.mod go.sum ./
RUN go mod download

WORKDIR /build/modules/eventsource_go
COPY --from=eventsource_go . .

WORKDIR /build/services/liaradb
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/liaradb cmd/main.go

FROM scratch AS app

WORKDIR /app
COPY --from=build /build/services/liaradb/bin/liaradb /build/services/liaradb/config.yml ./

EXPOSE 50051

CMD [ "/app/liaradb" ]
