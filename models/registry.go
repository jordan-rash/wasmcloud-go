package models

type RegistryCredential struct {
	Password     string `json:"password,omitempty"`
	Token        string `json:"token,omitempty"`
	Username     string `json:"username,omitempty"`
	RegistryType string `json:"registryType"`
}

type RegistrpCredentialMap map[string]RegistryCredential
