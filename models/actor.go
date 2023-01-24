package models

type ActorAucutionAck struct {
	ActorRef string `json:"actor_ref"`
	HostId   string `json:"host_id"`
}

type ActorAuctionRequest struct {
	ActorRef    string        `json:"actor_ref"`
	Constraints ConstraintMap `json:"constraints"`
}

type ActorDescription struct {
	Id        string         `json:"id"`
	ImageRef  string         `json:"image_ref,omitempty"`
	Instances ActorInstances `json:"instances"`
	Name      string         `json:"name"`
}

type ActorInstance struct {
	Annotations AnnotationMap `json:"annotations,omitempty"`
	InstanceId  string        `json:"intstance_id"`
	Revision    uint32        `json:"revision"`
}

type ActorInstances []ActorInstance
