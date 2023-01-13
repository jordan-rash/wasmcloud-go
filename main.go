package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/bombsimon/logrusr/v4"
	"github.com/go-logr/logr"
	"github.com/jordan-rash/wasmcloud-go/internal/cli"
	inats "github.com/jordan-rash/wasmcloud-go/internal/nats"
	"github.com/jordan-rash/wasmcloud-go/wasmbus"
	"github.com/sirupsen/logrus"
)

var log logr.Logger

func init() {
	logrusLog := logrus.New()
	log = logrusr.New(logrusLog)
}

func main() {

	// Read CLI inputs and parse to struct
	host := cli.WasmcloudHost{Context: context.Background()}
	host.Parse()

	logLvl := func() int {
		if (host.Verbose * 2) <= 10 {
			return host.Verbose * 2
		} else {
			return 10
		}
	}
	switch logLvl() {
	case 2:
		log.Info("Log level: panic")
	case 4:
		log.Info("Log level: error")
	case 6:
		log.Info("Log level: info")
	case 8:
		log.Info("Log level: debug")
	case 10:
		log.Info("Log level: trace")
	}

	log.V(logLvl())
	host.Logger = log

	// Start wazero
	wb, err := wasmbus.NewWasmbus(host)
	if err != nil {
		panic(err)
	}

	// Initialize NATs leaf node
	s, err := inats.InitLeafNode(host)
	if err != nil {
		panic(err)
	}

	// Start NATs leaf node
	err = s.Start()
	if err != nil {
		panic(err)
	}

	err = s.StartSubscriptions(wb)
	if err != nil {
		panic(err)
	}

	log.Info(fmt.Sprintf("host id: %s", host.HostId))
	log.Info(fmt.Sprintf("cluster issuer: %s", host.Issuer))
	log.Info(fmt.Sprintf("cluster seed: %s", host.WasmcloudClusterSeed))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	s.Close()
	os.Exit(1)
}
