package gococoa

// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework Cocoa
//# include "application.h"
import "C"
import "unsafe"

type Application struct {
	ptr unsafe.Pointer
}

func InitApplication() {
	C.InitApplication()
}

func RunApplication() {
	C.RunApplication()
}

