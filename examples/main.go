package main

import (
	"time"

	"github.com/jordan-rash/wasmcloud-go/client"
	"github.com/nats-io/nats.go"
)

var HOST_ID string

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	c := client.New(nc, "default", time.Second*2)

	hosts := c.GetHosts(time.Second * 1)
	HOST_ID = hosts[0].ID

	annotations := make(map[string]interface{})

	c.StartActor(
		HOST_ID,
		"wasmcloud.azurecr.io/echo:0.3.2",
		1,
		annotations,
	)

	c.StartProvider(
		HOST_ID,
		"wasmcloud.azurecr.io/httpserver:0.14.6",
		"default",
		annotations,
		"",
	)

	time.Sleep(time.Second * 10)
	hostInv := c.GetHostInventory(HOST_ID)

	linkValues := make(map[string]interface{})
	linkValues["address"] = "0.0.0.0:3000"
	c.AdvertiseLink(
		hostInv.Actors[0].ID,
		hostInv.Providers[0].ID,
		"wasmcloud:httpserver",
		"default",
		linkValues,
	)
}
