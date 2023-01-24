package models

type RemoveLinkDefinationRequest struct {
	ActorId    string `json:"actor_id"`
	ContractId string `json:"contract_id"`
	LinkName   string `json:"link_name"`
}
