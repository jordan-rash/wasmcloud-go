package models

import core "github.com/wasmcloud/interfaces/core/tinygo"

type KeyValueMap map[string]string
type CtlKVList []KeyValueMap

type GetClaimsResponse struct {
	Claims CtlKVList `json:"claims"`
}

type LinkDefinitionList struct {
	Links core.ActorLinks `json:"links"`
}
