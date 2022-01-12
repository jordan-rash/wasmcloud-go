package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/jordan-rash/wasmcloud-go/client"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

const (
	ECHO_ACTOR          = "wasmcloud.azurecr.io/echo:0.3.2"
	HTTPSERVER_PROVIDER = "wasmcloud.azurecr.io/httpserver:0.14.6"
)

var HOST_ID string

func main() {
	if len(os.Args[1:]) < 2 {
		log.Errorln("Input needs at least 2 arguments")
		os.Exit(1)
	}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	c := client.New(nc, "default", time.Second*2)
	annotations := make(map[string]string)

	hosts := c.GetHosts(time.Second * 1)
	HOST_ID = hosts[0].ID

	hostInv := c.GetHostInventory(HOST_ID)

	switch os.Args[1] {
	case "claims":
		switch os.Args[2] {
		case "ls":
			PrettyPrint(c.GetClaims())
		}
	case "actor":
		switch os.Args[2] {
		case "start":
			PrettyPrint(c.StartActor(
				HOST_ID,
				ECHO_ACTOR,
				1,
				annotations,
			))
		case "stop":
			PrettyPrint(c.StopActor(
				HOST_ID,
				ECHO_ACTOR,
				1,
				annotations,
			))
		case "update":
			PrettyPrint(c.UpdateActor(HOST_ID,
				hostInv.Actors[0].ID,
				ECHO_ACTOR,
				annotations,
			))
		case "ls":
			PrettyPrint(hostInv.Actors)
		default:
			log.Fatal("Invalid command")
		}
	case "provider":
		switch os.Args[2] {
		case "start":
			PrettyPrint(c.StartProvider(
				HOST_ID,
				HTTPSERVER_PROVIDER,
				"default",
				annotations,
				"",
			))
		case "stop":
			PrettyPrint(c.StopProvider(HOST_ID,
				hostInv.Providers[0].ID,
				"default",
				"wasmcloud:httpserver",
				annotations,
			))
		case "ls":
			PrettyPrint(hostInv.Providers)
		default:
			log.Fatal("Invalid command")
		}
	case "link":
		switch os.Args[2] {
		case "start":
			linkValues := make(map[string]string)
			linkValues["address"] = "0.0.0.0:3000"
			PrettyPrint(c.AdvertiseLink(
				hostInv.Actors[0].ID,
				hostInv.Providers[0].ID,
				"wasmcloud:httpserver",
				"default",
				linkValues,
			))
		case "stop":
			PrettyPrint(c.RemoveLink(
				hostInv.Actors[0].ID,
				"wasmcloud:httpserver",
				"default",
			))
		case "ls":
			PrettyPrint(c.QueryLinks())
		default:
			log.Fatal("Invalid command")
		}
	case "host":
		switch os.Args[2] {
		case "all":
			PrettyPrint(c.GetHosts(time.Second * 1))
		case "stop":
			PrettyPrint(c.StopHost(HOST_ID, time.Second*1))
		case "ls":
			PrettyPrint(c.GetHostInventory(HOST_ID))
		default:
			log.Fatal("Invalid command")
		}
	default:
		log.Fatal("Invalid command")
	}
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		log.Println(string(b))
	}
	return
}
