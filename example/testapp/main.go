package main

import (
	"errors"
	_ "github.com/jc-lab/go-wasm-helper/wmem"
	"github.com/jc-lab/go-wasm-helper/wret"
	"github.com/jc-lab/go-wasm-helper/wruntime"
	"log"
	"time"
)

func main() {
	wruntime.Main()
}

//export sampleData
func sampleData() wret.RefId {
	return wret.ReturnBuffer([]byte("HELLO WORLD!!!"))
}

//export sampleError
func sampleError() wret.RefId {
	return wret.ReturnError(errors.New("sample error"))
}

//export sampleObject
func sampleObject() wret.RefId {
	obj := &SampleObject{}
	return wret.ReturnObject(obj)
}

//export goroutineTestA
func goroutineTestA() wret.RefId {
	go func() {
		time.Sleep(time.Second)
		log.Println("HELLO WORLD")
	}()
	return wret.ReturnVoid()
}

type SampleObject struct {
}
