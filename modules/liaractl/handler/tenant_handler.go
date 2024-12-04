package handler

import (
	"context"
	"fmt"

	"github.com/cardboardrobots/liara"
	"github.com/cardboardrobots/liara/esgrpc"
)

type tenantHandler struct {
	es *esgrpc.EventSourceGRPC
}

func (th *tenantHandler) handle(cmd command, args []string) error {
	switch cmd {
	case commandTenantList:
		return th.listTenants()
	case commandTenantCreate:
		return th.createTenant(args)
	default:
		return errNoCommand
	}
}

func (th *tenantHandler) listTenants() error {
	count := 0
	for t, err := range th.es.ListTenants(context.Background()) {
		if err != nil {
			return err
		}

		fmt.Printf("%v\t%v\n", t.TenantId, t.Name)
		count++
	}
	if count == 0 {
		fmt.Println("no tenants")
	}
	return nil
}

func (th *tenantHandler) createTenant(args []string) error {
	name := ""
	if len(args) > 0 {
		name = args[0]
	}
	_, err := th.es.CreateTenant(context.Background(), liara.TenantName(name))
	return err
}
