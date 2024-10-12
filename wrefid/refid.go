package wrefid

type WasmPtr uint32
type RefId uint64

const (
	RetIsObject   RefId = 0x80000000
	RetIsError    RefId = 0x40000000
	RetIsBytes    RefId = 0x20000000
	RetLengthMask RefId = 0x0fffffff
)

func (r RefId) GetPointer() WasmPtr {
	return WasmPtr(r >> 32)
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

func NewRefId(ptr WasmPtr, flags RefId, length int) RefId {
	return RefId(uint64(ptr)<<uint64(32)) | flags | RefId(uint64(length)&uint64(RetLengthMask))
}
