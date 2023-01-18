package models

type HostInventory struct {
	HostID    string               `json:"host_id"`
	Labels    LabelsMap            `json:"labels"`
	Actors    ActorDescriptions    `json:"actors"`
	Providers ProviderDescriptions `json:"providers"`
}

type LabelsMap map[string]string
