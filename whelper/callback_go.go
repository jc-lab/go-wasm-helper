//go:build !tinygo
// +build !tinygo

package whelper

//go:wasmimport gojs github.com/jc-lab/go-wasm-helper.goCallbackJsHandler
func goCallbackJsHandler(refId RefId, argsPtr uint32, argsLen uint32) RefId

type CallbackFunc func(args ...uint64) RefId
