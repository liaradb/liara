package handler

import (
	"errors"

	"github.com/cardboardrobots/liara/esgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type object string
type command string

const (
	objectTenant        object  = "tenant"
	objectEvent         object  = "event"
	commandTenantList   command = "list"
	commandTenantCreate command = "create"
	commandTenantDelete command = "delete"
	commandEventList    command = "list"
)

var (
	errNoCommand = errors.New("no command specified")
	errNoObject  = errors.New("no object specified")
)

type commandHandler struct {
	eventHandler  eventHandler
	tenantHandler tenantHandler
}

func NewCommandHandler(url string) (*commandHandler, error) {
	grpcConn, err := grpc.NewClient(url,
		grpc.WithTransportCredentials(
			insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &commandHandler{
		eventHandler: eventHandler{
			es: esgrpc.NewEventSourceGRPC(grpcConn)},
		tenantHandler: tenantHandler{
			es: esgrpc.NewEventSourceGRPC(grpcConn)},
	}, nil
}

func (ch *commandHandler) Handle(args []string) error {
	obj, cmd, args := ch.getArgs(args)
	switch obj {
	case objectTenant:
		return ch.tenantHandler.handle(cmd, args)
	case objectEvent:
		return ch.eventHandler.handle(cmd)
	default:
		return errNoObject
	}
}

func (ch *commandHandler) getArgs(args []string) (object, command, []string) {
	switch len(args) {
	case 0:
		fallthrough
	case 1:
		return "", "", nil
	case 2:
		return object(args[1]), "", nil
	default:
		return object(args[1]), command(args[2]), args[3:]
	}
}
