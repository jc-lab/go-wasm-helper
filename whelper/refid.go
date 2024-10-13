package whelper

import (
	"fmt"
	"unsafe"
)

type WasmPtr uint32
type RefId uint64

const (
	RefVoid RefId = 0

	RefIsNonPointer        RefId = 0x80000000
	RefIsObject            RefId = 0x40000000
	RefIsError             RefId = 0x20000000
	RefIsBytes             RefId = 0x10000000
	RefLengthOrSubTypeMask RefId = 0x0fffffff
)

type RefSubType uint32

const (
	RefSubTypeMask RefSubType = 0xfff0000

	RefSubTypeNumber         RefSubType = 0x0010000
	RefSubTypeNumberUnsigned RefSubType = RefSubTypeNumber | 0x1000
	RefSubTypeInt8           RefSubType = RefSubTypeNumber | 0x0001
	RefSubTypeUint8          RefSubType = RefSubTypeNumberUnsigned | 0x0001
	RefSubTypeInt16          RefSubType = RefSubTypeNumber | 0x0002
	RefSubTypeUint16         RefSubType = RefSubTypeNumberUnsigned | 0x0002
	RefSubTypeInt32          RefSubType = RefSubTypeNumber | 0x0004
	RefSubTypeUint32         RefSubType = RefSubTypeNumberUnsigned | 0x0004

	RefSubTypeFunction RefSubType = 0x0010001
)

func (r RefId) GetPointer() WasmPtr {
	return WasmPtr(r >> 32)
}

func (r RefId) IsVoid() bool {
	return r == RefVoid
}

func (r RefId) IsNonPointer() bool {
	return (r & RefIsNonPointer) != 0
}

func (r RefId) IsObject() bool {
	return (r & RefIsObject) != 0
}

func (r RefId) GetObject() interface{} {
	if !r.IsObject() {
		panic(fmt.Errorf("refid(0x%d) is not a object", r))
	}
	obj, _ := GetKeepaliveObject(r.GetPointer())
	return obj
}

func (r RefId) IsError() bool {
	return (r & RefIsError) != 0
}

func (r RefId) IsBytes() bool {
	return (r & RefIsBytes) != 0
}

func (r RefId) GetLength() uint32 {
	if !r.IsBytes() {
		return 0
	}
	return uint32(r & RefLengthOrSubTypeMask)
}

func (r RefId) GetBuffer() []byte {
	if !r.IsBytes() {
		panic(fmt.Errorf("refid(0x%d) is not a buffer", r))
	}
	ptr := uintptr(r.GetPointer())
	length := r.GetLength()
	if ptr == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(ptr)), length)
}

func (r RefId) GetSubType() RefSubType {
	if !r.IsNonPointer() {
		return 0
	}
	return RefSubType(r & RefLengthOrSubTypeMask)
}

func (r RefId) IsNumber() bool {
	return (r.GetSubType() & RefSubTypeMask) == RefSubTypeNumber
}

func (r RefId) GetNumber() uint32 {
	if !r.IsNumber() {
		panic(fmt.Errorf("refid(0x%d) is not a number", r))
	}
	return uint32(r.GetPointer())
}

func (r RefId) GetInt8() int8 {
	if !r.IsNumber() {
		panic(fmt.Errorf("refid(0x%d) is not a number", r))
	}
	return int8(r.GetPointer())
}

func (r RefId) GetUint8() uint8 {
	if !r.IsNumber() {
		panic(fmt.Errorf("refid(0x%d) is not a number", r))
	}
	return uint8(r.GetPointer())
}

func (r RefId) GetInt16() int16 {
	if !r.IsNumber() {
		panic(fmt.Errorf("refid(0x%d) is not a number", r))
	}
	return int16(r.GetPointer())
}

func (r RefId) GetUint16() uint16 {
	if !r.IsNumber() {
		panic(fmt.Errorf("refid(0x%d) is not a number", r))
	}
	return uint16(r.GetPointer())
}

func (r RefId) GetInt32() int32 {
	if !r.IsNumber() {
		panic(fmt.Errorf("refid(0x%d) is not a number", r))
	}
	return int32(r.GetPointer())
}

func (r RefId) GetUint32() uint32 {
	if !r.IsNumber() {
		panic(fmt.Errorf("refid(0x%d) is not a number", r))
	}
	return uint32(r.GetPointer())
}

func (r RefId) IsFunction() bool {
	return r.GetSubType() == RefSubTypeFunction
}

func (r RefId) ToFunction() CallbackFunc {
	if !r.IsFunction() {
		panic(fmt.Errorf("refid(0x%d) is not a function", r))
	}
	return func(args ...uint64) RefId {
		return goCallbackJsHandler(r, args)
	}
}

func NewRefId(ptr WasmPtr, flags RefId, lengthOrSubType uint32) RefId {
	return RefId(uint64(ptr)<<uint64(32)) | flags | RefId(uint64(lengthOrSubType)&uint64(RefLengthOrSubTypeMask))
}

func NewRefIdWithBytes(ptr WasmPtr, flags RefId, lengthOrSubType int) RefId {
	return NewRefId(ptr, flags, uint32(lengthOrSubType))
}

func Number(subType RefSubType, v uint32) RefId {
	return NewRefId(WasmPtr(v), RefIsNonPointer, uint32(subType))
}

func Int8(v int8) RefId {
	return Number(RefSubTypeInt8, uint32(v))
}

func Uint8(v uint8) RefId {
	return Number(RefSubTypeUint8, uint32(v))
}

func Int16(v int16) RefId {
	return Number(RefSubTypeInt16, uint32(v))
}

func Uint16(v uint16) RefId {
	return Number(RefSubTypeUint16, uint32(v))
}

func Int32(v int32) RefId {
	return Number(RefSubTypeInt32, uint32(v))
}

func Uint32(v uint32) RefId {
	return Number(RefSubTypeUint32, uint32(v))
}
