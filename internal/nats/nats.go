package nats

import (
	"net/url"

	ns "github.com/nats-io/nats-server/v2/server"
	nats "github.com/nats-io/nats.go"
)

var (
	NATS_URL          string = nats.DefaultURL
	NATS_SUBJ         string = "development"
	NATS_DURABLE_NAME string = "development"
)

func StartLeafNode() (*ns.Server, error) {
	opt := ns.Options{}
	opt.JetStream = true
	opt.Port = 0
	opt.Host = "localhost"
	opt.JetStream = true
	opt.JetStreamDomain = "core"
	opt.LeafNode = ns.LeafNodeOpts{
		Remotes: []*ns.RemoteLeafOpts{
			&ns.RemoteLeafOpts{
				URLs: []*url.URL{
					&url.URL{
						Host: "192.168.150.183:7422",
					},
				},
			},
		},
	}

	return ns.NewServer(&opt)
}
