package client

type provider struct {
	ID         string `json:"id"`
	ImageRef   string `json:"image_ref"`
	InstanceID string `json:"instance_id"`
	LinkName   string `json:"link_name"`
	Name       string `json:"name"`
	Revision   int    `json:"revision"`
}
