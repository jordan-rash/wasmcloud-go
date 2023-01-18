package models

type ActorDescription struct {
	Name    string `json:"name,omitempty"`
	Hash    string `json:"hash"`
	Version string `json:"version"`

	Id        string         `json:"id"`
	ImageRef  string         `json:"image_ref,omitempty"`
	Instances ActorInstances `json:"instances"`
}

type ActorDescriptions []ActorDescription

type ActorInstance struct {
	Annotations AnnotationMap `json:"annotations,omitempty"`
	InstanceID  string        `json:"instance_id"`
	Revision    int32         `json:"revision"`
}
type ActorInstances []ActorInstance

type AnnotationMap map[string]string
