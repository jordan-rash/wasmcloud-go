//go:build embedded

package cli

import (
	"context"
	"time"

	"github.com/jordan-rash/nkeys"
)

func (c *WasmcloudHost) Parse() {
	cIssuer, err := nkeys.CreateCluster()
	if err != nil {
		panic(err)
	}

	c.Issuer, err = cIssuer.PublicKey()
	if err != nil {
		panic(err)
	}

	seed, err := cIssuer.Seed()
	if err != nil {
		panic(err)
	}
	cSeed, err := nkeys.CreateServer()
	if err != nil {
		panic(err)
	}
	pubClusterSeed, err := cSeed.PublicKey()
	if err != nil {
		panic(err)
	}

	c.HostId = string(pubClusterSeed)

	c.Context = context.Background()
	c.WasmcloudClusterSeed = string(seed)
	c.WasmcloudClusterIssuers = []string{}
	c.WasmcloudStructuredLogLevel = "error"
	c.WasmcloudLatticePrefix = "default"
	c.WasmcloudNatsRemoteUrl = "127.0.0.1:7422"
	c.WasmcloudNatsJsDomain = "core"

	c.Version = VERSION
	c.Friendly = "yolo-bro-1234"
	c.Actors = []string{}
	c.Providers = []string{}
	c.Uptime = time.Now()

}
