FROM golang:1.27.0-alpine AS build

WORKDIR /build/modules/eventsource_go
COPY --from=eventsource_go go.mod go.sum ./

WORKDIR /build/services/liaradb
COPY go.mod go.sum ./
RUN go mod download

WORKDIR /build/modules/eventsource_go
COPY --from=eventsource_go . .

WORKDIR /build/services/liaradb
COPY . .

# Extract target OS and architecture injected automatically by Docker Buildx
ARG TARGETOS
ARG TARGETARCH

# Build the statically linked Go binary utilizing native cross-compilation
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -o bin/liaradb cmd/liaradb/main.go

FROM scratch AS app

WORKDIR /app
COPY --from=build /build/services/liaradb/bin/liaradb /build/services/liaradb/config.yml ./

EXPOSE 50051

CMD [ "/app/liaradb" ]
