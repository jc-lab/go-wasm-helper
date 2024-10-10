package wret

import (
	"github.com/jc-lab/go-wasm-helper/wmem"
)

type RefId uint64

const (
	RetIsObject   RefId = 0x80000000
	RetIsError    RefId = 0x40000000
	RetIsBytes    RefId = 0x20000000
	RetLengthMask RefId = 0x0fffffff
)

func (r RefId) GetPointer() wmem.WasmPtr {
	return wmem.WasmPtr(r >> 32)
}

func (r RefId) IsObject() bool {
	return (r & RetIsObject) != 0
}

func (r RefId) IsError() bool {
	return (r & RetIsError) != 0
}

func (r RefId) IsBytes() bool {
	return (r & RetIsBytes) != 0
}

func (r RefId) GetLength() int {
	return int(r & RetLengthMask)
}

func newRefId(ptr wmem.WasmPtr, flags RefId, length int) RefId {
	return RefId(uint64(ptr)<<uint64(32)) | flags | RefId(uint64(length)&uint64(RetLengthMask))
}

func ReturnObject(data interface{}) RefId {
	ptr := wmem.KeepaliveObject(data)
	return newRefId(ptr, RetIsObject, -1)
}

func ReturnBuffer(data []byte) RefId {
	ptr := wmem.KeepaliveObject(data)
	return newRefId(ptr, RetIsBytes, len(data))
}

func ReturnBufferWithFlag(data []byte, flags RefId) RefId {
	ptr := wmem.KeepaliveObject(data)
	return newRefId(ptr, flags|RetIsBytes, len(data))
}

func ReturnVoid() RefId {
	return RefId(0)
}

func ReturnError(err error) RefId {
	wrapped := &Error{
		Message: err.Error(),
	}
	encoded, err := wrapped.MarshalMsg(nil)
	if err != nil {
		panic(err)
	}
	return ReturnBufferWithFlag(encoded, RetIsError)
}
