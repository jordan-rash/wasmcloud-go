package client_test

import (
	"testing"
	"time"

	"github.com/jordan-rash/wasmcloud-go/client"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

var nc *nats.Conn

func init() {
	opt := server.Options{
		JetStream:       true,
		Port:            4222,        // port 0 forces the server to use pipes instead of network
		Host:            "127.0.0.1", // must be localhost for embedded server
		JetStreamDomain: "tester",
	}

	n, err := server.NewServer(&opt)
	if err != nil {
		panic(err)
	}

	go n.Start()
	for !n.Running() {
		time.Sleep(50 * time.Millisecond)
	}

	nc, err = nats.Connect("127.0.0.1")
	if err != nil {
		panic(err)
	}
}

func TestCreateDefaultClient(t *testing.T) {
	client, err := client.New(nc).Build()
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
