package wmem

import (
	"fmt"
	"github.com/jc-lab/go-wasm-helper/wrefid"
	"reflect"
	"unsafe"
)

var alivePointers = map[wrefid.WasmPtr]interface{}{}

//export goPtrAllocate
func GoPtrAllocate(size uint32) wrefid.WasmPtr {
	return KeepaliveObject(make([]byte, size))
}

//export goPtrFree
func GoPtrFree(ptr wrefid.WasmPtr) {
	old := alivePointers[ptr]
	delete(alivePointers, ptr)
	tmp, ok := old.([]byte)
	if ok {
		for i := range tmp {
			tmp[i] = 0
		}
	}
}

func KeepaliveObject(obj interface{}) wrefid.WasmPtr {
	value := reflect.ValueOf(obj)
	wasmPtr := wrefid.WasmPtr(value.Pointer())
	alivePointers[wasmPtr] = obj
	return wasmPtr
}

func ParamBuffer(info wrefid.RefId) []byte {
	if !info.IsBytes() {
		panic(fmt.Errorf("refid 0x%x is not buffer", info))
	}
	ptr := uintptr(info.GetPointer())
	length := info.GetLength()
	if ptr == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(ptr)), length)
}
