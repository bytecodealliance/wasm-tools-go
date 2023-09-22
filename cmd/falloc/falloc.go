package main

import (
	"runtime"
	"time"
	"unsafe"

	"github.com/k0kubun/pp/v3"
	"github.com/ydnar/wasm-tools-go/wit"
)

func main() {
	p := leak()
	runtime.GC()
	time.Sleep(time.Second)
	runtime.GC()
	res := (*wit.Resolve)(unsafe.Pointer(uintptr(p)))
	pp.Println(res)
}

func leak() uint {
	res := alloc[wit.Resolve]()
	res.Worlds = make([]*wit.World, 3)
	pp.Println(res)
	return uint(uintptr(unsafe.Pointer(res)))
}

func alloc[T any]() *T {
	var v T
	b := make([]byte, unsafe.Sizeof(v))
	p := unsafe.Pointer(unsafe.SliceData(b))
	return (*T)(p)
}
