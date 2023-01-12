package cli

import (
	"context"
	"time"
)

var VERSION string = "v0.0.0"

type WasmcloudHost struct {
	Context context.Context `kong:"-"`

	// Host Settings
	WasmcloudClusterSeed        string   `kong:"name='cluster_seed',group='Host Settings',env=WASMCLOUD_CLUSTER_SEED"`
	WasmcloudClusterIssuers     []string `kong:"name='cluster_issuers',group='Host Settings',env=WASMCLOUD_CLUSTER_ISSUERS"`
	WasmcloudStructuredLogLevel string   `kong:"name='structured_log_level',group='Host Settings',env=WASMCLOUD_STRUCTURED_LOG_LEVEL"`
	WasmcloudLatticePrefix      string   `kong:"name='lattice_prefix',group='Host Settings',env=WASMCLOUD_LATTICE_PREFIX"`

	// NATs Configuration
	WasmcloudNatsRemoteUrl string `kong:"name='nats_remote_url',group='NATS Settings',env=WASMCLOUD_NATS_REMOTE_URL,default='127.0.0.1:7422'"`
	WasmcloudNatsJsDomain  string `kong:"name='nats_js_domain',group='NATS Settings',env=WASMCLOUD_NATS_JS_DOMAIN,default='core'"`

	// RPC Settings
	WasmcloudRpcHost      string `kong:"name='rpc_host',group='RPC Configuration',env=WASMCLOUD_RPC_HOST"`
	WasmcloudRpcPort      string `kong:"name='rpc_port',group='RPC Configuration',env=WASMCLOUD_RPC_PORT"`
	WasmcloudRpcSeed      string `kong:"name='rpc_seed',group='RPC Configuration',env=WASMCLOUD_RPC_SEED"`
	WasmcloudRpcJwt       string `kong:"name='rpc_jwt',group='RPC Configuration',env=WASMCLOUD_RPC_JWT"`
	WasmcloudRpcTimeoutMs string `kong:"name='rpc_timeout',group='RPC Configuration',env=WASMCLOUD_RPC_TIMEOUT_MS"`

	// Provider Settings
	WasmcloudProvRpcHost         string `kong:"name='prov_rpc_host',group='Provider Configuration',env=WASMCLOUD_PROV_RPC_HOST"`
	WasmcloudProvRpcPort         string `kong:"name='prov_rpc_port',group='Provider Configuration',env=WASMCLOUD_PROV_RPC_PORT"`
	WasmcloudProvRpcSeed         string `kong:"name='proc_rpc_seed',group='Provider Configuration',env=WASMCLOUD_PROV_RPC_SEED"`
	WasmcloudProvRpcJwt          string `kong:"name='prov_rpc_jwt',group='Provider Configuration',env=WASMCLOUD_PROV_RPC_JWT"`
	WasmcloudProvRpcTimeoutMs    string `kong:"name='prov_rpc_shutdown_delay',group='Provider Configuration',env=WASMCLOUD_PROV_RPC_TIMEOUT_MS"`
	WasmcloudProvShutdownDelayMs string `kong:"name='prov_shutdown_delay',group='Provider Configuration',env=WASMCLOUD_PROV_SHUTDOWN_DELAY_MS"`

	// Control Interface Settings
	WasmcloudCtlHost string `kong:"name='ctl_host',group='Control Configuration',env=WASMCLOUD_CTL_HOST"`
	WasmcloudCtlPort string `kong:"name='ctl_port',group='Control Configuration',env=WASMCLOUD_CTL_PORT"`
	WasmcloudCtlSeed string `kong:"name='ctl_seed',group='Control Configuration',env=WASMCLOUD_CTL_SEED"`
	WasmcloudCtlJwt  string `kong:"name='ctl_jwt',group='Control Configuration',env=WASMCLOUD_CTL_JWT"`

	// OCI Settings
	WasmcloudOciAllowLatest      bool   `kong:"name='oci_allow_latest',group='OCI Configuration',env=WASMCLOUD_OCI_ALLOW_LATEST"`
	WasmcloudOciAllowInsecure    bool   `kong:"name='oci_allow_insecure',group='OCI Configuration',env=WASMCLOUD_OCI_ALLOW_INSECURE"`
	WasmcloudOciRegistry         string `kong:"name='oci_registry',group='OCI Configuration',env=WASMCLOUD_OCI_REGISTRY"`
	WasmcloudOciRegistryUser     string `kong:"name='oci_registry_user',group='OCI Configuration',env=WASMCLOUD_OCI_REGISTRY_USER"`
	WasmcloudOciRegistryPassword string `kong:"name='oci_registry_password',group='OCI Configuration',env=WASMCLOUD_OCI_REGISTRY_PASSWORD"`

	// Non-cli vars
	HostId    string            `kong:"-"`
	Issuer    string            `kong:"-"`
	Version   string            `kong:"-"`
	Labels    map[string]string `kong:"-"`
	Friendly  string            `kong:"-"`
	Actors    []string          `kong:"-"`
	Providers []string          `kong:"-"`
	Uptime    time.Time         `kong:"-"`
}
