package wret

import (
	"github.com/jc-lab/go-wasm-helper/wmem"
	"github.com/jc-lab/go-wasm-helper/wrefid"
)

func ReturnObject(data interface{}) wrefid.RefId {
	ptr := wmem.KeepaliveObject(data)
	return wrefid.NewRefId(ptr, wrefid.RetIsObject, -1)
}

func ReturnBuffer(data []byte) wrefid.RefId {
	ptr := wmem.KeepaliveObject(data)
	return wrefid.NewRefId(ptr, wrefid.RetIsBytes, len(data))
}

func ReturnBufferWithFlag(data []byte, flags wrefid.RefId) wrefid.RefId {
	ptr := wmem.KeepaliveObject(data)
	return wrefid.NewRefId(ptr, flags|wrefid.RetIsBytes, len(data))
}

func ReturnVoid() wrefid.RefId {
	return wrefid.RefId(0)
}

func ReturnError(err error) wrefid.RefId {
	wrapped := &Error{
		Message: err.Error(),
	}
	encoded, err := wrapped.MarshalMsg(nil)
	if err != nil {
		panic(err)
	}
	return ReturnBufferWithFlag(encoded, wrefid.RetIsError)
}
