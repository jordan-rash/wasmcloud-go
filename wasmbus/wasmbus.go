package wasmbus

import (
	"context"
	"log"

	"github.com/jordan-rash/wasmcloud-go/internal/cli"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type Wasmbus struct {
	runtime wazero.Runtime
	module  api.Module

	Context    context.Context
	ActorBytes []byte
	Data       []byte
	Operation  string

	gr     []byte
	err    uint32
	errlen uint32
}

// TODO: Namespaces not implemented because it look like team is
// going a different direction with cache.
func NewWasmbus(host cli.WasmcloudHost) (*Wasmbus, error) {
	wb := new(Wasmbus)

	r := wazero.NewRuntime(host.Context)
	wasi_snapshot_preview1.MustInstantiate(host.Context, r)

	_, err := r.NewHostModuleBuilder("wasmbus").
		NewFunctionBuilder().WithFunc(wb.guestRequest).Export("__guest_request").
		NewFunctionBuilder().WithFunc(wb.guestResponse).Export("__guest_response").
		NewFunctionBuilder().WithFunc(wb.guestError).Export("__guest_error").
		NewFunctionBuilder().WithFunc(wb.consoleLog).Export("__console_log").
		NewFunctionBuilder().WithFunc(wb.hostCall).Export("__host_call").
		NewFunctionBuilder().WithFunc(wb.hostResponse).Export("__host_response").
		NewFunctionBuilder().WithFunc(wb.hostResponseLen).Export("__host_response_len").
		NewFunctionBuilder().WithFunc(wb.hostError).Export("__host_error").
		NewFunctionBuilder().WithFunc(wb.hostErrorLen).Export("__host_error_len").
		Instantiate(host.Context, r)
	if err != nil {
		return nil, err
	}

	wb.runtime = r
	wb.Context = host.Context
	return wb, nil
}

func (wb Wasmbus) GetGuestError() (uint32, uint32) {
	return wb.err, wb.errlen
}
func (wb Wasmbus) GetGuestResponse() []byte {
	return wb.gr
}

func (wb *Wasmbus) CreateModule(actorBytes []byte) (api.Module, error) {
	var err error

	log.Printf("actors bytes: %d", len(actorBytes))
	mod, err := wb.runtime.InstantiateModuleFromBinary(wb.Context, actorBytes)
	if err != nil {
		return nil, err
	}

	wb.module = mod
	return mod, nil
}
