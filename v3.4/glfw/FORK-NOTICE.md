# Fork de Sytue de go-gl/glfw (v3.4)

Este repositorio es el módulo `github.com/go-gl/glfw/v3.4/glfw` en la versión que
requiere Fyne v2.8.1 (`v0.1.0-pre.1.0.20260707082822-2a407d02d01a`, GLFW C
`7b6aead9fb88b3623e3b3725ebb42670cbe4c579`), con dos parches: **una ventana puede
aceptar TEXTO y URL que no son ficheros soltados sobre ella**, además de las rutas, y
**una ventana sabe cuándo le están arrastrando ficheros por encima** (para enseñar dónde
soltarlos).

- C: `glfwSetDropTextCallback` y `GLFWdroptextfun` (`glfw/include/GLFW/glfw3.h`), el
  reparto de `text/uri-list` entre ficheros y URL (`glfw/src/input.c`) y los cuatro
  backends que lo entregan: el XDnD de X11, el `wl_data_device` de Wayland, la vista de
  Cocoa y un `IDropTarget` de OLE en Win32 (que entrega también los ficheros, porque una
  ventana registrada para OLE deja de recibir `WM_DROPFILES` del Explorador).
- Go: `SetDropTextCallback` (`droptext.go`), una retrollamada para todas las ventanas.
- El arrastre por encima: `glfwSetDragCallback` y `GLFWdragfun`, avisados al entrar y al
  salir (o al cancelarse o soltarse, y en ese caso antes de la entrega), y solo cuando
  lo arrastrado son ficheros: `XdndEnter` con `text/uri-list` y `XdndLeave`/`XdndDrop` en
  X11; `draggingEntered:` con `NSPasteboardTypeFileURL`, `draggingExited:`,
  `draggingEnded:` y `performDragOperation:` en Cocoa; y `DragEnter` con `CF_HDROP`,
  `DragLeave` y `Drop` del `IDropTarget` en Win32. Cada cambio llega una sola vez. En Go,
  `SetDragCallback` (`drag.go`). En X11 una lista de URI cuenta antes de saber si trae
  ficheros o enlaces; Wayland no avisa; y una ventana de Win32 sin OLE (hilo ya en el
  apartamento multihilo) tampoco.

Una ventana sin retrollamada de texto se comporta EXACTAMENTE como sin el parche. Los
ficheros siguen llegando por la retrollamada de rutas de siempre.

Las fuentes alteradas llevan el comentario «Sytue patch» (la licencia zlib de GLFW pide
marcarlas). El código original sigue siendo de sus autores y bajo sus licencias:
BSD-3-Clause para el enlace de Go (`LICENSE`) y zlib para GLFW
(`glfw/LICENSE.md`); lo añadido por Sytue se distribuye bajo esas mismas licencias.

Este fork no está afiliado ni respaldado por go-gl ni por GLFW. Quien lo usa y por qué:
`apps/fyne/upstream/README.md` del monorepo de Sytue One.
