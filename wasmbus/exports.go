package wasmbus

import (
	"log"
)

func (Wasmbus) hostCall(
	bindingPtr uintptr, bindingLen uint32,
	namespacePtr uintptr, namespaceLen uint32,
	operationPtr uintptr, operationLen uint32,
	payloadPtr uintptr, payloadLen uint32) uint32 {
	log.Print("__host_call called")
	return 0
}
func (Wasmbus) hostResponse(ptr uint32) { log.Print("__host_response called") }
func (Wasmbus) hostResponseLen() uint32 { log.Print("__host_response_len called"); return 0 }
func (Wasmbus) hostError(ptr uint32)    { log.Print("__host_error called") }
func (Wasmbus) hostErrorLen() uint32    { log.Print("__host_error_len called"); return 0 }

func (a *Wasmbus) guestRequest(operationPtr uint32, payloadPtr uint32) {
	log.Print("__guest_request called")

	a.module.Memory().Write(operationPtr, []byte(a.Operation))
	a.module.Memory().Write(payloadPtr, a.Data)
	log.Printf("op ptr: %d / payload ptr: %d", operationPtr, payloadPtr)
}

func (a *Wasmbus) guestResponse(ptr, len uint32) {
	log.Print("__guest_response called")

	a.gr, _ = a.module.Memory().Read(ptr, len)
}

func (a *Wasmbus) guestError(ptr uint32, len uint32) {
	log.Print("__guest_error called")

	a.err = ptr
	a.errlen = len
	log.Printf("ptr: %d / len: %d", ptr, len)
}

func (a Wasmbus) consoleLog(ptr, len uint32) {
	log.Print("__console_log called")

	logLine, _ := a.module.Memory().Read(ptr, len)
	log.Printf("%s", logLine)
}
