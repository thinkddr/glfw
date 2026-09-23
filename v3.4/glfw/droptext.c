// Sytue patch: the C side of the text drop callback (droptext.go).
#include "_cgo_export.h"

void glfwSetDropTextCallbackCB(GLFWwindow *window) {
  glfwSetDropTextCallback(window, (GLFWdroptextfun)goDropTextCB);
}
