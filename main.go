package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/jordan-rash/wasmcloud-go/internal/cli"
	inats "github.com/jordan-rash/wasmcloud-go/internal/nats"
	"github.com/jordan-rash/wasmcloud-go/wasmbus"
)

func main() {
	// Read CLI inputs and parse to struct
	host := cli.WasmcloudHost{Context: context.Background()}
	host.Parse()

	// Start wazero
	wb, err := wasmbus.NewWasmbus(host)
	if err != nil {
		panic(err)
	}

	s, err := inats.InitLeafNode(host)
	if err != nil {
		panic(err)
	}
	err = s.Start()
	if err != nil {
		panic(err)
	}

	err = s.StartSubscriptions(wb)
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	s.Close()
	// r.Close(ctx)
	// invSub.Drain()
	// nc.Close()
	os.Exit(1)
}
