package client

import (
	"time"

	"github.com/jordan-rash/wasmcloud-go/kv"
	"github.com/nats-io/nats.go"
)

const (
	DEFAULT_NS_PREFIX       string        = "default"
	DEFAULT_TOPIC_PREFIX    string        = ""
	DEFAULT_JS_DOMAIN       string        = ""
	DEFAULT_TIMEOUT         time.Duration = time.Second * 2
	DEFAULT_AUCTION_TIMEOUT time.Duration = time.Second * 5
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

func New(nc *nats.Conn, options ...ClientBuilderOption) *ClientBuilder {
	cb := &ClientBuilder{
		nc:             nc,
		topicPrefix:    DEFAULT_TOPIC_PREFIX,
		nsPrefix:       DEFAULT_NS_PREFIX,
		timeout:        DEFAULT_TIMEOUT,
		auctionTimeout: DEFAULT_AUCTION_TIMEOUT,
		jsDomain:       DEFAULT_JS_DOMAIN,
	}
	for _, opt := range options {
		opt(cb)
	}
	return cb
}

type ClientBuilderOption func(*ClientBuilder)

func WithTopicPrefix(inTopic string) ClientBuilderOption {
	return func(bc *ClientBuilder) {
		bc.topicPrefix = inTopic
	}
}

func WithNSPrefix(inPrefix string) ClientBuilderOption {
	return func(bc *ClientBuilder) {
		bc.nsPrefix = inPrefix
	}
}

func WithTimeout(inTimeout time.Duration) ClientBuilderOption {
	return func(bc *ClientBuilder) {
		bc.timeout = inTimeout
	}
}

func WithAuctionTimeout(inTimeout time.Duration) ClientBuilderOption {
	return func(bc *ClientBuilder) {
		bc.auctionTimeout = inTimeout
	}
}

func WithJsDomain(inJsDomain string) ClientBuilderOption {
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
