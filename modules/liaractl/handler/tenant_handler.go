package handler

import (
	"context"
	"errors"
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
	case commandTenantDelete:
		return th.deleteTenant(args)
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
	id, name := "", ""
	switch len(args) {
	case 0:
		break
	case 1:
		id = args[0]
	default:
		id, name = args[0], args[1]
	}

	tenantID, err := th.es.CreateTenant(context.Background(),
		liara.TenantID(id),
		liara.TenantName(name))
	if err != nil {
		return err
	}

	fmt.Println(tenantID.String())

	return nil
}

func (th *tenantHandler) deleteTenant(args []string) error {
	if len(args) == 0 {
		return errors.New("no tenantID")
	}
	id := args[0]
	return th.es.DeleteTenant(context.Background(), liara.TenantID(id))
}
