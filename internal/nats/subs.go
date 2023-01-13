package nats

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/jordan-rash/wasmcloud-go/broker"
	"github.com/jordan-rash/wasmcloud-go/internal/oci"
	"github.com/jordan-rash/wasmcloud-go/wasmbus"
	nats "github.com/nats-io/nats.go"
	core "github.com/wasmcloud/interfaces/core/tinygo"
	msgpack "github.com/wasmcloud/tinygo-msgpack"
)

func (w *WasmcloudNats) StartSubscriptions(wb *wasmbus.Wasmbus) error {
	var err error
	w.wb = wb

	// Command based topics
	w.nc.Subscribe(broker.Commands{}.StartActor(w.host.WasmcloudLatticePrefix, w.host.HostId), w.startActor)
	w.nc.Subscribe(broker.Commands{}.StopActor(w.host.WasmcloudLatticePrefix, w.host.HostId), func(*nats.Msg) {})
	w.nc.Subscribe(broker.Commands{}.StartProvider(w.host.WasmcloudLatticePrefix, w.host.HostId), func(*nats.Msg) {})
	w.nc.Subscribe(broker.Commands{}.StopProvider(w.host.WasmcloudLatticePrefix, w.host.HostId), func(*nats.Msg) {})
	w.nc.Subscribe(broker.Commands{}.UpdateActor(w.host.WasmcloudLatticePrefix, w.host.HostId), func(*nats.Msg) {})
	w.nc.Subscribe(broker.Commands{}.StopHost(w.host.WasmcloudLatticePrefix, w.host.HostId), func(*nats.Msg) {})

	// Query based topic
	w.nc.Subscribe(broker.Queries{}.HostInventory(w.host.WasmcloudLatticePrefix, w.host.HostId), w.hostInventory)
	w.nc.Subscribe(broker.Queries{}.Hosts(w.host.WasmcloudLatticePrefix), w.hostPing)
	w.nc.Subscribe(broker.Queries{}.LinkDefinitions(w.host.WasmcloudLatticePrefix), func(*nats.Msg) {})
	w.nc.Subscribe(broker.Queries{}.Claims(w.host.WasmcloudLatticePrefix), func(*nats.Msg) {})

	// Other topics
	w.nc.Subscribe(broker.ActorAuctionSubject(w.host.WasmcloudLatticePrefix), w.actorAuction)
	w.nc.Subscribe(broker.ProviderAuctionSubject(w.host.WasmcloudLatticePrefix), func(*nats.Msg) {})

	return err
}

func (w WasmcloudNats) hostInventory(m *nats.Msg) {
	inv, _ := json.Marshal(w.host)
	w.nc.Publish(m.Reply, inv)
}

func (w WasmcloudNats) actorAuction(m *nats.Msg) {
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
	}

}

func (w WasmcloudNats) hostPing(m *nats.Msg) {
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
}

func (w WasmcloudNats) startActor(m *nats.Msg) {
	req := struct {
		ActorRef string `json:"actor_ref"`
		Count    int    `json:"count"`
		HostId   string `json:"host_id"`
	}{}

	json.Unmarshal(m.Data, &req)

	splitOCI := strings.Split(req.ActorRef, ":")
	aB, _, err := oci.PullOCIRef(w.host.Context, splitOCI[0], splitOCI[1], w.host.Logger)
	if err != nil {
		panic(err)
	}

	// TODO: wascap stuff here

	mod, err := w.wb.CreateModule(aB)
	if err != nil {
		panic(err)
	}

	subj, event := broker.Event{}.ActorStarted(broker.WASMCLOUD_DEFAULT_NSPREFIX, w.host.HostId, "")
	w.nc.Publish(subj, event)
	w.host.Actors = append(w.host.Actors, "MAVJ6A57BJA2IYXCJC2DJ64XUUPATDC3RBITEW5XNI4C73JFDLVKB2YE")

	w.nc.Subscribe("wasmbus.rpc.default.MBCFOPM6JW2APJLXJD3Z5O4CN7CPYJ2B4FTKLJUR5YR5MITIU7HD3WD5", func(m *nats.Msg) {
		d := msgpack.NewDecoder(m.Data)
		i, _ := core.MDecodeInvocation(&d)
		w.wb.Data = m.Data
		w.wb.Operation = i.Operation

		_, err = mod.ExportedFunction("__guest_call").
			Call(w.host.Context, uint64(len([]byte(w.wb.Operation))), uint64(len(w.wb.Data)))
		if err != nil {
			log.Panicln(err)
		}

		aerr, _ := mod.Memory().Read(w.wb.GetGuestError())

		ir := core.InvocationResponse{
			Msg:           w.wb.GetGuestResponse(),
			InvocationId:  i.Id,
			Error:         string(aerr),
			ContentLength: uint64(len(w.wb.GetGuestResponse())),
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

}
