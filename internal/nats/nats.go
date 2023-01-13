package nats

import (
	"net/url"
	"time"

	"github.com/jordan-rash/wasmcloud-go/broker"
	"github.com/jordan-rash/wasmcloud-go/internal/cli"
	"github.com/jordan-rash/wasmcloud-go/wasmbus"
	ns "github.com/nats-io/nats-server/v2/server"
	nats "github.com/nats-io/nats.go"
)

const (
	NATS_LEAF_HOST string = "localhost"
	NATS_LEAF_PORT int    = 0
)

var (
	NATS_URL          string = nats.DefaultURL
	NATS_SUBJ         string = "development"
	NATS_DURABLE_NAME string = "development"
)

type WasmcloudNats struct {
	ns   *ns.Server
	nc   *nats.Conn
	host cli.WasmcloudHost
	wb   *wasmbus.Wasmbus
}

func (w WasmcloudNats) Close() {
	w.nc.Drain()
	for w.nc.IsDraining() {
	}
	w.ns.Shutdown()
}

func (w *WasmcloudNats) Start() error {
	var err error

	go w.ns.Start()
	for !w.ns.Running() {
		time.Sleep(50 * time.Millisecond)
	}

	w.nc, err = nats.Connect("localhost", nats.InProcessServer(w.ns))
	if err != nil {
		return err
	}

	// Start Host Cloud Event
	subj, event := broker.Event{}.HostStart(
		w.host.WasmcloudLatticePrefix,
		w.host.HostId,
	)

	w.nc.Publish(subj, event)

	// Heartbeat every 30 seconds
	go func() {
		for {
			time.Sleep(30 * time.Second)

			subj, event := broker.Event{}.HostHeartbeat(
				w.host.WasmcloudLatticePrefix,
				w.host.HostId,
				w.host.Uptime,
			)
			w.nc.Publish(subj, event)
		}
	}()

	return nil
}

func InitLeafNode(host cli.WasmcloudHost) (WasmcloudNats, error) {
	opt := ns.Options{
		JetStream:       true,
		Port:            NATS_LEAF_PORT, // port 0 forces the server to use pipes instead of network
		Host:            NATS_LEAF_HOST, // must be localhost for embedded server
		JetStreamDomain: host.WasmcloudNatsJsDomain,
		LeafNode: ns.LeafNodeOpts{
			Remotes: []*ns.RemoteLeafOpts{
				&ns.RemoteLeafOpts{
					URLs: []*url.URL{
						&url.URL{
							Host: host.WasmcloudNatsRemoteUrl,
						},
					},
				},
			},
		},
	}
	n, err := ns.NewServer(&opt)
	return WasmcloudNats{ns: n, host: host}, err
}
