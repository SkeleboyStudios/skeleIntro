//+build js

package main

import "syscall/js"

func navigateToPageImpl(url string) {
	window := js.Global().Get("window")
	window.Call("open", url, "_blank")
}
