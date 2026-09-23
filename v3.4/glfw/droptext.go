package glfw

// Sytue patch: text, and URLs that are not files, dropped on a window.
//
// It is a single callback for every window, present and future, rather than a
// method on Window, because the toolkit that owns the windows (Fyne) keeps its
// *Window private: whoever sets this callback tells the windows apart by their
// native handle.

//#define GLFW_INCLUDE_NONE
//#include "glfw/include/GLFW/glfw3.h"
//void glfwSetDropTextCallbackCB(GLFWwindow *window);
import "C"

import "unsafe"

// DropTextCallback receives text, or URLs that are not files, dropped on w.
// Exactly one of text and urls is set. x and y are in the same units as the
// cursor position callback. File URLs still go to the window's DropCallback.
type DropTextCallback func(w *Window, x, y float64, text string, urls []string)

// fDropText is guarded by windows.l.
var fDropText DropTextCallback

// SetDropTextCallback sets the text drop callback of every window, the ones
// that exist and the ones created afterwards. While it is nil, windows accept
// only paths, exactly as without the patch.
//
// This function may only be called from the main thread.
func SetDropTextCallback(cbfun DropTextCallback) {
	windows.l.Lock()
	fDropText = cbfun
	for data := range windows.m {
		setDropText(data)
	}
	n := len(windows.m)
	windows.l.Unlock()
	if n > 0 {
		panicError()
	}
}

// setDropText applies fDropText to one window; windows.l must be held.
func setDropText(data *C.GLFWwindow) {
	if fDropText == nil {
		C.glfwSetDropTextCallback(data, nil)
	} else {
		C.glfwSetDropTextCallbackCB(data)
	}
}

//export goDropTextCB
func goDropTextCB(window unsafe.Pointer, x, y C.double, text *C.char, count C.int, urls **C.char) {
	w := windows.get((*C.GLFWwindow)(window))
	windows.l.Lock()
	cb := fDropText
	windows.l.Unlock()
	if w == nil || cb == nil {
		return
	}
	var t string
	if text != nil {
		t = C.GoString(text)
	}
	list := make([]string, 0, int(count))
	for _, p := range unsafe.Slice(urls, int(count)) {
		list = append(list, C.GoString(p))
	}
	cb(w, float64(x), float64(y), t, list)
}
