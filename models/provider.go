package models

type ProviderAcutionAck struct {
	HostId      string `json:"host_id"`
	LinkName    string `json:"link_name"`
	ProviderRef string `json:"provider_ref"`
}

type ProviderAuctionRequest struct {
	Constraints ConstraintMap `json:"constraints"`
	LinkName    string        `json:"link_name"`
	ProviderRef string        `json:"provider_ref"`
}

type ProviderDescription struct {
	Annotations AnnotationMap `json:"annotations,omitempty"`
	Id          string        `json:"id"`
	ImageRef    string        `json:"image_ref,omitempty"`
	ContractId  string        `json:"contract_id"`
	LinkName    string        `json:"link_name"`
	Name        string        `json:"name,omitempty"`
	Revision    int32         `json:"revision"`
}

type ProviderDescriptions []ProviderDescription
