//go:build !embedded

package cli

import (
	"runtime"
	"time"

	"github.com/alecthomas/kong"
	"github.com/jordan-rash/nkeys"
	"github.com/jordan-rash/wasmcloud-go/models"
)

func (c *WasmcloudHost) Parse() {
	_ = kong.Parse(c,
		kong.Name("wasmcloud-go"),
		kong.Description("wasmcloud host implementation written in Go using Wazero"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: false,
		}),
	)
}

func (c *WasmcloudHost) Validate() error {
	cSeed, err := nkeys.CreateServer()
	if err != nil {
		return err
	}
	pubClusterSeed, err := cSeed.PublicKey()
	if err != nil {
		return err
	}

	c.HostId = string(pubClusterSeed)

	if c.WasmcloudClusterSeed == "" {
		cIssuer, err := nkeys.CreateCluster()
		if err != nil {
			return err
		}

		pubClusterIssuer, err := cIssuer.PublicKey()
		if err != nil {
			return err
		}

		seedClusterIssuer, err := cIssuer.Seed()
		if err != nil {
			return err
		}
		c.Issuer = string(pubClusterIssuer)
		c.WasmcloudClusterSeed = string(seedClusterIssuer)
	}

	c.Labels = map[string]string{"hostcore.arch": runtime.GOARCH, "hostcore.os": runtime.GOOS, "hostcore.library": "wasmcloud_go", "hostcore.version": VERSION, "hostcore.runtime": "wazero"}

	c.Version = VERSION
	c.Friendly = "yolo-bro-1234"
	c.Actors = models.ActorDescriptions{}
	c.Providers = []string{}
	c.Uptime = time.Now()

	return nil
}
