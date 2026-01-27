package address_test

import (
	"context"
	"testing"
	"time"

	"go.anx.io/go-anxcloud/pkg"
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/ipam/address"
)

func TestCreateUpdateDelete(t *testing.T) {
	//export ANEXIA_TOKEN='YOURTOKEN'
	//consts
	const prefix = ""
	const name = "" //ip addr to add
	const rdnsName = ""
	const descriptionCustomer = ""

	//init
	client, err := client.New(client.TokenFromEnv(false))
	if err != nil {
		t.Fatal(err)
	}
	a := pkg.NewAPI(client)
	ctx := context.Background()

	//create
	create := address.NewCreate(prefix, name)
	create.RDNSName = rdnsName
	addr, err := a.IPAM().Address().Create(ctx, create)
	if err != nil {
		t.Fatal(err)
	}

	for {
		addr2, err2 := a.IPAM().Address().Get(ctx, addr.ID)
		if err2 != nil {
			t.Fatal(err2)
		}
		if addr2.Status == "Inaktiv" {
			break
		}
		time.Sleep(time.Second * 5)
	}

	//update
	update := address.Update{
		DescriptionCustomer: descriptionCustomer,
		RDNSName:            rdnsName,
	}
	_, err = a.IPAM().Address().Update(ctx, addr.ID, update)
	if err != nil {
		t.Fatal(err)
	}

	//delete
	err = a.IPAM().Address().Delete(ctx, addr.ID)
	if err != nil {
		t.Fatal(err)
	}
}
