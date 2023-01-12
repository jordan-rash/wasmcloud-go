package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/jordan-rash/wasmcloud-go/internal/cli"
	inats "github.com/jordan-rash/wasmcloud-go/internal/nats"
)

func main() {
	// Read CLI inputs and parse to struct
	host := cli.WasmcloudHost{Context: context.Background()}
	host.Parse()

	s, err := inats.InitLeafNode(host)
	if err != nil {
		panic(err)
	}
	s.Start()

	err = s.StartSubscriptions()
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	// r.Close(ctx)
	// invSub.Drain()
	// nc.Close()
	os.Exit(1)
}
