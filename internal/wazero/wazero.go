package wazero

import (
	"context"
	"log"

	"github.com/jordan-rash/wasmcloud-go/wasmbus"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type Actor struct {
	Context    context.Context
	ActorBytes []byte
	Data       []byte
	Operation  string

	gr     []byte
	err    uint32
	errlen uint32
	mod    api.Module
}

func (a Actor) GetGuestResponse() []byte {
	return a.gr
}

func (a Actor) GetGuestError() (uint32, uint32) {
	return a.err, a.errlen
}

func (a Actor) GetModule() api.Module {
	return a.mod
}

func (a *Actor) StartActor() error {
	r := wazero.NewRuntime(a.Context)
	wasi_snapshot_preview1.MustInstantiate(a.Context, r)
	wasmbus := wasmbus.Wasmbus{}

	_, err := r.NewHostModuleBuilder("wasmbus").
		NewFunctionBuilder().WithFunc(a.guestRequest).Export("__guest_request").
		NewFunctionBuilder().WithFunc(a.guestResponse).Export("__guest_response").
		NewFunctionBuilder().WithFunc(a.guestError).Export("__guest_error").
		NewFunctionBuilder().WithFunc(a.consoleLog).Export("__console_log").
		NewFunctionBuilder().WithFunc(wasmbus.HostCall).Export("__host_call").
		NewFunctionBuilder().WithFunc(wasmbus.HostResponse).Export("__host_response").
		NewFunctionBuilder().WithFunc(wasmbus.HostResponseLen).Export("__host_response_len").
		NewFunctionBuilder().WithFunc(wasmbus.HostError).Export("__host_error").
		NewFunctionBuilder().WithFunc(wasmbus.HostErrorLen).Export("__host_error_len").
		Instantiate(a.Context, r)
	if err != nil {
		return err
	}

	a.mod, err = r.InstantiateModuleFromBinary(a.Context, a.ActorBytes)
	if err != nil {
		return err
	}

	return nil
}

func (a *Actor) guestRequest(operationPtr uint32, payloadPtr uint32) {
	log.Print("__guest_request called")

	a.mod.Memory().Write(operationPtr, []byte(a.Operation))
	a.mod.Memory().Write(payloadPtr, a.Data)
	log.Printf("op ptr: %d / payload ptr: %d", operationPtr, payloadPtr)
}

func (a *Actor) guestResponse(ptr, len uint32) {
	log.Print("__guest_response called")

	a.gr, _ = a.mod.Memory().Read(ptr, len)
}

func (a *Actor) guestError(ptr uint32, len uint32) {
	log.Print("__guest_error called")

	a.err = ptr
	a.errlen = len
	log.Printf("ptr: %d / len: %d", ptr, len)
}

func (a Actor) consoleLog(ptr, len uint32) {
	log.Print("__console_log called")

	logLine, _ := a.mod.Memory().Read(ptr, len)
	log.Printf("%s", logLine)
}
