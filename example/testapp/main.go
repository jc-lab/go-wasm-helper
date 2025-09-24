package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jc-lab/go-wasm-helper/whelper"
	_ "github.com/jc-lab/go-wasm-helper/whelper"
	"github.com/jc-lab/go-wasm-helper/wret"
	"github.com/jc-lab/go-wasm-helper/wruntime"
)

func main() {
	wruntime.Main()
}

//export sampleData
//go:wasmexport sampleData
func sampleData() whelper.RefId {
	return wret.ReturnBuffer([]byte("HELLO WORLD!!!"))
}

//export sampleError
//go:wasmexport sampleError
func sampleError() whelper.RefId {
	return wret.ReturnError(errors.New("sample error"))
}

//export sampleObject
//go:wasmexport sampleObject
func sampleObject() whelper.RefId {
	obj := &SampleObject{}
	return wret.ReturnObject(obj)
}

//export goroutineTestA
//go:wasmexport goroutineTestA
func goroutineTestA() whelper.RefId {
	go func() {
		time.Sleep(time.Second)
		log.Println("HELLO WORLD")
	}()
	return wret.ReturnVoid()
}

//export callbackTest
//go:wasmexport callbackTest
func callbackTest(callbackRef whelper.RefId, a uint32) whelper.RefId {
	if a != 0x12 {
		panic(fmt.Errorf("param a is not 0x12"))
	}
	fn := callbackRef.ToFunction()
	result := fn(uint64(0x1000+a), 0x11111111, 0x22222222)
	return whelper.Uint32(result.GetNumber() + 0x2000)
}

//export bufferReadTest
//go:wasmexport bufferReadTest
func bufferReadTest(param whelper.RefId) whelper.RefId {
	data := param.GetBuffer()
	log.Printf("INPUT DATA: %s", string(data))
	return wret.ReturnBuffer(data)
}

type SampleObject struct {
}
