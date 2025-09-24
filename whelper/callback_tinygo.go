//go:build tinygo
// +build tinygo

package whelper

//export goCallbackJsHandler
func goCallbackJsHandler(refId RefId, argsPtr uint32, argsLen uint32) RefId

type CallbackFunc func(args ...uint64) RefId
