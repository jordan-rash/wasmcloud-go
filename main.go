package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/jordan-rash/wasmcloud-go/broker"
	"github.com/jordan-rash/wasmcloud-go/internal/cli"
	inats "github.com/jordan-rash/wasmcloud-go/internal/nats"
	"github.com/jordan-rash/wasmcloud-go/internal/oci"
	"github.com/jordan-rash/wasmcloud-go/wasmbus"

	"github.com/alecthomas/kong"
	"github.com/nats-io/nats.go"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"

	core "github.com/wasmcloud/interfaces/core/tinygo"
	msgpack "github.com/wasmcloud/tinygo-msgpack"
)

//embed:go.wasm
var wasmMod []byte

var mod api.Module
var tempInv core.Invocation
var tempData []byte
var resp []byte

func main() {
	cli := cli.Context{Context: context.Background()}
	_ = kong.Parse(&cli,
		kong.Name("wasmcloud_go"),
		kong.Description("wasmcloud host implementation written in Go using Wazero"),
	)

	uptime := time.Now()
	s, err := inats.StartLeafNode()
	if err != nil {
		panic(err)
	}

	go s.Start()
	for {
		if s.Running() {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	ctx, cancel := context.WithCancel(cli.Context)
	defer cancel()

	nc, err := nats.Connect("localhost", nats.InProcessServer(s))
	if err != nil {
		panic(err)
	}

	wch := WasmcloudHost{
		HostId:    cli.HostId,
		Issuer:    "CDCR323CUXYTZAF4ODMPULKLHFIKYLN3YAMXYN6PP32JHNARBPJ5DICE",
		Labels:    map[string]string{"hostcore.arch": "x86_64", "hostcore.os": "ubuntu", "hostcore.osfamily": "unix", "hostcore.version": "v0.0.0-wasmcloud_go", "hostcore.runtime": "wazero"},
		Friendly:  "yolo-bro-1234",
		Actors:    []string{},
		Providers: []string{},
	}

	// Provide Inventory request response
	invSub, _ := nc.Subscribe(
		broker.Queries{}.HostInventory(
			broker.WASMCLOUD_DEFAULT_NSPREFIX,
			cli.HostId), func(m *nats.Msg) {
			inv, _ := json.Marshal(wch)
			nc.Publish(m.Reply, inv)
		},
	)

	// Start Host Cloud Event
	subj, event := broker.Event{}.HostStart(
		broker.WASMCLOUD_DEFAULT_NSPREFIX,
		cli.HostId,
	)
	b, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}
	nc.Publish(subj, b)

	// Heartbeat every 30 seconds
	go func() {
		for {
			time.Sleep(30 * time.Second)

			subj, event := broker.Event{}.HostHeartbeat(
				broker.WASMCLOUD_DEFAULT_NSPREFIX,
				cli.HostId,
				uptime,
			)
			b, err := json.Marshal(event)
			if err != nil {
				panic(err)
			}
			nc.Publish(subj, b)
		}
	}()

	r := wazero.NewRuntime(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, r)
	wasmbus := wasmbus.Wasmbus{}
	_, err = r.NewHostModuleBuilder("wasmbus").
		NewFunctionBuilder().WithFunc(wasmbus.HostCall).Export("__host_call").
		NewFunctionBuilder().WithFunc(wasmbus.ConsoleLog).Export("__console_log").
		NewFunctionBuilder().WithFunc(GuestRequest).Export("__guest_request").
		NewFunctionBuilder().WithFunc(wasmbus.HostResponse).Export("__host_response").
		NewFunctionBuilder().WithFunc(wasmbus.HostResponseLen).Export("__host_response_len").
		NewFunctionBuilder().WithFunc(GuestResponse).Export("__guest_response").
		NewFunctionBuilder().WithFunc(wasmbus.GuestError).Export("__guest_error").
		NewFunctionBuilder().WithFunc(wasmbus.HostError).Export("__host_error").
		NewFunctionBuilder().WithFunc(wasmbus.HostErrorLen).Export("__host_error_len").
		Instantiate(ctx, r)
	if err != nil {
		log.Panicln(err)
	}

	nc.Subscribe("wasmbus.ctl.default.cmd."+cli.HostId+".la", func(m *nats.Msg) {
		ack := struct {
			Accepted bool   `json:"accepted"`
			Error    string `json:"error"`
		}{
			true, "",
		}

		sAck, _ := json.Marshal(ack)
		nc.Publish(m.Reply, sAck)

	})

	nc.Subscribe("wasmbus.ctl.default.ping.hosts", func(m *nats.Msg) {
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
			ClusterIssuers: "CDCR323CUXYTZAF4ODMPULKLHFIKYLN3YAMXYN6PP32JHNARBPJ5DICE",
			CtlHost:        "127.0.0.1",
			Friendly:       "yolo-bro-1234",
			Id:             cli.HostId,
			JsDomain:       "",
			Labels:         map[string]string{"hostcore.arch": "armv61", "hostcore.os": "raspbian", "hostcore.osfamily": "unix", "hostcore.version": "v0.0.0-wasmcloud_go", "hostcore.runtime": "wazero"},
			LatticePrefix:  "default",
			ProvRpcHost:    "127.0.0.1",
			RpcHost:        "127.0.0.1",
			UptimeHuman:    "FOREVER",
			Uptime:         int(time.Since(uptime).Seconds()),
			Version:        "v0.0.0-wasmcloud_go",
		}

		sAck, _ := json.Marshal(ack)
		nc.Publish(m.Reply, sAck)

	})

	// {"cluster_issuers":"CDCR323CUXYTZAF4ODMPULKLHFIKYLN3YAMXYN6PP32JHNARBPJ5DICE","ctl_host":"127.0.0.1","friendly_name":"aged-shape-5351","id":"NDTCAWUH7TECOYLOOP3E7SDQBKR6TU45G3L56GEDKGFLJL2L5N6IK6UL","js_domain":null,"labels":{"hostcore.arch":"x86_64","hostcore.os":"macos","hostcore.osfamily":"unix"},"lattice_prefix":"default","prov_rpc_host":"127.0.0.1","rpc_host":"127.0.0.1","uptime_human":"6 seconds","uptime_seconds":6,"version":"0.57.1"}

	nc.Subscribe("wasmbus.ctl.default.auction.actor", func(m *nats.Msg) {
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
				HostId:      cli.HostId,
			}

			aRb, _ := json.Marshal(auctionResponse)
			nc.Publish(m.Reply, aRb)

			splitOCI := strings.Split(auction.ActorRef, ":")
			actor, err := oci.PullOCIRef(ctx, splitOCI[0], splitOCI[1])
			if err != nil {
				log.Panicln(err)
			}

			mod, err = r.InstantiateModuleFromBinary(ctx, actor)
			if err != nil {
				log.Panicln(err)
			}

			nc.Subscribe("wasmbus.rpc.default.MBCFOPM6JW2APJLXJD3Z5O4CN7CPYJ2B4FTKLJUR5YR5MITIU7HD3WD5", func(m *nats.Msg) {
				d := msgpack.NewDecoder(m.Data)
				i, _ := core.MDecodeInvocation(&d)
				tempData = m.Data
				tempInv = i

				_, err := mod.ExportedFunction("__guest_call").Call(ctx, uint64(len([]byte(i.Operation))), uint64(len(m.Data)))
				if err != nil {
					log.Panicln(err)
				}

				aerr, _ := mod.Memory().Read(wasmbus.Err, wasmbus.ErrLen)

				ir := core.InvocationResponse{
					Msg:           resp,
					InvocationId:  i.Id,
					Error:         string(aerr),
					ContentLength: uint64(len(resp)),
				}

				var sizer msgpack.Sizer
				size_enc := &sizer
				ir.MEncode(size_enc)
				buf := make([]byte, sizer.Len())
				encoder := msgpack.NewEncoder(buf)
				enc := &encoder
				ir.MEncode(enc)

				nc.Publish(m.Reply, buf)
			})

		}
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	r.Close(ctx)
	invSub.Drain()
	nc.Close()
	os.Exit(1)
}

func GuestRequest(operationPtr uint32, payloadPtr uint32) {
	mod.Memory().Write(payloadPtr, tempData)
	mod.Memory().Write(operationPtr, []byte(tempInv.Operation))
	log.Printf("op ptr: %d / payload ptr: %d", operationPtr, payloadPtr)
	log.Print("__guest_request called")
}

func GuestResponse(ptr, len uint32) {
	resp, _ = mod.Memory().Read(ptr, len)
	log.Print("__guest_response called")
}

type WasmcloudHost struct {
	HostId    string            `json:"host_id"`
	Issuer    string            `json:"issuer"` // cluster key
	Labels    map[string]string `json:"labels"`
	Friendly  string            `json:"friendly_name"`
	Actors    []string          `json:"actors"`
	Providers []string          `json:"providers"`
}
