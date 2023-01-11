package cli

import (
	"context"
	"log"

	"github.com/nats-io/nkeys"
)

type Context struct {
	Context context.Context `kong:"-"`

	// Cluster Settings
	WasmcloudClusterSeed        string   `kong:"name='cluster_seed',env=WASMCLOUD_CLUSTER_SEED"`
	WasmcloudClusterIssuers     []string `kong:"name='cluster_issuers',env=WASMCLOUD_CLUSTER_ISSUERS"`
	WasmcloudJsDomain           string   `kong:"name='js_domain',env=WASMCLOUD_JS_DOMAIN"`
	WasmcloudStructuredLogLevel string   `kong:"name='structured_log_level',env=WASMCLOUD_STRUCTURED_LOG_LEVEL"`
	WasmcloudLatticePrefix      string   `kong:"name='lattice_prefix',env=WASMCLOUD_LATTICE_PREFIX"`
	// RPC Settings
	WasmcloudRpcHost      string `kong:"name='rpc_host',env=WASMCLOUD_RPC_HOST"`
	WasmcloudRpcPort      string `kong:"name='rpc_port',env=WASMCLOUD_RPC_PORT"`
	WasmcloudRpcSeed      string `kong:"name='rpc_seed',env=WASMCLOUD_RPC_SEED"`
	WasmcloudRpcJwt       string `kong:"name='rpc_jwt',env=WASMCLOUD_RPC_JWT"`
	WasmcloudRpcTimeoutMs string `kong:"name='rpc_timeout',env=WASMCLOUD_RPC_TIMEOUT_MS"`
	// Provider Settings
	WasmcloudProvRpcHost         string `kong:"name='prov_rpc_host',env=WASMCLOUD_PROV_RPC_HOST"`
	WasmcloudProvRpcPort         string `kong:"name='prov_rpc_port',env=WASMCLOUD_PROV_RPC_PORT"`
	WasmcloudProvRpcSeed         string `kong:"name='proc_rpc_seed',env=WASMCLOUD_PROV_RPC_SEED"`
	WasmcloudProvRpcJwt          string `kong:"name='prov_rpc_jwt',env=WASMCLOUD_PROV_RPC_JWT"`
	WasmcloudProvRpcTimeoutMs    string `kong:"name='prov_rpc_shutdown_delay',env=WASMCLOUD_PROV_RPC_TIMEOUT_MS"`
	WasmcloudProvShutdownDelayMs string `kong:"name='prov_shutdown_delay',env=WASMCLOUD_PROV_SHUTDOWN_DELAY_MS"`
	// Control Interface Settings
	WasmcloudCtlHost string `kong:"name='ctl_host',env=WASMCLOUD_CTL_HOST"`
	WasmcloudCtlPort string `kong:"name='ctl_port',env=WASMCLOUD_CTL_PORT"`
	WasmcloudCtlSeed string `kong:"name='ctl_seed',env=WASMCLOUD_CTL_SEED"`
	WasmcloudCtlJwt  string `kong:"name='ctl_jwt',env=WASMCLOUD_CTL_JWT"`
	// OCI Settings
	WasmcloudOciAllowLatest      bool   `kong:"name='oci_allow_latest',env=WASMCLOUD_OCI_ALLOW_LATEST"`
	WasmcloudOciAllowInsecure    bool   `kong:"name='oci_allow_insecure',env=WASMCLOUD_OCI_ALLOW_INSECURE"`
	WasmcloudOciRegistry         string `kong:"name='oci_registry',env=WASMCLOUD_OCI_REGISTRY"`
	WasmcloudOciRegistryUser     string `kong:"name='oci_registry_user',env=WASMCLOUD_OCI_REGISTRY_USER"`
	WasmcloudOciRegistryPassword string `kong:"name='oci_registry_password',env=WASMCLOUD_OCI_REGISTRY_PASSWORD"`

	HostId string `kong:"-"`
}

func (c *Context) Validate() error {
	if c.WasmcloudClusterSeed == "" || !nkeys.IsValidPublicServerKey(c.WasmcloudClusterSeed) {
		cSeed, err := nkeys.CreateServer()
		if err != nil {
			return err
		}
		pubClusterSeed, err := cSeed.PublicKey()
		if err != nil {
			return err
		}
		c.HostId = string(pubClusterSeed)
		log.Printf("host id: %s", pubClusterSeed)
	}

	return nil
}
