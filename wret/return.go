package wret

import (
	"github.com/jc-lab/go-wasm-helper/whelper"
)

func ReturnObject(data interface{}) whelper.RefId {
	ptr := whelper.KeepaliveObject(data)
	return whelper.NewRefIdWithBytes(ptr, whelper.RefIsObject, -1)
}

func ReturnBuffer(data []byte) whelper.RefId {
	ptr := whelper.KeepaliveObject(data)
	return whelper.NewRefIdWithBytes(ptr, whelper.RefIsBytes, len(data))
}

func ReturnBufferWithFlag(data []byte, flags whelper.RefId) whelper.RefId {
	ptr := whelper.KeepaliveObject(data)
	return whelper.NewRefIdWithBytes(ptr, flags|whelper.RefIsBytes, len(data))
}

func ReturnVoid() whelper.RefId {
	return whelper.RefId(0)
}

func ReturnError(err error) whelper.RefId {
	wrapped := &Error{
		Message: err.Error(),
	}
	encoded, err := wrapped.MarshalMsg(nil)
	if err != nil {
		panic(err)
	}
	return ReturnBufferWithFlag(encoded, whelper.RefIsError)
}
