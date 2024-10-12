package main

import (
	"errors"
	_ "github.com/jc-lab/go-wasm-helper/wmem"
	"github.com/jc-lab/go-wasm-helper/wrefid"
	"github.com/jc-lab/go-wasm-helper/wret"
	"github.com/jc-lab/go-wasm-helper/wruntime"
	"log"
	"time"
)

func main() {
	wruntime.Main()
}

//export sampleData
func sampleData() wrefid.RefId {
	return wret.ReturnBuffer([]byte("HELLO WORLD!!!"))
}

//export sampleError
func sampleError() wrefid.RefId {
	return wret.ReturnError(errors.New("sample error"))
}

//export sampleObject
func sampleObject() wrefid.RefId {
	obj := &SampleObject{}
	return wret.ReturnObject(obj)
}

//export goroutineTestA
func goroutineTestA() wrefid.RefId {
	go func() {
		time.Sleep(time.Second)
		log.Println("HELLO WORLD")
	}()
	return wret.ReturnVoid()
}

type SampleObject struct {
}
