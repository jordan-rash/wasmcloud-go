package client

import (
	"time"

	"github.com/jordan-rash/wasmcloud-go/kv"
	"github.com/nats-io/nats.go"
)

type Client struct {
	nc             *nats.Conn
	topicPrefix    string
	nsPrefix       string
	timeout        time.Duration
	auctionTimeout time.Duration
	jsDomain       string
	kvstore        nats.KeyValue
}

type ClientBuilder struct {
	nc             *nats.Conn
	topicPrefix    string
	nsPrefix       string
	timeout        time.Duration
	auctionTimeout time.Duration
	jsDomain       string
}

func New(nc *nats.Conn, options ...ClientBuilderOptions) *ClientBuilder {
	cb := &ClientBuilder{
		nc:             nc,
		topicPrefix:    "",
		nsPrefix:       "default",
		timeout:        time.Second * 2,
		auctionTimeout: time.Second * 5,
		jsDomain:       "",
	}
	for _, opt := range options {
		opt(cb)
	}
	return cb
}

type ClientBuilderOptions func(*ClientBuilder)

func WithTopicPrefix(inTopic string) ClientBuilderOptions {
	return func(bc *ClientBuilder) {
		bc.topicPrefix = inTopic
	}
}

func WithNSPrefix(inPrefix string) ClientBuilderOptions {
	return func(bc *ClientBuilder) {
		bc.nsPrefix = inPrefix
	}
}

func WithTimeout(inTimeout time.Duration) ClientBuilderOptions {
	return func(bc *ClientBuilder) {
		bc.timeout = inTimeout
	}
}

func WithAuctionTimeout(inTimeout time.Duration) ClientBuilderOptions {
	return func(bc *ClientBuilder) {
		bc.auctionTimeout = inTimeout
	}
}

func WithJsDomain(inJsDomain string) ClientBuilderOptions {
	return func(bc *ClientBuilder) {
		bc.jsDomain = inJsDomain
	}
}

func (cb ClientBuilder) Build() (*Client, error) {
	kvs, err := kv.GetKVStore(cb.nc, cb.nsPrefix, cb.jsDomain)
	if err != nil {
		return nil, err
	}

	c := Client{
		nc:             cb.nc,
		topicPrefix:    cb.topicPrefix,
		nsPrefix:       cb.nsPrefix,
		timeout:        cb.timeout,
		auctionTimeout: cb.auctionTimeout,
		jsDomain:       cb.jsDomain,
		kvstore:        kvs,
	}

	return &c, nil
}
