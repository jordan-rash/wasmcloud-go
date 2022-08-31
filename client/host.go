package client

type hostStatus struct {
	FriendlyName string                 `json:"friendly_name"`
	ID           string                 `json:"host_id"`
	Labels       map[string]interface{} `json:"labels,omitempty"`
	Actors       []actor                `json:"actors,omitempty"`
	Providers    []provider             `json:"providers,omitempty"`
}

type host struct {
	FriendlyName string `json:"friendly_name"`
	ID           string `json:"id"`
	Uptime       int    `json:"uptime_seconds"`
}
