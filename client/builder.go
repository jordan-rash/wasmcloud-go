// # Control Interface Client
//
// This library provides a client API for consuming the wasmCloud control interface over a
// NATS connection. This library can be used by multiple types of tools, and is also used
// by the control interface capability provider and the wash CLI
package client

import (
	"log"
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

// A client builder that can be used to fluently provide configuration settings used to construct
// the control interface client
type ClientBuilder struct {
	nc             *nats.Conn
	topicPrefix    string
	nsPrefix       string
	timeout        time.Duration
	auctionTimeout time.Duration
	jsDomain       string
}

// Creates a new client builder
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

// Sets the topic prefix for the NATS topic used for all control requests. Not to be confused with lattice ID/prefix
func WithTopicPrefix(inTopic string) ClientBuilderOption {
	return func(bc *ClientBuilder) {
		bc.topicPrefix = inTopic
	}
}

// The lattice ID/prefix used for this client. If this function is not invoked, the prefix will be set to `default`
func WithNSPrefix(inPrefix string) ClientBuilderOption {
	return func(bc *ClientBuilder) {
		bc.nsPrefix = inPrefix
	}
}

// Sets the timeout for standard calls and RPC invocations used by the client. If not set, the default will be 2 seconds
func WithTimeout(inTimeout time.Duration) ClientBuilderOption {
	return func(bc *ClientBuilder) {
		bc.timeout = inTimeout
	}
}

// Sets the timeout for auction (scatter/gather) operations. If not set, the default will be 5 seconds
func WithAuctionTimeout(inTimeout time.Duration) ClientBuilderOption {
	return func(bc *ClientBuilder) {
		bc.auctionTimeout = inTimeout
	}
}

// Sets the JetStream domain for this client, which can be critical for locating the right key-value bucket
// for lattice metadata storage. If this is skipped, then the JS domain will be ""
func WithJsDomain(inJsDomain string) ClientBuilderOption {
	return func(bc *ClientBuilder) {
		bc.jsDomain = inJsDomain
	}
}

// Completes the generation of a control interface client. This function will attempt
// to locate and attach to a metadata key-value bucket (`LATTICEDATA_{prefix}`) when starting
func (cb ClientBuilder) Build() (*Client, error) {
	kvs, err := kv.GetKVStore(cb.nc, cb.nsPrefix, cb.jsDomain)
	if err != nil {
		log.Print("kvstore not initalized, using legacy lattice communications")
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
