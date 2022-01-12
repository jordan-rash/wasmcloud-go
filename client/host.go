package client

type hostStatus struct {
	FriendlyName string                 `json:"friendly_name"`
	ID           string                 `json:"host_id"`
	Labels       map[string]interface{} `json:"labels"`
	Actors       []actor                `json:"actors"`
	Providers    []provider             `json:"providers"`
}

type host struct {
	FriendlyName string `json:"friendly_name"`
	ID           string `json:"id"`
	Uptime       int    `json:"uptime_seconds"`
}
