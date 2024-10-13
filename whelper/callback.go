package whelper

//export goCallbackJsHandler
func goCallbackJsHandler(refId RefId, args []uint64) RefId

type CallbackFunc func(args ...uint64) RefId
