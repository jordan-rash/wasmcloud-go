package cli

import (
	"context"
	"time"
)

var VERSION string = "v0.0.0"

type WasmcloudHost struct {
	Context context.Context `json:"-" kong:"-"`

	// Host Settings
	WasmcloudClusterSeed        string   `json:"-" kong:"name='cluster_seed',group='Host Settings',env=WASMCLOUD_CLUSTER_SEED"`
	WasmcloudClusterIssuers     []string `json:"-" kong:"name='cluster_issuers',group='Host Settings',env=WASMCLOUD_CLUSTER_ISSUERS"`
	WasmcloudStructuredLogLevel string   `json:"-" kong:"name='structured_log_level',group='Host Settings',env=WASMCLOUD_STRUCTURED_LOG_LEVEL"`
	WasmcloudLatticePrefix      string   `json:"-" kong:"name='lattice_prefix',group='Host Settings',env=WASMCLOUD_LATTICE_PREFIX,default='default'"`

	// NATs Configuration
	WasmcloudNatsRemoteUrl string `json:"-" kong:"name='nats_remote_url',group='NATS Settings',env=WASMCLOUD_NATS_REMOTE_URL,default='127.0.0.1:7422'"`
	WasmcloudNatsJsDomain  string `json:"-" kong:"name='nats_js_domain',group='NATS Settings',env=WASMCLOUD_NATS_JS_DOMAIN,default='core'"`

	// RPC Settings
	WasmcloudRpcHost      string `json:"-" kong:"name='rpc_host',group='RPC Configuration',env=WASMCLOUD_RPC_HOST"`
	WasmcloudRpcPort      string `json:"-" kong:"name='rpc_port',group='RPC Configuration',env=WASMCLOUD_RPC_PORT"`
	WasmcloudRpcSeed      string `json:"-" kong:"name='rpc_seed',group='RPC Configuration',env=WASMCLOUD_RPC_SEED"`
	WasmcloudRpcJwt       string `json:"-" kong:"name='rpc_jwt',group='RPC Configuration',env=WASMCLOUD_RPC_JWT"`
	WasmcloudRpcTimeoutMs string `json:"-" kong:"name='rpc_timeout',group='RPC Configuration',env=WASMCLOUD_RPC_TIMEOUT_MS"`

	// Provider Settings
	WasmcloudProvRpcHost         string `json:"-" kong:"name='prov_rpc_host',group='Provider Configuration',env=WASMCLOUD_PROV_RPC_HOST"`
	WasmcloudProvRpcPort         string `json:"-" kong:"name='prov_rpc_port',group='Provider Configuration',env=WASMCLOUD_PROV_RPC_PORT"`
	WasmcloudProvRpcSeed         string `json:"-" kong:"name='proc_rpc_seed',group='Provider Configuration',env=WASMCLOUD_PROV_RPC_SEED"`
	WasmcloudProvRpcJwt          string `json:"-" kong:"name='prov_rpc_jwt',group='Provider Configuration',env=WASMCLOUD_PROV_RPC_JWT"`
	WasmcloudProvRpcTimeoutMs    string `json:"-" kong:"name='prov_rpc_shutdown_delay',group='Provider Configuration',env=WASMCLOUD_PROV_RPC_TIMEOUT_MS"`
	WasmcloudProvShutdownDelayMs string `json:"-" kong:"name='prov_shutdown_delay',group='Provider Configuration',env=WASMCLOUD_PROV_SHUTDOWN_DELAY_MS"`

	// Control Interface Settings
	WasmcloudCtlHost string `json:"-" kong:"name='ctl_host',group='Control Configuration',env=WASMCLOUD_CTL_HOST"`
	WasmcloudCtlPort string `json:"-" kong:"name='ctl_port',group='Control Configuration',env=WASMCLOUD_CTL_PORT"`
	WasmcloudCtlSeed string `json:"-" kong:"name='ctl_seed',group='Control Configuration',env=WASMCLOUD_CTL_SEED"`
	WasmcloudCtlJwt  string `json:"-" kong:"name='ctl_jwt',group='Control Configuration',env=WASMCLOUD_CTL_JWT"`

	// OCI Settings
	WasmcloudOciAllowLatest      bool   `json:"-" kong:"name='oci_allow_latest',group='OCI Configuration',env=WASMCLOUD_OCI_ALLOW_LATEST"`
	WasmcloudOciAllowInsecure    bool   `json:"-" kong:"name='oci_allow_insecure',group='OCI Configuration',env=WASMCLOUD_OCI_ALLOW_INSECURE"`
	WasmcloudOciRegistry         string `json:"-" kong:"name='oci_registry',group='OCI Configuration',env=WASMCLOUD_OCI_REGISTRY"`
	WasmcloudOciRegistryUser     string `json:"-" kong:"name='oci_registry_user',group='OCI Configuration',env=WASMCLOUD_OCI_REGISTRY_USER"`
	WasmcloudOciRegistryPassword string `json:"-" kong:"name='oci_registry_password',group='OCI Configuration',env=WASMCLOUD_OCI_REGISTRY_PASSWORD"`

	// Non-cli vars
	HostId    string            `json:"host_id" kong:"-"`
	Issuer    string            `json:"issuer" kong:"-"`
	Version   string            `json:"version" kong:"-"`
	Labels    map[string]string `json:"labels" kong:"-"`
	Friendly  string            `json:"friendly_name" kong:"-"`
	Actors    []string          `json:"actors" kong:"-"`
	Providers []string          `json:"providers" kong:"-"`
	Uptime    time.Time         `json:"uptime" kong:"-"`
}
