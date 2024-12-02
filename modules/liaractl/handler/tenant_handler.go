package handler

import (
	"context"
	"fmt"

	"github.com/cardboardrobots/liara/esgrpc"
)

type tenantHandler struct {
	es *esgrpc.EventSourceGRPC
}

func (th *tenantHandler) handle(cmd command) error {
	switch cmd {
	case commandObjectList:
		return th.listTenants()
	default:
		return errNoCommand
	}
}

func (th *tenantHandler) listTenants() error {
	count := 0
	for _, err := range th.es.ListTenants(context.Background()) {
		if err != nil {
			return err
		}

		count++
	}
	if count == 0 {
		fmt.Println("no tenants")
	}
	return nil
}
