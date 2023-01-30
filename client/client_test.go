package client_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/jordan-rash/nkeys"
	"github.com/jordan-rash/wasmcloud-go/client"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

var nc *nats.Conn

func init() {
	opt := server.Options{
		JetStream:       true,
		Host:            nats.DefaultURL,
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

	nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
}

func TestCreateDefaultClient(t *testing.T) {
	client, err := client.New(nc).Build()
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestGetHosts(t *testing.T) {
	client, err := client.New(nc,
		client.WithJsDomain("tester"),
	).Build()
	assert.NoError(t, err)
	assert.NotNil(t, client)

	_, err = client.GetHosts(2 * time.Second)
	assert.ErrorIs(t, err, nats.ErrNoResponders)

	s, _, err := startNewHost()
	assert.NoError(t, err)
	assert.NotNil(t, s)

	hosts, err := client.GetHosts(2 * time.Second)
	assert.NoError(t, err)
	assert.Len(t, *hosts, 1)

	h := *hosts
	assert.Equal(t, "tester", h[0].JsDomain)
}

func TestGetHostInventory(t *testing.T) {
	client, err := client.New(nc,
		client.WithJsDomain("tester"),
	).Build()
	assert.NoError(t, err)
	assert.NotNil(t, client)

	s, hostid, err := startNewHost()
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.NotEmpty(t, hostid)

	inv, err := client.GetHostInventory(hostid)
	assert.NoError(t, err)
	assert.Len(t, inv.Actors, 1)
	assert.Len(t, inv.Providers, 1)
	assert.Equal(t, hostid, inv.HostId)
	assert.Contains(t, inv.Labels, "test_key")
}

func genFriendly() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	var numbers = []rune("1234567890")
	return fmt.Sprintf("%U-%U-%U", letters[rand.Intn(len(letters))], letters[rand.Intn(len(letters))], numbers[rand.Intn(len(numbers))])
}

func startNewHost() ([]*nats.Subscription, string, error) {
	subs := []*nats.Subscription{}
	cSeed, _ := nkeys.CreateServer()
	pubClusterSeed, _ := cSeed.PublicKey()

	s, _ := nc.Subscribe("wasmbus.ctl.default.ping.hosts", func(m *nats.Msg) {
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
			ClusterIssuers: "",
			CtlHost:        "",
			Friendly:       genFriendly(),
			Id:             pubClusterSeed,
			JsDomain:       "tester",
			Labels:         map[string]string{},
			LatticePrefix:  "default",
			ProvRpcHost:    "127.0.0.1",
			RpcHost:        "127.0.0.1",
			UptimeHuman:    "1 second",
			Uptime:         1,
			Version:        "wasmcloud-tester",
		}
		sAck, _ := json.Marshal(ack)
		_ = nc.Publish(m.Reply, sAck)
	})
	subs = append(subs, s)

	s, _ = nc.Subscribe(fmt.Sprintf("wasmbus.ctl.default.get.%s.inv", pubClusterSeed), func(m *nats.Msg) {
		hostInventory := struct {
			HostID    string            `json:"host_id"`
			Labels    map[string]string `json:"labels"`
			Actors    interface{}       `json:"actors"`
			Providers interface{}       `json:"providers"`
		}{
			HostID: pubClusterSeed,
			Labels: map[string]string{"test_key": "test_value"},
			Actors: []struct {
				Name      string      `json:"name,omitempty"`
				Hash      string      `json:"hash"`
				Version   string      `json:"version"`
				Id        string      `json:"id"`
				ImageRef  string      `json:"image_ref,omitempty"`
				Instances interface{} `json:"instances"`
			}{
				{Name: "myactor", Hash: "123", Version: "1", Id: "1", ImageRef: "derp.com",
					Instances: []struct {
						Annotations map[string]string `json:"annotations,omitempty"`
						InstanceID  string            `json:"instance_id"`
						Revision    int32             `json:"revision"`
					}{
						{Annotations: map[string]string{}, InstanceID: "1", Revision: 1},
					},
				},
			},
			Providers: []struct {
				Annotations map[string]string `json:"annotations,omitempty"`
				Id          string            `json:"id"`
				ImageRef    string            `json:"image_ref,omitempty"`
				ContractId  string            `json:"contract_id"`
				LinkName    string            `json:"link_name"`
				Name        string            `json:"name,omitempty"`
				Revision    int32             `json:"revision"`
			}{
				{Id: "2", ImageRef: "derp.com", ContractId: "wasmcloud:test", LinkName: "default", Name: "oompa_lumpa", Annotations: map[string]string{}},
			},
		}

		sAck, _ := json.Marshal(hostInventory)
		_ = nc.Publish(m.Reply, sAck)
	})
	subs = append(subs, s)

	return subs, string(pubClusterSeed), nil
}
