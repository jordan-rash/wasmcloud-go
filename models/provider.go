package models

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
