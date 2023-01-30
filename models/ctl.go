package models

type CtlOperationAck struct {
	Accepted bool   `json:"accepted"`
	Error    string `json:"error"`
}
