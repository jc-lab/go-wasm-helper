package whelper

import (
	"reflect"
)

var alivePointers = map[WasmPtr]interface{}{}

//export goPtrAllocate
func GoPtrAllocate(size uint32) WasmPtr {
	return KeepaliveObject(make([]byte, size))
}

//export goPtrFree
func GoPtrFree(ptr WasmPtr) {
	old := alivePointers[ptr]
	delete(alivePointers, ptr)
	tmp, ok := old.([]byte)
	if ok {
		for i := range tmp {
			tmp[i] = 0
		}
	}
}

func KeepaliveObject(obj interface{}) WasmPtr {
	value := reflect.ValueOf(obj)
	wasmPtr := WasmPtr(value.Pointer())
	alivePointers[wasmPtr] = obj
	return wasmPtr
}

func GetKeepaliveObject(ptr WasmPtr) (interface{}, bool) {
	obj, ok := alivePointers[ptr]
	return obj, ok
}
