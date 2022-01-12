package client

type claim struct {
	CallAlias string `json:"call_alias"`
	Caps      string `json:"caps"`
	Iss       string `json:"iss"`
	Name      string `json:"name"`
	Rev       string `json:"rev"`
	Sub       string `json:"sub"`
	Tags      string `json:"tags"`
	Version   string `json:"version"`
}

type claims struct {
	Claims []claim `json:"claims"`
}
