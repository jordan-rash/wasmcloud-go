package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/bombsimon/logrusr/v4"
	"github.com/go-logr/logr"
	"github.com/jordan-rash/wasmcloud-go/internal/cli"
	inats "github.com/jordan-rash/wasmcloud-go/internal/nats"
	"github.com/jordan-rash/wasmcloud-go/wasmbus"
	"github.com/sirupsen/logrus"
)

func configureLogger(v int, structured bool) logr.Logger {
	logLvls := []logrus.Level{logrus.InfoLevel, logrus.DebugLevel, logrus.TraceLevel}
	logLvl := func() logrus.Level {
		if v < 3 {
			return logLvls[v]
		} else {
			return logrus.TraceLevel
		}
	}

	logrusLog := logrus.New()
	logrusLog.SetLevel(logLvl())
	if !structured {
		logrusLog.SetFormatter(&logrus.TextFormatter{
			DisableColors:             true,
			FullTimestamp:             true,
			EnvironmentOverrideColors: true,
		})
	} else {
		logrusLog.SetFormatter(&logrus.JSONFormatter{})
	}

	log := logrusr.New(
		logrusLog,
		logrusr.WithReportCaller(),
	)

	return log
}

func main() {
	// Read CLI inputs and parse to struct
	host := cli.WasmcloudHost{Context: context.Background()}
	host.Parse() // many defaults set with Validate function called here

	host.Logger = configureLogger(host.Verbose, host.WasmcloudStructuredLoggingEnabled)

	// Start wazero
	wb, err := wasmbus.NewWasmbus(host)
	if err != nil {
		host.Logger.Error(err, "failed to initialize wazero")
		os.Exit(1)
	}

	// Initialize NATs leaf node
	s, err := inats.InitLeafNode(host)
	if err != nil {
		host.Logger.Error(err, "failed to initialize nats")
		os.Exit(1)
	}

	// Start NATs leaf node
	err = s.Start()
	if err != nil {
		host.Logger.Error(err, "failed to start nats leaf node")
		os.Exit(1)
	}

	err = s.StartSubscriptions(wb)
	if err != nil {
		host.Logger.Error(err, "failed to start nats subscriptions")
		os.Exit(1)
	}

	host.Logger.V(0).Info(
		"wasmcloud-go host started",
		"host_id", host.HostId,
		"cluster_issuer", host.Issuer,
		"cluster_seed", host.WasmcloudClusterSeed,
	)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	s.Close()
	os.Exit(0)
}
