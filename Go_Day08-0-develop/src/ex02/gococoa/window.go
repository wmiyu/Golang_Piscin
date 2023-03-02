package gococoa

// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework Cocoa
//#include "window.h"
import "C"
import "unsafe"

type Window struct {
	ptr unsafe.Pointer
}

func NewWindow(x, y, width, height int, title string) (*Window) {
	w := new(Window)
	w.ptr = C.Window_Create(C.int(x), C.int(y), C.int(width), C.int(height), C.CString(title))

	return w
}

func (self *Window) MakeKeyAndOrderFront() {
	C.Window_MakeKeyAndOrderFront(self.ptr)
}
