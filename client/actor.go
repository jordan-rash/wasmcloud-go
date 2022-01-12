package client

type instance struct {
	Annotations map[string]interface{} `json:"annotations"`
	InstanceID  string                 `json:"instance_id"`
	Revision    int                    `json:"revision"`
}

type actor struct {
	ID        string     `json:"id"`
	ImageRef  string     `json:"image_ref"`
	Name      string     `json:"name"`
	Instances []instance `json:"instances"`
}
