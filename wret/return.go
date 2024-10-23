package wret

import (
	"github.com/jc-lab/go-wasm-helper/internal/fmtstate"
	"github.com/jc-lab/go-wasm-helper/whelper"
	"github.com/pkg/errors"
	"strings"
)

type stackTraceable interface {
	StackTrace() errors.StackTrace
}

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

	stackErr, ok := err.(stackTraceable)
	if ok {
		fmtBuf := fmtstate.NewCustomState(0, 0, "+")
		stackErr.StackTrace().Format(fmtBuf, 'v')
		wrapped.Stack = strings.Split(fmtBuf.String(), "\n")
	}

	encoded, err := wrapped.MarshalMsg(nil)
	if err != nil {
		panic(err)
	}
	return ReturnBufferWithFlag(encoded, whelper.RefIsError)
}
