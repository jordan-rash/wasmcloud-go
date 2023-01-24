package models

type StopHostCommand struct {
	HostId  string `json:"host_id"`
	Timeout uint64 `json:"timeout,omitempty"`
}

type ScaleActorCommand struct {
	ActorId     string        `json:"actor_id"`
	ActorRef    string        `json:"actor_ref"`
	Annotations AnnotationMap `json:"annotations,omitempty"`
	Count       uint16        `json:"count"`
	HostId      string        `json:"host_id"`
}

type StartActorCommand struct {
	ActorRef    string        `json:"actor_ref"`
	Annotations AnnotationMap `json:"annotations,omitempty"`
	Count       uint16        `json:"count"`
	HostId      string        `json:"host_id"`
}

type StopActorCommand struct {
	ActorRef    string        `json:"actor_ref"`
	Annotations AnnotationMap `json:"annotations,omitempty"`
	Count       uint16        `json:"count"`
	HostId      string        `json:"host_id"`
}

type UpdateActorCommand struct {
	ActorRef    string        `json:"actor_ref"`
	Annotations AnnotationMap `json:"annotations,omitempty"`
	HostId      string        `json:"host_id"`
	NewActorRef string        `json:"new_actor_ref"`
}

type StartProviderCommand struct {
	Annotations   AnnotationMap `json:"annotations,omitempty"`
	Configuration string        `json:"configuration,omitempty"`
	HostId        string        `json:"host_id"`
	LinkName      string        `json:"link_name"`
	ProviderRef   string        `json:"provider_ref"`
}

type StopProviderCommand struct {
	Annotations AnnotationMap `json:"annotations,omitempty"`
	ContractId  string        `json:"contract_id"`
	HostId      string        `json:"host_id"`
	LinkName    string        `json:"link_name"`
	ProviderRef string        `json:"provider_ref"`
}
