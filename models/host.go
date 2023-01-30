package models

type Host struct {
	ClusterIssuers string      `json:"cluster_issuers,omitempty"`
	CtlHost        string      `json:"ctl_host,omitempty"`
	Id             string      `json:"id"`
	JsDomain       string      `json:"js_domain,omitempty"`
	Labels         KeyValueMap `json:"labels,omitempty"`
	LatticePrefix  string      `json:"lattice_prefix"`
	ProvRpcHost    string      `json:"prov_rpc_host,omitempty"`
	RpcHost        string      `json:"rpc_host,omitempty"`
	UptimeHuman    string      `json:"uptime_human"`
	UptimeSeconds  uint64      `json:"uptime_seconds"`
	Version        string      `json:"version,omitempty"`
}

type Hosts []Host

type HostInventory struct {
	Actors    []ActorDescription   `json:"actors"`
	HostId    string               `json:"host_id"`
	Labels    LabelsMap            `json:"labels"`
	Providers ProviderDescriptions `json:"providers"`
}
