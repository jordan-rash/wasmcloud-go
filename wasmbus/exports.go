package wasmbus

import (
	"log"
)

type Wasmbus struct {
	Err    uint32
	ErrLen uint32
}

func (Wasmbus) HostCall(
	bindingPtr uintptr, bindingLen uint32,
	namespacePtr uintptr, namespaceLen uint32,
	operationPtr uintptr, operationLen uint32,
	payloadPtr uintptr, payloadLen uint32) uint32 {
	log.Print("__host_call called")
	return 0
}
func (Wasmbus) ConsoleLog(ptr uint32, sz uint32) { log.Print("__console_log called") }
func (w *Wasmbus) GuestRequest(operationPtr uint32, payloadPtr uint32) {
	log.Printf("op ptr: %d / payload ptr: %d", operationPtr, payloadPtr)
	log.Print("__guest_request called")
}
func (Wasmbus) HostResponse(ptr uint32)              { log.Print("__host_response called") }
func (Wasmbus) HostResponseLen() uint32              { log.Print("__host_response_len called"); return 0 }
func (Wasmbus) GuestResponse(ptr uint32, len uint32) { log.Print("__guest_response called") }
func (w *Wasmbus) GuestError(ptr uint32, len uint32) {
	w.Err = ptr
	w.ErrLen = len
	log.Printf("ptr: %d / len: %d", ptr, len)
	log.Print("__guest_error called")
}
func (Wasmbus) HostError(ptr uint32) { log.Print("__host_error called") }
func (Wasmbus) HostErrorLen() uint32 { log.Print("__host_error_len called"); return 0 }
