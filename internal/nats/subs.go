package nats

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/jordan-rash/wasmcloud-go/broker"
	"github.com/jordan-rash/wasmcloud-go/internal/oci"
	iwazero "github.com/jordan-rash/wasmcloud-go/internal/wazero"
	nats "github.com/nats-io/nats.go"
	core "github.com/wasmcloud/interfaces/core/tinygo"
	msgpack "github.com/wasmcloud/tinygo-msgpack"
)

func (w WasmcloudNats) StartSubscriptions() error {
	var err error

	// Provide Inventory request response
	_, err = w.nc.Subscribe(
		broker.Queries{}.HostInventory(
			broker.WASMCLOUD_DEFAULT_NSPREFIX,
			w.host.HostId), func(m *nats.Msg) {
			inv, _ := json.Marshal(w.host)
			w.nc.Publish(m.Reply, inv)
		},
	)

	w.nc.Subscribe("wasmbus.ctl.default.cmd."+w.host.HostId+".la", func(m *nats.Msg) {
		req := struct {
			ActorRef string `json:"actor_ref"`
			Count    int    `json:"count"`
			HostId   string `json:"host_id"`
		}{}

		json.Unmarshal(m.Data, &req)

		splitOCI := strings.Split(req.ActorRef, ":")
		aB, _, err := oci.PullOCIRef(w.host.Context, splitOCI[0], splitOCI[1])
		if err != nil {
			panic(err)
		}

		// TODO: wascap stuff here

		actor := iwazero.Actor{
			Context:    w.host.Context,
			ActorBytes: aB,
		}

		err = actor.StartActor()
		if err != nil {
			log.Panicln(err)
		}

		subj, event := broker.Event{}.ActorStarted(broker.WASMCLOUD_DEFAULT_NSPREFIX, w.host.HostId, "")
		w.nc.Publish(subj, event)
		w.host.Actors = append(w.host.Actors, "MAVJ6A57BJA2IYXCJC2DJ64XUUPATDC3RBITEW5XNI4C73JFDLVKB2YE")

		w.nc.Subscribe("wasmbus.rpc.default.MBCFOPM6JW2APJLXJD3Z5O4CN7CPYJ2B4FTKLJUR5YR5MITIU7HD3WD5", func(m *nats.Msg) {
			d := msgpack.NewDecoder(m.Data)
			i, _ := core.MDecodeInvocation(&d)
			actor.Data = m.Data
			actor.Operation = i.Operation

			_, err = actor.GetModule().ExportedFunction("__guest_call").
				Call(w.host.Context, uint64(len([]byte(actor.Operation))), uint64(len(actor.Data)))
			if err != nil {
				log.Panicln(err)
			}

			aerr, _ := actor.GetModule().Memory().Read(actor.GetGuestError())

			ir := core.InvocationResponse{
				Msg:           actor.GetGuestResponse(),
				InvocationId:  i.Id,
				Error:         string(aerr),
				ContentLength: uint64(len(actor.GetGuestResponse())),
			}

			var sizer msgpack.Sizer
			size_enc := &sizer
			ir.MEncode(size_enc)
			buf := make([]byte, sizer.Len())
			encoder := msgpack.NewEncoder(buf)
			enc := &encoder
			ir.MEncode(enc)

			w.nc.Publish(m.Reply, buf)
		})

		ack := struct {
			Accepted bool   `json:"accepted"`
			Error    string `json:"error"`
		}{
			true, "",
		}

		sAck, _ := json.Marshal(ack)
		w.nc.Publish(m.Reply, sAck)

	})

	w.nc.Subscribe("wasmbus.ctl.default.ping.hosts", func(m *nats.Msg) {
		ack := struct {
			ClusterIssuers string            `json:"cluster_issuers"`
			CtlHost        string            `json:"ctl_host"`
			Friendly       string            `json:"friendly_name"`
			Id             string            `json:"id"`
			JsDomain       string            `json:"js_domain"`
			Labels         map[string]string `json:"labels"`
			LatticePrefix  string            `json:"lattice_prefix"`
			ProvRpcHost    string            `json:"prov_rpc_host"`
			RpcHost        string            `json:"rpc_host"`
			UptimeHuman    string            `json:"uptime_human"`
			Uptime         int               `json:"uptime_seconds"`
			Version        string            `json:"version"`
		}{
			ClusterIssuers: strings.Join(w.host.WasmcloudClusterIssuers, ","),
			CtlHost:        w.host.WasmcloudCtlHost,
			Friendly:       "yolo-bro-1234",
			Id:             w.host.HostId,
			JsDomain:       w.host.WasmcloudNatsJsDomain,
			Labels:         w.host.Labels,
			LatticePrefix:  w.host.WasmcloudLatticePrefix,
			ProvRpcHost:    w.host.WasmcloudProvRpcHost,
			RpcHost:        w.host.WasmcloudRpcHost,
			UptimeHuman:    "FOREVER",
			Uptime:         int(time.Since(w.host.Uptime).Seconds()),
			Version:        w.host.Version + "-wasmcloud_go",
		}

		sAck, _ := json.Marshal(ack)
		w.nc.Publish(m.Reply, sAck)

	})

	w.nc.Subscribe("wasmbus.ctl.default.auction.actor", func(m *nats.Msg) {
		auction := struct {
			ActorRef    string            `json:"actor_ref"`
			Constraints map[string]string `json:"constraints"`
		}{}
		json.Unmarshal(m.Data, &auction)

		if auction.Constraints["hostcore.version"] == "v0.0.0-wasmcloud_go" {
			auctionResponse := struct {
				ActorRef    string            `json:"actor_ref"`
				Constraints map[string]string `json:"constraints"`
				HostId      string            `json:"host_id"`
			}{
				ActorRef:    auction.ActorRef,
				Constraints: auction.Constraints,
				HostId:      w.host.HostId,
			}

			aRb, _ := json.Marshal(auctionResponse)
			w.nc.Publish(m.Reply, aRb)

			// pull actor here
		}
	})

	return err
}

// Provide Inventory request response
