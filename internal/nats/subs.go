package nats

import (
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jordan-rash/nkeys"
	"github.com/jordan-rash/wasmcloud-go/broker"
	"github.com/jordan-rash/wasmcloud-go/internal/oci"
	"github.com/jordan-rash/wasmcloud-go/models"
	"github.com/jordan-rash/wasmcloud-go/wasmbus"
	nats "github.com/nats-io/nats.go"
	core "github.com/wasmcloud/interfaces/core/tinygo"
	msgpack "github.com/wasmcloud/tinygo-msgpack"
)

type Claims struct {
	jwt.StandardClaims
	ID     string `json:"jti"`
	Wascap Wascap `json:"wascap"`
}

type Wascap struct {
	models.ActorDescription

	TargetURL string `json:"target_url"`
	OriginURL string `json:"origin_url"`
}

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
		ActorRef    string            `json:"actor_ref"`
		Count       int               `json:"count"`
		HostId      string            `json:"host_id"`
		Constraints map[string]string `json:"constraints"`
	}{}

	json.Unmarshal(m.Data, &req)

	splitOCI := strings.Split(req.ActorRef, ":")
	aB, metadata, err := oci.PullOCIRef(w.host.Context, splitOCI[0], splitOCI[1], w.host.Logger)
	if err != nil {
		panic(err)
	}

	var claims *Claims

	for _, i := range metadata.CustomSection {
		if i.Name == "jwt" {
			var Ed25519SigningMethod jwt.SigningMethodEd25519
			jwt.RegisterSigningMethod("Ed25519",
				func() jwt.SigningMethod { return &Ed25519SigningMethod })

			token, err := jwt.ParseWithClaims(
				string(i.Data),
				&Claims{},
				func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
						err := errors.New("invalid signing method")
						w.host.Logger.Error(err, "provided key was not valid ed25519")
						return nil, err
					}

					if claims, ok := token.Claims.(*Claims); ok {
						rKey, err := nkeys.Decode(
							nkeys.PrefixByteAccount,
							[]byte(claims.Issuer),
						)
						if err != nil {
							return nil, err
						}

						var pubKey ed25519.PublicKey = rKey
						return pubKey, nil
					}
					return nil, err
				})
			if err != nil {
				w.host.Logger.Error(err, "failed to validate jwt")
				return
			}

			var ok bool
			if claims, ok = token.Claims.(*Claims); !ok || !token.Valid {
				fmt.Println(err)
			}
		}
	}

	// TODO: need to compare wasm module hash, WITHOUT embedded signature, to
	// token.Claims.(*Claims).Wascap.Hash
	// sum := sha256.Sum256(aB)
	// fmt.Printf("SHA256: %x", sum)

	mod, err := w.wb.CreateModule(aB)
	if err != nil {
		panic(err)
	}

	subj, event := broker.Event{}.ActorStarted(broker.WASMCLOUD_DEFAULT_NSPREFIX, w.host.HostId, "")
	w.nc.Publish(subj, event)

	ai := models.ActorInstance{
		Annotations: req.Constraints,
		InstanceID:  claims.ID,
		Revision:    0,
	}

	a := models.ActorDescription{
		Name:      claims.Wascap.Name,
		Id:        claims.Subject,
		ImageRef:  req.ActorRef,
		Instances: models.ActorInstances{ai},
	}

	w.host.AddActor(a)

	w.nc.Subscribe("wasmbus.rpc.default."+claims.Subject, func(m *nats.Msg) {
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
