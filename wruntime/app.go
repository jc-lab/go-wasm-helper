package wruntime

import "sync"

type appContext struct {
	mutex    sync.Mutex
	shutdown chan struct{}
}

var def appContext

func Main() {
	def.shutdown = make(chan struct{})
	<-def.shutdown
}

//export goShutdown
func Shutdown() {
	def.mutex.Lock()
	defer def.mutex.Unlock()
	if def.shutdown == nil {
		return
	}
	def.shutdown <- struct{}{}
	def.shutdown = nil
}
