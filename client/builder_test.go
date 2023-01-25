package client

import (
	"bytes"
	"testing"
	"time"

	"github.com/bombsimon/logrusr/v4"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var nc *nats.Conn

func init() {
	opt := server.Options{
		JetStream:       true,
		Port:            4222,
		Host:            "127.0.0.1",
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

func TestNewDefaultBuilder(t *testing.T) {
	client, err := New(nc).Build()
	assert.NoError(t, err)

	assert.Equal(t, "default", client.nsPrefix)
	assert.Equal(t, "", client.topicPrefix)
	assert.Equal(t, "", client.jsDomain)
	assert.Equal(t, time.Second*2, client.timeout)
	assert.Equal(t, time.Second*5, client.auctionTimeout)
}

func TestNewBuilderWithoutKvStore(t *testing.T) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	log := logrusr.New(logger)

	client, err := New(nc, WithJsDomain("derp"), WithLogger(log)).Build()
	assert.NoError(t, err)
	assert.NotNil(t, client)

	assert.Contains(t, buf.String(), "kvstore not initalized, using legacy lattice communications")
	assert.Nil(t, client.kvstore)
}

func TestNewBuilderWithKvStore(t *testing.T) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	log := logrusr.New(logger)

	client, err := New(nc, WithLogger(log)).Build()
	assert.NoError(t, err)
	assert.NotNil(t, client)

	assert.NotContains(t, buf.String(), "kvstore not initalized, using legacy lattice communications")
	assert.NotNil(t, client.kvstore)
}

func TestNewBuilderWithTopicPrefix(t *testing.T) {
	client, err := New(nc, WithTopicPrefix("derp")).Build()
	assert.NoError(t, err)
	assert.Equal(t, "derp", client.topicPrefix)
}

func TestNewBuilderWithNSPrefix(t *testing.T) {
	client, err := New(nc, WithNSPrefix("derp")).Build()
	assert.NoError(t, err)
	assert.Equal(t, "derp", client.nsPrefix)
}

func TestNewBuilderWithTimeout(t *testing.T) {
	client, err := New(nc, WithTimeout(1*time.Minute)).Build()
	assert.NoError(t, err)
	assert.Equal(t, 1*time.Minute, client.timeout)
}

func TestNewBuilderWithAuctionTimeout(t *testing.T) {
	client, err := New(nc, WithAuctionTimeout(1*time.Minute)).Build()
	assert.NoError(t, err)
	assert.Equal(t, 1*time.Minute, client.auctionTimeout)
}

func TestNewBuilderWithJsDomain(t *testing.T) {
	client, err := New(nc, WithJsDomain("derp")).Build()
	assert.NoError(t, err)
	assert.Equal(t, "derp", client.jsDomain)
}
