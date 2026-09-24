package glfw

// Sytue patch: whether paths are being dragged over the window, so that the
// application can show where to drop them.
//
// Like the text drop callback (droptext.go), it is a single callback for every
// window, because the toolkit that owns the windows keeps its *Window private.

//#define GLFW_INCLUDE_NONE
//#include "glfw/include/GLFW/glfw3.h"
//void glfwSetDragCallbackCB(GLFWwindow *window);
import "C"

import "unsafe"

// DragCallback is told when a drag of paths starts hovering over w, and when
// it stops: it left, it was cancelled or it was dropped. When it is dropped,
// this comes before the drop callbacks. Each change is reported once. On X11
// a URI list counts before it is known whether it holds files; on Wayland
// nothing is reported.
type DragCallback func(w *Window, hovering bool)

// fDrag is guarded by windows.l.
var fDrag DragCallback

// SetDragCallback sets the drag hover callback of every window, the ones that
// exist and the ones created afterwards. It does not change what the windows
// accept.
//
// This function may only be called from the main thread.
func SetDragCallback(cbfun DragCallback) {
	windows.l.Lock()
	fDrag = cbfun
	for data := range windows.m {
		setDrag(data)
	}
	n := len(windows.m)
	windows.l.Unlock()
	if n > 0 {
		panicError()
	}
}

// setDrag applies fDrag to one window; windows.l must be held.
func setDrag(data *C.GLFWwindow) {
	if fDrag == nil {
		C.glfwSetDragCallback(data, nil)
	} else {
		C.glfwSetDragCallbackCB(data)
	}
}

//export goDragCB
func goDragCB(window unsafe.Pointer, hovering C.int) {
	w := windows.get((*C.GLFWwindow)(window))
	windows.l.Lock()
	cb := fDrag
	windows.l.Unlock()
	if w == nil || cb == nil {
		return
	}
	cb(w, hovering != 0)
}
