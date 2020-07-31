package main

import (
	"context"
	"log"
	"time"

	"github.com/elarian/elariango"
	elarian "github.com/elarian/elariango/com_elarian_hera_proto"
)

func main() {
	client, err := elariango.Initialize("some apikey", true)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	customerRequest := elarian.GetCustomerStateRequest{
		AppId: "soem id",
	}
	res, err := client.GetCustomerState(ctx, &customerRequest)
	if err != nil {
		log.Fatalf("could not get stat %v", err)
	}
	log.Println(res.GetData())

}
