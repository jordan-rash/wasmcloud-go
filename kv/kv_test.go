package kv_test

import (
	"testing"

	"github.com/jordan-rash/wasmcloud-go/kv"
	"github.com/jordan-rash/wasmcloud-go/models"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	core "github.com/wasmcloud/interfaces/core/tinygo"
)

const (
	CLAIMS_1 = "{\"call_alias\":\"\",\"caps\":\"wasmcloud:httpserver\",\"iss\":\"ABRIBHH54GM7QIEJBYYGZJUSDAMO34YM4SKWUQJGIILRB7JYGXEPWUVT\",\"name\":\"kvcounter\",\"rev\":\"1631624220\",\"sub\":\"MBW3UGAIONCX3RIDDUGDCQIRGBQQOWS643CVICQ5EZ7SWNQPZLZTSQKU\",\"tags\":\"\",\"version\":\"0.3.0\"}"
	CLAIMS_2 = "{\"call_alias\":\"\",\"caps\":\"\",\"iss\":\"ACOJJN6WUP4ODD75XEBKKTCCUJJCY5ZKQ56XVKYK4BEJWGVAOOQHZMCW\",\"name\":\"HTTP Server\",\"rev\":\"1644594344\",\"sub\":\"VAG3QITQQ2ODAOWB5TTQSDJ53XK3SHBEIFNK4AYJ5RKAX2UNSCAPHA5M\",\"tags\":\"\",\"version\":\"0.14.10\"}"
	LINK_1   = "{\"actor_id\":\"MBW3UGAIONCX3RIDDUGDCQIRGBQQOWS643CVICQ5EZ7SWNQPZLZTSQKU\",\"contract_id\":\"wasmcloud:httpserver\",\"id\":\"fb30deff-bbe7-4a28-a525-e53ebd4e8228\",\"link_name\":\"default\",\"provider_id\":\"VAG3QITQQ2ODAOWB5TTQSDJ53XK3SHBEIFNK4AYJ5RKAX2UNSCAPHA5M\",\"values\":{\"PORT\":\"8082\"}}"
	LINK_2   = "{\"actor_id\":\"MBW3UGAIONCX3RIDDUGDCQIRGBQQOWS643CVICQ5EZ7SWNQPZLZTSQKU\",\"contract_id\":\"wasmcloud:keyvalue\",\"id\":\"ff140106-dd0d-44ee-8241-a2158a528b1d\",\"link_name\":\"default\",\"provider_id\":\"VAZVC4RX54J2NVCMCW7BPCAHGGG5XZXDBXFUMDUXGESTMQEJLC3YVZWB\",\"values\":{\"URL\":\"redis://127.0.0.1:6379\"}}"
)

var client *nats.Conn

func init() {
	opt := server.Options{
		JetStream:       true,
		Port:            4224,
		Host:            "127.0.0.1",
		JetStreamDomain: "testy",
	}
	n, err := server.NewServer(&opt)
	if err != nil {
		panic(err)
	}

	n.Start()

	client, err = nats.Connect("127.0.0.1:4224")
	if err != nil {
		panic(err)
	}
}

func TestHashCompatibility(t *testing.T) {
	const ELIXIR_HASH = "B40411AD09B70A2C83D59923584F66BA2C5A3C274DC4F19416DA49CCD6531F9C"

	ld := core.LinkDefinition{}
	ld.ActorId = "Mbob"
	ld.ContractId = "wasmcloud:testy"
	ld.LinkName = "default"
	ld.ProviderId = "Valice"

	h1, err := kv.LDHash(&ld)
	assert.NoError(t, err)
	assert.Equal(t, ELIXIR_HASH, h1)
}

func TestGetReturnsNonForNonexistentStore(t *testing.T) {
	store, err := kv.GetKVStore(client, "this-lattice-shall-never-existeth", "")

	assert.Nil(t, store)
	assert.Error(t, err)
	// TODO: figure out what error this is suppose to be
	// assert.ErrorIs(t, err, nats.Err?????)
}

func TestGetClaimsReturnsResponse(t *testing.T) {
	js, err := client.JetStream(nats.PublishAsyncMaxPending(256))
	assert.NoError(t, err)

	k, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: "LATTICEDATA_mylattice1",
	})
	assert.NoError(t, err)

	_, err = k.Put("CLAIMS_VAG3QITQQ2ODAOWB5TTQSDJ53XK3SHBEIFNK4AYJ5RKAX2UNSCAPHA5M", []byte(CLAIMS_2))
	assert.NoError(t, err)
	_, err = k.Put("CLAIMS_MBW3UGAIONCX3RIDDUGDCQIRGBQQOWS643CVICQ5EZ7SWNQPZLZTSQKU", []byte(CLAIMS_1))
	assert.NoError(t, err)

	store, err := kv.GetKVStore(client, "mylattice1", "")
	assert.NoError(t, err)

	claims, err := kv.GetClaims(store)
	assert.NoError(t, err)

	err = js.DeleteKeyValue("LATTICEDATA_mylattice1")
	assert.NoError(t, err)

	assert.Len(t, claims.Claims, 2)
	assert.Contains(t, claims.Claims[0], "name")
	assert.Contains(t, claims.Claims[0], "rev")
	assert.Contains(t, claims.Claims[0], "sub")
	assert.Contains(t, claims.Claims[1], "call_alias")
}

func TestGetLinksReturnsResponse(t *testing.T) {
	js, err := client.JetStream(nats.PublishAsyncMaxPending(256))
	assert.NoError(t, err)

	k, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: "LATTICEDATA_mylattice2",
	})
	assert.NoError(t, err)

	_, err = k.Put("LINKDEF_ff140106-dd0d-44ee-8241-a2158a528b1d", []byte(LINK_2))
	assert.NoError(t, err)
	_, err = k.Put("LINKDEF_fb30deff-bbe7-4a28-a525-e53ebd4e822", []byte(LINK_1))
	assert.NoError(t, err)

	store, err := kv.GetKVStore(client, "mylattice2", "")
	assert.NoError(t, err)

	links, err := kv.GetLinks(store)
	assert.NoError(t, err)

	err = js.DeleteKeyValue("LATTICEDATA_mylattice2")
	assert.NoError(t, err)

	assert.Len(t, links.Links, 2)
}

func TestPutAndDeleteLink(t *testing.T) {
	js, err := client.JetStream(nats.PublishAsyncMaxPending(256))
	assert.NoError(t, err)

	k, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: "LATTICEDATA_mylattice3",
	})
	assert.NoError(t, err)

	ld1 := core.LinkDefinition{
		ActorId:    "Mbob",
		ProviderId: "Valice",
		ContractId: "wasmcloud:testy",
		LinkName:   "default",
	}
	ld2 := core.LinkDefinition{
		ActorId:    "Msteve",
		ProviderId: "Valice",
		ContractId: "wasmcloud:testy",
		LinkName:   "default",
	}

	err = kv.PutLink(k, ld1)
	assert.NoError(t, err)
	err = kv.PutLink(k, ld2)
	assert.NoError(t, err)

	err = kv.DeleteLink(k, models.RemoveLinkDefinationRequest{"Mbob", "wasmcloud:testy", "default"})
	assert.NoError(t, err)

	links, err := kv.GetLinks(k)
	assert.NoError(t, err)

	err = js.DeleteKeyValue("LATTICEDATA_mylattice3")
	assert.NoError(t, err)

	assert.Len(t, links.Links, 1)
}
